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

	"github.com/IgorBayerl/AdlerCov/filereader"
	"github.com/IgorBayerl/AdlerCov/internal/model"
	"github.com/IgorBayerl/AdlerCov/internal/utils"
	"golang.org/x/net/html"
)

// generateDetailsPages iterates through all file nodes and creates a separate HTML page for each.
func generateDetailsPages(b *HtmlReactReportBuilder, tree *model.SummaryTree) error {
	fileNodeMap := make(map[string]*model.FileNode)
	collectFiles(tree.Root, fileNodeMap)

	detailsHTMLContent, err := readEmbeddedDetailsHTML()
	if err != nil {
		return err
	}

	for _, fileNode := range fileNodeMap {
		if err := b.createDetailPage(fileNode, detailsHTMLContent, tree); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not generate details page for '%s': %v\n", fileNode.Path, err)
		}
	}

	return nil
}

// createDetailPage generates a single HTML file with coverage details for a given file node.
func (b *HtmlReactReportBuilder) createDetailPage(fileNode *model.FileNode, detailsHTMLContent []byte, tree *model.SummaryTree) error {
	// Transform the file node data into the format required by the UI.
	detailsData, err := b.transformFileNodeToDetails(fileNode, tree)
	if err != nil {
		return fmt.Errorf("failed to transform file node data: %w", err)
	}

	var jsonBuf bytes.Buffer
	enc := json.NewEncoder(&jsonBuf)
	enc.SetEscapeHTML(false) // Prevent characters like '<' from being escaped
	if err := enc.Encode(detailsData); err != nil {
		return fmt.Errorf("failed to marshal details data to JSON: %w", err)
	}
	scriptContent := "window.__ADLERCOV_DETAILS__ = " + jsonBuf.String()

	modifiedHTML, err := injectDataIntoHTML(detailsHTMLContent, scriptContent)
	if err != nil {
		return err
	}

	detailsFileName := strings.ReplaceAll(fileNode.Path, "/", "_") + ".html"
	detailsFilePath := filepath.Join(b.outputDir, detailsFileName)

	return os.WriteFile(detailsFilePath, []byte(modifiedHTML), 0644)
}

// transformFileNodeToDetails converts a model.FileNode into the rich detailsV1 structure.
func (b *HtmlReactReportBuilder) transformFileNodeToDetails(fileNode *model.FileNode, tree *model.SummaryTree) (*detailsV1, error) {
	reader := filereader.NewDefaultReader()
	absPath, err := utils.FindFileInSourceDirs(fileNode.Path, []string{fileNode.SourceDir}, reader, b.logger)
	var sourceLines []string
	if err == nil {
		sourceLines, err = reader.ReadFile(absPath)
		if err != nil {
			b.logger.Warn("Could not read source file for details page", "file", absPath, "error", err)
		}
	} else {
		b.logger.Warn("Could not find source file for details page", "file", fileNode.Path, "error", err)
	}

	detailsLines := make([]lineDetail, len(sourceLines))
	for i, lineContent := range sourceLines {
		lineNumber := i + 1
		lineMetric, hasMetric := fileNode.Lines[lineNumber]
		ld := lineDetail{LineNumber: lineNumber, Content: lineContent, Status: StatusNotCoverable}

		if hasMetric && lineMetric.Hits >= 0 {
			hits := lineMetric.Hits
			ld.Hits = &hits
			ld.Status = StatusUncovered
			if lineMetric.Hits > 0 {
				ld.Status = StatusCovered
			}

			if lineMetric.TotalBranches > 0 {
				ld.BranchInfo = &branchInfo{Covered: lineMetric.CoveredBranches, Total: lineMetric.TotalBranches}
				if lineMetric.CoveredBranches > 0 && lineMetric.CoveredBranches < lineMetric.TotalBranches {
					ld.Status = StatusPartial
				}
			}
		}
		detailsLines[i] = ld
	}

	var detailsMethods []methodDetail
	var maxCyclo int = 0
	var totalMethodBranches, coveredMethodBranches int = 0, 0

	for _, method := range fileNode.Methods {
		lineCovPct := utils.CalculatePercentage(method.LinesCovered, method.LinesValid, 0)
		branchCovPct := utils.CalculatePercentage(method.BranchesCovered, method.BranchesValid, 0)
		md := methodDetail{
			Name:      method.Name,
			StartLine: method.StartLine,
			EndLine:   method.EndLine,
			Metrics:   make(map[string]methodMetric),
		}

		lineMetric := methodMetric{Value: utils.FormatPercentage(lineCovPct, 0)}
		lineRisk := getRiskStatus(lineCovPct)
		if lineRisk == RiskDanger || lineRisk == RiskWarning {
			lineMetric.Status = lineRisk
		}
		md.Metrics["lineCoverage"] = lineMetric

		if method.BranchesValid > 0 {
			branchMetric := methodMetric{Value: utils.FormatPercentage(branchCovPct, 0)}
			branchRisk := getRiskStatus(branchCovPct)
			if branchRisk == RiskDanger || branchRisk == RiskWarning {
				branchMetric.Status = branchRisk
			}
			md.Metrics["branchCoverage"] = branchMetric
		}

		if method.CyclomaticComplexity != nil {
			md.Metrics["cyclomaticComplexity"] = methodMetric{Value: fmt.Sprintf("%d", *method.CyclomaticComplexity)}
		}
		detailsMethods = append(detailsMethods, md)

		// Aggregate values for the top card
		totalMethodBranches += method.BranchesValid
		coveredMethodBranches += method.BranchesCovered
		if method.CyclomaticComplexity != nil && *method.CyclomaticComplexity > maxCyclo {
			maxCyclo = *method.CyclomaticComplexity
		}
	}
	sort.Slice(detailsMethods, func(i, j int) bool { return detailsMethods[i].StartLine < detailsMethods[j].StartLine })

	fileMetrics, fileStatuses := b.buildMetricsMap(fileNode.Metrics)
	totalsData := totals{Files: 1, Folders: 0, Statuses: fileStatuses}
	if lc, ok := fileMetrics["lineCoverage"].(lineCoverageDetail); ok {
		totalsData.LineCoverage = &lc
	}
	if bc, ok := fileMetrics["branchCoverage"].(branchCoverageDetail); ok {
		totalsData.BranchCoverage = &bc
	}
	if mc, ok := fileMetrics["methodsCovered"].(methodsCoveredDetail); ok {
		totalsData.MethodsCovered = &mc
	}
	if mfc, ok := fileMetrics["methodsFullyCovered"].(methodsFullyCoveredDetail); ok {
		totalsData.MethodsFullyCovered = &mfc
	}

	// Populate new aggregated method metrics
	if totalMethodBranches > 0 {
		methodBranchPct := utils.CalculatePercentage(coveredMethodBranches, totalMethodBranches, 2)
		totalsData.MethodBranchCoverage = &branchCoverageDetail{
			Covered:    coveredMethodBranches,
			Total:      totalMethodBranches,
			Percentage: methodBranchPct,
		}
		fileStatuses["methodBranchCoverage"] = getRiskStatus(methodBranchPct)
	}

	if maxCyclo > 0 {
		// CORRECTED: Use lineCoverageDetail and map the value to the 'Total' field to satisfy the Zod schema.
		totalsData.MaxCyclomaticComplexity = &lineCoverageDetail{
			Covered:    0,
			Total:      maxCyclo,
			Percentage: 0,
		}
	}

	return &detailsV1{
		SchemaVersion:     1,
		GeneratedAt:       time.Now().UTC().Format(time.RFC3339),
		Title:             tree.ParserName,
		FileName:          fileNode.Path,
		Metadata:          []metadataItem{},
		Totals:            totalsData,
		MetricDefinitions: b.buildMetricDefinitions(),
		Methods:           detailsMethods,
		Lines:             detailsLines,
	}, nil
}

// readEmbeddedDetailsHTML reads the content of the details.html file from the embedded file system.
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

// injectDataIntoHTML parses the HTML, finds the placeholder script, and replaces its content.
func injectDataIntoHTML(htmlContent []byte, scriptContent string) (string, error) {
	doc, err := html.Parse(bytes.NewReader(htmlContent))
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %w", err)
	}

	var found bool
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "script" {
			if n.FirstChild != nil && strings.Contains(n.FirstChild.Data, "window.__ADLERCOV_DETAILS__") {
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
		return "", fmt.Errorf("placeholder script 'window.__ADLERCOV_DETAILS__' not found in details.html template")
	}

	var buf bytes.Buffer
	if err := html.Render(&buf, doc); err != nil {
		return "", fmt.Errorf("failed to render modified HTML: %w", err)
	}

	return buf.String(), nil
}

// collectFiles is a helper function to recursively gather all file nodes from the directory tree.
func collectFiles(dir *model.DirNode, fileMap map[string]*model.FileNode) {
	for _, file := range dir.Files {
		fileMap[file.Path] = file
	}
	for _, subDir := range dir.Subdirs {
		collectFiles(subDir, fileMap)
	}
}
