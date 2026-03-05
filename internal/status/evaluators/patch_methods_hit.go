package evaluators

import (
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/status"
	"github.com/IgorBayerl/nanovision/internal/utils"
)

// PatchMethodsHitEvaluator evaluates patch methods-hit coverage (higher is better).
type PatchMethodsHitEvaluator struct{}

func (PatchMethodsHitEvaluator) Key() config.MetricKey { return config.PatchMethodsHit }

func (PatchMethodsHitEvaluator) IsApplicable(_ status.Capabilities) bool { return true }

func (PatchMethodsHitEvaluator) Evaluate(m model.CoverageMetrics, band *config.Band) (status.RiskLevel, bool) {
	if m.PatchMethodsValid == 0 {
		return "", false
	}
	pct := utils.CalculatePercentage(m.PatchMethodsHit, m.PatchMethodsValid, 2)
	return status.ClassifyHigherIsBetter(pct, band)
}
