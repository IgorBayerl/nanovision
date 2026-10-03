package server_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/history"
	"github.com/IgorBayerl/nanovision/internal/reporter/htmlreact"
	"github.com/IgorBayerl/nanovision/internal/server"
	"github.com/IgorBayerl/nanovision/internal/store"
	"github.com/IgorBayerl/nanovision/internal/store/blob"
	"github.com/stretchr/testify/require"
)

// TestScale measures a synthetic run at studio size, for milestone M2 of
// docs/plans/run-store-and-server.md. It is skipped unless NANOVISION_SCALE_FILES
// is set, e.g. to 30000. With NANOVISION_SCALE_STORE set, the store stays in
// that folder, so `nanovision serve` can open it in a browser.
func TestScale(t *testing.T) {
	files, _ := strconv.Atoi(os.Getenv("NANOVISION_SCALE_FILES"))
	if files <= 0 {
		t.Skip("set NANOVISION_SCALE_FILES to run the scale measurement")
	}
	dir := os.Getenv("NANOVISION_SCALE_STORE")
	if dir == "" {
		dir = t.TempDir()
	}
	ctx := context.Background()
	s, err := store.Open(dir)
	require.NoError(t, err)
	defer s.Close()
	rng := rand.New(rand.NewPCG(1, 2))

	// 150 coverable lines, 15 functions and 130 statements a file, on average
	start := time.Now()
	manifest, blobs, lines := synthetic(rng, files)
	encoded := time.Since(start)

	var sizes = map[blob.Kind]int{}
	for _, b := range blobs {
		sizes[blob.Kind(b.Data[0])] += len(b.Data)
	}
	manifestData := blob.EncodeManifest(manifest)
	blobs = append(blobs, store.Blob{Hash: blob.Sum(manifestData), Data: manifestData})

	start = time.Now()
	require.NoError(t, s.PutBlobs(ctx, blobs))
	totals, _ := manifest.Totals()
	cfg := config.GetDefaultConfig()
	cfg.FileMetrics = []config.MetricKey{config.StatementCoverage, config.MethodsHit, config.MaxCyclomaticComplexity}
	cfg.MethodMetrics = []config.MetricKey{config.MethodStatementCoverage, config.CyclomaticComplexity, config.MethodCrapScore}
	cfg.StatusBands = config.StatusBands{config.StatementCoverage: {Min: 60, Max: 80}}
	meta, _ := json.Marshal(history.Meta{Title: fmt.Sprintf("Synthetic %d files", files), Reports: []string{"coverage.xml"},
		Parsers: []string{"Cobertura"}, Files: files, Totals: totals, Config: history.SnapshotConfig(cfg)})
	id, err := s.AddRun(ctx, store.Run{Project: "scale", Stream: "//scale/main", Profile: "unit", Kind: store.KindSubmit,
		Clean: true, Revision: "//scale/main@1000", Change: 1000, Tool: "test", Manifest: blob.Sum(manifestData), Meta: meta})
	require.NoError(t, err)
	stored := time.Since(start)

	// a second run where the coverage of 1% of the files changed
	for i := range manifest.Entries[:files/100] {
		e := &manifest.Entries[i*100]
		cov := blob.EncodeCoverage(blob.Coverage{Reports: 1, Lines: []blob.Line{{Number: 1, Hits: i + 7}}})
		e.Coverage = blob.Sum(cov)
		require.NoError(t, s.PutBlobs(ctx, []store.Blob{{Hash: e.Coverage, Data: cov}}))
	}
	second := blob.EncodeManifest(manifest)
	require.NoError(t, s.PutBlobs(ctx, []store.Blob{{Hash: blob.Sum(second), Data: second}}))
	_, err = s.AddRun(ctx, store.Run{Project: "scale", Stream: "//scale/main", Profile: "unit", Kind: store.KindSubmit,
		Clean: true, Revision: "//scale/main@1001", Change: 1001, Tool: "test", Manifest: blob.Sum(second), Meta: meta,
		CreatedAt: time.Now().Add(time.Second)})
	require.NoError(t, err)

	ui, err := htmlreact.DistFS()
	require.NoError(t, err)
	srv, err := server.New(server.Options{Store: s, Token: "secret", Version: "test", Build: "scale", UI: ui, Logger: quiet})
	require.NoError(t, err)
	ts := httptest.NewServer(srv)
	defer ts.Close()

	fetch := func() (time.Duration, int, int) {
		req, _ := http.NewRequest(http.MethodGet, ts.URL+"/api/v1/runs/"+strconv.FormatInt(id, 10)+"/summary", nil)
		req.Header.Set("Accept-Encoding", "gzip")
		start := time.Now()
		resp, err := http.DefaultTransport.RoundTrip(req)
		require.NoError(t, err)
		gz, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		elapsed := time.Since(start)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		plain, _ := http.Get(ts.URL + "/api/v1/runs/" + strconv.FormatInt(id, 10) + "/summary")
		raw, _ := io.ReadAll(plain.Body)
		plain.Body.Close()
		return elapsed, len(gz), len(raw)
	}
	runtime.GC()
	var before runtime.MemStats
	runtime.ReadMemStats(&before)
	cold, gzSize, rawSize := fetch()
	var after runtime.MemStats
	runtime.ReadMemStats(&after)
	warm, _, _ := fetch()

	t.Logf("files %d, coverable lines %d", files, lines)
	t.Logf("encode %s, store %s", encoded.Round(time.Millisecond), stored.Round(time.Millisecond))
	t.Logf("blob bytes: coverage %d (%.1f B/line), analysis %d, manifest %d (%.0f B/file)",
		sizes[blob.KindCoverage], float64(sizes[blob.KindCoverage])/float64(lines), sizes[blob.KindAnalysis],
		len(manifestData), float64(len(manifestData))/float64(files))
	if info, err := os.Stat(s.Path()); err == nil {
		t.Logf("store file after a second run with 1%% of the coverage changed: %d KB", info.Size()>>10)
	}
	t.Logf("summary: cold %s, warm %s, %d KB gzip, %d KB raw, %d MB allocated by the cold build",
		cold.Round(time.Millisecond), warm.Round(time.Millisecond), gzSize>>10, rawSize>>10, (after.TotalAlloc-before.TotalAlloc)>>20)
}

func synthetic(rng *rand.Rand, files int) (blob.Manifest, []store.Blob, int) {
	var m blob.Manifest
	var blobs []store.Blob
	seen := map[blob.Hash]bool{}
	add := func(data []byte) blob.Hash {
		h := blob.Sum(data)
		if !seen[h] {
			seen[h] = true
			blobs = append(blobs, store.Blob{Hash: h, Data: data})
		}
		return h
	}
	lines := 0
	for f := range files {
		n := 50 + rng.IntN(200)
		cov := blob.Coverage{Reports: 1}
		var c blob.Counts
		line := 1
		for range n {
			line += 1 + rng.IntN(3)
			hits := 0
			if rng.Float64() < 0.75 {
				hits = 1 + rng.IntN(50)
			}
			cov.Lines = append(cov.Lines, blob.Line{Number: line, Hits: hits})
			c.LinesValid++
			if hits > 0 {
				c.LinesCovered++
			}
		}
		lines += n
		var a blob.Analysis
		a.TotalLines = line + 10
		for fn := range 5 + rng.IntN(20) {
			start := 1 + fn*10
			cc := 1 + rng.IntN(12)
			a.Functions = append(a.Functions, blob.Function{Name: fmt.Sprintf("Class%d::Method%d", f%50, fn), StartLine: start, EndLine: start + 8, Complexity: &cc})
			c.MaxComplexity = max(c.MaxComplexity, cc)
		}
		for st := range n * 9 / 10 {
			a.Statements = append(a.Statements, blob.Statement{StartLine: 2 + st*2, EndLine: 2 + st*2, Type: "expression_statement"})
		}
		c.StatementsValid = len(a.Statements)
		c.StatementsCovered = c.StatementsValid * c.LinesCovered / max(c.LinesValid, 1)
		c.MethodsValid = len(a.Functions)
		c.MethodsHit = c.MethodsValid * 3 / 4
		m.Entries = append(m.Entries, blob.Entry{
			Path:     fmt.Sprintf("src/module%03d/sub%d/file_%05d.cpp", f%200, f%7, f),
			Coverage: add(blob.EncodeCoverage(cov)),
			Analysis: add(blob.EncodeAnalysis(a)),
			Content:  blob.Sum([]byte(strconv.Itoa(f))),
			Counts:   c,
		})
	}
	return m, blobs, lines
}
