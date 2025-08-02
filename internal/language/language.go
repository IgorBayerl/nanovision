package language

import (
	"errors"

	"github.com/IgorBayerl/AdlerCov/internal/model"
)

var ErrNotSupported = errors.New("feature not supported for this language")

type Processor interface {
	// Name returns the unique, human-readable name of the processor (e.g., "Go", "C#").
	Name() string

	// Detect checks if this processor should be used for a given source file path.
	Detect(filePath string) bool

	// AnalyzeFile performs static analysis on the source code content.
	// It should return a slice of MethodMetrics containing the name, start/end lines,
	// and cyclomatic complexity for each function/method found.
	// Other fields (like coverage) will be populated later by the Hydrator.
	AnalyzeFile(filePath string, sourceLines []string) ([]model.MethodMetrics, error)
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
