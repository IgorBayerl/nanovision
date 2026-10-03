package store

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/IgorBayerl/nanovision/internal/store/blob"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func openTemp(t *testing.T) *Store {
	t.Helper()
	s, err := Open(t.TempDir())
	require.NoError(t, err)
	t.Cleanup(func() { s.Close() })
	return s
}

// fakeRun stores a manifest of n files, each with its own coverage blob, and
// returns the run ready to add. seed changes the coverage of the first file.
func fakeRun(t *testing.T, s *Store, n, seed int) Run {
	t.Helper()
	ctx := context.Background()
	var blobs []Blob
	m := blob.Manifest{}
	for i := range n {
		hits := 1
		if i == 0 {
			hits = seed
		}
		cov := blob.EncodeCoverage(blob.Coverage{Reports: 1, Lines: []blob.Line{{Number: 1, Hits: hits}}})
		blobs = append(blobs, Blob{Hash: blob.Sum(cov), Data: cov})
		m.Entries = append(m.Entries, blob.Entry{
			Path:     fmt.Sprintf("src/file%03d.go", i),
			Content:  blob.Sum(fmt.Appendf(nil, "file %d", i)),
			Coverage: blob.Sum(cov),
			Counts:   blob.Counts{LinesValid: 1, LinesCovered: min(hits, 1)},
		})
	}
	data := blob.EncodeManifest(m)
	blobs = append(blobs, Blob{Hash: blob.Sum(data), Data: data})
	require.NoError(t, s.PutBlobs(ctx, blobs))
	return Run{
		Project: "game", Stream: "//game/main", Profile: "unit", Kind: KindSubmit, Clean: true,
		Tool: "test", Manifest: blob.Sum(data), Meta: json.RawMessage(`{"title":"t"}`),
	}
}

func TestPutBlobsVerifiesHashes(t *testing.T) {
	s := openTemp(t)
	ctx := context.Background()

	cov := blob.EncodeCoverage(blob.Coverage{Reports: 1})
	err := s.PutBlobs(ctx, []Blob{{Hash: blob.Sum([]byte("other")), Data: cov}})
	assert.ErrorContains(t, err, "hashes to")

	require.NoError(t, s.PutBlobs(ctx, []Blob{{Hash: blob.Sum(cov), Data: cov}}))
	require.NoError(t, s.PutBlobs(ctx, []Blob{{Hash: blob.Sum(cov), Data: cov}}), "storing twice is fine")

	missing, err := s.Missing(ctx, []blob.Hash{blob.Sum([]byte("nope")), blob.Sum(cov), blob.Sum([]byte("nope"))})
	require.NoError(t, err)
	assert.Equal(t, []blob.Hash{blob.Sum([]byte("nope"))}, missing)
}

func TestAddRunAndQuery(t *testing.T) {
	s := openTemp(t)
	ctx := context.Background()

	r := fakeRun(t, s, 3, 1)
	r.Revision, r.Change = "//game/main@100", 100
	id, err := s.AddRun(ctx, r)
	require.NoError(t, err)

	got, err := s.Run(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, "//game/main@100", got.Revision)
	assert.Equal(t, int64(100), got.Change)
	assert.True(t, got.Clean)
	assert.JSONEq(t, `{"title":"t"}`, string(got.Meta))

	_, err = s.Run(ctx, id+1)
	assert.ErrorIs(t, err, ErrNotFound)

	r2 := fakeRun(t, s, 3, 2)
	r2.Revision, r2.Change, r2.Kind, r2.Clean = "//game/main@105", 105, KindReview, false
	_, err = s.AddRun(ctx, r2)
	require.NoError(t, err)

	r3 := fakeRun(t, s, 3, 3)
	r3.Revision, r3.Change = "//game/main@110", 110
	id3, err := s.AddRun(ctx, r3)
	require.NoError(t, err)

	q := CandidateQuery{Project: "game", Stream: "//game/main", Profile: "unit", Kinds: []RunKind{KindSubmit}, Change: 105}
	near, err := s.Candidates(ctx, q)
	require.NoError(t, err)
	require.Len(t, near, 2, "the review run at 105 has local edits and is no candidate")
	assert.Equal(t, []int64{id3, id}, []int64{near[0].ID, near[1].ID}, "the run at or above the changelist, then the one below")

	q.Change, q.Revisions = 0, []string{"//game/main@100", "//game/main@105"}
	at, err := s.Candidates(ctx, q)
	require.NoError(t, err)
	require.Len(t, at, 1)
	assert.Equal(t, id, at[0].ID)

	latest, err := s.Latest(ctx)
	require.NoError(t, err)
	assert.Equal(t, id3, latest.ID)
}

func TestLatestPrefersSubmitRuns(t *testing.T) {
	s := openTemp(t)
	ctx := context.Background()
	_, err := s.Latest(ctx)
	assert.ErrorIs(t, err, ErrNotFound)

	submit := fakeRun(t, s, 1, 1)
	submit.CreatedAt = time.Now().Add(-time.Hour)
	id, err := s.AddRun(ctx, submit)
	require.NoError(t, err)
	review := fakeRun(t, s, 1, 2)
	review.Kind = KindReview
	_, err = s.AddRun(ctx, review)
	require.NoError(t, err)

	latest, err := s.Latest(ctx)
	require.NoError(t, err)
	assert.Equal(t, id, latest.ID)
}

func TestAddRunRejectsMissingBlobs(t *testing.T) {
	s := openTemp(t)
	ctx := context.Background()

	m := blob.EncodeManifest(blob.Manifest{Entries: []blob.Entry{{Path: "a.go", Coverage: blob.Sum([]byte("gone"))}}})
	require.NoError(t, s.PutBlobs(ctx, []Blob{{Hash: blob.Sum(m), Data: m}}))
	_, err := s.AddRun(ctx, Run{Project: "p", Profile: "d", Kind: KindLocal, Tool: "t", Manifest: blob.Sum(m)})
	assert.ErrorContains(t, err, "does not have")
}

func TestTrimDeletesTheOldestRuns(t *testing.T) {
	s := openTemp(t)
	ctx := context.Background()
	now := time.Now()

	var ids []int64
	for i := range 6 {
		r := fakeRun(t, s, 40, 100+i)
		r.CreatedAt = now.Add(time.Duration(i) * time.Minute)
		id, err := s.AddRun(ctx, r)
		require.NoError(t, err)
		ids = append(ids, id)
	}
	size, err := s.Size(ctx)
	require.NoError(t, err)
	require.Positive(t, size)

	deleted, err := s.Trim(ctx, size)
	require.NoError(t, err)
	assert.Zero(t, deleted, "a store within its limit keeps every run")

	deleted, err = s.Trim(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, 5, deleted, "a limit nothing fits in still keeps the newest run")
	left, err := s.runs(ctx, "1")
	require.NoError(t, err)
	require.Len(t, left, 1)
	assert.Equal(t, ids[5], left[0].ID)
	after, err := s.Size(ctx)
	require.NoError(t, err)
	assert.Less(t, after, size)
}

func TestCleanupAndGC(t *testing.T) {
	s := openTemp(t)
	ctx := context.Background()
	now := time.Now()

	var ids []int64
	for i := range 4 {
		r := fakeRun(t, s, 2, 10+i)
		r.Kind, r.Clean = KindLocal, i == 0
		r.CreatedAt = now.Add(time.Duration(i) * time.Minute)
		id, err := s.AddRun(ctx, r)
		require.NoError(t, err)
		ids = append(ids, id)
	}
	old := fakeRun(t, s, 2, 99)
	old.Kind, old.Clean, old.CreatedAt = KindReview, false, now.Add(-40*24*time.Hour)
	_, err := s.AddRun(ctx, old)
	require.NoError(t, err)

	deleted, err := s.Cleanup(ctx, Retention{ReviewDays: 30, LocalKeep: 2}, now.Add(time.Hour))
	require.NoError(t, err)
	assert.Equal(t, 2, deleted, "the old review run and the third-newest local run go")

	left, err := s.runs(ctx, "1")
	require.NoError(t, err)
	var leftIDs []int64
	for _, r := range left {
		leftIDs = append(leftIDs, r.ID)
	}
	assert.ElementsMatch(t, []int64{ids[0], ids[2], ids[3]}, leftIDs, "the only clean run is the base and stays")

	res, err := s.GC(ctx)
	require.NoError(t, err)
	assert.Positive(t, res.Blobs)

	for _, r := range left {
		m, err := s.Manifest(ctx, r.Manifest)
		require.NoError(t, err, "every kept run can still be read")
		hashes := []blob.Hash{}
		for _, e := range m.Entries {
			hashes = append(hashes, e.Coverage)
		}
		missing, err := s.Missing(ctx, hashes)
		require.NoError(t, err)
		assert.Empty(t, missing)
	}
}

func TestBackup(t *testing.T) {
	s := openTemp(t)
	ctx := context.Background()
	id, err := s.AddRun(ctx, fakeRun(t, s, 3, 1))
	require.NoError(t, err)

	path := filepath.Join(t.TempDir(), "copy.db")
	require.NoError(t, s.Backup(ctx, path))
	assert.Error(t, s.Backup(ctx, path), "a backup never overwrites")

	copyStore, err := OpenFile(path)
	require.NoError(t, err)
	defer copyStore.Close()
	_, err = copyStore.Run(ctx, id)
	require.NoError(t, err)
}

func TestParseRunKind(t *testing.T) {
	k, err := ParseRunKind("")
	require.NoError(t, err)
	assert.Equal(t, KindLocal, k)
	k, err = ParseRunKind("Submit")
	require.NoError(t, err)
	assert.Equal(t, KindSubmit, k)
	_, err = ParseRunKind("nightly")
	assert.Error(t, err)
}
