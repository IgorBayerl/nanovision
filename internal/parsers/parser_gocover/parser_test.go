package parser_gocover_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/IgorBayerl/AdlerCov/internal/parsers"
	"github.com/IgorBayerl/AdlerCov/internal/parsers/parser_gocover"
	"github.com/IgorBayerl/AdlerCov/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGoCoverParser_Parse(t *testing.T) {
	const reportFileName = "coverage.out"
	const sourceDir = "/app/src"

	testCases := []struct {
		name          string
		reportContent string
		sourceFiles   map[string]string // Used by the mock filereader to simulate file existence
		sourceDirs    []string
		asserter      func(t *testing.T, result *parsers.ParserResult, err error)
	}{
		{
			name: "Golden Path - Valid report with one file",
			reportContent: `mode: set
github.com/user/project/calculator.go:3.21,4.15 1 1
github.com/user/project/calculator.go:6.24,7.18 1 0
`,
			sourceFiles: map[string]string{
				"/app/src/github.com/user/project/calculator.go": "file content here",
			},
			sourceDirs: []string{sourceDir},
			asserter: func(t *testing.T, result *parsers.ParserResult, err error) {
				const expectedPath = "github.com/user/project/calculator.go"
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Empty(t, result.UnresolvedSourceFiles, "Source file should be resolved")
				assert.Equal(t, "GoCover", result.ParserName)

				require.Len(t, result.FileCoverage, 1, "Should produce coverage for one file")
				fileCov := result.FileCoverage[0]

				assert.Equal(t, expectedPath, fileCov.Path)
				require.Len(t, fileCov.Lines, 4, "Should have metrics for 4 distinct lines")

				// Assert line 3 and 4 were covered
				assert.Equal(t, 1, fileCov.Lines[3].Hits)
				assert.Equal(t, 1, fileCov.Lines[4].Hits)

				// Assert line 6 and 7 were not covered
				assert.Equal(t, 0, fileCov.Lines[6].Hits)
				assert.Equal(t, 0, fileCov.Lines[7].Hits)
			},
		},
		{
			name: "Source File Not Found - Should report as unresolved",
			reportContent: `mode: set
github.com/user/project/calculator.go:3.21,4.15 1 1
`,
			sourceFiles: map[string]string{}, // Mock filesystem is empty
			sourceDirs:  []string{sourceDir},
			asserter: func(t *testing.T, result *parsers.ParserResult, err error) {
				const unresolvedPath = "github.com/user/project/calculator.go"
				require.NoError(t, err)
				require.NotNil(t, result)

				// Still produces FileCoverage, but marks the file as unresolved
				require.Len(t, result.FileCoverage, 1)
				assert.Equal(t, unresolvedPath, result.FileCoverage[0].Path)

				require.Len(t, result.UnresolvedSourceFiles, 1)
				assert.Equal(t, unresolvedPath, result.UnresolvedSourceFiles[0])
			},
		},
		{
			name: "Report with multiple files",
			reportContent: `mode: set
project/file1.go:1.1,1.10 1 5
project/file2.go:2.1,2.12 1 0
`,
			sourceFiles: map[string]string{
				"/app/src/project/file1.go": "content",
				"/app/src/project/file2.go": "content",
			},
			sourceDirs: []string{sourceDir},
			asserter: func(t *testing.T, result *parsers.ParserResult, err error) {
				require.NoError(t, err)
				require.Len(t, result.FileCoverage, 2)
				assert.Empty(t, result.UnresolvedSourceFiles)

				// Use a map for easier lookup instead of relying on slice order
				coverageByPath := make(map[string]parsers.FileCoverage)
				for _, fc := range result.FileCoverage {
					coverageByPath[fc.Path] = fc
				}

				require.Contains(t, coverageByPath, "project/file1.go")
				assert.Equal(t, 5, coverageByPath["project/file1.go"].Lines[1].Hits)

				require.Contains(t, coverageByPath, "project/file2.go")
				assert.Equal(t, 0, coverageByPath["project/file2.go"].Lines[2].Hits)
			},
		},
		{
			name:          "Report is logically empty (only mode line)",
			reportContent: `mode: set`,
			sourceFiles:   map[string]string{},
			sourceDirs:    []string{sourceDir},
			asserter: func(t *testing.T, result *parsers.ParserResult, err error) {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Empty(t, result.FileCoverage, "Should produce no file coverage for an empty report")
				assert.Empty(t, result.UnresolvedSourceFiles)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			tmpDir := t.TempDir()
			reportPath := filepath.Join(tmpDir, reportFileName)
			err := os.WriteFile(reportPath, []byte(tc.reportContent), 0644)
			require.NoError(t, err)

			// The mock filesystem is now used to check for file existence
			mockFS := testutil.NewMockFilesystem("unix")
			for path, content := range tc.sourceFiles {
				mockFS.AddFile(path, content)
			}

			mockConfig := testutil.NewTestConfig(tc.sourceDirs)
			parser := parser_gocover.NewGoCoverParser(mockFS)

			// Act
			result, err := parser.Parse(reportPath, mockConfig)

			// Assert
			tc.asserter(t, result, err)
		})
	}
}
