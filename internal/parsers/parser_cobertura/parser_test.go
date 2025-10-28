package parser_cobertura_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/IgorBayerl/nanovision/internal/parsers"
	"github.com/IgorBayerl/nanovision/internal/parsers/parser_cobertura"
	"github.com/IgorBayerl/nanovision/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCoberturaParser_Parse(t *testing.T) {
	const reportFileName = "cobertura.xml"
	const sourceDir = "/app/src"
	const sourceFilePath = "MyProject/Calculator.cs" // This is the path inside the report
	const resolvedSourcePath = "/app/src/MyProject/Calculator.cs"

	testCases := []struct {
		name          string
		reportContent string
		sourceFiles   map[string]string // map[path]content, for mock filesystem
		sourceDirs    []string
		asserter      func(t *testing.T, result *parsers.ParserResult, err error)
	}{
		{
			name: "Golden Path - Valid report with branch coverage",
			reportContent: `<?xml version="1.0" encoding="utf-8"?>
<coverage lines-covered="6" lines-valid="8" branches-covered="1" branches-valid="2" timestamp="1672531200">
  <packages>
    <package name="MyProject.Core">
      <classes>
        <class name="MyProject.Core.Calculator" filename="MyProject/Calculator.cs">
          <lines>
            <line number="5" hits="1" branch="false" />
            <line number="6" hits="1" branch="false" />
            <line number="7" hits="1" branch="false" />
            <line number="10" hits="2" branch="true" condition-coverage="50% (1/2)" />
            <line number="11" hits="1" branch="false" />
            <line number="12" hits="0" branch="false" />
            <line number="13" hits="0" branch="false" />
            <line number="15" hits="1" branch="false" />
          </lines>
        </class>
      </classes>
    </package>
  </packages>
</coverage>`,
			sourceFiles: map[string]string{
				resolvedSourcePath: `// Dummy content`,
			},
			sourceDirs: []string{sourceDir},
			asserter: func(t *testing.T, result *parsers.ParserResult, err error) {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Empty(t, result.UnresolvedSourceFiles)
				assert.Equal(t, "Cobertura", result.ParserName)
				require.NotNil(t, result.Timestamp)

				require.Len(t, result.FileCoverage, 1, "Should produce coverage for one file")
				fileCov := result.FileCoverage[0]

				assert.Equal(t, sourceFilePath, fileCov.Path)
				require.Len(t, fileCov.Lines, 8, "Should have metrics for 8 distinct lines")

				// Assert line hits
				assert.Equal(t, 1, fileCov.Lines[5].Hits)
				assert.Equal(t, 1, fileCov.Lines[6].Hits)
				assert.Equal(t, 1, fileCov.Lines[7].Hits)
				assert.Equal(t, 2, fileCov.Lines[10].Hits)
				assert.Equal(t, 1, fileCov.Lines[11].Hits)
				assert.Equal(t, 0, fileCov.Lines[12].Hits)
				assert.Equal(t, 0, fileCov.Lines[13].Hits)
				assert.Equal(t, 1, fileCov.Lines[15].Hits)

				// Assert branch data on the specific line
				assert.Equal(t, 1, fileCov.Lines[10].CoveredBranches)
				assert.Equal(t, 2, fileCov.Lines[10].TotalBranches)

				// Assert no branch data on non-branch line
				assert.Zero(t, fileCov.Lines[5].TotalBranches)
			},
		},
		{
			name: "Source File Not Found - Should report as unresolved",
			reportContent: `<?xml version="1.0"?>
<coverage>
  <packages>
    <package name="MyProject.Core">
      <classes>
        <class name="MyProject.Core.Calculator" filename="MyProject/DoesNotExist.cs">
          <lines><line number="5" hits="1" /></lines>
        </class>
      </classes>
    </package>
  </packages>
</coverage>`,
			sourceFiles: map[string]string{}, // Mock filesystem is empty
			sourceDirs:  []string{sourceDir},
			asserter: func(t *testing.T, result *parsers.ParserResult, err error) {
				const unresolvedPath = "MyProject/DoesNotExist.cs"
				require.NoError(t, err)
				require.NotNil(t, result)

				require.Len(t, result.FileCoverage, 1)
				assert.Equal(t, unresolvedPath, result.FileCoverage[0].Path)

				require.Len(t, result.UnresolvedSourceFiles, 1)
				assert.Equal(t, unresolvedPath, result.UnresolvedSourceFiles[0])
			},
		},
		{
			name: "Report is logically empty (no packages)",
			reportContent: `<?xml version="1.0"?>
<coverage lines-covered="0" lines-valid="0" branches-covered="0" branches-valid="0">
  <packages />
</coverage>`,
			sourceFiles: map[string]string{},
			sourceDirs:  []string{sourceDir},
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

			mockFS := testutil.NewMockFilesystem("unix")
			for path, content := range tc.sourceFiles {
				mockFS.AddFile(path, content)
			}

			mockConfig := testutil.NewTestConfig(tc.sourceDirs)
			parser := parser_cobertura.NewCoberturaParser(mockFS)

			// Act
			result, err := parser.Parse(reportPath, mockConfig)

			// Assert
			tc.asserter(t, result, err)
		})
	}
}
