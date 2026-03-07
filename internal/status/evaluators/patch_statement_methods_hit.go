package evaluators

import (
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/status"
)

// PatchStatementMethodsHitEvaluator evaluates patch-statement-methods-hit
// coverage (higher is better).
type PatchStatementMethodsHitEvaluator struct{}

func (PatchStatementMethodsHitEvaluator) Key() config.MetricKey {
	return config.PatchStatementMethodsHit
}
func (PatchStatementMethodsHitEvaluator) Name() string { return "Patch Statement Methods Hit" }
func (PatchStatementMethodsHitEvaluator) Description() string {
	return "Percentage of patched methods with at least one statement hit."
}
func (PatchStatementMethodsHitEvaluator) SupportedScopes() status.MetricScope {
	return status.FileScope
}

func (PatchStatementMethodsHitEvaluator) IsApplicable(_ status.Capabilities) bool { return true }

func (e PatchStatementMethodsHitEvaluator) Evaluate(m model.CoverageMetrics, band *config.Band) (status.RiskLevel, bool) {
	calc, exists := m.Calculated[e.Key()]
	if !exists {
		return "", false
	}
	detail := calc.(model.CoverageDetail)
	return status.ClassifyHigherIsBetter(detail.Percentage, band)
}
