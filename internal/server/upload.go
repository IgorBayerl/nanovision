package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/IgorBayerl/nanovision/internal/store"
)

const (
	// hash lists: 16 bytes a hash, far more than one run needs
	maxHashList = 64 << 20
	// one upload request; the client sends about 8 MiB at a time
	maxBlobBody = 64 << 20
	// the blobs of one upload request, unpacked
	maxBlobBytes = 256 << 20
	// run metadata, including the list of changed files
	maxRunBody = 16 << 20
)

func (s *Server) handleMissing(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxHashList))
	if err != nil {
		writeError(w, http.StatusRequestEntityTooLarge, "%v", err)
		return
	}
	hashes, err := DecodeHashes(body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "%v", err)
		return
	}
	missing, err := s.opts.Store.Missing(r.Context(), hashes)
	if err != nil {
		writeError(w, statusFor(err), "%v", err)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(EncodeHashes(missing))
}

func (s *Server) handlePutBlobs(w http.ResponseWriter, r *http.Request) {
	var blobs []store.Blob
	size := 0
	err := DecodeBlobs(http.MaxBytesReader(w, r.Body, maxBlobBody), func(b store.Blob) error {
		// the body is gzip, so its limit does not bound what it unpacks to
		if size += len(b.Data); size > maxBlobBytes {
			return fmt.Errorf("upload unpacks to more than %d bytes; send smaller batches", maxBlobBytes)
		}
		blobs = append(blobs, b)
		return nil
	})
	if err == nil {
		err = s.opts.Store.PutBlobs(r.Context(), blobs)
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, "%v", err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"stored": len(blobs)})
}

func (s *Server) handleAddRun(w http.ResponseWriter, r *http.Request) {
	var run store.Run
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxRunBody))
	if err := dec.Decode(&run); err != nil {
		writeError(w, http.StatusBadRequest, "run: %v", err)
		return
	}
	run.ID = 0
	if run.CreatedAt.IsZero() || run.CreatedAt.After(time.Now().Add(time.Hour)) {
		run.CreatedAt = time.Now()
	}
	id, err := s.opts.Store.AddRun(r.Context(), run)
	if err != nil {
		writeError(w, http.StatusBadRequest, "%v", err)
		return
	}
	s.opts.Logger.Info("Run uploaded", "id", id, "project", run.Project, "stream", run.Stream, "kind", run.Kind, "revision", run.Revision)

	// the first open of a big run should not wait for its summary
	go func() {
		_, err := s.views.get(summaryKey(id), func() ([]byte, error) {
			ctx, cancel := context.WithTimeout(context.Background(), buildTimeout)
			defer cancel()
			return s.buildSummary(ctx, id)
		})
		if err != nil {
			s.opts.Logger.Warn("Preparing the run summary failed", "id", id, "error", err)
		}
	}()
	writeJSON(w, http.StatusCreated, map[string]any{"id": id, "url": s.RunURL(r, id)})
}
