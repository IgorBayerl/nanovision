package aggregator

import (
	"testing"

	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// resolve is the client-side rule, mirrored here so the tests pin the contract
// the UI relies on: entities count as covered when every mask in their bucket
// intersects the selection.
func resolve(buckets []ReportBucket, selection uint32) (covered, total int) {
	for _, bucket := range buckets {
		total += bucket.Count
		satisfied := true
		for _, mask := range bucket.Masks {
			if mask&selection == 0 {
				satisfied = false
				break
			}
		}
		if satisfied {
			covered += bucket.Count
		}
	}
	return covered, total
}

func allFileMetrics() map[config.MetricKey]bool {
	return map[config.MetricKey]bool{
		config.LineCoverage:                 true,
		config.PatchLineCoverage:            true,
		config.StatementCoverage:            true,
		config.PatchStatementCoverage:       true,
		config.MethodsHit:                   true,
		config.MethodsFullyCovered:          true,
		config.PatchMethodsHit:              true,
		config.StatementMethodsHit:          true,
		config.StatementMethodsFullyCovered: true,
		config.PatchStatementMethodsHit:     true,
	}
}

// two methods, two reports: report A covers the first method, report B covers
// one line of the second.
func twoReportFile() *model.FileNode {
	return &model.FileNode{
		Path: "app.go",
		Lines: map[int]model.LineMetrics{
			1: {Hits: 3, ReportHits: []int{3, 0}},
			2: {Hits: 1, ReportHits: []int{1, 0}},
			3: {Hits: 2, ReportHits: []int{0, 2}},
			4: {Hits: 0, ReportHits: []int{0, 0}},
			5: {Hits: -1}, // a comment, never counted
		},
		Statements: []model.Statement{
			{StartLine: 1, EndLine: 1},
			{StartLine: 2, EndLine: 2},
			{StartLine: 3, EndLine: 3},
			{StartLine: 4, EndLine: 4},
		},
		Methods: []model.MethodMetrics{
			{Name: "first", StartLine: 1, EndLine: 2},
			{Name: "second", StartLine: 3, EndLine: 4},
		},
	}
}

// enrich fills the per-method counts the enricher normally produces, so the
// aggregator has something to roll up.
func enrich(file *model.FileNode) {
	for i := range file.Methods {
		method := &file.Methods[i]
		for ln := method.StartLine; ln <= method.EndLine; ln++ {
			lm, ok := file.Lines[ln]
			if !ok || lm.Hits < 0 {
				continue
			}
			method.LinesValid++
			if lm.Hits > 0 {
				method.LinesCovered++
			}
		}
		for _, stmt := range file.Statements {
			if stmt.StartLine < method.StartLine || stmt.EndLine > method.EndLine {
				continue
			}
			method.StatementsValid++
			for ln := stmt.StartLine; ln <= stmt.EndLine; ln++ {
				if lm, ok := file.Lines[ln]; ok && lm.Hits > 0 {
					method.StatementsCovered++
					break
				}
			}
		}
	}
}

func TestBuildFileReportIndex_MatchesAggregatorWhenAllReportsSelected(t *testing.T) {
	file := twoReportFile()
	enrich(file)
	calculateFileMethodMetrics(file)
	// The tree builder normally fills these; the aggregator only owns the
	// method-derived counts.
	file.Metrics.LinesValid = 4
	file.Metrics.LinesCovered = 3
	file.Metrics.StatementsValid = 4
	file.Metrics.StatementsCovered = 3

	idx := BuildFileReportIndex(file, 2, allFileMetrics())
	require.NotEmpty(t, idx)

	const both uint32 = 0b11

	covered, total := resolve(idx[config.LineCoverage], both)
	assert.Equal(t, file.Metrics.LinesCovered, covered)
	assert.Equal(t, file.Metrics.LinesValid, total)

	covered, total = resolve(idx[config.StatementCoverage], both)
	assert.Equal(t, file.Metrics.StatementsCovered, covered)
	assert.Equal(t, file.Metrics.StatementsValid, total)

	covered, total = resolve(idx[config.MethodsHit], both)
	assert.Equal(t, file.Metrics.MethodsHit, covered)
	assert.Equal(t, file.Metrics.MethodsValid, total)

	covered, total = resolve(idx[config.MethodsFullyCovered], both)
	assert.Equal(t, file.Metrics.MethodsFullyCovered, covered)
	assert.Equal(t, file.Metrics.MethodsValid, total)

	covered, _ = resolve(idx[config.StatementMethodsHit], both)
	assert.Equal(t, file.Metrics.StatementMethodsHit, covered)

	covered, _ = resolve(idx[config.StatementMethodsFullyCovered], both)
	assert.Equal(t, file.Metrics.StatementMethodsFullyCovered, covered)
}

func TestBuildFileReportIndex_SubsetSelections(t *testing.T) {
	idx := BuildFileReportIndex(twoReportFile(), 2, allFileMetrics())
	require.NotEmpty(t, idx)

	const (
		none    uint32 = 0b00
		onlyA   uint32 = 0b01
		onlyB   uint32 = 0b10
		bothAaB uint32 = 0b11
	)

	// Lines 1 and 2 come from report A, line 3 from report B, line 4 from neither.
	covered, total := resolve(idx[config.LineCoverage], onlyA)
	assert.Equal(t, 2, covered)
	assert.Equal(t, 4, total, "totals stay fixed: the selection changes coverage, not what exists")

	covered, _ = resolve(idx[config.LineCoverage], onlyB)
	assert.Equal(t, 1, covered)

	covered, _ = resolve(idx[config.LineCoverage], none)
	assert.Equal(t, 0, covered)

	// "first" needs both its lines and both come from A; "second" has an
	// uncovered line, so no selection ever makes it fully covered.
	covered, total = resolve(idx[config.MethodsFullyCovered], onlyA)
	assert.Equal(t, 1, covered)
	assert.Equal(t, 2, total)

	covered, _ = resolve(idx[config.MethodsFullyCovered], bothAaB)
	assert.Equal(t, 1, covered)

	// "second" is only reached by report B.
	covered, _ = resolve(idx[config.MethodsHit], onlyA)
	assert.Equal(t, 1, covered)

	covered, _ = resolve(idx[config.MethodsHit], onlyB)
	assert.Equal(t, 1, covered)

	covered, _ = resolve(idx[config.MethodsHit], bothAaB)
	assert.Equal(t, 2, covered)
}

func TestBuildFileReportIndex_PatchMetrics(t *testing.T) {
	file := twoReportFile()
	file.Diff = &model.DiffInfo{
		Kind:          model.ChangeKindModified,
		AddedLines:    map[int]bool{3: true},
		ModifiedLines: map[int]bool{4: true},
	}
	enrich(file)
	calculateFileMethodMetrics(file)

	idx := BuildFileReportIndex(file, 2, allFileMetrics())
	require.NotEmpty(t, idx)

	const (
		onlyA uint32 = 0b01
		onlyB uint32 = 0b10
		both  uint32 = 0b11
	)

	covered, total := resolve(idx[config.PatchLineCoverage], both)
	assert.Equal(t, file.Metrics.PatchLinesCovered, covered)
	assert.Equal(t, file.Metrics.PatchLinesValid, total)

	// Only line 3 is in the patch and covered, and only report B covers it.
	covered, _ = resolve(idx[config.PatchLineCoverage], onlyA)
	assert.Equal(t, 0, covered)

	covered, _ = resolve(idx[config.PatchLineCoverage], onlyB)
	assert.Equal(t, 1, covered)

	// Only "second" has changed lines, so it is the only method in the patch.
	covered, total = resolve(idx[config.PatchMethodsHit], both)
	assert.Equal(t, file.Metrics.PatchMethodsHit, covered)
	assert.Equal(t, file.Metrics.PatchMethodsValid, total)
}

func TestBuildFileReportIndex_BucketsCompressRepeatedPatterns(t *testing.T) {
	lines := make(map[int]model.LineMetrics, 500)
	for ln := 1; ln <= 500; ln++ {
		lines[ln] = model.LineMetrics{Hits: 1, ReportHits: []int{1, 0}}
	}
	lines[500] = model.LineMetrics{Hits: 0, ReportHits: []int{0, 0}}

	idx := BuildFileReportIndex(&model.FileNode{Lines: lines}, 2, allFileMetrics())

	assert.Len(t, idx[config.LineCoverage], 2, "500 lines with two patterns collapse to two buckets")
}

func TestBuildFileReportIndex_SkippedWhenReportsCannotBeMasked(t *testing.T) {
	file := twoReportFile()

	assert.Nil(t, BuildFileReportIndex(file, 0, allFileMetrics()))
	assert.Nil(t, BuildFileReportIndex(file, MaxIndexedReports+1, allFileMetrics()))
	assert.Nil(t, BuildFileReportIndex(file, 2, map[config.MetricKey]bool{}), "no active metrics, no index")
}

func TestBuildFileReportIndex_FallsBackWhenPerReportHitsAreMissing(t *testing.T) {
	file := &model.FileNode{
		Lines: map[int]model.LineMetrics{
			1: {Hits: 4}, // merged before per-report hits existed
			2: {Hits: 0},
		},
	}

	idx := BuildFileReportIndex(file, 3, allFileMetrics())
	covered, total := resolve(idx[config.LineCoverage], 0b001)
	assert.Equal(t, 1, covered, "an executed line without a breakdown counts for any selection")
	assert.Equal(t, 2, total)
}
