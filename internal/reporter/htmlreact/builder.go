package htmlreact

import (
	"fmt"
	"log/slog"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/IgorBayerl/AdlerCov/internal/model"
	"github.com/IgorBayerl/AdlerCov/internal/reporter"
	"github.com/IgorBayerl/AdlerCov/internal/utils"
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
	b.logger.Info("Starting generation of new React HTML report.")

	summaryData, err := b.transformTree(tree)
	if err != nil {
		return fmt.Errorf("failed to transform coverage data: %w", err)
	}

	if err := GenerateSummary(b.outputDir, summaryData, nil); err != nil {
		return fmt.Errorf("failed to generate summary files: %w", err)
	}

	// *** ADD THIS CALL to generate the details pages ***
	if err := generateDetailsPages(b.outputDir, tree); err != nil {
		return fmt.Errorf("failed to generate details pages: %w", err)
	}

	b.logger.Info("Successfully generated React HTML report.", "directory", b.outputDir)
	return nil
}

func (b *HtmlReactReportBuilder) transformTree(tree *model.SummaryTree) (summaryV1, error) {
	generatedAt := time.Now().UTC()
	treeNodes := b.buildTreeChildren(tree.Root)
	totalFiles, totalFolders := countNodes(treeNodes)

	return summaryV1{
		SchemaVersion:     1,
		GeneratedAt:       generatedAt.Format(time.RFC3339),
		Title:             "Coverage Report",
		Totals:            b.buildTotals(tree, totalFiles, totalFolders),
		Tree:              treeNodes,
		MetricDefinitions: b.buildMetricDefinitions(),
		Metadata:          b.buildMetadata(tree, generatedAt),
	}, nil
}

// addMeta is a helper function for creating metadata items.
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

// buildMetadata creates the report information section.
func (b *HtmlReactReportBuilder) buildMetadata(tree *model.SummaryTree, generatedAt time.Time) []metadataItem {
	meta := make([]metadataItem, 0)

	addMeta(&meta, "Generated At", generatedAt.Format("2006-01-02 15:04:05"))
	if tree.Timestamp > 0 {
		coverageDate := time.Unix(tree.Timestamp, 0).Format("2006-01-02 15:04:05")
		addMeta(&meta, "Coverage Date", coverageDate)
	}
	addMeta(&meta, "Parser", tree.ParserName)
	addMeta(&meta, "Report Files", tree.ReportFiles, "large")

	return meta
}

func (b *HtmlReactReportBuilder) buildTreeChildren(dir *model.DirNode) []fileNode {
	children := make([]fileNode, 0, len(dir.Subdirs)+len(dir.Files))

	// Add subdirectories
	for _, subdir := range dir.Subdirs {
		nodeMetrics, nodeStatuses := b.buildMetricsMap(subdir.Metrics)
		children = append(children, fileNode{
			ID:       subdir.Path,
			Name:     subdir.Name,
			Type:     "folder",
			Path:     subdir.Path,
			Children: b.buildTreeChildren(subdir),
			Metrics:  nodeMetrics,
			Statuses: nodeStatuses,
		})
	}

	// Add files
	for _, file := range dir.Files {
		nodeMetrics, nodeStatuses := b.buildMetricsMap(file.Metrics)

		detailsFileName := strings.ReplaceAll(file.Path, "/", "_") + ".html"

		children = append(children, fileNode{
			ID:        file.Path,
			Name:      file.Name,
			Type:      "file",
			Path:      file.Path,
			Metrics:   nodeMetrics,
			Statuses:  nodeStatuses,
			TargetURL: detailsFileName,
		})
	}

	// Sort for consistent ordering
	sort.Slice(children, func(i, j int) bool {
		if children[i].Type != children[j].Type {
			return children[i].Type == "folder" // Folders first
		}
		return children[i].Name < children[j].Name
	})

	return children
}

func (b *HtmlReactReportBuilder) buildTotals(tree *model.SummaryTree, files, folders int) totals {
	metrics, totalStatuses := b.buildMetricsMap(tree.Metrics)

	t := totals{
		Files:    files,
		Folders:  folders,
		Statuses: totalStatuses,
	}

	if lc, ok := metrics["lineCoverage"].(lineCoverageDetail); ok {
		t.LineCoverage = &lc
	}
	if bc, ok := metrics["branchCoverage"].(branchCoverageDetail); ok {
		t.BranchCoverage = &bc
	}
	if mc, ok := metrics["methodsCovered"].(methodsCoveredDetail); ok {
		t.MethodsCovered = &mc
	}
	if mfc, ok := metrics["methodsFullyCovered"].(methodsFullyCoveredDetail); ok {
		t.MethodsFullyCovered = &mfc
	}
	return t
}

// getRiskStatus determines the risk level based on coverage percentage.
// TODO: get this information from the config file when implementing the components system
func getRiskStatus(percentage float64) riskLevel {
	if percentage >= 80 {
		return RiskSafe
	}
	if percentage >= 60 {
		return RiskWarning
	}
	return RiskDanger
}

func (b *HtmlReactReportBuilder) buildMetricsMap(m model.CoverageMetrics) (metricsMap, statuses) {
	linePct := utils.CalculatePercentage(m.LinesCovered, m.LinesValid, 2)
	if math.IsNaN(linePct) {
		linePct = 0.0
	}

	metrics := metricsMap{
		"lineCoverage": lineCoverageDetail{
			Covered:    m.LinesCovered,
			Uncovered:  m.LinesValid - m.LinesCovered,
			Coverable:  m.LinesValid,
			Total:      m.TotalLines,
			Percentage: linePct,
		},
	}

	nodeStatuses := statuses{
		"lineCoverage": getRiskStatus(linePct),
	}

	if m.BranchesValid > 0 {
		branchPct := utils.CalculatePercentage(m.BranchesCovered, m.BranchesValid, 2)
		if math.IsNaN(branchPct) {
			branchPct = 0.0
		}

		metrics["branchCoverage"] = branchCoverageDetail{
			Covered:    m.BranchesCovered,
			Total:      m.BranchesValid,
			Percentage: branchPct,
		}
		nodeStatuses["branchCoverage"] = getRiskStatus(branchPct)
	}

	if m.MethodsValid > 0 {
		methodsCoveredPct := utils.CalculatePercentage(m.MethodsCovered, m.MethodsValid, 2)
		if math.IsNaN(methodsCoveredPct) {
			methodsCoveredPct = 0.0
		}

		metrics["methodsCovered"] = methodsCoveredDetail{
			Covered:    m.MethodsCovered,
			Total:      m.MethodsValid,
			Percentage: methodsCoveredPct,
		}
		nodeStatuses["methodsCovered"] = getRiskStatus(methodsCoveredPct)

		methodsFullyCoveredPct := utils.CalculatePercentage(m.MethodsFullyCovered, m.MethodsValid, 2)
		if math.IsNaN(methodsFullyCoveredPct) {
			methodsFullyCoveredPct = 0.0
		}

		metrics["methodsFullyCovered"] = methodsFullyCoveredDetail{
			Covered:    m.MethodsFullyCovered,
			Total:      m.MethodsValid,
			Percentage: methodsFullyCoveredPct,
		}
		nodeStatuses["methodsFullyCovered"] = getRiskStatus(methodsFullyCoveredPct)
	}

	return metrics, nodeStatuses
}

func (b *HtmlReactReportBuilder) buildMetricDefinitions() metricDefinitions {
	return metricDefinitions{
		"lineCoverage": {
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
		"branchCoverage": {
			Label:      "Branches",
			ShortLabel: "Branches",
			SubMetrics: []subMetric{
				{ID: "covered", Label: "Covered", Width: 100},
				{ID: "total", Label: "Total", Width: 80},
				{ID: "percentage", Label: "Percentage %", Width: 160},
			},
		},
		"methodsCovered": {
			Label:      "Methods Covered",
			ShortLabel: "Methods Cov.",
			SubMetrics: []subMetric{
				{ID: "covered", Label: "Covered", Width: 80},
				{ID: "total", Label: "Total", Width: 80},
				{ID: "percentage", Label: "Percentage %", Width: 160},
			},
		},
		"methodsFullyCovered": {
			Label:      "Methods Fully Covered",
			ShortLabel: "Methods Full Cov.",
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
		} else {
			folders++
			f, fo := countNodes(node.Children)
			files += f
			folders += fo
		}
	}
	return
}
