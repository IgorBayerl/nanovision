package utils

// isValidUnixSeconds checks if a timestamp (in seconds) is within a reasonable range.
// E.g., between 1975-01-01 and 2100-01-01.
func IsValidUnixSeconds(ts int64) bool {
	const minValidSeconds int64 = 157766400  // Approx 1975-01-01 UTC
	const maxValidSeconds int64 = 4102444800 // Approx 2100-01-01 UTC
	return ts >= minValidSeconds && ts <= maxValidSeconds
}
