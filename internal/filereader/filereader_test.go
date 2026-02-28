package filereader

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadLinesInFile_VeryLongLine(t *testing.T) {
	// Create a string that explicitly exceeds bufio.MaxScanTokenSize (64KB)
	longLine := strings.Repeat("a", 100*1024) // 100KB string
	content := "short line\n" + longLine + "\nanother short line\n"

	// Create a temporary file
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "longline.txt")
	err := os.WriteFile(filePath, []byte(content), 0o644)
	require.NoError(t, err)

	// Test ReadLinesInFile
	lines, err := ReadLinesInFile(filePath)
	require.NoError(t, err, "Should not return 'token too long' error")
	assert.Len(t, lines, 3)
	assert.Equal(t, "short line", lines[0])
	assert.Equal(t, longLine, lines[1])
	assert.Equal(t, "another short line", lines[2])

	// Test CountLinesInFile
	count, err := CountLinesInFile(filePath)
	require.NoError(t, err)
	assert.Equal(t, 3, count)
}

func TestReadLinesInFile_ExceedsMaxLineLength(t *testing.T) {
	// Create a single line that exceeds our custom maxLineLength (50MB) + 1 byte
	// strings.Repeat with 50MB takes ~50MB RAM, safe for typical environments
	longLine := strings.Repeat("a", 50*1024*1024+10)

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "toolong.txt")
	err := os.WriteFile(filePath, []byte(longLine), 0o644)
	require.NoError(t, err)

	// Ensure the scanner properly returns the "token too long" error when exceeding 50MB
	_, err = ReadLinesInFile(filePath)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "line exceeds maximum allowed length of 50MB")

	_, err = CountLinesInFile(filePath)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "line exceeds maximum allowed length of 50MB")
}

func TestReadLinesInFile_DirectoryLog(t *testing.T) {
	tmpDir := t.TempDir()

	// Redirect os.Stderr to capture logs during the test
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Passing a directory path will typically succeed in os.Open but fail during Read/Seek
	// in DetectEncoding, triggering the warning log we want to test.
	_, err := ReadLinesInFile(tmpDir)

	// Expect ReadLinesInFile to ultimately return an error since directory isn't a text file
	require.Error(t, err)

	w.Close()
	os.Stderr = oldStderr

	var buf bytes.Buffer
	_, copyErr := io.Copy(&buf, r)
	require.NoError(t, copyErr)
	output := buf.String()

	// Verify that DetectEncoding logged the warning before the larger function failed
	assert.Contains(t, output, "Warning: could not detect encoding")
	assert.Contains(t, output, "Assuming UTF-8")
}

func TestReadLinesInFile_100MBLine(t *testing.T) {
	// Create a single line of 100MB
	longLine := strings.Repeat("a", 100*1024*1024)

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "100mb.txt")
	err := os.WriteFile(filePath, []byte(longLine), 0o644)
	require.NoError(t, err)

	_, err = ReadLinesInFile(filePath)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "line exceeds maximum allowed length of 50MB")

	_, err = CountLinesInFile(filePath)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "line exceeds maximum allowed length of 50MB")
}
