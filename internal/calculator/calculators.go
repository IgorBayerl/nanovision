package calculator

import (
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/utils"
)

// =============================================================================
// File-Level Calculators
// =============================================================================

type LineCoverageCalculator struct{}

func (LineCoverageCalculator) Key() config.MetricKey { return config.LineCoverage }
func (LineCoverageCalculator) Calculate(raw model.CoverageMetrics) (any, bool) {
	if raw.LinesValid == 0 {
		return nil, false
	}
	return model.CoverageDetail{
		Percentage: utils.CalculatePercentage(raw.LinesCovered, raw.LinesValid, 2),
		Covered:    raw.LinesCovered,
		Uncovered:  raw.LinesValid - raw.LinesCovered,
		Total:      raw.LinesValid,
	}, true
}

type StatementCoverageCalculator struct{}

func (StatementCoverageCalculator) Key() config.MetricKey { return config.StatementCoverage }
func (StatementCoverageCalculator) Calculate(raw model.CoverageMetrics) (any, bool) {
	if raw.StatementsValid == 0 {
		return nil, false
	}
	return model.CoverageDetail{
		Percentage: utils.CalculatePercentage(raw.StatementsCovered, raw.StatementsValid, 2),
		Covered:    raw.StatementsCovered,
		Uncovered:  raw.StatementsValid - raw.StatementsCovered,
		Total:      raw.StatementsValid,
	}, true
}

type BranchCoverageCalculator struct{}

func (BranchCoverageCalculator) Key() config.MetricKey { return config.BranchCoverage }
func (BranchCoverageCalculator) Calculate(raw model.CoverageMetrics) (any, bool) {
	if raw.BranchesValid == 0 {
		return nil, false
	}
	return model.CoverageDetail{
		Percentage: utils.CalculatePercentage(raw.BranchesCovered, raw.BranchesValid, 2),
		Covered:    raw.BranchesCovered,
		Uncovered:  raw.BranchesValid - raw.BranchesCovered,
		Total:      raw.BranchesValid,
	}, true
}

type MethodsHitCalculator struct{}

func (MethodsHitCalculator) Key() config.MetricKey { return config.MethodsHit }
func (MethodsHitCalculator) Calculate(raw model.CoverageMetrics) (any, bool) {
	if raw.MethodsValid == 0 {
		return nil, false
	}
	return model.CoverageDetail{
		Percentage: utils.CalculatePercentage(raw.MethodsHit, raw.MethodsValid, 2),
		Covered:    raw.MethodsHit,
		Uncovered:  raw.MethodsValid - raw.MethodsHit,
		Total:      raw.MethodsValid,
	}, true
}

type MethodsFullyCoveredCalculator struct{}

func (MethodsFullyCoveredCalculator) Key() config.MetricKey { return config.MethodsFullyCovered }
func (MethodsFullyCoveredCalculator) Calculate(raw model.CoverageMetrics) (any, bool) {
	if raw.MethodsValid == 0 {
		return nil, false
	}
	return model.CoverageDetail{
		Percentage: utils.CalculatePercentage(raw.MethodsFullyCovered, raw.MethodsValid, 2),
		Covered:    raw.MethodsFullyCovered,
		Uncovered:  raw.MethodsValid - raw.MethodsFullyCovered,
		Total:      raw.MethodsValid,
	}, true
}

type PatchLineCoverageCalculator struct{}

func (PatchLineCoverageCalculator) Key() config.MetricKey { return config.PatchLineCoverage }
func (PatchLineCoverageCalculator) Calculate(raw model.CoverageMetrics) (any, bool) {
	if raw.PatchLinesValid == 0 {
		return nil, false
	}
	return model.CoverageDetail{
		Percentage: utils.CalculatePercentage(raw.PatchLinesCovered, raw.PatchLinesValid, 2),
		Covered:    raw.PatchLinesCovered,
		Uncovered:  raw.PatchLinesValid - raw.PatchLinesCovered,
		Total:      raw.PatchLinesValid,
	}, true
}

type PatchStatementCoverageCalculator struct{}

func (PatchStatementCoverageCalculator) Key() config.MetricKey { return config.PatchStatementCoverage }
func (PatchStatementCoverageCalculator) Calculate(raw model.CoverageMetrics) (any, bool) {
	if raw.PatchStatementsValid == 0 {
		return nil, false
	}
	return model.CoverageDetail{
		Percentage: utils.CalculatePercentage(raw.PatchStatementsCovered, raw.PatchStatementsValid, 2),
		Covered:    raw.PatchStatementsCovered,
		Uncovered:  raw.PatchStatementsValid - raw.PatchStatementsCovered,
		Total:      raw.PatchStatementsValid,
	}, true
}

type PatchMethodsHitCalculator struct{}

func (PatchMethodsHitCalculator) Key() config.MetricKey { return config.PatchMethodsHit }
func (PatchMethodsHitCalculator) Calculate(raw model.CoverageMetrics) (any, bool) {
	if raw.PatchMethodsValid == 0 {
		return nil, false
	}
	return model.CoverageDetail{
		Percentage: utils.CalculatePercentage(raw.PatchMethodsHit, raw.PatchMethodsValid, 2),
		Covered:    raw.PatchMethodsHit,
		Uncovered:  raw.PatchMethodsValid - raw.PatchMethodsHit,
		Total:      raw.PatchMethodsValid,
	}, true
}

type MaxCyclomaticComplexityCalculator struct{}

func (MaxCyclomaticComplexityCalculator) Key() config.MetricKey {
	return config.MaxCyclomaticComplexity
}
func (MaxCyclomaticComplexityCalculator) Calculate(raw model.CoverageMetrics) (any, bool) {
	return model.ScoreDetail{
		Value: float64(raw.MaxCyclomaticComplexity),
	}, true
}

// =============================================================================
// Method-Level Calculators
// =============================================================================

type MethodLineCoverageCalculator struct{}

func (MethodLineCoverageCalculator) Key() config.MetricKey { return config.MethodLineCoverage }
func (MethodLineCoverageCalculator) Calculate(raw model.MethodMetrics) (any, bool) {
	if raw.LinesValid == 0 {
		return nil, false
	}
	return model.CoverageDetail{
		Percentage: utils.CalculatePercentage(raw.LinesCovered, raw.LinesValid, 2),
		Covered:    raw.LinesCovered,
		Uncovered:  raw.LinesValid - raw.LinesCovered,
		Total:      raw.LinesValid,
	}, true
}

type MethodStatementCoverageCalculator struct{}

func (MethodStatementCoverageCalculator) Key() config.MetricKey {
	return config.MethodStatementCoverage
}
func (MethodStatementCoverageCalculator) Calculate(raw model.MethodMetrics) (any, bool) {
	if raw.StatementsValid == 0 {
		return nil, false
	}
	return model.CoverageDetail{
		Percentage: utils.CalculatePercentage(raw.StatementsCovered, raw.StatementsValid, 2),
		Covered:    raw.StatementsCovered,
		Uncovered:  raw.StatementsValid - raw.StatementsCovered,
		Total:      raw.StatementsValid,
	}, true
}

type MethodBranchCoverageCalculator struct{}

func (MethodBranchCoverageCalculator) Key() config.MetricKey { return config.MethodBranchCoverage }
func (MethodBranchCoverageCalculator) Calculate(raw model.MethodMetrics) (any, bool) {
	if raw.BranchesValid == 0 {
		return nil, false
	}
	return model.CoverageDetail{
		Percentage: utils.CalculatePercentage(raw.BranchesCovered, raw.BranchesValid, 2),
		Covered:    raw.BranchesCovered,
		Uncovered:  raw.BranchesValid - raw.BranchesCovered,
		Total:      raw.BranchesValid,
	}, true
}

type MethodPatchLineCoverageCalculator struct{}

func (MethodPatchLineCoverageCalculator) Key() config.MetricKey {
	return config.MethodPatchLineCoverage
}
func (MethodPatchLineCoverageCalculator) Calculate(raw model.MethodMetrics) (any, bool) {
	if raw.DiffStatus == "" || raw.PatchLinesValid == 0 {
		return nil, false
	}
	return model.CoverageDetail{
		Percentage: utils.CalculatePercentage(raw.PatchLinesCovered, raw.PatchLinesValid, 2),
		Covered:    raw.PatchLinesCovered,
		Uncovered:  raw.PatchLinesValid - raw.PatchLinesCovered,
		Total:      raw.PatchLinesValid,
	}, true
}

type MethodPatchStatementCoverageCalculator struct{}

func (MethodPatchStatementCoverageCalculator) Key() config.MetricKey {
	return config.MethodPatchStatementCoverage
}
func (MethodPatchStatementCoverageCalculator) Calculate(raw model.MethodMetrics) (any, bool) {
	if raw.DiffStatus == "" || raw.PatchStatementsValid == 0 {
		return nil, false
	}
	return model.CoverageDetail{
		Percentage: utils.CalculatePercentage(raw.PatchStatementsCovered, raw.PatchStatementsValid, 2),
		Covered:    raw.PatchStatementsCovered,
		Uncovered:  raw.PatchStatementsValid - raw.PatchStatementsCovered,
		Total:      raw.PatchStatementsValid,
	}, true
}

type MethodCyclomaticComplexityCalculator struct{}

func (MethodCyclomaticComplexityCalculator) Key() config.MetricKey {
	return config.CyclomaticComplexity
}
func (MethodCyclomaticComplexityCalculator) Calculate(raw model.MethodMetrics) (any, bool) {
	if raw.CyclomaticComplexity == nil {
		return nil, false
	}
	return model.ScoreDetail{
		Value: float64(*raw.CyclomaticComplexity),
	}, true
}
