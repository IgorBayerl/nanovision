package htmlreact

import (
	"fmt"
	"log/slog"
	"sort"
	"strings"
	"time"

	"github.com/IgorBayerl/nanovision/internal/aggregator"
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/reporter"
	"github.com/IgorBayerl/nanovision/internal/review"
	"github.com/IgorBayerl/nanovision/internal/status/evaluators"
)

type HtmlReactReportBuilder struct {
	outputDir  string
	logger     *slog.Logger
	singleFile bool
	// emit only the changed files and the folders that hold them
	onlyChanged bool
	config      *config.AppConfig
	// set when a server builds the views of a stored run
	view *ViewOptions
}

func NewHtmlReactReportBuilder(outputDir string, logger *slog.Logger, singleFile bool, cfg *config.AppConfig) reporter.ReportBuilder {
	return &HtmlReactReportBuilder{
		outputDir:   outputDir,
		logger:      logger,
		singleFile:  singleFile,
		onlyChanged: cfg.Diff.OnlyChanged,
		config:      cfg,
	}
}

func (b *HtmlReactReportBuilder) ReportType() string {
	if b.singleFile {
		return "HtmlUnified"
	}
	return "Html"
}

func (b *HtmlReactReportBuilder) CreateReport(tree *model.SummaryTree) error {
	b.logger.Info("Starting generation of new React HTML report.", "directory", b.outputDir, "single_file", b.singleFile, "only_changed", b.onlyChanged)

	if b.onlyChanged && tree.Change == nil {
		return fmt.Errorf("diff.only_changed needs a diff: set diff.file in nanovision.yaml, pass -diff, or set vcs.type")
	}

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
		if b.onlyChanged && !isChangedFile(fileNode) {
			continue
		}
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
	if b.view != nil && !b.view.GeneratedAt.IsZero() {
		generatedAt = b.view.GeneratedAt.UTC()
	}
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

	defaultFilters := b.config.DefaultFilters

	// a run measured with a diff gets the verdict on its changed code
	var reviewResult *review.Result
	if tree.Change != nil {
		reviewResult = review.Evaluate(tree, b.config)
	}
	metadata := b.buildMetadata(tree, generatedAt)
	if b.view != nil {
		metadata = b.view.Metadata
	}

	reports, indexes := b.buildReportIndexes(tree)

	return summaryV1{
		SchemaVersion:     1,
		GeneratedAt:       generatedAt.Format(time.RFC3339),
		Title:             title,
		Totals:            totalsData,
		Nodes:             nodes,
		MetricDefinitions: b.buildMetricDefinitions(),
		MetricOrder:       b.metricOrder(),
		Metadata:          metadata,
		DefaultFilters:    defaultFilters,
		Review:            reviewResult,
		Comparison:        tree.Comparison,
		Comparing:         comparingItems(tree),
		OnlyChanged:       b.onlyChanged,
		Reports:           reports,
		ReportIndexes:     indexes,
		StatusBands:       b.buildStatusBands(),
	}, nil
}

// buildReportIndexes compresses every emitted file into per-report buckets so
// the UI can recompute metrics for any subset of reports. Returns nothing when
// a single report was parsed, since there is nothing to toggle.
func (b *HtmlReactReportBuilder) buildReportIndexes(tree *model.SummaryTree) ([]report, map[string]reportIndex) {
	if len(tree.ReportNames) < 2 {
		return nil, nil
	}

	fileMap := make(map[string]*model.FileNode)
	collectFiles(tree.Root, fileMap)

	indexes := make(map[string]reportIndex, len(fileMap))
	for path, file := range fileMap {
		if b.onlyChanged && !isChangedFile(file) {
			continue
		}
		idx := aggregator.BuildFileReportIndex(file, len(tree.ReportNames), b.config.ActiveFileMetrics)
		if len(idx) == 0 {
			continue
		}
		converted := make(reportIndex, len(idx))
		for key, buckets := range idx {
			converted[string(key)] = buckets
		}
		indexes[path] = converted
	}

	if len(indexes) == 0 {
		return nil, nil
	}
	return buildGlobalReports(tree), indexes
}

// buildGlobalReports lists every parsed report in the order the coverage masks
// address them, so one selection means the same thing on every screen.
func buildGlobalReports(tree *model.SummaryTree) []report {
	labels := uniqueReportLabels(tree.ReportNames)

	reports := make([]report, 0, len(tree.ReportNames))
	for i, name := range tree.ReportNames {
		reports = append(reports, report{
			Name: labels[i],
			Path: name,
		})
	}
	return reports
}

// uniqueReportLabels names each report by its file name, keeping just enough
// parent directories to tell apart reports that share one. Merging test shards
// routinely produces several coverage.out files, and "coverage.out" three times
// over is a list nobody can use. The full path is still on the report, for the
// UI to show on hover.
func uniqueReportLabels(paths []string) []string {
	segments := make([][]string, len(paths))
	depths := make([]int, len(paths))
	for i, path := range paths {
		segments[i] = splitPathSegments(path)
		depths[i] = 1
	}

	// Each pass lengthens every still-ambiguous label by one directory. Paths
	// that are genuinely identical stop growing once fully spelled out.
	for range paths {
		counts := make(map[string]int, len(paths))
		for i := range paths {
			counts[trailingSegments(segments[i], depths[i])]++
		}

		grew := false
		for i := range paths {
			if counts[trailingSegments(segments[i], depths[i])] > 1 && depths[i] < len(segments[i]) {
				depths[i]++
				grew = true
			}
		}
		if !grew {
			break
		}
	}

	labels := make([]string, len(paths))
	for i := range paths {
		labels[i] = trailingSegments(segments[i], depths[i])
	}
	return labels
}

func splitPathSegments(path string) []string {
	normalised := strings.ReplaceAll(path, "\\", "/")
	segments := make([]string, 0, strings.Count(normalised, "/")+1)
	for segment := range strings.SplitSeq(normalised, "/") {
		if segment != "" && segment != "." {
			segments = append(segments, segment)
		}
	}
	if len(segments) == 0 {
		return []string{path}
	}
	return segments
}

func trailingSegments(segments []string, depth int) string {
	if depth >= len(segments) {
		return strings.Join(segments, "/")
	}
	return strings.Join(segments[len(segments)-depth:], "/")
}

// metricOrder is file_metrics as configured, so every metric list in the UI
// shows them in the order the user wrote them.
func (b *HtmlReactReportBuilder) metricOrder() []string {
	order := make([]string, len(b.config.FileMetrics))
	for i, key := range b.config.FileMetrics {
		order[i] = string(key)
	}
	return order
}

func (b *HtmlReactReportBuilder) buildStatusBands() map[string]statusBand {
	if len(b.config.StatusBands) == 0 {
		return nil
	}
	bands := make(map[string]statusBand, len(b.config.StatusBands))
	for key, band := range b.config.StatusBands {
		bands[string(key)] = statusBand{Min: band.Min, Max: band.Max}
	}
	return bands
}

func (b *HtmlReactReportBuilder) convertStatuses(modelStatuses map[config.MetricKey]string) statuses {
	uiStatuses := make(statuses)
	for key, val := range modelStatuses {
		uiStatuses[string(key)] = riskLevel(val)
	}
	return uiStatuses
}

func addMeta(meta *[]MetadataItem, label string, value any, sizeHint ...string) {
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

	item := MetadataItem{Label: label, Value: value}
	if len(sizeHint) > 0 {
		item.SizeHint = sizeHint[0]
	}
	*meta = append(*meta, item)
}

func (b *HtmlReactReportBuilder) buildMetadata(tree *model.SummaryTree, generatedAt time.Time) []MetadataItem {
	meta := make([]MetadataItem, 0)

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

// comparingItems says which two states of the code the report compares: the
// revisions, the base run the coverage delta uses, and where the changed
// files come from.
func comparingItems(tree *model.SummaryTree) []MetadataItem {
	var meta []MetadataItem
	if v := tree.Versions; v != nil {
		addMeta(&meta, "Base", revisionLabel(v.Base))
		current := revisionLabel(v.Current)
		switch {
		case v.Current.Revision == "":
			current = strings.TrimSpace(v.Current.Stream + " working copy")
		case v.LocalEdits:
			current += " + edits"
		}
		addMeta(&meta, "Current", current)
	}
	if c := tree.Comparison; c != nil {
		base := shortRevision(c.Base.Revision)
		switch {
		case c.Exact:
		case c.Distance > 0:
			base += fmt.Sprintf(", %d before base", c.Distance)
		default:
			base += ", before base"
		}
		addMeta(&meta, "Coverage base", base)
	}
	if n := tree.Change.Total(); n > 0 {
		files := fmt.Sprintf("%d", n)
		if tree.Versions != nil && tree.Versions.Diff != "" {
			files += " from " + tree.Versions.Diff
		}
		addMeta(&meta, "Changed files", files)
	}
	return meta
}

// revisionLabel is "stream revision"; a Perforce revision already names its stream.
func revisionLabel(r model.RunRef) string {
	rev := shortRevision(r.Revision)
	if r.Stream == "" || strings.Contains(rev, r.Stream) {
		return rev
	}
	return strings.TrimSpace(r.Stream + " " + rev)
}

// shortRevision cuts a git hash to 8 characters and keeps anything else.
func shortRevision(rev string) string {
	if len(rev) != 40 || strings.Trim(rev, "0123456789abcdef") != "" {
		return rev
	}
	return rev[:8]
}

// depth-first, pre-ordered. siblings go folders first, then files, each sorted by name.
//
// With diff.only_changed only the changed files and the folders needed to
// reach them are emitted.
func (b *HtmlReactReportBuilder) buildFlatNodes(root *model.DirNode) []fileNode {
	nodes := make([]fileNode, 0, len(root.Subdirs)+len(root.Files))
	b.appendFlatNodes(&nodes, root, "", 0)
	return nodes
}

// returns false when the subtree emitted nothing, which drops empty folders with diff.only_changed.
func (b *HtmlReactReportBuilder) appendFlatNodes(out *[]fileNode, dir *model.DirNode, parentID string, depth int) bool {
	emitted := false

	// Subdirs/Files are maps (non-deterministic order), so collect and sort by name.
	subdirs := make([]*model.DirNode, 0, len(dir.Subdirs))
	for _, subdir := range dir.Subdirs {
		subdirs = append(subdirs, subdir)
	}
	sort.Slice(subdirs, func(i, j int) bool { return subdirs[i].Name < subdirs[j].Name })

	// Folders first, each followed immediately by its subtree (pre-order).
	for _, subdir := range subdirs {
		var children []fileNode
		hasChildren := b.appendFlatNodes(&children, subdir, subdir.Path, depth+1)
		if b.onlyChanged && !hasChildren {
			continue
		}
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
		*out = append(*out, children...)
		emitted = true
	}

	files := make([]*model.FileNode, 0, len(dir.Files))
	for _, file := range dir.Files {
		files = append(files, file)
	}
	sort.Slice(files, func(i, j int) bool { return files[i].Name < files[j].Name })

	for _, file := range files {
		if b.onlyChanged && !isChangedFile(file) {
			continue
		}

		target := ""
		if b.view != nil && b.view.FileURL != nil {
			target = b.view.FileURL(file.Path)
		} else if file.SourceDir != "" {
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
		emitted = true
	}

	return emitted
}

func isChangedFile(file *model.FileNode) bool {
	return file.Diff != nil && file.Diff.Kind != model.ChangeKindNone
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
	if mc, ok := metrics[string(config.MethodsHit)].(methodsHitDetail); ok {
		t.MethodsHit = &mc
	}
	if mfc, ok := metrics[string(config.MethodsFullyCovered)].(methodsFullyCoveredDetail); ok {
		t.MethodsFullyCovered = &mfc
	}

	if psc, ok := metrics[string(config.PatchStatementCoverage)].(lineCoverageDetail); ok {
		t.PatchStatementCoverage = &psc
	}
	if plc, ok := metrics[string(config.PatchLineCoverage)].(lineCoverageDetail); ok {
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

	if b.config.ActiveMethodMetrics[config.MethodExposedRisk] {
		defs[MethodUIExposedRisk] = metricDefinition{
			Label:      "Exposed Risk",
			ShortLabel: "Risk",
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

	for key, def := range defs {
		def.Description = describeMetric(key)
		defs[key] = def
	}

	return defs
}

// uiMetricKeys maps the sort-prefixed method metric keys the UI uses back to
// the config key whose evaluator owns the description.
var uiMetricKeys = map[string]config.MetricKey{
	MethodUIStmtCoverage:         config.MethodStatementCoverage,
	MethodUILineCoverage:         config.MethodLineCoverage,
	MethodUIPatchStmtCoverage:    config.MethodPatchStatementCoverage,
	MethodUIPatchLineCoverage:    config.MethodPatchLineCoverage,
	MethodUICyclomaticComplexity: config.CyclomaticComplexity,
	MethodUICrapScore:            config.MethodCrapScore,
	MethodUIExposedRisk:          config.MethodExposedRisk,
}

// describeMetric returns the one-line explanation shown in the UI tooltips. The
// evaluator is the source of truth, so a metric is documented in one place.
func describeMetric(key string) string {
	metricKey := config.MetricKey(key)
	if mapped, ok := uiMetricKeys[key]; ok {
		metricKey = mapped
	}

	if evaluator, ok := evaluators.Registry[metricKey]; ok {
		return evaluator.Description()
	}
	return ""
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
