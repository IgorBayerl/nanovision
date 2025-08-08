package utils_test

import (
	"io"
	"log/slog"
	"path/filepath"
	"testing"

	"github.com/IgorBayerl/AdlerCov/internal/testutil"
	"github.com/IgorBayerl/AdlerCov/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindFileInSourceDirs(t *testing.T) {
	testCases := []struct {
		name         string
		platform     string // "unix" or "windows"
		mockFiles    map[string]string
		sourceDirs   []string
		relativePath string
		expectedPath string
		expectError  bool
	}{
		{
			name:         "Success: Simple relative path found in first source dir",
			platform:     "unix",
			mockFiles:    map[string]string{"/project/src/app.go": "content"},
			sourceDirs:   []string{"/project"},
			relativePath: "src/app.go",
			expectedPath: "/project/src/app.go",
			expectError:  false,
		},
		{
			name:         "Success: The C# Cobertura scenario (filename only) found via recursive search",
			platform:     "windows",
			mockFiles:    map[string]string{"C:/dev/project/services/MyService.cs": "content"},
			sourceDirs:   []string{"C:\\dev\\project"},
			relativePath: "MyService.cs",
			expectedPath: "C:\\dev\\project\\services\\MyService.cs",
			expectError:  false,
		},
		{
			name:         "Success: The gcov scenario (absolute path in report)",
			platform:     "windows",
			mockFiles:    map[string]string{"C:/project/src/calculator.cpp": "content"},
			sourceDirs:   []string{"C:\\another\\dir"}, // Should be ignored
			relativePath: "C:/project/src/calculator.cpp",
			expectedPath: "C:\\project\\src\\calculator.cpp",
			expectError:  false,
		},
		{
			name:         "Success: The CI/CD scenario (suffix matching)",
			platform:     "unix",
			mockFiles:    map[string]string{"/local/dev/src/main.js": "content"},
			sourceDirs:   []string{"/local/dev"},
			relativePath: "/home/runner/work/project/src/main.js",
			expectedPath: "/local/dev/src/main.js",
			expectError:  false,
		},
		{
			name:         "Success: Found in second source directory",
			platform:     "unix",
			mockFiles:    map[string]string{"/project/lib/utils.go": "content"},
			sourceDirs:   []string{"/project/api", "/project/lib"},
			relativePath: "utils.go",
			expectedPath: "/project/lib/utils.go",
			expectError:  false,
		},
		{
			name:         "Success: Windows path with mixed separators",
			platform:     "windows",
			mockFiles:    map[string]string{"C:/dev/project/api/controllers/HealthController.cs": "content"},
			sourceDirs:   []string{"C:\\dev\\project"},
			relativePath: "api/controllers/HealthController.cs", // Report uses '/'
			expectedPath: "C:\\dev\\project\\api\\controllers\\HealthController.cs",
			expectError:  false,
		},
		{
			name:         "Failure: File does not exist in any directory",
			platform:     "unix",
			mockFiles:    map[string]string{"/project/src/app.go": "content"},
			sourceDirs:   []string{"/project"},
			relativePath: "nonexistent.go",
			expectError:  true,
		},
		{
			name:         "Failure: Empty source directories",
			platform:     "unix",
			mockFiles:    map[string]string{"/project/src/app.go": "content"},
			sourceDirs:   []string{},
			relativePath: "app.go",
			expectError:  true,
		},
		{
			name:         "Success: Path contains spaces",
			platform:     "windows",
			mockFiles:    map[string]string{"C:/My Project/My Source/app.cs": "content"},
			sourceDirs:   []string{"C:\\My Project"},
			relativePath: "My Source/app.cs",
			expectedPath: "C:\\My Project\\My Source\\app.cs",
			expectError:  false,
		},
	}

	nopLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockFS := testutil.NewMockFilesystem(tc.platform)
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
				// Clean both paths to ensure consistent separator comparison
				assert.Equal(t, filepath.Clean(tc.expectedPath), filepath.Clean(foundPath))
			}
		})
	}
}
