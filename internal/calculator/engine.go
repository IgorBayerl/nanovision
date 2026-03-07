package calculator

import (
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
)

// CalculateTree traverses the coverage tree and populates the Calculated map on every metric.
func CalculateTree(tree *model.SummaryTree, activeFileMetrics map[config.MetricKey]bool, activeMethodMetrics map[config.MetricKey]bool) {
	// First do the root
	if tree.Metrics.Calculated == nil {
		tree.Metrics.Calculated = make(map[config.MetricKey]any)
	}
	for key := range activeFileMetrics {
		if calc, ok := FileRegistry[key]; ok {
			if res, ok := calc.Calculate(tree.Metrics); ok {
				tree.Metrics.Calculated[key] = res
			}
		}
	}

	var walk func(n *model.DirNode)
	walk = func(n *model.DirNode) {
		if n.Metrics.Calculated == nil {
			n.Metrics.Calculated = make(map[config.MetricKey]any)
		}
		for key := range activeFileMetrics {
			if calc, ok := FileRegistry[key]; ok {
				if res, ok := calc.Calculate(n.Metrics); ok {
					n.Metrics.Calculated[key] = res
				}
			}
		}

		for _, file := range n.Files {
			if file.Metrics.Calculated == nil {
				file.Metrics.Calculated = make(map[config.MetricKey]any)
			}
			for key := range activeFileMetrics {
				if calc, ok := FileRegistry[key]; ok {
					if res, ok := calc.Calculate(file.Metrics); ok {
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
				for key := range activeMethodMetrics {
					if calc, ok := MethodRegistry[key]; ok {
						if res, ok := calc.Calculate(*m); ok {
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
