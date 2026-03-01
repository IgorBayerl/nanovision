package cpp

import (
	"testing"

	"github.com/IgorBayerl/nanovision/internal/analyzer"
	"github.com/stretchr/testify/assert"
)

func TestCppAnalyzer_Statements(t *testing.T) {
	code := `
int my_func() {        // Line 2: Ignored (Container)
    int x = 10;        // Line 3: Declaration (Counted)
    if (x > 5) {       // Line 4: If Condition (Counted - Atomic)
        x++;           // Line 5: Expression (Counted)
        return 0;      // Line 6: Return (Counted)
    }
    return 1;          // Line 8: Return (Counted)
}
`
	analyzerInstance := New()
	result, err := analyzerInstance.Analyze([]byte(code))

	assert.NoError(t, err)

	for _, s := range result.Statements {
		t.Logf("Line %d-%d: %s", s.StartLine, s.EndLine, s.Type)
	}

	assert.Equal(t, 5, len(result.Statements), "Should find exactly 5 atomic statements")

	var ifStmt *analyzer.StatementMetric
	for _, s := range result.Statements {
		sCopy := s
		if s.StartLine == 4 {
			ifStmt = &sCopy
			break
		}
	}

	assert.NotNil(t, ifStmt, "Should find the if condition on line 4")
	assert.Equal(t, 4, ifStmt.EndLine, "If statement should strictly cover only the condition line")
}
