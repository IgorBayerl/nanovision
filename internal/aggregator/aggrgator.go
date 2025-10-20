// internal/aggregator/aggregator.go

package aggregator

import "github.com/IgorBayerl/AdlerCov/internal/model"

// AggregateMetricsAfterEnrichment recalculates and aggregates all metrics for the entire
// tree now that method data from the enrichment phase is available. This should be
// called *after* the enricher has run and *before* reports are generated.
func AggregateMetricsAfterEnrichment(tree *model.SummaryTree) {
	// The recursive call on the root node will calculate and aggregate all metrics
	// from the bottom up. The return value is the final, correct total for the project.
	tree.Metrics = aggregateNodeMetrics(tree.Root)
}

// aggregateNodeMetrics performs a post-order traversal to correctly sum all metrics.
// It returns the total metrics for the subtree rooted at the given `dir` node.
func aggregateNodeMetrics(dir *model.DirNode) model.CoverageMetrics {
	// Start with a clean slate for this directory's aggregated metrics.
	currentDirTotals := model.CoverageMetrics{}

	// 1. Recurse into subdirectories first.
	for _, subDir := range dir.Subdirs {
		// The result of the recursive call is the total for that entire subtree.
		subDirTotals := aggregateNodeMetrics(subDir)
		// Store this total on the subdirectory node itself.
		subDir.Metrics = subDirTotals
		// Add the subdirectory's total to the running total for the current directory.
		addMetrics(&currentDirTotals, subDirTotals)
	}

	// 2. Process the files in the current directory.
	for _, file := range dir.Files {
		// First, update the file's own metrics with method coverage stats.
		// The line/branch metrics on file.Metrics are already correct from the BUILD step.
		calculateFileMethodMetrics(file)
		// Then, add the file's complete metrics to the current directory's total.
		addMetrics(&currentDirTotals, file.Metrics)
	}

	// 3. Return the final aggregated metrics for this directory and its children.
	return currentDirTotals
}

// calculateFileMethodMetrics updates a single file's metrics struct with method coverage
// statistics based on the enriched data.
func calculateFileMethodMetrics(file *model.FileNode) {
	// Reset only the method counters before recalculating to ensure freshness.
	file.Metrics.MethodsValid = 0
	file.Metrics.MethodsCovered = 0
	file.Metrics.MethodsFullyCovered = 0

	for _, method := range file.Methods {
		// A method is only valid if it has at least one coverable line.
		if method.LinesValid > 0 {
			file.Metrics.MethodsValid++
			if method.LinesCovered > 0 {
				file.Metrics.MethodsCovered++
			}
			if method.LinesCovered == method.LinesValid {
				file.Metrics.MethodsFullyCovered++
			}
		}
	}
}

// addMetrics is a helper function to safely sum two CoverageMetrics structs.
func addMetrics(dest *model.CoverageMetrics, src model.CoverageMetrics) {
	dest.LinesCovered += src.LinesCovered
	dest.LinesValid += src.LinesValid
	dest.BranchesCovered += src.BranchesCovered
	dest.BranchesValid += src.BranchesValid
	dest.TotalLines += src.TotalLines
	dest.MethodsValid += src.MethodsValid
	dest.MethodsCovered += src.MethodsCovered
	dest.MethodsFullyCovered += src.MethodsFullyCovered
}
