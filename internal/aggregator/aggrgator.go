// internal/aggregator/aggregator.go

package aggregator

import "github.com/IgorBayerl/nanovision/internal/model"

// AggregateMetricsAfterEnrichment recalculates and aggregates all metrics for the entire
// tree now that method data from the enrichment phase is available. This should be
// called *after* the enricher has run and *before* reports are generated.
func AggregateMetricsAfterEnrichment(tree *model.SummaryTree) {
	// The recursive call on the root node will calculate and aggregate all metrics
	// from the bottom up. The return value is the final, correct total for the project.
	totalMetrics := aggregateNodeMetrics(tree.Root)
	tree.Metrics = totalMetrics
	tree.Root.Metrics = totalMetrics // This ensures the root node itself has the final aggregated metrics.
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
		// First, update the file's own metrics with method coverage stats
		// and all patch/diff-based metrics.
		calculateFileMethodMetrics(file)

		// Then, add the file's now-complete metrics to the current directory's total.
		addMetrics(&currentDirTotals, file.Metrics)
	}

	// 3. Return the final aggregated metrics for this directory and its children.
	return currentDirTotals
}

// calculateFileMethodMetrics updates a single file's metrics struct with method coverage
// statistics based on the enriched data, and computes patch/diff-based metrics as well.
func calculateFileMethodMetrics(file *model.FileNode) {
	// Reset method-level and patch/diff-based counters before recalculating to ensure freshness.
	file.Metrics.MethodsValid = 0
	file.Metrics.MethodsHit = 0
	file.Metrics.MethodsFullyCovered = 0

	file.Metrics.PatchLinesCovered = 0
	file.Metrics.PatchLinesValid = 0
	file.Metrics.PatchLinesTotal = 0
	file.Metrics.PatchMethodsHit = 0
	file.Metrics.PatchMethodsValid = 0

	file.Metrics.PatchStatementsCovered = 0
	file.Metrics.PatchStatementsValid = 0

	file.Metrics.StatementMethodsValid = 0
	file.Metrics.StatementMethodsHit = 0
	file.Metrics.StatementMethodsFullyCovered = 0

	file.Metrics.PatchStatementMethodsHit = 0
	file.Metrics.PatchStatementMethodsValid = 0

	// --- Patch line metrics (per file) ---
	// We only consider coverable lines (Hits >= 0). A "patch line" is any
	// coverable line that was added or modified according to DiffInfo.
	if file.Diff != nil {
		// Set the total number of changed lines regardless of type.
		file.Metrics.PatchLinesTotal = len(file.Diff.AddedLines) + len(file.Diff.ModifiedLines)

		if file.Diff.Kind == model.ChangeKindAdded {
			// For added files, all coverable lines are considered part of the patch.
			file.Metrics.PatchLinesValid = file.Metrics.LinesValid
			file.Metrics.PatchLinesCovered = file.Metrics.LinesCovered
		} else {
			for lineNumber, lineMetric := range file.Lines {
				// Not coverable => not part of the patch coverage denominator.
				if lineMetric.Hits < 0 {
					continue
				}
				if file.Diff.AddedLines[lineNumber] || file.Diff.ModifiedLines[lineNumber] {
					file.Metrics.PatchLinesValid++
					if lineMetric.Hits > 0 {
						file.Metrics.PatchLinesCovered++
					}
				}
			}
		}

		// Patch Statements
		if file.Diff.Kind == model.ChangeKindAdded {
			file.Metrics.PatchStatementsValid = file.Metrics.StatementsValid
			file.Metrics.PatchStatementsCovered = file.Metrics.StatementsCovered
		} else {
			for _, stmt := range file.Statements {
				inPatch := false
				for i := stmt.StartLine; i <= stmt.EndLine; i++ {
					if file.Diff.AddedLines[i] || file.Diff.ModifiedLines[i] {
						inPatch = true
						break
					}
				}
				if inPatch {
					file.Metrics.PatchStatementsValid++
					covered := false
					for i := stmt.StartLine; i <= stmt.EndLine; i++ {
						if line, ok := file.Lines[i]; ok && line.Hits > 0 {
							if file.Diff.AddedLines[i] || file.Diff.ModifiedLines[i] {
								covered = true
								break
							}
						}
					}
					if covered {
						file.Metrics.PatchStatementsCovered++
					}
				}
			}
		}
	}

	hasDiff := file.Diff != nil

	for _, method := range file.Methods {
		// A method is only valid if it has at least one coverable line.
		if method.LinesValid > 0 {
			file.Metrics.MethodsValid++
			if method.LinesCovered > 0 {
				file.Metrics.MethodsHit++
			}
			if method.LinesCovered == method.LinesValid {
				file.Metrics.MethodsFullyCovered++
			}
		}

		if method.StatementsValid > 0 {
			file.Metrics.StatementMethodsValid++
			if method.StatementsCovered > 0 {
				file.Metrics.StatementMethodsHit++
			}
			if method.StatementsCovered == method.StatementsValid {
				file.Metrics.StatementMethodsFullyCovered++
			}
		}

		// If there is no diff associated with this file, there can be no
		// patch-level method metrics.
		if !hasDiff {
			continue
		}

		// For added files, every method with coverable lines is considered a
		// "patch method". It is "covered" if any of its lines are covered.
		if file.Diff.Kind == model.ChangeKindAdded {
			if method.LinesValid > 0 {
				file.Metrics.PatchMethodsValid++
				if method.LinesCovered > 0 {
					file.Metrics.PatchMethodsHit++
				}
			}
			if method.StatementsValid > 0 {
				file.Metrics.PatchStatementMethodsValid++
				if method.StatementsCovered > 0 {
					file.Metrics.PatchStatementMethodsHit++
				}
			}
			continue
		}

		// For modified files, determine if this method contains any changed,
		// coverable lines, and whether any of those changed lines were executed.
		patchLinesTotal := 0
		patchLinesCovered := 0

		for line := method.StartLine; line <= method.EndLine; line++ {
			if !file.Diff.AddedLines[line] && !file.Diff.ModifiedLines[line] {
				continue
			}

			if lm, ok := file.Lines[line]; ok && lm.Hits >= 0 {
				patchLinesTotal++
				if lm.Hits > 0 {
					patchLinesCovered++
				}
			}
		}

		if patchLinesTotal > 0 {
			// This method is affected by the patch.
			file.Metrics.PatchMethodsValid++
			// "Covered" for patch methods means: at least one changed, coverable
			// line in the method was executed.
			if patchLinesCovered > 0 {
				file.Metrics.PatchMethodsHit++
			}

			if method.StatementsValid > 0 {
				file.Metrics.PatchStatementMethodsValid++
				if method.StatementsCovered > 0 {
					file.Metrics.PatchStatementMethodsHit++
				}
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
	dest.MethodsHit += src.MethodsHit
	dest.MethodsFullyCovered += src.MethodsFullyCovered

	dest.PatchLinesCovered += src.PatchLinesCovered
	dest.PatchLinesValid += src.PatchLinesValid
	dest.PatchLinesTotal += src.PatchLinesTotal
	dest.PatchMethodsHit += src.PatchMethodsHit
	dest.PatchMethodsValid += src.PatchMethodsValid

	dest.StatementsValid += src.StatementsValid
	dest.StatementsCovered += src.StatementsCovered
	dest.PatchStatementsValid += src.PatchStatementsValid
	dest.PatchStatementsCovered += src.PatchStatementsCovered

	dest.StatementMethodsValid += src.StatementMethodsValid
	dest.StatementMethodsHit += src.StatementMethodsHit
	dest.StatementMethodsFullyCovered += src.StatementMethodsFullyCovered
	dest.PatchStatementMethodsValid += src.PatchStatementMethodsValid
	dest.PatchStatementMethodsHit += src.PatchStatementMethodsHit
}
