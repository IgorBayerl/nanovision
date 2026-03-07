package evaluators

import (
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/status"
	"github.com/IgorBayerl/nanovision/internal/utils"
)

// LineCoverageEvaluator evaluates line coverage (higher is better).
type LineCoverageEvaluator struct{}

func (LineCoverageEvaluator) Key() config.MetricKey { return config.LineCoverage }
func (LineCoverageEvaluator) Name() string          { return "Line Coverage" }
func (LineCoverageEvaluator) Description() string   { return "Percentage of executed code lines." }
func (LineCoverageEvaluator) SupportedScopes() status.MetricScope {
	return status.FileScope
}

// MethodLineCoverageEvaluator evaluates line coverage (higher is better) for methods.
type MethodLineCoverageEvaluator struct{}

func (MethodLineCoverageEvaluator) Key() config.MetricKey { return config.MethodLineCoverage }
func (MethodLineCoverageEvaluator) Name() string          { return "Line Coverage" }
func (MethodLineCoverageEvaluator) Description() string   { return "Percentage of executed code lines." }
func (MethodLineCoverageEvaluator) SupportedScopes() status.MetricScope {
	return status.MethodScope
}

func (MethodLineCoverageEvaluator) IsApplicable(_ status.Capabilities) bool { return true }

func (MethodLineCoverageEvaluator) Evaluate(m model.CoverageMetrics, band *config.Band) (status.RiskLevel, bool) {
	if m.LinesValid == 0 {
		return "", false
	}
	pct := utils.CalculatePercentage(m.LinesCovered, m.LinesValid, 2)
	return status.ClassifyHigherIsBetter(pct, band)
}

func (LineCoverageEvaluator) IsApplicable(_ status.Capabilities) bool { return true }

func (LineCoverageEvaluator) Evaluate(m model.CoverageMetrics, band *config.Band) (status.RiskLevel, bool) {
	if m.LinesValid == 0 {
		return "", false
	}
	pct := utils.CalculatePercentage(m.LinesCovered, m.LinesValid, 2)
	return status.ClassifyHigherIsBetter(pct, band)
}
