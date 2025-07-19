package glob_test

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/glob"
	assert "github.com/IgorBayerl/ReportGenerator/go_report_generator/testutil/asserts"
)

// Cross-platform execution helper
func forPlatforms(t *testing.T, fn func(t *testing.T, fs *MockFilesystem)) {
	t.Helper()

	t.Run("unix", func(t *testing.T) {
		t.Parallel()
		fn(t, setupLinuxFS())
	})

	t.Run("windows", func(t *testing.T) {
		t.Parallel()
		fn(t, setupWindowsFS())
	})
}

// Mock filesystem
type MockFileInfo struct {
	name    string
	size    int64
	mode    fs.FileMode
	modTime time.Time
	isDir   bool
}

func (m MockFileInfo) Name() string       { return m.name }
func (m MockFileInfo) Size() int64        { return m.size }
func (m MockFileInfo) Mode() fs.FileMode  { return m.mode }
func (m MockFileInfo) ModTime() time.Time { return m.modTime }
func (m MockFileInfo) IsDir() bool        { return m.isDir }
func (m MockFileInfo) Sys() any           { return nil }

type MockDirEntry struct {
	name  string
	isDir bool
	info  MockFileInfo
}

func (m MockDirEntry) Name() string               { return m.name }
func (m MockDirEntry) IsDir() bool                { return m.isDir }
func (m MockDirEntry) Type() fs.FileMode          { return m.info.Mode() }
func (m MockDirEntry) Info() (fs.FileInfo, error) { return m.info, nil }

type MockFilesystem struct {
	files     map[string]MockFileInfo
	dirs      map[string][]MockDirEntry
	cwd       string
	platform  string
	separator string
}

func NewMockFilesystem(platform string) *MockFilesystem {
	sep := "/"
	if platform == "windows" {
		sep = `\`
	}
	return &MockFilesystem{
		files:     map[string]MockFileInfo{},
		dirs:      map[string][]MockDirEntry{},
		cwd:       "/",
		platform:  platform,
		separator: sep,
	}
}

func (m *MockFilesystem) Platform() string { return m.platform }

func (m *MockFilesystem) normalizePath(p string) string {
	if m.platform == "windows" {
		p = strings.ReplaceAll(p, "/", `\`)
		switch {
		case len(p) >= 2 && p[1] == ':':
			return p // already absolute
		case strings.HasPrefix(p, `\\`):
			return p // UNC
		case strings.HasPrefix(p, `\`):
			return `C:` + p // rooted but drive-less
		default:
			return p // relative
		}
	}
	return strings.ReplaceAll(p, `\`, "/")
}

func (m *MockFilesystem) Stat(name string) (fs.FileInfo, error) {
	abs, err := m.Abs(name)
	if err != nil {
		return nil, &fs.PathError{Op: "stat", Path: name, Err: err}
	}
	if info, ok := m.files[abs]; ok {
		return info, nil
	}
	return nil, &fs.PathError{Op: "stat", Path: name, Err: fs.ErrNotExist}
}

func (m *MockFilesystem) ReadDir(name string) ([]fs.DirEntry, error) {
	abs, err := m.Abs(name)
	if err != nil {
		return nil, &fs.PathError{Op: "readdir", Path: name, Err: err}
	}
	entries, ok := m.dirs[abs]
	if !ok {
		return nil, &fs.PathError{Op: "readdir", Path: name, Err: fs.ErrNotExist}
	}
	out := make([]fs.DirEntry, len(entries))
	for i := range entries {
		out[i] = entries[i]
	}
	return out, nil
}

func (m *MockFilesystem) Getwd() (string, error) { return m.cwd, nil }

func (m *MockFilesystem) Abs(path string) (string, error) {
	path = m.normalizePath(path)
	if m.platform == "windows" {
		if filepath.IsAbs(path) || strings.HasPrefix(path, `C:`) {
			return path, nil
		}
		if strings.HasPrefix(path, `\`) {
			return `C:` + path, nil
		}
		return filepath.Join(m.cwd, path), nil
	}
	if filepath.IsAbs(path) {
		return path, nil
	}
	return filepath.Join(m.cwd, path), nil
}

func (m *MockFilesystem) AddFile(path string, isDir bool) {
	abs, _ := m.Abs(path)

	info := MockFileInfo{
		name:    filepath.Base(abs),
		size:    100,
		mode:    0o644,
		modTime: time.Now(),
		isDir:   isDir,
	}
	if isDir {
		info.mode = fs.ModeDir | 0o755
	}
	m.files[abs] = info

	parent := filepath.Dir(abs)
	if parent != abs {
		m.dirs[parent] = append(m.dirs[parent], MockDirEntry{
			name: info.name, isDir: isDir, info: info,
		})
	}
}

func (m *MockFilesystem) SetCwd(cwd string) { m.cwd = m.normalizePath(cwd) }

// Stubbed-out methods not needed in these tests.
func (*MockFilesystem) MkdirAll(string, fs.FileMode) error          { return nil }
func (*MockFilesystem) Create(string) (io.WriteCloser, error)       { return nil, nil }
func (*MockFilesystem) Open(string) (fs.File, error)                { return nil, nil }
func (*MockFilesystem) ReadFile(string) ([]byte, error)             { return nil, nil }
func (*MockFilesystem) WriteFile(string, []byte, fs.FileMode) error { return nil }

// Test helper functions
func setupLinuxFS() *MockFilesystem {
	fs := NewMockFilesystem("unix")
	fs.SetCwd("/home/user")

	// Create directory structure
	fs.AddFile("/", true)
	fs.AddFile("/home", true)
	fs.AddFile("/home/user", true)
	fs.AddFile("/home/user/documents", true)
	fs.AddFile("/home/user/documents/file1.txt", false)
	fs.AddFile("/home/user/documents/file2.txt", false)
	fs.AddFile("/home/user/documents/report.pdf", false)
	fs.AddFile("/home/user/documents/subdir", true)
	fs.AddFile("/home/user/documents/subdir/nested.txt", false)
	fs.AddFile("/home/user/documents/subdir/deep", true)
	fs.AddFile("/home/user/documents/subdir/deep/file.log", false)
	fs.AddFile("/home/user/pictures", true)
	fs.AddFile("/home/user/pictures/photo1.jpg", false)
	fs.AddFile("/home/user/pictures/photo2.png", false)
	fs.AddFile("/tmp", true)
	fs.AddFile("/tmp/temp1.tmp", false)
	fs.AddFile("/tmp/temp2.tmp", false)

	return fs
}

func setupWindowsFS() *MockFilesystem {
	fs := NewMockFilesystem("windows")
	fs.SetCwd("C:\\Users\\User")

	// Create directory structure
	fs.AddFile("C:\\", true)
	fs.AddFile("C:\\Users", true)
	fs.AddFile("C:\\Users\\User", true)
	fs.AddFile("C:\\Users\\User\\Documents", true)
	fs.AddFile("C:\\Users\\User\\Documents\\file1.txt", false)
	fs.AddFile("C:\\Users\\User\\Documents\\file2.txt", false)
	fs.AddFile("C:\\Users\\User\\Documents\\report.pdf", false)
	fs.AddFile("C:\\Users\\User\\Documents\\subdir", true)
	fs.AddFile("C:\\Users\\User\\Documents\\subdir\\nested.txt", false)
	fs.AddFile("C:\\Users\\User\\Documents\\subdir\\deep", true)
	fs.AddFile("C:\\Users\\User\\Documents\\subdir\\deep\\file.log", false)
	fs.AddFile("C:\\Users\\User\\Pictures", true)
	fs.AddFile("C:\\Users\\User\\Pictures\\photo1.jpg", false)
	fs.AddFile("C:\\Users\\User\\Pictures\\photo2.png", false)
	fs.AddFile("C:\\Temp", true)
	fs.AddFile("C:\\Temp\\temp1.tmp", false)
	fs.AddFile("C:\\Temp\\temp2.tmp", false)

	return fs
}

func TestExpandNames_BasicPatterns_ReturnExpected(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		pattern  string
		wantUnix []string
		wantWin  []string
	}{
		{
			name:    "single asterisk",
			pattern: "documents/*.txt",
			wantUnix: []string{
				"/home/user/documents/file1.txt",
				"/home/user/documents/file2.txt",
			},
			wantWin: []string{
				"C:/Users/User/Documents/file1.txt",
				"C:/Users/User/Documents/file2.txt",
			},
		},
		{
			name:    "question mark",
			pattern: "documents/file?.txt",
			wantUnix: []string{
				"/home/user/documents/file1.txt",
				"/home/user/documents/file2.txt",
			},
			wantWin: []string{
				"C:/Users/User/Documents/file1.txt",
				"C:/Users/User/Documents/file2.txt",
			},
		},
		{
			name:    "double asterisk recursive",
			pattern: "documents/**/*.txt",
			wantUnix: []string{
				"/home/user/documents/file1.txt",
				"/home/user/documents/file2.txt",
				"/home/user/documents/subdir/nested.txt",
			},
			wantWin: []string{
				"C:/Users/User/Documents/file1.txt",
				"C:/Users/User/Documents/file2.txt",
				"C:/Users/User/Documents/subdir/nested.txt",
			},
		},
		{
			name:    "character class",
			pattern: "documents/file[12].txt",
			wantUnix: []string{
				"/home/user/documents/file1.txt",
				"/home/user/documents/file2.txt",
			},
			wantWin: []string{
				"C:/Users/User/Documents/file1.txt",
				"C:/Users/User/Documents/file2.txt",
			},
		},
		{
			name:    "brace expansion",
			pattern: "documents/{file1,file2}.txt",
			wantUnix: []string{
				"/home/user/documents/file1.txt",
				"/home/user/documents/file2.txt",
			},
			wantWin: []string{
				"C:/Users/User/Documents/file1.txt",
				"C:/Users/User/Documents/file2.txt",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			forPlatforms(t, func(t *testing.T, fs *MockFilesystem) {
				// Arrange
				g := glob.NewGlob(tc.pattern, fs)

				// Act
				got, err := g.ExpandNames()

				// Assert
				if err != nil {
					t.Fatalf("unexpected err: %v", err)
				}
				want := tc.wantUnix
				if fs.platform == "windows" {
					want = tc.wantWin
				}
				assert.Equal(t, want, got, assert.CmpPaths...)
			})
		})
	}
}

func TestExpandNames_AbsolutePaths_CorrectPerPlatform(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		pattern  string
		wantUnix []string
		wantWin  []string
	}{
		{
			name:    "absolute path unix",
			pattern: "/home/user/documents/*.txt",
			wantUnix: []string{
				"/home/user/documents/file1.txt",
				"/home/user/documents/file2.txt",
			},
			wantWin: []string{}, // Should work differently on Windows
		},
		{
			name:     "absolute path windows",
			pattern:  "C:\\Users\\User\\Documents\\*.txt",
			wantUnix: []string{}, // Should work differently on Unix
			wantWin: []string{
				"C:/Users/User/Documents/file1.txt",
				"C:/Users/User/Documents/file2.txt",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			forPlatforms(t, func(t *testing.T, fs *MockFilesystem) {
				// Arrange
				g := glob.NewGlob(tc.pattern, fs)

				// Act
				got, err := g.ExpandNames()

				// Assert
				if err != nil {
					t.Fatalf("unexpected err: %v", err)
				}
				want := tc.wantUnix
				if fs.platform == "windows" {
					want = tc.wantWin
				}
				assert.Equal(t, want, got, assert.CmpPaths...)
			})
		})
	}
}

func TestExpandNames_RecursivePatterns_ReturnExpected(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		pattern  string
		wantUnix []string
		wantWin  []string
	}{
		{
			name:    "recursive all files",
			pattern: "documents/**",
			wantUnix: []string{
				"/home/user/documents",
				"/home/user/documents/file1.txt",
				"/home/user/documents/file2.txt",
				"/home/user/documents/report.pdf",
				"/home/user/documents/subdir",
				"/home/user/documents/subdir/nested.txt",
				"/home/user/documents/subdir/deep",
				"/home/user/documents/subdir/deep/file.log",
			},
			wantWin: []string{
				"C:/Users/User/Documents",
				"C:/Users/User/Documents/file1.txt",
				"C:/Users/User/Documents/file2.txt",
				"C:/Users/User/Documents/report.pdf",
				"C:/Users/User/Documents/subdir",
				"C:/Users/User/Documents/subdir/nested.txt",
				"C:/Users/User/Documents/subdir/deep",
				"C:/Users/User/Documents/subdir/deep/file.log",
			},
		},
		{
			name:    "recursive specific extension",
			pattern: "**/*.txt",
			wantUnix: []string{
				"/home/user/documents/file1.txt",
				"/home/user/documents/file2.txt",
				"/home/user/documents/subdir/nested.txt",
			},
			wantWin: []string{
				"C:/Users/User/Documents/file1.txt",
				"C:/Users/User/Documents/file2.txt",
				"C:/Users/User/Documents/subdir/nested.txt",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			forPlatforms(t, func(t *testing.T, fs *MockFilesystem) {
				// Arrange
				g := glob.NewGlob(tc.pattern, fs)

				// Act
				got, err := g.ExpandNames()

				// Assert
				if err != nil {
					t.Fatalf("unexpected err: %v", err)
				}
				want := tc.wantUnix
				if fs.platform == "windows" {
					want = tc.wantWin
				}
				assert.Equal(t, want, got, assert.CmpPaths...)
			})
		})
	}
}

func TestExpandNames_InvalidGlob_ReturnsError(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		pattern string
	}{
		{
			name:    "unbalanced braces",
			pattern: "documents/{file1,file2.txt",
		},
		{
			name:    "unterminated character class",
			pattern: "documents/file[12.txt",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			forPlatforms(t, func(t *testing.T, fs *MockFilesystem) {
				// Arrange
				g := glob.NewGlob(tc.pattern, fs)

				// Act
				got, err := g.ExpandNames()

				// Assert
				if err == nil {
					t.Errorf("expected error for malformed pattern %q", tc.pattern)
				}
				if len(got) != 0 {
					t.Errorf("expected empty results for malformed pattern, got %d", len(got))
				}
			})
		})
	}
}

func TestExpandNames_NonGlobInputs_ReturnsPathOrEmpty(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		pattern  string
		wantUnix []string
		wantWin  []string
	}{
		{
			name:     "empty pattern",
			pattern:  "",
			wantUnix: []string{},
			wantWin:  []string{},
		},
		{
			name:    "literal path exists",
			pattern: "documents/file1.txt",
			wantUnix: []string{
				"/home/user/documents/file1.txt",
			},
			wantWin: []string{
				"C:/Users/User/Documents/file1.txt",
			},
		},
		{
			name:     "literal path doesn't exist",
			pattern:  "documents/nonexistent.txt",
			wantUnix: []string{},
			wantWin:  []string{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			forPlatforms(t, func(t *testing.T, fs *MockFilesystem) {
				// Arrange
				g := glob.NewGlob(tc.pattern, fs)

				// Act
				got, err := g.ExpandNames()

				// Assert
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				want := tc.wantUnix
				if fs.platform == "windows" {
					want = tc.wantWin
				}
				assert.Equal(t, want, got, assert.CmpPaths...)
			})
		})
	}
}

func TestExpandNames_FilesystemError_ReturnsEmptySlice(t *testing.T) {
	t.Parallel()

	// Arrange
	fs := &MockFilesystem{
		files:     make(map[string]MockFileInfo),
		dirs:      make(map[string][]MockDirEntry),
		cwd:       "/",
		platform:  "unix",
		separator: "/",
	}
	g := glob.NewGlob("documents/*.txt", fs)

	// Act
	got, err := g.ExpandNames()

	// Assert
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty results, got %d", len(got))
	}
}

func TestExpandNames_LargeTree_FindsAllTxtFiles(t *testing.T) {
	t.Parallel()

	// Arrange
	fs := NewMockFilesystem("unix")
	fs.SetCwd("/home/user")
	fs.AddFile("/", true)
	fs.AddFile("/home", true)
	fs.AddFile("/home/user", true)
	fs.AddFile("/home/user/large", true)

	expectedCount := 0
	for i := 0; i < 100; i++ {
		fs.AddFile(fmt.Sprintf("/home/user/large/file%d.txt", i), false)
		expectedCount++
		if i%10 == 0 {
			fs.AddFile(fmt.Sprintf("/home/user/large/subdir%d", i), true)
			for j := 0; j < 10; j++ {
				fs.AddFile(fmt.Sprintf("/home/user/large/subdir%d/nested%d.txt", i, j), false)
				expectedCount++
			}
		}
	}

	g := glob.NewGlob("large/**/*.txt", fs)

	// Act
	got, err := g.ExpandNames()

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != expectedCount {
		t.Errorf("expected %d results, got %d", expectedCount, len(got))
	}
}

func TestExpandNames_CaseSensitivity_VariousFlags(t *testing.T) {
	t.Parallel()

	addExtraFiles := func(fs *MockFilesystem) {
		if fs.platform == "unix" {
			fs.AddFile("/home/user/documents/File1.TXT", false)
			fs.AddFile("/home/user/documents/FILE2.txt", false)
		} else {
			fs.AddFile(`C:\Users\User\Documents\File1.TXT`, false)
			fs.AddFile(`C:\Users\User\Documents\FILE2.txt`, false)
		}
	}

	cases := []struct {
		name       string
		pattern    string
		ignoreCase bool
		wantUnix   int
		wantWin    int
	}{
		{"mixed-case sensitive", "documents/*.TXT", false, 1, 1},
		{"mixed-case insensitive", "documents/*.TXT", true, 4, 4},
		{"lower-case sensitive", "documents/*.txt", false, 3, 3},
		{"lower-case insensitive", "documents/*.txt", true, 4, 4},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			forPlatforms(t, func(t *testing.T, fs *MockFilesystem) {
				// Arrange
				addExtraFiles(fs)
				g := glob.NewGlob(tc.pattern, fs, glob.WithIgnoreCase(tc.ignoreCase))

				// Act
				got, err := g.ExpandNames()

				// Assert
				if err != nil {
					t.Fatal(err)
				}
				want := tc.wantUnix
				if fs.platform == "windows" {
					want = tc.wantWin
				}
				if len(got) != want {
					t.Fatalf("want %d got %d", want, len(got))
				}
			})
		})
	}
}

func TestExpandNames_ComplexBraceExpansion_ReturnsExpected(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		pattern  string
		wantUnix []string
		wantWin  []string
	}{
		{
			name:    "nested braces",
			pattern: "documents/{file{1,2},report}.{txt,pdf}",
			wantUnix: []string{
				"/home/user/documents/file1.txt",
				"/home/user/documents/file2.txt",
				"/home/user/documents/report.pdf",
			},
			wantWin: []string{
				"C:/Users/User/Documents/file1.txt",
				"C:/Users/User/Documents/file2.txt",
				"C:/Users/User/Documents/report.pdf",
			},
		},
		{
			name:    "cross-separator braces",
			pattern: "{documents,pictures}/*.{txt,jpg}",
			wantUnix: []string{
				"/home/user/documents/file1.txt",
				"/home/user/documents/file2.txt",
				"/home/user/pictures/photo1.jpg",
			},
			wantWin: []string{
				"C:/Users/User/Documents/file1.txt",
				"C:/Users/User/Documents/file2.txt",
				"C:/Users/User/Pictures/photo1.jpg",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			forPlatforms(t, func(t *testing.T, fs *MockFilesystem) {
				// Arrange
				g := glob.NewGlob(tc.pattern, fs)

				// Act
				got, err := g.ExpandNames()

				// Assert
				if err != nil {
					t.Fatalf("unexpected err: %v", err)
				}
				want := tc.wantUnix
				if fs.platform == "windows" {
					want = tc.wantWin
				}
				assert.Equal(t, want, got, assert.CmpPaths...)
			})
		})
	}
}

func TestExpandNames_DotAndDotDot_Patterns(t *testing.T) {
	t.Parallel()

	t.Run("current_directory", func(t *testing.T) {
		t.Parallel()
		forPlatforms(t, func(t *testing.T, fs *MockFilesystem) {
			// Arrange
			g := glob.NewGlob(".", fs)

			// Act
			got, err := g.ExpandNames()

			// Assert
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			expectedCwd := "/home/user"
			if fs.platform == "windows" {
				expectedCwd = "C:/Users/User"
			}
			assert.Equal(t, []string{expectedCwd}, got, assert.CmpPaths...)
		})
	})

	t.Run("parent_directory", func(t *testing.T) {
		t.Parallel()
		forPlatforms(t, func(t *testing.T, fs *MockFilesystem) {
			// Arrange
			g := glob.NewGlob("..", fs)

			// Act
			got, err := g.ExpandNames()

			// Assert
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			expectedParent := "/home"
			if fs.platform == "windows" {
				expectedParent = "C:/Users"
			}
			assert.Equal(t, []string{expectedParent}, got, assert.CmpPaths...)
		})
	})
}

func TestGetFilesPublicAPI(t *testing.T) {
	t.Parallel()

	// Arrange
	dir := t.TempDir()
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	os.Chdir(dir)

	// Act
	results, err := glob.GetFiles("*.nonexistent")

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected empty results for non-matching pattern, got %d results", len(results))
	}

	// Act
	results, err = glob.GetFiles("")

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected empty results for empty pattern, got %d results", len(results))
	}
}
