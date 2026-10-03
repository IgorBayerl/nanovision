package htmlreact

import (
	"encoding/json"
	"maps"

	"github.com/IgorBayerl/nanovision/internal/aggregator"
	"github.com/IgorBayerl/nanovision/internal/model"
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

// a metric that counts methods
type countDetail struct {
	Covered    int     `json:"covered"`
	Total      int     `json:"total"`
	Percentage float64 `json:"percentage"`
}

// a scalar metric with no percentage, e.g. max cyclomatic complexity
type scoreDetail struct {
	Value float64 `json:"value"`
}

type metricsMap map[string]any

// totals is the metrics of the whole report or of one file, next to its
// counts. In JSON the metrics sit beside "files", "folders" and "statuses".
type totals struct {
	Metrics  metricsMap
	Files    int
	Folders  int
	Statuses statuses
}

func (t totals) MarshalJSON() ([]byte, error) {
	out := make(map[string]any, len(t.Metrics)+3)
	maps.Copy(out, t.Metrics)
	out["files"], out["folders"] = t.Files, t.Folders
	if len(t.Statuses) > 0 {
		out["statuses"] = t.Statuses
	}
	return json.Marshal(out)
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
	// a config file of the run, listed in its folder; it has no metrics
	Config bool `json:"config,omitempty"`
}

type lineStatus string

const (
	StatusCovered      lineStatus = "covered"
	StatusUncovered    lineStatus = "uncovered"
	StatusNotCoverable lineStatus = "not-coverable"
)

type lineDetail struct {
	LineNumber int        `json:"lineNumber"`
	Content    string     `json:"content"`
	Status     lineStatus `json:"status"`
	Hits       []int      `json:"hits,omitempty"`
	DiffStatus string     `json:"diffStatus,omitempty"`
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

type folderBands struct {
	Path  string                `json:"path"`
	Bands map[string]statusBand `json:"bands"`
}

// configFile is one source of settings: the root config file (Path ""), or
// the settings of one folder.
type configFile struct {
	Path   string `json:"path"`   // the folder it applies to
	Source string `json:"source"` // the file that holds the settings
}

// MetadataItem is one line of the information block of a report page.
type MetadataItem struct {
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
	// file_metrics in configured order; the UI lists metrics in this order
	MetricOrder []string       `json:"metricOrder,omitempty"`
	Metadata    []MetadataItem `json:"metadata,omitempty"`
	// URL query string applied on first load, e.g. "diff=changed&risk=danger"
	DefaultFilters string `json:"defaultFilters,omitempty"`
	// the verdict on the changed code, for a run measured with a diff
	Review *review.Result `json:"review,omitempty"`
	// the delta against the base run, when one was found
	Comparison *model.Comparison `json:"comparison,omitempty"`
	// the revisions being compared, as rows for the side panel
	Comparing []MetadataItem `json:"comparing,omitempty"`
	// the nodes hold only the changed files; the totals cover the whole run
	OnlyChanged bool `json:"onlyChanged,omitempty"`
	// every parsed report, indexed exactly as the masks in ReportIndexes
	Reports []report `json:"reports,omitempty"`
	// file path -> compressed per-report coverage; absent when a single report
	// was parsed, or when there are too many to address with a bitmask
	ReportIndexes map[string]reportIndex `json:"reportIndexes,omitempty"`
	StatusBands   map[string]statusBand  `json:"statusBands,omitempty"`
	// folders with their own warning ranges, outer folders first
	FolderBands []folderBands `json:"folderBands,omitempty"`
	// the config files of the run; absent when no folder has its own settings
	Configs []configFile `json:"configs,omitempty"`
}

type detailsV1 struct {
	SchemaVersion     int               `json:"schemaVersion"`
	GeneratedAt       string            `json:"generatedAt"`
	Title             string            `json:"title"`
	FileName          string            `json:"fileName"`
	Metadata          []MetadataItem    `json:"metadata"`
	Totals            totals            `json:"totals"`
	MetricDefinitions metricDefinitions `json:"metricDefinitions"`
	// file_metrics in configured order; the UI lists metrics in this order
	MetricOrder []string       `json:"metricOrder,omitempty"`
	Methods     []methodDetail `json:"methods,omitempty"`
	Lines       []lineDetail   `json:"lines"`
	Reports     []report       `json:"reports,omitempty"`
	// compressed per-report coverage for this file, keyed by metric
	ReportIndex reportIndex           `json:"reportIndex,omitempty"`
	StatusBands map[string]statusBand `json:"statusBands,omitempty"`
	// URL query string applied on first load, e.g. "diff=changed&risk=danger"
	DefaultFilters string `json:"defaultFilters,omitempty"`
}
