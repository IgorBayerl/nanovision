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
	PatchLinesTotal     int

	// Patch / diff-based metrics.
	// PatchLinesValid counts the number of coverable lines that were added or
	// modified in the diff. PatchLinesCovered counts how many of those lines
	// were executed at least once.
	PatchLinesCovered int
	PatchLinesValid   int

	// PatchMethodsValid counts how many methods contain at least one changed,
	// coverable line. PatchMethodsCovered counts how many of those methods have
	// at least one changed line that was executed at least once.
	PatchMethodsCovered int
	PatchMethodsValid   int

	StatementsCovered      int
	StatementsValid        int
	PatchStatementsCovered int
	PatchStatementsValid   int

	StatementMethodsCovered      int
	StatementMethodsFullyCovered int
	StatementMethodsValid        int
	PatchStatementMethodsCovered int
	PatchStatementMethodsValid   int
}

// LineMetrics holds the coverage data for a single line of code.
type LineMetrics struct {
	// Hits is the total hit count aggregated from all coverage reports
	// for this line. A value of -1 means the line is not coverable
	// (e.g., comments, blank lines).
	Hits int

	// ReportHits stores the individual hit count from each report. The index
	// corresponds to the ReportNames slice in the SummaryTree.
	// This is used by reporters that support per-report filtering (e.g., HTML).
	ReportHits []int

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

	StatementsValid   int
	StatementsCovered int
}
