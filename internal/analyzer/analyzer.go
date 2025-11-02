package analyzer

import "fmt"

type Position struct {
	StartLine int
	EndLine   int
}

type FunctionMetric struct {
	Name                 string
	Position             Position
	CyclomaticComplexity *int
}

type AnalysisResult struct {
	Functions []FunctionMetric
}

type Analyzer interface {
	// Name returns the human-readable name of the analyzer.
	Name() string
	// SupportsFile returns true if the analyzer can process the given file path.
	SupportsFile(filePath string) bool
	// Analyze takes source code as input and returns structured metrics.
	Analyze(sourceCode []byte) (AnalysisResult, error)
}

type AnalysisError struct {
	FilePath string
	Err      error
}

func (e *AnalysisError) Error() string {
	return fmt.Sprintf("failed to analyze %s: %v", e.FilePath, e.Err)
}
