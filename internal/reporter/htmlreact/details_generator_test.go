package htmlreact

import (
	"testing"

	"github.com/IgorBayerl/nanovision/internal/calculator"
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestBuildReportFilter(t *testing.T) {
	builder := &HtmlReactReportBuilder{}
	tree := &model.SummaryTree{
		ReportNames: []string{"report_A.json", "report_B.json", "report_C.json"},
	}

	tests := []struct {
		name                 string
		fileNode             *model.FileNode
		expectedReportsSize  int
		expectedGlobalMapLen int
		expectedFallback     bool
	}{
		{
			name: "Partial overlap - only report A and C have hits",
			fileNode: &model.FileNode{
				Lines: map[int]model.LineMetrics{
					1: {ReportHits: []int{1, 0, 5}}, // hits in A & C
					2: {ReportHits: []int{0, 0, 2}}, // hit in C
				},
			},
			expectedReportsSize:  2,
			expectedGlobalMapLen: 2,
		},
		{
			name: "No hits anywhere - trigger fallback to show all reports",
			fileNode: &model.FileNode{
				Lines: map[int]model.LineMetrics{
					1: {ReportHits: []int{0, 0, 0}},
				},
			},
			expectedReportsSize:  3,
			expectedGlobalMapLen: 0, // Fallback implies map is nil or not used
			expectedFallback:     true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			reportsList, globalMap := builder.buildReportFilter(tree, tc.fileNode)
			assert.Len(t, reportsList, tc.expectedReportsSize)
			if tc.expectedFallback {
				assert.Nil(t, globalMap)
			} else {
				assert.Len(t, globalMap, tc.expectedGlobalMapLen)
			}
		})
	}
}

func TestBuildLineDetails(t *testing.T) {
	builder := &HtmlReactReportBuilder{}
	fileNode := &model.FileNode{
		Diff: &model.DiffInfo{
			AddedLines:    map[int]bool{1: true},
			ModifiedLines: map[int]bool{2: true},
		},
		Lines: map[int]model.LineMetrics{
			1: {Hits: 5, ReportHits: []int{5, 0}},
			2: {Hits: 0, ReportHits: []int{0, 0}, TotalBranches: 2, CoveredBranches: 1},
			3: {Hits: -1}, // Not coverable
		},
	}
	sourceLines := []string{"add", "mod", "none"}
	reportsList := []report{{Name: "Report A"}, {Name: "Report B"}}
	globalToLocalMap := map[int]int{0: 0, 1: 1}

	details := builder.buildLineDetails(fileNode, sourceLines, reportsList, globalToLocalMap, 2)

	assert.Len(t, details, 3)

	// Line 1: Added, Covered
	assert.Equal(t, "added", details[0].DiffStatus)
	assert.Equal(t, StatusCovered, details[0].Status)
	assert.Equal(t, []int{5, 0}, details[0].Hits)

	// Line 2: Modified, Uncovered, Partial Branches
	assert.Equal(t, "modified", details[1].DiffStatus)
	assert.Equal(t, StatusPartial, details[1].Status) // due to branches (1 of 2)
	assert.NotNil(t, details[1].BranchInfo)
	assert.Equal(t, 1, details[1].BranchInfo.Covered)

	// Line 3: Not coverable
	assert.Equal(t, "", details[2].DiffStatus)
	assert.Equal(t, StatusNotCoverable, details[2].Status)
}

func TestBuildFileTotals(t *testing.T) {
	builder := &HtmlReactReportBuilder{
		config: &config.AppConfig{
			ActiveFileMetrics: map[config.MetricKey]bool{
				config.LineCoverage:            true,
				config.BranchCoverage:          true,
				config.PatchStatementCoverage:  true,
				config.MaxCyclomaticComplexity: true,
			},
		},
	}

	totalBranches := 10
	coveredBranches := 5
	maxCyclo := 12

	fileNode := &model.FileNode{
		Metrics: model.CoverageMetrics{
			TotalLines:      150,
			LinesValid:      100,
			LinesCovered:    50,
			BranchesValid:   10,
			BranchesCovered: 5,
			StatementsValid: 1,
		},
		Diff: &model.DiffInfo{}, // Trigger patch generation
		Statuses: map[config.MetricKey]string{
			config.LineCoverage: "danger",
		},
	}

	tree := &model.SummaryTree{Root: &model.DirNode{
		Metrics: fileNode.Metrics,
		Files:   map[string]*model.FileNode{"test.go": fileNode},
	}}
	calculator.CalculateTree(tree, builder.config.ActiveFileMetrics, nil)

	totalsData := builder.buildFileTotals(fileNode, totalBranches, coveredBranches, maxCyclo)

	// Explicit assignment tests
	assert.NotNil(t, totalsData.LineCoverage)
	assert.Equal(t, 150, totalsData.LineCoverage.Total)

	assert.NotNil(t, totalsData.MethodBranchCoverage)
	assert.Equal(t, 10, totalsData.MethodBranchCoverage.Total)

	assert.NotNil(t, totalsData.MaxCyclomaticComplexity)
	assert.Equal(t, 12.0, totalsData.MaxCyclomaticComplexity.Value)

	// Edge case: PatchStatementCoverage explicitly forced to 100% when nil
	assert.NotNil(t, totalsData.PatchStatementCoverage)
	assert.Equal(t, 100.0, totalsData.PatchStatementCoverage.Percentage)
	assert.Equal(t, 0, totalsData.PatchStatementCoverage.Total)
}
