// Package store keeps runs and their blobs in one SQLite file.
//
// Only one process writes to a store: the CLI for a local store, or
// `nanovision serve` for a team store. Do not put the file on a network share;
// SQLite locks are not reliable there.
package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

// FileName is the SQLite file inside a store folder.
const FileName = "store.db"

// ErrNotFound is returned for a run or blob the store does not have.
var ErrNotFound = errors.New("not found")

type Store struct {
	db   *sql.DB
	path string

	// SQLite has one writer; serializing here avoids busy retries
	writeMu sync.Mutex
}

// Open opens the store in dir, creating the folder and the file when needed.
func Open(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("create store folder: %w", err)
	}
	return OpenFile(filepath.Join(dir, FileName))
}

// OpenFile opens a store file directly, for example a backup.
func OpenFile(path string) (*Store, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	dsn := "file:" + filepath.ToSlash(abs) + "?" + url.Values{
		"_journal_mode": {"WAL"},
		"_synchronous":  {"NORMAL"},
		"_busy_timeout": {"10000"},
		"_foreign_keys": {"on"},
	}.Encode()

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("open store %s: %w", abs, err)
	}
	s := &Store{db: db, path: abs}
	if err := s.migrate(context.Background()); err != nil {
		db.Close()
		return nil, fmt.Errorf("open store %s: %w", abs, err)
	}
	return s, nil
}

func (s *Store) Path() string { return s.path }

func (s *Store) Close() error { return s.db.Close() }

var migrations = []string{
	`create table runs (
		id          integer primary key,
		project     text    not null,
		stream      text    not null,
		profile     text    not null,
		revision    text    not null,
		change      integer,
		kind        text    not null,
		clean       integer not null,
		review_id   text,
		author      text,
		ci_url      text,
		created_at  integer not null,
		tool        text    not null,
		manifest    blob    not null,
		meta        text    not null
	);
	create index runs_by_change on runs (project, stream, profile, change);
	create index runs_by_revision on runs (project, profile, revision);
	create index runs_by_created on runs (created_at);
	create table blobs (
		hash  blob    primary key,
		kind  integer not null,
		data  blob    not null
	) without rowid;`,
}

func (s *Store) migrate(ctx context.Context) error {
	var version int
	if err := s.db.QueryRowContext(ctx, "pragma user_version").Scan(&version); err != nil {
		return err
	}
	if version > len(migrations) {
		return fmt.Errorf("store schema %d is newer than this nanovision (%d); update nanovision", version, len(migrations))
	}
	for i := version; i < len(migrations); i++ {
		tx, err := s.db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, migrations[i]); err != nil {
			tx.Rollback()
			return fmt.Errorf("migration %d: %w", i+1, err)
		}
		if _, err := tx.ExecContext(ctx, fmt.Sprintf("pragma user_version = %d", i+1)); err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}

// Backup writes a consistent copy of the store to path, which must not exist.
func (s *Store) Backup(ctx context.Context, path string) error {
	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	if _, err := os.Stat(abs); err == nil {
		return fmt.Errorf("backup %s already exists", abs)
	}
	if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx, "vacuum into ?", abs)
	return err
}
