package history

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/diff"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/pipeline"
	"github.com/IgorBayerl/nanovision/internal/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, _ := runtime.Caller(0)
	return filepath.Dir(filepath.Dir(filepath.Dir(file)))
}

// demoTree runs the real pipeline on the demo projects: Go, and C++ measured
// by two reports, so lines carry per-report hits.
func demoTree(t *testing.T, withDiff bool) (*model.SummaryTree, *config.AppConfig) {
	t.Helper()
	root := repoRoot(t)
	demo := func(parts ...string) string {
		return filepath.Join(append([]string{root, "demo_projects"}, parts...)...)
	}

	// the test runs in this package's folder, which has no nanovision.yaml to pick up
	cfg, err := config.Load("", config.RawConfigInput{
		Reports: []string{demo("go", "report", "gocover", "coverage.out") + ";" +
			demo("cpp", "report", "gcov", "branch-probabilities", "*.gcov") + ";" +
			demo("cpp", "report", "cobertura", "cobertura.xml")},
		SourceDirs:  demo("go", "project") + ";" + demo("cpp", "project") + ";" + demo("cpp", "project"),
		ReportTypes: "TextSummary,Html", OutputDir: "coverage-report", LogFormat: "text", Verbosity: "Info",
		IgnoreCache: true,
		StatusBands: []string{"statement_coverage=60..80", "methods.complexity=0..5"},
	})
	require.NoError(t, err)
	cfg.ProjectRoot = root

	var dd *diff.DiffData
	if withDiff {
		dd = &diff.DiffData{Files: []diff.FileDiff{
			{NewPath: "demo_projects/go/project/calculator/calculator.go", Kind: "modified",
				Hunks: []diff.Hunk{{NewStart: 4, AddedLineOffsets: []int{0, 1, 2, 11, 12, 13}}}},
			{NewPath: "demo_projects/cpp/project/src/advanced_calculator.cpp", Kind: "added"},
			{NewPath: "demo_projects/go/project/calculator/calculator_test.go", Kind: "modified"},
		}}
	}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	tree, err := pipeline.Run(cfg, dd, pipeline.BuildInfo{Version: "test", Commit: "test"}, logger)
	require.NoError(t, err)
	return tree, cfg
}

// saveCapture stores a capture, the way the CLI does.
func saveCapture(t *testing.T, s *store.Store, c *Capture) int64 {
	t.Helper()
	ctx := context.Background()
	require.NoError(t, s.PutBlobs(ctx, c.Blobs))
	meta, err := json.Marshal(c.Meta)
	require.NoError(t, err)
	id, err := s.AddRun(ctx, store.Run{Project: "demo", Stream: "main", Profile: "default", Kind: store.KindLocal, Tool: "test", Manifest: c.ManifestHash, Meta: meta})
	require.NoError(t, err)
	return id
}

type fileView struct {
	Metrics    string
	Methods    string
	Statuses   string
	Lines      string
	Statements string
	Diff       string
	TotalLines int
}

func mustJSON(t *testing.T, v any) string {
	b, err := json.Marshal(v)
	require.NoError(t, err)
	return string(b)
}

// views flattens a tree to what the reporters read, so two trees compare field by field.
func views(t *testing.T, tree *model.SummaryTree) (map[string]fileView, map[string]string) {
	files := map[string]fileView{}
	dirs := map[string]string{}
	var walk func(d *model.DirNode)
	walk = func(d *model.DirNode) {
		dirs[d.Path] = mustJSON(t, d.Metrics) + mustJSON(t, d.Statuses)
		for _, f := range d.Files {
			statements := ""
			if len(f.Statements) > 0 {
				statements = mustJSON(t, f.Statements)
			}
			files[f.Path] = fileView{
				Metrics:    mustJSON(t, f.Metrics),
				Methods:    mustJSON(t, f.Methods),
				Statuses:   mustJSON(t, f.Statuses),
				Lines:      mustJSON(t, f.Lines),
				Statements: statements,
				Diff:       mustJSON(t, f.Diff),
				TotalLines: f.TotalLines,
			}
		}
		for _, sub := range d.Subdirs {
			walk(sub)
		}
	}
	walk(tree.Root)
	return files, dirs
}

func TestRebuildMatchesThePipeline(t *testing.T) {
	for _, withDiff := range []bool{false, true} {
		name := "without diff"
		if withDiff {
			name = "with diff"
		}
		t.Run(name, func(t *testing.T) {
			original, cfg := demoTree(t, withDiff)
			capture, err := NewCapture(original, cfg, "")
			require.NoError(t, err)
			require.NotEmpty(t, capture.Manifest.Entries)
			assert.Len(t, capture.Meta.Reports, 3)

			s, err := store.Open(t.TempDir())
			require.NoError(t, err)
			defer s.Close()
			id := saveCapture(t, s, capture)

			run, err := s.Run(context.Background(), id)
			require.NoError(t, err)
			meta, err := ParseMeta(run.Meta)
			require.NoError(t, err)
			manifest, err := s.Manifest(context.Background(), run.Manifest)
			require.NoError(t, err)

			rebuilt, _, err := Rebuild(context.Background(), s, meta, manifest, nil)
			require.NoError(t, err)

			wantFiles, wantDirs := views(t, original)
			gotFiles, gotDirs := views(t, rebuilt)
			require.GreaterOrEqual(t, len(wantFiles), 6)
			require.Equal(t, len(wantFiles), len(gotFiles))
			assertExercised(t, original)
			for path, want := range wantFiles {
				assert.Equalf(t, want, gotFiles[path], "file %s", path)
			}
			assert.Equal(t, wantDirs, gotDirs)
			assert.Equal(t, mustJSON(t, original.Metrics), mustJSON(t, rebuilt.Metrics))

			if withDiff {
				require.NotNil(t, meta.Diff)
				require.NotNil(t, meta.Change)
				assert.Equal(t, []string{"demo_projects/go/project/calculator/calculator_test.go"}, meta.Change.NotInReports)
				assert.Positive(t, meta.Patch.LinesValid)
			} else {
				assert.Nil(t, meta.Diff)
			}
		})
	}
}

func TestRebuildOneFile(t *testing.T) {
	original, cfg := demoTree(t, true)
	capture, err := NewCapture(original, cfg, "")
	require.NoError(t, err)
	s, err := store.Open(t.TempDir())
	require.NoError(t, err)
	defer s.Close()
	saveCapture(t, s, capture)

	const path = "demo_projects/go/project/calculator/calculator.go"
	tree, _, err := Rebuild(context.Background(), s, capture.Meta, capture.Manifest, func(p string) bool { return p == path })
	require.NoError(t, err)

	files, _ := views(t, tree)
	require.Len(t, files, 1)
	want, _ := views(t, original)
	assert.Equal(t, want[path].Methods, files[path].Methods)
	assert.Equal(t, want[path].Lines, files[path].Lines)
}

// assertExercised guards the round trip test against passing on empty data.
func assertExercised(t *testing.T, tree *model.SummaryTree) {
	t.Helper()
	var multiReportLines, methodsWithStatus, complexities int
	var walk func(d *model.DirNode)
	walk = func(d *model.DirNode) {
		for _, f := range d.Files {
			for _, l := range f.Lines {
				nonZero := 0
				for _, h := range l.ReportHits {
					if h != 0 {
						nonZero++
					}
				}
				if nonZero > 1 {
					multiReportLines++
				}
			}
			for _, m := range f.Methods {
				if len(m.Statuses) > 0 {
					methodsWithStatus++
				}
				if m.CyclomaticComplexity != nil {
					complexities++
				}
			}
		}
		for _, sub := range d.Subdirs {
			walk(sub)
		}
	}
	walk(tree.Root)
	assert.Positive(t, multiReportLines, "lines hit by more than one report")
	assert.Positive(t, methodsWithStatus)
	assert.Positive(t, complexities)
}
