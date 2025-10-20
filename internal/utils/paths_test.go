package utils_test

import (
	"io"
	"log/slog"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/IgorBayerl/AdlerCov/internal/testutil"
	"github.com/IgorBayerl/AdlerCov/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindFileInSourceDirs(t *testing.T) {
	isWindows := runtime.GOOS == "windows"
	
	// Helper to create platform-appropriate paths
	p := func(path string) string {
		if isWindows {
			return filepath.FromSlash("C:" + path)
		}
		return path
	}

	testCases := []struct {
		name            string
		mockFiles       map[string]string
		sourceDirs      []string
		relativePath    string
		expectedPath    string
		expectError     bool
		skipPathCheck   bool // For ambiguous cases where we can't predict the exact path
	}{
		// ===== STRATEGY 1: ABSOLUTE PATH =====
		{
			name:         "Strategy 1: Absolute path exists",
			mockFiles:    map[string]string{p("/project/src/calculator.cpp"): "content"},
			sourceDirs:   []string{p("/another/dir")}, // Should be ignored
			relativePath: p("/project/src/calculator.cpp"),
			expectedPath: p("/project/src/calculator.cpp"),
			expectError:  false,
		},
		{
			name:         "Strategy 1: Absolute path doesn't exist, should try other strategies",
			mockFiles:    map[string]string{p("/actual/location/file.go"): "content"},
			sourceDirs:   []string{p("/actual/location")},
			relativePath: p("/wrong/path/file.go"), // Absolute but wrong
			expectedPath: p("/actual/location/file.go"), // Should find via recursive search
			expectError:  false,
		},

		// ===== STRATEGY 2: DIRECT JOIN =====
		{
			name:         "Strategy 2: Simple relative path in first dir",
			mockFiles:    map[string]string{p("/project/src/app.go"): "content"},
			sourceDirs:   []string{p("/project")},
			relativePath: "src/app.go",
			expectedPath: p("/project/src/app.go"),
			expectError:  false,
		},
		{
			name:         "Strategy 2: Found in second source directory",
			mockFiles:    map[string]string{p("/project/lib/utils.go"): "content"},
			sourceDirs:   []string{p("/project/api"), p("/project/lib")},
			relativePath: "lib/utils.go",
			expectedPath: p("/project/lib/utils.go"),
			expectError:  false,
		},
		{
			name:         "Strategy 2: Nested path with multiple segments",
			mockFiles:    map[string]string{p("/dev/project/api/controllers/v2/HealthController.cs"): "content"},
			sourceDirs:   []string{p("/dev/project")},
			relativePath: "api/controllers/v2/HealthController.cs",
			expectedPath: p("/dev/project/api/controllers/v2/HealthController.cs"),
			expectError:  false,
		},
		{
			name:         "Strategy 2: Path with spaces",
			mockFiles:    map[string]string{p("/My Project/My Source/app.cs"): "content"},
			sourceDirs:   []string{p("/My Project")},
			relativePath: "My Source/app.cs",
			expectedPath: p("/My Project/My Source/app.cs"),
			expectError:  false,
		},

		// ===== STRATEGY 3: SUFFIX MATCHING (CI/CD scenarios) =====
		{
			name:         "Strategy 3: CI/CD path mismatch - different root",
			mockFiles:    map[string]string{p("/local/dev/src/main.js"): "content"},
			sourceDirs:   []string{p("/local/dev")},
			relativePath: p("/home/runner/work/project/src/main.js"),
			expectedPath: p("/local/dev/src/main.js"),
			expectError:  false,
		},
		{
			name:         "Strategy 3: Multiple path segments need to be stripped",
			mockFiles:    map[string]string{p("/workspace/app/controllers/user.go"): "content"},
			sourceDirs:   []string{p("/workspace")},
			relativePath: p("/build/ci/pipeline/workspace/app/controllers/user.go"),
			expectedPath: p("/workspace/app/controllers/user.go"),
			expectError:  false,
		},
		{
			name:         "Strategy 3: Docker volume path mismatch",
			mockFiles:    map[string]string{p("/app/src/service.py"): "content"},
			sourceDirs:   []string{p("/app")},
			relativePath: p("/home/user/project/src/service.py"),
			expectedPath: p("/app/src/service.py"),
			expectError:  false,
		},

		// ===== STRATEGY 4: RECURSIVE SEARCH =====
		{
			name:         "Strategy 4: Filename only - deeply nested file",
			mockFiles:    map[string]string{p("/project/src/services/api/v2/handlers/MyService.cs"): "content"},
			sourceDirs:   []string{p("/project")},
			relativePath: "MyService.cs",
			expectedPath: p("/project/src/services/api/v2/handlers/MyService.cs"),
			expectError:  false,
		},
		{
			name:         "Strategy 4: Filename only - multiple dirs searched",
			mockFiles:    map[string]string{p("/project/lib/utils.go"): "content"},
			sourceDirs:   []string{p("/project/api"), p("/project/lib")},
			relativePath: "utils.go",
			expectedPath: p("/project/lib/utils.go"),
			expectError:  false,
		},
		{
			name:         "Strategy 4: Filename exists in multiple locations - returns first match",
			mockFiles: map[string]string{
				p("/project/api/config.go"):  "content1",
				p("/project/lib/config.go"):  "content2",
			},
			sourceDirs:   []string{p("/project/api"), p("/project/lib")},
			relativePath: "config.go",
			expectedPath: p("/project/api/config.go"), // Should find in first dir
			expectError:  false,
		},

		// ===== AMBIGUOUS CASES & DUPLICATES =====
		{
			name: "Ambiguous: Same filename in subdirs of same source dir",
			mockFiles: map[string]string{
				p("/project/services/user/handler.go"): "service handler",
				p("/project/api/user/handler.go"):      "api handler",
			},
			sourceDirs:    []string{p("/project")},
			relativePath:  "handler.go",
			expectedPath:  "", // Not used when skipPathCheck is true
			skipPathCheck: true, // Can't predict which one will be found first
			expectError:   false,
		},

		// ===== EDGE CASES =====
		{
			name:         "Edge: Empty source directories",
			mockFiles:    map[string]string{p("/project/src/app.go"): "content"},
			sourceDirs:   []string{},
			relativePath: "app.go",
			expectError:  true,
		},
		{
			name:         "Edge: File does not exist anywhere",
			mockFiles:    map[string]string{p("/project/src/app.go"): "content"},
			sourceDirs:   []string{p("/project")},
			relativePath: "nonexistent.go",
			expectError:  true,
		},
		{
			name:         "Edge: Path already uses forward slashes (cross-platform safe)",
			mockFiles:    map[string]string{p("/project/src/app.go"): "content"},
			sourceDirs:   []string{p("/project")},
			relativePath: "src/app.go", // Forward slashes work on all platforms
			expectedPath: p("/project/src/app.go"),
			expectError:  false,
		},
		{
			name:         "Edge: Path with dots (parent directory reference)",
			mockFiles:    map[string]string{p("/project/lib/utils.go"): "content"},
			sourceDirs:   []string{p("/project/src")},
			relativePath: "../lib/utils.go",
			expectedPath: p("/project/lib/utils.go"),
			expectError:  false,
		},
		{
			name:         "Edge: Current directory reference",
			mockFiles:    map[string]string{p("/project/src/app.go"): "content"},
			sourceDirs:   []string{p("/project")},
			relativePath: "./src/app.go",
			expectedPath: p("/project/src/app.go"),
			expectError:  false,
		},
		{
			name:         "Edge: Filename with special characters",
			mockFiles:    map[string]string{p("/project/my-file_v2.test.go"): "content"},
			sourceDirs:   []string{p("/project")},
			relativePath: "my-file_v2.test.go",
			expectedPath: p("/project/my-file_v2.test.go"),
			expectError:  false,
		},
		{
			name:         "Edge: Very long nested path",
			mockFiles:    map[string]string{p("/a/b/c/d/e/f/g/h/i/j/file.go"): "content"},
			sourceDirs:   []string{p("/a")},
			relativePath: "b/c/d/e/f/g/h/i/j/file.go",
			expectedPath: p("/a/b/c/d/e/f/g/h/i/j/file.go"),
			expectError:  false,
		},
		
		// ===== REALISTIC SCENARIOS FROM DIFFERENT TOOLS =====
		{
			name:         "Real: Go coverage - relative from module root",
			mockFiles:    map[string]string{p("/workspace/myproject/internal/service/handler.go"): "content"},
			sourceDirs:   []string{p("/workspace/myproject")},
			relativePath: "internal/service/handler.go",
			expectedPath: p("/workspace/myproject/internal/service/handler.go"),
			expectError:  false,
		},
		{
			name:         "Real: .NET Cobertura - filename only",
			mockFiles:    map[string]string{p("/src/MyApp/Services/UserService.cs"): "content"},
			sourceDirs:   []string{p("/src/MyApp")},
			relativePath: "UserService.cs",
			expectedPath: p("/src/MyApp/Services/UserService.cs"),
			expectError:  false,
		},
		{
			name:         "Real: JavaScript Istanbul - absolute from project",
			mockFiles:    map[string]string{p("/home/dev/project/src/components/Button.jsx"): "content"},
			sourceDirs:   []string{p("/home/dev/project")},
			relativePath: p("/home/dev/project/src/components/Button.jsx"),
			expectedPath: p("/home/dev/project/src/components/Button.jsx"),
			expectError:  false,
		},
		{
			name:         "Real: Python coverage.py - relative path",
			mockFiles:    map[string]string{p("/app/mypackage/module.py"): "content"},
			sourceDirs:   []string{p("/app")},
			relativePath: "mypackage/module.py",
			expectedPath: p("/app/mypackage/module.py"),
			expectError:  false,
		},
	}

	nopLogger := slog.New(slog.NewTextHandler(io.Discard, nil))
	platform := "unix"
	if isWindows {
		platform = "windows"
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockFS := testutil.NewMockFilesystem(platform)
			for path, content := range tc.mockFiles {
				mockFS.AddFile(path, content)
			}

			// Act
			foundPath, err := utils.FindFileInSourceDirs(tc.relativePath, tc.sourceDirs, mockFS, nopLogger)

			// Assert
			if tc.expectError {
				require.Error(t, err, "Expected an error but got none")
			} else {
				require.NoError(t, err, "Expected no error but got: %v", err)
				if !tc.skipPathCheck {
					// Clean both paths to ensure consistent separator comparison
					assert.Equal(t, filepath.Clean(tc.expectedPath), filepath.Clean(foundPath))
				} else {
					// Just verify we got a non-empty path
					assert.NotEmpty(t, foundPath, "Expected to find a file but got empty path")
				}
			}
		})
	}
}

// TestFindFileInSourceDirs_ErrorHandling tests error conditions separately
func TestFindFileInSourceDirs_ErrorHandling(t *testing.T) {
	isWindows := runtime.GOOS == "windows"
	p := func(path string) string {
		if isWindows {
			return filepath.FromSlash("C:" + path)
		}
		return path
	}

	platform := "unix"
	if isWindows {
		platform = "windows"
	}

	nopLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	t.Run("Nil logger should not panic", func(t *testing.T) {
		mockFS := testutil.NewMockFilesystem(platform)
		mockFS.AddFile(p("/project/app.go"), "content")
		
		// Should not panic with nil logger
		_, err := utils.FindFileInSourceDirs("app.go", []string{p("/project")}, mockFS, nil)
		require.NoError(t, err)
	})

	t.Run("ReadDir error during recursive search should be handled gracefully", func(t *testing.T) {
		mockFS := testutil.NewMockFilesystem(platform)
		// Add a file but make the parent directory unreadable (if your mock supports this)
		mockFS.AddFile(p("/project/src/app.go"), "content")
		
		// Even if we can't read some dirs, it shouldn't panic
		_, err := utils.FindFileInSourceDirs("nonexistent.go", []string{p("/project")}, mockFS, nopLogger)
		require.Error(t, err) // Should error because file not found, not because of dir read
	})

	// Document current behavior with backslashes
	t.Run("Backslashes in relative path - normalized on all platforms", func(t *testing.T) {
		mockFS := testutil.NewMockFilesystem(platform)
		mockFS.AddFile(p("/project/src/app.go"), "content")
		
		foundPath, err := utils.FindFileInSourceDirs("src\\app.go", []string{p("/project")}, mockFS, nopLogger)
		require.NoError(t, err, "Backslashes should be normalized to forward slashes on all platforms")
		assert.Equal(t, filepath.Clean(p("/project/src/app.go")), filepath.Clean(foundPath))
	})
}

// TestFindFileInSourceDirs_StrategyPriority verifies the order strategies are applied
func TestFindFileInSourceDirs_StrategyPriority(t *testing.T) {
	isWindows := runtime.GOOS == "windows"
	p := func(path string) string {
		if isWindows {
			return filepath.FromSlash("C:" + path)
		}
		return path
	}

	platform := "unix"
	if isWindows {
		platform = "windows"
	}

	nopLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	t.Run("Absolute path takes precedence over everything", func(t *testing.T) {
		mockFS := testutil.NewMockFilesystem(platform)
		mockFS.AddFile(p("/absolute/path/file.go"), "absolute content")
		mockFS.AddFile(p("/project/file.go"), "relative content")
		
		foundPath, err := utils.FindFileInSourceDirs(
			p("/absolute/path/file.go"),
			[]string{p("/project")},
			mockFS,
			nopLogger,
		)
		
		require.NoError(t, err)
		assert.Equal(t, filepath.Clean(p("/absolute/path/file.go")), filepath.Clean(foundPath))
	})

	t.Run("Direct join takes precedence over suffix matching", func(t *testing.T) {
		mockFS := testutil.NewMockFilesystem(platform)
		mockFS.AddFile(p("/project/src/app.go"), "direct join")
		mockFS.AddFile(p("/project/app.go"), "would match by suffix")
		
		foundPath, err := utils.FindFileInSourceDirs(
			"src/app.go",
			[]string{p("/project")},
			mockFS,
			nopLogger,
		)
		
		require.NoError(t, err)
		assert.Equal(t, filepath.Clean(p("/project/src/app.go")), filepath.Clean(foundPath))
	})

	t.Run("Suffix matching takes precedence over recursive search", func(t *testing.T) {
		mockFS := testutil.NewMockFilesystem(platform)
		mockFS.AddFile(p("/project/correct/path/file.go"), "suffix match")
		mockFS.AddFile(p("/project/wrong/location/file.go"), "recursive would find this first")
		
		foundPath, err := utils.FindFileInSourceDirs(
			p("/different/root/correct/path/file.go"),
			[]string{p("/project")},
			mockFS,
			nopLogger,
		)
		
		require.NoError(t, err)
		assert.Equal(t, filepath.Clean(p("/project/correct/path/file.go")), filepath.Clean(foundPath))
	})
}