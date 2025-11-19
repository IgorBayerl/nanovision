package htmlreact

import (
	"fmt"
	"log/slog"
	"sort"
	"strings"
	"time"

	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/reporter"
	"github.com/IgorBayerl/nanovision/internal/utils"
)

type HtmlReactReportBuilder struct {
	outputDir string
	logger    *slog.Logger
}

func NewHtmlReactReportBuilder(outputDir string, logger *slog.Logger) reporter.ReportBuilder {
	return &HtmlReactReportBuilder{
		outputDir: outputDir,
		logger:    logger,
	}
}

func (b *HtmlReactReportBuilder) ReportType() string {
	return "Html"
}

func (b *HtmlReactReportBuilder) CreateReport(tree *model.SummaryTree) error {
	b.logger.Info("Starting generation of new React HTML report.", "directory", b.outputDir)

	summaryData, err := b.transformTree(tree)
	if err != nil {
		return fmt.Errorf("failed to transform coverage data: %w", err)
	}

	// NOTE: GenerateSummary expects a Logger interface (with Debugf/Infof/Errorf).
	// slog.Logger does not implement that, so we pass nil (logging inside
	// GenerateSummary/copyDist is optional).
	if err := GenerateSummary(b.outputDir, summaryData, nil); err != nil {
		return fmt.Errorf("failed to generate summary files: %w", err)
	}

	if err := generateDetailsPages(b, tree); err != nil {
		return fmt.Errorf("failed to generate details pages: %w", err)
	}

	b.logger.Info("Successfully generated React HTML report.", "directory", b.outputDir)
	return nil
}

func (b *HtmlReactReportBuilder) transformTree(tree *model.SummaryTree) (summaryV1, error) {
	generatedAt := time.Now().UTC()
	treeNodes := b.buildTreeChildren(tree.Root)
	totalFiles, totalFolders := countNodes(treeNodes)

	totalsData := b.buildTotals(tree, totalFiles, totalFolders)
	if tree.Root.Statuses != nil {
		totalsData.Statuses = b.convertStatuses(tree.Root.Statuses)
	}

	return summaryV1{
		SchemaVersion:     1,
		GeneratedAt:       generatedAt.Format(time.RFC3339),
		Title:             "Coverage Report",
		Totals:            totalsData,
		Tree:              treeNodes,
		MetricDefinitions: b.buildMetricDefinitions(),
		Metadata:          b.buildMetadata(tree, generatedAt),
	}, nil
}

func (b *HtmlReactReportBuilder) convertStatuses(modelStatuses map[config.MetricKey]string) statuses {
	uiStatuses := make(statuses)
	for key, val := range modelStatuses {
		uiStatuses[string(key)] = riskLevel(val)
	}
	return uiStatuses
}

func addMeta(meta *[]metadataItem, label string, value any, sizeHint ...string) {
	switch v := value.(type) {
	case string:
		if v == "" {
			return
		}
	case int:
		if v == 0 {
			return
		}
	case []string:
		if len(v) == 0 {
			return
		}
	}

	item := metadataItem{Label: label, Value: value}
	if len(sizeHint) > 0 {
		item.SizeHint = sizeHint[0]
	}
	*meta = append(*meta, item)
}

func (b *HtmlReactReportBuilder) buildMetadata(tree *model.SummaryTree, generatedAt time.Time) []metadataItem {
	meta := make([]metadataItem, 0)

	addMeta(&meta, "Generated At", generatedAt.Format("2006-01-02 15:04:05"))
	if tree.Timestamp > 0 {
		coverageDate := time.Unix(tree.Timestamp, 0).Format("2006-01-02 15:04:05")
		addMeta(&meta, "Coverage Date", coverageDate)
	}
	if len(tree.ParserNames) > 0 {
		parserValue := strings.Join(tree.ParserNames, " | ")
		addMeta(&meta, "Parser", parserValue)
	}
	addMeta(&meta, "Report Files", tree.ReportFiles, "large")

	return meta
}

func (b *HtmlReactReportBuilder) buildTreeChildren(dir *model.DirNode) []fileNode {
	children := make([]fileNode, 0, len(dir.Subdirs)+len(dir.Files))

	for _, subdir := range dir.Subdirs {
		nodeMetrics := b.buildMetricsMap(subdir.Metrics)
		statuses := b.convertStatuses(subdir.Statuses)

		child := fileNode{
			ID:       subdir.Path,
			Name:     subdir.Name,
			Type:     "folder",
			Path:     subdir.Path,
			Metrics:  nodeMetrics,
			Statuses: statuses,
			Children: b.buildTreeChildren(subdir),
		}
		children = append(children, child)
	}

	for _, file := range dir.Files {
		nodeMetrics := b.buildMetricsMap(file.Metrics)
		statuses := b.convertStatuses(file.Statuses)

		target := ""
		if file.SourceDir != "" {
			target = fmt.Sprintf("%s.html", strings.ReplaceAll(file.Path, "/", "_"))
		}

		diffStatus := ""
		if file.Diff != nil {
			diffStatus = file.Diff.Kind.String()
		}

		child := fileNode{
			ID:         file.Path,
			Name:       file.Name,
			Type:       "file",
			Path:       file.Path,
			Metrics:    nodeMetrics,
			Statuses:   statuses,
			TargetURL:  target,
			DiffStatus: diffStatus,
		}
		children = append(children, child)
	}

	// Sort children so folders come before files, then by name.
	sort.Slice(children, func(i, j int) bool {
		if children[i].Type != children[j].Type {
			return children[i].Type == "folder"
		}
		return children[i].Name < children[j].Name
	})

	return children
}

func (b *HtmlReactReportBuilder) buildTotals(tree *model.SummaryTree, files, folders int) totals {
	metrics := b.buildMetricsMap(tree.Metrics)

	t := totals{
		Files:   files,
		Folders: folders,
	}

	if lc, ok := metrics["line_coverage"].(lineCoverageDetail); ok {
		t.LineCoverage = &lc
	}
	if bc, ok := metrics["branch_coverage"].(branchCoverageDetail); ok {
		t.BranchCoverage = &bc
	}
	if mc, ok := metrics["methods_covered"].(methodsCoveredDetail); ok {
		t.MethodsCovered = &mc
	}
	if mfc, ok := metrics["methods_fully_covered"].(methodsFullyCoveredDetail); ok {
		t.MethodsFullyCovered = &mfc
	}

	if plc, ok := metrics["patch_line_coverage"].(lineCoverageDetail); ok { // Changed branchCoverageDetail to lineCoverageDetail
		t.PatchLineCoverage = &plc
	}

	if pmc, ok := metrics["patch_methods_covered"].(methodsCoveredDetail); ok {
		t.PatchMethodsCovered = &pmc
	}

	return t
}

func (b *HtmlReactReportBuilder) buildMetricsMap(m model.CoverageMetrics) metricsMap {
	linePct := utils.CalculatePercentage(m.LinesCovered, m.LinesValid, 2)
	metrics := metricsMap{
		"line_coverage": lineCoverageDetail{
			Covered:    m.LinesCovered,
			Uncovered:  m.LinesValid - m.LinesCovered,
			Coverable:  m.LinesValid,
			Total:      m.TotalLines,
			Percentage: linePct,
		},
	}

	if m.BranchesValid > 0 {
		branchPct := utils.CalculatePercentage(m.BranchesCovered, m.BranchesValid, 2)
		metrics["branch_coverage"] = branchCoverageDetail{
			Covered:    m.BranchesCovered,
			Total:      m.BranchesValid,
			Percentage: branchPct,
		}
	}

	if m.MethodsValid > 0 {
		methodsCoveredPct := utils.CalculatePercentage(m.MethodsCovered, m.MethodsValid, 2)
		metrics["methods_covered"] = methodsCoveredDetail{
			Covered:    m.MethodsCovered,
			Total:      m.MethodsValid,
			Percentage: methodsCoveredPct,
		}

		methodsFullyCoveredPct := utils.CalculatePercentage(m.MethodsFullyCovered, m.MethodsValid, 2)
		metrics["methods_fully_covered"] = methodsFullyCoveredDetail{
			Covered:    m.MethodsFullyCovered,
			Total:      m.MethodsValid,
			Percentage: methodsFullyCoveredPct,
		}
	}

	// Patch line coverage
	if m.PatchLinesTotal > 0 {
		patchLinePct := utils.CalculatePercentage(m.PatchLinesCovered, m.PatchLinesValid, 2)
		metrics["patch_line_coverage"] = lineCoverageDetail{
			Covered:    m.PatchLinesCovered,
			Uncovered:  m.PatchLinesValid - m.PatchLinesCovered,
			Coverable:  m.PatchLinesValid,
			Total:      m.PatchLinesTotal,
			Percentage: patchLinePct,
		}
	}

	// Patch methods coverage
	if m.PatchMethodsValid > 0 {
		patchMethodsPct := utils.CalculatePercentage(m.PatchMethodsCovered, m.PatchMethodsValid, 2)
		metrics["patch_methods_covered"] = methodsCoveredDetail{
			Covered:    m.PatchMethodsCovered,
			Total:      m.PatchMethodsValid,
			Percentage: patchMethodsPct,
		}
	}

	return metrics
}

func (b *HtmlReactReportBuilder) buildMetricDefinitions() metricDefinitions {
	return metricDefinitions{
		"line_coverage": {
			Label:      "Lines",
			ShortLabel: "Lines",
			SubMetrics: []subMetric{
				{ID: "covered", Label: "Covered", Width: 100},
				{ID: "uncovered", Label: "Uncovered", Width: 100},
				{ID: "coverable", Label: "Coverable", Width: 100},
				{ID: "total", Label: "Total", Width: 80},
				{ID: "percentage", Label: "Percentage %", Width: 160},
			},
		},
		"branch_coverage": {
			Label:      "Branches",
			ShortLabel: "Branches",
			SubMetrics: []subMetric{
				{ID: "covered", Label: "Covered", Width: 100},
				{ID: "total", Label: "Total", Width: 80},
				{ID: "percentage", Label: "Percentage %", Width: 160},
			},
		},
		"methods_covered": {
			Label:      "Methods Covered",
			ShortLabel: "Methods",
			SubMetrics: []subMetric{
				{ID: "covered", Label: "Covered", Width: 80},
				{ID: "total", Label: "Total", Width: 80},
				{ID: "percentage", Label: "Percentage %", Width: 160},
			},
		},
		"methods_fully_covered": {
			Label:      "Methods Fully Covered",
			ShortLabel: "Fully Covered",
			SubMetrics: []subMetric{
				{ID: "covered", Label: "Covered", Width: 80},
				{ID: "total", Label: "Total", Width: 80},
				{ID: "percentage", Label: "Percentage %", Width: 160},
			},
		},
		"method_branch_coverage": {
			Label:      "Method Branches",
			ShortLabel: "Method Branches",
			SubMetrics: []subMetric{
				{ID: "covered", Label: "Covered", Width: 100},
				{ID: "total", Label: "Total", Width: 80},
				{ID: "percentage", Label: "Percentage %", Width: 160},
			},
		},
		"max_cyclomatic_complexity": {
			Label:      "Max Cyclomatic Complexity",
			ShortLabel: "Max Complexity",
			SubMetrics: []subMetric{
				{ID: "total", Label: "Value", Width: 100},
			},
		},
		"patch_line_coverage": {
			Label:      "Patch Lines",
			ShortLabel: "Patch Lines",
			SubMetrics: []subMetric{
				{ID: "covered", Label: "Covered", Width: 100},
				{ID: "uncovered", Label: "Uncovered", Width: 100},
				{ID: "coverable", Label: "Coverable", Width: 100},
				{ID: "total", Label: "Total", Width: 80},
				{ID: "percentage", Label: "Percentage %", Width: 160},
			},
		},
		"patch_methods_covered": {
			Label:      "Patch Methods",
			ShortLabel: "Patch Methods",
			SubMetrics: []subMetric{
				{ID: "covered", Label: "Covered", Width: 80},
				{ID: "total", Label: "Total", Width: 80},
				{ID: "percentage", Label: "Percentage %", Width: 160},
			},
		},
	}
}

func countNodes(nodes []fileNode) (files, folders int) {
	for _, node := range nodes {
		if node.Type == "file" {
			files++
		} else { // folder
			folders++
			f, fo := countNodes(node.Children)
			files += f
			folders += fo
		}
	}
	return
}
