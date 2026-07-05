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

func (LineCoverageCalculator) Key() config.MetricKey         { return config.LineCoverage }
func (LineCoverageCalculator) DependsOn() []config.MetricKey { return nil }
func (LineCoverageCalculator) Calculate(raw model.CoverageMetrics, prior map[config.MetricKey]any) (any, bool) {
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

func (StatementCoverageCalculator) Key() config.MetricKey         { return config.StatementCoverage }
func (StatementCoverageCalculator) DependsOn() []config.MetricKey { return nil }
func (StatementCoverageCalculator) Calculate(raw model.CoverageMetrics, prior map[config.MetricKey]any) (any, bool) {
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

func (BranchCoverageCalculator) Key() config.MetricKey         { return config.BranchCoverage }
func (BranchCoverageCalculator) DependsOn() []config.MetricKey { return nil }
func (BranchCoverageCalculator) Calculate(raw model.CoverageMetrics, prior map[config.MetricKey]any) (any, bool) {
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

func (MethodsHitCalculator) Key() config.MetricKey         { return config.MethodsHit }
func (MethodsHitCalculator) DependsOn() []config.MetricKey { return nil }
func (MethodsHitCalculator) Calculate(raw model.CoverageMetrics, prior map[config.MetricKey]any) (any, bool) {
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

func (MethodsFullyCoveredCalculator) Key() config.MetricKey         { return config.MethodsFullyCovered }
func (MethodsFullyCoveredCalculator) DependsOn() []config.MetricKey { return nil }
func (MethodsFullyCoveredCalculator) Calculate(raw model.CoverageMetrics, prior map[config.MetricKey]any) (any, bool) {
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

func (PatchLineCoverageCalculator) Key() config.MetricKey         { return config.PatchLineCoverage }
func (PatchLineCoverageCalculator) DependsOn() []config.MetricKey { return nil }
func (PatchLineCoverageCalculator) Calculate(raw model.CoverageMetrics, prior map[config.MetricKey]any) (any, bool) {
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

func (PatchStatementCoverageCalculator) Key() config.MetricKey         { return config.PatchStatementCoverage }
func (PatchStatementCoverageCalculator) DependsOn() []config.MetricKey { return nil }
func (PatchStatementCoverageCalculator) Calculate(raw model.CoverageMetrics, prior map[config.MetricKey]any) (any, bool) {
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

func (PatchMethodsHitCalculator) Key() config.MetricKey         { return config.PatchMethodsHit }
func (PatchMethodsHitCalculator) DependsOn() []config.MetricKey { return nil }
func (PatchMethodsHitCalculator) Calculate(raw model.CoverageMetrics, prior map[config.MetricKey]any) (any, bool) {
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
func (MaxCyclomaticComplexityCalculator) DependsOn() []config.MetricKey { return nil }
func (MaxCyclomaticComplexityCalculator) Calculate(raw model.CoverageMetrics, prior map[config.MetricKey]any) (any, bool) {
	return model.ScoreDetail{
		Value: float64(raw.MaxCyclomaticComplexity),
	}, true
}

// =============================================================================
// Method-Level Calculators
// =============================================================================

type MethodLineCoverageCalculator struct{}

func (MethodLineCoverageCalculator) Key() config.MetricKey         { return config.MethodLineCoverage }
func (MethodLineCoverageCalculator) DependsOn() []config.MetricKey { return nil }
func (MethodLineCoverageCalculator) Calculate(raw model.MethodMetrics, prior map[config.MetricKey]any) (any, bool) {
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
func (MethodStatementCoverageCalculator) DependsOn() []config.MetricKey { return nil }
func (MethodStatementCoverageCalculator) Calculate(raw model.MethodMetrics, prior map[config.MetricKey]any) (any, bool) {
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

func (MethodBranchCoverageCalculator) Key() config.MetricKey         { return config.MethodBranchCoverage }
func (MethodBranchCoverageCalculator) DependsOn() []config.MetricKey { return nil }
func (MethodBranchCoverageCalculator) Calculate(raw model.MethodMetrics, prior map[config.MetricKey]any) (any, bool) {
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
func (MethodPatchLineCoverageCalculator) DependsOn() []config.MetricKey { return nil }
func (MethodPatchLineCoverageCalculator) Calculate(raw model.MethodMetrics, prior map[config.MetricKey]any) (any, bool) {
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
func (MethodPatchStatementCoverageCalculator) DependsOn() []config.MetricKey { return nil }
func (MethodPatchStatementCoverageCalculator) Calculate(raw model.MethodMetrics, prior map[config.MetricKey]any) (any, bool) {
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
func (MethodCyclomaticComplexityCalculator) DependsOn() []config.MetricKey { return nil }
func (MethodCyclomaticComplexityCalculator) Calculate(raw model.MethodMetrics, prior map[config.MetricKey]any) (any, bool) {
	if raw.CyclomaticComplexity == nil {
		return nil, false
	}
	return model.ScoreDetail{
		Value: float64(*raw.CyclomaticComplexity),
	}, true
}

type MethodCrapScoreCalculator struct{}

func (MethodCrapScoreCalculator) Key() config.MetricKey {
	return config.MethodCrapScore
}

func (MethodCrapScoreCalculator) DependsOn() []config.MetricKey {
	return []config.MetricKey{config.CyclomaticComplexity, config.MethodLineCoverage}
}

func (MethodCrapScoreCalculator) Calculate(raw model.MethodMetrics, prior map[config.MetricKey]any) (any, bool) {
	compRaw, hasComp := prior[config.CyclomaticComplexity]
	covRaw, hasCov := prior[config.MethodLineCoverage]

	if !hasComp || !hasCov {
		return nil, false
	}

	compScore, compOk := compRaw.(model.ScoreDetail)
	covDetail, covOk := covRaw.(model.CoverageDetail)

	if !compOk || !covOk {
		return nil, false
	}

	comp := compScore.Value
	cov := covDetail.Percentage

	// CRAP formula: comp(m)^2 * (1 - cov(m)/100)^3 + comp(m)
	compSquared := comp * comp
	uncoveredRatio := 1.0 - (cov / 100.0)
	uncoveredRatioCubed := uncoveredRatio * uncoveredRatio * uncoveredRatio
	crap := (compSquared * uncoveredRatioCubed) + comp

	return model.ScoreDetail{
		Value: crap,
	}, true
}

type MethodPatchCrapScoreCalculator struct{}

func (MethodPatchCrapScoreCalculator) Key() config.MetricKey {
	return config.MethodPatchCrapScore
}

func (MethodPatchCrapScoreCalculator) DependsOn() []config.MetricKey {
	return []config.MetricKey{config.CyclomaticComplexity, config.MethodPatchStatementCoverage}
}

func (MethodPatchCrapScoreCalculator) Calculate(raw model.MethodMetrics, prior map[config.MetricKey]any) (any, bool) {
	compRaw, hasComp := prior[config.CyclomaticComplexity]
	covRaw, hasCov := prior[config.MethodPatchStatementCoverage]

	if !hasComp || !hasCov {
		return nil, false
	}

	compScore, compOk := compRaw.(model.ScoreDetail)
	covDetail, covOk := covRaw.(model.CoverageDetail)

	if !compOk || !covOk {
		return nil, false
	}

	comp := compScore.Value
	cov := covDetail.Percentage

	// PCRAP formula: CC(m)^2 * (1 - PCov(m))^3 + CC(m)
	compSquared := comp * comp
	uncoveredRatio := 1.0 - (cov / 100.0)
	uncoveredRatioCubed := uncoveredRatio * uncoveredRatio * uncoveredRatio
	crap := (compSquared * uncoveredRatioCubed) + comp

	return model.ScoreDetail{
		Value: crap,
	}, true
}

type MethodExposedRiskCalculator struct{}

func (MethodExposedRiskCalculator) Key() config.MetricKey {
	return config.MethodExposedRisk
}

func (MethodExposedRiskCalculator) DependsOn() []config.MetricKey {
	return []config.MetricKey{config.CyclomaticComplexity, config.MethodStatementCoverage}
}

func (MethodExposedRiskCalculator) Calculate(raw model.MethodMetrics, prior map[config.MetricKey]any) (any, bool) {
	compRaw, hasComp := prior[config.CyclomaticComplexity]
	covRaw, hasCov := prior[config.MethodStatementCoverage]

	if !hasComp || !hasCov {
		return nil, false
	}

	compScore, compOk := compRaw.(model.ScoreDetail)
	covDetail, covOk := covRaw.(model.CoverageDetail)

	if !compOk || !covOk {
		return nil, false
	}

	comp := compScore.Value
	cov := covDetail.Percentage

	// ExposedRisk formula: CC(m) * (1 - Cov(m))
	uncoveredRatio := 1.0 - (cov / 100.0)
	exposedRisk := comp * uncoveredRatio

	return model.ScoreDetail{
		Value: exposedRisk,
	}, true
}

type MethodDefectProbabilityCalculator struct{}

func (MethodDefectProbabilityCalculator) Key() config.MetricKey {
	return config.MethodDefectProbability
}

func (MethodDefectProbabilityCalculator) DependsOn() []config.MetricKey {
	return []config.MetricKey{config.CyclomaticComplexity, config.MethodPatchStatementCoverage, config.MethodStatementCoverage}
}

func (MethodDefectProbabilityCalculator) Calculate(raw model.MethodMetrics, prior map[config.MetricKey]any) (any, bool) {
	compRaw, hasComp := prior[config.CyclomaticComplexity]
	pcovRaw, hasPcov := prior[config.MethodPatchStatementCoverage]
	covRaw, hasCov := prior[config.MethodStatementCoverage]

	if !hasComp || !hasPcov || !hasCov {
		return nil, false
	}

	compScore, compOk := compRaw.(model.ScoreDetail)
	pcovDetail, pcovOk := pcovRaw.(model.CoverageDetail)
	covDetail, covOk := covRaw.(model.CoverageDetail)

	if !compOk || !pcovOk || !covOk {
		return nil, false
	}

	comp := compScore.Value
	pcov := pcovDetail.Percentage
	cov := covDetail.Percentage

	// Trigger HIGH RISK if: CC(m) > 10 AND PCov(m) < 50% AND Cov(m) < 70%
	isHighRisk := comp > 10.0 && pcov < 50.0 && cov < 70.0

	val := 0.0
	if isHighRisk {
		val = 1.0
	}

	return model.ScoreDetail{
		Value: val,
	}, true
}
