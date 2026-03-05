package textsummary

import (
	"fmt"
	"io"

	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/utils"
)

// TextMetricProvider defines the strategy interface for the text summary reporter.
// Each implementation knows how to render one metric in both the top-level
// summary section and the per-node tree row.
type TextMetricProvider interface {
	// Key returns the config.MetricKey this provider handles.
	Key() config.MetricKey
	// PrintSummary writes the top-level summary lines (e.g., "Statement coverage: 80%").
	PrintSummary(w io.Writer, tree *model.SummaryTree)
	// NodePart returns the small string for the tree row (e.g., "80% (Stmt)").
	// Returns "" when the metric should be hidden for this node.
	NodePart(m model.CoverageMetrics) string
}

// =============================================================================
// Concrete Providers
// =============================================================================

// --- StatementCoverageTextProvider ---

type StatementCoverageTextProvider struct{}

func (StatementCoverageTextProvider) Key() config.MetricKey { return config.StatementCoverage }

func (StatementCoverageTextProvider) PrintSummary(w io.Writer, tree *model.SummaryTree) {
	if tree.Metrics.StatementsValid > 0 {
		statementCoverage := utils.CalculatePercentage(tree.Metrics.StatementsCovered, tree.Metrics.StatementsValid, 1)
		fmt.Fprintf(w, "  Statement coverage: %s\n", utils.FormatPercentage(statementCoverage, 0))
		fmt.Fprintf(w, "  Covered statements: %d\n", tree.Metrics.StatementsCovered)
		fmt.Fprintf(w, "  Uncovered statements: %d\n", tree.Metrics.StatementsValid-tree.Metrics.StatementsCovered)
		fmt.Fprintf(w, "  Total statements: %d\n", tree.Metrics.StatementsValid)
	}
}

func (StatementCoverageTextProvider) NodePart(m model.CoverageMetrics) string {
	if m.StatementsValid > 0 {
		stmtCov := utils.CalculatePercentage(m.StatementsCovered, m.StatementsValid, 1)
		return fmt.Sprintf("%s (Stmt)", utils.FormatPercentage(stmtCov, 0))
	}
	return ""
}

// --- LineCoverageTextProvider ---

type LineCoverageTextProvider struct{}

func (LineCoverageTextProvider) Key() config.MetricKey { return config.LineCoverage }

func (LineCoverageTextProvider) PrintSummary(w io.Writer, tree *model.SummaryTree) {
	lineCoverage := utils.CalculatePercentage(tree.Metrics.LinesCovered, tree.Metrics.LinesValid, 1)
	fmt.Fprintf(w, "  Line coverage: %s\n", utils.FormatPercentage(lineCoverage, 0))
	fmt.Fprintf(w, "  Covered lines: %d\n", tree.Metrics.LinesCovered)
	fmt.Fprintf(w, "  Uncovered lines: %d\n", tree.Metrics.LinesValid-tree.Metrics.LinesCovered)
	fmt.Fprintf(w, "  Coverable lines: %d\n", tree.Metrics.LinesValid)
}

func (LineCoverageTextProvider) NodePart(m model.CoverageMetrics) string {
	lineCov := utils.CalculatePercentage(m.LinesCovered, m.LinesValid, 1)
	return fmt.Sprintf("%s (Line)", utils.FormatPercentage(lineCov, 0))
}

// --- BranchCoverageTextProvider ---

type BranchCoverageTextProvider struct{}

func (BranchCoverageTextProvider) Key() config.MetricKey { return config.BranchCoverage }

func (BranchCoverageTextProvider) PrintSummary(w io.Writer, tree *model.SummaryTree) {
	if tree.Metrics.BranchesValid > 0 {
		branchCoverage := utils.CalculatePercentage(tree.Metrics.BranchesCovered, tree.Metrics.BranchesValid, 1)
		fmt.Fprintf(w, "  Branch coverage: %s (%d of %d)\n", utils.FormatPercentage(branchCoverage, 0), tree.Metrics.BranchesCovered, tree.Metrics.BranchesValid)
	}
}

func (BranchCoverageTextProvider) NodePart(m model.CoverageMetrics) string {
	if m.BranchesValid > 0 {
		branchCov := utils.CalculatePercentage(m.BranchesCovered, m.BranchesValid, 1)
		return fmt.Sprintf("%s (Branch)", utils.FormatPercentage(branchCov, 0))
	}
	return ""
}

// =============================================================================
// Text Provider Registry
// =============================================================================

// TextProviderRegistry maps each MetricKey to its TextMetricProvider.
var TextProviderRegistry = map[config.MetricKey]TextMetricProvider{
	config.StatementCoverage: StatementCoverageTextProvider{},
	config.LineCoverage:      LineCoverageTextProvider{},
	config.BranchCoverage:    BranchCoverageTextProvider{},
}
