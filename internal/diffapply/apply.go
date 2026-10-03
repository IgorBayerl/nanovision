package diffapply

import (
	"log/slog"
	"slices"

	"github.com/IgorBayerl/nanovision/internal/diff"
	"github.com/IgorBayerl/nanovision/internal/filtering"
	"github.com/IgorBayerl/nanovision/internal/model"
)

// Apply annotates a model.SummaryTree with diff information.
// It resolves paths from the diff data, finds corresponding file nodes in the tree,
// and attaches change details like kind (added/modified) and changed line numbers.
//
// It also sorts every diff file into a group of model.ChangeSet and stores it
// on the tree. A file the tree does not have is "ignored" when filter removes
// it, and "not in the reports" otherwise. filter may be nil.
func Apply(tree *model.SummaryTree, dd *diff.DiffData, filter filtering.IFilter, logger *slog.Logger) *model.ChangeSet {
	if tree == nil || tree.Root == nil || dd == nil || len(dd.Files) == 0 {
		return nil
	}

	// 1. Build file index and covered set from the tree.
	fileIndex := make(map[string]*model.FileNode)
	coveredSet := make(map[string]bool)
	collectFilesAndCoverage(tree.Root, fileIndex, coveredSet)

	// 2. Build the resolver with heuristics.
	resolver := BuildResolver(dd, fileIndex, coveredSet, logger)
	change := &model.ChangeSet{}
	excluded := func(rel string) bool { return filter != nil && !filter.IsElementIncludedInReport(rel) }

	// 3. Iterate through each file in the diff and apply its changes.
	for _, fileDiff := range dd.Files {
		if fileDiff.NewPath == "/dev/null" || fileDiff.NewPath == "" {
			old := resolver.Relative(fileDiff.OldPath)
			if excluded(old) {
				change.Ignored = append(change.Ignored, old)
			} else {
				change.Deleted = append(change.Deleted, old)
			}
			continue
		}

		// Resolve the diff path to a path in our coverage tree.
		treePath, ok := resolver.Resolve(fileDiff.NewPath)
		node, exists := fileIndex[treePath]
		if !ok || !exists {
			rel := resolver.Relative(fileDiff.NewPath)
			if excluded(rel) {
				logger.Debug("Diff file is removed by the coverage filters", "path", rel)
				change.Ignored = append(change.Ignored, rel)
			} else {
				logger.Debug("Diff file is not in any coverage report", "path", rel)
				change.NotInReports = append(change.NotInReports, rel)
			}
			continue
		}
		change.Measured = append(change.Measured, treePath)

		// Ensure the DiffInfo struct exists.
		if node.Diff == nil {
			node.Diff = &model.DiffInfo{
				AddedLines:    make(map[int]bool),
				ModifiedLines: make(map[int]bool),
			}
		}

		// Set the change kind.
		switch fileDiff.Kind {
		case "added":
			node.Diff.Kind = model.ChangeKindAdded
		default:
			node.Diff.Kind = model.ChangeKindModified
		}

		// Populate the line maps from hunks.
		for _, hunk := range fileDiff.Hunks {
			for _, offset := range hunk.AddedLineOffsets {
				lineNumber := hunk.NewStart + offset
				node.Diff.AddedLines[lineNumber] = true
			}
			for _, offset := range hunk.ModifiedLineOffsets {
				lineNumber := hunk.NewStart + offset
				node.Diff.ModifiedLines[lineNumber] = true
			}
		}
	}

	for _, group := range []*[]string{&change.Measured, &change.Ignored, &change.NotInReports, &change.Deleted} {
		slices.Sort(*group)
		*group = slices.Compact(*group)
	}

	logger.Info("Diff files sorted", "measured", len(change.Measured), "ignored_by_filters", len(change.Ignored),
		"not_in_reports", len(change.NotInReports), "deleted", len(change.Deleted))
	tree.Change = change
	return change
}

// collectFilesAndCoverage recursively traverses the tree to build a flat map of
// file paths to FileNode pointers and a set of paths for files with coverage.
func collectFilesAndCoverage(dir *model.DirNode, fileIndex map[string]*model.FileNode, coveredSet map[string]bool) {
	for _, file := range dir.Files {
		fileIndex[file.Path] = file // Use the full relative path stored in the node
		if file.Metrics.LinesValid > 0 {
			coveredSet[file.Path] = true
		}
	}
	for _, subDir := range dir.Subdirs {
		collectFilesAndCoverage(subDir, fileIndex, coveredSet)
	}
}
