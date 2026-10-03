package htmlreact

import (
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/IgorBayerl/nanovision/internal/calculator"
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHtmlReactReportBuilder_ActiveFileMetricsFilter(t *testing.T) {
	tempDir := t.TempDir()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg := &config.AppConfig{
		ActiveFileMetrics: map[config.MetricKey]bool{
			"line_coverage": true,
			// implicitly missing statement_coverage
		},
	}

	builder := NewHtmlReactReportBuilder(tempDir, logger, false, cfg)

	tree := &model.SummaryTree{
		Metrics: model.CoverageMetrics{
			StatementsValid:   10,
			StatementsCovered: 5,
			LinesValid:        10,
			LinesCovered:      5,
		},
		Root: &model.DirNode{
			Name: "root",
			Metrics: model.CoverageMetrics{
				StatementsValid:   10,
				StatementsCovered: 5,
				LinesValid:        10,
				LinesCovered:      5,
			},
		},
	}
	calculator.CalculateTree(tree, cfg.ActiveFileMetrics, nil)

	b, ok := builder.(*HtmlReactReportBuilder)
	if !ok {
		t.Fatalf("Expected HtmlReactReportBuilder")
	}

	totalsData := b.buildTotals(tree, 1, 0)

	if _, ok := totalsData.Metrics[string(config.LineCoverage)]; !ok {
		t.Errorf("Expected line_coverage to be present")
	}
	if _, ok := totalsData.Metrics[string(config.StatementCoverage)]; ok {
		t.Errorf("Expected statement_coverage to be absent: it is not an active metric")
	}

	// in JSON the metrics sit next to the counts
	data, err := json.Marshal(totalsData)
	require.NoError(t, err)
	var flat map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(data, &flat))
	assert.Contains(t, flat, "line_coverage")
	assert.JSONEq(t, "1", string(flat["files"]))
}

// TestBuildMetricsMap_GoldenFile verifies that the provider-loop output is
// structurally equivalent to the pre-refactoring inline output by snapshot-testing
// the JSON for a rich fixture with all metric fields populated.
func TestBuildMetricsMap_GoldenFile(t *testing.T) {
	cfg := &config.AppConfig{
		ActiveFileMetrics: map[config.MetricKey]bool{
			config.StatementCoverage:      true,
			config.LineCoverage:           true,
			config.MethodsHit:             true,
			config.MethodsFullyCovered:    true,
			config.PatchStatementCoverage: true,
			config.PatchLineCoverage:      true,
			config.PatchMethodsHit:        true,
		},
	}
	b := &HtmlReactReportBuilder{config: cfg}

	m := model.CoverageMetrics{
		StatementsCovered:      7,
		StatementsValid:        10,
		LinesCovered:           80,
		LinesValid:             100,
		TotalLines:             150,
		MethodsHit:             4,
		MethodsFullyCovered:    2,
		MethodsValid:           6,
		PatchStatementsCovered: 2,
		PatchStatementsValid:   3,
		PatchLinesCovered:      10,
		PatchLinesValid:        12,
		PatchLinesTotal:        20,
		PatchMethodsHit:        1,
		PatchMethodsValid:      2,
	}
	tree := &model.SummaryTree{Root: &model.DirNode{Metrics: m}}
	calculator.CalculateTree(tree, cfg.ActiveFileMetrics, nil)

	metrics := b.buildMetricsMap(tree.Root.Metrics)

	// Verify all 7 metric keys are present
	expectedKeys := []string{
		string(config.StatementCoverage),
		string(config.LineCoverage),
		string(config.MethodsHit),
		string(config.MethodsFullyCovered),
		string(config.PatchStatementCoverage),
		string(config.PatchLineCoverage),
		string(config.PatchMethodsHit),
	}
	for _, key := range expectedKeys {
		assert.Contains(t, metrics, key, "expected key %q in metrics map", key)
	}

	// Snapshot the JSON to verify structural shape
	jsonBytes, err := json.Marshal(metrics)
	require.NoError(t, err)

	var roundTrip map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(jsonBytes, &roundTrip))

	// Verify statement_coverage detail shape
	var sc lineCoverageDetail
	require.NoError(t, json.Unmarshal(roundTrip[string(config.StatementCoverage)], &sc))
	assert.Equal(t, 7, sc.Covered)
	assert.Equal(t, 3, sc.Uncovered)
	assert.Equal(t, 10, sc.Coverable)
	assert.Equal(t, 10, sc.Total)
	assert.Equal(t, 70.0, sc.Percentage)

	// Verify line_coverage detail shape
	var lc lineCoverageDetail
	require.NoError(t, json.Unmarshal(roundTrip[string(config.LineCoverage)], &lc))
	assert.Equal(t, 80, lc.Covered)
	assert.Equal(t, 20, lc.Uncovered)
	assert.Equal(t, 100, lc.Coverable)
	assert.Equal(t, 150, lc.Total)
	assert.Equal(t, 80.0, lc.Percentage)

	// Verify methods_hit detail shape
	var mh countDetail
	require.NoError(t, json.Unmarshal(roundTrip[string(config.MethodsHit)], &mh))
	assert.Equal(t, 4, mh.Covered)
	assert.Equal(t, 6, mh.Total)

	// Verify methods_fully_covered detail shape
	var mfc countDetail
	require.NoError(t, json.Unmarshal(roundTrip[string(config.MethodsFullyCovered)], &mfc))
	assert.Equal(t, 2, mfc.Covered)
	assert.Equal(t, 6, mfc.Total)

	// Verify patch_statement_coverage detail shape
	var psc lineCoverageDetail
	require.NoError(t, json.Unmarshal(roundTrip[string(config.PatchStatementCoverage)], &psc))
	assert.Equal(t, 2, psc.Covered)
	assert.Equal(t, 1, psc.Uncovered)
	assert.Equal(t, 3, psc.Total)

	// Verify patch_line_coverage detail shape
	var plc lineCoverageDetail
	require.NoError(t, json.Unmarshal(roundTrip[string(config.PatchLineCoverage)], &plc))
	assert.Equal(t, 10, plc.Covered)
	assert.Equal(t, 2, plc.Uncovered)
	assert.Equal(t, 12, plc.Coverable)
	assert.Equal(t, 20, plc.Total)

	// Verify patch_methods_hit detail shape
	var pmh countDetail
	require.NoError(t, json.Unmarshal(roundTrip[string(config.PatchMethodsHit)], &pmh))
	assert.Equal(t, 1, pmh.Covered)
	assert.Equal(t, 2, pmh.Total)
}

// TestBuildMetricsMap_GuardsSkipEmptyData confirms that providers with
// guards (e.g. StatementsValid > 0) skip writing when there's no data.
func TestBuildMetricsMap_GuardsSkipEmptyData(t *testing.T) {
	cfg := &config.AppConfig{
		ActiveFileMetrics: map[config.MetricKey]bool{
			config.StatementCoverage: true,
			config.MethodsHit:        true,
		},
	}
	b := &HtmlReactReportBuilder{config: cfg}

	// All zero values — guards should prevent any entries
	m := model.CoverageMetrics{}
	tree := &model.SummaryTree{Root: &model.DirNode{Metrics: m}}
	calculator.CalculateTree(tree, cfg.ActiveFileMetrics, nil)

	metrics := b.buildMetricsMap(tree.Root.Metrics)

	_, hasStmt := metrics[string(config.StatementCoverage)]
	assert.False(t, hasStmt, "statement_coverage should be skipped when StatementsValid == 0")

	_, hasMethods := metrics[string(config.MethodsHit)]
	assert.False(t, hasMethods, "methods_hit should be skipped when MethodsValid == 0")
}

// Every metric the UI can show needs a tooltip description, so a new metric
// cannot ship without one.
func TestBuildMetricDefinitions_EveryMetricIsDescribed(t *testing.T) {
	cfg := &config.AppConfig{
		FileMetrics:         config.DefaultFileMetrics,
		MethodMetrics:       config.DefaultMethodMetrics,
		ActiveFileMetrics:   map[config.MetricKey]bool{},
		ActiveMethodMetrics: map[config.MetricKey]bool{},
	}
	for _, k := range cfg.FileMetrics {
		cfg.ActiveFileMetrics[k] = true
	}
	for _, k := range cfg.MethodMetrics {
		cfg.ActiveMethodMetrics[k] = true
	}
	b := &HtmlReactReportBuilder{config: cfg}

	defs := b.buildMetricDefinitions()
	assert.Len(t, defs, len(config.FileMetricDefs)+len(config.MethodMetricDefs), "every metric of the tables has a definition")

	for key, def := range defs {
		assert.NotEmpty(t, def.Description, "metric %q has no tooltip description", key)
		assert.NotEmpty(t, def.SubMetrics, "metric %q has no columns", key)
	}

	// The text comes from the metric tables, not a copy kept in the reporter.
	statements, _ := config.Metric(config.StatementCoverage)
	assert.Equal(t, statements.Doc, defs[string(config.StatementCoverage)].Description)
	// method metrics are keyed by their configured position
	assert.Equal(t, "Statements", defs["a_statement_coverage"].Label)
	assert.Equal(t, "value", defs["e_complexity"].Kind)
}

func TestUniqueReportLabels(t *testing.T) {
	tests := []struct {
		name  string
		paths []string
		want  []string
	}{
		{
			name:  "distinct file names need no directories",
			paths: []string{"reports/coverage-unit.out", "reports/coverage-integration.out", "cobertura.xml"},
			want:  []string{"coverage-unit.out", "coverage-integration.out", "cobertura.xml"},
		},
		{
			name:  "shared file names grow until they differ",
			paths: []string{"shard1/coverage.out", "shard2/coverage.out", "reports/cobertura.xml"},
			want:  []string{"shard1/coverage.out", "shard2/coverage.out", "cobertura.xml"},
		},
		{
			name:  "only the ambiguous ones grow",
			paths: []string{"a/out/coverage.out", "b/out/coverage.out", "lcov.info"},
			want:  []string{"a/out/coverage.out", "b/out/coverage.out", "lcov.info"},
		},
		{
			name:  "windows separators are treated as paths",
			paths: []string{`C:\ci\unit\coverage.out`, `C:\ci\e2e\coverage.out`},
			want:  []string{"unit/coverage.out", "e2e/coverage.out"},
		},
		{
			name:  "identical paths stop once fully spelled out",
			paths: []string{"reports/coverage.out", "reports/coverage.out"},
			want:  []string{"reports/coverage.out", "reports/coverage.out"},
		},
		{
			name:  "a bare file name survives",
			paths: []string{"coverage.out"},
			want:  []string{"coverage.out"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, uniqueReportLabels(tc.paths))
		})
	}
}

func TestMetricOrder_FollowsConfiguredFileMetrics(t *testing.T) {
	cfg := &config.AppConfig{
		FileMetrics: []config.MetricKey{config.StatementCoverage, config.PatchMethodsHit, config.MaxCyclomaticComplexity, config.MethodsHit},
	}
	b := NewHtmlReactReportBuilder(t.TempDir(), slog.New(slog.NewTextHandler(os.Stdout, nil)), false, cfg).(*HtmlReactReportBuilder)

	assert.Equal(t, []string{"statement_coverage", "patch_methods_hit", "max_cyclomatic_complexity", "methods_hit"}, b.metricOrder())
}

func TestComparingItems(t *testing.T) {
	hash := "d26739fc86ec8b84bdf7b192e92288b8d616e583"
	tree := &model.SummaryTree{
		Versions: &model.Versions{
			Base:       model.RunRef{Revision: hash, Stream: "main"},
			Current:    model.RunRef{Revision: "//game/dev@120", Stream: "//game/dev"},
			LocalEdits: true,
			Diff:       "change.diff",
		},
		Comparison: &model.Comparison{Base: model.RunRef{Revision: hash}, Distance: 3},
		Change:     &model.ChangeSet{Measured: []string{"a.go"}, Ignored: []string{"a_test.go"}},
	}
	assert.Equal(t, []MetadataItem{
		{Label: "Base", Value: "main d26739fc"},
		{Label: "Current", Value: "//game/dev@120 + edits"},
		{Label: "Coverage base", Value: "d26739fc, 3 before base"},
		{Label: "Changed files", Value: "2 from change.diff"},
	}, comparingItems(tree))

	assert.Empty(t, comparingItems(&model.SummaryTree{}), "a report without a diff or revisions has no section")
}

func TestBuildConfigs(t *testing.T) {
	band := &config.Band{Min: 80, Max: 90}
	root := t.TempDir()
	b := &HtmlReactReportBuilder{config: &config.AppConfig{
		ProjectRoot: root,
		ConfigFile:  filepath.Join(root, "nanovision.yaml"),
		Folders: []config.FolderConfig{
			{Path: "cmd", File: filepath.Join(root, "cmd", "nanovision.yaml"), ScopedConfig: config.ScopedConfig{IgnoreFiles: []string{"gen/**"}}},
			{Path: "internal/store", File: filepath.Join(root, "internal", "store", "nanovision.yaml"), ScopedConfig: config.ScopedConfig{
				Metrics: config.MetricsConfig{Files: []config.MetricSetting{{Name: "statement_coverage", Warning: band}}},
				Reports: []config.ReportInput{{Path: "store.out", Name: "store"}},
			}},
		},
	}}

	assert.Equal(t, []configFile{
		{Source: "nanovision.yaml"},
		{Path: "cmd", Source: "cmd/nanovision.yaml"},
		{Path: "internal/store", Source: "internal/store/nanovision.yaml"},
	}, b.buildConfigs())

	assert.Equal(t, "internal/store/nanovision.yaml", b.configSourceFor("internal/store/blob/manifest.go"))
	assert.Equal(t, "", b.configSourceFor("internal/storekeeper/a.go"), "a folder name that only starts the same is another folder")

	// a config file is a row of its folder
	assert.Equal(t, "nanovision.yaml", b.configFileIn("."))
	assert.Equal(t, "nanovision.yaml", b.configFileIn("internal/store"))
	assert.Equal(t, "", b.configFileIn("internal"))
	files, folders := countFlatNodes([]fileNode{{Type: "folder"}, {Type: "file"}, {Type: "file", Config: true}})
	assert.Equal(t, []int{1, 1}, []int{files, folders}, "a config row is not a measured file")

	b.config.Folders = nil
	assert.Nil(t, b.buildConfigs(), "a run without folder settings lists no config files")
}
