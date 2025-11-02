package status

import "github.com/IgorBayerl/nanovision/internal/config"

// Classify takes a percentage and a risk band configuration and returns the
// corresponding risk level. It is the core of the status evaluation logic.
//
// The rules are as follows:
//   - If the `band` is nil (meaning the metric was not configured in `status_bands`
//     in the YAML file), it returns `show: false`, indicating that no status
//     should be displayed for this metric.
//   - If the percentage is BELOW the band's minimum, it is classified as `RiskDanger`.
//   - If the percentage is WITHIN the band's range [Min, Max] (inclusive), it is `RiskWarning`.
//   - If the percentage is ABOVE the band's maximum, it is `RiskSafe`.
//
// The returned `show` boolean is a signal to the caller (`Annotate`) on whether
// a status should be attached to the data model.
func Classify(pct float64, band *config.Band) (level RiskLevel, show bool) {
	if band == nil {
		// This metric is not configured for status reporting; do not show anything.
		return "", false
	}
	if pct < band.Min {
		return RiskDanger, true
	}
	if pct <= band.Max {
		return RiskWarning, true
	}
	return RiskSafe, true
}
