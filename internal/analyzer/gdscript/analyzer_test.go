package gdscript

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGdScriptAnalyzer_Statements(t *testing.T) {
	code := `func test_func():
	if true:
		print("hello")
		return
`
	analyzer := New()
	result, err := analyzer.Analyze([]byte(code))

	assert.NoError(t, err)

	assert.NotEmpty(t, result.Statements)

	foundIf := false
	foundReturn := false

	for _, s := range result.Statements {
		if s.Type == "if_statement" {
			foundIf = true
			assert.Equal(t, 2, s.StartLine)
		}
		if s.Type == "return_statement" {
			foundReturn = true
			assert.Equal(t, 4, s.StartLine)
		}
	}

	assert.True(t, foundIf, "should find if_statement")
	assert.True(t, foundReturn, "should find return_statement")
}
