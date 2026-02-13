package cache

import (
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"io"
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
}

// initializes the cache from a binary file.
func NewManager(cacheDir string) (*Manager, error) {
	fileName := "analysis_v1.bin" // Versioned filename to avoid conflicts in future updates
	path := filepath.Join(cacheDir, fileName)

	m := &Manager{
		filePath: path,
		entries:  make(map[string]CachedData),
	}

	if err := m.load(); err != nil {
		// If load fails (e.g., file doesn't exist or corruption), start with empty cache
		return m, nil
	}
	return m, nil
}

// reads the GOB encoded file from disk.
func (m *Manager) load() error {
	f, err := os.Open(m.filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Use a buffered reader for performance
	decoder := gob.NewDecoder(f)
	if err := decoder.Decode(&m.entries); err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}
	return nil
}

// persists the current cache state to disk using GOB encoding.
func (m *Manager) Save() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.dirty {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(m.filePath), 0755); err != nil {
		return err
	}

	// Write to a temporary file first to avoid corruption on crash
	tmpPath := m.filePath + ".tmp"
	f, err := os.Create(tmpPath)
	if err != nil {
		return err
	}

	// ENCODE
	encoder := gob.NewEncoder(f)
	if err := encoder.Encode(m.entries); err != nil {
		f.Close() // Close if encoding fails
		return err
	}

	// Close the file explicitly BEFORE calling os.Rename
	if err := f.Close(); err != nil {
		return err
	}

	return os.Rename(tmpPath, m.filePath)
}

// retrieves analysis data if the content hash exists.
func (m *Manager) Get(content []byte) (CachedData, bool) {
	hash := computeHash(content)

	m.mu.RLock()
	defer m.mu.RUnlock()

	val, ok := m.entries[hash]
	return val, ok
}

// stores analysis data for the given content.
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
