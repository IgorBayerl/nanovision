package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/IgorBayerl/nanovision/internal/store/blob"
)

// RunKind is how a run was made. Only a clean run can be a base run.
type RunKind string

const (
	// CI, after a submit or a merge; kept forever
	KindSubmit RunKind = "submit"
	// CI, for a shelved changelist or a pull request
	KindReview RunKind = "review"
	// a developer
	KindLocal RunKind = "local"
)

func ParseRunKind(s string) (RunKind, error) {
	switch k := RunKind(strings.ToLower(strings.TrimSpace(s))); k {
	case KindSubmit, KindReview, KindLocal:
		return k, nil
	case "":
		return KindLocal, nil
	default:
		return "", fmt.Errorf("unknown run kind %q (want submit, review or local)", s)
	}
}

// Run is one saved execution of nanovision. A saved run never changes.
type Run struct {
	ID      int64   `json:"id"`
	Project string  `json:"project"`
	Stream  string  `json:"stream"`
	Profile string  `json:"profile"`
	Kind    RunKind `json:"kind"`
	// "//game/main@118432" for Perforce, a commit hash for git
	Revision string `json:"revision"`
	// the changelist number, 0 when the revision has none
	Change    int64     `json:"change,omitempty"`
	Clean     bool      `json:"clean"`
	ReviewID  string    `json:"reviewId,omitempty"`
	Author    string    `json:"author,omitempty"`
	CIURL     string    `json:"ciUrl,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	Tool      string    `json:"tool"`
	Manifest  blob.Hash `json:"manifest"`
	// totals, thresholds and report patterns; see history.Meta
	Meta json.RawMessage `json:"meta"`
}

const runColumns = "id, project, stream, profile, revision, change, kind, clean, review_id, author, ci_url, created_at, tool, manifest, meta"

func scanRun(row interface{ Scan(...any) error }) (Run, error) {
	var r Run
	var change sql.NullInt64
	var reviewID, author, ciURL sql.NullString
	var created int64
	var clean int
	var manifest []byte
	var meta string
	err := row.Scan(&r.ID, &r.Project, &r.Stream, &r.Profile, &r.Revision, &change, &r.Kind, &clean,
		&reviewID, &author, &ciURL, &created, &r.Tool, &manifest, &meta)
	if err != nil {
		return Run{}, err
	}
	r.Change = change.Int64
	r.Clean = clean != 0
	r.ReviewID, r.Author, r.CIURL = reviewID.String, author.String, ciURL.String
	r.CreatedAt = time.UnixMilli(created).UTC()
	copy(r.Manifest[:], manifest)
	r.Meta = json.RawMessage(meta)
	return r, nil
}

func nullString(s string) sql.NullString { return sql.NullString{String: s, Valid: s != ""} }

// AddRun saves a run whose blobs are already stored and returns its ID.
func (s *Store) AddRun(ctx context.Context, r Run) (int64, error) {
	if r.Project == "" || r.Profile == "" {
		return 0, errors.New("a run needs a project and a profile")
	}
	if _, err := ParseRunKind(string(r.Kind)); err != nil || r.Kind == "" {
		return 0, fmt.Errorf("run kind %q: want submit, review or local", r.Kind)
	}
	if len(r.Meta) == 0 {
		r.Meta = json.RawMessage("{}")
	}
	if !json.Valid(r.Meta) {
		return 0, errors.New("run meta is not valid JSON")
	}
	if r.CreatedAt.IsZero() {
		r.CreatedAt = time.Now()
	}

	manifest, err := s.Manifest(ctx, r.Manifest)
	if err != nil {
		return 0, fmt.Errorf("run manifest: %w", err)
	}
	if err := s.checkReferences(ctx, manifest, r.Meta); err != nil {
		return 0, err
	}

	s.writeMu.Lock()
	defer s.writeMu.Unlock()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var change sql.NullInt64
	if r.Change > 0 {
		change = sql.NullInt64{Int64: r.Change, Valid: true}
	}
	res, err := tx.ExecContext(ctx, `insert into runs
		(project, stream, profile, revision, change, kind, clean, review_id, author, ci_url, created_at, tool, manifest, meta)
		values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		r.Project, r.Stream, r.Profile, r.Revision, change, string(r.Kind), boolInt(r.Clean),
		nullString(r.ReviewID), nullString(r.Author), nullString(r.CIURL), r.CreatedAt.UnixMilli(),
		r.Tool, r.Manifest[:], string(r.Meta))
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, tx.Commit()
}

// checkReferences makes sure every blob a run needs is stored. Source text is
// optional: it can be excluded on purpose or be unreadable.
func (s *Store) checkReferences(ctx context.Context, m blob.Manifest, meta json.RawMessage) error {
	var need []blob.Hash
	for _, e := range m.Entries {
		if !e.Coverage.IsZero() {
			need = append(need, e.Coverage)
		}
		if !e.Analysis.IsZero() {
			need = append(need, e.Analysis)
		}
	}
	var refs struct {
		Diff *blob.Hash `json:"diff"`
	}
	if err := json.Unmarshal(meta, &refs); err == nil && refs.Diff != nil && !refs.Diff.IsZero() {
		need = append(need, *refs.Diff)
	}
	missing, err := s.Missing(ctx, need)
	if err != nil {
		return err
	}
	if len(missing) > 0 {
		return fmt.Errorf("run references %d blobs the store does not have, for example %s", len(missing), missing[0])
	}
	return nil
}

func (s *Store) Run(ctx context.Context, id int64) (Run, error) {
	r, err := scanRun(s.db.QueryRowContext(ctx, "select "+runColumns+" from runs where id = ?", id))
	if errors.Is(err, sql.ErrNoRows) {
		return Run{}, fmt.Errorf("run %d: %w", id, ErrNotFound)
	}
	return r, err
}

// CandidateQuery selects the clean runs a delta may compare with. Set
// Revisions for git or Change for Perforce.
type CandidateQuery struct {
	Project string
	Profile string
	// empty matches every stream
	Stream string
	// empty matches every kind
	Kinds []RunKind
	// runs at any of these revisions
	Revisions []string
	// the nearest run at or above this changelist, and the nearest below it
	Change int64
}

// Candidates returns the runs of q, newest first for revisions; for a
// changelist at most two, the one at or above it first.
func (s *Store) Candidates(ctx context.Context, q CandidateQuery) ([]Run, error) {
	where := "project = ? and profile = ? and clean = 1"
	args := []any{q.Project, q.Profile}
	if q.Stream != "" {
		where += " and stream = ?"
		args = append(args, q.Stream)
	}
	if len(q.Kinds) > 0 {
		where += " and kind in (" + placeholders(len(q.Kinds)) + ")"
		for _, k := range q.Kinds {
			args = append(args, string(k))
		}
	}
	if len(q.Revisions) > 0 {
		for _, r := range q.Revisions {
			args = append(args, r)
		}
		return s.runs(ctx, where+" and revision in ("+placeholders(len(q.Revisions))+") order by created_at desc, id desc", args...)
	}
	args = append(args, q.Change)
	above, err := s.runs(ctx, where+" and change >= ? order by change, created_at desc, id desc limit 1", args...)
	if err != nil {
		return nil, err
	}
	below, err := s.runs(ctx, where+" and change < ? order by change desc, created_at desc, id desc limit 1", args...)
	return append(above, below...), err
}

// Latest is the newest run, a submit run when the store has one.
func (s *Store) Latest(ctx context.Context) (Run, error) {
	runs, err := s.runs(ctx, "1 order by kind = 'submit' desc, created_at desc, id desc limit 1")
	if err != nil {
		return Run{}, err
	}
	if len(runs) == 0 {
		return Run{}, fmt.Errorf("the store has no runs: %w", ErrNotFound)
	}
	return runs[0], nil
}

// runs reads the runs matching a where clause, which may end in order and limit.
func (s *Store) runs(ctx context.Context, where string, args ...any) ([]Run, error) {
	rows, err := s.db.QueryContext(ctx, "select "+runColumns+" from runs where "+where, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var runs []Run
	for rows.Next() {
		r, err := scanRun(rows)
		if err != nil {
			return nil, err
		}
		runs = append(runs, r)
	}
	return runs, rows.Err()
}

// DeleteRuns removes runs. Their blobs stay until the next GC.
func (s *Store) DeleteRuns(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	s.writeMu.Lock()
	defer s.writeMu.Unlock()
	args := make([]any, len(ids))
	for i, id := range ids {
		args[i] = id
	}
	_, err := s.db.ExecContext(ctx, "delete from runs where id in ("+placeholders(len(ids))+")", args...)
	return err
}

func boolInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
