package gdscript

import (
	"testing"

	"github.com/IgorBayerl/nanovision/internal/analyzer"
	"github.com/stretchr/testify/assert"
)

func TestGdScriptAnalyzer_Statements(t *testing.T) {
	code := `
func my_func():        # Line 2: Ignored (Container)
    var x = 10         # Line 3: Assignment (Counted)
    if x > 5:          # Line 4: If Condition (Counted - Atomic)
        print("hi")    # Line 5: Call (Counted)
        return         # Line 6: Return (Counted)
`
	analyzerInstance := New()
	result, err := analyzerInstance.Analyze([]byte(code))

	assert.NoError(t, err)

	for _, s := range result.Statements {
		t.Logf("Line %d-%d: %s", s.StartLine, s.EndLine, s.Type)
	}

	assert.Equal(t, 4, len(result.Statements), "Should find exactly 4 statements")

	var ifStmt *analyzer.StatementMetric
	for _, s := range result.Statements {
		sCopy := s
		if s.StartLine == 4 {
			ifStmt = &sCopy
			break
		}
	}

	assert.NotNil(t, ifStmt, "Should find the if condition")
	assert.Equal(t, 4, ifStmt.EndLine, "If statement should strictly cover only the condition line")
}
