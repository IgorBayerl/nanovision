package aggregator

import (
	"testing"

	"github.com/IgorBayerl/nanovision/internal/model"
)

func TestCalculateFileMethodMetrics_StrictPatchStatementCoverage(t *testing.T) {
	file := &model.FileNode{
		Path: "test.go",
		Lines: map[int]model.LineMetrics{
			10: {Hits: 1}, // Condition: Covered
			11: {Hits: 0}, // Body: Uncovered
		},
		Statements: []model.Statement{
			{StartLine: 10, EndLine: 11, Type: "if_statement"},
		},
		Diff: &model.DiffInfo{
			Kind:          model.ChangeKindModified,
			ModifiedLines: map[int]bool{11: true}, // Only the body was changed
			AddedLines:    map[int]bool{},
		},
		Metrics: model.CoverageMetrics{},
	}

	calculateFileMethodMetrics(file)

	// Assertions
	if file.Metrics.PatchLinesValid != 1 {
		t.Errorf("Expected 1 valid patch line, got %d", file.Metrics.PatchLinesValid)
	}
	if file.Metrics.PatchLinesCovered != 0 {
		t.Errorf("Expected 0 covered patch lines, got %d", file.Metrics.PatchLinesCovered)
	}

	// THE CRITICAL CHECK
	if file.Metrics.PatchStatementsCovered != 0 {
		t.Errorf("BUG DETECTED: Expected 0 covered patch statements because the modified line (11) was not hit, but got %d", file.Metrics.PatchStatementsCovered)
	}
}
