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
	// --- Test Data ---
	const reportFileName = "calculator.cpp.gcov"
	const sourceDir = "C:/www/AdlerCov/demo_projects/cpp/project/src"
	const resolvedSourcePath = "C:/www/AdlerCov/demo_projects/cpp/project/src/calculator.cpp"

	const gcovReportContent = `        -:    0:Source:C:/www/AdlerCov/demo_projects/cpp/project/src/calculator.cpp
        -:    0:Graph:C:\www\AdlerCov\demo_projects\cpp\project\build\CMakeFiles\app_lib.dir\src\calculator.cpp.gcno
        -:    0:Data:C:\www\AdlerCov\demo_projects\cpp\project\build\CMakeFiles\app_lib.dir\src\calculator.cpp.gcda
        -:    0:Runs:2
function _ZN10Calculator3addEii called 2 returned 100% blocks executed 100%
        2:    3:int Calculator::add(int a, int b) {
        2:    4:    return a + b;
        -:    5:}
        -:    6:
function _ZN10Calculator8subtractEii called 1 returned 100% blocks executed 100%
        1:    7:int Calculator::subtract(int a, int b) {
        1:    8:    return a - b;
        -:    9:}
        -:   10:
function _ZN10Calculator8multiplyEii called 0 returned 0% blocks executed 0%
    #####:   11:int Calculator::multiply(int a, int b) {
        -:   12:    // This function is not called by any test.
    #####:   13:    return a * b;
        -:   14:}
        -:   15:
function _ZN10Calculator6divideEdd called 3 returned 67% blocks executed 88%
        3:   16:double Calculator::divide(double a, double b) {
        3:   17:    if (b == 0.0) {
branch  0 taken 33% (fallthrough)
branch  1 taken 67%
        1:   18:        throw std::invalid_argument("Division by zero is not allowed.");
        -:   19:    }
        2:   20:    return a / b;
        -:   21:}
        -:   22:
function _ZN10Calculator4signEi called 2 returned 100% blocks executed 83%
        2:   23:int Calculator::sign(int x) {
        2:   24:    if (x > 0) {
branch  0 taken 50% (fallthrough)
branch  1 taken 50%
        1:   25:        return 1;
        1:   26:    } else if (x < 0) {
branch  0 taken 100% (fallthrough)
branch  1 taken 0%
        1:   27:        return -1;
        -:   28:    } else {
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
				class := result.Assemblies[0].Classes[0]

				assert.Equal(t, 5, class.TotalMethods)
				assert.Equal(t, 4, class.CoveredMethods)
				assert.Equal(t, 3, class.FullyCoveredMethods)

				require.Len(t, class.Methods, 5)

				// Assert against the final, unmangled names
				addMethod := testutil.FindMethod(t, class.Methods, "Calculator::add")
				assert.Equal(t, "(int a, int b)", addMethod.Signature)
				assert.Equal(t, "Calculator::add(int a, int b)", addMethod.DisplayName)
				assert.Equal(t, 3, addMethod.FirstLine)
				assert.Equal(t, 6, addMethod.LastLine) // Ends before subtract starts on line 7
				assert.InDelta(t, 1.0, addMethod.LineRate, 0.001)

				divideMethod := testutil.FindMethod(t, class.Methods, "Calculator::divide")
				assert.Equal(t, "(double a, double b)", divideMethod.Signature)
				assert.Equal(t, 16, divideMethod.FirstLine)
				assert.Equal(t, 22, divideMethod.LastLine) // Ends before sign starts on line 23
				assert.InDelta(t, 1.0, divideMethod.LineRate, 0.001)
				require.NotNil(t, divideMethod.BranchRate)
				assert.InDelta(t, 1.0, *divideMethod.BranchRate, 0.001)

				// CodeElements were created correctly
				require.Len(t, class.Files, 1)
				file := class.Files[0]
				require.Len(t, file.CodeElements, 5)

				var signCodeElement model.CodeElement
				for _, ce := range file.CodeElements {
					if ce.FullName == "Calculator::sign(int x)" {
						signCodeElement = ce
						break
					}
				}
				require.NotEmpty(t, signCodeElement.FullName, "Could not find CodeElement for 'Calculator::sign'")
				assert.Equal(t, "Calculator::sign(...)", signCodeElement.Name)
				assert.Equal(t, 23, signCodeElement.FirstLine)
				require.NotNil(t, signCodeElement.CoverageQuota)
				assert.InDelta(t, 83.3, *signCodeElement.CoverageQuota, 0.1)
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
