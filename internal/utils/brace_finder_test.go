package utils

import (
	"testing"
)

func TestFindMatchingBrace_SimpleCase_ReturnsCorrectLineNumber(t *testing.T) {
	// Arrange
	sourceLines := []string{
		"function test() {",
		"    return true;",
		"}",
	}
	startLineIndex := 0

	// Act
	lineNumber, found := FindMatchingBrace(sourceLines, startLineIndex)

	// Assert
	if !found {
		t.Error("Expected to find matching brace, but none was found")
	}
	if lineNumber != 3 {
		t.Errorf("Expected line number 3, got %d", lineNumber)
	}
}

func TestFindMatchingBrace_NestedBraces_ReturnsOutermostClosingBrace(t *testing.T) {
	// Arrange
	sourceLines := []string{
		"class Test {",
		"    method() {",
		"        if (condition) {",
		"            doSomething();",
		"        }",
		"    }",
		"}",
	}
	startLineIndex := 0

	// Act
	lineNumber, found := FindMatchingBrace(sourceLines, startLineIndex)

	// Assert
	if !found {
		t.Error("Expected to find matching brace, but none was found")
	}
	if lineNumber != 7 {
		t.Errorf("Expected line number 7, got %d", lineNumber)
	}
}

func TestFindMatchingBrace_StartFromMiddle_FindsCorrectMatch(t *testing.T) {
	// Arrange
	sourceLines := []string{
		"outer {",
		"    inner {",
		"        content",
		"    }",
		"}",
	}
	startLineIndex := 1

	// Act
	lineNumber, found := FindMatchingBrace(sourceLines, startLineIndex)

	// Assert
	if !found {
		t.Error("Expected to find matching brace, but none was found")
	}
	if lineNumber != 4 {
		t.Errorf("Expected line number 4, got %d", lineNumber)
	}
}

func TestFindMatchingBrace_NoOpeningBrace_ReturnsFalse(t *testing.T) {
	// Arrange
	sourceLines := []string{
		"some code",
		"more code",
		"}",
	}
	startLineIndex := 0

	// Act
	_, found := FindMatchingBrace(sourceLines, startLineIndex)

	// Assert
	if found {
		t.Error("Expected not to find matching brace when no opening brace exists")
	}
}

func TestFindMatchingBrace_NoClosingBrace_ReturnsFalse(t *testing.T) {
	// Arrange
	sourceLines := []string{
		"function test() {",
		"    return true;",
		"    // missing closing brace",
	}
	startLineIndex := 0

	// Act
	_, found := FindMatchingBrace(sourceLines, startLineIndex)

	// Assert
	if found {
		t.Error("Expected not to find matching brace when no closing brace exists")
	}
}

func TestFindMatchingBrace_BraceInLineComment_IgnoresBrace(t *testing.T) {
	// Arrange
	sourceLines := []string{
		"function test() {",
		"    // this } should be ignored",
		"    return true;",
		"}",
	}
	startLineIndex := 0

	// Act
	lineNumber, found := FindMatchingBrace(sourceLines, startLineIndex)

	// Assert
	if !found {
		t.Error("Expected to find matching brace, but none was found")
	}
	if lineNumber != 4 {
		t.Errorf("Expected line number 4, got %d", lineNumber)
	}
}

func TestFindMatchingBrace_BraceInBlockComment_IgnoresBrace(t *testing.T) {
	// Arrange
	sourceLines := []string{
		"function test() {",
		"    /* this } should be ignored */",
		"    return true;",
		"}",
	}
	startLineIndex := 0

	// Act
	lineNumber, found := FindMatchingBrace(sourceLines, startLineIndex)

	// Assert
	if !found {
		t.Error("Expected to find matching brace, but none was found")
	}
	if lineNumber != 4 {
		t.Errorf("Expected line number 4, got %d", lineNumber)
	}
}

func TestFindMatchingBrace_MultilineBlockComment_IgnoresBraces(t *testing.T) {
	// Arrange
	sourceLines := []string{
		"function test() {",
		"    /*",
		"     * this } should be ignored",
		"     * and this { too",
		"     */",
		"    return true;",
		"}",
	}
	startLineIndex := 0

	// Act
	lineNumber, found := FindMatchingBrace(sourceLines, startLineIndex)

	// Assert
	if !found {
		t.Error("Expected to find matching brace, but none was found")
	}
	if lineNumber != 7 {
		t.Errorf("Expected line number 7, got %d", lineNumber)
	}
}

func TestFindMatchingBrace_BraceInString_IgnoresBrace(t *testing.T) {
	// Arrange
	sourceLines := []string{
		"function test() {",
		"    let msg = \"this } should be ignored\";",
		"    return true;",
		"}",
	}
	startLineIndex := 0

	// Act
	lineNumber, found := FindMatchingBrace(sourceLines, startLineIndex)

	// Assert
	if !found {
		t.Error("Expected to find matching brace, but none was found")
	}
	if lineNumber != 4 {
		t.Errorf("Expected line number 4, got %d", lineNumber)
	}
}

func TestFindMatchingBrace_EscapedQuoteInString_HandlesCorrectly(t *testing.T) {
	// Arrange
	sourceLines := []string{
		"function test() {",
		"    let msg = \"escaped quote \\\" and brace }\";",
		"    return true;",
		"}",
	}
	startLineIndex := 0

	// Act
	lineNumber, found := FindMatchingBrace(sourceLines, startLineIndex)

	// Assert
	if !found {
		t.Error("Expected to find matching brace, but none was found")
	}
	if lineNumber != 4 {
		t.Errorf("Expected line number 4, got %d", lineNumber)
	}
}

func TestFindMatchingBrace_MultipleBracesOnSameLine_FindsCorrectMatch(t *testing.T) {
	// Arrange
	sourceLines := []string{
		"if (condition) { doSomething(); }",
		"function test() {",
		"    return true;",
		"}",
	}
	startLineIndex := 0

	// Act
	lineNumber, found := FindMatchingBrace(sourceLines, startLineIndex)

	// Assert
	if !found {
		t.Error("Expected to find matching brace, but none was found")
	}
	if lineNumber != 1 {
		t.Errorf("Expected line number 1, got %d", lineNumber)
	}
}

func TestFindMatchingBrace_StartIndexBeyondArray_ReturnsFalse(t *testing.T) {
	// Arrange
	sourceLines := []string{
		"function test() {",
		"}",
	}
	startLineIndex := 5

	// Act
	_, found := FindMatchingBrace(sourceLines, startLineIndex)

	// Assert
	if found {
		t.Error("Expected not to find matching brace when start index is beyond array bounds")
	}
}

func TestFindMatchingBrace_EmptySourceLines_ReturnsFalse(t *testing.T) {
	// Arrange
	sourceLines := []string{}
	startLineIndex := 0

	// Act
	_, found := FindMatchingBrace(sourceLines, startLineIndex)

	// Assert
	if found {
		t.Error("Expected not to find matching brace in empty source lines")
	}
}

func TestFindMatchingBrace_ComplexNestedStructure_FindsCorrectMatch(t *testing.T) {
	// Arrange
	sourceLines := []string{
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
	}
	startLineIndex := 0

	// Act
	lineNumber, found := FindMatchingBrace(sourceLines, startLineIndex)

	// Assert
	if !found {
		t.Error("Expected to find matching brace, but none was found")
	}
	if lineNumber != 15 {
		t.Errorf("Expected line number 15, got %d", lineNumber)
	}
}

func TestFindMatchingBrace_LineCommentAfterBrace_IgnoresComment(t *testing.T) {
	// Arrange
	sourceLines := []string{
		"function test() { // comment with }",
		"    return true;",
		"}",
	}
	startLineIndex := 0

	// Act
	lineNumber, found := FindMatchingBrace(sourceLines, startLineIndex)

	// Assert
	if !found {
		t.Error("Expected to find matching brace, but none was found")
	}
	if lineNumber != 3 {
		t.Errorf("Expected line number 3, got %d", lineNumber)
	}
}

func TestFindMatchingBrace_BlockCommentAfterBrace_IgnoresComment(t *testing.T) {
	// Arrange
	sourceLines := []string{
		"function test() { /* comment with } */",
		"    return true;",
		"}",
	}
	startLineIndex := 0

	// Act
	lineNumber, found := FindMatchingBrace(sourceLines, startLineIndex)

	// Assert
	if !found {
		t.Error("Expected to find matching brace, but none was found")
	}
	if lineNumber != 3 {
		t.Errorf("Expected line number 3, got %d", lineNumber)
	}
}

func TestFindMatchingBrace_OnlyClosingBraces_ReturnsFalse(t *testing.T) {
	// Arrange
	sourceLines := []string{
		"    }",
		"    }",
		"}",
	}
	startLineIndex := 0

	// Act
	_, found := FindMatchingBrace(sourceLines, startLineIndex)

	// Assert
	if found {
		t.Error("Expected not to find matching brace when only closing braces exist")
	}
}

func TestFindMatchingBrace_MixedQuotesAndComments_HandlesCorrectly(t *testing.T) {
	// Arrange
	sourceLines := []string{
		"function test() {",
		"    let str1 = \"quote with // comment\";",
		"    /* block comment with \" quote */",
		"    let str2 = \"another } brace\";",
		"    // line comment with }",
		"}",
	}
	startLineIndex := 0

	// Act
	lineNumber, found := FindMatchingBrace(sourceLines, startLineIndex)

	// Assert
	if !found {
		t.Error("Expected to find matching brace, but none was found")
	}
	if lineNumber != 6 {
		t.Errorf("Expected line number 6, got %d", lineNumber)
	}
}

func TestFindMatchingBrace_BraceAtLineEnd_FindsCorrectMatch(t *testing.T) {
	// Arrange
	sourceLines := []string{
		"if (condition)",
		"{",
		"    doSomething();",
		"}",
	}
	startLineIndex := 1

	// Act
	lineNumber, found := FindMatchingBrace(sourceLines, startLineIndex)

	// Assert
	if !found {
		t.Error("Expected to find matching brace, but none was found")
	}
	if lineNumber != 4 {
		t.Errorf("Expected line number 4, got %d", lineNumber)
	}
}

// Table-driven test for multiple edge cases
func TestFindMatchingBrace_EdgeCases(t *testing.T) {
	testCases := []struct {
		name           string
		sourceLines    []string
		startLineIndex int
		expectedLine   int
		expectedFound  bool
	}{
		{
			name: "SingleLineWithBraces",
			sourceLines: []string{
				"{ code here }",
			},
			startLineIndex: 0,
			expectedLine:   1,
			expectedFound:  true,
		},
		{
			name: "BraceInMiddleOfLine",
			sourceLines: []string{
				"code { more code",
				"    content",
				"} end",
			},
			startLineIndex: 0,
			expectedLine:   3,
			expectedFound:  true,
		},
		{
			name: "MultipleOpeningBraces",
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
			name: "UnbalancedBraces",
			sourceLines: []string{
				"{ { {",
				"content",
				"} }",
			},
			startLineIndex: 0,
			expectedLine:   -1,
			expectedFound:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			lineNumber, found := FindMatchingBrace(tc.sourceLines, tc.startLineIndex)

			// Assert
			if found != tc.expectedFound {
				t.Errorf("Expected found=%v, got found=%v", tc.expectedFound, found)
			}
			if tc.expectedFound && lineNumber != tc.expectedLine {
				t.Errorf("Expected line number %d, got %d", tc.expectedLine, lineNumber)
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
		FindMatchingBrace(sourceLines, 0)
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
		FindMatchingBrace(sourceLines, 0)
	}
}
