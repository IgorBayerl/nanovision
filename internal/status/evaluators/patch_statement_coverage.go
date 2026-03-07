package evaluators

import (
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/status"
	"github.com/IgorBayerl/nanovision/internal/utils"
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

func (PatchStatementCoverageEvaluator) Evaluate(m model.CoverageMetrics, band *config.Band) (status.RiskLevel, bool) {
	if m.PatchStatementsValid == 0 {
		return "", false
	}
	pct := utils.CalculatePercentage(m.PatchStatementsCovered, m.PatchStatementsValid, 2)
	return status.ClassifyHigherIsBetter(pct, band)
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

func (MethodPatchStatementCoverageEvaluator) Evaluate(m model.CoverageMetrics, band *config.Band) (status.RiskLevel, bool) {
	if m.PatchStatementsValid == 0 {
		return "", false
	}
	pct := utils.CalculatePercentage(m.PatchStatementsCovered, m.PatchStatementsValid, 2)
	return status.ClassifyHigherIsBetter(pct, band)
}
