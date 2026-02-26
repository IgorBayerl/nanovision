package golang

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoAnalyzer_Statements(t *testing.T) {
	code := `package main

func main() {
	x := 10        // short_var_declaration, line 4
	if x > 5 {     // if_statement, line 5
		return     // return_statement, line 6
	}
}
`
	analyzer := New()
	result, err := analyzer.Analyze([]byte(code))
	
	assert.NoError(t, err)

	// we expect 3 statements inside main
	assert.Equal(t, 3, len(result.Statements))

	// check first statement
	assert.Equal(t, "short_var_declaration", result.Statements[0].Type)
	assert.Equal(t, 4, result.Statements[0].StartLine)
	assert.Equal(t, 4, result.Statements[0].EndLine)

	// check second statement
	assert.Equal(t, "if_statement", result.Statements[1].Type)
	assert.Equal(t, 5, result.Statements[1].StartLine)
	assert.Equal(t, 7, result.Statements[1].EndLine)

	// check third statement
	assert.Equal(t, "return_statement", result.Statements[2].Type)
	assert.Equal(t, 6, result.Statements[2].StartLine)
	assert.Equal(t, 6, result.Statements[2].EndLine)
}
