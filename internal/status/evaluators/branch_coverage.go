package evaluators

import (
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/status"
	"github.com/IgorBayerl/nanovision/internal/utils"
)

// BranchCoverageEvaluator evaluates branch coverage (higher is better).
// It is only applicable when the dataset contains branch data.
type BranchCoverageEvaluator struct{}

func (BranchCoverageEvaluator) Key() config.MetricKey { return config.BranchCoverage }
func (BranchCoverageEvaluator) Name() string          { return "Branch Coverage" }
func (BranchCoverageEvaluator) Description() string   { return "Percentage of covered code branches." }
func (BranchCoverageEvaluator) SupportedScopes() status.MetricScope {
	return status.FileScope
}

func (BranchCoverageEvaluator) IsApplicable(caps status.Capabilities) bool {
	return caps.HasBranchCoverage
}

func (BranchCoverageEvaluator) Evaluate(m model.CoverageMetrics, band *config.Band) (status.RiskLevel, bool) {
	if m.BranchesValid == 0 {
		return "", false
	}
	pct := utils.CalculatePercentage(m.BranchesCovered, m.BranchesValid, 2)
	return status.ClassifyHigherIsBetter(pct, band)
}
