package language

import (
	"errors"

	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/model"
)

// ErrNotSupported is a sentinel error returned by a Processor when a feature
// like cyclomatic complexity is not applicable to that language.
var ErrNotSupported = errors.New("feature not supported for this language")

type Processor interface {
	// Name returns the unique, human-readable name of the processor (e.g., "C#", "Go").
	Name() string

	// Detect checks if this processor should be used for a given source file path.
	Detect(filePath string) bool

	// GetLogicalClassName determines the grouping key for a class from a raw name.
	GetLogicalClassName(rawClassName string) string

	// FormatClassName transforms a raw class name into a display-friendly version.
	FormatClassName(class *model.Class) string

	// FormatMethodName transforms a raw method name and signature into a display-friendly version.
	FormatMethodName(method *model.Method, class *model.Class) string

	// CategorizeCodeElement determines if a method is a standard method, property, etc.
	CategorizeCodeElement(method *model.Method) model.CodeElementType

	// IsCompilerGeneratedClass determines if a class is a compiler-generated artifact
	// that should be filtered out from the final report.
	IsCompilerGeneratedClass(class *model.Class) bool

	// CalculateCyclomaticComplexity analyzes a file and returns the metric for each function.
	// If the language does not support this metric, it must return language.ErrNotSupported.
	CalculateCyclomaticComplexity(filePath string) ([]model.MethodMetric, error)
}

type ProcessorFactory struct {
	processors       []Processor
	defaultProcessor Processor
}

func NewProcessorFactory(processors ...Processor) *ProcessorFactory {
	factory := &ProcessorFactory{
		processors: make([]Processor, 0, len(processors)),
	}

	for _, p := range processors {
		if p.Name() == "Default" {
			factory.defaultProcessor = p
		} else {
			factory.processors = append(factory.processors, p)
		}
	}

	if factory.defaultProcessor == nil {
		panic("FATAL: Default language processor was not provided to the factory.")
	}

	return factory
}

func (f *ProcessorFactory) FindProcessorForFile(filePath string) Processor {
	for _, p := range f.processors {
		if p.Detect(filePath) {
			return p
		}
	}

	return f.defaultProcessor
}
