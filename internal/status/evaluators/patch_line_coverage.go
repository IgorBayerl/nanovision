package evaluators

import (
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/status"
	"github.com/IgorBayerl/nanovision/internal/utils"
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

func (PatchLineCoverageEvaluator) Evaluate(m model.CoverageMetrics, band *config.Band) (status.RiskLevel, bool) {
	if m.PatchLinesValid == 0 {
		return "", false
	}
	pct := utils.CalculatePercentage(m.PatchLinesCovered, m.PatchLinesValid, 2)
	return status.ClassifyHigherIsBetter(pct, band)
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

func (MethodPatchLineCoverageEvaluator) Evaluate(m model.CoverageMetrics, band *config.Band) (status.RiskLevel, bool) {
	if m.PatchLinesValid == 0 {
		return "", false
	}
	pct := utils.CalculatePercentage(m.PatchLinesCovered, m.PatchLinesValid, 2)
	return status.ClassifyHigherIsBetter(pct, band)
}
