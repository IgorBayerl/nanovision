// Package status centralizes the logic for classifying coverage metrics into risk
// levels (e.g., "danger", "warning", "safe"). Its primary goal is to decouple
// the risk classification rules from the report generators.
//
// The process works in a dedicated "ANNOTATE" stage in the main pipeline:
//  1. The `AppConfig` loads risk thresholds from `nanovision.yaml` into `status_bands`.
//  2. The `DeriveCapabilities` function checks which metrics (like statement coverage)
//     are actually present in the parsed report data.
//  3. The `Annotate` function is called once. It traverses the entire in-memory
//     data tree (`model.SummaryTree`).
//  4. For each node (file or directory), it looks up evaluators from the registry
//     and falls back to generic percentage specs for unregistered metrics.
//  5. The risk level is attached directly to the node in its `Statuses` map.
//
// By pre-computing these statuses, all report generators (HTML, text summary, etc.)
// can simply read the status from the model without needing to know about the
// specific thresholds or calculation logic. This makes the reporters simpler and
// ensures consistent status reporting everywhere.
package status

// RiskLevel represents the classification of a coverage metric based on predefined thresholds.
type RiskLevel string

const (
	RiskDanger  RiskLevel = "danger"  // Indicates a metric is below the configured minimum threshold.
	RiskWarning RiskLevel = "warning" // Indicates a metric is within the configured warning range.
	RiskSafe    RiskLevel = "safe"    // Indicates a metric is above the configured maximum threshold.
)

// Capabilities informs the annotation process about which metrics are semantically
// available in the current dataset. This is crucial for handling different report
// formats that may lack certain features.
//
// For example, a file without an analyzer has no statements. With
// `HasStatementCoverage` false, the `Annotate` function skips the status of
// `statement_coverage`, so a metric that was never measured does not show a
// misleading "0% danger".
//
// This struct is populated by `DeriveCapabilities` after all data has been aggregated.
type Capabilities struct {
	HasMethodCoverage    bool
	HasStatementCoverage bool
}
