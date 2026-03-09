package calculator

import (
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
)

// topSortFileCalculators returns a slice of MetricKeys sorted such that
// for every calculator C, its dependencies appear before C in the slice.
func topSortFileCalculators(active map[config.MetricKey]bool) []config.MetricKey {
	var sorted []config.MetricKey
	visited := make(map[config.MetricKey]bool)
	visiting := make(map[config.MetricKey]bool)

	var visit func(key config.MetricKey)
	visit = func(key config.MetricKey) {
		if visited[key] {
			return
		}
		if visiting[key] {
			// Circular dependency detected, ignore for now
			return
		}
		visiting[key] = true

		if calc, ok := FileRegistry[key]; ok {
			for _, dep := range calc.DependsOn() {
				// Auto-enable dependencies if they aren't explicitly requested
				active[dep] = true
				visit(dep)
			}
		}

		visiting[key] = false
		visited[key] = true
		sorted = append(sorted, key)
	}

	// We copy the initial keys to avoid map iteration issues while modifying active
	var initialKeys []config.MetricKey
	for k := range active {
		initialKeys = append(initialKeys, k)
	}

	for _, key := range initialKeys {
		visit(key)
	}
	return sorted
}

// topSortMethodCalculators returns a slice of MetricKeys sorted such that
// for every calculator C, its dependencies appear before C in the slice.
func topSortMethodCalculators(active map[config.MetricKey]bool) []config.MetricKey {
	var sorted []config.MetricKey
	visited := make(map[config.MetricKey]bool)
	visiting := make(map[config.MetricKey]bool)

	var visit func(key config.MetricKey)
	visit = func(key config.MetricKey) {
		if visited[key] {
			return
		}
		if visiting[key] {
			// Circular dependency detected, ignore for now
			return
		}
		visiting[key] = true

		if calc, ok := MethodRegistry[key]; ok {
			for _, dep := range calc.DependsOn() {
				// Auto-enable dependencies if they aren't explicitly requested
				active[dep] = true
				visit(dep)
			}
		}

		visiting[key] = false
		visited[key] = true
		sorted = append(sorted, key)
	}

	// We copy the initial keys to avoid map iteration issues while modifying active
	var initialKeys []config.MetricKey
	for k := range active {
		initialKeys = append(initialKeys, k)
	}

	for _, key := range initialKeys {
		visit(key)
	}
	return sorted
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
