package gcov_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/IgorBayerl/AdlerCov/internal/model"
	"github.com/IgorBayerl/AdlerCov/internal/parsers"
	"github.com/IgorBayerl/AdlerCov/internal/parsers/gcov"
	"github.com/IgorBayerl/AdlerCov/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGCovParser_Parse(t *testing.T) {
	// --- Test Data ---
	const reportFileName = "calculator.cpp.gcov"
	const sourceDir = "C:/www/AdlerCov/demo_projects/cpp/project/src"
	const sourceFileName = "calculator.cpp"
	const resolvedSourcePath = "C:/www/AdlerCov/demo_projects/cpp/project/src/calculator.cpp"

	const gcovReportContent = `        -:    0:Source:C:/www/AdlerCov/demo_projects/cpp/project/src/calculator.cpp
        -:    0:Graph:C:\www\AdlerCov\demo_projects\cpp\project\build\CMakeFiles\app_lib.dir\src\calculator.cpp.gcno
        -:    0:Data:C:\www\AdlerCov\demo_projects\cpp\project\build\CMakeFiles\app_lib.dir\src\calculator.cpp.gcda
        -:    0:Runs:2
        -:    1:#include "calculator.h"
        -:    2:
        2:    3:int Calculator::add(int a, int b) {
        2:    4:    return a + b;
        -:    5:}
        -:    6:
        1:    7:int Calculator::subtract(int a, int b) {
        1:    8:    return a - b;
        -:    9:}
        -:   10:
    #####:   11:int Calculator::multiply(int a, int b) {
        -:   12:    // This function is not called by any test.
    #####:   13:    return a * b;
        -:   14:}
        -:   15:
        3:   16:double Calculator::divide(double a, double b) {
        3:   17:    if (b == 0.0) {
branch  0 taken 1
branch  1 taken 2
        1:   18:        throw std::invalid_argument("Division by zero is not allowed.");
        -:   19:    }
        2:   20:    return a / b;
        -:   21:}
        -:   22:
        2:   23:int Calculator::sign(int x) {
        2:   24:    if (x > 0) {
branch  0 taken 1
branch  1 taken 1
        1:   25:        return 1;
        1:   26:    } else if (x < 0) {
branch  0 taken 1
branch  1 taken 0
        1:   27:        return -1;
        -:   28:    } else {
        -:   29:        // This branch will be deliberately missed by the tests.
    #####:   30:        return 0;
        -:   31:    }
        -:   32:}`

	const calculatorCppContent = `#include "calculator.h"

int Calculator::add(int a, int b) {
    return a + b;
}

int Calculator::subtract(int a, int b) {
    return a - b;
}

int Calculator::multiply(int a, int b) {
    // This function is not called by any test.
    return a * b;
}

double Calculator::divide(double a, double b) {
    if (b == 0.0) {
        throw std::invalid_argument("Division by zero is not allowed.");
    }
    return a / b;
}

int Calculator::sign(int x) {
    if (x > 0) {
        return 1;
    } else if (x < 0) {
        return -1;
    } else {
        // This branch will be deliberately missed by the tests.
        return 0;
    }
}`

	testCases := []struct {
		name          string
		reportContent string
		sourceFiles   map[string]string
		sourceDirs    []string
		asserter      func(t *testing.T, result *parsers.ParserResult, err error)
	}{
		{
			name:          "Golden Path - Valid report with branch coverage",
			reportContent: gcovReportContent,
			sourceFiles: map[string]string{
				resolvedSourcePath: calculatorCppContent,
			},
			sourceDirs: []string{sourceDir},
			asserter: func(t *testing.T, result *parsers.ParserResult, err error) {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Empty(t, result.UnresolvedSourceFiles)
				assert.Equal(t, "GCov", result.ParserName)
				assert.True(t, result.SupportsBranchCoverage, "Branch data was present, so this should be true")

				require.Len(t, result.Assemblies, 1)
				assembly := result.Assemblies[0]
				assert.Equal(t, "Default", assembly.Name)
				assert.Equal(t, 13, assembly.LinesCovered) // CORRECTED: 13 lines have hits > 0
				assert.Equal(t, 16, assembly.LinesValid)   // CORRECTED: 16 lines are not marked with '-'
				require.NotNil(t, assembly.BranchesCovered)
				assert.Equal(t, 5, *assembly.BranchesCovered) // CORRECTED: 2+2+1=5
				require.NotNil(t, assembly.BranchesValid)
				assert.Equal(t, 6, *assembly.BranchesValid) // CORRECTED: 2+2+2=6

				require.Len(t, assembly.Classes, 1)
				class := assembly.Classes[0]
				assert.Equal(t, sourceFileName, class.DisplayName)

				require.Len(t, class.Files, 1)
				file := class.Files[0]
				// CORRECTED: Normalize path separators for comparison
				assert.Equal(t, filepath.ToSlash(resolvedSourcePath), filepath.ToSlash(file.Path))
				assert.Equal(t, 32, file.TotalLines)

				line4 := testutil.FindLine(t, file.Lines, 4)
				assert.Equal(t, 2, line4.Hits)
				assert.Equal(t, model.Covered, line4.LineVisitStatus)

				line13 := testutil.FindLine(t, file.Lines, 13)
				assert.Equal(t, 0, line13.Hits)
				assert.Equal(t, model.NotCovered, line13.LineVisitStatus)

				// This is now the most important assertion for branches
				line17 := testutil.FindLine(t, file.Lines, 17)
				assert.Equal(t, 3, line17.Hits)
				assert.True(t, line17.IsBranchPoint)
				assert.Equal(t, 2, line17.CoveredBranches)             // CORRECTED: Both branches taken
				assert.Equal(t, 2, line17.TotalBranches)               // CORRECTED: Two branches on this line
				assert.Equal(t, model.Covered, line17.LineVisitStatus) // Both covered -> Covered, not Partially

				line30 := testutil.FindLine(t, file.Lines, 30)
				assert.Equal(t, 0, line30.Hits)
				assert.Equal(t, model.NotCovered, line30.LineVisitStatus)
			},
		},
		{
			name:          "Source File Not Found",
			reportContent: gcovReportContent,
			sourceFiles:   map[string]string{},
			sourceDirs:    []string{"/another/dir"},
			asserter: func(t *testing.T, result *parsers.ParserResult, err error) {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Empty(t, result.Assemblies, "No assemblies should be created if the source file is not found")
				require.Len(t, result.UnresolvedSourceFiles, 1)
				// CORRECTED: Normalize path for assertion
				assert.Equal(t, filepath.ToSlash(resolvedSourcePath), filepath.ToSlash(result.UnresolvedSourceFiles[0]))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			reportPath := filepath.Join(tmpDir, reportFileName)
			err := os.WriteFile(reportPath, []byte(tc.reportContent), 0644)
			require.NoError(t, err)

			mockReader := testutil.NewMockFilesystem("windows") // Simulating windows
			for _, dir := range tc.sourceDirs {
				mockReader.AddDir(dir)
			}
			for path, content := range tc.sourceFiles {
				mockReader.AddFile(path, content)
			}

			mockConfig := testutil.NewTestConfig(tc.sourceDirs)
			parser := gcov.NewGCovParser(mockReader)

			result, err := parser.Parse(reportPath, mockConfig)

			tc.asserter(t, result, err)
		})
	}
}
