package evaluators

import (
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/status"
)

// PatchMethodsHitEvaluator evaluates patch methods-hit coverage (higher is better).
type PatchMethodsHitEvaluator struct{}

func (PatchMethodsHitEvaluator) Key() config.MetricKey { return config.PatchMethodsHit }
func (PatchMethodsHitEvaluator) Name() string          { return "Patch Methods Hit" }
func (PatchMethodsHitEvaluator) Description() string {
	return "Percentage of patched methods with at least one hit."
}
func (PatchMethodsHitEvaluator) SupportedScopes() status.MetricScope {
	return status.FileScope
}

func (PatchMethodsHitEvaluator) IsApplicable(_ status.Capabilities) bool { return true }

func (e PatchMethodsHitEvaluator) Evaluate(m model.CoverageMetrics, band *config.Band) (status.RiskLevel, bool) {
	calc, exists := m.Calculated[e.Key()]
	if !exists {
		return "", false
	}
	detail := calc.(model.CoverageDetail)
	return status.ClassifyHigherIsBetter(detail.Percentage, band)
}
