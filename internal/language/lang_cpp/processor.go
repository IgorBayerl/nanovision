package lang_cpp

import (
	"regexp"
	"strings"

	"github.com/IgorBayerl/AdlerCov/internal/language"
	"github.com/IgorBayerl/AdlerCov/internal/model"
	"github.com/IgorBayerl/AdlerCov/internal/utils"
)

// Matches simple standalone functions
// Examples: "void foo()", "int calculate(int x)", "bool isValid(const char* str)"
var standaloneFunction = regexp.MustCompile(
	`^\s*(?:template\s*<[^>]*>\s*)?` + // Optional template
		`(?:static\s+|inline\s+|extern\s+)*` + // Optional storage specifiers
		`([a-zA-Z_]\w*(?:\s*[*&]+)*)\s+` + // Return type with optional pointers/references
		`([a-zA-Z_]\w*)\s*` + // Function name
		`\([^)]*\)\s*` + // Parameters
		`(?:const\s*)?` + // Optional const
		`\{`, // Must have opening brace
)

// Matches class methods defined inside class body
// Examples: "void method()", "int getValue() const", "static bool check()"
var classMethod = regexp.MustCompile(
	`^\s*(?:public|private|protected)?\s*:?\s*` + // Optional access specifier
		`(?:static\s+|virtual\s+|inline\s+)*` + // Optional method specifiers
		`([a-zA-Z_]\w*(?:\s*[*&]+)*)\s+` + // Return type
		`([a-zA-Z_~]\w*)\s*` + // Method name (including destructors)
		`\([^)]*\)\s*` + // Parameters
		`(?:const\s+|override\s+|final\s+|noexcept\s+)*` + // Optional trailing specifiers
		`\{`, // Must have opening brace
)

// Matches scoped method definitions (defined outside class)
// Examples: "MyClass::method()", "namespace::Class::function()", "MyClass::MyClass()"
var scopedMethod = regexp.MustCompile(
	`^\s*(?:template\s*<[^>]*>\s*)?` + // Optional template
		`(?:inline\s+)?` + // Optional inline
		`([a-zA-Z_]\w*(?:\s*[*&]+)*)\s+` + // Return type
		`([a-zA-Z_]\w*(?:<[^>]*>)?(?:::[a-zA-Z_~]\w*(?:<[^>]*>)?)+)\s*` + // Scoped method name with template support
		`\([^)]*\)\s*` + // Parameters
		`(?:const\s+|noexcept\s+)*` + // Optional trailing specifiers
		`\{`, // Must have opening brace
)

// Matches functions with complex return types (templates, namespaces)
// Examples: "std::vector<int> getData()", "const MyClass& getRef()", "auto getValue() -> int"
var complexReturnFunction = regexp.MustCompile(
	`^\s*(?:template\s*<[^>]*>\s*)?` + // Optional template
		`(?:static\s+|inline\s+)*` + // Optional specifiers
		`((?:const\s+)?(?:\w+::)*\w+(?:\s*<[^>]*>)?\s*[*&]*)\s+` + // Complex return type
		`([a-zA-Z_]\w*)\s*` + // Function name
		`\([^)]*\)\s*` + // Parameters
		`(?:const\s+|noexcept\s+|override\s+)*` + // Optional trailing specifiers
		`(?:->.*?)?\s*` + // Optional trailing return type
		`\{`, // Must have opening brace
)

// Matches constructors (no return type)
// Examples: "MyClass()", "explicit MyClass(int x)", "MyClass::MyClass()"
var constructor = regexp.MustCompile(
	`^\s*(?:explicit\s+)?` + // Optional explicit keyword
		`([a-zA-Z_]\w*(?:::[a-zA-Z_]\w*)*)\s*` + // Constructor name (possibly scoped)
		`\([^)]*\)\s*` + // Parameters
		`(?::\s*[^{]*)?` + // Optional initializer list
		`\{`, // Must have opening brace
)

// Matches operator overloads
// Examples: "bool operator==(const MyClass& other)", "MyClass& operator++()"
var operatorOverload = regexp.MustCompile(
	`^\s*(?:virtual\s+|static\s+|inline\s+|friend\s+)*` + // Optional specifiers
		`([a-zA-Z_][\w:]*(?:\s*[*&]+)*)\s+` + // Return type (allow :: in return type)
		`(operator(?:[+\-*/%=<>!&|^~]+|<<|>>|\[\]|\(\)|new|delete))\s*` + // Operator name
		`\([^)]*\)\s*` + // Parameters
		`(?:const\s+|noexcept\s+|override\s+)*` + // Optional trailing specifiers
		`\{`, // Must have opening brace
)

// Matches destructors specifically
var destructor = regexp.MustCompile(
	`^\s*(?:virtual\s+)?` + // Optional virtual
		`(~[a-zA-Z_]\w*)\s*` + // Destructor name
		`\([^)]*\)\s*` + // Parameters (usually empty)
		`(?:override\s+|noexcept\s+)*` + // Optional trailing specifiers
		`\{`, // Must have opening brace
)

type CppProcessor struct{}

func NewCppProcessor() language.Processor {
	return &CppProcessor{}
}

func (p *CppProcessor) Name() string {
	return "C/C++"
}

func (p *CppProcessor) Detect(filePath string) bool {
	lowerPath := strings.ToLower(filePath)
	extensions := []string{".c", ".cpp", ".h", ".hpp", ".cxx", ".hxx", ".cc"}
	for _, ext := range extensions {
		if strings.HasSuffix(lowerPath, ext) {
			return true
		}
	}
	return false
}

type RegexMatcher struct {
	regex *regexp.Regexp
	name  string
}

func (p *CppProcessor) AnalyzeFile(filePath string, sourceLines []string) ([]model.MethodMetrics, error) {
	methods := make([]model.MethodMetrics, 0) // Always initialize as empty slice, never nil
	if len(sourceLines) == 0 {
		return methods, nil // Return empty slice, not nil
	}

	processedLines := make(map[int]bool)

	// List of regex matchers to try in order
	matchers := []RegexMatcher{
		{destructor, "destructor"},
		{scopedMethod, "scoped method"},
		{constructor, "constructor"},
		{operatorOverload, "operator overload"},
		{complexReturnFunction, "complex return function"},
		{classMethod, "class method"},
		{standaloneFunction, "standalone function"},
	}

	for i, line := range sourceLines {
		lineNum := i + 1
		if processedLines[lineNum] {
			continue
		}

		// Skip obvious non-function lines
		trimmedLine := strings.TrimSpace(line)
		if p.shouldSkipLine(trimmedLine) {
			continue
		}

		// Check for multiline function signatures
		combinedLine := p.getCombinedLine(sourceLines, i)

		// Try each regex matcher
		for _, matcher := range matchers {
			if match := matcher.regex.FindStringSubmatch(combinedLine); match != nil {
				var methodName string

				// Extract method name based on the matcher type
				switch matcher.name {
				case "destructor":
					if len(match) >= 2 {
						methodName = strings.TrimSpace(match[1])
					}
				case "constructor":
					methodName = strings.TrimSpace(match[1])
				case "scoped method", "complex return function", "class method", "standalone function", "operator overload":
					if len(match) >= 3 {
						methodName = strings.TrimSpace(match[2])
					}
				}

				if methodName == "" || p.isInvalidMethodName(methodName) {
					continue
				}

				// Find the actual start line where the function begins
				startLine := p.findFunctionStartLine(sourceLines, i, combinedLine)
				endLine, found := p.findMethodEnd(sourceLines, startLine-1)

				if found && endLine >= startLine {
					methods = append(methods, model.MethodMetrics{
						Name:      methodName,
						StartLine: startLine,
						EndLine:   endLine,
					})

					// Mark all lines of this method as processed
					for j := startLine; j <= endLine; j++ {
						processedLines[j] = true
					}

					// Break after first match to avoid multiple matches on same line
					break
				}
			}
		}
	}

	return methods, nil
}

func (p *CppProcessor) findFunctionStartLine(sourceLines []string, currentIndex int, combinedLine string) int {
	// Check if the current line contains a template declaration
	currentLine := strings.TrimSpace(sourceLines[currentIndex])

	// If current line starts with template, we need to find the actual function signature
	if strings.HasPrefix(currentLine, "template") {
		// Look for the next line that contains the actual function signature
		for i := currentIndex + 1; i < len(sourceLines) && i < currentIndex+5; i++ {
			nextLine := strings.TrimSpace(sourceLines[i])
			// Skip empty lines and comments
			if nextLine == "" || strings.HasPrefix(nextLine, "//") || strings.HasPrefix(nextLine, "/*") {
				continue
			}
			// This should be the function signature line
			return i + 1
		}
	}

	// For all other cases, the function signature starts at the current line
	// The regex matching already ensures we found a valid function signature
	return currentIndex + 1
}

// getCombinedLine combines the current line with following lines for multiline signatures
func (p *CppProcessor) getCombinedLine(sourceLines []string, startIndex int) string {
	if startIndex >= len(sourceLines) {
		return ""
	}

	combined := sourceLines[startIndex]

	// If line already contains opening brace, return as is
	if strings.Contains(combined, "{") {
		return combined
	}

	// Look ahead up to 10 lines to find the opening brace (increased for complex multiline signatures)
	for i := startIndex + 1; i < len(sourceLines) && i < startIndex+10; i++ {
		nextLine := strings.TrimSpace(sourceLines[i])
		combined += " " + nextLine

		// Stop if we find opening brace
		if strings.Contains(nextLine, "{") {
			break
		}

		// Stop if we find semicolon (declaration only)
		if strings.HasSuffix(nextLine, ";") {
			break
		}
	}

	return combined
}

// shouldSkipLine determines if a line should be skipped during analysis
func (p *CppProcessor) shouldSkipLine(trimmedLine string) bool {
	if trimmedLine == "" {
		return true
	}

	skipPrefixes := []string{
		"//", "#", "/*", "*/",
		"struct ", "enum ", "typedef ",
		"namespace ", "using ", "extern ",
		"public:", "private:", "protected:",
	}

	for _, prefix := range skipPrefixes {
		if strings.HasPrefix(trimmedLine, prefix) {
			return true
		}
	}

	// Don't skip class declarations as they might contain inline methods
	if strings.HasPrefix(trimmedLine, "class ") && strings.HasSuffix(trimmedLine, ";") {
		return true
	}

	// Skip variable declarations (simple heuristic)
	if strings.Contains(trimmedLine, "=") && !strings.Contains(trimmedLine, "(") && !strings.Contains(trimmedLine, "operator") {
		return true
	}

	// Correctly handle comments before checking for function declarations.
	lineWithoutComment := trimmedLine
	if commentIndex := strings.Index(lineWithoutComment, "//"); commentIndex != -1 {
		lineWithoutComment = lineWithoutComment[:commentIndex]
	}
	lineWithoutComment = strings.TrimSpace(lineWithoutComment)

	if strings.HasSuffix(lineWithoutComment, ";") && strings.Contains(lineWithoutComment, "(") && !strings.Contains(lineWithoutComment, "{") {
		return true
	}

	return false
}

// isInvalidMethodName checks if the extracted method name is valid
func (p *CppProcessor) isInvalidMethodName(name string) bool {
	// Skip C++ keywords that might be incorrectly matched
	keywords := []string{
		"if", "for", "while", "switch", "return", "break", "continue",
		"int", "char", "float", "double", "void", "bool", "auto",
		"const", "static", "virtual", "inline", "public", "private", "protected",
	}

	for _, keyword := range keywords {
		if name == keyword {
			return true
		}
	}

	// Skip if name contains invalid characters for C++ identifiers
	if strings.Contains(name, " ") && !strings.Contains(name, "::") && !strings.HasPrefix(name, "operator") {
		return true
	}

	return false
}

// findMethodEnd locates the end of a method by finding matching braces
func (p *CppProcessor) findMethodEnd(sourceLines []string, currentIndex int) (int, bool) {
	// Look for opening brace starting from current line
	braceLineIndex := -1

	for j := currentIndex; j < len(sourceLines) && j < currentIndex+10; j++ {
		if strings.Contains(sourceLines[j], "{") {
			braceLineIndex = j
			break
		}

		// If we hit a semicolon, this is just a declaration
		if strings.Contains(strings.TrimSpace(sourceLines[j]), ";") {
			return 0, false
		}
	}

	if braceLineIndex == -1 {
		return 0, false
	}

	return utils.FindMatchingBrace(sourceLines, braceLineIndex)
}
