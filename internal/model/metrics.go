package model

// CoverageMetrics holds the aggregated coverage data for a node (project, dir, or file).
type CoverageMetrics struct {
	LinesCovered    int
	LinesValid      int
	TotalLines      int
	BranchesCovered int
	BranchesValid   int
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
	CyclomaticComplexity int    // Calculated complexity, 0 if not supported/calculated.

	// These will be calculated by the Hydrator
	LinesValid     int     // Number of coverable lines within this method.
	LinesCovered   int     // Number of covered lines within this method.
	LineCoverage   float64 // (LinesCovered / LinesValid) * 100
	BranchCoverage float64 // Branch coverage for this specific method.
	IsFullyCovered bool    // True if all LinesValid are covered.
}
