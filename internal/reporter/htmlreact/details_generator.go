package htmlreact

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/filereader"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/utils"
	"golang.org/x/net/html"
)

// -----------------------------------------------------------------------------
// Orchestration
// -----------------------------------------------------------------------------

func generateDetailsPages(b *HtmlReactReportBuilder, tree *model.SummaryTree) error {
	fileMap := make(map[string]*model.FileNode)
	collectFiles(tree.Root, fileMap)

	detailsHTMLContent, err := readEmbeddedDetailsHTML()
	if err != nil {
		return fmt.Errorf("failed to read embedded details.html: %w", err)
	}

	for _, fileNode := range fileMap {
		detailsData, err := b.transformFileNodeToDetails(tree, fileNode)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not transform details for '%s': %v\n", fileNode.Path, err)
			continue
		}

		if err := writeDetailsPage(b.outputDir, fileNode.Path, detailsHTMLContent, detailsData); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not generate details page for '%s': %v\n", fileNode.Path, err)
		}
	}

	return nil
}

// transformFileNodeToDetails executes a clean pipeline to assemble the view model.
// Cyclomatic Complexity: 1
func (b *HtmlReactReportBuilder) transformFileNodeToDetails(tree *model.SummaryTree, fileNode *model.FileNode) (*detailsV1, error) {
	// 1. I/O Phase
	sourceLines := b.readSourceLines(fileNode)

	// 2. Data Filtering Phase
	reportsList, globalToLocalMap := b.buildReportFilter(tree, fileNode)

	// 3. Line & Method Mapping Phase
	detailsLines := b.buildLineDetails(fileNode, sourceLines, reportsList, globalToLocalMap, len(tree.ReportNames))
	detailsMethods, totalBranches, coveredBranches, maxCyclo := b.buildMethodDetails(fileNode)

	// 4. Totals Calculation Phase
	totalsData := b.buildFileTotals(fileNode, totalBranches, coveredBranches, maxCyclo)

	return &detailsV1{
		SchemaVersion:     1,
		GeneratedAt:       time.Now().UTC().Format(time.RFC3339),
		Title:             strings.Join(tree.ParserNames, " | "),
		FileName:          fileNode.Path,
		Metadata:          []metadataItem{},
		Totals:            totalsData,
		MetricDefinitions: b.buildMetricDefinitions(),
		Methods:           detailsMethods,
		Lines:             detailsLines,
		Reports:           reportsList,
	}, nil
}

// -----------------------------------------------------------------------------
// Pipeline Stage 1: I/O & Source Resolution
// -----------------------------------------------------------------------------

func (b *HtmlReactReportBuilder) readSourceLines(fileNode *model.FileNode) []string {
	if fileNode.SourceDir == "" {
		return generateEmptyLines(fileNode)
	}

	reader := filereader.NewDefaultReader()
	absPath, err := utils.FindFileInSourceDirs(fileNode.Path, []string{fileNode.SourceDir}, reader, b.logger)
	if err != nil {
		b.logger.Warn("Could not find source file for details page", "file", fileNode.Path, "error", err)
		return generateEmptyLines(fileNode)
	}

	lines, err := reader.ReadFile(absPath)
	if err != nil {
		b.logger.Warn("Could not read source file for details page", "file", absPath, "error", err)
		return generateEmptyLines(fileNode)
	}

	return lines
}

func generateEmptyLines(fileNode *model.FileNode) []string {
	maxLine := 0
	for ln := range fileNode.Lines {
		if ln > maxLine {
			maxLine = ln
		}
	}
	return make([]string, maxLine)
}

// -----------------------------------------------------------------------------
// Pipeline Stage 2: Report Filtering
// -----------------------------------------------------------------------------

func (b *HtmlReactReportBuilder) buildReportFilter(tree *model.SummaryTree, fileNode *model.FileNode) ([]report, map[int]int) {
	isReportRelevant := make(map[int]bool)
	for _, line := range fileNode.Lines {
		for reportIdx, hits := range line.ReportHits {
			if hits > 0 {
				isReportRelevant[reportIdx] = true
			}
		}
	}

	globalToLocalMap := make(map[int]int)
	var reportsList []report

	for globalIdx, name := range tree.ReportNames {
		if isReportRelevant[globalIdx] {
			globalToLocalMap[globalIdx] = len(reportsList)
			reportsList = append(reportsList, report{
				Name: fmt.Sprintf("Report %d: %s", globalIdx+1, filepath.Base(name)),
				Path: name,
			})
		}
	}

	// Fallback: If no report has covered lines, show all reports to avoid empty UI
	if len(reportsList) == 0 {
		for i, name := range tree.ReportNames {
			reportsList = append(reportsList, report{
				Name: fmt.Sprintf("Report %d: %s", i+1, filepath.Base(name)),
				Path: name,
			})
		}
		return reportsList, nil
	}

	return reportsList, globalToLocalMap
}

// -----------------------------------------------------------------------------
// Pipeline Stage 3: Data Mapping (Lines & Methods)
// -----------------------------------------------------------------------------

func (b *HtmlReactReportBuilder) buildLineDetails(fileNode *model.FileNode, sourceLines []string, reportsList []report, globalToLocalMap map[int]int, totalReports int) []lineDetail {
	detailsLines := make([]lineDetail, len(sourceLines))

	for i, content := range sourceLines {
		lineNumber := i + 1
		ld := lineDetail{
			LineNumber: lineNumber,
			Content:    content,
			Status:     StatusNotCoverable,
		}

		if fileNode.Diff != nil {
			if fileNode.Diff.AddedLines[lineNumber] {
				ld.DiffStatus = "added"
			} else if fileNode.Diff.ModifiedLines[lineNumber] {
				ld.DiffStatus = "modified"
			}
		}

		if lm, ok := fileNode.Lines[lineNumber]; ok && lm.Hits >= 0 {
			ld.Status = StatusUncovered
			if lm.Hits > 0 {
				ld.Status = StatusCovered
			}

			ld.Hits = mapHitsToLocal(lm.ReportHits, reportsList, globalToLocalMap, totalReports)

			if lm.TotalBranches > 0 {
				ld.BranchInfo = &branchInfo{Covered: lm.CoveredBranches, Total: lm.TotalBranches}
				if lm.CoveredBranches > 0 && lm.CoveredBranches < lm.TotalBranches {
					ld.Status = StatusPartial
				}
			}
		}

		detailsLines[i] = ld
	}

	return detailsLines
}

func mapHitsToLocal(reportHits []int, reportsList []report, globalToLocal map[int]int, totalReports int) []int {
	if len(reportHits) == 0 {
		return nil
	}
	if globalToLocal != nil {
		localHits := make([]int, len(reportsList))
		for globalIdx, hitCount := range reportHits {
			if localIdx, ok := globalToLocal[globalIdx]; ok {
				localHits[localIdx] = hitCount
			}
		}
		return localHits
	}

	n := min(totalReports, len(reportHits))
	dense := make([]int, n)
	copy(dense, reportHits[:n])
	return dense
}

func (b *HtmlReactReportBuilder) buildMethodDetails(fileNode *model.FileNode) ([]methodDetail, int, int, int) {
	var activeProviders []MethodMetricProvider
	for _, key := range b.config.MethodMetrics {
		if p, ok := MethodProviderRegistry[key]; ok && b.config.ActiveMethodMetrics[key] {
			activeProviders = append(activeProviders, p)
		}
	}

	var detailsMethods []methodDetail
	var totalMethodBranches, coveredMethodBranches, maxCyclo int

	for _, method := range fileNode.Methods {
		md := methodDetail{
			Name:       method.Name,
			StartLine:  method.StartLine,
			EndLine:    method.EndLine,
			DiffStatus: method.DiffStatus,
			Metrics:    make(map[string]methodMetric),
		}

		for _, provider := range activeProviders {
			provider.Apply(&method, &md)
		}

		if method.BranchesValid > 0 {
			totalMethodBranches += method.BranchesValid
			coveredMethodBranches += method.BranchesCovered
		}
		if method.CyclomaticComplexity != nil && *method.CyclomaticComplexity > maxCyclo {
			maxCyclo = *method.CyclomaticComplexity
		}

		detailsMethods = append(detailsMethods, md)
	}

	sort.Slice(detailsMethods, func(i, j int) bool {
		return detailsMethods[i].StartLine < detailsMethods[j].StartLine
	})

	return detailsMethods, totalMethodBranches, coveredMethodBranches, maxCyclo
}

// -----------------------------------------------------------------------------
// Pipeline Stage 4: Totals Assembly
// -----------------------------------------------------------------------------

func (b *HtmlReactReportBuilder) buildFileTotals(fileNode *model.FileNode, totalBranches, coveredBranches, maxCyclo int) totals {
	fileMetrics := b.buildMetricsMap(fileNode.Metrics)

	t := totals{
		Files:    1,
		Folders:  0,
		Statuses: b.convertStatuses(fileNode.Statuses),
	}

	// Dynamic assignment helpers heavily reduce Cyclomatic Complexity here
	assignLineMetric(&t.StatementCoverage, fileMetrics, string(config.StatementCoverage))
	assignLineMetric(&t.LineCoverage, fileMetrics, string(config.LineCoverage))
	assignBranchMetric(&t.BranchCoverage, fileMetrics, string(config.BranchCoverage))
	assignMethodHitMetric(&t.MethodsHit, fileMetrics, string(config.MethodsHit))
	assignMethodFullMetric(&t.MethodsFullyCovered, fileMetrics, string(config.MethodsFullyCovered))
	assignLineMetric(&t.PatchStatementCoverage, fileMetrics, string(config.PatchStatementCoverage))
	assignLineMetric(&t.PatchLineCoverage, fileMetrics, string(config.PatchLineCoverage))
	assignMethodHitMetric(&t.PatchMethodsHit, fileMetrics, string(config.PatchMethodsHit))

	// Overrides for specific edge cases
	if b.config.ActiveFileMetrics[config.PatchStatementCoverage] && fileNode.Diff != nil && t.PatchStatementCoverage == nil && fileNode.Metrics.StatementsValid > 0 {
		t.PatchStatementCoverage = &lineCoverageDetail{Percentage: 100.0} // Fallback to safe when modified but no statements changed
	}

	if totalBranches > 0 && b.config.ActiveFileMetrics[config.BranchCoverage] {
		t.MethodBranchCoverage = &branchCoverageDetail{
			Covered:    coveredBranches,
			Total:      totalBranches,
			Percentage: utils.CalculatePercentage(coveredBranches, totalBranches, 2),
		}
	}

	if maxCyclo > 0 && b.config.ActiveFileMetrics[config.MaxCyclomaticComplexity] {
		t.MaxCyclomaticComplexity = &lineCoverageDetail{Total: maxCyclo}
	}

	return t
}

// --- Reflection-free Dynamic Assignment Helpers ---

func assignLineMetric(target **lineCoverageDetail, fm metricsMap, key string) {
	if val, ok := fm[key].(lineCoverageDetail); ok {
		*target = &val
	}
}

func assignBranchMetric(target **branchCoverageDetail, fm metricsMap, key string) {
	if val, ok := fm[key].(branchCoverageDetail); ok {
		*target = &val
	}
}

func assignMethodHitMetric(target **methodsHitDetail, fm metricsMap, key string) {
	if val, ok := fm[key].(methodsHitDetail); ok {
		*target = &val
	}
}

func assignMethodFullMetric(target **methodsFullyCoveredDetail, fm metricsMap, key string) {
	if val, ok := fm[key].(methodsFullyCoveredDetail); ok {
		*target = &val
	}
}

// -----------------------------------------------------------------------------
// HTML & File Writing Utilities
// -----------------------------------------------------------------------------

func writeDetailsPage(outDir, logicalPath string, template []byte, detailsData *detailsV1) error {
	var jsonBuf bytes.Buffer
	enc := json.NewEncoder(&jsonBuf)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(detailsData); err != nil {
		return fmt.Errorf("failed to marshal details data to JSON: %w", err)
	}
	scriptContent := "window.__NANOVISION_DETAILS__ = " + jsonBuf.String()

	modifiedHTML, err := injectDataIntoHTML(template, scriptContent)
	if err != nil {
		return err
	}

	detailsFileName := strings.ReplaceAll(logicalPath, "/", "_") + ".html"
	return os.WriteFile(filepath.Join(outDir, detailsFileName), []byte(modifiedHTML), 0o644)
}

func readEmbeddedDetailsHTML() ([]byte, error) {
	distFS, err := getReactDist()
	if err != nil {
		return nil, fmt.Errorf("failed to get embedded dist FS: %w", err)
	}

	file, err := distFS.Open("details.html")
	if err != nil {
		return nil, fmt.Errorf("failed to open embedded details.html: %w", err)
	}
	defer file.Close()

	return io.ReadAll(file)
}

func injectDataIntoHTML(htmlContent []byte, scriptContent string) (string, error) {
	doc, err := html.Parse(bytes.NewReader(htmlContent))
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %w", err)
	}

	var found bool
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "script" {
			if n.FirstChild != nil && strings.Contains(n.FirstChild.Data, "window.__NANOVISION_DETAILS__") {
				n.FirstChild.Data = scriptContent
				found = true
				return
			}
		}
		if found {
			return
		}
		for c := n.FirstChild; c != nil && !found; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(doc)

	if !found {
		return "", fmt.Errorf("placeholder script 'window.__NANOVISION_DETAILS__' not found")
	}

	var buf bytes.Buffer
	if err := html.Render(&buf, doc); err != nil {
		return "", fmt.Errorf("failed to render modified HTML: %w", err)
	}

	return buf.String(), nil
}

func collectFiles(dir *model.DirNode, fileMap map[string]*model.FileNode) {
	for _, file := range dir.Files {
		fileMap[file.Path] = file
	}
	for _, subDir := range dir.Subdirs {
		collectFiles(subDir, fileMap)
	}
}
