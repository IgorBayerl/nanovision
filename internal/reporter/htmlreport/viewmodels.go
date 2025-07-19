package htmlreport

import "html/template"

// AngularAssemblyViewModel corresponds to the data structure for window.assemblies.
type AngularAssemblyViewModel struct {
	Name    string                  `json:"name"`
	Classes []AngularClassViewModel `json:"classes"`
}

// AngularClassViewModel corresponds to the data structure for classes within window.assemblies.
type AngularClassViewModel struct {
	Name                      string                             `json:"name"`
	ReportPath                string                             `json:"rp"`
	CoveredLines              int                                `json:"cl"`
	UncoveredLines            int                                `json:"ucl"`
	CoverableLines            int                                `json:"cal"`
	TotalLines                int                                `json:"tl"`
	CoveredBranches           int                                `json:"cb"`
	TotalBranches             int                                `json:"tb"`
	CoveredMethods            int                                `json:"cm"`
	FullyCoveredMethods       int                                `json:"fcm"`
	TotalMethods              int                                `json:"tm"`
	LineCoverageHistory       []float64                          `json:"lch"`
	BranchCoverageHistory     []float64                          `json:"bch"`
	MethodCoverageHistory     []float64                          `json:"mch"`
	FullMethodCoverageHistory []float64                          `json:"mfch"`
	HistoricCoverages         []AngularHistoricCoverageViewModel `json:"hc"`
	Metrics                   map[string]float64                 `json:"metrics,omitempty"`
}

// AngularHistoricCoverageViewModel corresponds to individual historic coverage data points.
type AngularHistoricCoverageViewModel struct {
	ExecutionTime           string  `json:"et"`
	CoveredLines            int     `json:"cl"`
	UncoveredLines          int     `json:"ucl"`
	CoverableLines          int     `json:"cal"`
	TotalLines              int     `json:"tl"`
	LineCoverageQuota       float64 `json:"lcq"`
	CoveredBranches         int     `json:"cb"`
	TotalBranches           int     `json:"tb"`
	BranchCoverageQuota     float64 `json:"bcq"`
	CoveredMethods          int     `json:"cm"`
	FullyCoveredMethods     int     `json:"fcm"`
	TotalMethods            int     `json:"tm"`
	MethodCoverageQuota     float64 `json:"mcq"`
	FullMethodCoverageQuota float64 `json:"mfcq"`
}

// AngularMetricViewModel corresponds to the data structure for window.metrics.
type AngularMetricViewModel struct {
	Name           string `json:"name"`
	Abbreviation   string `json:"abbreviation"` // Corrected from C# model (abbreviation, not abbr)
	ExplanationURL string `json:"explanationUrl"`
}

// AngularCodeElementViewModel represents an item in the "Methods/Properties" sidebar
type AngularCodeElementViewModel struct {
	Name      string   `json:"name"`  // Display name (e.g., MyMethod(...))
	FullName  string   `json:"fname"` // Full unique name
	Type      string   `json:"type"`  // "Method" or "Property"
	FileIndex int      `json:"fidx"`  // Index of the file this element belongs to
	Line      int      `json:"line"`  // First line
	Coverage  *float64 `json:"cov"`   // Coverage quota (nullable percentage)
}

// AngularLineAnalysisViewModel represents the analysis of a single line of code for Angular.
type AngularLineAnalysisViewModel struct {
	LineNumber      int    `json:"ln"`
	LineContent     string `json:"lc"` // Will be empty for now (Phase 2.3)
	Hits            int    `json:"h"`
	LineVisitStatus string `json:"lvs"` // e.g., "covered", "uncovered", "partiallycovered"
	CoveredBranches int    `json:"cb"`
	TotalBranches   int    `json:"tb"`
}

// AngularCodeFileViewModel represents a code file within a class for Angular.
type AngularCodeFileViewModel struct {
	Path                string                             `json:"p"`
	Lines               []AngularLineAnalysisViewModel     `json:"ls"`
	CoveredLines        int                                `json:"cl"`  // File-specific line coverage
	CoverableLines      int                                `json:"cal"` // File-specific
	TotalLines          int                                `json:"tl"`  // File-specific
	MetricsTableHeaders []AngularMetricDefinitionViewModel `json:"mmh"` // Headers for this file's metrics table
	MetricsTableRows    []AngularMethodMetricsViewModel    `json:"mmr"` // Rows for this file's metrics table
	CodeElements        []AngularCodeElementViewModel      `json:"ce"`  // For the "Methods/Properties" sidebar
}

// AngularClassDetailViewModel represents the detailed data for a single class page for Angular.
type AngularClassDetailViewModel struct {
	Class AngularClassViewModel      `json:"class"` // Contains overall class stats
	Files []AngularCodeFileViewModel `json:"files"` // Contains per-file line data, method metrics, code elements

	// Optional: If a separate class-level aggregated metrics table is needed by Angular:
	// ClassMethodMetricsHeaders []string `json:"classMetricsHeaders,omitempty"`
	// ClassMethodMetricsRows    []AngularMethodMetricRowViewModel `json:"classMetricsRows,omitempty"`
}

// AngularRiskHotspotViewModel corresponds to the data structure for window.riskHotspots.
type AngularRiskHotspotViewModel struct {
	Assembly        string                                    `json:"assembly"`        // Corrected from C# (assembly, not ass)
	Class           string                                    `json:"class"`           // Corrected from C# (class, not cls)
	ReportPath      string                                    `json:"reportPath"`      // Corrected from C# (reportPath, not rp)
	MethodName      string                                    `json:"methodName"`      // Corrected from C# (methodName, not meth)
	MethodShortName string                                    `json:"methodShortName"` // Corrected from C# (methodShortName, not methsn)
	FileIndex       int                                       `json:"fileIndex"`       // Corrected from C# (fileIndex, not fi)
	Line            int                                       `json:"line"`
	Metrics         []AngularRiskHotspotStatusMetricViewModel `json:"metrics"`
}

// AngularRiskHotspotStatusMetricViewModel represents a single metric's status for a risk hotspot.
type AngularRiskHotspotStatusMetricViewModel struct {
	Value    float64 `json:"value"` // C# has this as decimal?, let's use float64 in Go
	Exceeded bool    `json:"exceeded"`
}

// AngularRiskHotspotMetricHeaderViewModel corresponds to the data structure for window.riskHotspotMetrics (headers).
type AngularRiskHotspotMetricHeaderViewModel struct {
	Name           string `json:"name"`
	Abbreviation   string `json:"abbreviation,omitempty"` // Make consistent with Angular model
	ExplanationURL string `json:"explanationUrl"`
}

// AngularMetricDefinitionViewModel describes a metric type for table headers
type AngularMetricDefinitionViewModel struct {
	Name           string `json:"name"`           // e.g., "Cyclomatic Complexity"
	ExplanationURL string `json:"explanationUrl"` // URL for the info icon
}

// AngularMethodMetricsViewModel represents a single method's row in the metrics table
type AngularMethodMetricsViewModel struct {
	Name           string   `json:"name"`                     // Display name of the method/property
	FullName       string   `json:"fullName"`                 // NEW: Full unique name (for title attributes, etc.)
	FileIndex      int      `json:"fileIndex"`                // Index of the file (for linking)
	FileIndexPlus1 int      `json:"fileIndexPlus1,omitempty"` // NEW: 1-based index for display
	FileShortPath  string   `json:"fileShortPath,omitempty"`  // NEW: Sanitized file path for href ID
	Line           int      `json:"line"`                     // First line of the method (for linking)
	MetricValues   []string `json:"metricValues"`             // Metric values as strings, in order of headers
	IsProperty     bool     `json:"isProperty"`               // To choose icon (wrench vs cube)
	CoverageQuota  *float64 `json:"coverageQuota"`            // Method's own line coverage quota
}

// ClassDetailData is the top-level struct for the class_detail_layout.gohtml template
type ClassDetailData struct {
	ReportTitle     string
	AppVersion      string // e.g., "5.4.7.0"
	CurrentDateTime string

	Class                                 ClassViewModelForDetail // Specific view model for the class being detailed
	BranchCoverageAvailable               bool
	MethodCoverageAvailable               bool // To decide whether to show PRO version message
	Tag                                   string
	Translations                          map[string]string
	MaximumDecimalPlacesForCoverageQuotas int // Needed for JS if any Angular components on page use it

	// For JS script includes
	AngularCssFile         string
	AngularRuntimeJsFile   string
	AngularPolyfillsJsFile string
	AngularMainJsFile      string
	CombinedAngularJsFile  string

	// For window.* JSON objects
	ClassDetailJSON                    template.JS // This will contain AngularClassDetailViewModel
	AssembliesJSON                     template.JS
	RiskHotspotsJSON                   template.JS
	MetricsJSON                        template.JS
	RiskHotspotMetricsJSON             template.JS
	HistoricCoverageExecutionTimesJSON template.JS
	TranslationsJSON                   template.JS // The map itself, already marshaled
}

// ClassViewModelForDetail holds data for the main class being displayed
type ClassViewModelForDetail struct {
	Name                                   string
	AssemblyName                           string
	Files                                  []FileViewModelForDetail
	IsMultiFile                            bool
	CoveragePercentageForDisplay           string
	CoveragePercentageBarValue             int
	CoveredLines                           int
	UncoveredLines                         int
	CoverableLines                         int
	TotalLines                             int
	CoverageRatioTextForDisplay            string
	BranchCoveragePercentageForDisplay     string
	BranchCoveragePercentageBarValue       int
	CoveredBranches                        int
	TotalBranches                          int
	BranchCoverageRatioTextForDisplay      string
	MethodCoveragePercentageForDisplay     string
	MethodCoveragePercentageBarValue       int
	FullMethodCoveragePercentageForDisplay string
	CoveredMethods                         int
	FullyCoveredMethods                    int
	TotalMethods                           int
	MethodCoverageRatioTextForDisplay      string
	FullMethodCoverageRatioTextForDisplay  string
	MetricsTable                           MetricsTableViewModel
	FilesWithMetrics                       bool
	SidebarElements                        []SidebarElementViewModel
	// Fields for JS data, if needed by Angular components directly via this struct (less likely with server-side template)
	HistoricCoverages         []AngularHistoricCoverageViewModel `json:"hc,omitempty"`
	LineCoverageHistory       []float64                          `json:"lch,omitempty"`
	BranchCoverageHistory     []float64                          `json:"bch,omitempty"`
	MethodCoverageHistory     []float64                          `json:"mch,omitempty"`
	FullMethodCoverageHistory []float64                          `json:"mfch,omitempty"`
	Metrics                   map[string]float64                 `json:"metrics,omitempty"` // Class-level aggregated metrics
}

// FileViewModelForDetail represents a source file within a class for server-side rendering
type FileViewModelForDetail struct {
	Path      string
	ShortPath string // For use in href IDs (sanitized)
	Lines     []LineViewModelForDetail
}

// LineViewModelForDetail represents a single line of code for server-side rendering
type LineViewModelForDetail struct {
	LineNumber      int
	LineContent     string // Raw content, template will escape and handle spaces
	LineVisitStatus string // CSS class: "green", "red", "orange", "gray"
	Hits            string // Formatted hits, or empty for not coverable
	IsBranch        bool
	BranchBarValue  int // For percentagebar CSS class (0-100 for uncovered part)
	Tooltip         string
	DataCoverage    template.JS // JSON string for data-coverage attribute
}

// MetricsTableViewModel holds data for the "Metrics" table
type MetricsTableViewModel struct {
	Headers []AngularMetricDefinitionViewModel // Re-use from existing viewmodels.go if it fits
	Rows    []AngularMethodMetricsViewModel    // Re-use from existing viewmodels.go if it fits
}

// SidebarElementViewModel holds data for the "Methods/Properties" sidebar links
type SidebarElementViewModel struct {
	Name             string // Display name for the link (short, e.g., Method())
	FullName         string // Full cleaned name (e.g., Namespace.MyClass.Method(Params)) for title
	FileShortPath    string // Sanitized file path for href ID
	FileIndexPlus1   int    // 1-based index of the file if class is multi-file
	Line             int    // First line of the method/property
	Icon             string // "cube" for method, "wrench" for property
	CoverageBarValue int    // For percentagebar CSS (0-100 for uncovered part)
	CoverageTitle    string // e.g., "Line coverage: 50% - Namespace.MyClass.Method(Params)"
}

// SummaryPageData is the top-level struct for the summaryPageLayoutTemplate
type SummaryPageData struct {
	ReportTitle     string
	AppVersion      string
	CurrentDateTime string
	Translations    map[string]string // For direct use in template

	SummaryCards            []CardViewModel
	OverallHistoryChartData HistoryChartDataViewModel

	// For JS script includes
	AngularCssFile         string
	AngularRuntimeJsFile   string
	AngularPolyfillsJsFile string
	AngularMainJsFile      string
	CombinedAngularJsFile  string

	// For window.* JSON objects - These should be template.JS
	AssembliesJSON                     template.JS
	RiskHotspotsJSON                   template.JS
	MetricsJSON                        template.JS
	RiskHotspotMetricsJSON             template.JS
	HistoricCoverageExecutionTimesJSON template.JS
	TranslationsJSON                   template.JS

	BranchCoverageAvailable               bool
	MethodCoverageAvailable               bool
	MaximumDecimalPlacesForCoverageQuotas int
	HasRiskHotspots                       bool
	HasAssemblies                         bool
}

// CardViewModel represents a summary card for the Go template
type CardViewModel struct {
	Title                      string
	SubTitle                   string // e.g., "72%"
	SubTitlePercentageBarValue int    // e.g., 27 for 72% coverage (100-72)
	Rows                       []CardRowViewModel
	ProRequired                bool // For the "Method Coverage" card
}

// CardRowViewModel represents a row in a summary card
type CardRowViewModel struct {
	Header    string
	Text      string
	Tooltip   string
	Alignment string // "left" or "right" (or empty for default)
}

// HistoryChartDataViewModel holds data for rendering a history chart with Go templates
type HistoryChartDataViewModel struct {
	Series     bool        // True if there's data to render the chart
	SVGContent string      // Pre-rendered SVG string
	JSONData   template.JS // JSON data for chart interactivity (if custom.js uses it)
}
