package gcc

import (
	"strings"

	"github.com/IgorBayerl/AdlerCov/internal/language"
	"github.com/IgorBayerl/AdlerCov/internal/model"
)

type GccProcessor struct{}

func NewGccProcessor() language.Processor {
	return &GccProcessor{}
}

func (p *GccProcessor) Name() string {
	return "C/C++"
}

func (p *GccProcessor) Detect(filePath string) bool {
	lowerPath := strings.ToLower(filePath)
	return strings.HasSuffix(lowerPath, ".c") || strings.HasSuffix(lowerPath, ".cpp") || strings.HasSuffix(lowerPath, ".h") || strings.HasSuffix(lowerPath, ".hpp")
}

// For C++, the class name is often the file name without extension.
func (p *GccProcessor) GetLogicalClassName(rawClassName string) string {
	return rawClassName
}

func (p *GccProcessor) FormatClassName(class *model.Class) string {
	return class.Name
}

func (p *GccProcessor) FormatMethodName(method *model.Method, class *model.Class) string {
	return method.Name + method.Signature
}

func (p *GccProcessor) CategorizeCodeElement(method *model.Method) model.CodeElementType {
	return model.MethodElementType
}

func (p *GccProcessor) IsCompilerGeneratedClass(class *model.Class) bool {
	return false
}

func (p *GccProcessor) CalculateCyclomaticComplexity(filePath string) ([]model.MethodMetric, error) {
	return nil, language.ErrNotSupported
}
