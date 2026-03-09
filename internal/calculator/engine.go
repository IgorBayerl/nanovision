package calculator

import (
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
)

// Calculator is a constraint for types that have a DependsOn method.
type Calculator interface {
	DependsOn() []config.MetricKey
}

// topSortCalculators returns a slice of MetricKeys sorted such that
// for every calculator C, its dependencies appear before C in the slice.
func topSortCalculators[C Calculator](registry map[config.MetricKey]C, active map[config.MetricKey]bool) []config.MetricKey {
	var sorted []config.MetricKey
	visited := make(map[config.MetricKey]bool)
	visiting := make(map[config.MetricKey]bool)

	var visit func(key config.MetricKey)
	visit = func(key config.MetricKey) {
		if visited[key] || visiting[key] {
			return
		}
		visiting[key] = true

		if calc, ok := registry[key]; ok {
			for _, dep := range calc.DependsOn() {
				visit(dep)
			}
		}

		visiting[key] = false
		visited[key] = true
		sorted = append(sorted, key)
	}

	// Copy initial keys to avoid map iteration issues while modifying active.
	initialKeys := make([]config.MetricKey, 0, len(active))
	for k := range active {
		initialKeys = append(initialKeys, k)
	}

	for _, key := range initialKeys {
		visit(key)
	}
	return sorted
}

func topSortFileCalculators(active map[config.MetricKey]bool) []config.MetricKey {
	return topSortCalculators(FileRegistry, active)
}

func topSortMethodCalculators(active map[config.MetricKey]bool) []config.MetricKey {
	return topSortCalculators(MethodRegistry, active)
}

// CalculateTree traverses the coverage tree and populates the Calculated map on every metric.
func CalculateTree(tree *model.SummaryTree, activeFileMetrics map[config.MetricKey]bool, activeMethodMetrics map[config.MetricKey]bool) {
	sortedFileKeys := topSortFileCalculators(activeFileMetrics)
	sortedMethodKeys := topSortMethodCalculators(activeMethodMetrics)

	// First do the root
	if tree.Metrics.Calculated == nil {
		tree.Metrics.Calculated = make(map[config.MetricKey]any)
	}
	for _, key := range sortedFileKeys {
		if calc, ok := FileRegistry[key]; ok {
			if res, ok := calc.Calculate(tree.Metrics, tree.Metrics.Calculated); ok {
				tree.Metrics.Calculated[key] = res
			}
		}
	}

	var walk func(n *model.DirNode)
	walk = func(n *model.DirNode) {
		if n.Metrics.Calculated == nil {
			n.Metrics.Calculated = make(map[config.MetricKey]any)
		}
		for _, key := range sortedFileKeys {
			if calc, ok := FileRegistry[key]; ok {
				if res, ok := calc.Calculate(n.Metrics, n.Metrics.Calculated); ok {
					n.Metrics.Calculated[key] = res
				}
			}
		}

		for _, file := range n.Files {
			if file.Metrics.Calculated == nil {
				file.Metrics.Calculated = make(map[config.MetricKey]any)
			}
			for _, key := range sortedFileKeys {
				if calc, ok := FileRegistry[key]; ok {
					if res, ok := calc.Calculate(file.Metrics, file.Metrics.Calculated); ok {
						file.Metrics.Calculated[key] = res
					}
				}
			}

			// Methods
			for i := range file.Methods {
				m := &file.Methods[i] // need pointer to modify the slice element
				if m.Calculated == nil {
					m.Calculated = make(map[config.MetricKey]any)
				}
				for _, key := range sortedMethodKeys {
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
