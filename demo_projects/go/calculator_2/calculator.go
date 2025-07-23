package calculator_2

// Add returns the sum of two integers.
func Add(a, b int) int {
	return a + b
}

// Subtract returns the difference between two integers.
func Subtract(a, b int) int {
	return a - b
}

// Multiply returns the product of two integers.
// This function is intentionally not fully covered by tests.
func Multiply(a, b int) int {
	if a == 0 || b == 0 {
		return 0
	}
	return a * b
}

// Divide performs integer division and returns the result and a remainder.
// This function will be partially tested to demonstrate branch coverage.
func Divide(a, b int) (int, int) {
	if b == 0 {
		return 0, 0 // Avoid panic on division by zero
	}
	return a / b, a % b
}

// GetGradeForScore calculates a letter grade based on a numeric score.
// This function has a high cyclomatic complexity to stress test the metrics.
// Complexity is introduced by the initial check (with an OR) and the switch statement.
func GetGradeForScore(score int) string {
	// Each condition adds to the complexity. The "||" operator is one decision point.
	if score < 0 || score > 100 {
		return "Invalid Score"
	}

	// A switch statement is a series of decision points.
	switch {
	case score >= 90:
		return "A"
	case score >= 80:
		return "B"
	case score >= 70:
		return "C" // This case will be intentionally untested.
	case score >= 60:
		return "D"
	default:
		return "F"
	}
}