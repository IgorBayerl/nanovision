package cache

import (
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"errors"
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
}

type Manager struct {
	mu       sync.RWMutex
	filePath string
	// entries maps ContentHash -> Analysis Data
	entries map[string]CachedData
	dirty   bool
	logger  *slog.Logger
}

// ensureWritableFn is a variable so tests can replace it to simulate failures.
var ensureWritableFn = ensureWritable

// NewManager initializes the cache from a binary file.
// It performs a write-permission check on cacheDir. If the directory is not writable,
// it falls back to the system temp directory to ensure the application doesn't crash.
func NewManager(cacheDir string, logger *slog.Logger) (*Manager, error) {
	finalDir := cacheDir

	// Check if the requested cache directory is writable
	if err := ensureWritableFn(finalDir); err != nil {
		logger.Warn("Default cache directory is not writable or accessible",
			"dir", finalDir,
			"error", err,
		)

		// Fallback to system temp directory
		fallbackDir := filepath.Join(os.TempDir(), "nanovision_cache")
		logger.Info("Attempting to use fallback cache directory", "dir", fallbackDir)

		if err := ensureWritableFn(fallbackDir); err != nil {
			return nil, errors.New("failed to initialize cache: neither default nor fallback directories are writable")
		}
		finalDir = fallbackDir
	}

	fileName := "analysis_v1.bin" // Versioned filename to avoid conflicts in future updates
	path := filepath.Join(finalDir, fileName)

	m := &Manager{
		filePath: path,
		entries:  make(map[string]CachedData),
		logger:   logger,
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
	return val, ok
}

// Put stores analysis data for the given content.
func (m *Manager) Put(content []byte, data CachedData) {
	hash := computeHash(content)

	m.mu.Lock()
	defer m.mu.Unlock()

	m.entries[hash] = data
	m.dirty = true
}

func computeHash(content []byte) string {
	h := sha256.Sum256(content)
	return hex.EncodeToString(h[:])
}
