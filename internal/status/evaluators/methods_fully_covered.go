package evaluators

import (
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/status"
)

// MethodsFullyCoveredEvaluator evaluates methods-fully-covered coverage (higher is better).
// It is only applicable when the dataset contains method-level data.
type MethodsFullyCoveredEvaluator struct{}

func (MethodsFullyCoveredEvaluator) Key() config.MetricKey { return config.MethodsFullyCovered }
func (MethodsFullyCoveredEvaluator) Name() string          { return "Methods Fully Covered" }
func (MethodsFullyCoveredEvaluator) Description() string {
	return "Percentage of methods with 100% line coverage."
}
func (MethodsFullyCoveredEvaluator) SupportedScopes() status.MetricScope {
	return status.FileScope
}

func (MethodsFullyCoveredEvaluator) IsApplicable(caps status.Capabilities) bool {
	return caps.HasMethodCoverage
}

func (e MethodsFullyCoveredEvaluator) Evaluate(m model.CoverageMetrics, band *config.Band) (status.RiskLevel, bool) {
	calc, exists := m.Calculated[e.Key()]
	if !exists {
		return "", false
	}
	detail := calc.(model.CoverageDetail)
	return status.ClassifyHigherIsBetter(detail.Percentage, band)
}
