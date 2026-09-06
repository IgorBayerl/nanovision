package htmlreact

import (
	"encoding/json"
	"log/slog"
	"os"
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
			// implicitly missing statement_coverage, branch_coverage
		},
	}

	builder := NewHtmlReactReportBuilder(tempDir, logger, false, cfg)

	tree := &model.SummaryTree{
		Metrics: model.CoverageMetrics{
			StatementsValid:   10,
			StatementsCovered: 5,
			LinesValid:        10,
			LinesCovered:      5,
			BranchesValid:     10,
			BranchesCovered:   5,
		},
		Root: &model.DirNode{
			Name: "root",
			Metrics: model.CoverageMetrics{
				StatementsValid:   10,
				StatementsCovered: 5,
				LinesValid:        10,
				LinesCovered:      5,
				BranchesValid:     10,
				BranchesCovered:   5,
			},
		},
	}
	calculator.CalculateTree(tree, cfg.ActiveFileMetrics, nil)

	b, ok := builder.(*HtmlReactReportBuilder)
	if !ok {
		t.Fatalf("Expected HtmlReactReportBuilder")
	}

	totalsData := b.buildTotals(tree, 1, 0)

	if totalsData.LineCoverage == nil {
		t.Errorf("Expected LineCoverage to be present, got nil")
	}
	if totalsData.StatementCoverage != nil {
		t.Errorf("Expected StatementCoverage to be nil, got %+v", totalsData.StatementCoverage)
	}
	if totalsData.BranchCoverage != nil {
		t.Errorf("Expected BranchCoverage to be nil, got %+v", totalsData.BranchCoverage)
	}
}

// TestBuildMetricsMap_GoldenFile verifies that the provider-loop output is
// structurally equivalent to the pre-refactoring inline output by snapshot-testing
// the JSON for a rich fixture with all metric fields populated.
func TestBuildMetricsMap_GoldenFile(t *testing.T) {
	cfg := &config.AppConfig{
		ActiveFileMetrics: map[config.MetricKey]bool{
			config.StatementCoverage:      true,
			config.LineCoverage:           true,
			config.BranchCoverage:         true,
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
		BranchesCovered:        3,
		BranchesValid:          5,
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

	// Verify all 8 metric keys are present
	expectedKeys := []string{
		string(config.StatementCoverage),
		string(config.LineCoverage),
		string(config.BranchCoverage),
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

	// Verify branch_coverage detail shape
	var bc branchCoverageDetail
	require.NoError(t, json.Unmarshal(roundTrip[string(config.BranchCoverage)], &bc))
	assert.Equal(t, 3, bc.Covered)
	assert.Equal(t, 5, bc.Total)
	assert.Equal(t, 60.0, bc.Percentage)

	// Verify methods_hit detail shape
	var mh methodsHitDetail
	require.NoError(t, json.Unmarshal(roundTrip[string(config.MethodsHit)], &mh))
	assert.Equal(t, 4, mh.Covered)
	assert.Equal(t, 6, mh.Total)

	// Verify methods_fully_covered detail shape
	var mfc methodsFullyCoveredDetail
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
	var pmh methodsHitDetail
	require.NoError(t, json.Unmarshal(roundTrip[string(config.PatchMethodsHit)], &pmh))
	assert.Equal(t, 1, pmh.Covered)
	assert.Equal(t, 2, pmh.Total)
}

// TestBuildMetricsMap_RemoveBranchCoverage confirms that removing branch_coverage
// from ActiveFileMetrics removes branch data from the rendered file detail view.
func TestBuildMetricsMap_RemoveBranchCoverage(t *testing.T) {
	cfg := &config.AppConfig{
		ActiveFileMetrics: map[config.MetricKey]bool{
			config.StatementCoverage: true,
			config.LineCoverage:      true,
			// branch_coverage intentionally omitted
		},
	}
	b := &HtmlReactReportBuilder{config: cfg}

	m := model.CoverageMetrics{
		StatementsCovered: 5,
		StatementsValid:   10,
		LinesCovered:      50,
		LinesValid:        100,
		TotalLines:        100,
		BranchesCovered:   3,
		BranchesValid:     5,
	}
	tree := &model.SummaryTree{Metrics: m, Root: &model.DirNode{Metrics: m}}
	calculator.CalculateTree(tree, cfg.ActiveFileMetrics, nil)

	metrics := b.buildMetricsMap(tree.Root.Metrics)

	// Branch data must be absent
	_, hasBranch := metrics[string(config.BranchCoverage)]
	assert.False(t, hasBranch, "branch_coverage should be absent when not in ActiveFileMetrics")

	// Statement and line data should be present
	assert.Contains(t, metrics, string(config.StatementCoverage))
	assert.Contains(t, metrics, string(config.LineCoverage))

	// Also verify totals mirrors the same behavior
	// Tree is already calculated above.
	totalsData := b.buildTotals(tree, 1, 0)
	assert.Nil(t, totalsData.BranchCoverage, "totals.BranchCoverage should be nil")
	assert.NotNil(t, totalsData.StatementCoverage, "totals.StatementCoverage should be present")
	assert.NotNil(t, totalsData.LineCoverage, "totals.LineCoverage should be present")
}

// TestBuildMetricsMap_GuardsSkipEmptyData confirms that providers with
// guards (e.g. StatementsValid > 0) skip writing when there's no data.
func TestBuildMetricsMap_GuardsSkipEmptyData(t *testing.T) {
	cfg := &config.AppConfig{
		ActiveFileMetrics: map[config.MetricKey]bool{
			config.StatementCoverage: true,
			config.BranchCoverage:    true,
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

	_, hasBranch := metrics[string(config.BranchCoverage)]
	assert.False(t, hasBranch, "branch_coverage should be skipped when BranchesValid == 0")

	_, hasMethods := metrics[string(config.MethodsHit)]
	assert.False(t, hasMethods, "methods_hit should be skipped when MethodsValid == 0")
}

// Every metric the UI can show needs a tooltip description, so a new metric
// cannot ship without one.
func TestBuildMetricDefinitions_EveryMetricIsDescribed(t *testing.T) {
	b := &HtmlReactReportBuilder{config: &config.AppConfig{
		ActiveFileMetrics: map[config.MetricKey]bool{
			config.StatementCoverage:       true,
			config.LineCoverage:            true,
			config.BranchCoverage:          true,
			config.MethodsHit:              true,
			config.MethodsFullyCovered:     true,
			config.PatchStatementCoverage:  true,
			config.PatchLineCoverage:       true,
			config.PatchMethodsHit:         true,
			config.MaxCyclomaticComplexity: true,
		},
		ActiveMethodMetrics: map[config.MetricKey]bool{
			config.MethodStatementCoverage:      true,
			config.MethodLineCoverage:           true,
			config.MethodPatchStatementCoverage: true,
			config.MethodPatchLineCoverage:      true,
			config.MethodBranchCoverage:         true,
			config.CyclomaticComplexity:         true,
			config.MethodCrapScore:              true,
			config.MethodPatchCrapScore:         true,
			config.MethodExposedRisk:            true,
			config.MethodDefectProbability:      true,
		},
	}}

	defs := b.buildMetricDefinitions()
	assert.NotEmpty(t, defs)

	for key, def := range defs {
		assert.NotEmpty(t, def.Description, "metric %q has no tooltip description", key)
	}

	// Descriptions come from the evaluators, not a copy kept in the reporter.
	assert.Equal(t, "Percentage of covered code branches.", defs[string(config.BranchCoverage)].Description)
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
