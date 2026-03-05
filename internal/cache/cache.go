package cache

import (
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	"github.com/IgorBayerl/nanovision/internal/analyzer"
)

// CachedData represents the static analysis result for a specific file content.
// We store this in a binary format.
type CachedData struct {
	TotalLines int
	Result     analyzer.AnalysisResult
	Metadata   CacheMetadata
}

type Manager struct {
	mu       sync.RWMutex
	filePath string
	// entries maps ContentHash -> Analysis Data
	entries   map[string]CachedData
	validator CacheValidator
	dirty     bool
	logger    *slog.Logger
}

// ensureWritableFn is a variable so tests can replace it to simulate failures.
var ensureWritableFn = ensureWritable

// DetermineCacheDir tries 3 different locations for the cache and returns the first writable one.
func DetermineCacheDir(logger *slog.Logger) (string, error) {
	var candidates []string

	// Try the user's standard cache dir (e.g. ~/.cache/nanovision)
	if userCache, err := os.UserCacheDir(); err == nil {
		candidates = append(candidates, filepath.Join(userCache, "nanovision"))
	}

	// Fallback to system temp directory
	candidates = append(candidates, filepath.Join(os.TempDir(), "nanovision_cache"))

	// Last resort: current working directory
	candidates = append(candidates, ".nanovision_cache")

	for _, dir := range candidates {
		if err := ensureWritableFn(dir); err == nil {
			return dir, nil
		} else if logger != nil {
			logger.Debug("Candidate cache directory is not writable", "dir", dir, "error", err)
		}
	}

	return "", errors.New("no writable cache directory found among candidates")
}

// NewManager initializes the cache from a binary file in the given directory.
func NewManager(cacheDir string, logger *slog.Logger, validator CacheValidator) (*Manager, error) {
	if err := ensureWritableFn(cacheDir); err != nil {
		return nil, fmt.Errorf("cache directory %s is not writable: %w", cacheDir, err)
	}

	fileName := "analysis_v1.bin" // Versioned filename to avoid conflicts in future updates
	path := filepath.Join(cacheDir, fileName)

	m := &Manager{
		filePath:  path,
		entries:   make(map[string]CachedData),
		validator: validator,
		logger:    logger,
	}

	if err := m.load(); err != nil {
		if os.IsNotExist(err) {
			m.logger.Debug("No existing cache file found, starting with empty cache", "path", path)
		} else {
			m.logger.Warn("Failed to load existing cache (file may be corrupted), starting fresh",
				"path", path,
				"error", err,
			)
		}
		return m, nil
	}

	m.logger.Debug("Cache loaded successfully", "entries", len(m.entries), "path", path)
	return m, nil
}

// ensureWritable checks if a directory exists and is writable.
// If it doesn't exist, it tries to create it.
func ensureWritable(dir string) error {
	// Try to create the directory if it doesn't exist
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Try to write a temporary file to verify actual write permissions
	// (MkdirAll might succeed if dir exists, but we might not have write access to files inside)
	testFile := filepath.Join(dir, ".permcheck")
	f, err := os.Create(testFile)
	if err != nil {
		return err
	}
	f.Close()
	return os.Remove(testFile)
}

// load reads the GOB encoded file from disk.
func (m *Manager) load() error {
	f, err := os.Open(m.filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := gob.NewDecoder(f)
	if err := decoder.Decode(&m.entries); err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}
	return nil
}

// Save persists the current cache state to disk using GOB encoding.
func (m *Manager) Save() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.dirty {
		m.logger.Debug("Cache is clean, skipping save")
		return nil
	}

	m.logger.Debug("Persisting cache to disk", "path", m.filePath, "entries", len(m.entries))

	if err := os.MkdirAll(filepath.Dir(m.filePath), 0755); err != nil {
		return err
	}

	// Write to a temporary file first to avoid corruption on crash
	tmpPath := m.filePath + ".tmp"
	f, err := os.Create(tmpPath)
	if err != nil {
		return err
	}

	encoder := gob.NewEncoder(f)
	if err := encoder.Encode(m.entries); err != nil {
		f.Close()
		return err
	}

	// Close the file explicitly BEFORE calling os.Rename
	if err := f.Close(); err != nil {
		return err
	}

	if err := os.Rename(tmpPath, m.filePath); err != nil {
		return err
	}

	m.dirty = false
	return nil
}

// Get retrieves analysis data if the content hash exists.
func (m *Manager) Get(content []byte) (CachedData, bool) {
	hash := computeHash(content)

	m.mu.RLock()
	defer m.mu.RUnlock()

	val, ok := m.entries[hash]
	if !ok {
		return CachedData{}, false
	}

	// Validate metadata
	if m.validator != nil && !m.validator.IsValid(val.Metadata, BuildMetadata{}) {
		m.logger.Debug("Cache entry invalid due to metadata mismatch", "hash", hash)
		return CachedData{}, false
	}

	return val, true
}

// Put stores analysis data for the given content.
func (m *Manager) Put(content []byte, data CachedData) {
	hash := computeHash(content)

	m.mu.Lock()
	defer m.mu.Unlock()

	if data.Metadata.CommitHash == "" && m.logger != nil {
		m.logger.Warn("Putting cache entry without commit metadata")
	}

	m.entries[hash] = data
	m.dirty = true
}

func computeHash(content []byte) string {
	h := sha256.Sum256(content)
	return hex.EncodeToString(h[:])
}
