package evaluators

import (
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/status"
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

func (e MethodLineCoverageEvaluator) Evaluate(m model.CoverageMetrics, band *config.Band) (status.RiskLevel, bool) {
	calc, exists := m.Calculated[e.Key()]
	if !exists {
		return "", false
	}
	detail := calc.(model.CoverageDetail)
	return status.ClassifyHigherIsBetter(detail.Percentage, band)
}

func (LineCoverageEvaluator) IsApplicable(_ status.Capabilities) bool { return true }

func (e LineCoverageEvaluator) Evaluate(m model.CoverageMetrics, band *config.Band) (status.RiskLevel, bool) {
	calc, exists := m.Calculated[e.Key()]
	if !exists {
		return "", false
	}
	detail := calc.(model.CoverageDetail)
	return status.ClassifyHigherIsBetter(detail.Percentage, band)
}
