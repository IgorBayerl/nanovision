package evaluators

import (
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/status"
)

// PatchStatementCoverageEvaluator evaluates patch statement coverage (higher is better).
type PatchStatementCoverageEvaluator struct{}

func (PatchStatementCoverageEvaluator) Key() config.MetricKey { return config.PatchStatementCoverage }
func (PatchStatementCoverageEvaluator) Name() string          { return "Patch Statement Coverage" }
func (PatchStatementCoverageEvaluator) Description() string {
	return "Statement coverage of changed (patched) code only."
}
func (PatchStatementCoverageEvaluator) SupportedScopes() status.MetricScope {
	return status.FileScope
}

func (PatchStatementCoverageEvaluator) IsApplicable(_ status.Capabilities) bool { return true }

func (e PatchStatementCoverageEvaluator) Evaluate(m model.CoverageMetrics, band *config.Band) (status.RiskLevel, bool) {
	calc, exists := m.Calculated[e.Key()]
	if !exists {
		return "", false
	}
	detail := calc.(model.CoverageDetail)
	return status.ClassifyHigherIsBetter(detail.Percentage, band)
}

// MethodPatchStatementCoverageEvaluator evaluates patch statement coverage (higher is better) for methods.
type MethodPatchStatementCoverageEvaluator struct{}

func (MethodPatchStatementCoverageEvaluator) Key() config.MetricKey {
	return config.MethodPatchStatementCoverage
}
func (MethodPatchStatementCoverageEvaluator) Name() string { return "Patch Statement Coverage" }
func (MethodPatchStatementCoverageEvaluator) Description() string {
	return "Statement coverage of changed (patched) code only."
}
func (MethodPatchStatementCoverageEvaluator) SupportedScopes() status.MetricScope {
	return status.MethodScope
}

func (MethodPatchStatementCoverageEvaluator) IsApplicable(_ status.Capabilities) bool { return true }

func (e MethodPatchStatementCoverageEvaluator) Evaluate(m model.CoverageMetrics, band *config.Band) (status.RiskLevel, bool) {
	calc, exists := m.Calculated[e.Key()]
	if !exists {
		return "", false
	}
	detail := calc.(model.CoverageDetail)
	return status.ClassifyHigherIsBetter(detail.Percentage, band)
}
