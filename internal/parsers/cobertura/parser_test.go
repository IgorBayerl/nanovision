package cobertura_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/IgorBayerl/AdlerCov/internal/model"
	"github.com/IgorBayerl/AdlerCov/internal/parsers"
	"github.com/IgorBayerl/AdlerCov/internal/parsers/cobertura"
	"github.com/IgorBayerl/AdlerCov/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCoberturaParser_Parse(t *testing.T) {
	const reportFileName = "cobertura.xml"
	const sourceDir = "/app/src"
	const sourceFilePath = "MyProject/Calculator.cs"
	const resolvedSourcePath = "/app/src/MyProject/Calculator.cs"

	testCases := []struct {
		name          string
		reportContent string
		sourceFiles   map[string]string // map[path]content
		sourceDirs    []string
		asserter      func(t *testing.T, result *parsers.ParserResult, err error)
	}{
		// The "Golden Path" - valid report with branches.
		{
			name: "Golden Path - Valid report with branch coverage",
			reportContent: `<?xml version="1.0" encoding="utf-8"?>
<coverage lines-covered="5" lines-valid="8" branches-covered="1" branches-valid="2" complexity="3" version="1.0" timestamp="1672531200">
  <sources>
    <source>/app/src</source>
  </sources>
  <packages>
    <package name="MyProject.Core" line-rate="0.625" branch-rate="0.5">
      <classes>
        <class name="MyProject.Core.Calculator" filename="MyProject/Calculator.cs" line-rate="0.625" branch-rate="0.5" complexity="3">
          <methods>
            <method name="Add" signature="(System.Int32,System.Int32)" line-rate="1.0" branch-rate="1.0" complexity="1">
              <lines>
                <line number="5" hits="1" branch="false" />
                <line number="6" hits="1" branch="false" />
                <line number="7" hits="1" branch="false" />
              </lines>
            </method>
            <method name="Divide" signature="(System.Int32,System.Int32)" line-rate="0.4" branch-rate="0.5" complexity="2">
              <lines>
                <line number="10" hits="2" branch="true" condition-coverage="50% (1/2)" />
                <line number="11" hits="1" branch="false" />
                <line number="12" hits="0" branch="false" />
                <line number="13" hits="0" branch="false" />
                <line number="15" hits="1" branch="false" />
              </lines>
            </method>
          </methods>
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
				resolvedSourcePath: `namespace MyProject.Core;

public class Calculator
{
    public int Add(int a, int b) // Line 5
    {
        return a + b; // Line 7
    }

    public int Divide(int a, int b) // Line 10
    {
        if (b == 0) // Line 11
            return 0; // Line 12

        return a / b; // Line 15
    }
}`,
			},
			sourceDirs: []string{sourceDir},
			asserter: func(t *testing.T, result *parsers.ParserResult, err error) {
				// Assert - No errors, valid result
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Empty(t, result.UnresolvedSourceFiles)
				assert.Equal(t, "Cobertura", result.ParserName)
				assert.True(t, result.SupportsBranchCoverage)

				// Assert - Assembly level
				require.Len(t, result.Assemblies, 1)
				assembly := result.Assemblies[0]
				assert.Equal(t, "MyProject.Core", assembly.Name)
				// Lines with hits > 0 are: 5,6,7,10,11,15 -> 6 lines
				assert.Equal(t, 6, assembly.LinesCovered)
				assert.Equal(t, 8, assembly.LinesValid)
				require.NotNil(t, assembly.BranchesCovered)
				assert.Equal(t, 1, *assembly.BranchesCovered)
				require.NotNil(t, assembly.BranchesValid)
				assert.Equal(t, 2, *assembly.BranchesValid)

				// Assert - Class level
				require.Len(t, assembly.Classes, 1)
				class := assembly.Classes[0]
				assert.Equal(t, "MyProject.Core.Calculator", class.DisplayName)
				assert.Equal(t, 2, class.TotalMethods)
				assert.Equal(t, 2, class.CoveredMethods)
				assert.Equal(t, 1, class.FullyCoveredMethods)

				// Assert - File level
				require.Len(t, class.Files, 1)
				file := class.Files[0]
				assert.Equal(t, resolvedSourcePath, filepath.ToSlash(file.Path))
				assert.Equal(t, 6, file.CoveredLines)
				assert.Equal(t, 8, file.CoverableLines)
				assert.Equal(t, 17, file.TotalLines)

				// Assert - Method level
				addMethod := testutil.FindMethod(t, class.Methods, "Add")
				assert.InDelta(t, 1.0, addMethod.LineRate, 0.001)

				divideMethod := testutil.FindMethod(t, class.Methods, "Divide")
				assert.InDelta(t, 0.6, divideMethod.LineRate, 0.001)
				require.NotNil(t, divideMethod.BranchRate)
				assert.InDelta(t, 0.5, *divideMethod.BranchRate, 0.001)

				// Assert - Line level (with branches)
				line10 := testutil.FindLine(t, file.Lines, 10)
				assert.Equal(t, 2, line10.Hits)
				assert.True(t, line10.IsBranchPoint)
				assert.Equal(t, 1, line10.CoveredBranches)
				assert.Equal(t, 2, line10.TotalBranches)
				assert.Equal(t, model.PartiallyCovered, line10.LineVisitStatus)

				line11 := testutil.FindLine(t, file.Lines, 11)
				assert.Equal(t, 1, line11.Hits) // XML says hits="1"
				assert.Equal(t, model.Covered, line11.LineVisitStatus)

				line12 := testutil.FindLine(t, file.Lines, 12)
				assert.Equal(t, 0, line12.Hits) // XML says hits="0"
				assert.Equal(t, model.NotCovered, line12.LineVisitStatus)

				line15 := testutil.FindLine(t, file.Lines, 15)
				assert.Equal(t, 1, line15.Hits) // XML says hits="1"
				assert.Equal(t, model.Covered, line15.LineVisitStatus)
			},
		},
		// Source File Edge Case - file cannot be found.
		{
			name: "Source File Edge Case - Unresolved source file",
			reportContent: `<?xml version="1.0"?>
<coverage lines-covered="1" lines-valid="1" branches-covered="0" branches-valid="0">
  <packages>
    <package name="MyProject.Core">
      <classes>
        <class name="MyProject.Core.Calculator" filename="MyProject/Calculator.cs" line-rate="1.0" branch-rate="1.0">
          <methods/>
          <lines><line number="5" hits="1" branch="false" /></lines>
        </class>
      </classes>
    </package>
  </packages>
</coverage>`,
			sourceFiles: map[string]string{
				// Source file is missing
			},
			sourceDirs: []string{sourceDir},
			asserter: func(t *testing.T, result *parsers.ParserResult, err error) {
				require.NoError(t, err)
				require.NotNil(t, result)
				require.Len(t, result.UnresolvedSourceFiles, 1)
				assert.Equal(t, sourceFilePath, result.UnresolvedSourceFiles[0])

				// gets filtered out when no source files can be resolved
				require.Len(t, result.Assemblies, 1)
				assembly := result.Assemblies[0]
				assert.Equal(t, "MyProject.Core", assembly.Name)
				assert.Empty(t, assembly.Classes, "Assembly should have no classes when source files are unresolved")

				// Verify metrics are zeroed out
				assert.Equal(t, 0, assembly.LinesCovered)
				assert.Equal(t, 0, assembly.LinesValid)
				assert.Equal(t, 0, assembly.TotalLines)
				assert.Nil(t, assembly.BranchesCovered)
				assert.Nil(t, assembly.BranchesValid)
			},
		},
		// Report File Edge Case - valid report but no coverage data.
		{
			name: "Report File Edge Case - Report is logically empty",
			reportContent: `<?xml version="1.0"?>
<coverage lines-covered="0" lines-valid="0" branches-covered="0" branches-valid="0">
  <packages />
</coverage>`,
			sourceFiles: map[string]string{},
			sourceDirs:  []string{sourceDir},
			asserter: func(t *testing.T, result *parsers.ParserResult, err error) {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Empty(t, result.Assemblies)
			},
		},
		// Filtering out a compiler-generated class.
		{
			name: "Language Processor - Filters compiler-generated class",
			// FIX: Properly escape the < and > characters as &lt; and &gt;
			reportContent: `<?xml version="1.0"?>
<coverage lines-covered="1" lines-valid="1" branches-covered="0" branches-valid="0">
  <packages>
    <package name="MyProject.Core">
      <classes>
        <class name="MyProject.Core.Calculator" filename="MyProject/Calculator.cs" line-rate="1.0" branch-rate="1.0">
           <methods/>
           <lines><line number="5" hits="1" branch="false" /></lines>
        </class>
        <class name="MyProject.Core.Calculator/&lt;&gt;c" filename="MyProject/Calculator.cs" line-rate="0" branch-rate="1.0">
           <methods/>
           <lines/>
        </class>
      </classes>
    </package>
  </packages>
</coverage>`,
			sourceFiles: map[string]string{
				resolvedSourcePath: `namespace MyProject.Core; public class Calculator { public int Add(int a, int b) => a + b; }`,
			},
			sourceDirs: []string{sourceDir},
			asserter: func(t *testing.T, result *parsers.ParserResult, err error) {
				require.NoError(t, err)
				require.Len(t, result.Assemblies, 1)
				assembly := result.Assemblies[0]
				require.Len(t, assembly.Classes, 1)
				assert.Equal(t, "MyProject.Core.Calculator", assembly.Classes[0].DisplayName)
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
			parser := cobertura.NewCoberturaParser(mockFS)

			// Act
			result, err := parser.Parse(reportPath, mockConfig)

			// Assert
			tc.asserter(t, result, err)
		})
	}
}

// Add this new function to internal/parsers/cobertura/parser_test.go

func TestCoberturaParser_SupportsFile(t *testing.T) {
	testCases := []struct {
		name        string
		fileName    string
		fileContent string
		shouldMatch bool
	}{
		{
			name:        "Valid Cobertura file with .xml extension",
			fileName:    "cobertura.xml",
			fileContent: `<?xml version="1.0" ?><coverage line-rate="1.0">...</coverage>`,
			shouldMatch: true,
		},
		{
			name:        "File with .xml extension but not Cobertura format",
			fileName:    "report.xml",
			fileContent: `<?xml version="1.0" ?><notcoverage></notcoverage>`,
			shouldMatch: false,
		},
		{
			name:        "File with correct content but wrong extension",
			fileName:    "cobertura.txt",
			fileContent: `<?xml version="1.0" ?><coverage></coverage>`,
			shouldMatch: false,
		},
		{
			name:        "Empty file with .xml extension",
			fileName:    "empty.xml",
			fileContent: "",
			shouldMatch: false,
		},
		{
			name:        "Malformed XML file",
			fileName:    "malformed.xml",
			fileContent: `<coverage><unclosed>`,
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

			// The parser's filereader isn't used by SupportsFile, which uses os.Open,
			// so we can pass a nil or mock reader.
			parser := cobertura.NewCoberturaParser(nil)

			// Act
			actual := parser.SupportsFile(filePath)

			// Assert
			assert.Equal(t, tc.shouldMatch, actual)
		})
	}
}
