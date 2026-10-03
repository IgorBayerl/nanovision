package config

// MetricKey is the internal id of a metric: in the report data and in stored
// runs. The config file uses MetricDef.Name, which is unique only within its
// scope.
type MetricKey string

const (
	LineCoverage            MetricKey = "line_coverage"
	MethodsHit              MetricKey = "methods_hit"
	MethodsFullyCovered     MetricKey = "methods_fully_covered"
	PatchLineCoverage       MetricKey = "patch_line_coverage"
	PatchMethodsHit         MetricKey = "patch_methods_hit"
	StatementCoverage       MetricKey = "statement_coverage"
	PatchStatementCoverage  MetricKey = "patch_statement_coverage"
	MaxCyclomaticComplexity MetricKey = "max_cyclomatic_complexity"

	StatementMethodsHit          MetricKey = "statement_methods_hit"
	StatementMethodsFullyCovered MetricKey = "statement_methods_fully_covered"
	PatchStatementMethodsHit     MetricKey = "patch_statement_methods_hit"

	MethodLineCoverage           MetricKey = "method_line_coverage"
	MethodStatementCoverage      MetricKey = "method_statement_coverage"
	MethodPatchLineCoverage      MetricKey = "method_patch_line_coverage"
	MethodPatchStatementCoverage MetricKey = "method_patch_statement_coverage"
	CyclomaticComplexity         MetricKey = "cyclomatic_complexity"
	MethodCrapScore              MetricKey = "method_crap_score"
)

// Needs names the data a metric depends on. A report format that lacks it
// gives the metric no status instead of a misleading 0%.
type Needs string

const (
	NeedsMethods    Needs = "methods"
	NeedsStatements Needs = "statements"
)

// MetricDef describes one metric. The tables below are the only place a
// metric is named and documented: the config validation, the report labels
// and tooltips, the status rules and the generated docs all read them. The
// calculation is in internal/calculator, under the same Key.
type MetricDef struct {
	Key   MetricKey
	Name  string // in the config file, under metrics.files or metrics.methods
	Label string // in the report: side panel and tooltips
	Short string // column header; Label when empty
	Title string // in sentences, e.g. "Statement Coverage is 52%"; Label when empty
	Doc   string
	// a plain number such as complexity, not a percentage
	Value bool
	// a higher number is worse; the warning range then has danger above it
	LowerIsBetter bool
	Needs         Needs
}

// FileMetricDefs are the metrics of a file or folder. Their order is the
// default column order.
var FileMetricDefs = []MetricDef{
	{
		Key:   StatementCoverage,
		Name:  "statement_coverage",
		Label: "Statements",
		Title: "Statement Coverage",
		Doc:   "Statements run by tests."},
	{
		Key:   LineCoverage,
		Name:  "line_coverage",
		Label: "Lines",
		Title: "Line Coverage",
		Doc:   "Lines run by tests."},
	{
		Key:   MethodsHit,
		Name:  "methods_hit",
		Label: "Methods Hit",
		Needs: NeedsMethods,
		Doc:   "Methods entered by tests."},
	{
		Key:   MethodsFullyCovered,
		Name:  "methods_fully_covered",
		Label: "Methods Fully Covered",
		Short: "Fully Covered",
		Needs: NeedsMethods,
		Doc:   "Methods with every line run."},
	{
		Key:   StatementMethodsHit,
		Name:  "statement_methods_hit",
		Label: "Statement Methods Hit",
		Short: "Stmt Methods Hit",
		Needs: NeedsStatements,
		Doc:   "Methods with a statement run."},
	{
		Key:   StatementMethodsFullyCovered,
		Name:  "statement_methods_fully_covered",
		Label: "Statement Methods Fully Covered",
		Short: "Stmt Fully Covered",
		Needs: NeedsStatements,
		Doc:   "Methods with every statement run."},
	{
		Key:   PatchStatementCoverage,
		Name:  "patch_statement_coverage",
		Label: "Patch Statements",
		Title: "Patch Statement Coverage",
		Doc:   "Changed statements run by tests. Needs a diff."},
	{
		Key:   PatchLineCoverage,
		Name:  "patch_line_coverage",
		Label: "Patch Lines",
		Title: "Patch Line Coverage",
		Doc:   "Changed lines run by tests. Needs a diff."},
	{
		Key:   PatchMethodsHit,
		Name:  "patch_methods_hit",
		Label: "Patch Methods Hit",
		Doc:   "Changed methods entered by tests. Needs a diff."},
	{
		Key:   PatchStatementMethodsHit,
		Name:  "patch_statement_methods_hit",
		Label: "Patch Statement Methods Hit",
		Short: "Patch Stmt Methods Hit",
		Doc:   "Changed methods with a statement run. Needs a diff."},
	{
		Key:           MaxCyclomaticComplexity,
		Name:          "max_complexity",
		Label:         "Max Cyclomatic Complexity",
		Title:         "Max Complexity",
		Short:         "Max Complexity",
		Value:         true,
		LowerIsBetter: true,
		Doc:           "Highest method complexity. Lower is better."},
}

// MethodMetricDefs are the metrics of one method, shown on the file page.
var MethodMetricDefs = []MetricDef{
	{
		Key:   MethodStatementCoverage,
		Name:  "statement_coverage",
		Label: "Statements",
		Title: "Statement Coverage",
		Doc:   "Statements run by tests."},
	{
		Key:   MethodLineCoverage,
		Name:  "line_coverage",
		Label: "Lines",
		Title: "Line Coverage",
		Doc:   "Lines run by tests."},
	{
		Key:   MethodPatchStatementCoverage,
		Name:  "patch_statement_coverage",
		Label: "Patch Statements",
		Title: "Patch Statement Coverage",
		Short: "Patch Stmts",
		Doc:   "Changed statements run by tests. Needs a diff."},
	{
		Key:   MethodPatchLineCoverage,
		Name:  "patch_line_coverage",
		Label: "Patch Lines",
		Title: "Patch Line Coverage",
		Doc:   "Changed lines run by tests. Needs a diff."},
	{
		Key:           CyclomaticComplexity,
		Name:          "complexity",
		Label:         "Cyclomatic Complexity",
		Short:         "Complexity",
		Value:         true,
		LowerIsBetter: true,
		Doc:           "Paths through the method. Lower is better."},
	{
		Key:           MethodCrapScore,
		Name:          "crap_score",
		Label:         "CRAP Score",
		Short:         "CRAP",
		Value:         true,
		LowerIsBetter: true,
		Doc:           "Complexity weighted by missing coverage. Lower is better."},
}

// DefaultFileMetrics and DefaultMethodMetrics are every metric, in table
// order: what a config without a metrics section shows.
var (
	DefaultFileMetrics   = keysOf(FileMetricDefs)
	DefaultMethodMetrics = keysOf(MethodMetricDefs)
)

func keysOf(defs []MetricDef) []MetricKey {
	keys := make([]MetricKey, len(defs))
	for i, d := range defs {
		keys[i] = d.Key
	}
	return keys
}

// Metric returns the definition of key, from either table.
func Metric(key MetricKey) (MetricDef, bool) {
	for _, defs := range [][]MetricDef{FileMetricDefs, MethodMetricDefs} {
		for _, d := range defs {
			if d.Key == key {
				return d, true
			}
		}
	}
	return MetricDef{}, false
}

// IsMethodMetric reports whether key is in the method table.
func IsMethodMetric(key MetricKey) bool {
	for _, d := range MethodMetricDefs {
		if d.Key == key {
			return true
		}
	}
	return false
}

func (d MetricDef) TitleLabel() string {
	if d.Title != "" {
		return d.Title
	}
	return d.Label
}

func (d MetricDef) ShortLabel() string {
	if d.Short != "" {
		return d.Short
	}
	return d.Label
}

func metricNames(defs []MetricDef) []string {
	names := make([]string, len(defs))
	for i, d := range defs {
		names[i] = d.Name
	}
	return names
}

func metricByName(defs []MetricDef,
	name string) (MetricDef, bool) {
	for _, d := range defs {
		if d.Name == name {
			return d, true
		}
	}
	return MetricDef{}, false
}
