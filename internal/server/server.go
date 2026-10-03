// Package server is `nanovision serve`: the team server. It runs on one
// machine in the studio network; CI uploads runs to it and the web UI shows them.
package server

import (
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/IgorBayerl/nanovision/internal/store"
)

// Options configure a server.
type Options struct {
	Store *store.Store
	// the upload token that CI sends
	Token string
	// nanovision version and build; the build keys the browser cache
	Version string
	Build   string
	// the address people use, e.g. "http://nanovision.studio.local:7070";
	// links returned to the CLI use it
	PublicURL string
	// the web UI files
	UI     fs.FS
	Logger *slog.Logger
	// bytes of built views kept in memory; 0 means 256 MiB
	CacheBytes int64
}

type Server struct {
	opts  Options
	mux   *http.ServeMux
	views *viewCache
	// full-run views rebuild the whole tree; two at a time keeps memory bounded
	heavy chan struct{}
	index []byte
}

func New(opts Options) (*Server, error) {
	if opts.Store == nil {
		return nil, errors.New("server needs a store")
	}
	if opts.Token == "" {
		return nil, errors.New("server needs an upload token")
	}
	if opts.Logger == nil {
		opts.Logger = slog.Default()
	}
	if opts.CacheBytes <= 0 {
		opts.CacheBytes = 256 << 20
	}
	opts.PublicURL = strings.TrimSuffix(opts.PublicURL, "/")

	s := &Server{
		opts:  opts,
		mux:   http.NewServeMux(),
		views: newViewCache(opts.CacheBytes),
		heavy: make(chan struct{}, 2),
	}
	if opts.UI != nil {
		index, err := serverIndex(opts.UI, opts.Build)
		if err != nil {
			return nil, err
		}
		s.index = index
	}
	s.routes()
	return s, nil
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /api/v1/info", s.handleInfo)
	s.mux.HandleFunc("GET /api/v1/base", s.handleBase)
	s.mux.HandleFunc("GET /api/v1/runs/{id}/summary", s.handleSummary)
	s.mux.HandleFunc("GET /api/v1/runs/{id}/file", s.handleFile)
	s.mux.HandleFunc("GET /api/v1/blobs/{hash}", s.handleBlob)

	s.mux.HandleFunc("POST /api/v1/blobs/missing", s.upload(s.handleMissing))
	s.mux.HandleFunc("POST /api/v1/blobs", s.upload(s.handlePutBlobs))
	s.mux.HandleFunc("POST /api/v1/runs", s.upload(s.handleAddRun))

	s.mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotFound, "no such API: %s %s", r.Method, r.URL.Path)
	})
	s.mux.HandleFunc("/", s.handleUI)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
	s.mux.ServeHTTP(rec, r)
	if strings.HasPrefix(r.URL.Path, "/api/") {
		s.opts.Logger.Debug("request", "method", r.Method, "path", r.URL.Path, "status", rec.status, "duration", time.Since(start))
	}
}

// upload guards the write API with the token.
func (s *Server) upload(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		got, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
		if !ok || subtle.ConstantTimeCompare([]byte(got), []byte(s.opts.Token)) != 1 {
			writeError(w, http.StatusUnauthorized, "missing or wrong upload token")
			return
		}
		next(w, r)
	}
}

// RunURL is the page of a run, for the CLI to print.
func (s *Server) RunURL(r *http.Request, id int64) string {
	base := s.opts.PublicURL
	if base == "" && r != nil {
		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}
		base = scheme + "://" + r.Host
	}
	return fmt.Sprintf("%s/runs/%d", base, id)
}

func (s *Server) handleInfo(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"version": s.opts.Version, "build": s.opts.Build})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(v)
}

func writeError(w http.ResponseWriter, status int, format string, args ...any) {
	writeJSON(w, status, map[string]string{"error": fmt.Sprintf(format, args...)})
}

// statusFor maps store errors to HTTP statuses.
func statusFor(err error) int {
	if errors.Is(err, store.ErrNotFound) {
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}

func pathID(r *http.Request, name string) (int64, error) {
	id, err := strconv.ParseInt(r.PathValue(name), 10, 64)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("%s %q is not a run ID", name, r.PathValue(name))
	}
	return id, nil
}
