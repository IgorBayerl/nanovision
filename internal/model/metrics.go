package model

import "github.com/IgorBayerl/nanovision/internal/config"

// CoverageDetail represents a standard metric expressed as a percentage,
// tracking covered out of total possible.
type CoverageDetail struct {
	Percentage float64 `json:"percentage"`
	Covered    int     `json:"covered"`
	Uncovered  int     `json:"uncovered"`
	Total      int     `json:"total"`
}

// ScoreDetail represents a standalone score (e.g. Max Cyclomatic Complexity).
type ScoreDetail struct {
	Value float64 `json:"value"`
}

// CoverageMetrics holds the aggregated coverage data for a node (project, dir, or file).
type CoverageMetrics struct {
	LinesCovered    int
	LinesValid      int
	BranchesCovered int
	BranchesValid   int
	TotalLines      int

	MethodsHit          int
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
	// coverable line. PatchMethodsHit counts how many of those methods have
	// at least one changed line that was executed at least once.
	PatchMethodsHit   int
	PatchMethodsValid int

	StatementsCovered      int
	StatementsValid        int
	PatchStatementsCovered int
	PatchStatementsValid   int

	StatementMethodsHit          int
	StatementMethodsFullyCovered int
	StatementMethodsValid        int
	PatchStatementMethodsHit     int
	PatchStatementMethodsValid   int

	MaxCyclomaticComplexity int

	// Calculated contains all computed metrics, populated by the calculation pipeline.
	Calculated map[config.MetricKey]any `json:"calculated,omitempty"`
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

	// Patch-specific metrics directly on the method
	PatchLinesValid        int
	PatchLinesCovered      int
	PatchStatementsValid   int
	PatchStatementsCovered int
	DiffStatus             string // "added", "modified", or ""

	Statuses map[config.MetricKey]string // Per-method risk statuses

	// Calculated contains all computed metrics for this method.
	Calculated map[config.MetricKey]any `json:"calculated,omitempty"`
}
