package utils_test

import (
	"testing"

	"github.com/IgorBayerl/nanovision/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestFindMatchingBrace(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name           string
		sourceLines    []string
		startLineIndex int
		expectedLine   int
		expectedFound  bool
	}{
		{
			name: "SimpleCase_ReturnsCorrectLineNumber",
			sourceLines: []string{
				"function test() {",
				"    return true;",
				"}",
			},
			startLineIndex: 0,
			expectedLine:   3,
			expectedFound:  true,
		},
		{
			name: "NestedBraces_ReturnsOutermostClosingBrace",
			sourceLines: []string{
				"class Test {",
				"    method() {",
				"        if (condition) {",
				"            doSomething();",
				"        }",
				"    }",
				"}",
			},
			startLineIndex: 0,
			expectedLine:   7,
			expectedFound:  true,
		},
		{
			name: "StartFromMiddle_FindsCorrectMatch",
			sourceLines: []string{
				"outer {",
				"    inner {",
				"        content",
				"    }",
				"}",
			},
			startLineIndex: 1,
			expectedLine:   4,
			expectedFound:  true,
		},
		{
			name: "NoOpeningBraceOnStartLine_ReturnsFalse",
			sourceLines: []string{
				"some code",
				"more code",
				"}",
			},
			startLineIndex: 0,
			expectedFound:  false,
		},
		{
			name: "NoClosingBrace_ReturnsFalse",
			sourceLines: []string{
				"function test() {",
				"    return true;",
				"    // missing closing brace",
			},
			startLineIndex: 0,
			expectedFound:  false,
		},
		{
			name: "BraceInLineComment_IgnoresBrace",
			sourceLines: []string{
				"function test() {",
				"    // this } should be ignored",
				"    return true;",
				"}",
			},
			startLineIndex: 0,
			expectedLine:   4,
			expectedFound:  true,
		},
		{
			name: "BraceInBlockComment_IgnoresBrace",
			sourceLines: []string{
				"function test() {",
				"    /* this } should be ignored */",
				"    return true;",
				"}",
			},
			startLineIndex: 0,
			expectedLine:   4,
			expectedFound:  true,
		},
		{
			name: "MultilineBlockComment_IgnoresBraces",
			sourceLines: []string{
				"function test() {",
				"    /*",
				"     * this } should be ignored",
				"     * and this { too",
				"     */",
				"    return true;",
				"}",
			},
			startLineIndex: 0,
			expectedLine:   7,
			expectedFound:  true,
		},
		{
			name: "BraceInString_IgnoresBrace",
			sourceLines: []string{
				"function test() {",
				"    let msg = \"this } should be ignored\";",
				"    return true;",
				"}",
			},
			startLineIndex: 0,
			expectedLine:   4,
			expectedFound:  true,
		},
		{
			name: "EscapedQuoteInString_HandlesCorrectly",
			sourceLines: []string{
				"function test() {",
				"    let msg = \"escaped quote \\\" and brace }\";",
				"    return true;",
				"}",
			},
			startLineIndex: 0,
			expectedLine:   4,
			expectedFound:  true,
		},
		{
			name: "MultipleBracesOnSameLine_FindsCorrectMatch",
			sourceLines: []string{
				"if (condition) { doSomething(); }",
				"function test() {",
				"    return true;",
				"}",
			},
			startLineIndex: 0,
			expectedLine:   1,
			expectedFound:  true,
		},
		{
			name: "StartIndexBeyondArray_ReturnsFalse",
			sourceLines: []string{
				"function test() {",
				"}",
			},
			startLineIndex: 5,
			expectedFound:  false,
		},
		{
			name:           "EmptySourceLines_ReturnsFalse",
			sourceLines:    []string{},
			startLineIndex: 0,
			expectedFound:  false,
		},
		{
			name: "ComplexNestedStructure_FindsCorrectMatch",
			sourceLines: []string{
				"class MyClass {",
				"    constructor() {",
				"        this.data = {",
				"            nested: {",
				"                value: \"test\"",
				"            }",
				"        };",
				"    }",
				"    method() {",
				"        if (condition) {",
				"            // comment with }",
				"            return \"string with }\";",
				"        }",
				"    }",
				"}",
			},
			startLineIndex: 0,
			expectedLine:   15,
			expectedFound:  true,
		},
		{
			name: "MixedQuotesAndComments_HandlesCorrectly",
			sourceLines: []string{
				"function test() {",
				"    let str1 = \"quote with // comment\";",
				"    /* block comment with \\\" quote */",
				"    let str2 = \"another } brace\";",
				"    // line comment with }",
				"}",
			},
			startLineIndex: 0,
			expectedLine:   6,
			expectedFound:  true,
		},
		{
			name: "MultipleOpeningBracesOnSameLine",
			sourceLines: []string{
				"{ { {",
				"content",
				"} } }",
			},
			startLineIndex: 0,
			expectedLine:   3,
			expectedFound:  true,
		},
		{
			name: "UnbalancedBraces_ReturnsFalse",
			sourceLines: []string{
				"{ { {",
				"content",
				"} }",
			},
			startLineIndex: 0,
			expectedFound:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			lineNumber, found := utils.FindMatchingBrace(tc.sourceLines, tc.startLineIndex)

			// Assert
			assert.Equal(t, tc.expectedFound, found)
			if tc.expectedFound {
				assert.Equal(t, tc.expectedLine, lineNumber)
			}
		})
	}
}

// Benchmark tests for performance evaluation
func BenchmarkFindMatchingBrace_SimpleCase(b *testing.B) {
	sourceLines := []string{
		"function test() {",
		"    return true;",
		"}",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		utils.FindMatchingBrace(sourceLines, 0)
	}
}

func BenchmarkFindMatchingBrace_LargeFile(b *testing.B) {
	// Create a large source file simulation
	sourceLines := make([]string, 1000)
	sourceLines[0] = "class LargeClass {"
	for i := 1; i < 999; i++ {
		sourceLines[i] = "    // line " + string(rune(i))
	}
	sourceLines[999] = "}"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		utils.FindMatchingBrace(sourceLines, 0)
	}
}
