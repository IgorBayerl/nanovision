package cpp

import (
	"regexp"
	"strings"

	"github.com/IgorBayerl/AdlerCov/internal/language"
	"github.com/IgorBayerl/AdlerCov/internal/model"
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
	return strings.HasSuffix(lowerPath, ".c") || strings.HasSuffix(lowerPath, ".cpp") || strings.HasSuffix(lowerPath, ".h") || strings.HasSuffix(lowerPath, ".hpp")
}

func (p *CppProcessor) GetLogicalClassName(rawClassName string) string {
	return rawClassName
}

func (p *CppProcessor) FormatClassName(class *model.Class) string {
	return class.Name
}

func (p *CppProcessor) FormatMethodName(method *model.Method, class *model.Class) string {
	return method.Name + method.Signature
}

// MapCsharpMangledNames reads the source file and maps the mangled names from the gcov report
// to the clean, human-readable names from the source code.
func (p *CppProcessor) MapCsharpMangledNames(methods []model.Method, sourceLines []string) {
	// This improved regex handles:
	// - Optional return types (including templates, pointers, const etc.)
	// - ClassName::FunctionName patterns
	// - Standalone FunctionName patterns (for namespaces)
	// - Complex parameters with templates, const, and references.
	funcRegex := regexp.MustCompile(
		`(?:\s*[\w\s:&<>\*,]+[\s\*&]*\s+)?` + // Optional return type (non-capturing)
			`([\w:]+(?:::\w+)*)` + // Group 1: The full function name (e.g., MyClass::MyFunc, my_func)
			`\s*(\(.*\))`, // Group 2: The full parameter list, including parentheses
	)

	cleanNamesByLine := make(map[int]string)

	for i, line := range sourceLines {
		if match := funcRegex.FindStringSubmatch(line); match != nil {
			cleanName := match[1] + match[2]
			cleanNamesByLine[i+1] = cleanName
		}
	}

	for i := range methods {
		method := &methods[i]
		if cleanName, ok := cleanNamesByLine[method.FirstLine]; ok {
			// Found a clean name. Separate the name from the signature.
			parenIndex := strings.Index(cleanName, "(")
			if parenIndex != -1 {
				method.Name = cleanName[:parenIndex]
				method.Signature = cleanName[parenIndex:]
			} else {
				method.Name = cleanName
				method.Signature = ""
			}
			method.DisplayName = cleanName
		} else {
			// Fallback for safety, though it should not be hit
			method.DisplayName = method.Name
		}
	}
}

func (p *CppProcessor) CategorizeCodeElement(method *model.Method) model.CodeElementType {
	return model.MethodElementType
}

func (p *CppProcessor) IsCompilerGeneratedClass(class *model.Class) bool {
	return false
}

func (p *CppProcessor) CalculateCyclomaticComplexity(filePath string) ([]model.MethodMetric, error) {
	return nil, language.ErrNotSupported
}
