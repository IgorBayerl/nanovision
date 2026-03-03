package htmlreact

import (
	"fmt"

	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
)

// MethodMetricProvider defines the interface for mapping model.MethodMetrics
// into the UI-specific methodDetail format using the Strategy Pattern.
type MethodMetricProvider interface {
	Key() config.MetricKey
	Apply(method *model.MethodMetrics, ui *methodDetail)
}

// StatementCoverageProvider Strategy
type StatementCoverageProvider struct{}

func (p StatementCoverageProvider) Key() config.MetricKey { return config.StatementCoverage }
func (p StatementCoverageProvider) Apply(m *model.MethodMetrics, ui *methodDetail) {
	if m.StatementsValid > 0 {
		ui.Metrics[MethodUIStmtCoverage] = methodMetric{
			Value: fmt.Sprintf("%d / %d", m.StatementsCovered, m.StatementsValid),
		}
	}
}

// LineCoverageProvider Strategy
type LineCoverageProvider struct{}

func (p LineCoverageProvider) Key() config.MetricKey { return config.LineCoverage }
func (p LineCoverageProvider) Apply(m *model.MethodMetrics, ui *methodDetail) {
	ui.Metrics[MethodUILineCoverage] = methodMetric{
		Value: fmt.Sprintf("%d / %d", m.LinesCovered, m.LinesValid),
	}
}

// BranchCoverageProvider Strategy
type BranchCoverageProvider struct{}

func (p BranchCoverageProvider) Key() config.MetricKey { return config.BranchCoverage }
func (p BranchCoverageProvider) Apply(m *model.MethodMetrics, ui *methodDetail) {
	if m.BranchesValid > 0 {
		ui.Metrics[MethodUIBranchCoverage] = methodMetric{
			Value: fmt.Sprintf("%d / %d", m.BranchesCovered, m.BranchesValid),
		}
	}
}

// PatchLineCoverageProvider Strategy
type PatchLineCoverageProvider struct{}

func (p PatchLineCoverageProvider) Key() config.MetricKey { return config.PatchLineCoverage }
func (p PatchLineCoverageProvider) Apply(m *model.MethodMetrics, ui *methodDetail) {
	// Only apply if there is patch data
	if m.DiffStatus != "" {
		if m.PatchLinesValid > 0 {
			ui.NewLinesCoverage = &newLinesCoverage{
				Covered: m.PatchLinesCovered,
				Total:   m.PatchLinesValid,
			}
		}
		ui.Metrics[MethodUIPatchLineCoverage] = methodMetric{
			Value: fmt.Sprintf("%d / %d", m.PatchLinesCovered, m.PatchLinesValid),
		}
	}
}

// PatchStatementProvider Strategy
type PatchStatementProvider struct{}

func (p PatchStatementProvider) Key() config.MetricKey { return config.PatchStatementCoverage }
func (p PatchStatementProvider) Apply(m *model.MethodMetrics, ui *methodDetail) {
	// Only apply if there is patch data
	if m.DiffStatus != "" && m.StatementsValid > 0 {
		if m.PatchStatementsValid > 0 || m.PatchLinesValid > 0 {
			cov := &newLinesCoverage{
				Covered: m.PatchStatementsCovered,
				Total:   m.PatchStatementsValid,
			}
			ui.NewStatementsCoverage = cov
			ui.NewStatementCoverage = cov
		}
		ui.Metrics[MethodUIPatchStmtCoverage] = methodMetric{
			Value: fmt.Sprintf("%d / %d", m.PatchStatementsCovered, m.PatchStatementsValid),
		}
	}
}

// CyclomaticComplexityProvider Strategy
type CyclomaticComplexityProvider struct{}

func (p CyclomaticComplexityProvider) Key() config.MetricKey { return config.MaxCyclomaticComplexity }
func (p CyclomaticComplexityProvider) Apply(m *model.MethodMetrics, ui *methodDetail) {
	if m.CyclomaticComplexity != nil {
		ui.Metrics[MethodUICyclomaticComplexity] = methodMetric{
			Value: fmt.Sprintf("%d", *m.CyclomaticComplexity),
		}
	}
}
