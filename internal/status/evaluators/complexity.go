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

func (MaxComplexityEvaluator) Evaluate(m model.CoverageMetrics, band *config.Band) (status.RiskLevel, bool) {
	if m.MaxCyclomaticComplexity == 0 {
		return "", false
	}
	return status.ClassifyLowerIsBetter(float64(m.MaxCyclomaticComplexity), band)
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

func (CyclomaticComplexityEvaluator) Evaluate(m model.CoverageMetrics, band *config.Band) (status.RiskLevel, bool) {
	if m.MaxCyclomaticComplexity == 0 {
		return "", false
	}
	return status.ClassifyLowerIsBetter(float64(m.MaxCyclomaticComplexity), band)
}
