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

// generateDetailsPages iterates through all file nodes and creates a separate HTML page for each.
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

// writeDetailsPage renders the details HTML file with embedded JSON data.
func writeDetailsPage(outDir, logicalPath string, template []byte, detailsData *detailsV1) error {
	var jsonBuf bytes.Buffer
	enc := json.NewEncoder(&jsonBuf)
	enc.SetEscapeHTML(false) // Prevent characters like '<' from being escaped.
	if err := enc.Encode(detailsData); err != nil {
		return fmt.Errorf("failed to marshal details data to JSON: %w", err)
	}
	scriptContent := "window.__NANOVISION_DETAILS__ = " + jsonBuf.String()

	modifiedHTML, err := injectDataIntoHTML(template, scriptContent)
	if err != nil {
		return err
	}

	detailsFileName := strings.ReplaceAll(logicalPath, "/", "_") + ".html"
	detailsFilePath := filepath.Join(outDir, detailsFileName)

	return os.WriteFile(detailsFilePath, []byte(modifiedHTML), 0o644)
}

// transformFileNodeToDetails converts a model.FileNode into the rich detailsV1 structure.
func (b *HtmlReactReportBuilder) transformFileNodeToDetails(tree *model.SummaryTree, fileNode *model.FileNode) (*detailsV1, error) {
	reader := filereader.NewDefaultReader()

	var sourceLines []string
	if fileNode.SourceDir != "" {
		if absPath, err := utils.FindFileInSourceDirs(fileNode.Path, []string{fileNode.SourceDir}, reader, b.logger); err == nil {
			if lines, err := reader.ReadFile(absPath); err == nil {
				sourceLines = lines
			} else {
				b.logger.Warn("Could not read source file for details page", "file", absPath, "error", err)
			}
		} else {
			b.logger.Warn("Could not find source file for details page", "file", fileNode.Path, "error", err)
		}
	}

	if len(sourceLines) == 0 {
		maxLine := 0
		for ln := range fileNode.Lines {
			if ln > maxLine {
				maxLine = ln
			}
		}
		sourceLines = make([]string, maxLine)
	}

	// 1. Identify which reports have coverage data for this specific file.
	isReportRelevant := make([]bool, len(tree.ReportNames))
	for _, line := range fileNode.Lines {
		for reportIdx, hits := range line.ReportHits {
			// A report is relevant if it results in at least one covered line.
			if hits > 0 {
				isReportRelevant[reportIdx] = true
			}
		}
	}

	// 2. Build the filtered list of reports and create a mapping from the
	// old global report index to the new local index for this page.
	globalToLocalIndexMap := make(map[int]int)
	reportsList := make([]report, 0)
	for globalIdx, name := range tree.ReportNames {
		if isReportRelevant[globalIdx] {
			localIdx := len(reportsList)
			globalToLocalIndexMap[globalIdx] = localIdx
			reportsList = append(reportsList, report{
				// Use globalIdx+1 for consistent "Report X" numbering across the site.
				Name: fmt.Sprintf("Report %d: %s", globalIdx+1, filepath.Base(name)),
				Path: name,
			})
		}
	}

	// Fallback: If no report has covered lines (e.g., file is fully uncovered),
	// show all reports to avoid a confusing empty list.
	if len(reportsList) == 0 {
		globalToLocalIndexMap = nil // Signal to use the global list
		for i, name := range tree.ReportNames {
			reportsList = append(reportsList, report{
				Name: fmt.Sprintf("Report %d: %s", i+1, filepath.Base(name)),
				Path: name,
			})
		}
	}

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

			// 3. Build the dense `Hits` array for the UI, mapping global hits to the local, filtered list.
			if len(lm.ReportHits) > 0 {
				if globalToLocalIndexMap != nil {
					// Create a dense hits array that matches the filtered reportsList.
					localHits := make([]int, len(reportsList))
					for globalIdx, hitCount := range lm.ReportHits {
						if localIdx, ok := globalToLocalIndexMap[globalIdx]; ok {
							localHits[localIdx] = hitCount
						}
					}
					ld.Hits = localHits
				} else {
					// Fallback case: use the full, original list.
					n := len(tree.ReportNames)
					if len(lm.ReportHits) < n {
						n = len(lm.ReportHits)
					}
					dense := make([]int, n)
					copy(dense, lm.ReportHits[:n])
					ld.Hits = dense
				}
			}

			if lm.TotalBranches > 0 {
				ld.BranchInfo = &branchInfo{
					Covered: lm.CoveredBranches,
					Total:   lm.TotalBranches,
				}
				if lm.CoveredBranches > 0 && lm.CoveredBranches < lm.TotalBranches {
					ld.Status = StatusPartial
				}
			}
		}

		detailsLines[i] = ld
	}

	detailsMethods := make([]methodDetail, 0, len(fileNode.Methods))
	var totalMethodBranches, coveredMethodBranches int
	maxCyclo := 0

	for _, method := range fileNode.Methods {
		md := methodDetail{
			Name:      method.Name,
			StartLine: method.StartLine,
			EndLine:   method.EndLine,
			Metrics:   make(map[string]methodMetric),
		}

		if fileNode.Diff != nil {
			newLinesTotal := 0
			newLinesCovered := 0

			if fileNode.Diff.Kind == model.ChangeKindAdded {
				md.DiffStatus = "added"
				newLinesTotal = method.LinesValid
				newLinesCovered = method.LinesCovered
			} else {
				isModified := false
				for ln := method.StartLine; ln <= method.EndLine; ln++ {
					if fileNode.Diff.AddedLines[ln] || fileNode.Diff.ModifiedLines[ln] {
						isModified = true
						if lm, ok := fileNode.Lines[ln]; ok && lm.Hits >= 0 {
							newLinesTotal++
							if lm.Hits > 0 {
								newLinesCovered++
							}
						}
					}
				}

				if isModified {
					if newLinesTotal > 0 && newLinesTotal == method.LinesValid {
						md.DiffStatus = "added"
					} else {
						md.DiffStatus = "modified"
					}
				}
			}

			// 1. Guard the Lines assignment with ActiveMetrics
			if newLinesTotal > 0 {
				if b.config.ActiveMetrics[config.PatchLineCoverage] {
					md.NewLinesCoverage = &newLinesCoverage{
						Covered: newLinesCovered,
						Total:   newLinesTotal,
					}
					// Show patch lines as a fraction
					md.Metrics[MethodUIPatchLineCoverage] = methodMetric{
						Value: fmt.Sprintf("%d / %d", newLinesCovered, newLinesTotal),
					}
				}
			}

			// Compute Patch Statements for this method
			patchStmtsTotal := 0
			patchStmtsCovered := 0

			// 2. Properly handle newly added files for Patch Statements
			if fileNode.Diff.Kind == model.ChangeKindAdded {
				patchStmtsTotal = method.StatementsValid
				patchStmtsCovered = method.StatementsCovered
			} else {
				for _, stmt := range fileNode.Statements {
					if stmt.StartLine >= method.StartLine && stmt.EndLine <= method.EndLine {
						inPatch := false
						for i := stmt.StartLine; i <= stmt.EndLine; i++ {
							if fileNode.Diff.AddedLines[i] || fileNode.Diff.ModifiedLines[i] {
								inPatch = true
								break
							}
						}
						if inPatch {
							patchStmtsTotal++
							stmtCovered := false
							for i := stmt.StartLine; i <= stmt.EndLine; i++ {
								if lm, ok := fileNode.Lines[i]; ok && lm.Hits > 0 {
									if fileNode.Diff.AddedLines[i] || fileNode.Diff.ModifiedLines[i] {
										stmtCovered = true
										break
									}
								}
							}
							if stmtCovered {
								patchStmtsCovered++
							}
						}
					}
				}
			}

			// 3. Fallback to 0/0 if lines changed but no statements were caught (e.g. function signature changed)
			// and Guard the Statements assignment with ActiveMetrics
			if patchStmtsTotal > 0 || newLinesTotal > 0 {
				if b.config.ActiveMetrics[config.PatchStatementCoverage] {
					cov := &newLinesCoverage{
						Covered: patchStmtsCovered,
						Total:   patchStmtsTotal,
					}
					md.NewStatementsCoverage = cov
					md.NewStatementCoverage = cov
					md.Metrics[MethodUIPatchStmtCoverage] = methodMetric{
						Value: fmt.Sprintf("%d / %d", patchStmtsCovered, patchStmtsTotal),
					}
				}
			}
		}

		if b.config.ActiveMetrics[config.StatementCoverage] {
			md.Metrics[MethodUIStmtCoverage] = methodMetric{
				Value: fmt.Sprintf("%d / %d", method.StatementsCovered, method.StatementsValid),
			}
		}

		if b.config.ActiveMetrics[config.LineCoverage] {
			md.Metrics[MethodUILineCoverage] = methodMetric{
				Value: fmt.Sprintf("%d / %d", method.LinesCovered, method.LinesValid),
			}
		}

		if method.BranchesValid > 0 {
			if b.config.ActiveMetrics[config.BranchCoverage] {
				md.Metrics[MethodUIBranchCoverage] = methodMetric{
					Value: fmt.Sprintf("%d / %d", method.BranchesCovered, method.BranchesValid),
				}
			}
			totalMethodBranches += method.BranchesValid
			coveredMethodBranches += method.BranchesCovered
		}

		if method.CyclomaticComplexity != nil {
			if *method.CyclomaticComplexity > maxCyclo {
				maxCyclo = *method.CyclomaticComplexity
			}
			if b.config.ActiveMetrics[config.MaxCyclomaticComplexity] {
				md.Metrics[MethodUICyclomaticComplexity] = methodMetric{
					Value: fmt.Sprintf("%d", *method.CyclomaticComplexity),
				}
			}
		}

		detailsMethods = append(detailsMethods, md)
	}

	sort.Slice(detailsMethods, func(i, j int) bool {
		return detailsMethods[i].StartLine < detailsMethods[j].StartLine
	})

	fileMetrics := b.buildMetricsMap(fileNode.Metrics)
	totalsData := totals{
		Files:    1,
		Folders:  0,
		Statuses: b.convertStatuses(fileNode.Statuses),
	}

	if sc, ok := fileMetrics["statement_coverage"].(lineCoverageDetail); ok {
		totalsData.StatementCoverage = &sc
	}
	if lc, ok := fileMetrics["line_coverage"].(lineCoverageDetail); ok {
		totalsData.LineCoverage = &lc
	}
	if bc, ok := fileMetrics["branch_coverage"].(branchCoverageDetail); ok {
		totalsData.BranchCoverage = &bc
	}
	if mc, ok := fileMetrics["methods_hit"].(methodsHitDetail); ok {
		totalsData.MethodsHit = &mc
	}
	if mfc, ok := fileMetrics["methods_fully_covered"].(methodsFullyCoveredDetail); ok {
		totalsData.MethodsFullyCovered = &mfc
	}

	if psc, ok := fileMetrics["patch_statement_coverage"].(lineCoverageDetail); ok {
		totalsData.PatchStatementCoverage = &psc
	}
	if plc, ok := fileMetrics["patch_line_coverage"].(lineCoverageDetail); ok {
		totalsData.PatchLineCoverage = &plc
	}
	if pmc, ok := fileMetrics["patch_methods_hit"].(methodsHitDetail); ok {
		totalsData.PatchMethodsHit = &pmc
	}

	if totalMethodBranches > 0 && b.config.ActiveMetrics[config.BranchCoverage] {
		methodBranchPct := utils.CalculatePercentage(coveredMethodBranches, totalMethodBranches, 2)
		totalsData.MethodBranchCoverage = &branchCoverageDetail{
			Covered:    coveredMethodBranches,
			Total:      totalMethodBranches,
			Percentage: methodBranchPct,
		}
	}

	if maxCyclo > 0 && b.config.ActiveMetrics[config.MaxCyclomaticComplexity] {
		totalsData.MaxCyclomaticComplexity = &lineCoverageDetail{
			Total: maxCyclo,
		}
	}

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
		return "", fmt.Errorf("placeholder script 'window.__NANOVISION_DETAILS__' not found in details.html template")
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
