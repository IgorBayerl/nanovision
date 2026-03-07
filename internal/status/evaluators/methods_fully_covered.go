package evaluators

import (
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/status"
	"github.com/IgorBayerl/nanovision/internal/utils"
)

// MethodsFullyCoveredEvaluator evaluates methods-fully-covered coverage (higher is better).
// It is only applicable when the dataset contains method-level data.
type MethodsFullyCoveredEvaluator struct{}

func (MethodsFullyCoveredEvaluator) Key() config.MetricKey { return config.MethodsFullyCovered }
func (MethodsFullyCoveredEvaluator) Name() string          { return "Methods Fully Covered" }
func (MethodsFullyCoveredEvaluator) Description() string {
	return "Percentage of methods with 100% line coverage."
}
func (MethodsFullyCoveredEvaluator) SupportedScopes() status.MetricScope {
	return status.FileScope
}

func (MethodsFullyCoveredEvaluator) IsApplicable(caps status.Capabilities) bool {
	return caps.HasMethodCoverage
}

func (MethodsFullyCoveredEvaluator) Evaluate(m model.CoverageMetrics, band *config.Band) (status.RiskLevel, bool) {
	if m.MethodsValid == 0 {
		return "", false
	}
	pct := utils.CalculatePercentage(m.MethodsFullyCovered, m.MethodsValid, 2)
	return status.ClassifyHigherIsBetter(pct, band)
}
