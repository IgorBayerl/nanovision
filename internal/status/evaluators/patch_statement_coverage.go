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

func (PatchStatementCoverageEvaluator) IsApplicable(_ status.Capabilities) bool { return true }

func (PatchStatementCoverageEvaluator) Evaluate(m model.CoverageMetrics, band *config.Band) (status.RiskLevel, bool) {
	if m.PatchStatementsValid == 0 {
		return "", false
	}
	pct := utils.CalculatePercentage(m.PatchStatementsCovered, m.PatchStatementsValid, 2)
	return status.ClassifyHigherIsBetter(pct, band)
}
