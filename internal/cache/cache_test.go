package cache

import (
	"path/filepath"
	"sync"
	"testing"

	"github.com/IgorBayerl/nanovision/internal/analyzer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewManager_Initialization(t *testing.T) {
	tmpDir := t.TempDir()

	// Test initializing in a non-existent subdirectory (should succeed lazily)
	cacheDir := filepath.Join(tmpDir, "nested", "cache")
	m, err := NewManager(cacheDir)

	require.NoError(t, err)
	require.NotNil(t, m)
}

func TestManager_PutAndGet(t *testing.T) {
	tmpDir := t.TempDir()
	m, err := NewManager(tmpDir)
	require.NoError(t, err)

	// Setup dummy data
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

	// Verify Cache Miss on new content
	_, hit := m.Get(content)
	assert.False(t, hit, "Expected cache miss for new content")

	// Put data into cache
	m.Put(content, dataToCache)

	// Verify Cache Hit
	cached, hit := m.Get(content)
	assert.True(t, hit, "Expected cache hit after Put")
	assert.Equal(t, dataToCache, cached, "Cached data should match original")

	// Verify Cache Miss on different content
	otherContent := []byte("function other() {}")
	_, hit = m.Get(otherContent)
	assert.False(t, hit, "Expected cache miss for different content")
}

func TestManager_Persistence(t *testing.T) {
	tmpDir := t.TempDir()

	contentA := []byte("file A content")
	dataA := CachedData{TotalLines: 100}

	contentB := []byte("file B content")
	dataB := CachedData{TotalLines: 200}

	// Create manager, populate, and save
	func() {
		m1, err := NewManager(tmpDir)
		require.NoError(t, err)

		m1.Put(contentA, dataA)
		m1.Put(contentB, dataB)

		err = m1.Save()
		require.NoError(t, err, "Save should succeed")
	}()

	// Verify file was created
	expectedFile := filepath.Join(tmpDir, "analysis_v1.bin")
	assert.FileExists(t, expectedFile)

	// Create NEW manager and verify it loads the data
	func() {
		m2, err := NewManager(tmpDir)
		require.NoError(t, err)

		// Check Data A
		gotA, hit := m2.Get(contentA)
		require.True(t, hit, "Should find data A loaded from disk")
		assert.Equal(t, dataA, gotA)

		// Check Data B
		gotB, hit := m2.Get(contentB)
		require.True(t, hit, "Should find data B loaded from disk")
		assert.Equal(t, dataB, gotB)

		// Check Unknown
		_, hit = m2.Get([]byte("unknown"))
		assert.False(t, hit)
	}()
}

func TestManager_Concurrency(t *testing.T) {
	// This test ensures that the map and mutex usage is correct
	// and doesn't panic under load.
	tmpDir := t.TempDir()
	m, _ := NewManager(tmpDir)

	var wg sync.WaitGroup
	workers := 50
	iterations := 100

	wg.Add(workers)
	for w := 0; w < workers; w++ {
		go func(workerID int) {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				// Create simulated content based on iteration
				// This ensures some overlap between workers (same keys)
				// and some unique keys.
				char := byte((workerID + i) % 10) // 10 unique content variations
				content := []byte{char, char, char}

				// Write
				m.Put(content, CachedData{TotalLines: int(char)})

				// Read
				_, _ = m.Get(content)
			}
		}(w)
	}

	wg.Wait()
}

func TestManager_SaveWhenClean(t *testing.T) {
	tmpDir := t.TempDir()
	m, err := NewManager(tmpDir)
	require.NoError(t, err)

	// Save without putting anything
	err = m.Save()
	require.NoError(t, err)

	// File should NOT exist because we return early if !dirty
	// (Assuming implementation adheres to the `dirty` flag logic proposed)
	expectedFile := filepath.Join(tmpDir, "analysis_v1.bin")
	assert.NoFileExists(t, expectedFile, "Should not write file if cache is empty/clean")
}
