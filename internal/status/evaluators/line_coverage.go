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

func (LineCoverageEvaluator) IsApplicable(_ status.Capabilities) bool { return true }

func (LineCoverageEvaluator) Evaluate(m model.CoverageMetrics, band *config.Band) (status.RiskLevel, bool) {
	if m.LinesValid == 0 {
		return "", false
	}
	pct := utils.CalculatePercentage(m.LinesCovered, m.LinesValid, 2)
	return status.ClassifyHigherIsBetter(pct, band)
}
