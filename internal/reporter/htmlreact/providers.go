package htmlreact

import (
	"fmt"

	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/utils"
)

// =============================================================================
// Method-level providers (MethodMetricProvider)
// =============================================================================

// MethodMetricProvider defines the interface for mapping model.MethodMetrics
// into the UI-specific methodDetail format using the Strategy Pattern.
type MethodMetricProvider interface {
	Key() config.MetricKey
	Apply(method *model.MethodMetrics, ui *methodDetail)
}

// StatementCoverageMethodProvider Strategy
type StatementCoverageMethodProvider struct{}

func (p StatementCoverageMethodProvider) Key() config.MetricKey {
	return config.MethodStatementCoverage
}
func (p StatementCoverageMethodProvider) Apply(m *model.MethodMetrics, ui *methodDetail) {
	if m.StatementsValid > 0 {
		ui.Metrics[MethodUIStmtCoverage] = methodMetric{
			Value: fmt.Sprintf("%d / %d", m.StatementsCovered, m.StatementsValid),
		}
	}
}

// LineCoverageMethodProvider Strategy
type LineCoverageMethodProvider struct{}

func (p LineCoverageMethodProvider) Key() config.MetricKey { return config.MethodLineCoverage }
func (p LineCoverageMethodProvider) Apply(m *model.MethodMetrics, ui *methodDetail) {
	ui.Metrics[MethodUILineCoverage] = methodMetric{
		Value: fmt.Sprintf("%d / %d", m.LinesCovered, m.LinesValid),
	}
}

// BranchCoverageMethodProvider Strategy
type BranchCoverageMethodProvider struct{}

func (p BranchCoverageMethodProvider) Key() config.MetricKey { return config.MethodBranchCoverage }
func (p BranchCoverageMethodProvider) Apply(m *model.MethodMetrics, ui *methodDetail) {
	if m.BranchesValid > 0 {
		ui.Metrics[MethodUIBranchCoverage] = methodMetric{
			Value: fmt.Sprintf("%d / %d", m.BranchesCovered, m.BranchesValid),
		}
	}
}

// PatchLineCoverageMethodProvider Strategy
type PatchLineCoverageMethodProvider struct{}

func (p PatchLineCoverageMethodProvider) Key() config.MetricKey {
	return config.MethodPatchLineCoverage
}
func (p PatchLineCoverageMethodProvider) Apply(m *model.MethodMetrics, ui *methodDetail) {
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

// PatchStatementMethodProvider Strategy
type PatchStatementMethodProvider struct{}

func (p PatchStatementMethodProvider) Key() config.MetricKey {
	return config.MethodPatchStatementCoverage
}
func (p PatchStatementMethodProvider) Apply(m *model.MethodMetrics, ui *methodDetail) {
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

// CyclomaticComplexityMethodProvider Strategy
type CyclomaticComplexityMethodProvider struct{}

func (p CyclomaticComplexityMethodProvider) Key() config.MetricKey {
	return config.CyclomaticComplexity
}
func (p CyclomaticComplexityMethodProvider) Apply(m *model.MethodMetrics, ui *methodDetail) {
	if m.CyclomaticComplexity != nil {
		ui.Metrics[MethodUICyclomaticComplexity] = methodMetric{
			Value: fmt.Sprintf("%d", *m.CyclomaticComplexity),
		}
	}
}

// =============================================================================
// File-level providers (FileMetricProvider)
// =============================================================================

// FileMetricProvider defines the interface for mapping model.CoverageMetrics
// into the UI-specific metricsMap format using the Strategy Pattern.
type FileMetricProvider interface {
	Key() config.MetricKey
	Apply(metrics model.CoverageMetrics, ui metricsMap)
}

// --- StatementCoverageFileProvider ---

type StatementCoverageFileProvider struct{}

func (p StatementCoverageFileProvider) Key() config.MetricKey { return config.StatementCoverage }
func (p StatementCoverageFileProvider) Apply(m model.CoverageMetrics, ui metricsMap) {
	if m.StatementsValid > 0 {
		statementPct := utils.CalculatePercentage(m.StatementsCovered, m.StatementsValid, 2)
		ui[string(config.StatementCoverage)] = lineCoverageDetail{
			Covered:    m.StatementsCovered,
			Uncovered:  m.StatementsValid - m.StatementsCovered,
			Coverable:  m.StatementsValid,
			Total:      m.StatementsValid,
			Percentage: statementPct,
		}
	}
}

// --- LineCoverageFileProvider ---

type LineCoverageFileProvider struct{}

func (p LineCoverageFileProvider) Key() config.MetricKey { return config.LineCoverage }
func (p LineCoverageFileProvider) Apply(m model.CoverageMetrics, ui metricsMap) {
	linePct := utils.CalculatePercentage(m.LinesCovered, m.LinesValid, 2)
	ui[string(config.LineCoverage)] = lineCoverageDetail{
		Covered:    m.LinesCovered,
		Uncovered:  m.LinesValid - m.LinesCovered,
		Coverable:  m.LinesValid,
		Total:      m.TotalLines,
		Percentage: linePct,
	}
}

// --- BranchCoverageFileProvider ---

type BranchCoverageFileProvider struct{}

func (p BranchCoverageFileProvider) Key() config.MetricKey { return config.BranchCoverage }
func (p BranchCoverageFileProvider) Apply(m model.CoverageMetrics, ui metricsMap) {
	if m.BranchesValid > 0 {
		branchPct := utils.CalculatePercentage(m.BranchesCovered, m.BranchesValid, 2)
		ui[string(config.BranchCoverage)] = branchCoverageDetail{
			Covered:    m.BranchesCovered,
			Total:      m.BranchesValid,
			Percentage: branchPct,
		}
	}
}

// --- MethodsHitFileProvider ---

type MethodsHitFileProvider struct{}

func (p MethodsHitFileProvider) Key() config.MetricKey { return config.MethodsHit }
func (p MethodsHitFileProvider) Apply(m model.CoverageMetrics, ui metricsMap) {
	if m.MethodsValid > 0 {
		methodsHitPct := utils.CalculatePercentage(m.MethodsHit, m.MethodsValid, 2)
		ui[string(config.MethodsHit)] = methodsHitDetail{
			Covered:    m.MethodsHit,
			Total:      m.MethodsValid,
			Percentage: methodsHitPct,
		}
	}
}

// --- MethodsFullyCoveredFileProvider ---

type MethodsFullyCoveredFileProvider struct{}

func (p MethodsFullyCoveredFileProvider) Key() config.MetricKey { return config.MethodsFullyCovered }
func (p MethodsFullyCoveredFileProvider) Apply(m model.CoverageMetrics, ui metricsMap) {
	if m.MethodsValid > 0 {
		methodsFullyCoveredPct := utils.CalculatePercentage(m.MethodsFullyCovered, m.MethodsValid, 2)
		ui[string(config.MethodsFullyCovered)] = methodsFullyCoveredDetail{
			Covered:    m.MethodsFullyCovered,
			Total:      m.MethodsValid,
			Percentage: methodsFullyCoveredPct,
		}
	}
}

// --- PatchStatementCoverageFileProvider ---

type PatchStatementCoverageFileProvider struct{}

func (p PatchStatementCoverageFileProvider) Key() config.MetricKey {
	return config.PatchStatementCoverage
}
func (p PatchStatementCoverageFileProvider) Apply(m model.CoverageMetrics, ui metricsMap) {
	if m.PatchStatementsValid > 0 {
		patchStatementPct := utils.CalculatePercentage(m.PatchStatementsCovered, m.PatchStatementsValid, 2)
		ui[string(config.PatchStatementCoverage)] = lineCoverageDetail{
			Covered:    m.PatchStatementsCovered,
			Uncovered:  m.PatchStatementsValid - m.PatchStatementsCovered,
			Coverable:  m.PatchStatementsValid,
			Total:      m.PatchStatementsValid,
			Percentage: patchStatementPct,
		}
	}
}

// --- PatchLineCoverageFileProvider ---

type PatchLineCoverageFileProvider struct{}

func (p PatchLineCoverageFileProvider) Key() config.MetricKey { return config.PatchLineCoverage }
func (p PatchLineCoverageFileProvider) Apply(m model.CoverageMetrics, ui metricsMap) {
	if m.PatchLinesTotal > 0 {
		patchLinePct := utils.CalculatePercentage(m.PatchLinesCovered, m.PatchLinesValid, 2)
		ui[string(config.PatchLineCoverage)] = lineCoverageDetail{
			Covered:    m.PatchLinesCovered,
			Uncovered:  m.PatchLinesValid - m.PatchLinesCovered,
			Coverable:  m.PatchLinesValid,
			Total:      m.PatchLinesTotal,
			Percentage: patchLinePct,
		}
	}
}

// --- PatchMethodsHitFileProvider ---

type PatchMethodsHitFileProvider struct{}

func (p PatchMethodsHitFileProvider) Key() config.MetricKey { return config.PatchMethodsHit }
func (p PatchMethodsHitFileProvider) Apply(m model.CoverageMetrics, ui metricsMap) {
	if m.PatchMethodsValid > 0 {
		patchMethodsPct := utils.CalculatePercentage(m.PatchMethodsHit, m.PatchMethodsValid, 2)
		ui[string(config.PatchMethodsHit)] = methodsHitDetail{
			Covered:    m.PatchMethodsHit,
			Total:      m.PatchMethodsValid,
			Percentage: patchMethodsPct,
		}
	}
}

// =============================================================================
// File Provider Registry
// =============================================================================

// FileProviderRegistry maps each file-scoped MetricKey to its FileMetricProvider.
var FileProviderRegistry = map[config.MetricKey]FileMetricProvider{
	config.StatementCoverage:      StatementCoverageFileProvider{},
	config.LineCoverage:           LineCoverageFileProvider{},
	config.BranchCoverage:         BranchCoverageFileProvider{},
	config.MethodsHit:             MethodsHitFileProvider{},
	config.MethodsFullyCovered:    MethodsFullyCoveredFileProvider{},
	config.PatchStatementCoverage: PatchStatementCoverageFileProvider{},
	config.PatchLineCoverage:      PatchLineCoverageFileProvider{},
	config.PatchMethodsHit:        PatchMethodsHitFileProvider{},
}

// =============================================================================
// Method Provider Registry
// =============================================================================

// MethodProviderRegistry maps each method-scoped MetricKey to its MethodMetricProvider.
var MethodProviderRegistry = map[config.MetricKey]MethodMetricProvider{
	config.MethodStatementCoverage:      StatementCoverageMethodProvider{},
	config.MethodLineCoverage:           LineCoverageMethodProvider{},
	config.MethodBranchCoverage:         BranchCoverageMethodProvider{},
	config.MethodPatchLineCoverage:      PatchLineCoverageMethodProvider{},
	config.MethodPatchStatementCoverage: PatchStatementMethodProvider{},
	config.CyclomaticComplexity:         CyclomaticComplexityMethodProvider{},
}
