package reporter_rawjson_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/IgorBayerl/AdlerCov/internal/model"
	"github.com/IgorBayerl/AdlerCov/internal/reporter/reporter_rawjson"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRawJsonReportBuilder_CreateReport(t *testing.T) {
	// Arrange
	tmpDir := t.TempDir()
	builder := reporter_rawjson.NewRawJsonReportBuilder(tmpDir)

	// Create a sample hydrated tree with parent pointers
	sampleTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC).Unix()
	rootNode := &model.DirNode{
		Name: "Root",
		Path: ".",
		Metrics: model.CoverageMetrics{
			LinesCovered:  5,
			LinesValid:    10,
			BranchesValid: 2,
		},
		Subdirs: make(map[string]*model.DirNode),
		Files:   make(map[string]*model.FileNode),
	}

	subdirNode := &model.DirNode{
		Name:   "subdir",
		Path:   "subdir",
		Parent: rootNode,
		Metrics: model.CoverageMetrics{
			LinesCovered: 2,
			LinesValid:   4,
		},
		Files: make(map[string]*model.FileNode),
	}
	rootNode.Subdirs["subdir"] = subdirNode

	file1Node := &model.FileNode{
		Name:   "file1.go",
		Path:   "file1.go",
		Parent: rootNode,
		Metrics: model.CoverageMetrics{
			LinesCovered:  3,
			LinesValid:    6,
			BranchesValid: 2,
		},
		Lines: map[int]model.LineMetrics{
			5: {Hits: 5},
		},
		Methods: []model.MethodMetrics{
			{Name: "myFunc", StartLine: 4, EndLine: 8, CyclomaticComplexity: 2},
		},
		SourceDir: "/src",
	}
	rootNode.Files["file1.go"] = file1Node

	file2Node := &model.FileNode{
		Name:   "file2.go",
		Path:   "subdir/file2.go",
		Parent: subdirNode,
		Metrics: model.CoverageMetrics{
			LinesCovered: 2,
			LinesValid:   4,
		},
		Lines: map[int]model.LineMetrics{
			10: {Hits: 1},
			12: {Hits: 0},
		},
		// *** FIX: Add sample methods data to simulate hydration ***
		Methods: []model.MethodMetrics{
			{Name: "helperFunc", StartLine: 9, EndLine: 13, CyclomaticComplexity: 1},
		},
	}
	subdirNode.Files["file2.go"] = file2Node

	tree := &model.SummaryTree{
		Root:        rootNode,
		Metrics:     rootNode.Metrics,
		Timestamp:   sampleTime,
		ParserName:  "GoCover",
		ReportFiles: []string{"coverage.out"},
		SourceFiles: []string{"/src"},
	}

	// Act
	err := builder.CreateReport(tree)

	// Assert
	require.NoError(t, err, "CreateReport should not fail with the cycle-breaking json tag")

	// Verify the file was created
	reportPath := filepath.Join(tmpDir, "RawJson.json")
	fileContent, err := os.ReadFile(reportPath)
	require.NoError(t, err, "Report file should exist and be readable")

	// Unmarshal the actual content and verify its structure
	var actualTree model.SummaryTree
	err = json.Unmarshal(fileContent, &actualTree)
	require.NoError(t, err, "Failed to unmarshal actual JSON content")

	// We can now compare the deserialized tree with our expected structure.
	// The deserialized tree will have `nil` for Parent fields, which is correct.
	assert.Equal(t, "Root", actualTree.Root.Name)
	assert.Equal(t, 5, actualTree.Root.Metrics.LinesCovered)
	require.Contains(t, actualTree.Root.Files, "file1.go")
	assert.Equal(t, "file1.go", actualTree.Root.Files["file1.go"].Name)
	assert.Equal(t, 2, actualTree.Root.Files["file1.go"].Methods[0].CyclomaticComplexity)
	require.Contains(t, actualTree.Root.Subdirs, "subdir")
	assert.Nil(t, actualTree.Root.Subdirs["subdir"].Parent, "Parent field should not be serialized")
	require.Contains(t, actualTree.Root.Subdirs["subdir"].Files, "file2.go")
	assert.Equal(t, 1, actualTree.Root.Subdirs["subdir"].Files["file2.go"].Lines[10].Hits)
	assert.Nil(t, actualTree.Root.Subdirs["subdir"].Files["file2.go"].Parent, "Parent field should not be serialized")

	// *** FIX: Add assertion to verify the methods field for file2.go is now present ***
	require.NotNil(t, actualTree.Root.Subdirs["subdir"].Files["file2.go"].Methods)
	require.Len(t, actualTree.Root.Subdirs["subdir"].Files["file2.go"].Methods, 1)
	assert.Equal(t, "helperFunc", actualTree.Root.Subdirs["subdir"].Files["file2.go"].Methods[0].Name)
}
