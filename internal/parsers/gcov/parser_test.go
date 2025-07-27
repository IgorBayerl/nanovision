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

func Test_GCovParser_Parse(t *testing.T) {
	const reportFileName = "calculator.cpp.gcov"
	const sourceDir = "C:/www/AdlerCov/demo_projects/cpp/project/src"
	const sourceFileName = "calculator.cpp"
	const resolvedSourcePath = "C:/www/AdlerCov/demo_projects/cpp/project/src/calculator.cpp"

	const gcovReportContent = `        -:    0:Source:C:/www/AdlerCov/demo_projects/cpp/project/src/calculator.cpp
        -:    0:Graph:C:\www\AdlerCov\demo_projects\cpp\project\build\CMakeFiles\app_lib.dir\src\calculator.cpp.gcno
        -:    0:Data:C:\www\AdlerCov\demo_projects\cpp\project\build\CMakeFiles\app_lib.dir\src\calculator.cpp.gcda
        -:    0:Runs:2
function Calculator::add line 3
function Calculator::subtract line 7
function Calculator::multiply line 11
function Calculator::divide line 16
function Calculator::sign line 23
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

				require.Len(t, result.Assemblies, 1)
				assembly := result.Assemblies[0]

				require.Len(t, assembly.Classes, 1)
				class := assembly.Classes[0]

				assert.Equal(t, 5, class.TotalMethods)
				assert.Equal(t, 4, class.CoveredMethods)
				assert.Equal(t, 3, class.FullyCoveredMethods)

				require.Len(t, class.Methods, 5)

				divideMethod := testutil.FindMethod(t, class.Methods, "Calculator::divide")
				assert.Equal(t, 16, divideMethod.FirstLine)
				assert.Equal(t, 22, divideMethod.LastLine)
				assert.InDelta(t, 1.0, divideMethod.LineRate, 0.001)
				require.NotNil(t, divideMethod.BranchRate)
				assert.InDelta(t, 1.0, *divideMethod.BranchRate, 0.001)

				// Assert - NEW: Check for CodeElement creation
				require.Len(t, class.Files, 1)
				file := class.Files[0]
				require.Len(t, file.CodeElements, 5, "Should create 5 code elements for the 5 methods")

				var addCodeElement model.CodeElement
				for _, ce := range file.CodeElements {
					if ce.FullName == "Calculator::add" {
						addCodeElement = ce
						break
					}
				}
				require.NotNil(t, addCodeElement.FullName, "Could not find CodeElement for 'Calculator::add'")

				assert.Equal(t, "add", addCodeElement.Name)
				assert.Equal(t, "Calculator::add", addCodeElement.FullName)
				assert.Equal(t, model.MethodElementType, addCodeElement.Type)
				assert.Equal(t, 3, addCodeElement.FirstLine)
				require.NotNil(t, addCodeElement.CoverageQuota)
				assert.InDelta(t, 100.0, *addCodeElement.CoverageQuota, 0.001)
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
				assert.Empty(t, result.Assemblies)
				require.Len(t, result.UnresolvedSourceFiles, 1)
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

			mockReader := testutil.NewMockFilesystem("windows")
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
