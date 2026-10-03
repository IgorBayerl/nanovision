package history

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/store"
	"github.com/IgorBayerl/nanovision/internal/store/blob"
	"github.com/IgorBayerl/nanovision/internal/vcs"
)

// Target is a run store the CLI talks to: a local store, or a team server.
// *store.Store is one.
type Target interface {
	Missing(ctx context.Context, hashes []blob.Hash) ([]blob.Hash, error)
	PutBlobs(ctx context.Context, blobs []store.Blob) error
	AddRun(ctx context.Context, run store.Run) (int64, error)
	Candidates(ctx context.Context, q store.CandidateQuery) ([]store.Run, error)
	Manifest(ctx context.Context, hash blob.Hash) (blob.Manifest, error)
}

// BaseQuery describes the run to compare with.
type BaseQuery struct {
	Project, Stream, Profile string
	// the revision the change starts from
	Base vcs.Revision
	// the workspace is clean and at Base, so the search starts below it:
	// a run never compares with its own revision
	OlderOnly   bool
	MaxDistance int
	// the kinds of run that can be a base run: submit runs on a team server,
	// submit and clean local runs in a local store
	Kinds []store.RunKind
}

// Base is the run a delta compares with.
type Base struct {
	Run   store.Run
	Exact bool
	// revisions between the base revision and the base run, -1 when unknown
	Distance int
}

// FindBase applies the rules of section 6 of the delta plan. It returns nil
// when no run qualifies. v may be nil: then only exact revisions match, plus
// older changelists of the same stream for Perforce revisions.
func FindBase(ctx context.Context, t Target, v vcs.VCS, q BaseQuery) (*Base, error) {
	if q.Base.ID == "" {
		return nil, nil
	}
	if q.Base.Change > 0 {
		return findPerforceBase(ctx, t, v, q)
	}
	return findCommitBase(ctx, t, v, q)
}

func (q BaseQuery) candidates() store.CandidateQuery {
	return store.CandidateQuery{Project: q.Project, Profile: q.Profile, Kinds: q.Kinds}
}

func findPerforceBase(ctx context.Context, t Target, v vcs.VCS, q BaseQuery) (*Base, error) {
	// changelist numbers are global to the server; the stream keeps a release
	// stream's runs out of the main stream's search
	stream := q.Base.Branch
	if stream == "" {
		stream = q.Stream
	}
	counter, _ := v.(vcs.ChangeCounter)
	change := q.Base.Change
	cq := q.candidates()
	cq.Stream, cq.Change = stream, change
	runs, err := t.Candidates(ctx, cq)
	if err != nil {
		return nil, err
	}
	var below *store.Run
	for i, r := range runs {
		if r.Change < change {
			below = &runs[i]
			continue
		}
		if q.OlderOnly {
			continue
		}
		if r.Change == change {
			return &Base{Run: r, Exact: true}, nil
		}
		// a run at a higher changelist can hold the same code, when the
		// changelists between changed only files outside the workspace view
		if counter != nil {
			n, err := counter.ChangesBetween(change, r.Change)
			if err != nil {
				return nil, err
			}
			if n == 0 {
				return &Base{Run: r, Exact: true}, nil
			}
		}
	}
	if below == nil {
		return nil, nil
	}
	base := &Base{Run: *below, Distance: -1}
	if counter != nil {
		n, err := counter.ChangesBetween(below.Change, change)
		if err != nil {
			return nil, err
		}
		if q.OlderOnly {
			n-- // the current changelist is not a revision between the two
		}
		base.Distance = max(n, 0)
		base.Exact = base.Distance == 0
		if base.Distance > q.MaxDistance {
			return nil, nil
		}
	}
	return base, nil
}

func findCommitBase(ctx context.Context, t Target, v vcs.VCS, q BaseQuery) (*Base, error) {
	candidates := []string{q.Base.ID}
	if ancestry, ok := v.(vcs.Ancestry); ok {
		list, err := ancestry.Ancestors(q.Base, q.MaxDistance+1)
		if err != nil {
			return nil, err
		}
		if len(list) > 0 {
			candidates = list
		}
	}
	if q.OlderOnly {
		if len(candidates) < 2 {
			return nil, nil
		}
		candidates = candidates[1:]
	}
	if len(candidates) > q.MaxDistance+1 {
		candidates = candidates[:q.MaxDistance+1]
	}

	cq := q.candidates()
	cq.Revisions = candidates
	runs, err := t.Candidates(ctx, cq)
	if err != nil {
		return nil, err
	}
	newest := make(map[string]store.Run)
	for _, r := range runs {
		if _, ok := newest[r.Revision]; !ok {
			newest[r.Revision] = r // runs come newest first
		}
	}
	for i, rev := range candidates {
		if r, ok := newest[rev]; ok {
			return &Base{Run: r, Exact: i == 0, Distance: i}, nil
		}
	}
	return nil, nil
}

// Ref names a stored run inside a comparison.
func Ref(r store.Run) model.RunRef {
	return model.RunRef{
		ID: r.ID, Revision: r.Revision, Stream: r.Stream, Profile: r.Profile,
		Kind: string(r.Kind), CreatedAt: r.CreatedAt.UnixMilli(), Tool: r.Tool,
	}
}

// how many bytes one upload request carries
const uploadBatch = 8 << 20

// Save uploads the blobs the target is missing and registers the run. The meta and manifest of run come from c.
func Save(ctx context.Context, t Target, c *Capture, run store.Run, logger *slog.Logger) (int64, error) {
	meta, err := json.Marshal(c.Meta)
	if err != nil {
		return 0, err
	}
	run.Manifest, run.Meta = c.ManifestHash, meta

	id, err := save(ctx, t, c, run, logger)
	// a GC on the server can remove an orphan blob between the check and the
	// registration; one retry uploads it again
	if err != nil && strings.Contains(err.Error(), "does not have") {
		logger.Debug("Retrying the upload", "error", err)
		id, err = save(ctx, t, c, run, logger)
	}
	return id, err
}

func save(ctx context.Context, t Target, c *Capture, run store.Run, logger *slog.Logger) (int64, error) {
	hashes := make([]blob.Hash, len(c.Blobs))
	byHash := make(map[blob.Hash][]byte, len(c.Blobs))
	for i, b := range c.Blobs {
		hashes[i] = b.Hash
		byHash[b.Hash] = b.Data
	}
	missing, err := t.Missing(ctx, hashes)
	if err != nil {
		return 0, err
	}

	var batch []store.Blob
	size := 0
	flush := func() error {
		if len(batch) == 0 {
			return nil
		}
		err := t.PutBlobs(ctx, batch)
		batch, size = nil, 0
		return err
	}
	for _, h := range missing {
		batch = append(batch, store.Blob{Hash: h, Data: byHash[h]})
		if size += len(byHash[h]); size >= uploadBatch {
			if err := flush(); err != nil {
				return 0, err
			}
		}
	}
	if err := flush(); err != nil {
		return 0, err
	}
	logger.Debug("Blobs uploaded", "missing", len(missing), "total", len(hashes))

	id, err := t.AddRun(ctx, run)
	if err != nil {
		return 0, fmt.Errorf("register run: %w", err)
	}
	return id, nil
}
