package evaluators

import (
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/status"
	"github.com/IgorBayerl/nanovision/internal/utils"
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

func (StatementMethodsHitEvaluator) Evaluate(m model.CoverageMetrics, band *config.Band) (status.RiskLevel, bool) {
	if m.StatementMethodsValid == 0 {
		return "", false
	}
	pct := utils.CalculatePercentage(m.StatementMethodsHit, m.StatementMethodsValid, 2)
	return status.ClassifyHigherIsBetter(pct, band)
}
