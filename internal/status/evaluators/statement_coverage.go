package evaluators

import (
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/status"
	"github.com/IgorBayerl/nanovision/internal/utils"
)

// StatementCoverageEvaluator evaluates statement coverage (higher is better).
type StatementCoverageEvaluator struct{}

func (StatementCoverageEvaluator) Key() config.MetricKey { return config.StatementCoverage }
func (StatementCoverageEvaluator) Name() string          { return "Statement Coverage" }
func (StatementCoverageEvaluator) Description() string   { return "Percentage of executed statements." }
func (StatementCoverageEvaluator) SupportedScopes() status.MetricScope {
	return status.FileScope
}

// MethodStatementCoverageEvaluator evaluates statement coverage (higher is better) for methods.
type MethodStatementCoverageEvaluator struct{}

func (MethodStatementCoverageEvaluator) Key() config.MetricKey { return config.MethodStatementCoverage }
func (MethodStatementCoverageEvaluator) Name() string          { return "Statement Coverage" }
func (MethodStatementCoverageEvaluator) Description() string {
	return "Percentage of executed statements."
}
func (MethodStatementCoverageEvaluator) SupportedScopes() status.MetricScope {
	return status.MethodScope
}

func (MethodStatementCoverageEvaluator) IsApplicable(_ status.Capabilities) bool { return true }

func (MethodStatementCoverageEvaluator) Evaluate(m model.CoverageMetrics, band *config.Band) (status.RiskLevel, bool) {
	if m.StatementsValid == 0 {
		return "", false
	}
	pct := utils.CalculatePercentage(m.StatementsCovered, m.StatementsValid, 2)
	return status.ClassifyHigherIsBetter(pct, band)
}

func (StatementCoverageEvaluator) IsApplicable(_ status.Capabilities) bool { return true }

func (StatementCoverageEvaluator) Evaluate(m model.CoverageMetrics, band *config.Band) (status.RiskLevel, bool) {
	if m.StatementsValid == 0 {
		return "", false
	}
	pct := utils.CalculatePercentage(m.StatementsCovered, m.StatementsValid, 2)
	return status.ClassifyHigherIsBetter(pct, band)
}
