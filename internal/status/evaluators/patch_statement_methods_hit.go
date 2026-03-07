package evaluators

import (
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/status"
	"github.com/IgorBayerl/nanovision/internal/utils"
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

func (PatchStatementMethodsHitEvaluator) Evaluate(m model.CoverageMetrics, band *config.Band) (status.RiskLevel, bool) {
	if m.PatchStatementMethodsValid == 0 {
		return "", false
	}
	pct := utils.CalculatePercentage(m.PatchStatementMethodsHit, m.PatchStatementMethodsValid, 2)
	return status.ClassifyHigherIsBetter(pct, band)
}
