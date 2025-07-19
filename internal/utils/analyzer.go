package utils

import (
	"math"
	"strconv"
	"strings"
)

// parseInt is a utility function to parse string to int, ignoring errors for simplicity.
func ParseInt(s string, fallback int) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		return fallback
	}
	return val
}

// parseFloat is a utility function to parse string to float64, ignoring errors.
func ParseFloat(s string) float64 {
	if strings.ToLower(s) == "nan" { // Handle "NaN" case-insensitively
		return math.NaN() // Or return 0 or a specific indicator if preferred
	}
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

// isValidUnixSeconds checks if a timestamp (in seconds) is within a reasonable range.
// E.g., between 1975-01-01 and 2100-01-01.
func IsValidUnixSeconds(ts int64) bool {
	const minValidSeconds int64 = 157766400  // Approx 1975-01-01 UTC
	const maxValidSeconds int64 = 4102444800 // Approx 2100-01-01 UTC
	return ts >= minValidSeconds && ts <= maxValidSeconds
}
