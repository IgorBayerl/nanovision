package lang_default

import (
	"github.com/IgorBayerl/AdlerCov/internal/language"
	"github.com/IgorBayerl/AdlerCov/internal/model"
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

func (p *DefaultProcessor) AnalyzeFile(filePath string, sourceLines []string) ([]model.MethodMetrics, error) {
	return []model.MethodMetrics{}, nil
}
