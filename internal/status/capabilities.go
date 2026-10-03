package status

import "github.com/IgorBayerl/nanovision/internal/model"

// DeriveCapabilities reports which kinds of coverage the aggregated tree holds.
func DeriveCapabilities(tree *model.SummaryTree) Capabilities {
	caps := Capabilities{}
	note := func(m model.CoverageMetrics) {
		caps.HasMethodCoverage = caps.HasMethodCoverage || m.MethodsValid > 0
		caps.HasStatementCoverage = caps.HasStatementCoverage || m.StatementsValid > 0
	}
	var walk func(n *model.DirNode)
	walk = func(n *model.DirNode) {
		note(n.Metrics)
		for _, c := range n.Subdirs {
			walk(c)
		}
		for _, f := range n.Files {
			note(f.Metrics)
		}
	}
	if tree != nil && tree.Root != nil {
		walk(tree.Root)
	}
	return caps
}
