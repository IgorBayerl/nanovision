package lang_go_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/IgorBayerl/AdlerCov/internal/language/lang_go"
	"github.com/IgorBayerl/AdlerCov/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestFile(t *testing.T, content string) (string, func()) {
	t.Helper()
	tmpFile, err := os.CreateTemp("", "test*.go")
	require.NoError(t, err, "Failed to create temp file")
	_, err = tmpFile.WriteString(content)
	require.NoError(t, err, "Failed to write to temp file")
	err = tmpFile.Close()
	require.NoError(t, err, "Failed to close temp file")
	cleanup := func() { os.Remove(tmpFile.Name()) }
	absPath, err := filepath.Abs(tmpFile.Name())
	require.NoError(t, err)
	return absPath, cleanup
}

func TestGoProcessor_Detect(t *testing.T) {
	p := lang_go.NewGoProcessor()
	assert.True(t, p.Detect("file.go"))
	assert.False(t, p.Detect("file.txt"))
}

func TestGoProcessor_Name(t *testing.T) {
	p := lang_go.NewGoProcessor()
	assert.Equal(t, "Go", p.Name())
}

func TestGoProcessor_AnalyzeFile(t *testing.T) {
	testCases := []struct {
		name            string
		sourceCode      string
		expectedMetrics []model.MethodMetrics
		expectError     bool
	}{
		{
			name: "GoldenPath_FunctionAndMethods",
			sourceCode: `package main

import "fmt"

func simpleFunction() {
	fmt.Println("Hello")
}

type MyStruct struct{}

func (s MyStruct) ValueReceiverMethod() int {
	return 1
}

func (s *MyStruct) PointerReceiverMethod(x int) {
	if x > 0 {
		fmt.Println("Positive")
	}
}
`,
			expectedMetrics: []model.MethodMetrics{
				{Name: "simpleFunction", StartLine: 5, EndLine: 7, CyclomaticComplexity: 1},
				{Name: "(MyStruct).ValueReceiverMethod", StartLine: 11, EndLine: 13, CyclomaticComplexity: 1},
				{Name: "(*MyStruct).PointerReceiverMethod", StartLine: 15, EndLine: 19, CyclomaticComplexity: 2},
			},
		},
		{
			name: "ComplexFunction_ForCyclomaticComplexity",
			sourceCode: `package main

func complexFunc(a, b int) int {
	if a > 0 && b > 0 {
		for i := 0; i < a; i++ {
			switch b {
			case 1:
				return 1
			case 2:
				return 2
			default:
				return 3
			}
		}
	}
	return 0
}
`,
			expectedMetrics: []model.MethodMetrics{
				{Name: "complexFunc", StartLine: 3, EndLine: 17, CyclomaticComplexity: 6},
			},
		},
		{
			name:            "NoFunctions_ShouldReturnEmpty",
			sourceCode:      `package main; var X = 1`,
			expectedMetrics: []model.MethodMetrics{},
		},
		{
			name:        "InvalidSyntax_ShouldReturnError",
			sourceCode:  `package main; func main() { oops`,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filePath, cleanup := setupTestFile(t, tc.sourceCode)
			defer cleanup()

			p := lang_go.NewGoProcessor()
			sourceLines := strings.Split(tc.sourceCode, "\n")

			methods, err := p.AnalyzeFile(filePath, sourceLines)

			if tc.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, methods)
			assert.ElementsMatch(t, tc.expectedMetrics, methods, "The discovered methods did not match the expected metrics.")
		})
	}
}
