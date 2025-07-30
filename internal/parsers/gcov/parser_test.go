package gcov_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/IgorBayerl/AdlerCov/internal/parsers"
	"github.com/IgorBayerl/AdlerCov/internal/parsers/gcov"
	"github.com/IgorBayerl/AdlerCov/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGCovParser_Parse(t *testing.T) {
	const reportFileName = "calculator.cpp.gcov"
	const sourceDir = "C:/project/src"
	const normalizedSourcePath = "C:/project/src/calculator.cpp" // Path is absolute in this gcov report

	const gcovReportContent = `        -:    0:Source:C:/project/src/calculator.cpp
        -:    0:Runs:1
function _ZN10Calculator6divideEdd called 3 returned 100% blocks executed 88%
        3:   16:double Calculator::divide(double a, double b) {
        3:   17:    if (b == 0.0) {
branch  0 taken 33%
branch  1 taken 67%
        1:   18:        throw std::invalid_argument("Division by zero");
        -:   19:    }
        2:   20:    return a / b;
        -:   21:}
function _ZN10Calculator4signEi called 2 returned 100% blocks executed 83%
        2:   23:int Calculator::sign(int x) {
        2:   24:    if (x > 0) {
branch  0 taken 50%
branch  1 taken 50%
        1:   25:        return 1;
        1:   26:    } else if (x < 0) {
branch  0 taken 100%
branch  1 never executed
        1:   27:        return -1;
        -:   28:    } else {
    #####:   30:        return 0;
        -:   31:    }
        -:   32:}`

	testCases := []struct {
		name          string
		reportContent string
		sourceFiles   map[string]string // For mock filesystem to check existence
		sourceDirs    []string
		asserter      func(t *testing.T, result *parsers.ParserResult, err error)
	}{
		{
			name:          "Golden Path - Valid report with branch coverage",
			reportContent: gcovReportContent,
			sourceFiles: map[string]string{
				normalizedSourcePath: "// C++ source content",
			},
			sourceDirs: []string{sourceDir},
			asserter: func(t *testing.T, result *parsers.ParserResult, err error) {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Empty(t, result.UnresolvedSourceFiles)
				assert.Equal(t, "GCov", result.ParserName)

				require.Len(t, result.FileCoverage, 1)
				fileCov := result.FileCoverage[0]

				assert.Equal(t, normalizedSourcePath, fileCov.Path)
				require.NotNil(t, fileCov.Lines)

				// Assert line hits
				assert.Equal(t, 3, fileCov.Lines[16].Hits)
				assert.Equal(t, 3, fileCov.Lines[17].Hits)
				assert.Equal(t, 1, fileCov.Lines[18].Hits)
				assert.Equal(t, 2, fileCov.Lines[20].Hits)
				assert.Equal(t, 0, fileCov.Lines[30].Hits) // '#####' is 0 hits

				// Assert branch metrics
				// Line 17: two branches, both "taken"
				assert.Equal(t, 2, fileCov.Lines[17].TotalBranches)
				assert.Equal(t, 2, fileCov.Lines[17].CoveredBranches)

				// Line 24: two branches, both "taken"
				assert.Equal(t, 2, fileCov.Lines[24].TotalBranches)
				assert.Equal(t, 2, fileCov.Lines[24].CoveredBranches)

				// Line 26: two branches, one "taken", one "never executed"
				assert.Equal(t, 2, fileCov.Lines[26].TotalBranches)
				assert.Equal(t, 1, fileCov.Lines[26].CoveredBranches)
			},
		},
		{
			name:          "Source File Not Found",
			reportContent: gcovReportContent,
			sourceFiles:   map[string]string{}, // Mock filesystem is empty
			sourceDirs:    []string{"/another/dir"},
			asserter: func(t *testing.T, result *parsers.ParserResult, err error) {
				require.NoError(t, err)
				require.NotNil(t, result)

				// The parser should still produce coverage data
				require.Len(t, result.FileCoverage, 1)
				assert.Equal(t, normalizedSourcePath, result.FileCoverage[0].Path)

				// But it should also report the file as unresolved
				require.Len(t, result.UnresolvedSourceFiles, 1)
				assert.Equal(t, normalizedSourcePath, result.UnresolvedSourceFiles[0])
			},
		},
		{
			name:          "Invalid gcov file (missing source line)",
			reportContent: "function main called 1",
			sourceFiles:   map[string]string{},
			sourceDirs:    []string{},
			asserter: func(t *testing.T, result *parsers.ParserResult, err error) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "invalid gcov format")
				assert.Nil(t, result)
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

			mockFS := testutil.NewMockFilesystem("windows") // Test with windows-like paths
			for path, content := range tc.sourceFiles {
				mockFS.AddFile(path, content)
			}

			mockConfig := testutil.NewTestConfig(tc.sourceDirs)
			parser := gcov.NewGCovParser(mockFS)

			// Act
			result, err := parser.Parse(reportPath, mockConfig)

			// Assert
			tc.asserter(t, result, err)
		})
	}
}
