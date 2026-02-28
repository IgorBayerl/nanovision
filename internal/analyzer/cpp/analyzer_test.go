package cpp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCppAnalyzer_Statements(t *testing.T) {
	code := `
int main() {
    int x = 10;
    while (x > 0) {
        x--;
    }
    return 0;
}
`
	analyzer := New()
	result, err := analyzer.Analyze([]byte(code))

	assert.NoError(t, err)

	var types []string
	for _, stmt := range result.Statements {
		types = append(types, stmt.Type)
	}

	assert.Contains(t, types, "declaration")
	assert.Contains(t, types, "while_statement")
	assert.Contains(t, types, "expression_statement")
	assert.Contains(t, types, "return_statement")

	foundReturn := false
	for _, s := range result.Statements {
		if s.Type == "return_statement" {
			foundReturn = true
			assert.Equal(t, 7, s.StartLine)
		}
	}
	assert.True(t, foundReturn)
}
