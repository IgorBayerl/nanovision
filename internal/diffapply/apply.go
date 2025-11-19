package diffapply

import (
	"log/slog"

	"github.com/IgorBayerl/nanovision/internal/diff"
	"github.com/IgorBayerl/nanovision/internal/model"
)

// Apply annotates a model.SummaryTree with diff information.
// It resolves paths from the diff data, finds corresponding file nodes in the tree,
// and attaches change details like kind (added/modified) and changed line numbers.
func Apply(tree *model.SummaryTree, dd *diff.DiffData, logger *slog.Logger) {
	if tree == nil || tree.Root == nil || dd == nil || len(dd.Files) == 0 {
		return
	}

	// 1. Build file index and covered set from the tree.
	fileIndex := make(map[string]*model.FileNode)
	coveredSet := make(map[string]bool)
	collectFilesAndCoverage(tree.Root, fileIndex, coveredSet)

	// 2. Build the resolver with heuristics.
	resolver := BuildResolver(dd, fileIndex, coveredSet, logger)

	// 3. Iterate through each file in the diff and apply its changes.
	for _, fileDiff := range dd.Files {
		// Resolve the diff path to a path in our coverage tree.
		treePath, ok := resolver.Resolve(fileDiff.NewPath)
		if !ok {
			logger.Debug("Skipping unmapped diff file", "path", fileDiff.NewPath)
			continue
		}

		// Find the corresponding node in the tree.
		node, exists := fileIndex[treePath]
		if !exists {
			// This can happen if the resolver finds a match but the file was filtered
			// out during the initial tree build. It's safe to skip.
			logger.Debug("Skipping diff for file not present in final tree", "path", treePath)
			continue
		}

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
