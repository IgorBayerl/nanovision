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

type metricsMap map[string]any

type totals struct {
	// Define fields in the desired output order for the summary cards and columns.
	LineCoverage   *lineCoverageDetail   `json:"lineCoverage,omitempty"`
	BranchCoverage *branchCoverageDetail `json:"branchCoverage,omitempty"`

	// Non-metric fields last
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
	ComponentName string     `json:"componentName,omitempty"`
	TargetURL     string     `json:"targetUrl,omitempty"`
}

type totalsMap map[string]any

type metadataItem struct {
	Label    string `json:"label"`
	Value    any    `json:"value"` // string or []string
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
