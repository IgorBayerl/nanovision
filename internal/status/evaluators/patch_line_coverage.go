package evaluators

import (
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/status"
)

// PatchLineCoverageEvaluator evaluates patch line coverage (higher is better).
type PatchLineCoverageEvaluator struct{}

func (PatchLineCoverageEvaluator) Key() config.MetricKey { return config.PatchLineCoverage }
func (PatchLineCoverageEvaluator) Name() string          { return "Patch Line Coverage" }
func (PatchLineCoverageEvaluator) Description() string {
	return "Line coverage of changed (patched) code only."
}
func (PatchLineCoverageEvaluator) SupportedScopes() status.MetricScope {
	return status.FileScope
}

func (PatchLineCoverageEvaluator) IsApplicable(_ status.Capabilities) bool { return true }

func (e PatchLineCoverageEvaluator) Evaluate(m model.CoverageMetrics, band *config.Band) (status.RiskLevel, bool) {
	calc, exists := m.Calculated[e.Key()]
	if !exists {
		return "", false
	}
	detail := calc.(model.CoverageDetail)
	return status.ClassifyHigherIsBetter(detail.Percentage, band)
}

// MethodPatchLineCoverageEvaluator evaluates patch line coverage (higher is better) for methods.
type MethodPatchLineCoverageEvaluator struct{}

func (MethodPatchLineCoverageEvaluator) Key() config.MetricKey { return config.MethodPatchLineCoverage }
func (MethodPatchLineCoverageEvaluator) Name() string          { return "Patch Line Coverage" }
func (MethodPatchLineCoverageEvaluator) Description() string {
	return "Line coverage of changed (patched) code only."
}
func (MethodPatchLineCoverageEvaluator) SupportedScopes() status.MetricScope {
	return status.MethodScope
}

func (MethodPatchLineCoverageEvaluator) IsApplicable(_ status.Capabilities) bool { return true }

func (e MethodPatchLineCoverageEvaluator) Evaluate(m model.CoverageMetrics, band *config.Band) (status.RiskLevel, bool) {
	calc, exists := m.Calculated[e.Key()]
	if !exists {
		return "", false
	}
	detail := calc.(model.CoverageDetail)
	return status.ClassifyHigherIsBetter(detail.Percentage, band)
}
