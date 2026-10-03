// Package calculator turns the raw counts of a file, folder or method into
// the metrics the reports show. Each metric has one function here, under the
// key of its definition in internal/config/metrics.go.
package calculator

import (
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/utils"
)

// FileCalc calculates one metric of a file or folder. It returns false when
// the data for the metric is absent.
type FileCalc func(raw model.CoverageMetrics) (any, bool)

// MethodCalc calculates one metric of a method. prior holds the metrics in
// Deps, which are calculated first.
type MethodCalc struct {
	Deps      []config.MetricKey
	Calculate func(raw model.MethodMetrics, prior map[config.MetricKey]any) (any, bool)
}

// ratio is a coverage percentage; it has no value when nothing can be covered.
func ratio(covered, valid int) (any, bool) {
	if valid == 0 {
		return nil, false
	}
	return model.CoverageDetail{
		Percentage: utils.CalculatePercentage(covered, valid, 2),
		Covered:    covered,
		Uncovered:  valid - covered,
		Total:      valid,
	}, true
}

var FileRegistry = map[config.MetricKey]FileCalc{
	config.LineCoverage: func(m model.CoverageMetrics) (any, bool) {
		return ratio(m.LinesCovered, m.LinesValid)
	},
	config.StatementCoverage: func(m model.CoverageMetrics) (any, bool) {
		return ratio(m.StatementsCovered, m.StatementsValid)
	},
	config.MethodsHit: func(m model.CoverageMetrics) (any, bool) {
		return ratio(m.MethodsHit, m.MethodsValid)
	},
	config.MethodsFullyCovered: func(m model.CoverageMetrics) (any, bool) {
		return ratio(m.MethodsFullyCovered, m.MethodsValid)
	},
	config.StatementMethodsHit: func(m model.CoverageMetrics) (any, bool) {
		return ratio(m.StatementMethodsHit, m.StatementMethodsValid)
	},
	config.StatementMethodsFullyCovered: func(m model.CoverageMetrics) (any, bool) {
		return ratio(m.StatementMethodsFullyCovered, m.StatementMethodsValid)
	},
	config.PatchLineCoverage: func(m model.CoverageMetrics) (any, bool) {
		return ratio(m.PatchLinesCovered, m.PatchLinesValid)
	},
	config.PatchStatementCoverage: func(m model.CoverageMetrics) (any, bool) {
		return ratio(m.PatchStatementsCovered, m.PatchStatementsValid)
	},
	config.PatchMethodsHit: func(m model.CoverageMetrics) (any, bool) {
		return ratio(m.PatchMethodsHit, m.PatchMethodsValid)
	},
	config.PatchStatementMethodsHit: func(m model.CoverageMetrics) (any, bool) {
		return ratio(m.PatchStatementMethodsHit, m.PatchStatementMethodsValid)
	},
	config.MaxCyclomaticComplexity: func(m model.CoverageMetrics) (any, bool) {
		return model.ScoreDetail{Value: float64(m.MaxCyclomaticComplexity)}, true
	},
}

var MethodRegistry = map[config.MetricKey]MethodCalc{
	config.MethodLineCoverage: {Calculate: func(m model.MethodMetrics, _ map[config.MetricKey]any) (any, bool) {
		return ratio(m.LinesCovered, m.LinesValid)
	}},
	config.MethodStatementCoverage: {Calculate: func(m model.MethodMetrics, _ map[config.MetricKey]any) (any, bool) {
		return ratio(m.StatementsCovered, m.StatementsValid)
	}},
	// the patch metrics exist only for a method the change touches
	config.MethodPatchLineCoverage: {Calculate: func(m model.MethodMetrics, _ map[config.MetricKey]any) (any, bool) {
		if m.DiffStatus == "" {
			return nil, false
		}
		return ratio(m.PatchLinesCovered, m.PatchLinesValid)
	}},
	config.MethodPatchStatementCoverage: {Calculate: func(m model.MethodMetrics, _ map[config.MetricKey]any) (any, bool) {
		if m.DiffStatus == "" {
			return nil, false
		}
		return ratio(m.PatchStatementsCovered, m.PatchStatementsValid)
	}},
	config.CyclomaticComplexity: {Calculate: func(m model.MethodMetrics, _ map[config.MetricKey]any) (any, bool) {
		if m.CyclomaticComplexity == nil {
			return nil, false
		}
		return model.ScoreDetail{Value: float64(*m.CyclomaticComplexity)}, true
	}},
	// CRAP: comp^2 * (1 - line coverage)^3 + comp
	config.MethodCrapScore: {
		Deps: []config.MetricKey{config.CyclomaticComplexity, config.MethodLineCoverage},
		Calculate: func(_ model.MethodMetrics, prior map[config.MetricKey]any) (any, bool) {
			comp, uncovered, ok := complexityAndUncovered(prior, config.MethodLineCoverage)
			if !ok {
				return nil, false
			}
			return model.ScoreDetail{Value: comp*comp*uncovered*uncovered*uncovered + comp}, true
		},
	},
}

// complexityAndUncovered reads the complexity of a method and the uncovered
// share (0 to 1) of the given coverage metric from the metrics calculated before.
func complexityAndUncovered(prior map[config.MetricKey]any, coverage config.MetricKey) (comp, uncovered float64, ok bool) {
	score, hasComp := prior[config.CyclomaticComplexity].(model.ScoreDetail)
	cov, hasCov := prior[coverage].(model.CoverageDetail)
	if !hasComp || !hasCov {
		return 0, 0, false
	}
	return score.Value, 1 - cov.Percentage/100, true
}
