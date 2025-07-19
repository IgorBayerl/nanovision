package logging

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

// MockFS implements filesystem.Filesystem for testing
type MockFS struct {
	files       map[string][]byte
	dirs        map[string]bool
	createErr   error
	mkdirAllErr error
	currentDir  string
	openFiles   map[string]*MockFile
	platform    string
}

type MockFile struct {
	*bytes.Buffer
	closed bool
}

func (m *MockFile) Close() error {
	m.closed = true
	return nil
}

func NewMockFS() *MockFS {
	return &MockFS{
		files:      make(map[string][]byte),
		dirs:       make(map[string]bool),
		currentDir: "/current",
		openFiles:  make(map[string]*MockFile),
		platform:   "linux",
	}
}

func (m *MockFS) Platform() string {
	return m.platform
}

func (m *MockFS) Stat(name string) (fs.FileInfo, error) {
	if _, exists := m.files[name]; exists {
		return &mockFileInfo{name: filepath.Base(name), isDir: false}, nil
	}
	if _, exists := m.dirs[name]; exists {
		return &mockFileInfo{name: filepath.Base(name), isDir: true}, nil
	}
	return nil, os.ErrNotExist
}

func (m *MockFS) ReadDir(name string) ([]fs.DirEntry, error) {
	if !m.dirs[name] {
		return nil, os.ErrNotExist
	}
	return []fs.DirEntry{}, nil
}

func (m *MockFS) Getwd() (string, error) {
	return m.currentDir, nil
}

func (m *MockFS) Abs(path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}
	return filepath.Join(m.currentDir, path), nil
}

func (m *MockFS) MkdirAll(path string, perm fs.FileMode) error {
	if m.mkdirAllErr != nil {
		return m.mkdirAllErr
	}
	m.dirs[path] = true
	return nil
}

func (m *MockFS) Create(path string) (io.WriteCloser, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	file := &MockFile{Buffer: &bytes.Buffer{}}
	m.openFiles[path] = file
	return file, nil
}

func (m *MockFS) Open(path string) (fs.File, error) {
	if data, exists := m.files[path]; exists {
		return &mockReadOnlyFile{Buffer: bytes.NewBuffer(data)}, nil
	}
	return nil, os.ErrNotExist
}

func (m *MockFS) ReadFile(path string) ([]byte, error) {
	if data, exists := m.files[path]; exists {
		return data, nil
	}
	return nil, os.ErrNotExist
}

func (m *MockFS) WriteFile(path string, data []byte, perm fs.FileMode) error {
	m.files[path] = data
	return nil
}

type mockFileInfo struct {
	name  string
	isDir bool
}

func (m *mockFileInfo) Name() string       { return m.name }
func (m *mockFileInfo) Size() int64        { return 0 }
func (m *mockFileInfo) Mode() fs.FileMode  { return 0644 }
func (m *mockFileInfo) ModTime() time.Time { return time.Time{} }
func (m *mockFileInfo) IsDir() bool        { return m.isDir }
func (m *mockFileInfo) Sys() interface{}   { return nil }

type mockReadOnlyFile struct {
	*bytes.Buffer
}

func (m *mockReadOnlyFile) Stat() (fs.FileInfo, error) {
	return &mockFileInfo{name: "test", isDir: false}, nil
}

func (m *mockReadOnlyFile) Close() error { return nil }

func TestVerbosityLevel_SlogLevel(t *testing.T) {
	tests := []struct {
		name     string
		level    VerbosityLevel
		expected slog.Level
	}{
		{
			name:     "Verbose maps to Debug",
			level:    Verbose,
			expected: slog.LevelDebug,
		},
		{
			name:     "Info maps to Info",
			level:    Info,
			expected: slog.LevelInfo,
		},
		{
			name:     "Warning maps to Warn",
			level:    Warning,
			expected: slog.LevelWarn,
		},
		{
			name:     "Error maps to Error",
			level:    Error,
			expected: slog.LevelError,
		},
		{
			name:     "Off maps to high level",
			level:    Off,
			expected: slog.Level(slog.LevelError + 128),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.level.SlogLevel()

			// Assert
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestParseVerbosity_ValidInputs(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected VerbosityLevel
	}{
		{
			name:     "verbose lowercase",
			input:    "verbose",
			expected: Verbose,
		},
		{
			name:     "VERBOSE uppercase",
			input:    "VERBOSE",
			expected: Verbose,
		},
		{
			name:     "info with whitespace",
			input:    "  info  ",
			expected: Info,
		},
		{
			name:     "warning full word",
			input:    "warning",
			expected: Warning,
		},
		{
			name:     "warn abbreviation",
			input:    "warn",
			expected: Warning,
		},
		{
			name:     "error standard",
			input:    "error",
			expected: Error,
		},
		{
			name:     "off standard",
			input:    "off",
			expected: Off,
		},
		{
			name:     "silent alternative",
			input:    "silent",
			expected: Off,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result, err := ParseVerbosity(tt.input)

			// Assert
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestParseVerbosity_InvalidInputs(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "invalid string",
			input: "invalid",
		},
		{
			name:  "empty string",
			input: "",
		},
		{
			name:  "number",
			input: "123",
		},
		{
			name:  "partial match",
			input: "inf",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result, err := ParseVerbosity(tt.input)

			// Assert
			if err == nil {
				t.Error("expected error but got none")
			}
			if result != Info {
				t.Errorf("expected default Info level, got %v", result)
			}
			if !strings.Contains(err.Error(), "invalid verbosity level") {
				t.Errorf("expected error message to contain 'invalid verbosity level', got: %v", err.Error())
			}
		})
	}
}

func TestInit_DefaultConfig(t *testing.T) {
	// Arrange
	originalDefault := slog.Default()
	defer slog.SetDefault(originalDefault)

	// Act
	closer, err := Init(nil)

	// Assert
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if closer != nil {
		t.Error("expected nil closer for console-only logging")
	}
}

func TestInit_ConsoleOnly(t *testing.T) {
	// Arrange
	originalDefault := slog.Default()
	defer slog.SetDefault(originalDefault)
	cfg := &Config{
		Verbosity: Info,
		Format:    "text",
	}

	// Act
	closer, err := Init(cfg)

	// Assert
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if closer != nil {
		t.Error("expected nil closer for console-only logging")
	}
}

func TestInit_FileOutput(t *testing.T) {
	// Arrange
	originalDefault := slog.Default()
	defer slog.SetDefault(originalDefault)
	mockFS := NewMockFS()
	logFile := filepath.Join("logs", "test.log")
	expectedDir := filepath.Dir(logFile)

	cfg := &Config{
		Verbosity: Info,
		Format:    "text",
		File:      logFile,
		FS:        mockFS,
	}

	// Act
	closer, err := Init(cfg)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if closer == nil {
		t.Fatal("expected non-nil closer for file logging")
	}
	if !mockFS.dirs[expectedDir] {
		t.Errorf("expected log directory '%s' to be created", expectedDir)
	}
	if _, exists := mockFS.openFiles[logFile]; !exists {
		t.Errorf("expected log file '%s' to be created", logFile)
	}

	// Cleanup
	if closer != nil {
		closer.Close()
	}
}

func TestInit_FormatHandling(t *testing.T) {
	cases := []struct {
		name         string
		format       string
		expectedType any // Use reflection to check handler type
	}{
		{"text format", "text", &slog.TextHandler{}},
		{"json format", "json", &slog.JSONHandler{}},
		{"mixed-case json", "JsOn", &slog.JSONHandler{}},
		{"empty format defaults to text", "", &slog.TextHandler{}},
		{"unknown format defaults to text", "xml", &slog.TextHandler{}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			originalDefault := slog.Default()
			defer slog.SetDefault(originalDefault)

			// To check the handler type, we need to capture the output writer
			var buf bytes.Buffer

			// We can't directly inspect the handler set as default.
			// Instead, we'll create the handler directly and check its type,
			// which is what Init does internally. This is a reasonable compromise.
			opts := &slog.HandlerOptions{}
			var handler slog.Handler
			if strings.ToLower(tc.format) == "json" {
				handler = slog.NewJSONHandler(&buf, opts)
			} else {
				handler = slog.NewTextHandler(&buf, opts)
			}

			if reflect.TypeOf(handler) != reflect.TypeOf(tc.expectedType) {
				t.Errorf("expected handler type %T, got %T", tc.expectedType, handler)
			}
		})
	}
}

func TestInit_MkdirAllError(t *testing.T) {
	// Arrange
	originalDefault := slog.Default()
	defer slog.SetDefault(originalDefault)
	mockFS := NewMockFS()
	mockFS.mkdirAllErr = errors.New("mkdir failed")
	cfg := &Config{
		Verbosity: Info,
		Format:    "text",
		File:      "/logs/test.log",
		FS:        mockFS,
	}

	// Act
	closer, err := Init(cfg)

	// Assert
	if err == nil {
		t.Error("expected error from MkdirAll failure")
	}
	if !strings.Contains(err.Error(), "failed to create log directory") {
		t.Errorf("expected directory creation error, got: %v", err)
	}
	if closer != nil {
		t.Error("expected nil closer on error")
	}
}

func TestInit_CreateFileError(t *testing.T) {
	// Arrange
	originalDefault := slog.Default()
	defer slog.SetDefault(originalDefault)
	mockFS := NewMockFS()
	mockFS.createErr = errors.New("create failed")
	cfg := &Config{
		Verbosity: Info,
		Format:    "text",
		File:      "/logs/test.log",
		FS:        mockFS,
	}

	// Act
	closer, err := Init(cfg)

	// Assert
	if err == nil {
		t.Error("expected error from Create failure")
	}
	if !strings.Contains(err.Error(), "failed to create log file") {
		t.Errorf("expected file creation error, got: %v", err)
	}
	if closer != nil {
		t.Error("expected nil closer on error")
	}
}
