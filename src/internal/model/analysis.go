package model

// SummaryResult is the top-level analyzed report, similar to C#'s SummaryResult
type SummaryResult struct {
	ParserName      string
	Timestamp       int64
	SourceDirs      []string
	Assemblies      []Assembly
	LinesCovered    int  // Overall
	LinesValid      int  // Overall
	BranchesCovered *int // Overall - Pointer to indicate presence
	BranchesValid   *int // Overall - Pointer to indicate presence
	TotalLines      int  // Grand total physical lines from unique source files
}

type Assembly struct {
	Name            string
	Classes         []Class
	LinesCovered    int
	LinesValid      int
	BranchesCovered *int // Pointer
	BranchesValid   *int // Pointer
	TotalLines      int  // Sum of unique file TotalLines in this assembly
}

type Class struct {
	Name                string
	DisplayName         string
	Files               []CodeFile
	Methods             []Method
	LinesCovered        int
	LinesValid          int
	BranchesCovered     *int // Pointer
	BranchesValid       *int // Pointer
	TotalLines          int
	CoveredMethods      int
	FullyCoveredMethods int
	TotalMethods        int
	Metrics             map[string]float64 // Aggregated metrics (e.g., sum of complexities)
	HistoricCoverages   []HistoricCoverage // Historical coverage data for this class
}

type CodeFile struct {
	Path           string
	Lines          []Line
	CoveredLines   int
	CoverableLines int
	TotalLines     int
	MethodMetrics  []MethodMetric // Metrics for methods within this file
	CodeElements   []CodeElement  // Code elements (methods/properties) in this file
}

type CodeElementType int

const (
	PropertyElementType CodeElementType = iota
	MethodElementType
)

// BranchCoverageDetail provides details about a specific branch on a line.
type BranchCoverageDetail struct {
	Identifier string // Unique identifier for the branch, e.g., "0", "1", "true", "false"
	Visits     int    // Number of times this specific branch was visited
}

type Line struct {
	Number                   int
	Hits                     int
	IsBranchPoint            bool                   // True if the line is a branch point (from XML branch="true")
	Branch                   []BranchCoverageDetail // Details of branches on this line
	ConditionCoverage        string
	Content                  string         // The actual source code content of the line
	CoveredBranches          int            // Number of branches on this line that were covered
	TotalBranches            int            // Total number of branches on this line
	LineCoverageByTestMethod map[string]int // Tracks hits for this line by TestMethod.ID
	LineVisitStatus          LineVisitStatus
}

type CodeElement struct {
	Name          string
	FullName      string // For uniqueness, e.g., with signature
	Type          CodeElementType
	FirstLine     int
	LastLine      int
	CoverageQuota *float64 // Nullable (percentage 0-100)
}

// GetFirstLine implements utils.SortableByLineAndName for CodeElement
func (ce CodeElement) GetFirstLine() int { return ce.FirstLine }

// GetSortableName implements utils.SortableByLineAndName for CodeElement
// For CodeElement, FullName is the cleaned full name, suitable for consistent sorting.
func (ce CodeElement) GetSortableName() string { return ce.FullName }

type Method struct {
	Name          string
	Signature     string
	DisplayName   string
	LineRate      float64
	BranchRate    *float64
	Complexity    float64
	Lines         []Line
	FirstLine     int
	LastLine      int
	MethodMetrics []MethodMetric
}

// GetFirstLine implements utils.SortableByLineAndName for Method
func (m Method) GetFirstLine() int { return m.FirstLine }

// GetSortableName implements utils.SortableByLineAndName for Method
// For Method, DisplayName is the cleaned full name, suitable for consistent sorting.
func (m Method) GetSortableName() string { return m.DisplayName }
