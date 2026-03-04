package aggregator

import (
	"testing"

	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestCalculateFileMethodMetrics_StrictPatchStatementCoverage(t *testing.T) {
	file := &model.FileNode{
		Path: "test.go",
		Lines: map[int]model.LineMetrics{
			10: {Hits: 1}, // Condition: Covered
			11: {Hits: 0}, // Body: Uncovered, but part of the same statement
		},
		Statements: []model.Statement{
			{StartLine: 10, EndLine: 11, Type: "call_expression"},
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

	// With atomic statements, if a statement is hit on any line, the statement is hit.
	if file.Metrics.PatchStatementsCovered != 1 {
		t.Errorf("Expected 1 covered patch statement, got %d", file.Metrics.PatchStatementsCovered)
	}
}

func TestIsLineInPatch(t *testing.T) {
	diff := &model.DiffInfo{
		AddedLines:    map[int]bool{10: true},
		ModifiedLines: map[int]bool{20: true},
	}

	assert.True(t, isLineInPatch(10, diff), "Added line should be in patch")
	assert.True(t, isLineInPatch(20, diff), "Modified line should be in patch")
	assert.False(t, isLineInPatch(15, diff), "Unchanged line should not be in patch")
	assert.False(t, isLineInPatch(10, nil), "Nil diff means no lines are in patch")
}

func TestEvaluateStatementPatchStatus(t *testing.T) {
	file := &model.FileNode{
		Diff: &model.DiffInfo{
			AddedLines:    map[int]bool{10: true},
			ModifiedLines: map[int]bool{12: true},
		},
		Lines: map[int]model.LineMetrics{
			10: {Hits: 0}, // Added but uncovered
			12: {Hits: 5}, // Modified and covered
		},
	}

	tests := []struct {
		name            string
		stmt            model.Statement
		expectedInPatch bool
		expectedCovered bool
	}{
		{
			name:            "Statement fully outside patch",
			stmt:            model.Statement{StartLine: 1, EndLine: 5},
			expectedInPatch: false,
			expectedCovered: false,
		},
		{
			name:            "Statement intersects an added line but is uncovered",
			stmt:            model.Statement{StartLine: 8, EndLine: 10},
			expectedInPatch: true,
			expectedCovered: false,
		},
		{
			name:            "Statement intersects a modified line and is covered",
			stmt:            model.Statement{StartLine: 11, EndLine: 15},
			expectedInPatch: true,
			expectedCovered: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			inPatch, isCovered := evaluateStatementPatchStatus(tc.stmt, file)
			assert.Equal(t, tc.expectedInPatch, inPatch, "inPatch mismatch")
			assert.Equal(t, tc.expectedCovered, isCovered, "isCovered mismatch")
		})
	}
}

func TestEvaluateStatementPatchStatus_TrailingEmptyLine(t *testing.T) {
	file := &model.FileNode{
		Diff: &model.DiffInfo{
			AddedLines: map[int]bool{10: true}, // Empty line added
		},
		Lines: map[int]model.LineMetrics{
			9: {Hits: 1}, // Coverable line, not in patch
			// Line 10 is NOT coverable (not in map)
		},
	}

	stmt := model.Statement{StartLine: 9, EndLine: 10}

	inPatch, isCovered := evaluateStatementPatchStatus(stmt, file)
	assert.False(t, inPatch, "Statement should not be in patch if only non-coverable lines are modified")
	assert.False(t, isCovered, "Covered doesn't matter if not in patch, should be false")
}

func TestEvaluateStatementPatchStatus_CoveredOnUnpatchedLine(t *testing.T) {
	file := &model.FileNode{
		Diff: &model.DiffInfo{
			AddedLines: map[int]bool{10: true},
		},
		Lines: map[int]model.LineMetrics{
			9:  {Hits: 1}, // Covered line, NOT in patch
			10: {Hits: 0}, // Uncovered line, IN patch
		},
	}

	stmt := model.Statement{StartLine: 9, EndLine: 10}

	inPatch, isCovered := evaluateStatementPatchStatus(stmt, file)
	assert.True(t, inPatch, "Statement is in patch because line 10 is coverable and patched")
	assert.True(t, isCovered, "Statement should be considered covered because line 9 is covered")
}

func TestCalculateMethodPatchMetrics_AddedFileEdgeCase(t *testing.T) {
	file := &model.FileNode{
		Diff: &model.DiffInfo{
			Kind: model.ChangeKindAdded, // The entire file is new
		},
	}

	method := &model.MethodMetrics{
		StartLine:         1,
		EndLine:           10,
		LinesValid:        5,
		LinesCovered:      3,
		StatementsValid:   2,
		StatementsCovered: 2,
	}

	calculateMethodPatchMetrics(file, method)

	// In an added file, standard metrics should be copied exactly to patch metrics.
	assert.Equal(t, "added", method.DiffStatus)
	assert.Equal(t, 5, method.PatchLinesValid)
	assert.Equal(t, 3, method.PatchLinesCovered)
	assert.Equal(t, 2, method.PatchStatementsValid)
	assert.Equal(t, 2, method.PatchStatementsCovered)
}

func TestCalculateMethodPatchMetrics_ModifiedFile(t *testing.T) {
	file := &model.FileNode{
		Diff: &model.DiffInfo{
			Kind:          model.ChangeKindModified,
			AddedLines:    map[int]bool{12: true},
			ModifiedLines: map[int]bool{15: true},
		},
		Lines: map[int]model.LineMetrics{
			10: {Hits: 1}, // Untouched base line
			12: {Hits: 0}, // Added, not covered
			15: {Hits: 1}, // Modified, covered
		},
		Statements: []model.Statement{
			{StartLine: 12, EndLine: 15}, // Span covering both changes
		},
	}

	method := &model.MethodMetrics{
		StartLine: 10,
		EndLine:   20,
		// Standard metrics represent total reality over the whole file
		LinesValid:        3,
		LinesCovered:      2,
		StatementsValid:   1,
		StatementsCovered: 1,
	}

	calculateMethodPatchMetrics(file, method)

	assert.Equal(t, "modified", method.DiffStatus)
	// Only 12 and 15 are in the patch and coverable
	assert.Equal(t, 2, method.PatchLinesValid)
	// Only 15 is covered
	assert.Equal(t, 1, method.PatchLinesCovered)

	// Our single statement spans the patch and reaches a covered mod line
	assert.Equal(t, 1, method.PatchStatementsValid)
	assert.Equal(t, 1, method.PatchStatementsCovered)
}

func TestAggregatePatchMethodMetrics_AddedFile(t *testing.T) {
	file := &model.FileNode{
		Diff: &model.DiffInfo{Kind: model.ChangeKindAdded},
	}
	method := &model.MethodMetrics{
		LinesValid:        10,
		LinesCovered:      5,
		StatementsValid:   5,
		StatementsCovered: 5,
	}

	aggregatePatchMethodMetrics(file, method)

	assert.Equal(t, 1, file.Metrics.PatchMethodsValid)
	assert.Equal(t, 1, file.Metrics.PatchMethodsHit)
	assert.Equal(t, 1, file.Metrics.PatchStatementMethodsValid)
	assert.Equal(t, 1, file.Metrics.PatchStatementMethodsHit)
}

func TestAggregatePatchMethodMetrics_ModifiedFile(t *testing.T) {
	file := &model.FileNode{
		Diff: &model.DiffInfo{Kind: model.ChangeKindModified},
	}
	method := &model.MethodMetrics{
		PatchLinesValid:   5,
		PatchLinesCovered: 0,
		StatementsValid:   3,
		StatementsCovered: 2,
	}

	aggregatePatchMethodMetrics(file, method)

	assert.Equal(t, 1, file.Metrics.PatchMethodsValid)
	assert.Equal(t, 0, file.Metrics.PatchMethodsHit)
	assert.Equal(t, 1, file.Metrics.PatchStatementMethodsValid)
	assert.Equal(t, 1, file.Metrics.PatchStatementMethodsHit)
}

func TestAggregatePatchMethodMetrics_ModifiedFile_NoPatchLines(t *testing.T) {
	file := &model.FileNode{
		Diff: &model.DiffInfo{Kind: model.ChangeKindModified},
	}
	method := &model.MethodMetrics{
		PatchLinesValid:   0,
		PatchLinesCovered: 0,
		StatementsValid:   3,
		StatementsCovered: 2,
	}

	aggregatePatchMethodMetrics(file, method)

	// Since there are no patch lines, no metrics should be summed.
	assert.Equal(t, 0, file.Metrics.PatchMethodsValid)
	assert.Equal(t, 0, file.Metrics.PatchMethodsHit)
	assert.Equal(t, 0, file.Metrics.PatchStatementMethodsValid)
	assert.Equal(t, 0, file.Metrics.PatchStatementMethodsHit)
}

func TestMaxCyclomaticComplexity_FileLevel(t *testing.T) {
	// Build a file with two methods having different complexity values
	cyclo5 := 5
	cyclo12 := 12
	file := &model.FileNode{
		Path: "test.go",
		Lines: map[int]model.LineMetrics{
			1: {Hits: 1}, 2: {Hits: 1}, 3: {Hits: 1}, 4: {Hits: 1},
		},
		Methods: []model.MethodMetrics{
			{
				Name:                 "FuncA",
				StartLine:            1,
				EndLine:              2,
				CyclomaticComplexity: &cyclo5,
				LinesValid:          2,
				LinesCovered:         2,
			},
			{
				Name:                 "FuncB",
				StartLine:            3,
				EndLine:              4,
				CyclomaticComplexity: &cyclo12,
				LinesValid:          2,
				LinesCovered:         1,
			},
		},
		Metrics: model.CoverageMetrics{
			LinesValid:   4,
			LinesCovered: 3,
		},
	}

	calculateFileMethodMetrics(file)

	assert.Equal(t, 12, file.Metrics.MaxCyclomaticComplexity,
		"File-level MaxCyclomaticComplexity should be the max across methods")
}

func TestMaxCyclomaticComplexity_DirPropagation(t *testing.T) {
	// Build a mini tree: dir with two files having different max complexities
	cyclo3 := 3
	cyclo7 := 7
	cyclo15 := 15

	dir := &model.DirNode{
		Name:    "root",
		Subdirs: map[string]*model.DirNode{},
		Files: map[string]*model.FileNode{
			"a.go": {
				Path: "a.go",
				Lines: map[int]model.LineMetrics{
					1: {Hits: 1}, 2: {Hits: 1},
				},
				Methods: []model.MethodMetrics{
					{
						Name: "FuncA", StartLine: 1, EndLine: 2,
						CyclomaticComplexity: &cyclo7,
						LinesValid: 2, LinesCovered: 2,
					},
				},
				Metrics: model.CoverageMetrics{LinesValid: 2, LinesCovered: 2},
			},
			"b.go": {
				Path: "b.go",
				Lines: map[int]model.LineMetrics{
					1: {Hits: 1}, 2: {Hits: 0},
				},
				Methods: []model.MethodMetrics{
					{
						Name: "FuncB", StartLine: 1, EndLine: 1,
						CyclomaticComplexity: &cyclo3,
						LinesValid: 1, LinesCovered: 1,
					},
					{
						Name: "FuncC", StartLine: 2, EndLine: 2,
						CyclomaticComplexity: &cyclo15,
						LinesValid: 1, LinesCovered: 0,
					},
				},
				Metrics: model.CoverageMetrics{LinesValid: 2, LinesCovered: 1},
			},
		},
	}

	tree := &model.SummaryTree{Root: dir}
	AggregateMetricsAfterEnrichment(tree)

	// a.go should have max 7
	assert.Equal(t, 7, dir.Files["a.go"].Metrics.MaxCyclomaticComplexity)
	// b.go should have max 15 (from FuncC)
	assert.Equal(t, 15, dir.Files["b.go"].Metrics.MaxCyclomaticComplexity)
	// Dir should have max 15 (propagated via max, not sum)
	assert.Equal(t, 15, dir.Metrics.MaxCyclomaticComplexity)
	assert.Equal(t, 15, tree.Metrics.MaxCyclomaticComplexity)
}

