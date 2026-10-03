package server

import (
	"context"
	"time"

	"github.com/IgorBayerl/nanovision/internal/store"
)

// Upkeep is what RunMaintenance keeps the store to.
type Upkeep struct {
	// review runs older than this are deleted; 0 keeps them all
	ReviewDays int
	// above this size the oldest runs are deleted; 0 means no limit
	MaxBytes int64
}

// RunMaintenance applies u at the start and then once a day until ctx ends.
// Backups are a job for the machine's scheduler: `nanovision store backup`.
func (s *Server) RunMaintenance(ctx context.Context, u Upkeep) {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()
	for {
		s.maintain(ctx, u)
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (s *Server) maintain(ctx context.Context, u Upkeep) {
	log := s.opts.Logger
	if u.ReviewDays > 0 {
		deleted, err := s.opts.Store.Cleanup(ctx, store.Retention{ReviewDays: u.ReviewDays}, time.Now())
		if err != nil {
			log.Error("Cleanup failed", "error", err)
			return
		}
		if deleted > 0 {
			res, err := s.opts.Store.GC(ctx)
			if err != nil {
				log.Error("Blob cleanup failed", "error", err)
				return
			}
			log.Info("Old review runs deleted", "runs", deleted, "blobs", res.Blobs, "bytes", res.Bytes)
		}
	}
	if u.MaxBytes > 0 {
		deleted, err := s.opts.Store.Trim(ctx, u.MaxBytes)
		if err != nil {
			log.Error("Keeping the store within its size limit failed", "error", err)
			return
		}
		if deleted > 0 {
			log.Info("Oldest runs deleted to keep the store within its size limit", "runs", deleted, "limit", u.MaxBytes)
		}
	}
}
