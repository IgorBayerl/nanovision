package model

// CoverageMetrics holds the aggregated coverage data for a node (project, dir, or file).
type CoverageMetrics struct {
	LinesCovered    int
	LinesValid      int // Coverable lines
	BranchesCovered int
	BranchesValid   int
}

// LineMetrics holds the coverage data for a single line of code.
type LineMetrics struct {
	Hits            int
	CoveredBranches int
	TotalBranches   int
}

// MethodMetrics holds metric data for a single function or method.
type MethodMetrics struct {
	Name                 string
	Line                 int
	CyclomaticComplexity int
	LineRate             float64
	BranchRate           float64
}
