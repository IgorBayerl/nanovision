package parser_lcov_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/IgorBayerl/nanovision/internal/parsers"
	"github.com/IgorBayerl/nanovision/internal/parsers/parser_lcov"
	"github.com/IgorBayerl/nanovision/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLcovParser_Parse(t *testing.T) {
	const reportFileName = "lcov.info"
	const sourceDir = "/app/src"

	testCases := []struct {
		name          string
		reportContent string
		sourceFiles   map[string]string // Used by the mock filereader
		sourceDirs    []string
		asserter      func(t *testing.T, result *parsers.ParserResult, err error)
	}{
		{
			name: "Golden Path - Valid report with line and branch coverage",
			reportContent: `TN:
SF:/app/src/utils/math.js
FN:1,add
FNDA:5,add
FNF:1
FNH:1
DA:1,5
DA:2,5
DA:3,0
BRDA:2,0,0,5
BRDA:2,0,1,-
LF:3
LH:2
end_of_record
`,
			sourceFiles: map[string]string{
				"/app/src/utils/math.js": "function add(a,b){...}",
			},
			sourceDirs: []string{sourceDir},
			asserter: func(t *testing.T, result *parsers.ParserResult, err error) {
				const expectedPath = "/app/src/utils/math.js"
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, "LCOV", result.ParserName)
				assert.Empty(t, result.UnresolvedSourceFiles)

				require.Len(t, result.FileCoverage, 1)
				fileCov := result.FileCoverage[0]

				assert.Equal(t, expectedPath, fileCov.Path)
				require.Len(t, fileCov.Lines, 3)

				// Check Lines
				assert.Equal(t, 5, fileCov.Lines[1].Hits)
				assert.Equal(t, 5, fileCov.Lines[2].Hits)
				assert.Equal(t, 0, fileCov.Lines[3].Hits)

				// Check Branches (Line 2 has 2 branches: one hit(5), one miss(-))
				assert.Equal(t, 2, fileCov.Lines[2].TotalBranches)
				assert.Equal(t, 1, fileCov.Lines[2].CoveredBranches)
			},
		},
		{
			name: "Source File Not Found",
			reportContent: `SF:src/missing.ts
DA:1,1
end_of_record`,
			sourceFiles: map[string]string{}, // Empty FS
			sourceDirs:  []string{sourceDir},
			asserter: func(t *testing.T, result *parsers.ParserResult, err error) {
				require.NoError(t, err)
				require.Len(t, result.FileCoverage, 1)
				require.Len(t, result.UnresolvedSourceFiles, 1)
				assert.Equal(t, "src/missing.ts", result.UnresolvedSourceFiles[0])
			},
		},
		{
			name: "Multiple Records",
			reportContent: `SF:/app/src/a.js
DA:10,1
end_of_record
SF:/app/src/b.js
DA:20,2
end_of_record`,
			sourceFiles: map[string]string{
				"/app/src/a.js": "content",
				"/app/src/b.js": "content",
			},
			sourceDirs: []string{sourceDir},
			asserter: func(t *testing.T, result *parsers.ParserResult, err error) {
				require.NoError(t, err)
				require.Len(t, result.FileCoverage, 2)

				// Order isn't guaranteed by map iteration usually, but slice append is ordered
				assert.Equal(t, "/app/src/a.js", result.FileCoverage[0].Path)
				assert.Equal(t, 1, result.FileCoverage[0].Lines[10].Hits)

				assert.Equal(t, "/app/src/b.js", result.FileCoverage[1].Path)
				assert.Equal(t, 2, result.FileCoverage[1].Lines[20].Hits)
			},
		},
		{
			name: "Malformed - Missing end_of_record (should verify resilience)",
			reportContent: `SF:/app/src/a.js
DA:1,1`,
			sourceFiles: map[string]string{
				"/app/src/a.js": "content",
			},
			sourceDirs: []string{sourceDir},
			asserter: func(t *testing.T, result *parsers.ParserResult, err error) {
				require.NoError(t, err)
				require.Len(t, result.FileCoverage, 1)
				assert.Equal(t, "/app/src/a.js", result.FileCoverage[0].Path)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			reportPath := filepath.Join(tmpDir, reportFileName)
			err := os.WriteFile(reportPath, []byte(tc.reportContent), 0644)
			require.NoError(t, err)

			mockFS := testutil.NewMockFilesystem("unix")
			for path, content := range tc.sourceFiles {
				mockFS.AddFile(path, content)
			}

			mockConfig := testutil.NewTestConfig(tc.sourceDirs)
			parser := parser_lcov.NewLcovParser(mockFS)

			result, err := parser.Parse(reportPath, mockConfig)
			tc.asserter(t, result, err)
		})
	}
}
