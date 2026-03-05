package status_test

import (
	"testing"

	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/status"
	"github.com/IgorBayerl/nanovision/internal/status/evaluators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAnnotate_NilTree ensures Annotate handles nil inputs gracefully.
func TestAnnotate_NilTree(t *testing.T) {
	cfg := &config.AppConfig{}
	assert.NotPanics(t, func() {
		status.Annotate(nil, cfg, status.Capabilities{}, evaluators.Registry)
	})
	assert.NotPanics(t, func() {
		status.Annotate(&model.SummaryTree{Root: nil}, cfg, status.Capabilities{}, evaluators.Registry)
	})
}

// TestAnnotate_Integration builds a minimal in-memory tree and verifies that
// the new annotation engine produces correct statuses on file nodes and
// method nodes.
func TestAnnotate_Integration(t *testing.T) {
	complexity := 15

	fileNode := &model.FileNode{
		Name: "main.go",
		Path: "src/main.go",
		Metrics: model.CoverageMetrics{
			LinesCovered:            90,
			LinesValid:              100,
			BranchesCovered:         40,
			BranchesValid:           100,
			StatementsCovered:       70,
			StatementsValid:         100,
			MaxCyclomaticComplexity: 15,
		},
		Methods: []model.MethodMetrics{
			{
				Name:                 "DoWork",
				StartLine:            10,
				EndLine:              50,
				CyclomaticComplexity: &complexity,
				LinesValid:           20,
				LinesCovered:         18,
				StatementsValid:      10,
				StatementsCovered:    8,
			},
		},
	}

	rootDir := &model.DirNode{
		Name:    "root",
		Path:    "/",
		Subdirs: map[string]*model.DirNode{},
		Files:   map[string]*model.FileNode{"main.go": fileNode},
		Metrics: model.CoverageMetrics{
			LinesCovered:            90,
			LinesValid:              100,
			BranchesCovered:         40,
			BranchesValid:           100,
			StatementsCovered:       70,
			StatementsValid:         100,
			MaxCyclomaticComplexity: 15,
		},
	}

	tree := &model.SummaryTree{Root: rootDir}

	cfg := &config.AppConfig{
		StatusBands: config.StatusBands{
			config.LineCoverage:            config.Band{Min: 60, Max: 80},
			config.BranchCoverage:          config.Band{Min: 60, Max: 80},
			config.StatementCoverage:       config.Band{Min: 60, Max: 80},
			config.MaxCyclomaticComplexity: config.Band{Min: 5, Max: 10},
		},
		ActiveFileMetrics: map[config.MetricKey]bool{
			config.LineCoverage:            true,
			config.BranchCoverage:          true,
			config.StatementCoverage:       true,
			config.MaxCyclomaticComplexity: true,
		},
		ActiveMethodMetrics: map[config.MetricKey]bool{
			config.LineCoverage:            true,
			config.StatementCoverage:       true,
			config.MaxCyclomaticComplexity: true,
		},
	}

	caps := status.Capabilities{
		HasBranchCoverage:    true,
		HasMethodCoverage:    true,
		HasStatementCoverage: true,
	}

	status.Annotate(tree, cfg, caps, evaluators.Registry)

	// ---------- Verify file node statuses ----------
	require.NotNil(t, fileNode.Statuses)
	// 90% line coverage > 80 max → safe
	assert.Equal(t, "safe", fileNode.Statuses[config.LineCoverage])
	// 40% branch coverage < 60 min → danger
	assert.Equal(t, "danger", fileNode.Statuses[config.BranchCoverage])
	// 70% statement coverage in [60,80] → warning
	assert.Equal(t, "warning", fileNode.Statuses[config.StatementCoverage])
	// complexity 15 > 10 max → danger (lower is better)
	assert.Equal(t, "danger", fileNode.Statuses[config.MaxCyclomaticComplexity])

	// ---------- Verify dir node statuses ----------
	require.NotNil(t, rootDir.Statuses)
	assert.Equal(t, "safe", rootDir.Statuses[config.LineCoverage])
	assert.Equal(t, "danger", rootDir.Statuses[config.BranchCoverage])

	// ---------- Verify method node statuses ----------
	require.Len(t, fileNode.Methods, 1)
	method := fileNode.Methods[0]
	require.NotNil(t, method.Statuses)
	// 18/20 = 90% line coverage → safe
	assert.Equal(t, "safe", method.Statuses[config.LineCoverage])
	// 8/10 = 80% statement coverage → warning (at max)
	assert.Equal(t, "warning", method.Statuses[config.StatementCoverage])
	// complexity 15 > 10 → danger
	assert.Equal(t, "danger", method.Statuses[config.MaxCyclomaticComplexity])
}

// TestAnnotate_BranchNotApplicable ensures branch coverage is skipped when
// HasBranchCoverage is false.
func TestAnnotate_BranchNotApplicable(t *testing.T) {
	fileNode := &model.FileNode{
		Name: "main.go",
		Path: "src/main.go",
		Metrics: model.CoverageMetrics{
			LinesCovered:    90,
			LinesValid:      100,
			BranchesCovered: 40,
			BranchesValid:   100,
		},
	}

	rootDir := &model.DirNode{
		Name:    "root",
		Path:    "/",
		Subdirs: map[string]*model.DirNode{},
		Files:   map[string]*model.FileNode{"main.go": fileNode},
		Metrics: model.CoverageMetrics{
			LinesCovered:    90,
			LinesValid:      100,
			BranchesCovered: 40,
			BranchesValid:   100,
		},
	}

	tree := &model.SummaryTree{Root: rootDir}

	cfg := &config.AppConfig{
		StatusBands: config.StatusBands{
			config.LineCoverage:   config.Band{Min: 60, Max: 80},
			config.BranchCoverage: config.Band{Min: 60, Max: 80},
		},
		ActiveFileMetrics: map[config.MetricKey]bool{
			config.LineCoverage:   true,
			config.BranchCoverage: true,
		},
		ActiveMethodMetrics: map[config.MetricKey]bool{},
	}

	// HasBranchCoverage = false → branch coverage should not appear
	caps := status.Capabilities{HasBranchCoverage: false}

	status.Annotate(tree, cfg, caps, evaluators.Registry)

	assert.Equal(t, "safe", fileNode.Statuses[config.LineCoverage])
	_, hasBranch := fileNode.Statuses[config.BranchCoverage]
	assert.False(t, hasBranch, "branch_coverage status should not be set when HasBranchCoverage is false")
}

// TestAnnotate_MethodsHitViaRegistry verifies that metrics previously handled
// by the old fallback spec table (e.g. methods_hit) are now correctly
// evaluated through their dedicated Evaluator in the registry.
func TestAnnotate_MethodsHitViaRegistry(t *testing.T) {
	fileNode := &model.FileNode{
		Name: "main.go",
		Path: "src/main.go",
		Metrics: model.CoverageMetrics{
			LinesCovered: 90,
			LinesValid:   100,
			MethodsHit:   5,
			MethodsValid: 10,
		},
	}

	rootDir := &model.DirNode{
		Name:    "root",
		Path:    "/",
		Subdirs: map[string]*model.DirNode{},
		Files:   map[string]*model.FileNode{"main.go": fileNode},
		Metrics: model.CoverageMetrics{
			LinesCovered: 90,
			LinesValid:   100,
			MethodsHit:   5,
			MethodsValid: 10,
		},
	}

	tree := &model.SummaryTree{Root: rootDir}

	cfg := &config.AppConfig{
		StatusBands: config.StatusBands{
			config.MethodsHit: config.Band{Min: 60, Max: 80},
		},
		ActiveFileMetrics: map[config.MetricKey]bool{
			config.MethodsHit: true,
		},
		ActiveMethodMetrics: map[config.MetricKey]bool{},
	}

	caps := status.Capabilities{HasMethodCoverage: true}

	status.Annotate(tree, cfg, caps, evaluators.Registry)

	// 5/10 = 50% < 60 min → danger
	assert.Equal(t, "danger", fileNode.Statuses[config.MethodsHit])
}
