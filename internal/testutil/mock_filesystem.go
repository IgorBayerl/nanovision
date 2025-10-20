package testutil

import (
	"bufio"
	"io"
	"io/fs"
	"os"
	"path" // Use the generic 'path' package for all internal logic
	"strings"
	"time"

	"github.com/IgorBayerl/AdlerCov/filesystem"
)

// MockFilesystem implements filesystem.Filesystem and filereader.Reader for testing.
type MockFilesystem struct {
	filesystem.Platformer
	files    map[string]string // keys are always Unix-style ('/')
	dirs     map[string][]fs.DirEntry
	cwd      string // always Unix-style ('/')
	platform string
}

func NewMockFilesystem(platform string) *MockFilesystem {
	cwd := "/home/test"
	if platform == "windows" {
		cwd = "C:/Users/Test"
	}
	mfs := &MockFilesystem{
		files:    make(map[string]string),
		dirs:     make(map[string][]fs.DirEntry),
		cwd:      cwd,
		platform: platform,
	}
	mfs.AddDir(cwd)
	return mfs
}

func (m *MockFilesystem) Platform() string {
	return m.platform
}

// fromPlatform normalizes an incoming path to the mock's internal Unix-style representation.
func (m *MockFilesystem) fromPlatform(p string) string {
	return strings.ReplaceAll(p, "\\", "/")
}

// toPlatform converts an internal Unix-style path to the simulated platform's format.
func (m *MockFilesystem) toPlatform(p string) string {
	if m.platform == "windows" {
		return strings.ReplaceAll(p, "/", "\\")
	}
	return p
}

// isAbs robustly checks if a path is absolute for the mock's simulated platform.
func (m *MockFilesystem) isAbs(p string) bool {
	p = m.fromPlatform(p) // ALWAYS normalize before checking
	if m.platform == "windows" {
		// A windows path is absolute if it's like "C:/..."
		return len(p) > 2 && p[1] == ':' && p[2] == '/'
	}
	// A unix path is absolute if it starts with "/"
	return strings.HasPrefix(p, "/")
}

// Abs correctly calculates the absolute path within the mock filesystem.
func (m *MockFilesystem) Abs(p string) (string, error) {
	if m.isAbs(p) {
		// Clean the path to handle things like "C:/dev/../dev/project"
		cleaned := path.Clean(m.fromPlatform(p))
		return m.toPlatform(cleaned), nil
	}
	// Path is relative, join with CWD
	joined := path.Join(m.cwd, m.fromPlatform(p))
	return m.toPlatform(joined), nil
}

func (m *MockFilesystem) AddFile(p string, content string) {
	abs, _ := m.Abs(p)
	key := m.fromPlatform(abs)
	m.files[key] = content

	m.AddDir(path.Dir(key))

	dirKey := path.Dir(key)
	entry := &mockFileInfo{
		name:  path.Base(key),
		isDir: false,
	}

	// Avoid adding duplicate entries
	found := false
	for _, e := range m.dirs[dirKey] {
		if e.Name() == entry.Name() {
			found = true
			break
		}
	}
	if !found {
		m.dirs[dirKey] = append(m.dirs[dirKey], entry)
	}
}

func (m *MockFilesystem) AddDir(p string) {
	abs, _ := m.Abs(p)
	key := m.fromPlatform(abs)

	if _, exists := m.dirs[key]; exists {
		return
	}

	m.dirs[key] = []fs.DirEntry{}
	parentKey := path.Dir(key)

	if parentKey != key {
		m.AddDir(parentKey)
		entry := &mockFileInfo{
			name:  path.Base(key),
			isDir: true,
		}
		// Avoid adding duplicate entries
		found := false
		for _, e := range m.dirs[parentKey] {
			if e.Name() == entry.Name() {
				found = true
				break
			}
		}
		if !found {
			m.dirs[parentKey] = append(m.dirs[parentKey], entry)
		}
	}
}

// --- filesystem.Filesystem AND filereader.Reader implementation ---

func (m *MockFilesystem) Stat(name string) (fs.FileInfo, error) {
	abs, _ := m.Abs(name)
	key := m.fromPlatform(abs)
	if _, exists := m.files[key]; exists {
		return &mockFileInfo{name: path.Base(key), isDir: false}, nil
	}
	if _, exists := m.dirs[key]; exists {
		return &mockFileInfo{name: path.Base(key), isDir: true}, nil
	}
	return nil, os.ErrNotExist
}

// ReadDir is required for both filesystem.Filesystem and filereader.Reader
func (m *MockFilesystem) ReadDir(name string) ([]fs.DirEntry, error) {
	abs, _ := m.Abs(name)
	key := m.fromPlatform(abs)
	if entries, exists := m.dirs[key]; exists {
		return entries, nil
	}
	return nil, os.ErrNotExist
}

func (m *MockFilesystem) Getwd() (string, error) {
	return m.toPlatform(m.cwd), nil
}

// --- filereader.Reader specific implementation ---

func (m *MockFilesystem) ReadFile(p string) ([]string, error) {
	abs, _ := m.Abs(p)
	key := m.fromPlatform(abs)
	content, ok := m.files[key]
	if !ok {
		return nil, os.ErrNotExist
	}

	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func (m *MockFilesystem) CountLines(p string) (int, error) {
	lines, err := m.ReadFile(p)
	if err != nil {
		return 0, err
	}
	return len(lines), nil
}

// --- Unused interface methods ---
func (m *MockFilesystem) MkdirAll(path string, perm fs.FileMode) error               { return nil }
func (m *MockFilesystem) Create(path string) (io.WriteCloser, error)                 { return nil, nil }
func (m *MockFilesystem) Open(path string) (fs.File, error)                          { return nil, nil }
func (m *MockFilesystem) WriteFile(path string, data []byte, perm fs.FileMode) error { return nil }

// mockFileInfo implements fs.FileInfo and fs.DirEntry
type mockFileInfo struct {
	name  string
	isDir bool
}

func (m *mockFileInfo) Name() string               { return m.name }
func (m *mockFileInfo) Size() int64                { return 0 }
func (m *mockFileInfo) Mode() fs.FileMode          { return 0 }
func (m *mockFileInfo) ModTime() time.Time         { return time.Now() }
func (m *mockFileInfo) IsDir() bool                { return m.isDir }
func (m *mockFileInfo) Sys() interface{}           { return nil }
func (m *mockFileInfo) Type() fs.FileMode          { return m.Mode() }
func (m *mockFileInfo) Info() (fs.FileInfo, error) { return m, nil }
