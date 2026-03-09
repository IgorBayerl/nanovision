package evaluators

import (
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/status"
)

// MaxComplexityEvaluator evaluates cyclomatic complexity (lower is better).
type MaxComplexityEvaluator struct{}

func (MaxComplexityEvaluator) Key() config.MetricKey { return config.MaxCyclomaticComplexity }
func (MaxComplexityEvaluator) Name() string          { return "Max Complexity" }
func (MaxComplexityEvaluator) Description() string {
	return "Maximum cyclomatic complexity of a function (lower is better)."
}
func (MaxComplexityEvaluator) SupportedScopes() status.MetricScope {
	return status.FileScope
}

func (MaxComplexityEvaluator) IsApplicable(_ status.Capabilities) bool { return true }

func (e MaxComplexityEvaluator) Evaluate(m model.CoverageMetrics, band *config.Band) (status.RiskLevel, bool) {
	calc, exists := m.Calculated[e.Key()]
	if !exists {
		return "", false
	}
	detail := calc.(model.ScoreDetail)
	return status.ClassifyLowerIsBetter(detail.Value, band)
}

// CyclomaticComplexityEvaluator evaluates cyclomatic complexity (lower is better) for methods.
type CyclomaticComplexityEvaluator struct{}

func (CyclomaticComplexityEvaluator) Key() config.MetricKey { return config.CyclomaticComplexity }
func (CyclomaticComplexityEvaluator) Name() string          { return "Cyclomatic Complexity" }
func (CyclomaticComplexityEvaluator) Description() string {
	return "Cyclomatic complexity of a function (lower is better)."
}
func (CyclomaticComplexityEvaluator) SupportedScopes() status.MetricScope {
	return status.MethodScope
}

func (CyclomaticComplexityEvaluator) IsApplicable(_ status.Capabilities) bool { return true }

func (e CyclomaticComplexityEvaluator) Evaluate(m model.CoverageMetrics, band *config.Band) (status.RiskLevel, bool) {
	calc, exists := m.Calculated[e.Key()]
	if !exists {
		return "", false
	}
	detail := calc.(model.ScoreDetail)
	return status.ClassifyLowerIsBetter(detail.Value, band)
}

// MethodCrapScoreEvaluator

type MethodCrapScoreEvaluator struct{}

func (MethodCrapScoreEvaluator) Key() config.MetricKey { return config.MethodCrapScore }
func (MethodCrapScoreEvaluator) Name() string          { return "CRAP Score" }
func (MethodCrapScoreEvaluator) Description() string {
	return "Change Risk Anti-Pattern (CRAP) score combining complexity and coverage (lower is better)."
}
func (MethodCrapScoreEvaluator) SupportedScopes() status.MetricScope {
	return status.MethodScope
}

func (MethodCrapScoreEvaluator) IsApplicable(_ status.Capabilities) bool { return true }

func (e MethodCrapScoreEvaluator) Evaluate(m model.CoverageMetrics, band *config.Band) (status.RiskLevel, bool) {
	calc, exists := m.Calculated[e.Key()]
	if !exists {
		return "", false
	}
	detail := calc.(model.ScoreDetail)
	return status.ClassifyLowerIsBetter(detail.Value, band)
}

// MethodPatchCrapScoreEvaluator

type MethodPatchCrapScoreEvaluator struct{}

func (MethodPatchCrapScoreEvaluator) Key() config.MetricKey { return config.MethodPatchCrapScore }
func (MethodPatchCrapScoreEvaluator) Name() string          { return "Patch CRAP Score" }
func (MethodPatchCrapScoreEvaluator) Description() string {
	return "CRAP score applied only to patched statements (lower is better)."
}
func (MethodPatchCrapScoreEvaluator) SupportedScopes() status.MetricScope {
	return status.MethodScope
}

func (MethodPatchCrapScoreEvaluator) IsApplicable(_ status.Capabilities) bool { return true }

func (e MethodPatchCrapScoreEvaluator) Evaluate(m model.CoverageMetrics, band *config.Band) (status.RiskLevel, bool) {
	calc, exists := m.Calculated[e.Key()]
	if !exists {
		return "", false
	}
	detail := calc.(model.ScoreDetail)
	return status.ClassifyLowerIsBetter(detail.Value, band)
}

// MethodExposedRiskEvaluator

type MethodExposedRiskEvaluator struct{}

func (MethodExposedRiskEvaluator) Key() config.MetricKey { return config.MethodExposedRisk }
func (MethodExposedRiskEvaluator) Name() string          { return "Exposed Risk" }
func (MethodExposedRiskEvaluator) Description() string {
	return "Absolute volume of complexity that is unprotected by tests (lower is better)."
}
func (MethodExposedRiskEvaluator) SupportedScopes() status.MetricScope {
	return status.MethodScope
}

func (MethodExposedRiskEvaluator) IsApplicable(_ status.Capabilities) bool { return true }

func (e MethodExposedRiskEvaluator) Evaluate(m model.CoverageMetrics, band *config.Band) (status.RiskLevel, bool) {
	calc, exists := m.Calculated[e.Key()]
	if !exists {
		return "", false
	}
	detail := calc.(model.ScoreDetail)
	return status.ClassifyLowerIsBetter(detail.Value, band)
}

// MethodDefectProbabilityEvaluator

type MethodDefectProbabilityEvaluator struct{}

func (MethodDefectProbabilityEvaluator) Key() config.MetricKey { return config.MethodDefectProbability }
func (MethodDefectProbabilityEvaluator) Name() string          { return "Defect Probability Index" }
func (MethodDefectProbabilityEvaluator) Description() string {
	return "Index representing the probability of defects based on complexity and patch coverage (lower is better)."
}
func (MethodDefectProbabilityEvaluator) SupportedScopes() status.MetricScope {
	return status.MethodScope
}

func (MethodDefectProbabilityEvaluator) IsApplicable(_ status.Capabilities) bool { return true }

func (e MethodDefectProbabilityEvaluator) Evaluate(m model.CoverageMetrics, band *config.Band) (status.RiskLevel, bool) {
	calc, exists := m.Calculated[e.Key()]
	if !exists {
		return "", false
	}
	detail := calc.(model.ScoreDetail)
	return status.ClassifyLowerIsBetter(detail.Value, band)
}
