// Path: internal/analyzer/tree_builder_test.go
package analyzer

import (
	"testing"

	"github.com/IgorBayerl/AdlerCov/internal/model"
	"github.com/IgorBayerl/AdlerCov/internal/parsers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildTree_SimpleCase(t *testing.T) {
	// Arrange
	parserResults := []*parsers.ParserResult{
		{
			ParserName: "TestParser",
			FileCoverage: []parsers.FileCoverage{
				{
					Path: "pkg/a.go",
					Lines: map[int]model.LineMetrics{
						10: {Hits: 1, TotalBranches: 0, CoveredBranches: 0}, // Covered
						11: {Hits: 0, TotalBranches: 0, CoveredBranches: 0}, // Not Covered
						12: {Hits: -1},                                      // Not Coverable
					},
				},
				{
					Path: "pkg/b.go",
					Lines: map[int]model.LineMetrics{
						5: {Hits: 5, TotalBranches: 2, CoveredBranches: 1}, // Partially Covered
					},
				},
			},
		},
	}

	builder := NewTreeBuilder()

	// Act
	tree, err := builder.BuildTree(parserResults)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, tree)

	// --- Assert Structure ---
	require.Contains(t, tree.Root.Subdirs, "pkg", "Root should contain 'pkg' subdirectory")
	pkgDir := tree.Root.Subdirs["pkg"]
	require.NotNil(t, pkgDir)
	assert.Equal(t, "pkg", pkgDir.Name)
	// CORRECTED: path.Join(".", "pkg") simplifies to "pkg". The test should expect this canonical form.
	assert.Equal(t, "pkg", pkgDir.Path)
	assert.Equal(t, tree.Root, pkgDir.Parent)

	require.Len(t, pkgDir.Files, 2)
	require.Contains(t, pkgDir.Files, "a.go")
	require.Contains(t, pkgDir.Files, "b.go")

	fileA := pkgDir.Files["a.go"]
	assert.Equal(t, "a.go", fileA.Name)
	assert.Equal(t, "pkg/a.go", fileA.Path)
	assert.Equal(t, pkgDir, fileA.Parent)

	// --- Assert File Metrics ---
	assert.Equal(t, 1, fileA.Metrics.LinesCovered, "File A should have 1 covered line")
	assert.Equal(t, 2, fileA.Metrics.LinesValid, "File A should have 2 coverable lines")
	assert.Equal(t, 0, fileA.Metrics.BranchesCovered, "File A should have 0 covered branches")
	assert.Equal(t, 0, fileA.Metrics.BranchesValid, "File A should have 0 total branches")

	fileB := pkgDir.Files["b.go"]
	assert.Equal(t, 1, fileB.Metrics.LinesCovered, "File B should have 1 covered line")
	assert.Equal(t, 1, fileB.Metrics.LinesValid, "File B should have 1 coverable line")
	assert.Equal(t, 1, fileB.Metrics.BranchesCovered, "File B should have 1 covered branch")
	assert.Equal(t, 2, fileB.Metrics.BranchesValid, "File B should have 2 total branches")

	// --- Assert Directory & Root Metrics ---
	expectedLinesCovered := 2    // 1 from A + 1 from B
	expectedLinesValid := 3      // 2 from A + 1 from B
	expectedBranchesCovered := 1 // 0 from A + 1 from B
	expectedBranchesValid := 2   // 0 from A + 2 from B

	assert.Equal(t, expectedLinesCovered, pkgDir.Metrics.LinesCovered, "Pkg dir lines covered mismatch")
	assert.Equal(t, expectedLinesValid, pkgDir.Metrics.LinesValid, "Pkg dir lines valid mismatch")
	assert.Equal(t, expectedBranchesCovered, pkgDir.Metrics.BranchesCovered, "Pkg dir branches covered mismatch")
	assert.Equal(t, expectedBranchesValid, pkgDir.Metrics.BranchesValid, "Pkg dir branches valid mismatch")

	assert.Equal(t, pkgDir.Metrics, tree.Root.Metrics, "Root metrics should match the single subdirectory")
	assert.Equal(t, tree.Root.Metrics, tree.Metrics, "Top-level tree metrics should match root node metrics")
}

func TestBuildTree_MergeCase(t *testing.T) {
	// Arrange: Two reports covering the same file
	result1 := &parsers.ParserResult{
		ParserName: "UnitTests",
		FileCoverage: []parsers.FileCoverage{
			{
				Path: "service/user.go",
				Lines: map[int]model.LineMetrics{
					20: {Hits: 1, TotalBranches: 2, CoveredBranches: 1}, // Hit branch 1
					25: {Hits: 1},                                       // Hit line
				},
			},
		},
	}
	result2 := &parsers.ParserResult{
		ParserName: "IntegrationTests",
		FileCoverage: []parsers.FileCoverage{
			{
				Path: "service/user.go",
				Lines: map[int]model.LineMetrics{
					20: {Hits: 1, TotalBranches: 2, CoveredBranches: 1}, // Hit branch 2 (in a real scenario, this would be represented differently, but for testing the sum, this works)
					30: {Hits: 1},                                       // Hit another line
				},
			},
		},
	}

	builder := NewTreeBuilder()

	// Act
	tree, err := builder.BuildTree([]*parsers.ParserResult{result1, result2})

	// Assert
	require.NoError(t, err)
	serviceDir := tree.Root.Subdirs["service"]
	userFile := serviceDir.Files["user.go"]

	// --- Assert Merged Line Metrics ---
	require.Len(t, userFile.Lines, 3, "Should have 3 unique lines with coverage")
	// Line 20 should have its hits and branches summed up
	assert.Equal(t, 2, userFile.Lines[20].Hits, "Hits for line 20 should be summed")
	assert.Equal(t, 2, userFile.Lines[20].CoveredBranches, "Covered branches for line 20 should be summed")
	assert.Equal(t, 2, userFile.Lines[20].TotalBranches, "Total branches should be consistent")
	// Other lines should exist
	assert.Equal(t, 1, userFile.Lines[25].Hits)
	assert.Equal(t, 1, userFile.Lines[30].Hits)

	// --- Assert Final Aggregated Metrics ---
	assert.Equal(t, 3, userFile.Metrics.LinesCovered, "File should have 3 covered lines after merge")
	assert.Equal(t, 3, userFile.Metrics.LinesValid, "File should have 3 valid lines after merge")
	assert.Equal(t, 2, userFile.Metrics.BranchesCovered, "File should have 2 covered branches after merge")
	assert.Equal(t, 2, userFile.Metrics.BranchesValid, "File should have 2 total branches after merge")

	assert.Equal(t, userFile.Metrics, tree.Metrics, "Root metrics should reflect the fully merged file metrics")
}

func TestBuildTree_ComplexHierarchy(t *testing.T) {
	// Arrange
	results := []*parsers.ParserResult{
		{
			ParserName: "Test",
			FileCoverage: []parsers.FileCoverage{
				{Path: "a.go", Lines: map[int]model.LineMetrics{1: {Hits: 1}}},                                                // 1/1
				{Path: "pkg1/b.go", Lines: map[int]model.LineMetrics{1: {Hits: 1}}},                                           // 1/1
				{Path: "pkg1/sub/c.go", Lines: map[int]model.LineMetrics{1: {Hits: 0, TotalBranches: 4, CoveredBranches: 1}}}, // 0/1 lines, 1/4 branches
			},
		},
	}
	builder := NewTreeBuilder()

	// Act
	tree, err := builder.BuildTree(results)

	// Assert
	require.NoError(t, err)

	// --- Assert Structure ---
	require.Contains(t, tree.Root.Files, "a.go")
	require.Contains(t, tree.Root.Subdirs, "pkg1")
	pkg1Dir := tree.Root.Subdirs["pkg1"]
	require.Contains(t, pkg1Dir.Files, "b.go")
	require.Contains(t, pkg1Dir.Subdirs, "sub")
	subDir := pkg1Dir.Subdirs["sub"]
	require.Contains(t, subDir.Files, "c.go")

	// --- Assert Aggregation at each level ---
	// Level 2: pkg1/sub
	assert.Equal(t, 0, subDir.Metrics.LinesCovered)
	assert.Equal(t, 1, subDir.Metrics.LinesValid)
	assert.Equal(t, 1, subDir.Metrics.BranchesCovered)
	assert.Equal(t, 4, subDir.Metrics.BranchesValid)

	// Level 1: pkg1
	assert.Equal(t, 1, pkg1Dir.Metrics.LinesCovered)    // 1 (b.go) + 0 (sub)
	assert.Equal(t, 2, pkg1Dir.Metrics.LinesValid)      // 1 (b.go) + 1 (sub)
	assert.Equal(t, 1, pkg1Dir.Metrics.BranchesCovered) // 0 (b.go) + 1 (sub)
	assert.Equal(t, 4, pkg1Dir.Metrics.BranchesValid)   // 0 (b.go) + 4 (sub)

	// Level 0: root
	assert.Equal(t, 2, tree.Root.Metrics.LinesCovered)    // 1 (a.go) + 1 (pkg1)
	assert.Equal(t, 3, tree.Root.Metrics.LinesValid)      // 1 (a.go) + 2 (pkg1)
	assert.Equal(t, 1, tree.Root.Metrics.BranchesCovered) // 0 (a.go) + 1 (pkg1)
	assert.Equal(t, 4, tree.Root.Metrics.BranchesValid)   // 0 (a.go) + 4 (pkg1)
}

func TestBuildTree_EmptyInput(t *testing.T) {
	// Arrange
	builder := NewTreeBuilder()

	// Act
	_, err := builder.BuildTree([]*parsers.ParserResult{})

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot build tree with no parser results")
}

func TestBuildTree_ParserResultWithNoFiles(t *testing.T) {
	// Arrange
	results := []*parsers.ParserResult{{ParserName: "Test", FileCoverage: []parsers.FileCoverage{}}}
	builder := NewTreeBuilder()

	// Act
	tree, err := builder.BuildTree(results)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, tree)
	assert.Empty(t, tree.Root.Files)
	assert.Empty(t, tree.Root.Subdirs)
	assert.Zero(t, tree.Metrics)
}
