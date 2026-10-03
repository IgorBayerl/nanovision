package calculator

import (
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
)

// methodOrder returns the active method metrics and their dependencies, each
// after the metrics it depends on.
func methodOrder(active map[config.MetricKey]bool) []config.MetricKey {
	var sorted []config.MetricKey
	visited := make(map[config.MetricKey]bool)

	var visit func(key config.MetricKey)
	visit = func(key config.MetricKey) {
		if visited[key] {
			return
		}
		visited[key] = true
		for _, dep := range MethodRegistry[key].Deps {
			visit(dep)
		}
		sorted = append(sorted, key)
	}
	for key := range active {
		visit(key)
	}
	return sorted
}

// CalculateTree traverses the coverage tree and populates the Calculated map on every metric.
func CalculateTree(tree *model.SummaryTree, activeFileMetrics map[config.MetricKey]bool, activeMethodMetrics map[config.MetricKey]bool) {
	methodKeys := methodOrder(activeMethodMetrics)

	calculateFile := func(m *model.CoverageMetrics) {
		if m.Calculated == nil {
			m.Calculated = make(map[config.MetricKey]any)
		}
		for key := range activeFileMetrics {
			if calc, ok := FileRegistry[key]; ok {
				if res, ok := calc(*m); ok {
					m.Calculated[key] = res
				}
			}
		}
	}

	calculateFile(&tree.Metrics)

	var walk func(n *model.DirNode)
	walk = func(n *model.DirNode) {
		calculateFile(&n.Metrics)
		for _, file := range n.Files {
			calculateFile(&file.Metrics)
			for i := range file.Methods {
				m := &file.Methods[i]
				if m.Calculated == nil {
					m.Calculated = make(map[config.MetricKey]any)
				}
				for _, key := range methodKeys {
					if calc, ok := MethodRegistry[key]; ok {
						if res, ok := calc.Calculate(*m, m.Calculated); ok {
							m.Calculated[key] = res
						}
					}
				}
			}
		}
		for _, sub := range n.Subdirs {
			walk(sub)
		}
	}
	if tree.Root != nil {
		walk(tree.Root)
	}
}
