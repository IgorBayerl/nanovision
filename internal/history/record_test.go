package history

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/IgorBayerl/nanovision/internal/diff"
	"github.com/IgorBayerl/nanovision/internal/store"
	"github.com/IgorBayerl/nanovision/internal/store/blob"
	"github.com/IgorBayerl/nanovision/internal/vcs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// emptyRun stores a run with an empty manifest.
func emptyRun(t *testing.T, s *store.Store, r store.Run) store.Run {
	t.Helper()
	ctx := context.Background()
	m := blob.EncodeManifest(blob.Manifest{})
	require.NoError(t, s.PutBlobs(ctx, []store.Blob{{Hash: blob.Sum(m), Data: m}}))
	r.Manifest, r.Tool = blob.Sum(m), "test"
	if r.Project == "" {
		r.Project = "game"
	}
	if r.Profile == "" {
		r.Profile = "unit"
	}
	if r.Kind == "" {
		r.Kind, r.Clean = store.KindSubmit, true
	}
	id, err := s.AddRun(ctx, r)
	require.NoError(t, err)
	r.ID = id
	return r
}

type fakeVCS struct {
	changes   map[[2]int64]int
	ancestors []string
}

func (f *fakeVCS) Name() string                              { return "fake" }
func (f *fakeVCS) Current() (vcs.Revision, error)            { return vcs.Revision{}, nil }
func (f *fakeVCS) HasLocalEdits() (bool, error)              { return false, nil }
func (f *fakeVCS) Base() (vcs.Revision, error)               { return vcs.Revision{}, nil }
func (f *fakeVCS) Diff(vcs.Revision) (*diff.DiffData, error) { return &diff.DiffData{}, nil }
func (f *fakeVCS) Author(bool) (string, error)               { return "", nil }

type perforceFake struct{ fakeVCS }

func (f *perforceFake) ChangesBetween(low, high int64) (int, error) {
	n, ok := f.changes[[2]int64{low, high}]
	if !ok {
		return 0, fmt.Errorf("unexpected range %d..%d", low, high)
	}
	return n, nil
}

type gitFake struct{ fakeVCS }

func (f *gitFake) Ancestors(from vcs.Revision, max int) ([]string, error) {
	return f.ancestors[:min(max+1, len(f.ancestors))], nil
}

func perforceStore(t *testing.T) *store.Store {
	s, err := store.Open(t.TempDir())
	require.NoError(t, err)
	t.Cleanup(func() { s.Close() })
	for _, c := range []int64{100, 105, 120} {
		emptyRun(t, s, store.Run{Stream: "//game/main", Revision: vcs.PerforceRevision("//game/main", c), Change: c})
	}
	// neither of these can be a base run of //game/main
	emptyRun(t, s, store.Run{Stream: "//game/release", Revision: "//game/release@109", Change: 109})
	emptyRun(t, s, store.Run{Stream: "//game/main", Revision: "//game/main@108", Change: 108, Kind: store.KindReview})
	return s
}

func query(change int64) BaseQuery {
	return BaseQuery{
		Project: "game", Profile: "unit", MaxDistance: 50, Kinds: []store.RunKind{store.KindSubmit},
		Base: vcs.ParseRevision(vcs.PerforceRevision("//game/main", change)),
	}
}

func TestFindBasePerforce(t *testing.T) {
	s := perforceStore(t)
	ctx := context.Background()

	t.Run("a higher changelist with no change in the view is exact", func(t *testing.T) {
		v := &perforceFake{fakeVCS{changes: map[[2]int64]int{{110, 120}: 0}}}
		base, err := FindBase(ctx, s, v, query(110))
		require.NoError(t, err)
		require.NotNil(t, base)
		assert.Equal(t, int64(120), base.Run.Change)
		assert.True(t, base.Exact)
	})

	t.Run("otherwise the nearest older run, with the distance", func(t *testing.T) {
		v := &perforceFake{fakeVCS{changes: map[[2]int64]int{{110, 120}: 2, {105, 110}: 3}}}
		base, err := FindBase(ctx, s, v, query(110))
		require.NoError(t, err)
		assert.Equal(t, int64(105), base.Run.Change)
		assert.False(t, base.Exact)
		assert.Equal(t, 3, base.Distance)
	})

	t.Run("a submit build compares with the previous run", func(t *testing.T) {
		v := &perforceFake{fakeVCS{changes: map[[2]int64]int{{105, 120}: 1}}}
		q := query(120)
		q.OlderOnly = true
		base, err := FindBase(ctx, s, v, q)
		require.NoError(t, err)
		assert.Equal(t, int64(105), base.Run.Change)
		assert.True(t, base.Exact, "no view change between the two runs")
	})

	t.Run("too far back is no base", func(t *testing.T) {
		v := &perforceFake{fakeVCS{changes: map[[2]int64]int{{110, 120}: 2, {105, 110}: 3}}}
		q := query(110)
		q.MaxDistance = 2
		base, err := FindBase(ctx, s, v, q)
		require.NoError(t, err)
		assert.Nil(t, base)
	})

	t.Run("flags only: exact or older, without a distance", func(t *testing.T) {
		base, err := FindBase(ctx, s, nil, query(105))
		require.NoError(t, err)
		assert.True(t, base.Exact)

		base, err = FindBase(ctx, s, nil, query(110))
		require.NoError(t, err)
		assert.Equal(t, int64(105), base.Run.Change)
		assert.Equal(t, -1, base.Distance)
		assert.False(t, base.Exact)

		base, err = FindBase(ctx, s, nil, query(99))
		require.NoError(t, err)
		assert.Nil(t, base)
	})
}

func TestFindBaseGit(t *testing.T) {
	s, err := store.Open(t.TempDir())
	require.NoError(t, err)
	defer s.Close()
	ctx := context.Background()
	emptyRun(t, s, store.Run{Stream: "main", Revision: "c2"})
	emptyRun(t, s, store.Run{Stream: "main", Revision: "c3", Kind: store.KindLocal, Clean: false})

	q := BaseQuery{Project: "game", Profile: "unit", MaxDistance: 50, Kinds: []store.RunKind{store.KindSubmit, store.KindLocal},
		Base: vcs.Revision{ID: "c3"}}
	v := &gitFake{fakeVCS{ancestors: []string{"c3", "c2", "c1"}}}

	base, err := FindBase(ctx, s, v, q)
	require.NoError(t, err)
	require.NotNil(t, base)
	assert.Equal(t, "c2", base.Run.Revision, "the run at c3 has local edits and is no base run")
	assert.Equal(t, 1, base.Distance)

	q.OlderOnly = true
	base, err = FindBase(ctx, s, v, q)
	require.NoError(t, err)
	assert.Equal(t, "c2", base.Run.Revision)
	assert.True(t, base.Exact)

	q.OlderOnly, q.MaxDistance = false, 0
	base, err = FindBase(ctx, s, v, q)
	require.NoError(t, err)
	assert.Nil(t, base)

	base, err = FindBase(ctx, s, nil, BaseQuery{Project: "game", Profile: "unit", Kinds: q.Kinds, Base: vcs.Revision{ID: "c2"}})
	require.NoError(t, err)
	assert.True(t, base.Exact, "without the adapter only the exact commit matches")
}

type countingTarget struct {
	*store.Store
	uploaded int
}

func (c *countingTarget) PutBlobs(ctx context.Context, blobs []store.Blob) error {
	c.uploaded += len(blobs)
	return c.Store.PutBlobs(ctx, blobs)
}

func TestSaveUploadsOnlyMissingBlobs(t *testing.T) {
	tree, cfg := demoTree(t, false)
	capture, err := NewCapture(tree, cfg, "")
	require.NoError(t, err)

	s, err := store.Open(t.TempDir())
	require.NoError(t, err)
	defer s.Close()
	target := &countingTarget{Store: s}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	run := store.Run{Project: "demo", Stream: "main", Profile: "default", Kind: store.KindLocal, Tool: "test"}
	id, err := Save(context.Background(), target, capture, run, logger)
	require.NoError(t, err)
	first := target.uploaded
	assert.Equal(t, len(capture.Blobs), first)

	run.CreatedAt = time.Now().Add(time.Second)
	id2, err := Save(context.Background(), target, capture, run, logger)
	require.NoError(t, err)
	assert.NotEqual(t, id, id2)
	assert.Equal(t, first, target.uploaded, "the second run of the same code uploads nothing")

	stored, err := s.Run(context.Background(), id2)
	require.NoError(t, err)
	var meta Meta
	require.NoError(t, json.Unmarshal(stored.Meta, &meta))
	assert.Equal(t, capture.Meta.Files, meta.Files)
}
