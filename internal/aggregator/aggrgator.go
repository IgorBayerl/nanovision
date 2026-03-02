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

// -----------------------------------------------------------------------------
// Core Pipeline
// -----------------------------------------------------------------------------

// calculateFileMethodMetrics orchestrates the calculation of standard and patch metrics
// for both the file and its enclosed methods. Cyclomatic Complexity is now 1.
func calculateFileMethodMetrics(file *model.FileNode) {
	resetFileMetrics(file)

	if file.Diff != nil {
		calculateFilePatchMetrics(file)
	}

	calculateAllMethodMetrics(file)
}

// -----------------------------------------------------------------------------
// Pipeline Stage 1: File-Level Patch Metrics
// -----------------------------------------------------------------------------

func calculateFilePatchMetrics(file *model.FileNode) {
	file.Metrics.PatchLinesTotal = len(file.Diff.AddedLines) + len(file.Diff.ModifiedLines)

	// Optimization: If the whole file is new, copy standard metrics directly.
	if file.Diff.Kind == model.ChangeKindAdded {
		file.Metrics.PatchLinesValid = file.Metrics.LinesValid
		file.Metrics.PatchLinesCovered = file.Metrics.LinesCovered
		file.Metrics.PatchStatementsValid = file.Metrics.StatementsValid
		file.Metrics.PatchStatementsCovered = file.Metrics.StatementsCovered
		return
	}

	// 1. Calculate Patch Lines
	for ln, lm := range file.Lines {
		if lm.Hits >= 0 && isLineInPatch(ln, file.Diff) {
			file.Metrics.PatchLinesValid++
			if lm.Hits > 0 {
				file.Metrics.PatchLinesCovered++
			}
		}
	}

	// 2. Calculate Patch Statements
	for _, stmt := range file.Statements {
		inPatch, isCovered := evaluateStatementPatchStatus(stmt, file)
		if inPatch {
			file.Metrics.PatchStatementsValid++
			if isCovered {
				file.Metrics.PatchStatementsCovered++
			}
		}
	}
}

// -----------------------------------------------------------------------------
// Pipeline Stage 2: Method-Level Metrics (Standard & Patch)
// -----------------------------------------------------------------------------

func calculateAllMethodMetrics(file *model.FileNode) {
	for i := range file.Methods {
		method := &file.Methods[i]

		resetMethodMetrics(method)
		aggregateStandardMethodMetrics(file, method)

		if file.Diff != nil {
			calculateMethodPatchMetrics(file, method)
			aggregatePatchMethodMetrics(file, method)
		}
	}
}

func aggregateStandardMethodMetrics(file *model.FileNode, method *model.MethodMetrics) {
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
}

func calculateMethodPatchMetrics(file *model.FileNode, method *model.MethodMetrics) {
	if file.Diff.Kind == model.ChangeKindAdded {
		method.DiffStatus = "added"
		method.PatchLinesValid = method.LinesValid
		method.PatchLinesCovered = method.LinesCovered
		method.PatchStatementsValid = method.StatementsValid
		method.PatchStatementsCovered = method.StatementsCovered
		return
	}

	isModified := false

	// 1. Calculate Patch Lines for Method
	for ln := method.StartLine; ln <= method.EndLine; ln++ {
		if isLineInPatch(ln, file.Diff) {
			isModified = true
			if lm, ok := file.Lines[ln]; ok && lm.Hits >= 0 {
				method.PatchLinesValid++
				if lm.Hits > 0 {
					method.PatchLinesCovered++
				}
			}
		}
	}

	if isModified {
		if method.PatchLinesValid > 0 && method.PatchLinesValid == method.LinesValid {
			method.DiffStatus = "added"
		} else {
			method.DiffStatus = "modified"
		}
	}

	// 2. Calculate Patch Statements for Method
	for _, stmt := range file.Statements {
		if stmt.StartLine >= method.StartLine && stmt.EndLine <= method.EndLine {
			inPatch, isCovered := evaluateStatementPatchStatus(stmt, file)
			if inPatch {
				method.PatchStatementsValid++
				if isCovered {
					method.PatchStatementsCovered++
				}
			}
		}
	}
}

func aggregatePatchMethodMetrics(file *model.FileNode, method *model.MethodMetrics) {
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
		return
	}

	// For modified files, a method is in the patch if it has patch lines
	if method.PatchLinesValid > 0 {
		file.Metrics.PatchMethodsValid++
		if method.PatchLinesCovered > 0 {
			file.Metrics.PatchMethodsHit++
		}

		if method.StatementsValid > 0 {
			file.Metrics.PatchStatementMethodsValid++
			// Original logic check: if ANY standard statements are covered while the method was in the patch
			if method.StatementsCovered > 0 {
				file.Metrics.PatchStatementMethodsHit++
			}
		}
	}
}

// -----------------------------------------------------------------------------
// Core Helpers & Resetters
// -----------------------------------------------------------------------------

func isLineInPatch(line int, diff *model.DiffInfo) bool {
	return diff != nil && (diff.AddedLines[line] || diff.ModifiedLines[line])
}

// evaluateStatementPatchStatus checks if a statement intersects with a diff
// and if the statement is covered. It ensures that statements are only considered
// "in patch" if the patch modifies a semantically meaningful (coverable) line,
// preventing trailing empty lines from falsely pulling adjacent statements into the patch.
func evaluateStatementPatchStatus(stmt model.Statement, file *model.FileNode) (inPatch bool, isCovered bool) {
	statementHasCoverableLines := false
	patchIntersectsCoverableLine := false
	patchIntersectsAnyLine := false

	for i := stmt.StartLine; i <= stmt.EndLine; i++ {
		lm, ok := file.Lines[i]
		isCoverable := ok && lm.Hits >= 0

		if isCoverable {
			statementHasCoverableLines = true
			if lm.Hits > 0 {
				isCovered = true
			}
		}

		if isLineInPatch(i, file.Diff) {
			patchIntersectsAnyLine = true
			if isCoverable {
				patchIntersectsCoverableLine = true
			}
		}
	}

	// A statement is only part of the patch if a coverable line was modified.
	// If the statement has no coverable lines at all, we fallback to any line match.
	if statementHasCoverableLines {
		inPatch = patchIntersectsCoverableLine
	} else {
		inPatch = patchIntersectsAnyLine
	}

	if !inPatch {
		isCovered = false
	}

	return inPatch, isCovered
}

func resetFileMetrics(file *model.FileNode) {
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
}

func resetMethodMetrics(method *model.MethodMetrics) {
	method.PatchLinesValid = 0
	method.PatchLinesCovered = 0
	method.PatchStatementsValid = 0
	method.PatchStatementsCovered = 0
	method.DiffStatus = ""
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
