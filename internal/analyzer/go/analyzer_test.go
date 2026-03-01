package golang

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoAnalyzer_Statements(t *testing.T) {
	code := `package main

func main() {
	x := 10        // short_var_declaration (Line 4)
	if x > 5 {     // binary_expression (Line 5) - Atomic Condition
		return     // return_statement (Line 6)
	}
}
`
	analyzer := New()
	result, err := analyzer.Analyze([]byte(code))

	assert.NoError(t, err)

	// Expect 3 statements:
	// 1. x := 10
	// 2. x > 5
	// 3. return
	assert.Equal(t, 3, len(result.Statements))

	// Variable Declaration
	assert.Equal(t, "short_var_declaration", result.Statements[0].Type)
	assert.Equal(t, 4, result.Statements[0].StartLine)

	// If Condition (Atomic Logic)
	// We now catch the expression 'x > 5' directly, not the container.
	assert.Equal(t, "binary_expression", result.Statements[1].Type)
	assert.Equal(t, 5, result.Statements[1].StartLine)
	assert.Equal(t, 5, result.Statements[1].EndLine) // Strict atomic line

	// Return
	assert.Equal(t, "return_statement", result.Statements[2].Type)
	assert.Equal(t, 6, result.Statements[2].StartLine)
}
