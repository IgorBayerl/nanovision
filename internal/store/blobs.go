package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/IgorBayerl/nanovision/internal/store/blob"
)

// Blob is one encoded blob and its key.
type Blob struct {
	Hash blob.Hash
	Data []byte
}

// Verify checks that data is a well-formed blob of an uploadable kind and that
// hash is its key. Blobs never change, so a verified blob can be trusted forever.
func Verify(hash blob.Hash, data []byte) (blob.Kind, error) {
	kind, err := blob.PeekKind(data)
	if err != nil {
		return 0, err
	}
	if !kind.Valid() {
		return 0, fmt.Errorf("blob %s: kind %s cannot be stored", hash, kind)
	}

	switch kind {
	case blob.KindAnalysis:
		_, err = blob.DecodeAnalysis(data)
	case blob.KindCoverage:
		_, err = blob.DecodeCoverage(data)
	case blob.KindManifest:
		_, err = blob.DecodeManifest(data)
	case blob.KindDiff:
		_, err = blob.DecodeDiff(data)
	}
	if err != nil {
		return 0, fmt.Errorf("blob %s: %w", hash, err)
	}
	if key := blob.Sum(data); key != hash {
		return 0, fmt.Errorf("blob %s: content hashes to %s", hash, key)
	}
	return kind, nil
}

// PutBlobs verifies and stores blobs in one transaction. Blobs the store
// already has are skipped.
func (s *Store) PutBlobs(ctx context.Context, blobs []Blob) error {
	kinds := make([]blob.Kind, len(blobs))
	for i, b := range blobs {
		kind, err := Verify(b.Hash, b.Data)
		if err != nil {
			return err
		}
		kinds[i] = kind
	}

	s.writeMu.Lock()
	defer s.writeMu.Unlock()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "insert or ignore into blobs (hash, kind, data) values (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for i, b := range blobs {
		if _, err := stmt.ExecContext(ctx, b.Hash[:], int(kinds[i]), b.Data); err != nil {
			return fmt.Errorf("store blob %s: %w", b.Hash, err)
		}
	}
	return tx.Commit()
}

// Missing returns the hashes the store does not have, in the given order.
func (s *Store) Missing(ctx context.Context, hashes []blob.Hash) ([]blob.Hash, error) {
	have := make(map[blob.Hash]bool, len(hashes))
	err := forEachChunk(hashes, func(chunk []blob.Hash) error {
		rows, err := s.db.QueryContext(ctx, "select hash from blobs where hash in ("+placeholders(len(chunk))+")", hashArgs(chunk)...)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var raw []byte
			if err := rows.Scan(&raw); err != nil {
				return err
			}
			var h blob.Hash
			copy(h[:], raw)
			have[h] = true
		}
		return rows.Err()
	})
	if err != nil {
		return nil, err
	}

	var missing []blob.Hash
	seen := make(map[blob.Hash]bool)
	for _, h := range hashes {
		if !have[h] && !seen[h] {
			missing = append(missing, h)
			seen[h] = true
		}
	}
	return missing, nil
}

// Blob returns an encoded blob.
func (s *Store) Blob(ctx context.Context, hash blob.Hash) ([]byte, error) {
	var data []byte
	err := s.db.QueryRowContext(ctx, "select data from blobs where hash = ?", hash[:]).Scan(&data)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("blob %s: %w", hash, ErrNotFound)
	}
	return data, err
}

// Blobs returns the stored blobs among hashes. Hashes the store does not have
// are left out of the map.
func (s *Store) Blobs(ctx context.Context, hashes []blob.Hash) (map[blob.Hash][]byte, error) {
	out := make(map[blob.Hash][]byte, len(hashes))
	err := forEachChunk(dedupe(hashes), func(chunk []blob.Hash) error {
		rows, err := s.db.QueryContext(ctx, "select hash, data from blobs where hash in ("+placeholders(len(chunk))+")", hashArgs(chunk)...)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var raw, data []byte
			if err := rows.Scan(&raw, &data); err != nil {
				return err
			}
			var h blob.Hash
			copy(h[:], raw)
			out[h] = data
		}
		return rows.Err()
	})
	return out, err
}

// Manifest reads and decodes the manifest of a run.
func (s *Store) Manifest(ctx context.Context, hash blob.Hash) (blob.Manifest, error) {
	data, err := s.Blob(ctx, hash)
	if err != nil {
		return blob.Manifest{}, err
	}
	m, err := blob.DecodeManifest(data)
	if err != nil {
		return blob.Manifest{}, fmt.Errorf("manifest %s: %w", hash, err)
	}
	return m, nil
}

// the SQLite limit on bound parameters is 32766; stay well below it
const chunkSize = 500

func forEachChunk(hashes []blob.Hash, fn func([]blob.Hash) error) error {
	for start := 0; start < len(hashes); start += chunkSize {
		if err := fn(hashes[start:min(start+chunkSize, len(hashes))]); err != nil {
			return err
		}
	}
	return nil
}

func placeholders(n int) string {
	return strings.TrimSuffix(strings.Repeat("?,", n), ",")
}

func hashArgs(hashes []blob.Hash) []any {
	args := make([]any, len(hashes))
	for i := range hashes {
		args[i] = hashes[i][:]
	}
	return args
}

func dedupe(hashes []blob.Hash) []blob.Hash {
	seen := make(map[blob.Hash]bool, len(hashes))
	out := make([]blob.Hash, 0, len(hashes))
	for _, h := range hashes {
		if !h.IsZero() && !seen[h] {
			seen[h] = true
			out = append(out, h)
		}
	}
	return out
}
