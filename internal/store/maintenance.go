package store

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/IgorBayerl/nanovision/internal/store/blob"
)

// Retention says which runs Cleanup deletes. Submit runs are always kept.
type Retention struct {
	// review runs older than this are deleted; 0 keeps them all
	ReviewDays int
	// newest local runs kept for each stream; 0 keeps them all
	LocalKeep int
}

func DefaultRetention() Retention {
	return Retention{ReviewDays: 30, LocalKeep: 20}
}

// Cleanup deletes the runs the policy does not keep and returns how many. The
// newest clean local run of a stream is always kept, because it is the base run
// of the next local run. Run GC afterwards to free the blobs.
func (s *Store) Cleanup(ctx context.Context, policy Retention, now time.Time) (int, error) {
	var doomed []int64

	if policy.ReviewDays > 0 {
		cutoff := now.Add(-time.Duration(policy.ReviewDays) * 24 * time.Hour).UnixMilli()
		ids, err := s.ids(ctx, "select id from runs where kind = ? and created_at < ?", string(KindReview), cutoff)
		if err != nil {
			return 0, err
		}
		doomed = append(doomed, ids...)
	}

	if policy.LocalKeep > 0 {
		ids, err := s.ids(ctx, `select id from (
			select id, clean,
				row_number() over (partition by project, stream, profile order by created_at desc, id desc) as age,
				row_number() over (partition by project, stream, profile, clean order by created_at desc, id desc) as clean_age
			from runs where kind = ?)
			where age > ? and not (clean = 1 and clean_age = 1)`, string(KindLocal), policy.LocalKeep)
		if err != nil {
			return 0, err
		}
		doomed = append(doomed, ids...)
	}

	if err := s.DeleteRuns(ctx, doomed); err != nil {
		return 0, err
	}
	return len(doomed), nil
}

func (s *Store) ids(ctx context.Context, query string, args ...any) ([]int64, error) {
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// GCResult reports what GC freed.
type GCResult struct {
	Blobs int   `json:"blobs"`
	Bytes int64 `json:"bytes"`
}

// GC deletes the blobs that no run uses. A blob uploaded for a run that is
// not registered yet can be deleted too; the client then uploads it again.
func (s *Store) GC(ctx context.Context) (GCResult, error) {
	live, err := s.liveBlobs(ctx)
	if err != nil {
		return GCResult{}, err
	}

	rows, err := s.db.QueryContext(ctx, "select hash, length(data) from blobs")
	if err != nil {
		return GCResult{}, err
	}
	var dead []blob.Hash
	var res GCResult
	for rows.Next() {
		var raw []byte
		var size int64
		if err := rows.Scan(&raw, &size); err != nil {
			rows.Close()
			return GCResult{}, err
		}
		var h blob.Hash
		copy(h[:], raw)
		if !live[h] {
			dead = append(dead, h)
			res.Bytes += size
		}
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return GCResult{}, err
	}

	s.writeMu.Lock()
	defer s.writeMu.Unlock()
	err = forEachChunk(dead, func(chunk []blob.Hash) error {
		_, err := s.db.ExecContext(ctx, "delete from blobs where hash in ("+placeholders(len(chunk))+")", hashArgs(chunk)...)
		return err
	})
	if err != nil {
		return GCResult{}, err
	}
	res.Blobs = len(dead)
	return res, nil
}

// liveBlobs marks every blob a run reaches: its manifest, the blobs of the
// entries and the diff.
func (s *Store) liveBlobs(ctx context.Context) (map[blob.Hash]bool, error) {
	rows, err := s.db.QueryContext(ctx, "select manifest, meta from runs")
	if err != nil {
		return nil, err
	}
	type runRefs struct {
		manifest blob.Hash
		meta     string
	}
	var runs []runRefs
	for rows.Next() {
		var raw []byte
		var r runRefs
		if err := rows.Scan(&raw, &r.meta); err != nil {
			rows.Close()
			return nil, err
		}
		copy(r.manifest[:], raw)
		runs = append(runs, r)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	live := make(map[blob.Hash]bool)
	for _, r := range runs {
		if !live[r.manifest] {
			live[r.manifest] = true
			m, err := s.Manifest(ctx, r.manifest)
			if err != nil {
				return nil, fmt.Errorf("gc: %w", err)
			}
			for _, e := range m.Entries {
				// the content hash names the source text, which is not a blob
				for _, h := range []blob.Hash{e.Analysis, e.Coverage} {
					if !h.IsZero() {
						live[h] = true
					}
				}
			}
		}
		var refs struct {
			Diff *blob.Hash `json:"diff"`
		}
		if json.Unmarshal([]byte(r.meta), &refs) == nil && refs.Diff != nil {
			live[*refs.Diff] = true
		}
	}
	return live, nil
}

// Size is how many bytes of its file the store uses. The file does not shrink
// when runs are deleted; SQLite fills the free pages again.
func (s *Store) Size(ctx context.Context) (int64, error) {
	var size int64
	err := s.db.QueryRowContext(ctx, `select (page_count - freelist_count) * page_size
		from pragma_page_count, pragma_freelist_count, pragma_page_size`).Scan(&size)
	return size, err
}

// Trim deletes the oldest runs, and the blobs only they use, until the store
// uses at most maxBytes. It returns how many runs it deleted. The newest run
// always stays.
func (s *Store) Trim(ctx context.Context, maxBytes int64) (int, error) {
	deleted := 0
	for {
		size, err := s.Size(ctx)
		if err != nil || size <= maxBytes {
			return deleted, err
		}
		// a twentieth of the runs each round: few rounds, and not far below the limit
		oldest, err := s.ids(ctx, `select id from runs
			where id != (select id from runs order by created_at desc, id desc limit 1)
			order by created_at, id limit max(1, (select count(*) from runs) / 20)`)
		if err != nil || len(oldest) == 0 {
			return deleted, err
		}
		if err := s.DeleteRuns(ctx, oldest); err != nil {
			return deleted, err
		}
		if _, err := s.GC(ctx); err != nil {
			return deleted, err
		}
		deleted += len(oldest)
	}
}
