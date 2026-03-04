package status

import (
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/utils"
)

// percentageSpec describes a generic "higher is better" percentage metric
// that is not in the evaluator registry. It maps a MetricKey to field
// extractors so annotate can evaluate them without hardcoding keys inline.
type percentageSpec struct {
	numerator   func(model.CoverageMetrics) int
	denominator func(model.CoverageMetrics) int
	applicable  func(Capabilities) bool
}

// fallbackSpecs contains every percentage-based metric that is NOT handled
// by a registered Evaluator. Each entry maps a MetricKey to its field
// extractors and applicability check. All use ClassifyHigherIsBetter.
var fallbackSpecs = map[config.MetricKey]percentageSpec{
	config.MethodsHit: {
		numerator:   func(m model.CoverageMetrics) int { return m.MethodsHit },
		denominator: func(m model.CoverageMetrics) int { return m.MethodsValid },
		applicable:  func(c Capabilities) bool { return c.HasMethodCoverage },
	},
	config.MethodsFullyCovered: {
		numerator:   func(m model.CoverageMetrics) int { return m.MethodsFullyCovered },
		denominator: func(m model.CoverageMetrics) int { return m.MethodsValid },
		applicable:  func(c Capabilities) bool { return c.HasMethodCoverage },
	},
	config.PatchLineCoverage: {
		numerator:   func(m model.CoverageMetrics) int { return m.PatchLinesCovered },
		denominator: func(m model.CoverageMetrics) int { return m.PatchLinesValid },
		applicable:  func(_ Capabilities) bool { return true },
	},
	config.PatchMethodsHit: {
		numerator:   func(m model.CoverageMetrics) int { return m.PatchMethodsHit },
		denominator: func(m model.CoverageMetrics) int { return m.PatchMethodsValid },
		applicable:  func(_ Capabilities) bool { return true },
	},
	config.PatchStatementCoverage: {
		numerator:   func(m model.CoverageMetrics) int { return m.PatchStatementsCovered },
		denominator: func(m model.CoverageMetrics) int { return m.PatchStatementsValid },
		applicable:  func(_ Capabilities) bool { return true },
	},
	config.StatementMethodsHit: {
		numerator:   func(m model.CoverageMetrics) int { return m.StatementMethodsHit },
		denominator: func(m model.CoverageMetrics) int { return m.StatementMethodsValid },
		applicable:  func(c Capabilities) bool { return c.HasStatementCoverage },
	},
	config.StatementMethodsFullyCovered: {
		numerator:   func(m model.CoverageMetrics) int { return m.StatementMethodsFullyCovered },
		denominator: func(m model.CoverageMetrics) int { return m.StatementMethodsValid },
		applicable:  func(c Capabilities) bool { return c.HasStatementCoverage },
	},
	config.PatchStatementMethodsHit: {
		numerator:   func(m model.CoverageMetrics) int { return m.PatchStatementMethodsHit },
		denominator: func(m model.CoverageMetrics) int { return m.PatchStatementMethodsValid },
		applicable:  func(_ Capabilities) bool { return true },
	},
}

// Annotate traverses the entire model.SummaryTree and attaches risk statuses
// to every file, directory, and method node based on the active metric
// configuration and the evaluator registry.
//
// It loops over cfg.ActiveFileMetrics for file/directory nodes and
// cfg.ActiveMethodMetrics for method nodes. For each active metric key it
// looks up a registered Evaluator first, then falls back to the generic
// percentage-based fallbackSpecs table. This design means zero hardcoded
// metric keys appear in the annotation logic.
//
// The `registry` parameter is the evaluator registry (typically
// evaluators.Registry) passed in by the caller to avoid an import cycle.
func Annotate(tree *model.SummaryTree, cfg *config.AppConfig, caps Capabilities, registry map[config.MetricKey]Evaluator) {
	if tree == nil || tree.Root == nil {
		return
	}
	walk(tree.Root, func(dir *model.DirNode) {
		annotateMetrics(dir.Metrics, &dir.Statuses, cfg.ActiveFileMetrics, cfg.StatusBands, caps, registry)
		for _, file := range dir.Files {
			annotateMetrics(file.Metrics, &file.Statuses, cfg.ActiveFileMetrics, cfg.StatusBands, caps, registry)
			annotateMethodNodes(file.Methods, cfg.ActiveMethodMetrics, cfg.StatusBands, caps, registry)
		}
	})
}

// annotateMetrics evaluates all active metrics for a single node (file or directory).
func annotateMetrics(
	metrics model.CoverageMetrics,
	statuses *map[config.MetricKey]string,
	activeMetrics map[config.MetricKey]bool,
	bands config.StatusBands,
	caps Capabilities,
	registry map[config.MetricKey]Evaluator,
) {
	*statuses = make(map[config.MetricKey]string)

	for key := range activeMetrics {
		// Try the evaluator registry first.
		if ev, ok := registry[key]; ok {
			if !ev.IsApplicable(caps) {
				continue
			}
			if lvl, show := ev.Evaluate(metrics, bandPtr(bands, key)); show {
				(*statuses)[key] = string(lvl)
			}
			continue
		}

		// Fall back to the generic percentage spec table.
		if spec, ok := fallbackSpecs[key]; ok {
			if !spec.applicable(caps) {
				continue
			}
			denom := spec.denominator(metrics)
			if denom == 0 {
				continue
			}
			pct := utils.CalculatePercentage(spec.numerator(metrics), denom, 2)
			if lvl, show := ClassifyHigherIsBetter(pct, bandPtr(bands, key)); show {
				(*statuses)[key] = string(lvl)
			}
		}
	}
}

// annotateMethodNodes evaluates active method-level metrics for every method
// in a file node.
func annotateMethodNodes(
	methods []model.MethodMetrics,
	activeMetrics map[config.MetricKey]bool,
	bands config.StatusBands,
	caps Capabilities,
	registry map[config.MetricKey]Evaluator,
) {
	for i := range methods {
		m := &methods[i]
		m.Statuses = make(map[config.MetricKey]string)

		// Build a CoverageMetrics from the method-level fields so that
		// evaluators can be reused without knowing about MethodMetrics.
		cm := methodToCoverageMetrics(m)

		for key := range activeMetrics {
			if ev, ok := registry[key]; ok {
				if !ev.IsApplicable(caps) {
					continue
				}
				if lvl, show := ev.Evaluate(cm, bandPtr(bands, key)); show {
					m.Statuses[key] = string(lvl)
				}
				continue
			}

			if spec, ok := fallbackSpecs[key]; ok {
				if !spec.applicable(caps) {
					continue
				}
				denom := spec.denominator(cm)
				if denom == 0 {
					continue
				}
				pct := utils.CalculatePercentage(spec.numerator(cm), denom, 2)
				if lvl, show := ClassifyHigherIsBetter(pct, bandPtr(bands, key)); show {
					m.Statuses[key] = string(lvl)
				}
			}
		}
	}
}

// methodToCoverageMetrics maps MethodMetrics fields into CoverageMetrics so
// that registered evaluators can operate uniformly on both scopes.
func methodToCoverageMetrics(m *model.MethodMetrics) model.CoverageMetrics {
	cm := model.CoverageMetrics{
		LinesValid:             m.LinesValid,
		LinesCovered:           m.LinesCovered,
		BranchesValid:          m.BranchesValid,
		BranchesCovered:        m.BranchesCovered,
		StatementsValid:        m.StatementsValid,
		StatementsCovered:      m.StatementsCovered,
		PatchLinesValid:        m.PatchLinesValid,
		PatchLinesCovered:      m.PatchLinesCovered,
		PatchStatementsValid:   m.PatchStatementsValid,
		PatchStatementsCovered: m.PatchStatementsCovered,
	}
	if m.CyclomaticComplexity != nil {
		cm.MaxCyclomaticComplexity = *m.CyclomaticComplexity
	}
	return cm
}

// bandPtr is a small helper to get a pointer to a band from the map, which is
// needed for the classifier's nil check.
func bandPtr(bands config.StatusBands, k config.MetricKey) *config.Band {
	if b, ok := bands[k]; ok {
		return &b
	}
	return nil
}

// walk performs a recursive traversal of the directory tree.
func walk(n *model.DirNode, fn func(*model.DirNode)) {
	for _, c := range n.Subdirs {
		walk(c, fn)
	}
	fn(n)
}
