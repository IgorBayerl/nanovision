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
	tmpDir := t.TempDir()
	builder := reporter_rawjson.NewRawJsonReportBuilder(tmpDir)

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

	// CORRECTED: Helper variables for pointer assignment
	cyclo1 := 2
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
			{Name: "myFunc", StartLine: 4, EndLine: 8, CyclomaticComplexity: &cyclo1},
		},
		SourceDir: "/src",
	}
	rootNode.Files["file1.go"] = file1Node

	// CORRECTED: Helper variables for pointer assignment
	cyclo2 := 1
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
		Methods: []model.MethodMetrics{
			{Name: "helperFunc", StartLine: 9, EndLine: 13, CyclomaticComplexity: &cyclo2},
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

	err := builder.CreateReport(tree)
	require.NoError(t, err)

	reportPath := filepath.Join(tmpDir, "RawJson.json")
	fileContent, err := os.ReadFile(reportPath)
	require.NoError(t, err)

	var actualTree model.SummaryTree
	err = json.Unmarshal(fileContent, &actualTree)
	require.NoError(t, err)

	assert.Equal(t, "Root", actualTree.Root.Name)
	assert.NotNil(t, actualTree.Root.Files["file1.go"].Methods[0].CyclomaticComplexity)
	assert.Equal(t, 2, *actualTree.Root.Files["file1.go"].Methods[0].CyclomaticComplexity)
	assert.NotNil(t, actualTree.Root.Subdirs["subdir"].Files["file2.go"].Methods[0].CyclomaticComplexity)
	assert.Equal(t, 1, *actualTree.Root.Subdirs["subdir"].Files["file2.go"].Methods[0].CyclomaticComplexity)
}
