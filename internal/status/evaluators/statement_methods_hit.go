package evaluators

import (
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/status"
)

// StatementMethodsHitEvaluator evaluates statement-methods-hit coverage (higher is better).
// It is only applicable when the dataset contains statement coverage data.
type StatementMethodsHitEvaluator struct{}

func (StatementMethodsHitEvaluator) Key() config.MetricKey { return config.StatementMethodsHit }
func (StatementMethodsHitEvaluator) Name() string          { return "Statement Methods Hit" }
func (StatementMethodsHitEvaluator) Description() string {
	return "Percentage of methods with at least one statement hit."
}
func (StatementMethodsHitEvaluator) SupportedScopes() status.MetricScope {
	return status.FileScope
}

func (StatementMethodsHitEvaluator) IsApplicable(caps status.Capabilities) bool {
	return caps.HasStatementCoverage
}

func (e StatementMethodsHitEvaluator) Evaluate(m model.CoverageMetrics, band *config.Band) (status.RiskLevel, bool) {
	calc, exists := m.Calculated[e.Key()]
	if !exists {
		return "", false
	}
	detail := calc.(model.CoverageDetail)
	return status.ClassifyHigherIsBetter(detail.Percentage, band)
}
