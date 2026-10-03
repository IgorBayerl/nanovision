package server_test

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/IgorBayerl/nanovision/internal/calculator"
	"github.com/IgorBayerl/nanovision/internal/client"
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/diff"
	"github.com/IgorBayerl/nanovision/internal/history"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/pipeline"
	"github.com/IgorBayerl/nanovision/internal/reporter/htmlreact"
	"github.com/IgorBayerl/nanovision/internal/server"
	"github.com/IgorBayerl/nanovision/internal/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var quiet = slog.New(slog.NewTextHandler(io.Discard, nil))

func demoRun(t *testing.T, withDiff bool) (*model.SummaryTree, *config.AppConfig) {
	t.Helper()
	var files, methods []config.MetricKey
	for k := range calculator.FileRegistry {
		files = append(files, k)
	}
	for k := range calculator.MethodRegistry {
		methods = append(methods, k)
	}
	sort.Slice(files, func(i, j int) bool { return files[i] < files[j] })
	sort.Slice(methods, func(i, j int) bool { return methods[i] < methods[j] })
	config.RegisterDefaultMetrics(files, methods)

	_, here, _, _ := runtime.Caller(0)
	root := filepath.Dir(filepath.Dir(filepath.Dir(here)))
	demo := func(parts ...string) string {
		return filepath.Join(append([]string{root, "demo_projects"}, parts...)...)
	}
	cfg, err := config.Load("", config.RawConfigInput{
		ReportPatterns: demo("go", "report", "gocover", "coverage.out") + ";" + demo("cpp", "report", "cobertura", "cobertura.xml"),
		SourceDirs:     demo("go", "project") + ";" + demo("cpp", "project"),
		ReportTypes:    "TextSummary,Html", OutputDir: "coverage-report", LogFormat: "text", Verbosity: "Info",
		IgnoreCache: true, StatusBands: []string{"statement_coverage=60..80"},
	})
	require.NoError(t, err)
	cfg.ProjectRoot = root

	var dd *diff.DiffData
	if withDiff {
		dd = &diff.DiffData{Files: []diff.FileDiff{{NewPath: "demo_projects/go/project/calculator/calculator.go", Kind: "modified",
			Hunks: []diff.Hunk{{NewStart: 4, AddedLineOffsets: []int{0, 1, 2, 11, 12}}}}}}
	}
	tree, err := pipeline.Run(cfg, dd, pipeline.BuildInfo{Version: "test", Commit: "test"}, quiet)
	require.NoError(t, err)
	return tree, cfg
}

func newServer(t *testing.T, token string) (*httptest.Server, *store.Store) {
	t.Helper()
	s, err := store.Open(t.TempDir())
	require.NoError(t, err)
	t.Cleanup(func() { s.Close() })
	ui, err := htmlreact.DistFS()
	require.NoError(t, err)
	srv, err := server.New(server.Options{Store: s, Token: token, Version: "test", Build: "test-build", UI: ui, Logger: quiet})
	require.NoError(t, err)
	ts := httptest.NewServer(srv)
	t.Cleanup(ts.Close)
	return ts, s
}

func getJSON(t *testing.T, url string, v any) int {
	t.Helper()
	resp, err := http.Get(url)
	require.NoError(t, err)
	defer resp.Body.Close()
	if v != nil && resp.StatusCode == http.StatusOK {
		require.NoError(t, json.NewDecoder(resp.Body).Decode(v))
	}
	return resp.StatusCode
}

// upload saves a run through the HTTP client, as CI does. baseID is the base
// run the CLI found for it, 0 for none.
func upload(t *testing.T, ts *httptest.Server, token string, tree *model.SummaryTree, cfg *config.AppConfig, run store.Run, baseID int64) int64 {
	t.Helper()
	capture, err := history.NewCapture(tree, cfg, "")
	require.NoError(t, err)
	if baseID > 0 {
		capture.Meta.Base = &history.BaseRef{RunRef: model.RunRef{ID: baseID}, Exact: true}
	}
	c := client.New(ts.URL, token)
	require.NoError(t, c.Ping(context.Background()))
	id, err := history.Save(context.Background(), c, capture, run, quiet)
	require.NoError(t, err)
	assert.Equal(t, ts.URL+"/runs/"+itoa(id), c.RunURL(id))
	return id
}

func itoa(id int64) string { return strconv.FormatInt(id, 10) }

func TestUploadAndRead(t *testing.T) {
	ts, _ := newServer(t, "secret")
	tree, cfg := demoRun(t, false)
	base := upload(t, ts, "secret", tree, cfg, store.Run{Project: "demo", Stream: "//demo/main", Profile: "unit",
		Kind: store.KindSubmit, Clean: true, Revision: "//demo/main@100", Change: 100, Tool: "test", CreatedAt: time.Now().Add(-time.Hour)}, 0)

	changed, cfg2 := demoRun(t, true)
	head := upload(t, ts, "secret", changed, cfg2, store.Run{Project: "demo", Stream: "//demo/main", Profile: "unit",
		Kind: store.KindReview, Revision: "//demo/main@100", Change: 100, ReviewID: "105", Author: "ada", Tool: "test"}, base)

	var bases []map[string]any
	baseURL := ts.URL + "/api/v1/base?project=demo&profile=unit&kind=submit&stream=" + url.QueryEscape("//demo/main")
	require.Equal(t, http.StatusOK, getJSON(t, baseURL+"&change=100", &bases))
	require.Len(t, bases, 1, "the review run has local edits and is no base run")
	assert.Equal(t, float64(base), bases[0]["id"])
	assert.NotNil(t, bases[0]["meta"], "the CLI reads the report patterns of the base run")
	require.Equal(t, http.StatusOK, getJSON(t, baseURL+"&revision="+url.QueryEscape("//demo/main@100"), &bases))
	require.Len(t, bases, 1)
	assert.Equal(t, http.StatusBadRequest, getJSON(t, baseURL, nil))

	// the server has no list of runs: its root opens the newest submit run
	noRedirect := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	resp, err := noRedirect.Get(ts.URL + "/")
	require.NoError(t, err)
	resp.Body.Close()
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, "/runs/"+itoa(base), resp.Header.Get("Location"))

	// the summary of a run carries the delta against its base run
	var summary map[string]any
	require.Equal(t, http.StatusOK, getJSON(t, ts.URL+"/api/v1/runs/"+itoa(head)+"/summary", &summary))
	cmp, ok := summary["comparison"].(map[string]any)
	require.True(t, ok, "the review run has a base run")
	assert.Equal(t, "statement_coverage", cmp["headline"])
	assert.NotEmpty(t, cmp["metrics"])
	assert.NotNil(t, cmp["patch"])
	assert.Equal(t, float64(base), cmp["base"].(map[string]any)["id"])
	assert.Equal(t, true, cmp["exact"])
	assert.NotNil(t, summary["review"], "a run with a diff carries the verdict on its changed code")

	var first map[string]any
	require.Equal(t, http.StatusOK, getJSON(t, ts.URL+"/api/v1/runs/"+itoa(base)+"/summary", &first))
	assert.Nil(t, first["comparison"], "the first run has no base run")
	assert.Nil(t, first["review"], "a run without a diff has no changed code")

	assert.Equal(t, http.StatusNotFound, getJSON(t, ts.URL+"/api/v1/runs/999/summary", nil))
	assert.Equal(t, http.StatusNotFound, getJSON(t, ts.URL+"/api/v1/runs/"+itoa(head)+"/file?path=nope.go", nil))
}

func TestServedSummaryMatchesTheStaticReport(t *testing.T) {
	ts, _ := newServer(t, "secret")
	tree, cfg := demoRun(t, true)
	id := upload(t, ts, "secret", tree, cfg, store.Run{Project: "demo", Stream: "main", Profile: "default", Kind: store.KindLocal, Tool: "test"}, 0)

	var served map[string]any
	require.Equal(t, http.StatusOK, getJSON(t, ts.URL+"/api/v1/runs/"+itoa(id)+"/summary", &served))

	staticJSON, err := htmlreact.SummaryJSON(tree, cfg, htmlreact.ViewOptions{})
	require.NoError(t, err)
	var static map[string]any
	require.NoError(t, json.Unmarshal(staticJSON, &static))

	for _, key := range []string{"totals", "metricDefinitions", "metricOrder", "reports", "reportIndexes", "statusBands"} {
		assert.Equalf(t, static[key], served[key], "summary field %s", key)
	}
	staticNodes, servedNodes := static["nodes"].([]any), served["nodes"].([]any)
	require.Equal(t, len(staticNodes), len(servedNodes))
	for i := range staticNodes {
		a, b := staticNodes[i].(map[string]any), servedNodes[i].(map[string]any)
		delete(a, "targetUrl")
		link, _ := b["targetUrl"].(string)
		delete(b, "targetUrl")
		assert.Equal(t, a, b)
		if b["type"] == "file" {
			assert.True(t, strings.HasPrefix(link, "/runs/"+itoa(id)+"/files/"), link)
		}
	}
	assert.NotNil(t, served["review"], "a run with a diff carries the review verdict")

	var details map[string]any
	path := "demo_projects/go/project/calculator/calculator.go"
	require.Equal(t, http.StatusOK, getJSON(t, ts.URL+"/api/v1/runs/"+itoa(id)+"/file?path="+url.QueryEscape(path), &details))
	assert.NotEmpty(t, details["lines"], "the coverage of each line, without the source text")
	assert.NotEmpty(t, details["methods"])
}

func TestServerCachesViews(t *testing.T) {
	ts, _ := newServer(t, "secret")
	tree, cfg := demoRun(t, false)
	id := upload(t, ts, "secret", tree, cfg, store.Run{Project: "demo", Profile: "default", Kind: store.KindLocal, Tool: "test"}, 0)

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/api/v1/runs/"+itoa(id)+"/summary?b=test-build", nil)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	resp.Body.Close()
	assert.Contains(t, resp.Header.Get("Cache-Control"), "immutable", "a URL with the server build never changes")
	etag := resp.Header.Get("ETag")
	require.NotEmpty(t, etag)

	req, _ = http.NewRequest(http.MethodGet, ts.URL+"/api/v1/runs/"+itoa(id)+"/summary", nil)
	req.Header.Set("If-None-Match", etag)
	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	resp.Body.Close()
	assert.Equal(t, http.StatusNotModified, resp.StatusCode)
}

func TestUploadsNeedTheToken(t *testing.T) {
	ts, _ := newServer(t, "secret")
	tree, cfg := demoRun(t, false)
	capture, err := history.NewCapture(tree, cfg, "")
	require.NoError(t, err)

	_, err = history.Save(context.Background(), client.New(ts.URL, "wrong"), capture, store.Run{Project: "demo", Profile: "default", Kind: store.KindLocal}, quiet)
	assert.ErrorContains(t, err, "upload token")
}

func TestUIIsServedForEveryPage(t *testing.T) {
	ts, _ := newServer(t, "secret")
	for _, path := range []string{"/", "/runs/12", "/runs/12/files/src/a.go"} {
		resp, err := http.Get(ts.URL + path)
		require.NoError(t, err)
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode, path)
		html := string(body)
		assert.Contains(t, html, `<base href="/">`, path)
		assert.Contains(t, html, `window.__NANOVISION_MODE__="server"`, path)
		assert.NotContains(t, html, "data.js", path)
	}

	resp, err := http.Get(ts.URL + "/api/v1/nothing")
	require.NoError(t, err)
	resp.Body.Close()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "unknown API paths are not the app shell")
}
