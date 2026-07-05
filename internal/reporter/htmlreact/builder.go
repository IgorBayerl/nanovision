package htmlreact

import (
	"fmt"
	"log/slog"
	"sort"
	"strings"
	"time"

	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/diagnostics"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/reporter"
	"github.com/IgorBayerl/nanovision/internal/status/evaluators"
)

type HtmlReactReportBuilder struct {
	outputDir  string
	logger     *slog.Logger
	singleFile bool
	config     *config.AppConfig
}

func NewHtmlReactReportBuilder(outputDir string, logger *slog.Logger, singleFile bool, cfg *config.AppConfig) reporter.ReportBuilder {
	return &HtmlReactReportBuilder{
		outputDir:  outputDir,
		logger:     logger,
		singleFile: singleFile,
		config:     cfg,
	}
}

func (b *HtmlReactReportBuilder) ReportType() string {
	if b.singleFile {
		return "HtmlUnified"
	}
	return "Html"
}

func (b *HtmlReactReportBuilder) CreateReport(tree *model.SummaryTree) error {
	b.logger.Info("Starting generation of new React HTML report.", "directory", b.outputDir, "single_file", b.singleFile)

	if b.singleFile {
		return b.createSingleFileReport(tree)
	}

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

func (b *HtmlReactReportBuilder) createSingleFileReport(tree *model.SummaryTree) error {
	summaryData, err := b.transformTree(tree)
	if err != nil {
		return fmt.Errorf("failed to transform summary data: %w", err)
	}

	allDetails := make(map[string]*detailsV1)
	fileMap := make(map[string]*model.FileNode)

	collectFiles(tree.Root, fileMap)

	for path, fileNode := range fileMap {
		details, err := b.transformFileNodeToDetails(tree, fileNode)
		if err != nil {
			b.logger.Warn("Failed to generate details for file", "path", path, "error", err)
			continue
		}
		allDetails[path] = details
	}

	if err := GenerateSingleFile(b.outputDir, summaryData, allDetails, nil); err != nil {
		return fmt.Errorf("failed to generate single file report: %w", err)
	}

	b.logger.Info("Successfully generated React HTML report (Single File Mode).", "directory", b.outputDir)
	return nil
}

func (b *HtmlReactReportBuilder) transformTree(tree *model.SummaryTree) (summaryV1, error) {
	generatedAt := time.Now().UTC()
	nodes := b.buildFlatNodes(tree.Root)
	totalFiles, totalFolders := countFlatNodes(nodes)

	totalsData := b.buildTotals(tree, totalFiles, totalFolders)
	if tree.Root.Statuses != nil {
		totalsData.Statuses = b.convertStatuses(tree.Root.Statuses)
	}

	title := b.config.Title
	if title == "" {
		title = "Coverage Report"
	}

	return summaryV1{
		SchemaVersion:     1,
		GeneratedAt:       generatedAt.Format(time.RFC3339),
		Title:             title,
		Totals:            totalsData,
		Nodes:             nodes,
		MetricDefinitions: b.buildMetricDefinitions(),
		Metadata:          b.buildMetadata(tree, generatedAt),
		Diagnostics:       diagnostics.Extract(tree, b.config, evaluators.Registry),
		DefaultFilters:    b.config.DefaultFilters,
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

// buildFlatNodes walks the directory tree depth-first and returns a pre-ordered
// flat slice of nodes. Each node records its ParentID and Depth so the client can
// rebuild the hierarchy in a single linear pass. Sibling order matches the old
// nested layout: folders before files, each group sorted by name.
func (b *HtmlReactReportBuilder) buildFlatNodes(root *model.DirNode) []fileNode {
	nodes := make([]fileNode, 0, len(root.Subdirs)+len(root.Files))
	b.appendFlatNodes(&nodes, root, "", 0)
	return nodes
}

func (b *HtmlReactReportBuilder) appendFlatNodes(out *[]fileNode, dir *model.DirNode, parentID string, depth int) {
	// Subdirs/Files are maps (non-deterministic order), so collect and sort by name.
	subdirs := make([]*model.DirNode, 0, len(dir.Subdirs))
	for _, subdir := range dir.Subdirs {
		subdirs = append(subdirs, subdir)
	}
	sort.Slice(subdirs, func(i, j int) bool { return subdirs[i].Name < subdirs[j].Name })

	// Folders first, each followed immediately by its subtree (pre-order).
	for _, subdir := range subdirs {
		*out = append(*out, fileNode{
			ID:       subdir.Path,
			Name:     subdir.Name,
			Type:     "folder",
			Path:     subdir.Path,
			ParentID: parentID,
			Depth:    depth,
			Metrics:  b.buildMetricsMap(subdir.Metrics),
			Statuses: b.convertStatuses(subdir.Statuses),
		})
		b.appendFlatNodes(out, subdir, subdir.Path, depth+1)
	}

	files := make([]*model.FileNode, 0, len(dir.Files))
	for _, file := range dir.Files {
		files = append(files, file)
	}
	sort.Slice(files, func(i, j int) bool { return files[i].Name < files[j].Name })

	for _, file := range files {
		target := ""
		if file.SourceDir != "" {
			if b.singleFile {
				target = fmt.Sprintf("#/details/%s", file.Path)
			} else {
				target = fmt.Sprintf("%s.html", strings.ReplaceAll(file.Path, "/", "_"))
			}
		}

		diffStatus := ""
		if file.Diff != nil {
			diffStatus = file.Diff.Kind.String()
		}

		*out = append(*out, fileNode{
			ID:         file.Path,
			Name:       file.Name,
			Type:       "file",
			Path:       file.Path,
			ParentID:   parentID,
			Depth:      depth,
			Metrics:    b.buildMetricsMap(file.Metrics),
			Statuses:   b.convertStatuses(file.Statuses),
			TargetURL:  target,
			DiffStatus: diffStatus,
		})
	}
}

func (b *HtmlReactReportBuilder) buildTotals(tree *model.SummaryTree, files, folders int) totals {
	metrics := b.buildMetricsMap(tree.Metrics)

	t := totals{
		Files:   files,
		Folders: folders,
	}

	if sc, ok := metrics[string(config.StatementCoverage)].(lineCoverageDetail); ok {
		t.StatementCoverage = &sc
	}
	if lc, ok := metrics[string(config.LineCoverage)].(lineCoverageDetail); ok {
		t.LineCoverage = &lc
	}
	if bc, ok := metrics[string(config.BranchCoverage)].(branchCoverageDetail); ok {
		t.BranchCoverage = &bc
	}
	if mc, ok := metrics[string(config.MethodsHit)].(methodsHitDetail); ok {
		t.MethodsHit = &mc
	}
	if mfc, ok := metrics[string(config.MethodsFullyCovered)].(methodsFullyCoveredDetail); ok {
		t.MethodsFullyCovered = &mfc
	}

	if psc, ok := metrics[string(config.PatchStatementCoverage)].(lineCoverageDetail); ok {
		t.PatchStatementCoverage = &psc
	}
	if plc, ok := metrics[string(config.PatchLineCoverage)].(lineCoverageDetail); ok { // Changed branchCoverageDetail to lineCoverageDetail
		t.PatchLineCoverage = &plc
	}

	if pmc, ok := metrics[string(config.PatchMethodsHit)].(methodsHitDetail); ok {
		t.PatchMethodsHit = &pmc
	}

	if mcc, ok := metrics[string(config.MaxCyclomaticComplexity)].(scoreDetail); ok {
		t.MaxCyclomaticComplexity = &mcc
	}

	return t
}

func (b *HtmlReactReportBuilder) buildMetricsMap(m model.CoverageMetrics) metricsMap {
	metrics := metricsMap{}
	for key := range b.config.ActiveFileMetrics {
		if calcData, exists := m.Calculated[key]; exists {
			switch key {
			case config.LineCoverage:
				if detail, ok := calcData.(model.CoverageDetail); ok {
					metrics[string(key)] = lineCoverageDetail{Covered: detail.Covered, Uncovered: detail.Uncovered, Coverable: detail.Total, Total: m.TotalLines, Percentage: detail.Percentage}
				}
			case config.StatementCoverage, config.PatchStatementCoverage:
				if detail, ok := calcData.(model.CoverageDetail); ok {
					metrics[string(key)] = lineCoverageDetail{Covered: detail.Covered, Uncovered: detail.Uncovered, Coverable: detail.Total, Total: detail.Total, Percentage: detail.Percentage}
				}
			case config.PatchLineCoverage:
				if detail, ok := calcData.(model.CoverageDetail); ok {
					total := m.PatchLinesTotal
					if total == 0 {
						total = detail.Total
					}
					metrics[string(key)] = lineCoverageDetail{Covered: detail.Covered, Uncovered: detail.Uncovered, Coverable: detail.Total, Total: total, Percentage: detail.Percentage}
				}
			case config.BranchCoverage:
				if detail, ok := calcData.(model.CoverageDetail); ok {
					metrics[string(key)] = branchCoverageDetail{Covered: detail.Covered, Total: detail.Total, Percentage: detail.Percentage}
				}
			case config.MethodsHit, config.PatchMethodsHit:
				if detail, ok := calcData.(model.CoverageDetail); ok {
					metrics[string(key)] = methodsHitDetail{Covered: detail.Covered, Total: detail.Total, Percentage: detail.Percentage}
				}
			case config.MethodsFullyCovered:
				if detail, ok := calcData.(model.CoverageDetail); ok {
					metrics[string(key)] = methodsFullyCoveredDetail{Covered: detail.Covered, Total: detail.Total, Percentage: detail.Percentage}
				}
			case config.MaxCyclomaticComplexity:
				if score, ok := calcData.(model.ScoreDetail); ok {
					metrics[string(key)] = scoreDetail{Value: score.Value}
				}
			default:
				metrics[string(key)] = calcData
			}
		}
	}
	return metrics
}

func (b *HtmlReactReportBuilder) buildMetricDefinitions() metricDefinitions {
	defs := metricDefinitions{}

	if b.config.ActiveFileMetrics[config.StatementCoverage] {
		defs[string(config.StatementCoverage)] = metricDefinition{
			Label:      "Statements",
			ShortLabel: "Statements",
			SubMetrics: []subMetric{
				{ID: "covered", Label: "Covered", Width: 100},
				{ID: "uncovered", Label: "Uncovered", Width: 100},
				{ID: "total", Label: "Total", Width: 80},
				{ID: "percentage", Label: "Percentage %", Width: 160},
			},
		}
	}
	if b.config.ActiveMethodMetrics[config.MethodStatementCoverage] {
		defs[MethodUIStmtCoverage] = metricDefinition{
			Label:      "Statements",
			ShortLabel: "Statements",
			SubMetrics: []subMetric{{ID: "total", Label: "Value", Width: 100}},
		}
	}

	if b.config.ActiveMethodMetrics[config.MethodCrapScore] {
		defs[MethodUICrapScore] = metricDefinition{
			Label:      "CRAP Score",
			ShortLabel: "CRAP",
			SubMetrics: []subMetric{{ID: "total", Label: "Value", Width: 100}},
		}
	}

	if b.config.ActiveMethodMetrics[config.MethodPatchCrapScore] {
		defs[MethodUIPatchCrapScore] = metricDefinition{
			Label:      "Patch CRAP Score",
			ShortLabel: "PCRAP",
			SubMetrics: []subMetric{{ID: "total", Label: "Value", Width: 100}},
		}
	}

	if b.config.ActiveMethodMetrics[config.MethodExposedRisk] {
		defs[MethodUIExposedRisk] = metricDefinition{
			Label:      "Exposed Risk",
			ShortLabel: "Risk",
			SubMetrics: []subMetric{{ID: "total", Label: "Value", Width: 100}},
		}
	}

	if b.config.ActiveMethodMetrics[config.MethodDefectProbability] {
		defs[MethodUIDefectProbability] = metricDefinition{
			Label:      "Defect Probability",
			ShortLabel: "DPI",
			SubMetrics: []subMetric{{ID: "total", Label: "Value", Width: 100}},
		}
	}

	if b.config.ActiveFileMetrics[config.LineCoverage] {
		defs[string(config.LineCoverage)] = metricDefinition{
			Label:      "Lines",
			ShortLabel: "Lines",
			SubMetrics: []subMetric{
				{ID: "covered", Label: "Covered", Width: 100},
				{ID: "uncovered", Label: "Uncovered", Width: 100},
				{ID: "coverable", Label: "Coverable", Width: 100},
				{ID: "total", Label: "Total", Width: 80},
				{ID: "percentage", Label: "Percentage %", Width: 160},
			},
		}
	}
	if b.config.ActiveMethodMetrics[config.MethodLineCoverage] {
		defs[MethodUILineCoverage] = metricDefinition{
			Label:      "Lines",
			ShortLabel: "Lines",
			SubMetrics: []subMetric{{ID: "total", Label: "Value", Width: 100}},
		}
	}

	if b.config.ActiveFileMetrics[config.BranchCoverage] {
		defs[string(config.BranchCoverage)] = metricDefinition{
			Label:      "Branches",
			ShortLabel: "Branches",
			SubMetrics: []subMetric{
				{ID: "covered", Label: "Covered", Width: 100},
				{ID: "total", Label: "Total", Width: 80},
				{ID: "percentage", Label: "Percentage %", Width: 160},
			},
		}
	}
	if b.config.ActiveMethodMetrics[config.MethodBranchCoverage] {
		defs[string(config.MethodBranchCoverage)] = metricDefinition{
			Label:      "Method Branches",
			ShortLabel: "Method Branches",
			SubMetrics: []subMetric{
				{ID: "covered", Label: "Covered", Width: 100},
				{ID: "total", Label: "Total", Width: 80},
				{ID: "percentage", Label: "Percentage %", Width: 160},
			},
		}
		defs[MethodUIBranchCoverage] = metricDefinition{
			Label:      "Branches",
			ShortLabel: "Branches",
			SubMetrics: []subMetric{{ID: "total", Label: "Value", Width: 100}},
		}
	}

	if b.config.ActiveFileMetrics[config.PatchStatementCoverage] {
		defs[string(config.PatchStatementCoverage)] = metricDefinition{
			Label:      "Patch Statements",
			ShortLabel: "Patch Statements",
			SubMetrics: []subMetric{
				{ID: "covered", Label: "Covered", Width: 100},
				{ID: "uncovered", Label: "Uncovered", Width: 100},
				{ID: "total", Label: "Total", Width: 80},
				{ID: "percentage", Label: "Percentage %", Width: 160},
			},
		}
	}
	if b.config.ActiveMethodMetrics[config.MethodPatchStatementCoverage] {
		defs[MethodUIPatchStmtCoverage] = metricDefinition{
			Label:      "Patch Statements",
			ShortLabel: "Patch Stmts",
			SubMetrics: []subMetric{{ID: "total", Label: "Value", Width: 100}},
		}
	}

	if b.config.ActiveFileMetrics[config.PatchLineCoverage] {
		defs[string(config.PatchLineCoverage)] = metricDefinition{
			Label:      "Patch Lines",
			ShortLabel: "Patch Lines",
			SubMetrics: []subMetric{
				{ID: "covered", Label: "Covered", Width: 100},
				{ID: "uncovered", Label: "Uncovered", Width: 100},
				{ID: "coverable", Label: "Coverable", Width: 100},
				{ID: "total", Label: "Total", Width: 80},
				{ID: "percentage", Label: "Percentage %", Width: 160},
			},
		}
	}
	if b.config.ActiveMethodMetrics[config.MethodPatchLineCoverage] {
		defs[MethodUIPatchLineCoverage] = metricDefinition{
			Label:      "Patch Lines",
			ShortLabel: "Patch Lines",
			SubMetrics: []subMetric{{ID: "total", Label: "Value", Width: 100}},
		}
	}

	if b.config.ActiveFileMetrics[config.MaxCyclomaticComplexity] {
		defs[string(config.MaxCyclomaticComplexity)] = metricDefinition{
			Label:      "Max Cyclomatic Complexity",
			ShortLabel: "Max Complexity",
			Kind:       "value",
			SubMetrics: []subMetric{
				// Wide enough that the "Max Complexity" header never wraps.
				{ID: "value", Label: "Value", Width: 140},
			},
		}
	}
	if b.config.ActiveMethodMetrics[config.CyclomaticComplexity] {
		defs[MethodUICyclomaticComplexity] = metricDefinition{
			Label:      "Cyclomatic Complexity",
			ShortLabel: "Complexity",
			Kind:       "value",
			SubMetrics: []subMetric{{ID: "value", Label: "Value", Width: 100}},
		}
	}

	if b.config.ActiveFileMetrics[config.MethodsHit] {
		defs[string(config.MethodsHit)] = metricDefinition{
			Label:      "Methods Hit",
			ShortLabel: "Methods Hit",
			SubMetrics: []subMetric{
				{ID: "covered", Label: "Hit", Width: 80},
				{ID: "total", Label: "Total", Width: 80},
				{ID: "percentage", Label: "Percentage %", Width: 160},
			},
		}
	}

	if b.config.ActiveFileMetrics[config.MethodsFullyCovered] {
		defs[string(config.MethodsFullyCovered)] = metricDefinition{
			Label:      "Methods Fully Covered",
			ShortLabel: "Fully Covered",
			SubMetrics: []subMetric{
				{ID: "covered", Label: "Covered", Width: 80},
				{ID: "total", Label: "Total", Width: 80},
				{ID: "percentage", Label: "Percentage %", Width: 160},
			},
		}
	}

	if b.config.ActiveFileMetrics[config.PatchMethodsHit] {
		defs[string(config.PatchMethodsHit)] = metricDefinition{
			Label:      "Patch Methods Hit",
			ShortLabel: "Patch Methods Hit",
			SubMetrics: []subMetric{
				{ID: "covered", Label: "Hit", Width: 80},
				{ID: "total", Label: "Total", Width: 80},
				{ID: "percentage", Label: "Percentage %", Width: 160},
			},
		}
	}

	return defs
}

func countFlatNodes(nodes []fileNode) (files, folders int) {
	for _, node := range nodes {
		if node.Type == "file" {
			files++
		} else {
			folders++
		}
	}
	return
}
