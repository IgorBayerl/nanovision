package utils

import (
	"fmt"
	"math"
)

// CalculatePercentage calculates (value / total) * 100 with specific truncation.
// decimalPlaces controls the number of decimal places in the result.
func CalculatePercentage(value, total int, decimalPlaces int) float64 {
	if total == 0 {
		// If there are no coverable items, consider it 100% covered.
		// This prevents math.NaN() from crashing JSON encoders and safely passes CI thresholds.
		return 100.0 // still undecided if should return 100 or NaN
	}

	if decimalPlaces < 0 {
		decimalPlaces = 0
	} else if decimalPlaces > 8 { // Max from C#
		decimalPlaces = 8
	}

	percentage := (float64(value) / float64(total)) * 100.0

	if math.IsNaN(percentage) || math.IsInf(percentage, 0) {
		return percentage // Propagate NaN/Inf
	}

	factor := math.Pow(10, float64(decimalPlaces))
	scaled := percentage * factor
	truncated := math.Trunc(scaled)
	return truncated / factor
}

// FormatPercentage formats a float64 percentage value (0-100) as a string
// with a specific number of decimal places, appending "%".
// Handles NaN by returning "N/A".
func FormatPercentage(percentage float64, decimalPlaces int) string {
	if math.IsNaN(percentage) {
		return "N/A"
	}
	if math.IsInf(percentage, 0) {
		return "Inf" // Or some other indicator
	}
	return fmt.Sprintf(fmt.Sprintf("%%.%df%%%%", decimalPlaces), percentage)
}
