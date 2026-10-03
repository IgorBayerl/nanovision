package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// serverIndex turns the report's index.html into the app shell: every UI path
// serves it, a base element keeps the relative asset links working below
// /runs/..., and the mode flag makes the app read the API instead of data.js.
func serverIndex(ui fs.FS, build string) ([]byte, error) {
	raw, err := fs.ReadFile(ui, "index.html")
	if err != nil {
		return nil, fmt.Errorf("web UI: %w", err)
	}
	buildJSON, _ := json.Marshal(build)
	inject := `<base href="/"><script>window.__NANOVISION_MODE__="server";window.__NANOVISION_BUILD__=` +
		strings.ReplaceAll(string(buildJSON), "<", `<`) + `;</script>`

	html := string(raw)
	if i := strings.Index(html, "<head>"); i >= 0 {
		html = html[:i+len("<head>")] + inject + html[i+len("<head>"):]
	} else {
		return nil, fmt.Errorf("web UI: index.html has no <head>")
	}
	html = strings.Replace(html, `<script src="./data.js"></script>`, "", 1)
	html = devComment.ReplaceAllString(html, "")
	html = strings.Replace(html, "<title>Coverage Report</title>", "<title>nanovision</title>", 1)
	return []byte(html), nil
}

// assets never change while the server runs
var startTime = time.Now()

// the comment next to the data.js tag of the report's index.html
var devComment = regexp.MustCompile(`<!--[^>]*data\.js[^>]*-->`)

func (s *Server) handleUI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Allow", "GET, HEAD")
		writeError(w, http.StatusMethodNotAllowed, "the web UI only answers GET")
		return
	}
	if s.opts.UI == nil {
		writeError(w, http.StatusNotFound, "this server has no web UI")
		return
	}
	path := strings.TrimPrefix(r.URL.Path, "/")
	if strings.HasPrefix(path, "assets/") || path == "vite.svg" {
		data, err := fs.ReadFile(s.opts.UI, path)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		switch {
		case strings.HasSuffix(path, ".js"):
			w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		case strings.HasSuffix(path, ".css"):
			w.Header().Set("Content-Type", "text/css; charset=utf-8")
		case strings.HasSuffix(path, ".svg"):
			w.Header().Set("Content-Type", "image/svg+xml")
		}
		// asset names carry a content hash
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		http.ServeContent(w, r, path, startTime, bytes.NewReader(data))
		return
	}
	if path == "" {
		// the server has no list of runs: it opens the newest one
		if run, err := s.opts.Store.Latest(r.Context()); err == nil {
			http.Redirect(w, r, fmt.Sprintf("/runs/%d", run.ID), http.StatusFound)
			return
		}
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write(s.index)
}
