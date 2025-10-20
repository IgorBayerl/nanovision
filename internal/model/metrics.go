package model

// CoverageMetrics holds the aggregated coverage data for a node (project, dir, or file).
type CoverageMetrics struct {
	LinesCovered    int
	LinesValid      int
	BranchesCovered int
	BranchesValid   int
	TotalLines      int

	MethodsCovered      int
	MethodsFullyCovered int
	MethodsValid        int
}

// LineMetrics holds the coverage data for a single line of code.
type LineMetrics struct {
	Hits            int
	CoveredBranches int
	TotalBranches   int
}

// MethodMetrics holds all analysis and coverage data for a single function or method.
type MethodMetrics struct {
	Name                 string // e.g., "MyFunction", "(MyType).MyMethod"
	StartLine            int    // The starting line number of the method.
	EndLine              int    // The ending line number of the method.
	CyclomaticComplexity *int   // Is now a pointer.

	// Per-method coverage metrics.
	LinesValid      int
	LinesCovered    int
	BranchesValid   int
	BranchesCovered int
}
