package review_test

import (
	"testing"

	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/review"
	"github.com/stretchr/testify/assert"
)

func intPtr(v int) *int           { return &v }
func floatPtr(v float64) *float64 { return &v }

// one changed file, three methods: untested+added (CC 12), covered+modified
// (CC 4), and an unchanged one that must be ignored.
func buildTree() *model.SummaryTree {
	file := &model.FileNode{
		Name: "svc.go",
		Path: "pkg/svc.go",
		Diff: &model.DiffInfo{Kind: model.ChangeKindAdded},
		Methods: []model.MethodMetrics{
			{
				Name: "Untested", StartLine: 10, EndLine: 30,
				DiffStatus:           "added",
				CyclomaticComplexity: intPtr(12),
				StatementsValid:      10, StatementsCovered: 0,
				PatchStatementsValid: 10, PatchStatementsCovered: 0,
			},
			{
				Name: "Covered", StartLine: 40, EndLine: 60,
				DiffStatus:           "modified",
				CyclomaticComplexity: intPtr(4),
				StatementsValid:      8, StatementsCovered: 8,
				PatchStatementsValid: 4, PatchStatementsCovered: 4,
			},
			{
				Name: "Unchanged", StartLine: 70, EndLine: 90,
				CyclomaticComplexity: intPtr(30),
				StatementsValid:      5, StatementsCovered: 0,
			},
		},
	}

	root := &model.DirNode{
		Name:  "root",
		Files: map[string]*model.FileNode{"svc.go": file},
	}

	return &model.SummaryTree{
		Root: root,
		Metrics: model.CoverageMetrics{
			PatchStatementsValid:   14,
			PatchStatementsCovered: 4,
		},
	}
}

func TestEvaluate_StatsAndHotspots(t *testing.T) {
	cfg := &config.AppConfig{}
	res := review.Evaluate(buildTree(), cfg)

	assert.True(t, res.Passed, "no gate configured means nothing can fail")
	assert.Empty(t, res.Checks)

	assert.Equal(t, 1, res.Stats.ChangedFiles)
	assert.Equal(t, 1, res.Stats.MethodsAdded)
	assert.Equal(t, 1, res.Stats.MethodsModified)
	assert.Equal(t, 1, res.Stats.UntestedChangedMethods)
	assert.Equal(t, 12, res.Stats.MaxChangedComplexity, "unchanged method's CC 30 must not count")
	assert.Equal(t, 14, res.Stats.PatchStatementsValid)
	assert.Equal(t, 4, res.Stats.PatchStatementsCovered)

	// Only the two changed methods are hotspots, riskiest first.
	assert.Len(t, res.Hotspots, 2)
	assert.Equal(t, "Untested", res.Hotspots[0].Method)
	assert.InDelta(t, 12.0, res.Hotspots[0].Risk, 0.001) // CC 12 * fully uncovered
	assert.Equal(t, "Covered", res.Hotspots[1].Method)
	assert.InDelta(t, 0.0, res.Hotspots[1].Risk, 0.001)
	if assert.NotNil(t, res.Hotspots[0].PatchCoverage) {
		assert.InDelta(t, 0.0, *res.Hotspots[0].PatchCoverage, 0.001)
	}
}

func TestEvaluate_GateFailsAndPasses(t *testing.T) {
	failing := &config.AppConfig{
		Review: config.ReviewConfig{
			Gate: config.ReviewGate{
				PatchStatementCoverage:     floatPtr(80),
				MaxChangedMethodComplexity: intPtr(10),
			},
		},
	}
	res := review.Evaluate(buildTree(), failing)
	assert.False(t, res.Passed)
	assert.Len(t, res.Checks, 2)
	for _, c := range res.Checks {
		assert.False(t, c.Passed, c.Key)
	}
	// 4/14 covered.
	assert.InDelta(t, 28.571, res.Checks[0].Value, 0.01)

	passing := &config.AppConfig{
		Review: config.ReviewConfig{
			Gate: config.ReviewGate{
				PatchStatementCoverage:     floatPtr(20),
				MaxChangedMethodComplexity: intPtr(15),
			},
		},
	}
	res = review.Evaluate(buildTree(), passing)
	assert.True(t, res.Passed)
}

func TestEvaluate_HotspotLimit(t *testing.T) {
	cfg := &config.AppConfig{Review: config.ReviewConfig{Hotspots: 1}}
	res := review.Evaluate(buildTree(), cfg)
	assert.Len(t, res.Hotspots, 1)
	assert.Equal(t, "Untested", res.Hotspots[0].Method)
}

func TestEvaluate_NilTree(t *testing.T) {
	res := review.Evaluate(nil, &config.AppConfig{})
	assert.True(t, res.Passed)
	assert.Equal(t, review.Stats{}, res.Stats)
}
