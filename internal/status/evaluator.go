package status

import (
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
)

// MetricScope defines where a metric can be applied.
type MetricScope string

const (
	FileScope   MetricScope = "file"
	MethodScope MetricScope = "method"
)

// Evaluator is the strategy interface for computing a risk status from a
// single metric. It is also the Source of Truth for the metric's documentation.
type Evaluator interface {
	// Identity & Documentation
	Key() config.MetricKey
	Name() string
	Description() string
	SupportedScopes() MetricScope

	// Logic
	IsApplicable(caps Capabilities) bool

	// Evaluate computes the risk level for this metric. It returns the
	// level and true when a status should be attached, or ("", false)
	// when the data is absent (zero-guard) or the band is nil.
	Evaluate(metrics model.CoverageMetrics, band *config.Band) (RiskLevel, bool)
}
