package evaluators

import (
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/status"
)

// MethodsHitEvaluator evaluates methods-hit coverage (higher is better).
// It is only applicable when the dataset contains method-level data.
type MethodsHitEvaluator struct{}

func (MethodsHitEvaluator) Key() config.MetricKey { return config.MethodsHit }
func (MethodsHitEvaluator) Name() string          { return "Methods Hit" }
func (MethodsHitEvaluator) Description() string {
	return "Percentage of methods with at least one hit."
}
func (MethodsHitEvaluator) SupportedScopes() status.MetricScope {
	return status.FileScope
}

func (MethodsHitEvaluator) IsApplicable(caps status.Capabilities) bool {
	return caps.HasMethodCoverage
}

func (e MethodsHitEvaluator) Evaluate(m model.CoverageMetrics, band *config.Band) (status.RiskLevel, bool) {
	calc, exists := m.Calculated[e.Key()]
	if !exists {
		return "", false
	}
	detail := calc.(model.CoverageDetail)
	return status.ClassifyHigherIsBetter(detail.Percentage, band)
}
