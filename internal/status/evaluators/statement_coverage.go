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

func (StatementCoverageEvaluator) IsApplicable(_ status.Capabilities) bool { return true }

func (StatementCoverageEvaluator) Evaluate(m model.CoverageMetrics, band *config.Band) (status.RiskLevel, bool) {
	if m.StatementsValid == 0 {
		return "", false
	}
	pct := utils.CalculatePercentage(m.StatementsCovered, m.StatementsValid, 2)
	return status.ClassifyHigherIsBetter(pct, band)
}
