package utils

import (
	"fmt"
	"math"
)

// CalculatePercentage calculates (value / total) * 100 with specific truncation.
// decimalPlaces controls the number of decimal places in the result.
// Mimics C# ReportGenerator.Core.Common.MathExtensions.CalculatePercentage.
func CalculatePercentage(value, total int, decimalPlaces int) float64 {
	if total == 0 {
		// Depending on C# behavior for 0 total: NaN, 0, or specific handling.
		// C# MathExtensions.CalculatePercentage throws ArgumentException if number2 (total) is 0.
		// Let's return NaN for undefined, or 0.0 if preferred for "no coverage".
		// For display, "N/A" is often used.
		return math.NaN() // Or 0.0
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
	// C# uses Truncate: Math.Truncate(factor * (double)number1 / (double)number2) / divisor;
	// which means Math.Truncate(percentage * 10^dp_for_calc_internal) / 10^dp_for_calc_internal
	// The factor/divisor in C# are related to the internal calculation before final scaling.
	// MathExtensions.factor = 100 * Math.Pow(10, maximumDecimalPlaces);
	// MathExtensions.divisor = (int)Math.Pow(10, maximumDecimalPlaces);
	// It effectively calculates (value/total) then truncates to `maximumDecimalPlaces` for the *percentage value*.

	// Let's re-evaluate the C# logic:
	// return (decimal)Math.Truncate(factor * (double)number1 / (double)number2) / divisor;
	// factor = 100 * Math.Pow(10, maximumDecimalPlaces);
	// divisor = (int)Math.Pow(10, maximumDecimalPlaces);
	// This becomes: Math.Truncate(100 * ( (float64)value / (float64)total ) * 10^dp ) / 10^dp
	// This is equivalent to: Truncate(percentage * 10^dp) / 10^dp
	// E.g., if percentage = 75.12345 and dp = 1:
	// Truncate(75.12345 * 10) / 10 = Truncate(751.2345) / 10 = 751 / 10 = 75.1

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
