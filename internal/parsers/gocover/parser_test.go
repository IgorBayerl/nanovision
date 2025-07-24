package gocover_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/IgorBayerl/AdlerCov/internal/model"
	"github.com/IgorBayerl/AdlerCov/internal/parsers"
	"github.com/IgorBayerl/AdlerCov/internal/parsers/gocover"
	"github.com/IgorBayerl/AdlerCov/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGoCoverParser_Parse(t *testing.T) {
	const reportFileName = "coverage.out"
	const sourceDir = "/app/src"
	const goModPath = "/app/src/go.mod"
	const goModContent = "module github.com/user/project"

	testCases := []struct {
		name          string
		reportContent string
		sourceFiles   map[string]string
		sourceDirs    []string
		asserter      func(t *testing.T, result *parsers.ParserResult, err error)
	}{
		// The "Golden Path" - a valid report with all sources found.
		{
			name: "Golden Path - Valid report with found source file",
			reportContent: `mode: set
github.com/user/project/calculator.go:3.34,5.2 1 1
github.com/user/project/calculator.go:7.38,9.2 1 0
`,
			sourceFiles: map[string]string{
				"/app/src/github.com/user/project/calculator.go": `package calculator

func Add(a, b int) int {
	return a + b
}

func Subtract(a, b int) int {
	return a - b
}
`,
				goModPath: goModContent,
			},
			sourceDirs: []string{sourceDir},
			asserter: func(t *testing.T, result *parsers.ParserResult, err error) {
				// Arrange - Expected values
				const resolvedSourcePath = "/app/src/github.com/user/project/calculator.go"

				// Assert - No errors, valid result
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Empty(t, result.UnresolvedSourceFiles)
				assert.Equal(t, "GoCover", result.ParserName)

				// Assert - Assembly level
				require.Len(t, result.Assemblies, 1)
				assembly := result.Assemblies[0]
				assert.Equal(t, "github.com/user/project", assembly.Name)
				assert.Equal(t, 1, assembly.LinesCovered)
				assert.Equal(t, 2, assembly.LinesValid)

				// Assert - Class (Go Package) level
				require.Len(t, assembly.Classes, 1)
				class := assembly.Classes[0]
				assert.Equal(t, "(root)", class.DisplayName)
				assert.Equal(t, 2, class.TotalMethods)
				assert.Equal(t, 1, class.CoveredMethods)
				assert.Equal(t, 1, class.FullyCoveredMethods)

				// Assert - File level
				require.Len(t, class.Files, 1)
				file := class.Files[0]
				assert.Equal(t, resolvedSourcePath, filepath.ToSlash(file.Path))
				assert.Equal(t, 1, file.CoveredLines)
				assert.Equal(t, 2, file.CoverableLines)
				assert.Equal(t, 9, file.TotalLines)

				// Assert - Method level
				require.Len(t, class.Methods, 2)
				addMethod := testutil.FindMethod(t, class.Methods, "Add")
				assert.Equal(t, 3, addMethod.FirstLine)
				assert.Equal(t, 5, addMethod.LastLine)
				assert.InDelta(t, 1.0, addMethod.LineRate, 0.001)

				subtractMethod := testutil.FindMethod(t, class.Methods, "Subtract")
				assert.Equal(t, 7, subtractMethod.FirstLine)
				assert.Equal(t, 9, subtractMethod.LastLine)
				assert.InDelta(t, 0.0, subtractMethod.LineRate, 0.001)

				// Assert - Line level
				line4 := testutil.FindLine(t, file.Lines, 4)
				assert.Equal(t, 1, line4.Hits)
				assert.Equal(t, model.Covered, line4.LineVisitStatus)

				line8 := testutil.FindLine(t, file.Lines, 8)
				assert.Equal(t, 0, line8.Hits)
				assert.Equal(t, model.NotCovered, line8.LineVisitStatus)
			},
		},
		// Source File Edge Case - file cannot be found.
		{
			name: "Source File Edge Case - Unresolved source file",
			reportContent: `mode: set
github.com/user/project/calculator.go:3.34,5.2 1 1
`,
			sourceFiles: map[string]string{
				// Source file is missing from our mock filesystem
			},
			sourceDirs: []string{sourceDir},
			asserter: func(t *testing.T, result *parsers.ParserResult, err error) {
				const sourceFilePath = "github.com/user/project/calculator.go"
				require.NoError(t, err)
				require.NotNil(t, result)
				require.Len(t, result.UnresolvedSourceFiles, 1)
				assert.Equal(t, sourceFilePath, result.UnresolvedSourceFiles[0])
				// The parser should still produce a result, but it will be empty.
				assert.Empty(t, result.Assemblies)
			},
		},
		// Report File Edge Case - valid report but no coverage data.
		{
			name:          "Report File Edge Case - Report is logically empty",
			reportContent: `mode: set`,
			sourceFiles:   map[string]string{},
			sourceDirs:    []string{sourceDir},
			asserter: func(t *testing.T, result *parsers.ParserResult, err error) {
				require.NoError(t, err)
				require.NotNil(t, result)
				// The parser should succeed but produce no assemblies.
				assert.Empty(t, result.Assemblies)
				assert.Empty(t, result.UnresolvedSourceFiles)
			},
		},
		// Multiple files in the same package (class).
		{
			name: "Golden Path - Multiple files in one package",
			reportContent: `mode: set
github.com/user/project/calculator.go:3.34,5.2 1 1
github.com/user/project/greeter.go:3.19,5.2 1 1
`,
			sourceFiles: map[string]string{
				"/app/src/github.com/user/project/calculator.go": `package calculator
func Add(a, b int) int {
	return a + b
}`,
				"/app/src/github.com/user/project/greeter.go": `package calculator
func Greet(name string) string {
	return "Hello, " + name
}`,
				goModPath: goModContent,
			},
			sourceDirs: []string{sourceDir},
			asserter: func(t *testing.T, result *parsers.ParserResult, err error) {
				require.NoError(t, err)
				require.NotNil(t, result)
				require.Len(t, result.Assemblies, 1)
				assembly := result.Assemblies[0]
				require.Len(t, assembly.Classes, 1)
				class := assembly.Classes[0]

				assert.Equal(t, 2, class.TotalMethods)
				assert.Equal(t, 2, class.CoveredMethods)
				assert.Len(t, class.Files, 2, "Should have processed two separate files")
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
			for _, dir := range tc.sourceDirs {
				mockFS.AddDir(dir)
			}
			for path, content := range tc.sourceFiles {
				mockFS.AddFile(path, content)
			}

			mockConfig := testutil.NewTestConfig(tc.sourceDirs)
			parser := gocover.NewGoCoverParser(mockFS)

			// Act
			result, err := parser.Parse(reportPath, mockConfig)

			// Assert
			tc.asserter(t, result, err)
		})
	}
}

// Add this new function to internal/parsers/gocover/parser_test.go

func TestGoCoverParser_SupportsFile(t *testing.T) {
	testCases := []struct {
		name        string
		fileName    string
		fileContent string
		shouldMatch bool
	}{
		{
			name:        "Valid Go cover file with 'mode:' prefix",
			fileName:    "coverage.out",
			fileContent: "mode: set\ngithub.com/user/project/file.go:1.1,2.2 1 1",
			shouldMatch: true,
		},
		{
			name:        "Valid Go cover file with different name",
			fileName:    "gocover.txt",
			fileContent: "mode: atomic\ngithub.com/user/project/file.go:1.1,2.2 1 1",
			shouldMatch: true,
		},
		{
			name:        "File with 'mode:' prefix but with leading whitespace",
			fileName:    "coverage.out",
			fileContent: "  mode: set\n...",
			shouldMatch: false,
		},
		{
			name:        "File without 'mode:' prefix",
			fileName:    "coverage.out",
			fileContent: "github.com/user/project/file.go:1.1,2.2 1 1",
			shouldMatch: false,
		},
		{
			name:        "Empty file",
			fileName:    "empty.out",
			fileContent: "",
			shouldMatch: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			tmpDir := t.TempDir()
			filePath := filepath.Join(tmpDir, tc.fileName)
			err := os.WriteFile(filePath, []byte(tc.fileContent), 0644)
			require.NoError(t, err, "Failed to set up test file")

			parser := gocover.NewGoCoverParser(nil)

			// Act
			actual := parser.SupportsFile(filePath)

			// Assert
			assert.Equal(t, tc.shouldMatch, actual)
		})
	}
}
