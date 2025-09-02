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

// Represents the metrics for methods that have at least one line covered.
type methodsCoveredDetail struct {
	Covered    int     `json:"covered"`
	Total      int     `json:"total"`
	Percentage float64 `json:"percentage"`
}

// Represents the metrics for methods that have 100% line coverage.
type methodsFullyCoveredDetail struct {
	Covered    int     `json:"covered"`
	Total      int     `json:"total"`
	Percentage float64 `json:"percentage"`
}

type metricsMap map[string]any

type totals struct {
	LineCoverage            *lineCoverageDetail        `json:"lineCoverage,omitempty"`
	BranchCoverage          *branchCoverageDetail      `json:"branchCoverage,omitempty"`
	MethodsCovered          *methodsCoveredDetail      `json:"methodsCovered,omitempty"`
	MethodsFullyCovered     *methodsFullyCoveredDetail `json:"methodsFullyCovered,omitempty"`
	MethodBranchCoverage    *branchCoverageDetail      `json:"methodBranchCoverage,omitempty"`
	MaxCyclomaticComplexity *lineCoverageDetail        `json:"maxCyclomaticComplexity,omitempty"`
	Files                   int                        `json:"files"`
	Folders                 int                        `json:"folders"`
	Statuses                statuses                   `json:"statuses,omitempty"`
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
	ComponentName string     `json:"componentName,omitempty"`
	TargetURL     string     `json:"targetUrl,omitempty"`
}

// REMOVED: type totalsMap map[string]any (was unused)

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

// --- NEW STRUCTS FOR DETAILS PAGE ---

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
	Hits       *int        `json:"hits,omitempty"`
	BranchInfo *branchInfo `json:"branchInfo,omitempty"`
}

type methodMetric struct {
	Value  string    `json:"value"`
	Status riskLevel `json:"status,omitempty"`
}

type methodDetail struct {
	Name      string                  `json:"name"`
	StartLine int                     `json:"startLine"`
	EndLine   int                     `json:"endLine"`
	Metrics   map[string]methodMetric `json:"metrics"`
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
}
