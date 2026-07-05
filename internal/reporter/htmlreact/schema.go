package htmlreact

import "github.com/IgorBayerl/nanovision/internal/diagnostics"

type riskLevel string

const (
	RiskSafe    riskLevel = "safe"
	RiskWarning riskLevel = "warning"
	RiskDanger  riskLevel = "danger"
)

type lineCoverageDetail struct {
	Covered    int     `json:"covered"`
	Uncovered  int     `json:"uncovered"`
	Coverable  int     `json:"coverable"`
	Total      int     `json:"total"`
	Percentage float64 `json:"percentage"`
}

type branchCoverageDetail struct {
	Covered    int     `json:"covered"`
	Total      int     `json:"total"`
	Percentage float64 `json:"percentage"`
}

type methodsHitDetail struct {
	Covered    int     `json:"covered"`
	Total      int     `json:"total"`
	Percentage float64 `json:"percentage"`
}

type methodsFullyCoveredDetail struct {
	Covered    int     `json:"covered"`
	Total      int     `json:"total"`
	Percentage float64 `json:"percentage"`
}

// scoreDetail represents a standalone scalar metric (e.g. Max Cyclomatic
// Complexity). Unlike the coverage details it has no percentage; the UI renders
// it with a "value" card and treats it as a numeric column.
type scoreDetail struct {
	Value float64 `json:"value"`
}

// UI specific method metric keys to enforce alphabetical sorting
const (
	MethodUIStmtCoverage         = "a_statement_coverage"
	MethodUILineCoverage         = "b_line_coverage"
	MethodUIPatchStmtCoverage    = "c_patch_statement_coverage"
	MethodUIPatchLineCoverage    = "d_patch_line_coverage"
	MethodUIBranchCoverage       = "e_branch_coverage"
	MethodUICyclomaticComplexity = "f_cyclomatic_complexity"
)

type metricsMap map[string]any

type totals struct {
	StatementCoverage       *lineCoverageDetail        `json:"statement_coverage,omitempty"`
	LineCoverage            *lineCoverageDetail        `json:"line_coverage,omitempty"`
	BranchCoverage          *branchCoverageDetail      `json:"branch_coverage,omitempty"`
	MethodsHit              *methodsHitDetail          `json:"methods_hit,omitempty"`
	MethodsFullyCovered     *methodsFullyCoveredDetail `json:"methods_fully_covered,omitempty"`
	MethodBranchCoverage    *branchCoverageDetail      `json:"method_branch_coverage,omitempty"`
	MaxCyclomaticComplexity *scoreDetail               `json:"max_cyclomatic_complexity,omitempty"`

	// Patch / diff-based metrics.
	PatchStatementCoverage *lineCoverageDetail `json:"patch_statement_coverage,omitempty"`
	PatchLineCoverage      *lineCoverageDetail `json:"patch_line_coverage,omitempty"`
	PatchMethodsHit        *methodsHitDetail   `json:"patch_methods_hit,omitempty"`

	Files    int      `json:"files"`
	Folders  int      `json:"folders"`
	Statuses statuses `json:"statuses,omitempty"`
}

type statuses map[string]riskLevel

// fileNode is a single entry in the flat node list emitted to the UI.
//
// The report is delivered as a pre-order (depth-first) flat slice rather than a
// nested tree: each node carries its ParentID and Depth so the client can rebuild
// parent/child relationships in a single linear pass instead of walking a tree.
// This keeps client-side filtering and virtualization O(n).
type fileNode struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	Type          string     `json:"type"`
	Path          string     `json:"path"`
	ParentID      string     `json:"parentId,omitempty"` // "" for top-level nodes
	Depth         int        `json:"depth"`              // structural depth, root children = 0
	Metrics       metricsMap `json:"metrics,omitempty"`
	Statuses      statuses   `json:"statuses,omitempty"`
	ComponentID   string     `json:"componentId,omitempty"`
	TargetURL     string     `json:"targetUrl,omitempty"`
	DiffStatus    string     `json:"diffStatus,omitempty"`
	ComponentName string     `json:"componentName,omitempty"`
}

type lineStatus string

const (
	StatusCovered      lineStatus = "covered"
	StatusUncovered    lineStatus = "uncovered"
	StatusNotCoverable lineStatus = "not-coverable"
	StatusPartial      lineStatus = "partial"
)

type branchInfo struct {
	Covered int `json:"covered"`
	Total   int `json:"total"`
}

type lineDetail struct {
	LineNumber int         `json:"lineNumber"`
	Content    string      `json:"content"`
	Status     lineStatus  `json:"status"`
	Hits       []int       `json:"hits,omitempty"`
	BranchInfo *branchInfo `json:"branchInfo,omitempty"`
	DiffStatus string      `json:"diffStatus,omitempty"`
}

type methodMetric struct {
	Value  string    `json:"value"`
	Status riskLevel `json:"status,omitempty"`
}

type methodDetail struct {
	Name       string                  `json:"name"`
	StartLine  int                     `json:"startLine"`
	EndLine    int                     `json:"endLine"`
	Metrics    map[string]methodMetric `json:"metrics"`
	DiffStatus string                  `json:"diffStatus,omitempty"`
}

type report struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type metadataItem struct {
	Label    string `json:"label"`
	Value    any    `json:"value"`
	SizeHint string `json:"sizeHint,omitempty"`
}

type subMetric struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Width int    `json:"width"`
}

type metricDefinition struct {
	Label      string      `json:"label"`
	ShortLabel string      `json:"shortLabel,omitempty"`
	// Kind selects how the UI renders this metric. "" or "percentage" (default)
	// renders a percentage card with a progress bar; "value" renders a plain
	// scalar number card (e.g. cyclomatic complexity, CRAP score).
	Kind       string      `json:"kind,omitempty"`
	SubMetrics []subMetric `json:"subMetrics"`
}

type metricDefinitions map[string]metricDefinition

type summaryV1 struct {
	SchemaVersion     int               `json:"schemaVersion"`
	GeneratedAt       string            `json:"generatedAt"`
	ReportID          string            `json:"reportId,omitempty"`
	Title             string            `json:"title"`
	Totals            totals            `json:"totals"`
	Nodes             []fileNode        `json:"nodes"`
	MetricDefinitions metricDefinitions `json:"metricDefinitions"`
	Metadata          []metadataItem    `json:"metadata,omitempty"`
	// Diagnostics is the flat list of editor-style problems (coverage
	// warnings/errors) produced by the central diagnostics engine. The UI
	// renders these in a collapsible "Problems" panel.
	Diagnostics []diagnostics.Diagnostic `json:"diagnostics,omitempty"`
	// DefaultFilters is a raw URL query string (e.g. "diff=changed&risk=danger")
	// that the UI applies on first load when no query string is already present.
	DefaultFilters string `json:"defaultFilters,omitempty"`
}

type detailsV1 struct {
	SchemaVersion     int               `json:"schemaVersion"`
	GeneratedAt       string            `json:"generatedAt"`
	Title             string            `json:"title"`
	FileName          string            `json:"fileName"`
	Metadata          []metadataItem    `json:"metadata"`
	Totals            totals            `json:"totals"`
	MetricDefinitions metricDefinitions `json:"metricDefinitions"`
	Methods           []methodDetail    `json:"methods,omitempty"`
	Lines             []lineDetail      `json:"lines"`
	Reports           []report          `json:"reports,omitempty"`
	// DefaultFilters is a raw URL query string applied by the UI on first load
	// when no query string is already present.
	DefaultFilters string `json:"defaultFilters,omitempty"`
}
