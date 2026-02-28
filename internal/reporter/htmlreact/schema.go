package htmlreact

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
	MaxCyclomaticComplexity *lineCoverageDetail        `json:"max_cyclomatic_complexity,omitempty"`

	// Patch / diff-based metrics.
	PatchStatementCoverage *lineCoverageDetail `json:"patch_statement_coverage,omitempty"`
	PatchLineCoverage      *lineCoverageDetail `json:"patch_line_coverage,omitempty"`
	PatchMethodsHit        *methodsHitDetail   `json:"patch_methods_hit,omitempty"`

	Files    int      `json:"files"`
	Folders  int      `json:"folders"`
	Statuses statuses `json:"statuses,omitempty"`
}

type statuses map[string]riskLevel

type fileNode struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	Type          string     `json:"type"`
	Path          string     `json:"path"`
	Children      []fileNode `json:"children,omitempty"`
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

type newLinesCoverage struct {
	Covered int `json:"covered"`
	Total   int `json:"total"`
}

type methodDetail struct {
	Name                  string                  `json:"name"`
	StartLine             int                     `json:"startLine"`
	EndLine               int                     `json:"endLine"`
	Metrics               map[string]methodMetric `json:"metrics"`
	DiffStatus            string                  `json:"diffStatus,omitempty"`
	NewLinesCoverage      *newLinesCoverage       `json:"newLinesCoverage,omitempty"`
	NewStatementCoverage  *newLinesCoverage       `json:"newStatementCoverage,omitempty"`
	NewStatementsCoverage *newLinesCoverage       `json:"newStatementsCoverage,omitempty"`
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
	SubMetrics []subMetric `json:"subMetrics"`
}

type metricDefinitions map[string]metricDefinition

type summaryV1 struct {
	SchemaVersion     int               `json:"schemaVersion"`
	GeneratedAt       string            `json:"generatedAt"`
	ReportID          string            `json:"reportId,omitempty"`
	Title             string            `json:"title"`
	Totals            totals            `json:"totals"`
	Tree              []fileNode        `json:"tree"`
	MetricDefinitions metricDefinitions `json:"metricDefinitions"`
	Metadata          []metadataItem    `json:"metadata,omitempty"`
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
}
