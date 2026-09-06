package htmlreact

import (
	"github.com/IgorBayerl/nanovision/internal/aggregator"
	"github.com/IgorBayerl/nanovision/internal/diagnostics"
	"github.com/IgorBayerl/nanovision/internal/review"
)

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

// a scalar metric with no percentage, e.g. max cyclomatic complexity
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
	MethodUICrapScore            = "g_crap_score"
	MethodUIPatchCrapScore       = "h_patch_crap_score"
	MethodUIExposedRisk          = "i_exposed_risk"
	MethodUIDefectProbability    = "j_defect_probability"
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

// one entry in the node list.
//
// the list is flat and pre-order, not a tree: ParentID and Depth let a client
// rebuild the hierarchy in one linear pass, which keeps filtering O(n).
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
	// set on details pages only: this report contributed a hit to the file
	Relevant bool `json:"relevant,omitempty"`
}

// reportIndex is the compressed per-report coverage data the UI needs to
// recompute metrics for any subset of reports, keyed by metric.
// See aggregator.BuildFileReportIndex for the bucket semantics.
type reportIndex map[string][]aggregator.ReportBucket

// statusBand mirrors config.Band so the UI can re-classify risk after the
// report selection changes the underlying numbers.
type statusBand struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
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
	Label      string `json:"label"`
	ShortLabel string `json:"shortLabel,omitempty"`
	// one-line explanation of the metric, shown in the UI's hover tooltips
	Description string `json:"description,omitempty"`
	// "percentage" (default) or "value" for plain scalars like complexity
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
	// flat list of problems (coverage warnings/errors)
	Diagnostics []diagnostics.Diagnostic `json:"diagnostics,omitempty"`
	// URL query string applied on first load, e.g. "diff=changed&risk=danger"
	DefaultFilters string `json:"defaultFilters,omitempty"`
	// gate verdict from review.Evaluate; set only for changelist reports
	Review *review.Result `json:"review,omitempty"`
	// every parsed report, indexed exactly as the masks in ReportIndexes
	Reports []report `json:"reports,omitempty"`
	// file path -> compressed per-report coverage; absent when a single report
	// was parsed, or when there are too many to address with a bitmask
	ReportIndexes map[string]reportIndex `json:"reportIndexes,omitempty"`
	StatusBands   map[string]statusBand  `json:"statusBands,omitempty"`
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
	// compressed per-report coverage for this file, keyed by metric
	ReportIndex reportIndex           `json:"reportIndex,omitempty"`
	StatusBands map[string]statusBand `json:"statusBands,omitempty"`
	// URL query string applied on first load, e.g. "diff=changed&risk=danger"
	DefaultFilters string `json:"defaultFilters,omitempty"`
}
