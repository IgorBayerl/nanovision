package server

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/IgorBayerl/nanovision/internal/compare"
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/history"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/reporter/htmlreact"
	"github.com/IgorBayerl/nanovision/internal/store"
	"github.com/IgorBayerl/nanovision/internal/store/blob"
)

// how long one view build may take before the request gives up
const buildTimeout = 2 * time.Minute

// handleBase lists the runs a delta may compare with, for the CLI's search
// for a base run. The answer carries the meta of each run.
func (s *Server) handleBase(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	q := store.CandidateQuery{Project: v.Get("project"), Profile: v.Get("profile"), Stream: v.Get("stream"), Revisions: v["revision"]}
	for _, k := range v["kind"] {
		kind, err := store.ParseRunKind(k)
		if err != nil {
			writeError(w, http.StatusBadRequest, "%v", err)
			return
		}
		q.Kinds = append(q.Kinds, kind)
	}
	if c := v.Get("change"); c != "" {
		n, err := strconv.ParseInt(c, 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "change: %q is not a number", c)
			return
		}
		q.Change = n
	}
	if len(q.Revisions) == 0 && q.Change <= 0 {
		writeError(w, http.StatusBadRequest, "base needs revision or change")
		return
	}
	runs, err := s.opts.Store.Candidates(r.Context(), q)
	if err != nil {
		writeError(w, statusFor(err), "%v", err)
		return
	}
	if runs == nil {
		runs = []store.Run{}
	}
	writeJSON(w, http.StatusOK, runs)
}

func (s *Server) handleSummary(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "%v", err)
		return
	}
	s.serveView(w, r, summaryKey(id), func(ctx context.Context) ([]byte, error) {
		return s.buildSummary(ctx, id)
	})
}

func (s *Server) handleFile(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "%v", err)
		return
	}
	path := r.URL.Query().Get("path")
	if path == "" {
		writeError(w, http.StatusBadRequest, "path is missing")
		return
	}
	s.serveView(w, r, fmt.Sprintf("file/%d/%s", id, path), func(ctx context.Context) ([]byte, error) {
		return s.buildFile(ctx, id, path)
	})
}

// serveView answers with a cached view. A saved run never changes, so a URL
// that names the server build is cached by the browser for good.
func (s *Server) serveView(w http.ResponseWriter, r *http.Request, key string, build func(context.Context) ([]byte, error)) {
	etag := `"` + s.opts.Build + "/" + key + `"`
	if match := r.Header.Get("If-None-Match"); match != "" && match == etag {
		w.WriteHeader(http.StatusNotModified)
		return
	}
	gz, err := s.views.get(key, func() ([]byte, error) {
		ctx, cancel := context.WithTimeout(context.Background(), buildTimeout)
		defer cancel()
		return build(ctx)
	})
	if err != nil {
		writeError(w, statusFor(err), "%v", err)
		return
	}

	h := w.Header()
	h.Set("Content-Type", "application/json; charset=utf-8")
	h.Set("ETag", etag)
	h.Set("Vary", "Accept-Encoding")
	if s.opts.Build != "" && r.URL.Query().Get("b") == s.opts.Build {
		h.Set("Cache-Control", "public, max-age=31536000, immutable")
	} else {
		h.Set("Cache-Control", "no-cache")
	}
	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		h.Set("Content-Encoding", "gzip")
		h.Set("Content-Length", strconv.Itoa(len(gz)))
		w.Write(gz)
		return
	}
	zr, err := gzip.NewReader(bytes.NewReader(gz))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "%v", err)
		return
	}
	io.Copy(w, zr)
}

// load reads a run with its meta and manifest.
func (s *Server) load(ctx context.Context, id int64) (store.Run, history.Meta, blob.Manifest, error) {
	run, err := s.opts.Store.Run(ctx, id)
	if err != nil {
		return run, history.Meta{}, blob.Manifest{}, err
	}
	meta, err := history.ParseMeta(run.Meta)
	if err != nil {
		return run, meta, blob.Manifest{}, err
	}
	manifest, err := s.opts.Store.Manifest(ctx, run.Manifest)
	return run, meta, manifest, err
}

func (s *Server) buildSummary(ctx context.Context, id int64) ([]byte, error) {
	run, meta, manifest, err := s.load(ctx, id)
	if err != nil {
		return nil, err
	}
	select {
	case s.heavy <- struct{}{}:
		defer func() { <-s.heavy }()
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	tree, cfg, err := history.Rebuild(ctx, s.opts.Store, meta, manifest, nil)
	if err != nil {
		return nil, err
	}
	tree.Comparison = s.comparison(ctx, run, meta, manifest, cfg)
	return htmlreact.SummaryJSON(tree, cfg, htmlreact.ViewOptions{
		FileURL:     func(path string) string { return FileURL(id, path) },
		Metadata:    runMetadata(run, meta),
		GeneratedAt: run.CreatedAt,
	})
}

// comparison is the delta of a run against the base run found when it was
// saved. It is nil when the run had no base run, or the base run is gone.
func (s *Server) comparison(ctx context.Context, run store.Run, meta history.Meta, manifest blob.Manifest, cfg *config.AppConfig) *model.Comparison {
	if meta.Base == nil {
		return nil
	}
	baseRun, baseMeta, baseManifest, err := s.load(ctx, meta.Base.ID)
	if err != nil {
		s.opts.Logger.Debug("No delta: the base run cannot be read", "run", run.ID, "base", meta.Base.ID, "error", err)
		return nil
	}
	cmp := compare.Compare(baseManifest, manifest, compare.Options{
		Include:        cfg.FileFilterInstance.IsElementIncludedInReport,
		FileMetrics:    cfg.FileMetrics,
		HasDiff:        meta.Diff != nil,
		Change:         meta.Change,
		BaseReports:    baseMeta.Reports,
		CurrentReports: meta.Reports,
		BaseTool:       baseRun.Tool,
		CurrentTool:    run.Tool,
	})
	cmp.Base, cmp.Current = history.Ref(baseRun), history.Ref(run)
	cmp.Exact, cmp.Distance = meta.Base.Exact, meta.Base.Distance
	return cmp
}

func (s *Server) buildFile(ctx context.Context, id int64, path string) ([]byte, error) {
	run, meta, manifest, err := s.load(ctx, id)
	if err != nil {
		return nil, err
	}
	entry := manifest.Find(path)
	if entry == nil {
		return nil, fmt.Errorf("run %d has no file %q: %w", id, path, store.ErrNotFound)
	}
	tree, cfg, err := history.Rebuild(ctx, s.opts.Store, meta, manifest, func(p string) bool { return p == path })
	if err != nil {
		return nil, err
	}
	file := findFile(tree.Root, path)
	if file == nil {
		return nil, fmt.Errorf("run %d has no file %q: %w", id, path, store.ErrNotFound)
	}

	// ponytail: no source text until the server reads it from the version
	// control system; the page shows the coverage of each line without it
	metadata := []htmlreact.MetadataItem{
		{Label: "Run", Value: runLabel(run)},
		{Label: "Source", Value: "The store does not keep source text."},
	}
	return htmlreact.DetailsJSON(tree, file, nil, cfg, htmlreact.ViewOptions{Metadata: metadata, GeneratedAt: run.CreatedAt})
}

func (s *Server) handleBlob(w http.ResponseWriter, r *http.Request) {
	h, err := blob.ParseHash(r.PathValue("hash"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "%v", err)
		return
	}
	data, err := s.opts.Store.Blob(r.Context(), h)
	if err != nil {
		writeError(w, statusFor(err), "%v", err)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	w.Write(data)
}

func summaryKey(id int64) string { return fmt.Sprintf("summary/%d", id) }

// FileURL is the page of one file of a run in the web UI.
func FileURL(id int64, path string) string {
	parts := strings.Split(path, "/")
	for i, p := range parts {
		parts[i] = url.PathEscape(p)
	}
	return fmt.Sprintf("/runs/%d/files/%s", id, strings.Join(parts, "/"))
}

func runLabel(run store.Run) string {
	if run.Revision != "" {
		return fmt.Sprintf("%d at %s", run.ID, run.Revision)
	}
	return strconv.FormatInt(run.ID, 10)
}

// runMetadata fills the information block of a run's summary page.
func runMetadata(run store.Run, meta history.Meta) []htmlreact.MetadataItem {
	var items []htmlreact.MetadataItem
	add := func(label string, value any) {
		switch v := value.(type) {
		case string:
			if v == "" {
				return
			}
		case []string:
			if len(v) == 0 {
				return
			}
		}
		items = append(items, htmlreact.MetadataItem{Label: label, Value: value})
	}
	add("Run", strconv.FormatInt(run.ID, 10))
	add("Revision", run.Revision)
	add("Stream", run.Stream)
	add("Profile", run.Profile)
	kind := string(run.Kind)
	if !run.Clean {
		kind += ", local edits"
	}
	add("Kind", kind)
	add("Review", run.ReviewID)
	add("Author", run.Author)
	add("Created", run.CreatedAt.Local().Format("2006-01-02 15:04:05"))
	add("CI build", run.CIURL)
	if meta.Base != nil {
		base := fmt.Sprintf("run %d", meta.Base.ID)
		if meta.Base.Revision != "" {
			base += " at " + meta.Base.Revision
		}
		add("Base run", base)
	}
	add("Parser", strings.Join(meta.Parsers, " | "))
	add("nanovision", run.Tool)
	if len(meta.Reports) > 0 {
		items = append(items, htmlreact.MetadataItem{Label: "Report Files", Value: meta.Reports, SizeHint: "large"})
	}
	return items
}

func findFile(dir *model.DirNode, path string) *model.FileNode {
	for _, f := range dir.Files {
		if f.Path == path {
			return f
		}
	}
	for _, sub := range dir.Subdirs {
		if f := findFile(sub, path); f != nil {
			return f
		}
	}
	return nil
}
