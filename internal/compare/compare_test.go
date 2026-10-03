package compare

import (
	"bytes"
	"strings"
	"testing"

	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/reporter/textsummary"
	"github.com/IgorBayerl/nanovision/internal/store/blob"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func entry(path string, stmtCovered, stmtValid int) blob.Entry {
	return blob.Entry{Path: path, Counts: blob.Counts{
		StatementsCovered: stmtCovered, StatementsValid: stmtValid,
		LinesCovered: stmtCovered, LinesValid: stmtValid,
	}}
}

func manifest(entries ...blob.Entry) blob.Manifest { return blob.Manifest{Entries: entries} }

var metrics = []config.MetricKey{config.StatementCoverage, config.LineCoverage, config.MethodsHit, config.MaxCyclomaticComplexity}

func TestCompareCodeChange(t *testing.T) {
	base := manifest(
		entry("src/core/inventory.cpp", 88, 100),
		entry("src/net/replication.cpp", 61, 100),
		entry("src/old.cpp", 5, 10),
	)
	changed := entry("src/net/replication.cpp", 80, 110)
	changed.Diff = blob.DiffModified
	changed.Patch = &blob.PatchCounts{StatementsCovered: 8, StatementsValid: 10}
	current := manifest(
		entry("src/core/inventory.cpp", 70, 100),
		changed,
		entry("src/new.cpp", 1, 1),
	)

	cmp := Compare(base, current, Options{FileMetrics: metrics, HasDiff: true,
		Change: &model.ChangeSet{Measured: []string{"src/net/replication.cpp"}}})

	assert.Equal(t, config.StatementCoverage, cmp.Headline)
	require.Len(t, cmp.Metrics, 2, "methods hit has no data and complexity is not coverage")
	head := cmp.Metrics[0]
	assert.Equal(t, 154, head.BaseCovered)
	assert.Equal(t, 210, head.BaseTotal)
	assert.Equal(t, 151, head.CurrentCovered)
	assert.Equal(t, 211, head.CurrentTotal)
	assert.Equal(t, 73.33, head.Base)
	assert.Equal(t, 71.56, head.Current)
	assert.Equal(t, -1.77, head.Delta)
	assert.Equal(t, "Line Coverage", cmp.Metrics[1].Label)

	// every file and folder whose coverage moved carries the delta of each metric
	assert.Equal(t, -18.0, cmp.Deltas["src/core/inventory.cpp"][config.StatementCoverage])
	assert.Equal(t, -18.0, cmp.Deltas["src/core"][config.StatementCoverage], "the folder holds only this file")
	assert.Equal(t, 11.72, cmp.Deltas["src/net/replication.cpp"][config.LineCoverage])
	assert.InDelta(t, -1.77, cmp.Deltas["src"][config.StatementCoverage], 0.001, "src/old.cpp and src/new.cpp move the folder")
	assert.NotContains(t, cmp.Deltas, "src/new.cpp", "a file the base run does not have has no delta")
	assert.NotContains(t, cmp.Deltas, "src/old.cpp")

	require.Len(t, cmp.Files, 2)
	assert.Equal(t, "src/core/inventory.cpp", cmp.Files[0].Path, "largest change first")
	assert.Equal(t, -18.0, cmp.Files[0].Delta)
	assert.Equal(t, "not in diff", cmp.Files[0].Change)
	assert.Equal(t, "modified", cmp.Files[1].Change)
	assert.Equal(t, 1, cmp.AddedFiles)
	assert.Equal(t, 1, cmp.RemovedFiles)
	require.NotNil(t, cmp.Patch)
	assert.Equal(t, 8, cmp.Patch.Covered)
	assert.Empty(t, cmp.Warnings)
}

func TestCompareTestOnlyChange(t *testing.T) {
	base := manifest(entry("src/inventory.cpp", 61, 100), entry("src/item.cpp", 70, 100))
	current := manifest(entry("src/inventory.cpp", 88, 100), entry("src/item.cpp", 70, 100))
	change := &model.ChangeSet{Ignored: []string{"tests/inventory_test.cpp", "tests/helpers.cpp"}}

	cmp := Compare(base, current, Options{FileMetrics: metrics, HasDiff: true, Change: change})
	assert.True(t, cmp.Change.TestOnly())
	assert.Equal(t, 13.5, cmp.Metrics[0].Delta)
	assert.Equal(t, 0, cmp.Metrics[0].CurrentTotal-cmp.Metrics[0].BaseTotal, "no product code was added")
	require.Len(t, cmp.Files, 1)
	assert.Equal(t, "not in diff", cmp.Files[0].Change)
	require.NotNil(t, cmp.Patch)
	assert.Zero(t, cmp.Patch.Total)

	var out bytes.Buffer
	cmp.Base = model.RunRef{Revision: "//game/main@118432", Profile: "unit-win64"}
	cmp.Exact = true
	textsummary.WriteComparison(&out, cmp)
	text := out.String()
	assert.Contains(t, text, "//game/main@118432  (exact)")
	assert.Contains(t, text, "2 files, all ignored by the coverage filters (tests)")
	assert.Contains(t, text, "none: the change has no measured lines")
	assert.Contains(t, text, "+13.50")
	assert.Contains(t, text, "src/inventory.cpp")
	for _, r := range text {
		if r > 127 {
			t.Fatalf("the block must be ASCII, found %q in:\n%s", r, text)
		}
	}
}

func TestCompareAppliesCurrentFiltersToBase(t *testing.T) {
	base := manifest(entry("src/a.go", 5, 10), entry("vendor/lib.go", 0, 90))
	current := manifest(entry("src/a.go", 5, 10))

	cmp := Compare(base, current, Options{Include: func(p string) bool { return !strings.HasPrefix(p, "vendor/") }})
	assert.Equal(t, 1, cmp.FilteredBaseFiles)
	assert.Equal(t, 0, cmp.RemovedFiles)
	assert.Equal(t, 0.0, cmp.Metrics[0].Delta)
	assert.Empty(t, cmp.Warnings, "filtered files do not count as missing")
}

func TestCompareHeadlineFallsBackToLines(t *testing.T) {
	base := manifest(blob.Entry{Path: "a.cs", Counts: blob.Counts{LinesCovered: 1, LinesValid: 4}})
	current := manifest(blob.Entry{Path: "a.cs", Counts: blob.Counts{LinesCovered: 2, LinesValid: 4, StatementsCovered: 1, StatementsValid: 1}})

	cmp := Compare(base, current, Options{})
	assert.Equal(t, config.LineCoverage, cmp.Headline, "statements need both runs")
	assert.Equal(t, 25.0, cmp.Metrics[0].Delta)
}

func TestCompareZeroTotals(t *testing.T) {
	cmp := Compare(manifest(), manifest(), Options{FileMetrics: metrics})
	require.Len(t, cmp.Metrics, 1)
	assert.Equal(t, 0.0, cmp.Metrics[0].Base)
	assert.Equal(t, 0.0, cmp.Metrics[0].Delta)
}

func TestFairnessWarnings(t *testing.T) {
	base := manifest(entry("a.go", 10, 10), entry("b.go", 10, 30))
	current := manifest(entry("a.go", 10, 10))

	cmp := Compare(base, current, Options{
		BaseReports: []string{"unit.out", "integration.out"}, CurrentReports: []string{"unit.out"},
		BaseTool: "v1.3.0", CurrentTool: "v1.4.1",
	})
	require.Len(t, cmp.Warnings, 3)
	assert.Contains(t, cmp.Warnings[0], "report patterns differ")
	assert.Contains(t, cmp.Warnings[1], "75% of the base run's coverable lines")
	assert.Contains(t, cmp.Warnings[2], "v1.3.0")

	same := Compare(base, base, Options{BaseReports: []string{"b", "a"}, CurrentReports: []string{"a", "b"}, BaseTool: "dev", CurrentTool: "v1.4.1"})
	assert.Empty(t, same.Warnings, "same patterns in another order, and a dev build, are fine")
}
