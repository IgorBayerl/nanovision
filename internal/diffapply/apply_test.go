package diffapply

import (
	"io"
	"log/slog"
	"testing"

	"github.com/IgorBayerl/nanovision/internal/diff"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createTestTree creates a simple model.SummaryTree for testing.
func createTestTree() *model.SummaryTree {
	root := &model.DirNode{
		Name:    "Root",
		Path:    ".",
		Subdirs: make(map[string]*model.DirNode),
		Files:   make(map[string]*model.FileNode),
	}
	tree := &model.SummaryTree{Root: root}

	// Create src directory
	srcDir := &model.DirNode{
		Name:    "src",
		Path:    "src",
		Parent:  root,
		Subdirs: make(map[string]*model.DirNode),
		Files:   make(map[string]*model.FileNode),
	}
	root.Subdirs["src"] = srcDir

	// File to be modified
	srcDir.Files["main.go"] = &model.FileNode{
		Name:   "main.go",
		Path:   "src/main.go",
		Parent: srcDir,
		Metrics: model.CoverageMetrics{
			LinesValid: 10, // Indicates it has coverage
		},
	}

	// File to be added (already in tree from a previous commit, now changed)
	srcDir.Files["new.go"] = &model.FileNode{
		Name:   "new.go",
		Path:   "src/new.go",
		Parent: srcDir,
	}

	// A file not mentioned in the diff
	srcDir.Files["untouched.go"] = &model.FileNode{
		Name:   "untouched.go",
		Path:   "src/untouched.go",
		Parent: srcDir,
	}

	return tree
}

func TestApply(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	t.Run("applies diff to new and modified files", func(t *testing.T) {
		// Arrange
		tree := createTestTree()
		diffData := &diff.DiffData{
			Files: []diff.FileDiff{
				{ // A new file
					OldPath: "/dev/null",
					NewPath: "src/new.go",
					Kind:    "added",
					Hunks: []diff.Hunk{
						{NewStart: 1, NewLines: 3, AddedLineOffsets: []int{0, 1, 2}}, // Lines 1, 2, 3
					},
				},
				{ // A modified file
					OldPath: "src/main.go",
					NewPath: "src/main.go",
					Kind:    "modified",
					Hunks: []diff.Hunk{
						// Hunk construction in test now assumes everything is AddedLineOffsets per new parser logic
						// line 10 modified -> treated as Added
						// line 11 added -> treated as Added
						// line 20 added -> treated as Added
						{NewStart: 10, NewLines: 2, ModifiedLineOffsets: []int{}, AddedLineOffsets: []int{0, 1}},
						{NewStart: 20, NewLines: 1, AddedLineOffsets: []int{0}},
					},
				},
			},
		}

		// Act
		Apply(tree, diffData, logger)

		// Assert
		newFileNode := tree.Root.Subdirs["src"].Files["new.go"]
		require.NotNil(t, newFileNode.Diff, "DiffInfo should be added to new file")
		assert.Equal(t, model.ChangeKindAdded, newFileNode.Diff.Kind)
		assert.Equal(t, map[int]bool{1: true, 2: true, 3: true}, newFileNode.Diff.AddedLines)
		assert.Empty(t, newFileNode.Diff.ModifiedLines)

		modifiedFileNode := tree.Root.Subdirs["src"].Files["main.go"]
		require.NotNil(t, modifiedFileNode.Diff, "DiffInfo should be added to modified file")
		assert.Equal(t, model.ChangeKindModified, modifiedFileNode.Diff.Kind)
		assert.Equal(t, map[int]bool{10: true, 11: true, 20: true}, modifiedFileNode.Diff.AddedLines)
		assert.Empty(t, modifiedFileNode.Diff.ModifiedLines)

		untouchedFileNode := tree.Root.Subdirs["src"].Files["untouched.go"]
		assert.Nil(t, untouchedFileNode.Diff, "Untouched file should not have DiffInfo")
	})

	t.Run("skips unmapped and filtered files", func(t *testing.T) {
		// Arrange
		tree := createTestTree()
		diffData := &diff.DiffData{
			Files: []diff.FileDiff{
				{ // This file is in the diff but not in our tree (e.g., filtered out)
					NewPath: "src/filtered.go",
					Kind:    "modified",
					Hunks:   []diff.Hunk{{NewStart: 1, AddedLineOffsets: []int{0}}},
				},
				{ // This file is in the tree and diff
					NewPath: "src/main.go",
					Kind:    "modified",
					Hunks:   []diff.Hunk{{NewStart: 5, AddedLineOffsets: []int{0}}},
				},
			},
		}

		// Act
		Apply(tree, diffData, logger)

		// Assert
		// No node for 'filtered.go' exists, so we just check it didn't crash
		mainFileNode := tree.Root.Subdirs["src"].Files["main.go"]
		require.NotNil(t, mainFileNode.Diff)
		assert.Equal(t, map[int]bool{5: true}, mainFileNode.Diff.AddedLines)
	})

	t.Run("is idempotent", func(t *testing.T) {
		// Arrange
		tree := createTestTree()
		diffData := &diff.DiffData{
			Files: []diff.FileDiff{
				{
					NewPath: "src/main.go",
					Kind:    "modified",
					Hunks:   []diff.Hunk{{NewStart: 8, AddedLineOffsets: []int{0}}},
				},
			},
		}

		// Act
		Apply(tree, diffData, logger)
		firstRunDiff := *tree.Root.Subdirs["src"].Files["main.go"].Diff

		// Apply a second time
		Apply(tree, diffData, logger)
		secondRunDiff := *tree.Root.Subdirs["src"].Files["main.go"].Diff

		// Assert
		assert.Equal(t, firstRunDiff, secondRunDiff, "Applying the same diff twice should yield the same result")
	})

	t.Run("handles nil and empty inputs gracefully", func(t *testing.T) {
		tree := createTestTree()

		// This should not panic
		Apply(nil, &diff.DiffData{}, logger)
		Apply(tree, nil, logger)
		Apply(tree, &diff.DiffData{}, logger)

		// Check that nothing was modified
		assert.Nil(t, tree.Root.Subdirs["src"].Files["main.go"].Diff)
	})
}
