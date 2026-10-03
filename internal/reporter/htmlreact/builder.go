package htmlreact

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/IgorBayerl/nanovision/internal/aggregator"
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/reporter"
	"github.com/IgorBayerl/nanovision/internal/review"
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
		StatusBands:       buildStatusBands(b.config.StatusBands),
		FolderBands:       b.buildFolderBands(),
		Configs:           b.buildConfigs(),
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
	return b.buildGlobalReports(tree), indexes
}

// buildGlobalReports lists every parsed report in the order the coverage masks
// address them, so one selection means the same thing on every screen.
func (b *HtmlReactReportBuilder) buildGlobalReports(tree *model.SummaryTree) []report {
	labels := uniqueReportLabels(tree.ReportNames)
	// a name the config gives a report wins over its file name
	for i, pattern := range tree.ReportNames {
		for _, pair := range b.config.InputPairs {
			if pair.ReportPattern == pattern && pair.Name != "" {
				labels[i] = pair.Name
			}
		}
	}

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

func buildStatusBands(source config.StatusBands) map[string]statusBand {
	if len(source) == 0 {
		return nil
	}
	bands := make(map[string]statusBand, len(source))
	for key, band := range source {
		bands[string(key)] = statusBand{Min: band.Min, Max: band.Max}
	}
	return bands
}

// buildConfigs lists the config files that shape the report, so the UI can
// show which folders have their own settings and where they come from.
func (b *HtmlReactReportBuilder) buildConfigs() []configFile {
	cfg := b.config
	if len(cfg.Folders) == 0 {
		return nil
	}
	root := filepath.Base(cfg.ConfigFile)
	out := []configFile{{Source: root}}
	for _, o := range cfg.Folders {
		out = append(out, configFile{Path: o.Path, Source: folderConfigSource(cfg, o)})
	}
	return out
}

// configFileIn returns the name of the config file that sits in a folder of
// the tree, or "" when the folder has none. parentless is the tree root.
func (b *HtmlReactReportBuilder) configFileIn(dirPath string) string {
	cfg := b.config
	if dirPath == "." || dirPath == "" {
		if cfg.ConfigFile == "" {
			return ""
		}
		return filepath.Base(cfg.ConfigFile)
	}
	for _, o := range cfg.Folders {
		if o.Path == dirPath {
			return filepath.Base(o.File)
		}
	}
	return ""
}

// folderConfigSource is the path of a folder config, relative to the project root.
func folderConfigSource(cfg *config.AppConfig, o config.FolderConfig) string {
	if rel, err := filepath.Rel(cfg.ProjectRoot, o.File); err == nil {
		return filepath.ToSlash(rel)
	}
	return o.Path + "/" + filepath.Base(o.File)
}

// configSourceFor names the config of the nearest folder with its own
// settings above path, or "" when only the root config applies.
func (b *HtmlReactReportBuilder) configSourceFor(path string) string {
	source, longest := "", -1
	for _, o := range b.config.Folders {
		if (path == o.Path || strings.HasPrefix(path, o.Path+"/")) && len(o.Path) > longest {
			source, longest = folderConfigSource(b.config, o), len(o.Path)
		}
	}
	return source
}

// buildFolderBands lists the folders with their own warning ranges, outer
// folders first, so the UI can re-classify their rows like the annotator did.
func (b *HtmlReactReportBuilder) buildFolderBands() []folderBands {
	var out []folderBands
	for _, o := range b.config.FolderBands {
		out = append(out, folderBands{Path: o.Path, Bands: buildStatusBands(o.Bands)})
	}
	return out
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

	// the config file of the folder, as a row that only shows it is there
	if name := b.configFileIn(dir.Path); name != "" {
		path := name
		if parentID != "" {
			path = dir.Path + "/" + name
		}
		*out = append(*out, fileNode{ID: path, Name: name, Type: "file", Path: path, ParentID: parentID, Depth: depth, Config: true})
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
	return totals{Metrics: b.buildMetricsMap(tree.Metrics), Files: files, Folders: folders}
}

// buildMetricsMap shapes the calculated file metrics for the UI. A metric
// needs no code here: a percentage and a plain value each have one shape.
func (b *HtmlReactReportBuilder) buildMetricsMap(m model.CoverageMetrics) metricsMap {
	metrics := metricsMap{}
	for key := range b.config.ActiveFileMetrics {
		switch detail := m.Calculated[key].(type) {
		case model.ScoreDetail:
			metrics[string(key)] = scoreDetail{Value: detail.Value}
		case model.CoverageDetail:
			if !countsCode[key] {
				metrics[string(key)] = countDetail{Covered: detail.Covered, Total: detail.Total, Percentage: detail.Percentage}
				continue
			}
			d := lineCoverageDetail{Covered: detail.Covered, Uncovered: detail.Uncovered, Coverable: detail.Total, Total: detail.Total, Percentage: detail.Percentage}
			// lines also have a total that includes the lines nothing can cover
			switch key {
			case config.LineCoverage:
				d.Total = m.TotalLines
			case config.PatchLineCoverage:
				if m.PatchLinesTotal > 0 {
					d.Total = m.PatchLinesTotal
				}
			}
			metrics[string(key)] = d
		}
	}
	return metrics
}

// metrics that count lines or statements; the others count methods
var countsCode = map[config.MetricKey]bool{
	config.LineCoverage:           true,
	config.StatementCoverage:      true,
	config.PatchLineCoverage:      true,
	config.PatchStatementCoverage: true,
}

// buildMetricDefinitions tells the UI how to label and lay out each active
// metric. The text comes from the metric tables in internal/config.
func (b *HtmlReactReportBuilder) buildMetricDefinitions() metricDefinitions {
	defs := metricDefinitions{}
	for _, def := range config.FileMetricDefs {
		if b.config.ActiveFileMetrics[def.Key] {
			defs[string(def.Key)] = fileMetricDefinition(def)
		}
	}
	for key, def := range b.methodMetricKeys() {
		defs[key] = metricDefinition{
			Label:       def.Label,
			ShortLabel:  def.ShortLabel(),
			Description: def.Doc,
			Kind:        metricKind(def),
			SubMetrics:  []subMetric{{ID: "value", Label: "Value", Width: 100}},
		}
	}
	return defs
}

func metricKind(def config.MetricDef) string {
	if def.Value {
		return "value"
	}
	return ""
}

func fileMetricDefinition(def config.MetricDef) metricDefinition {
	out := metricDefinition{Label: def.Label, ShortLabel: def.ShortLabel(), Description: def.Doc, Kind: metricKind(def)}
	percentage := subMetric{ID: "percentage", Label: "Percentage %", Width: 160}
	switch {
	case def.Value:
		// wide enough that a header like "Max Complexity" never wraps
		out.SubMetrics = []subMetric{{ID: "value", Label: "Value", Width: 140}}
	case def.Key == config.LineCoverage || def.Key == config.PatchLineCoverage:
		out.SubMetrics = []subMetric{
			{ID: "covered", Label: "Covered", Width: 100},
			{ID: "uncovered", Label: "Uncovered", Width: 100},
			{ID: "coverable", Label: "Coverable", Width: 100},
			{ID: "total", Label: "Total", Width: 80},
			percentage,
		}
	case countsCode[def.Key]:
		out.SubMetrics = []subMetric{
			{ID: "covered", Label: "Covered", Width: 100},
			{ID: "uncovered", Label: "Uncovered", Width: 100},
			{ID: "total", Label: "Total", Width: 80},
			percentage,
		}
	default:
		covered := "Hit"
		if strings.Contains(def.Name, "fully_covered") {
			covered = "Covered"
		}
		out.SubMetrics = []subMetric{{ID: "covered", Label: covered, Width: 80}, {ID: "total", Label: "Total", Width: 80}, percentage}
	}
	return out
}

// methodMetricKeys gives each active method metric its key in the UI data.
// The UI lists method metrics sorted by key, so the key starts with a letter
// for the configured position: "a_statement_coverage", "b_complexity", ...
func (b *HtmlReactReportBuilder) methodMetricKeys() map[string]config.MetricDef {
	keys := make(map[string]config.MetricDef, len(b.config.MethodMetrics))
	for i, key := range b.config.MethodMetrics {
		if def, ok := config.Metric(key); ok {
			keys[methodUIKey(i, def)] = def
		}
	}
	return keys
}

func methodUIKey(position int, def config.MetricDef) string {
	return string(rune('a'+position)) + "_" + def.Name
}

func countFlatNodes(nodes []fileNode) (files, folders int) {
	for _, node := range nodes {
		switch {
		case node.Config:
		case node.Type == "file":
			files++
		default:
			folders++
		}
	}
	return
}
