package calculator

import (
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
)

// MetricCalculator defines the calculation logic for a single file/directory-level coverage metric.
type MetricCalculator interface {
	Key() config.MetricKey
	// Calculate takes raw coverage metrics and returns a shaped detail (like model.CoverageDetail or model.ScoreDetail),
	// and a boolean indicating if calculation was possible (e.g. valid total > 0).
	Calculate(raw model.CoverageMetrics) (any, bool)
}

// MethodMetricCalculator defines the calculation logic for a single method-level metric.
type MethodMetricCalculator interface {
	Key() config.MetricKey
	Calculate(raw model.MethodMetrics) (any, bool)
}
