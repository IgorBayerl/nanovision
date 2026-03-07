package evaluators

import (
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/status"
)

// StatementMethodsFullyCoveredEvaluator evaluates statement-methods-fully-covered
// coverage (higher is better). Only applicable when statement coverage data exists.
type StatementMethodsFullyCoveredEvaluator struct{}

func (StatementMethodsFullyCoveredEvaluator) Key() config.MetricKey {
	return config.StatementMethodsFullyCovered
}
func (StatementMethodsFullyCoveredEvaluator) Name() string { return "Statement Methods Fully Covered" }
func (StatementMethodsFullyCoveredEvaluator) Description() string {
	return "Percentage of methods with 100% statement coverage."
}
func (StatementMethodsFullyCoveredEvaluator) SupportedScopes() status.MetricScope {
	return status.FileScope
}

func (StatementMethodsFullyCoveredEvaluator) IsApplicable(caps status.Capabilities) bool {
	return caps.HasStatementCoverage
}

func (e StatementMethodsFullyCoveredEvaluator) Evaluate(m model.CoverageMetrics, band *config.Band) (status.RiskLevel, bool) {
	calc, exists := m.Calculated[e.Key()]
	if !exists {
		return "", false
	}
	detail := calc.(model.CoverageDetail)
	return status.ClassifyHigherIsBetter(detail.Percentage, band)
}
