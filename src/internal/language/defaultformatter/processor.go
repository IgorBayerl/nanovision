package defaultformatter

import (
	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/language"
	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/model"
)

type DefaultProcessor struct{}

func NewDefaultProcessor() language.Processor {
	return &DefaultProcessor{}
}

func (p *DefaultProcessor) Name() string {
	return "Default"
}

func (p *DefaultProcessor) Detect(filePath string) bool {
	return false
}

func (p *DefaultProcessor) GetLogicalClassName(rawClassName string) string {
	return rawClassName
}

func (p *DefaultProcessor) FormatClassName(class *model.Class) string {
	return class.Name
}

func (p *DefaultProcessor) FormatMethodName(method *model.Method, class *model.Class) string {
	return method.Name + method.Signature
}

func (p *DefaultProcessor) CategorizeCodeElement(method *model.Method) model.CodeElementType {
	return model.MethodElementType
}

func (p *DefaultProcessor) IsCompilerGeneratedClass(class *model.Class) bool {
	return false
}

func (p *DefaultProcessor) CalculateCyclomaticComplexity(filePath string) ([]model.MethodMetric, error) {
	return nil, language.ErrNotSupported
}
