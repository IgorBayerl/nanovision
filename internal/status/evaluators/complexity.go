package evaluators

import (
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/status"
)

// MaxComplexityEvaluator evaluates cyclomatic complexity (lower is better).
type MaxComplexityEvaluator struct{}

func (MaxComplexityEvaluator) Key() config.MetricKey { return config.MaxCyclomaticComplexity }

func (MaxComplexityEvaluator) IsApplicable(_ status.Capabilities) bool { return true }

func (MaxComplexityEvaluator) Evaluate(m model.CoverageMetrics, band *config.Band) (status.RiskLevel, bool) {
	if m.MaxCyclomaticComplexity == 0 {
		return "", false
	}
	return status.ClassifyLowerIsBetter(float64(m.MaxCyclomaticComplexity), band)
}
