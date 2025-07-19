package golang

import (
	"strings"

	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/language"
	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/model"
	"github.com/fzipp/gocyclo"
)

type GoProcessor struct{}

func NewGoProcessor() language.Processor {
	return &GoProcessor{}
}

func (p *GoProcessor) Name() string {
	return "Go"
}

func (p *GoProcessor) Detect(filePath string) bool {
	return strings.HasSuffix(strings.ToLower(filePath), ".go")
}

func (p *GoProcessor) GetLogicalClassName(rawClassName string) string {
	return rawClassName
}

func (p *GoProcessor) FormatClassName(class *model.Class) string {
	return class.Name
}

func (p *GoProcessor) FormatMethodName(method *model.Method, class *model.Class) string {
	return method.Name + method.Signature
}

func (p *GoProcessor) CategorizeCodeElement(method *model.Method) model.CodeElementType {
	return model.MethodElementType
}

func (p *GoProcessor) IsCompilerGeneratedClass(class *model.Class) bool {
	return false
}

func (p *GoProcessor) CalculateCyclomaticComplexity(filePath string) ([]model.MethodMetric, error) {
	stats := gocyclo.Analyze([]string{filePath}, nil)

	metrics := make([]model.MethodMetric, 0, len(stats))
	for _, s := range stats {
		metric := model.MethodMetric{
			Name: s.FuncName, // e.g., "(MyType).myFunc" or "myFunc"
			Line: s.Pos.Line,
			Metrics: []model.Metric{
				{
					Name:   "Cyclomatic complexity",
					Value:  float64(s.Complexity),
					Status: model.StatusOk,
				},
			},
		}
		metrics = append(metrics, metric)
	}

	return metrics, nil
}
