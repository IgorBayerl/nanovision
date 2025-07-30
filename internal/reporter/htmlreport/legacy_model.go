// Path: internal/reporter/htmlreport/legacy_model.go
package htmlreport

// This file contains a snapshot of the legacy data model. These types are
// deprecated and should only be used by the Tree-to-Legacy adapter to maintain
// compatibility with the existing Angular HTML report. Do not use these types
// in any new development.

// SummaryResult is the top-level analyzed report.
type SummaryResult struct {
	ParserName      string
	Timestamp       int64
	Assemblies      []Assembly
	LinesCovered    int
	LinesValid      int
	BranchesCovered *int
	BranchesValid   *int
	TotalLines      int
}

// Assembly represents a logical grouping of classes.
type Assembly struct {
	Name            string
	Classes         []Class
	LinesCovered    int
	LinesValid      int
	BranchesCovered *int // Pointer
	BranchesValid   *int // Pointer
	TotalLines      int  // Sum of unique file TotalLines in this assembly
}

// Class represents a single class or a logical group of methods.
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

// CodeFile represents a single source file's coverage data.
type CodeFile struct {
	Path           string
	Lines          []Line
	CoveredLines   int
	CoverableLines int
	TotalLines     int
	MethodMetrics  []MethodMetric // Metrics for methods within this file
	CodeElements   []CodeElement  // Code elements (methods/properties) in this file
}

// LineVisitStatus indicates the coverage status of a line.
type LineVisitStatus int

const (
	NotCoverable LineVisitStatus = iota
	NotCovered
	PartiallyCovered
	Covered
)

// Line represents a single line of code.
type Line struct {
	Number          int
	Hits            int
	IsBranchPoint   bool
	CoveredBranches int
	TotalBranches   int
	LineVisitStatus LineVisitStatus
	Content         string
}

// CodeElement represents a method or property.
type CodeElement struct {
	Name          string
	FullName      string // For uniqueness, e.g., with signature
	Type          CodeElementType
	FirstLine     int
	LastLine      int
	CoverageQuota *float64 // Nullable (percentage 0-100)
}

// Method represents a function or method.
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

// MetricStatus represents the status of a metric.
type MetricStatus int

const (
	// StatusOk indicates a normal/good status.
	StatusOk MetricStatus = iota
	// StatusWarning indicates a warning status.
	StatusWarning
	// StatusError indicates an error/critical status.
	StatusError
)

// Metric represents a single metric with a name, value, and status.
type Metric struct {
	Name   string
	Value  interface{} // To allow for different types of metric values (int, float, string)
	Status MetricStatus
}

// MethodMetric represents metrics associated with a specific method or code line.
type MethodMetric struct {
	Name    string   // Typically the method's name or a specific metric name for that method
	Line    int      // The line number where the method is defined or this metric applies
	Metrics []Metric // A slice of Metric structs associated with this method/entry
}

type CodeElementType int

const (
	PropertyElementType CodeElementType = iota
	MethodElementType
)

type HistoricCoverage struct {
	ExecutionTime   int64
	Tag             string
	CoveredLines    int
	CoverableLines  int
	TotalLines      int
	CoveredBranches int
	TotalBranches   int
}
