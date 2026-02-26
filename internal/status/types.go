// Package status centralizes the logic for classifying coverage metrics into risk
// levels (e.g., "danger", "warning", "safe"). Its primary goal is to decouple
// the risk classification rules from the report generators.
//
// The process works in a dedicated "ANNOTATE" stage in the main pipeline:
//  1. The `AppConfig` loads risk thresholds from `nanovision.yaml` into `status_bands`.
//  2. The `deriveCapabilities` function checks which metrics (like branch coverage)
//     are actually present in the parsed report data.
//  3. The `Annotate` function is called once. It traverses the entire in-memory
//     data tree (`model.SummaryTree`).
//  4. For each node (file or directory), it calculates the metric percentages and uses
//     the `Classify` function to determine a risk level.
//  5. This risk level is attached directly to the node in its `Statuses` map.
//
// By pre-computing these statuses, all report generators (HTML, text summary, etc.)
// can simply read the status from the model without needing to know about the
// specific thresholds or calculation logic. This makes the reporters simpler and
// ensures consistent status reporting everywhere.
package status

import "github.com/IgorBayerl/nanovision/internal/config"

// RiskLevel represents the classification of a coverage metric based on predefined thresholds.
type RiskLevel string

const (
	RiskDanger  RiskLevel = "danger"  // Indicates a metric is below the configured minimum threshold.
	RiskWarning RiskLevel = "warning" // Indicates a metric is within the configured warning range.
	RiskSafe    RiskLevel = "safe"    // Indicates a metric is above the configured maximum threshold.
)

// MetricKey is the canonical snake_case key used consistently across the configuration
// (nanovision.yaml), the data model, and the status logic to identify a metric.
type MetricKey = config.MetricKey

const (
	LineCoverage                 MetricKey = config.LineCoverage
	BranchCoverage               MetricKey = config.BranchCoverage
	MethodsCovered               MetricKey = config.MethodsCovered
	MethodsFullyCovered          MetricKey = config.MethodsFullyCovered
	PatchLineCoverage            MetricKey = config.PatchLineCoverage
	PatchMethodsCovered          MetricKey = config.PatchMethodsCovered
	StatementCoverage            MetricKey = config.StatementCoverage
	PatchStatementCoverage       MetricKey = config.PatchStatementCoverage
	StatementMethodsCovered      MetricKey = config.StatementMethodsCovered
	StatementMethodsFullyCovered MetricKey = config.StatementMethodsFullyCovered
	PatchStatementMethodsCovered MetricKey = config.PatchStatementMethodsCovered
)

// Capabilities informs the annotation process about which metrics are semantically
// available in the current dataset. This is crucial for handling different report
// formats that may lack certain features.
//
// For example, Go's native 'gocover' format does not produce branch coverage data.
// By setting `HasBranchCoverage` to false, the `Annotate` function will know to skip
// adding a status for `branch_coverage`, preventing a misleading "0% danger" status
// from appearing in reports for a metric that was never measured.
//
// This struct is populated by the `deriveCapabilities` function in `main.go` after
// all data has been aggregated.
type Capabilities struct {
	HasBranchCoverage    bool
	HasMethodCoverage    bool
	HasStatementCoverage bool
}
