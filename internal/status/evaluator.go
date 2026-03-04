package status

import (
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
)

// Evaluator is the strategy interface for computing a risk status from a
// single metric. Each implementation knows how to extract its own value
// from CoverageMetrics, guard against invalid data, and classify the
// result against a Band.
type Evaluator interface {
	// Key returns the config.MetricKey this evaluator handles.
	Key() config.MetricKey

	// IsApplicable returns whether this metric should be evaluated given
	// the dataset's capabilities (e.g., branch coverage may not be present).
	IsApplicable(caps Capabilities) bool

	// Evaluate computes the risk level for this metric. It returns the
	// level and true when a status should be attached, or ("", false)
	// when the data is absent (zero-guard) or the band is nil.
	Evaluate(metrics model.CoverageMetrics, band *config.Band) (RiskLevel, bool)
}
