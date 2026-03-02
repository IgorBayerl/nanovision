package cache

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/IgorBayerl/nanovision/internal/analyzer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDetermineCacheDir_SuccessFirst(t *testing.T) {
	logger := slog.Default()
	originalFn := ensureWritableFn
	t.Cleanup(func() { ensureWritableFn = originalFn })

	ensureWritableFn = func(dir string) error {
		return nil // Everything is writable
	}

	dir, err := DetermineCacheDir(logger)
	require.NoError(t, err)
	assert.NotEmpty(t, dir)
}

func TestDetermineCacheDir_Fallback(t *testing.T) {
	logger := slog.Default()
	originalFn := ensureWritableFn
	t.Cleanup(func() { ensureWritableFn = originalFn })

	callCount := 0
	ensureWritableFn = func(dir string) error {
		callCount++
		// Simulate first two paths failing, third succeeding
		if callCount <= 2 {
			return errors.New("simulated permission denied")
		}
		return nil
	}

	dir, err := DetermineCacheDir(logger)
	require.NoError(t, err)
	assert.Equal(t, ".nanovision_cache", dir)
}

func TestDetermineCacheDir_AllFail(t *testing.T) {
	logger := slog.Default()
	originalFn := ensureWritableFn
	t.Cleanup(func() { ensureWritableFn = originalFn })

	ensureWritableFn = func(dir string) error {
		return errors.New("simulated permission denied")
	}

	dir, err := DetermineCacheDir(logger)
	require.Error(t, err)
	assert.Empty(t, dir)
}

func TestNewManager_Initialization(t *testing.T) {
	tmpDir := t.TempDir()
	logger := slog.Default()

	// Test initializing in a non-existent subdirectory (should create it)
	cacheDir := filepath.Join(tmpDir, "nested", "cache")
	m, err := NewManager(cacheDir, logger, nil)

	require.NoError(t, err)
	require.NotNil(t, m)

	// Verify the directory was actually created
	_, statErr := os.Stat(cacheDir)
	assert.NoError(t, statErr, "Cache directory should have been created")
}

func TestNewManager_UnwritableDir(t *testing.T) {
	logger := slog.Default()
	tmpDir := t.TempDir()

	originalFn := ensureWritableFn
	t.Cleanup(func() { ensureWritableFn = originalFn })

	ensureWritableFn = func(dir string) error {
		return errors.New("simulated permission denied")
	}

	m, err := NewManager(tmpDir, logger, nil)

	require.Error(t, err, "Manager should return error if directory is unwritable")
	assert.Nil(t, m)
}

func TestNewManager_StartsFreshOnMissingFile(t *testing.T) {
	tmpDir := t.TempDir()
	logger := slog.Default()

	// Point to a dir with no existing cache file
	m, err := NewManager(tmpDir, logger, nil)

	require.NoError(t, err)
	require.NotNil(t, m)
	assert.Empty(t, m.entries, "Should start with empty cache when no file exists")
}

func TestNewManager_StartsFreshOnCorruptFile(t *testing.T) {
	tmpDir := t.TempDir()
	logger := slog.Default()

	// Write garbage into the cache file to simulate corruption
	corruptPath := filepath.Join(tmpDir, "analysis_v1.bin")
	require.NoError(t, os.WriteFile(corruptPath, []byte("not valid gob data!!!"), 0644))

	m, err := NewManager(tmpDir, logger, nil)

	// Should recover gracefully with an empty cache
	require.NoError(t, err)
	require.NotNil(t, m)
	assert.Empty(t, m.entries, "Should start with empty cache on corrupt file")
}

func TestManager_PutAndGet(t *testing.T) {
	tmpDir := t.TempDir()
	logger := slog.Default()
	m, err := NewManager(tmpDir, logger, nil)
	require.NoError(t, err)

	content := []byte("function main() { return 0; }")
	complexity := 5
	analysisResult := analyzer.AnalysisResult{
		Functions: []analyzer.FunctionMetric{
			{
				Name: "main",
				Position: analyzer.Position{
					StartLine: 1,
					EndLine:   1,
				},
				CyclomaticComplexity: &complexity,
			},
		},
	}
	dataToCache := CachedData{
		TotalLines: 15,
		Result:     analysisResult,
	}

	// Verify cache miss on new content
	_, hit := m.Get(content)
	assert.False(t, hit, "Expected cache miss for new content")

	// Put data into cache
	m.Put(content, dataToCache)

	// Verify cache hit
	cached, hit := m.Get(content)
	assert.True(t, hit, "Expected cache hit after Put")
	assert.Equal(t, dataToCache, cached, "Cached data should match original")

	// Verify cache miss on different content
	_, hit = m.Get([]byte("function other() {}"))
	assert.False(t, hit, "Expected cache miss for different content")
}

func TestManager_Persistence(t *testing.T) {
	tmpDir := t.TempDir()
	logger := slog.Default()

	contentA := []byte("file A content")
	dataA := CachedData{TotalLines: 100}

	contentB := []byte("file B content")
	dataB := CachedData{TotalLines: 200}

	// Create manager, populate, and save
	func() {
		m1, err := NewManager(tmpDir, logger, nil)
		require.NoError(t, err)

		m1.Put(contentA, dataA)
		m1.Put(contentB, dataB)

		err = m1.Save()
		require.NoError(t, err, "Save should succeed")
	}()

	// Verify file was created
	expectedFile := filepath.Join(tmpDir, "analysis_v1.bin")
	assert.FileExists(t, expectedFile)

	// Create a NEW manager and verify it loads the persisted data
	func() {
		m2, err := NewManager(tmpDir, logger, nil)
		require.NoError(t, err)

		gotA, hit := m2.Get(contentA)
		require.True(t, hit, "Should find data A loaded from disk")
		assert.Equal(t, dataA, gotA)

		gotB, hit := m2.Get(contentB)
		require.True(t, hit, "Should find data B loaded from disk")
		assert.Equal(t, dataB, gotB)

		_, hit = m2.Get([]byte("unknown"))
		assert.False(t, hit)
	}()
}

func TestManager_SaveClearsDirtyFlag(t *testing.T) {
	tmpDir := t.TempDir()
	logger := slog.Default()

	m, err := NewManager(tmpDir, logger, nil)
	require.NoError(t, err)

	m.Put([]byte("some content"), CachedData{TotalLines: 42})
	assert.True(t, m.dirty, "Cache should be dirty after Put")

	require.NoError(t, m.Save())
	assert.False(t, m.dirty, "Cache should be clean after Save")

	// A second Save should be a no-op (no file write needed)
	// We verify this doesn't error and the file still exists
	require.NoError(t, m.Save())
	assert.FileExists(t, filepath.Join(tmpDir, "analysis_v1.bin"))
}

func TestManager_SaveWhenClean(t *testing.T) {
	tmpDir := t.TempDir()
	logger := slog.Default()

	m, err := NewManager(tmpDir, logger, nil)
	require.NoError(t, err)

	// Save without putting anything — dirty flag is false
	err = m.Save()
	require.NoError(t, err)

	// File should NOT exist because we skip writing when not dirty
	expectedFile := filepath.Join(tmpDir, "analysis_v1.bin")
	assert.NoFileExists(t, expectedFile, "Should not write file if cache is clean")
}

func TestManager_Concurrency(t *testing.T) {
	tmpDir := t.TempDir()
	logger := slog.Default()
	m, err := NewManager(tmpDir, logger, nil)
	require.NoError(t, err)

	var wg sync.WaitGroup
	workers := 50
	iterations := 100

	wg.Add(workers)
	for w := 0; w < workers; w++ {
		go func(workerID int) {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				char := byte((workerID + i) % 10) // 10 unique content variations — some overlap intentional
				content := []byte{char, char, char}

				m.Put(content, CachedData{TotalLines: int(char)})
				_, _ = m.Get(content)
			}
		}(w)
	}

	wg.Wait()
}

func TestManager_ConcurrentSave(t *testing.T) {
	// Ensure Save is safe to call from multiple goroutines simultaneously.
	tmpDir := t.TempDir()
	logger := slog.Default()

	m, err := NewManager(tmpDir, logger, nil)
	require.NoError(t, err)

	m.Put([]byte("data"), CachedData{TotalLines: 1})

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = m.Save()
		}()
	}
	wg.Wait()

	// After all concurrent saves, the file should exist and be loadable
	m2, err := NewManager(tmpDir, logger, nil)
	require.NoError(t, err)
	_, hit := m2.Get([]byte("data"))
	assert.True(t, hit, "Data should have been persisted by one of the concurrent saves")
}
