package evaluators

import (
	"testing"

	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/status"
	"github.com/IgorBayerl/nanovision/internal/utils"
	"github.com/stretchr/testify/assert"
)

// ---------------------------------------------------------------------------
// LineCoverageEvaluator
// ---------------------------------------------------------------------------

func TestLineCoverageEvaluator_Key(t *testing.T) {
	assert.Equal(t, config.LineCoverage, LineCoverageEvaluator{}.Key())
}

func TestLineCoverageEvaluator_IsApplicable(t *testing.T) {
	assert.True(t, LineCoverageEvaluator{}.IsApplicable(status.Capabilities{}))
	assert.True(t, LineCoverageEvaluator{}.IsApplicable(status.Capabilities{HasBranchCoverage: true}))
}

func TestLineCoverageEvaluator_ZeroGuard(t *testing.T) {
	m := model.CoverageMetrics{LinesValid: 0}
	band := &config.Band{Min: 60, Max: 80}
	lvl, show := LineCoverageEvaluator{}.Evaluate(m, band)
	assert.Equal(t, status.RiskLevel(""), lvl)
	assert.False(t, show)
}

func TestLineCoverageEvaluator_Thresholds(t *testing.T) {
	band := &config.Band{Min: 60, Max: 80}
	tests := []struct {
		name    string
		covered int
		valid   int
		want    status.RiskLevel
	}{
		{"danger - below min", 50, 100, status.RiskDanger},
		{"warning - in range", 70, 100, status.RiskWarning},
		{"safe - above max", 90, 100, status.RiskSafe},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := model.CoverageMetrics{LinesCovered: tt.covered, LinesValid: tt.valid}
			m.Calculated = map[config.MetricKey]any{config.LineCoverage: model.CoverageDetail{Percentage: utils.CalculatePercentage(tt.covered, tt.valid, 2)}}
			lvl, show := LineCoverageEvaluator{}.Evaluate(m, band)
			assert.True(t, show)
			assert.Equal(t, tt.want, lvl)
		})
	}
}

// ---------------------------------------------------------------------------
// BranchCoverageEvaluator
// ---------------------------------------------------------------------------

func TestBranchCoverageEvaluator_Key(t *testing.T) {
	assert.Equal(t, config.BranchCoverage, BranchCoverageEvaluator{}.Key())
}

func TestBranchCoverageEvaluator_IsApplicable(t *testing.T) {
	assert.False(t, BranchCoverageEvaluator{}.IsApplicable(status.Capabilities{HasBranchCoverage: false}))
	assert.True(t, BranchCoverageEvaluator{}.IsApplicable(status.Capabilities{HasBranchCoverage: true}))
}

func TestBranchCoverageEvaluator_ZeroGuard(t *testing.T) {
	m := model.CoverageMetrics{BranchesValid: 0}
	band := &config.Band{Min: 60, Max: 80}
	lvl, show := BranchCoverageEvaluator{}.Evaluate(m, band)
	assert.Equal(t, status.RiskLevel(""), lvl)
	assert.False(t, show)
}

func TestBranchCoverageEvaluator_Thresholds(t *testing.T) {
	band := &config.Band{Min: 60, Max: 80}
	tests := []struct {
		name    string
		covered int
		valid   int
		want    status.RiskLevel
	}{
		{"danger", 50, 100, status.RiskDanger},
		{"warning", 70, 100, status.RiskWarning},
		{"safe", 90, 100, status.RiskSafe},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := model.CoverageMetrics{BranchesCovered: tt.covered, BranchesValid: tt.valid}
			m.Calculated = map[config.MetricKey]any{config.BranchCoverage: model.CoverageDetail{Percentage: utils.CalculatePercentage(tt.covered, tt.valid, 2)}}
			lvl, show := BranchCoverageEvaluator{}.Evaluate(m, band)
			assert.True(t, show)
			assert.Equal(t, tt.want, lvl)
		})
	}
}

// ---------------------------------------------------------------------------
// StatementCoverageEvaluator
// ---------------------------------------------------------------------------

func TestStatementCoverageEvaluator_Key(t *testing.T) {
	assert.Equal(t, config.StatementCoverage, StatementCoverageEvaluator{}.Key())
}

func TestStatementCoverageEvaluator_IsApplicable(t *testing.T) {
	assert.True(t, StatementCoverageEvaluator{}.IsApplicable(status.Capabilities{}))
}

func TestStatementCoverageEvaluator_ZeroGuard(t *testing.T) {
	m := model.CoverageMetrics{StatementsValid: 0}
	band := &config.Band{Min: 60, Max: 80}
	lvl, show := StatementCoverageEvaluator{}.Evaluate(m, band)
	assert.Equal(t, status.RiskLevel(""), lvl)
	assert.False(t, show)
}

func TestStatementCoverageEvaluator_Thresholds(t *testing.T) {
	band := &config.Band{Min: 60, Max: 80}
	tests := []struct {
		name    string
		covered int
		valid   int
		want    status.RiskLevel
	}{
		{"danger", 50, 100, status.RiskDanger},
		{"warning", 70, 100, status.RiskWarning},
		{"safe", 90, 100, status.RiskSafe},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := model.CoverageMetrics{StatementsCovered: tt.covered, StatementsValid: tt.valid}
			m.Calculated = map[config.MetricKey]any{config.StatementCoverage: model.CoverageDetail{Percentage: utils.CalculatePercentage(tt.covered, tt.valid, 2)}}
			lvl, show := StatementCoverageEvaluator{}.Evaluate(m, band)
			assert.True(t, show)
			assert.Equal(t, tt.want, lvl)
		})
	}
}

// ---------------------------------------------------------------------------
// MaxComplexityEvaluator
// ---------------------------------------------------------------------------

func TestMaxComplexityEvaluator_Key(t *testing.T) {
	assert.Equal(t, config.MaxCyclomaticComplexity, MaxComplexityEvaluator{}.Key())
}

func TestMaxComplexityEvaluator_IsApplicable(t *testing.T) {
	assert.True(t, MaxComplexityEvaluator{}.IsApplicable(status.Capabilities{}))
}

func TestMaxComplexityEvaluator_ZeroGuard(t *testing.T) {
	m := model.CoverageMetrics{MaxCyclomaticComplexity: 0}
	band := &config.Band{Min: 5, Max: 10}
	lvl, show := MaxComplexityEvaluator{}.Evaluate(m, band)
	assert.Equal(t, status.RiskLevel(""), lvl)
	assert.False(t, show)
}

func TestMaxComplexityEvaluator_Thresholds(t *testing.T) {
	band := &config.Band{Min: 5, Max: 10}
	tests := []struct {
		name       string
		complexity int
		want       status.RiskLevel
	}{
		{"safe - below min", 3, status.RiskSafe},
		{"warning - in range", 7, status.RiskWarning},
		{"danger - above max", 15, status.RiskDanger},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := model.CoverageMetrics{MaxCyclomaticComplexity: tt.complexity}
			m.Calculated = map[config.MetricKey]any{config.MaxCyclomaticComplexity: model.ScoreDetail{Value: float64(tt.complexity)}}
			lvl, show := MaxComplexityEvaluator{}.Evaluate(m, band)
			assert.True(t, show)
			assert.Equal(t, tt.want, lvl)
		})
	}
}

// Specific assertion from Definition of Done: complexity above band.Max returns RiskDanger.
func TestMaxComplexityEvaluator_AboveMaxReturnsDanger(t *testing.T) {
	band := &config.Band{Min: 5, Max: 10}
	m := model.CoverageMetrics{MaxCyclomaticComplexity: 20}
	m.Calculated = map[config.MetricKey]any{config.MaxCyclomaticComplexity: model.ScoreDetail{Value: 20.0}}
	lvl, show := MaxComplexityEvaluator{}.Evaluate(m, band)
	assert.True(t, show)
	assert.Equal(t, status.RiskDanger, lvl)
}

// ---------------------------------------------------------------------------
// Registry
// ---------------------------------------------------------------------------

func TestRegistryContainsAllEvaluators(t *testing.T) {
	expectedKeys := []config.MetricKey{
		config.LineCoverage,
		config.BranchCoverage,
		config.StatementCoverage,
		config.MaxCyclomaticComplexity,
		config.MethodsHit,
		config.MethodsFullyCovered,
		config.PatchLineCoverage,
		config.PatchMethodsHit,
		config.PatchStatementCoverage,
		config.StatementMethodsHit,
		config.StatementMethodsFullyCovered,
		config.PatchStatementMethodsHit,
	}
	for _, key := range expectedKeys {
		ev, ok := Registry[key]
		assert.True(t, ok, "Registry should contain %s", key)
		assert.Equal(t, key, ev.Key())
	}
}

// ---------------------------------------------------------------------------
// MethodsHitEvaluator (representative of the new percentage evaluators)
// ---------------------------------------------------------------------------

func TestMethodsHitEvaluator_Key(t *testing.T) {
	assert.Equal(t, config.MethodsHit, MethodsHitEvaluator{}.Key())
}

func TestMethodsHitEvaluator_IsApplicable(t *testing.T) {
	assert.False(t, MethodsHitEvaluator{}.IsApplicable(status.Capabilities{HasMethodCoverage: false}))
	assert.True(t, MethodsHitEvaluator{}.IsApplicable(status.Capabilities{HasMethodCoverage: true}))
}

func TestMethodsHitEvaluator_ZeroGuard(t *testing.T) {
	m := model.CoverageMetrics{MethodsValid: 0}
	band := &config.Band{Min: 60, Max: 80}
	lvl, show := MethodsHitEvaluator{}.Evaluate(m, band)
	assert.Equal(t, status.RiskLevel(""), lvl)
	assert.False(t, show)
}

func TestMethodsHitEvaluator_Thresholds(t *testing.T) {
	band := &config.Band{Min: 60, Max: 80}
	tests := []struct {
		name  string
		hit   int
		total int
		want  status.RiskLevel
	}{
		{"danger - below min", 50, 100, status.RiskDanger},
		{"warning - in range", 70, 100, status.RiskWarning},
		{"safe - above max", 90, 100, status.RiskSafe},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := model.CoverageMetrics{MethodsHit: tt.hit, MethodsValid: tt.total}
			m.Calculated = map[config.MetricKey]any{config.MethodsHit: model.CoverageDetail{Percentage: utils.CalculatePercentage(tt.hit, tt.total, 2)}}
			lvl, show := MethodsHitEvaluator{}.Evaluate(m, band)
			assert.True(t, show)
			assert.Equal(t, tt.want, lvl)
		})
	}
}

// ---------------------------------------------------------------------------
// PatchLineCoverageEvaluator (representative always-applicable evaluator)
// ---------------------------------------------------------------------------

func TestPatchLineCoverageEvaluator_Key(t *testing.T) {
	assert.Equal(t, config.PatchLineCoverage, PatchLineCoverageEvaluator{}.Key())
}

func TestPatchLineCoverageEvaluator_IsApplicable(t *testing.T) {
	assert.True(t, PatchLineCoverageEvaluator{}.IsApplicable(status.Capabilities{}))
}

func TestPatchLineCoverageEvaluator_ZeroGuard(t *testing.T) {
	m := model.CoverageMetrics{PatchLinesValid: 0}
	band := &config.Band{Min: 60, Max: 80}
	lvl, show := PatchLineCoverageEvaluator{}.Evaluate(m, band)
	assert.Equal(t, status.RiskLevel(""), lvl)
	assert.False(t, show)
}

func TestPatchLineCoverageEvaluator_Thresholds(t *testing.T) {
	band := &config.Band{Min: 60, Max: 80}
	tests := []struct {
		name    string
		covered int
		valid   int
		want    status.RiskLevel
	}{
		{"danger", 50, 100, status.RiskDanger},
		{"warning", 70, 100, status.RiskWarning},
		{"safe", 90, 100, status.RiskSafe},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := model.CoverageMetrics{PatchLinesCovered: tt.covered, PatchLinesValid: tt.valid}
			m.Calculated = map[config.MetricKey]any{config.PatchLineCoverage: model.CoverageDetail{Percentage: utils.CalculatePercentage(tt.covered, tt.valid, 2)}}
			lvl, show := PatchLineCoverageEvaluator{}.Evaluate(m, band)
			assert.True(t, show)
			assert.Equal(t, tt.want, lvl)
		})
	}
}
