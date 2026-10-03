package blob

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func intPtr(v int) *int { return &v }

func sampleCoverage() Coverage {
	return Coverage{
		Reports: 3,
		Lines: []Line{
			{Number: 12, Hits: 4, ReportHits: []int{4, 0, 0}},
			{Number: 3, Hits: 0, ReportHits: []int{0, 0, 0}},
			{Number: 4, Hits: -1, ReportHits: []int{-1, 0, 0}},
			{Number: 5, Hits: 7, ReportHits: []int{2, 5, 0}},
			{Number: 900, Hits: 1 << 40, ReportHits: []int{0, 0, 1 << 40}},
		},
	}
}

func TestCoverageRoundTrip(t *testing.T) {
	in := sampleCoverage()
	data := EncodeCoverage(in)

	out, err := DecodeCoverage(data)
	require.NoError(t, err)
	assert.Equal(t, 3, out.Reports)
	require.Len(t, out.Lines, 5)
	assert.Equal(t, []int{3, 4, 5, 12, 900}, []int{out.Lines[0].Number, out.Lines[1].Number, out.Lines[2].Number, out.Lines[3].Number, out.Lines[4].Number})
	assert.Equal(t, Line{Number: 5, Hits: 7, ReportHits: []int{2, 5, 0}}, out.Lines[2])
	assert.Equal(t, Line{Number: 4, Hits: -1, ReportHits: []int{-1, 0, 0}}, out.Lines[1])
	assert.Equal(t, 1<<40, out.Lines[4].Hits)

	assert.Equal(t, data, EncodeCoverage(out), "encoding must be deterministic")
}

func TestCoverageSingleReportStoresNoReportHits(t *testing.T) {
	lines := make([]Line, 0, 100)
	for i := 1; i <= 100; i++ {
		lines = append(lines, Line{Number: i, Hits: i % 3, ReportHits: []int{i % 3}})
	}
	data := EncodeCoverage(Coverage{Reports: 1, Lines: lines})

	// two bytes per line: the line gap and the hits
	assert.LessOrEqual(t, len(data), 2*100+8)

	out, err := DecodeCoverage(data)
	require.NoError(t, err)
	assert.Equal(t, lines, out.Lines)
}

func TestCoverageNormalizesReportHits(t *testing.T) {
	data := EncodeCoverage(Coverage{Reports: 2, Lines: []Line{{Number: 1, Hits: 3, ReportHits: []int{3}}}})
	out, err := DecodeCoverage(data)
	require.NoError(t, err)
	assert.Equal(t, []int{3, 0}, out.Lines[0].ReportHits, "missing report entries are zero")
}

func TestCoverageRejectsCorruptInput(t *testing.T) {
	data := EncodeCoverage(sampleCoverage())

	for name, bad := range map[string][]byte{
		"empty":          {},
		"wrong kind":     append([]byte{byte(KindAnalysis)}, data[1:]...),
		"future version": append([]byte{data[0], 99}, data[2:]...),
		"truncated":      data[:len(data)-1],
		"trailing bytes": append(bytes.Clone(data), 0),
		"huge line list": {byte(KindCoverage), 1, 1, 0xff, 0xff, 0xff, 0x7f},
	} {
		_, err := DecodeCoverage(bad)
		assert.Truef(t, errors.Is(err, ErrCorrupt), "%s: got %v", name, err)
	}
}

func TestAnalysisRoundTrip(t *testing.T) {
	in := Analysis{
		TotalLines: 120,
		Functions: []Function{
			{Name: "(*Store).Put", StartLine: 10, EndLine: 40, Complexity: intPtr(7)},
			{Name: "helper", StartLine: 50, EndLine: 52},
			{Name: "zero", StartLine: 60, EndLine: 60, Complexity: intPtr(0)},
		},
		Statements: []Statement{
			{StartLine: 11, EndLine: 11, Type: "expression_statement"},
			{StartLine: 12, EndLine: 15, Type: "if_statement"},
			{StartLine: 9, EndLine: 9, Type: "expression_statement"},
			{StartLine: 51, EndLine: 51, Type: "return_statement"},
		},
	}
	data := EncodeAnalysis(in)
	out, err := DecodeAnalysis(data)
	require.NoError(t, err)
	assert.Equal(t, in, out)
	assert.Equal(t, data, EncodeAnalysis(out))
}

func TestAnalysisEmpty(t *testing.T) {
	out, err := DecodeAnalysis(EncodeAnalysis(Analysis{TotalLines: 3}))
	require.NoError(t, err)
	assert.Equal(t, Analysis{TotalLines: 3}, out)
}

func sampleManifest() Manifest {
	return Manifest{Entries: []Entry{
		{
			Path:     "src/net/replication.cpp",
			Content:  Sum([]byte("replication")),
			Analysis: Sum([]byte("a1")),
			Coverage: Sum([]byte("c1")),
			Diff:     DiffModified,
			Counts:   Counts{TotalLines: 300, LinesCovered: 80, LinesValid: 120, StatementsCovered: 70, StatementsValid: 100, MethodsHit: 4, MethodsValid: 6, MaxComplexity: 12},
			Patch:    &PatchCounts{LinesTotal: 9, LinesCovered: 3, LinesValid: 5, StatementsCovered: 2, StatementsValid: 4, MethodsHit: 1, MethodsValid: 2},
		},
		{
			Path:     "src/core/inventory.cpp",
			Content:  Sum([]byte("inventory")),
			Coverage: Sum([]byte("c2")),
			Counts:   Counts{LinesCovered: 1, LinesValid: 2},
		},
		{Path: "src/core/item.cpp", Coverage: Sum([]byte("c3"))},
	}}
}

func TestManifestRoundTripSortsEntries(t *testing.T) {
	data := EncodeManifest(sampleManifest())
	out, err := DecodeManifest(data)
	require.NoError(t, err)

	require.Len(t, out.Entries, 3)
	assert.Equal(t, "src/core/inventory.cpp", out.Entries[0].Path)
	assert.Equal(t, "src/core/item.cpp", out.Entries[1].Path)
	assert.Equal(t, "src/net/replication.cpp", out.Entries[2].Path)
	assert.True(t, out.Entries[1].Content.IsZero())
	assert.Equal(t, DiffModified, out.Entries[2].Diff)
	require.NotNil(t, out.Entries[2].Patch)
	assert.Equal(t, 4, out.Entries[2].Patch.StatementsValid)
	assert.Equal(t, 12, out.Entries[2].Counts.MaxComplexity)

	assert.Equal(t, data, EncodeManifest(out))
	assert.NotNil(t, out.Find("src/core/item.cpp"))
	assert.Nil(t, out.Find("src/core/missing.cpp"))
}

func TestManifestTotals(t *testing.T) {
	counts, patch := sampleManifest().Totals()
	assert.Equal(t, 81, counts.LinesCovered)
	assert.Equal(t, 122, counts.LinesValid)
	assert.Equal(t, 12, counts.MaxComplexity)
	assert.Equal(t, 4, patch.StatementsValid)
}

func TestManifestRejectsUnsortedPaths(t *testing.T) {
	e := newEncoder(KindManifest, manifestVersion, 0)
	e.count(2)
	encodeEntry(e, "", &Entry{Path: "b"})
	encodeEntry(e, "b", &Entry{Path: "a"})
	_, err := DecodeManifest(e.buf)
	assert.ErrorIs(t, err, ErrCorrupt)
}

func TestDiffRoundTrip(t *testing.T) {
	in := Diff{Files: []DiffFile{
		{Path: "b.go", Kind: DiffModified, Added: []int{10, 11, 12, 20, 3, 3}, Modified: []int{40}},
		{Path: "a.go", Kind: DiffAdded, Added: []int{1, 2, 3, 4}},
	}}
	data := EncodeDiff(in)
	out, err := DecodeDiff(data)
	require.NoError(t, err)
	require.Len(t, out.Files, 2)
	assert.Equal(t, "a.go", out.Files[0].Path)
	assert.Equal(t, []int{1, 2, 3, 4}, out.Files[0].Added)
	assert.Equal(t, []int{3, 10, 11, 12, 20}, out.Files[1].Added)
	assert.Equal(t, []int{40}, out.Files[1].Modified)
	assert.Equal(t, data, EncodeDiff(out))
	assert.Equal(t, DiffAdded, out.Find("a.go").Kind)
}

func TestHashText(t *testing.T) {
	h := Sum([]byte("x"))
	parsed, err := ParseHash(h.String())
	require.NoError(t, err)
	assert.Equal(t, h, parsed)

	_, err = ParseHash("abc")
	assert.Error(t, err)
}
