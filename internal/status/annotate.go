package status

import (
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
)

// Annotate traverses the entire model.SummaryTree and attaches risk statuses
// to every file, directory, and method node based on the active metric
// configuration and the evaluator registry.
//
// It loops over cfg.ActiveFileMetrics for file/directory nodes and
// cfg.ActiveMethodMetrics for method nodes. For each active metric key it
// looks up a registered Evaluator from the registry. This design means zero
// hardcoded metric keys appear in the annotation logic.
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
		ev, ok := registry[key]
		if !ok {
			continue
		}
		if !ev.IsApplicable(caps) {
			continue
		}
		if lvl, show := ev.Evaluate(metrics, bandPtr(bands, key)); show {
			(*statuses)[key] = string(lvl)
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
		// evaluators can operate uniformly on both scopes.
		cm := methodToCoverageMetrics(m)

		for key := range activeMetrics {
			ev, ok := registry[key]
			if !ok {
				continue
			}
			if !ev.IsApplicable(caps) {
				continue
			}
			if lvl, show := ev.Evaluate(cm, bandPtr(bands, key)); show {
				m.Statuses[key] = string(lvl)
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
		Calculated:             m.Calculated,
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
