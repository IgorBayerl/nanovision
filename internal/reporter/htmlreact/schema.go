package htmlreact

type riskLevel string

const (
	RiskSafe    riskLevel = "safe"
	RiskWarning riskLevel = "warning"
	RiskDanger  riskLevel = "danger"
)

type metrics struct {
	LineCoverage      int `json:"lineCoverage"`
	BranchCoverage    int `json:"branchCoverage"`
	MethodCoverage    int `json:"methodCoverage"`
	StatementCoverage int `json:"statementCoverage"`
	FunctionCoverage  int `json:"functionCoverage"`
}

type statuses map[string]riskLevel // keys mirror Metrics JSON names

type fileNode struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	Type          string     `json:"type"` // "folder" | "file"
	Path          string     `json:"path"`
	Children      []fileNode `json:"children,omitempty"`
	Metrics       *metrics   `json:"metrics,omitempty"`
	Statuses      statuses   `json:"statuses,omitempty"`
	ComponentID   string     `json:"componentId,omitempty"`
	ComponentName string     `json:"componentName,omitempty"`
	TargetURL     string     `json:"targetUrl,omitempty"`
}

type totals struct {
	Files                int `json:"files"`
	Folders              int `json:"folders"`
	LinesCoveredPct      int `json:"linesCoveredPct"`
	BranchesCoveredPct   int `json:"branchesCoveredPct"`
	MethodsCoveredPct    int `json:"methodsCoveredPct"`
	StatementsCoveredPct int `json:"statementsCoveredPct"`
	FunctionsCoveredPct  int `json:"functionsCoveredPct"`
}

type summaryV1 struct {
	SchemaVersion int        `json:"schemaVersion"`
	GeneratedAt   string     `json:"generatedAt"`
	ReportID      string     `json:"reportId,omitempty"`
	Title         string     `json:"title"`
	Totals        totals     `json:"totals"`
	Tree          []fileNode `json:"tree"`
}
