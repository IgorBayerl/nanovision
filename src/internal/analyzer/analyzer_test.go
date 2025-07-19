package analyzer_test

import (
	"fmt"
	"log/slog"
	"sort"
	"testing"
	"time"

	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/analyzer"
	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/filtering"
	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/model"
	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/parsers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockMergerConfig implements MergerConfig for testing
type mockMergerConfig struct {
	sourceDirs      []string
	assemblyFilters filtering.IFilter
	logger          *slog.Logger
}

func (m *mockMergerConfig) SourceDirectories() []string        { return m.sourceDirs }
func (m *mockMergerConfig) AssemblyFilters() filtering.IFilter { return m.assemblyFilters }
func (m *mockMergerConfig) Logger() *slog.Logger               { return m.logger }

// =============================================================================
// ERROR HANDLING TESTS
// =============================================================================

func TestMergeParserResults_WhenNoResults_ShouldReturnError(t *testing.T) {
	// Arrange
	var results []*parsers.ParserResult
	config := &mockMergerConfig{logger: slog.Default()}

	// Act
	summary, err := analyzer.MergeParserResults(results, config)

	// Assert
	require.Error(t, err, "Expected error when no results provided")
	assert.Nil(t, summary, "Expected nil summary when error occurs")
}

func TestMergeParserResults_WhenEmptySlice_ShouldReturnError(t *testing.T) {
	// Arrange
	results := []*parsers.ParserResult{}
	config := &mockMergerConfig{logger: slog.Default()}

	// Act
	summary, err := analyzer.MergeParserResults(results, config)

	// Assert
	require.Error(t, err, "Expected error when empty slice provided")
	assert.Nil(t, summary, "Expected nil summary when error occurs")
}

// =============================================================================
// SINGLE RESULT TESTS
// =============================================================================

func TestMergeParserResults_WhenSingleResult_ShouldPreserveAllData(t *testing.T) {
	// Arrange
	timestamp := time.Now()
	branchesCovered := 10
	branchesValid := 20

	results := []*parsers.ParserResult{
		{
			ParserName:        "TestParser",
			MinimumTimeStamp:  &timestamp,
			SourceDirectories: []string{"/src/app"},
			Assemblies: []model.Assembly{
				{
					Name:            "TestAssembly",
					LinesCovered:    50,
					LinesValid:      100,
					BranchesCovered: &branchesCovered,
					BranchesValid:   &branchesValid,
					Classes: []model.Class{
						{
							Name:         "TestClass",
							LinesCovered: 30,
							LinesValid:   60,
							Files: []model.CodeFile{
								{Path: "/src/app/test.go", TotalLines: 200},
							},
						},
					},
				},
			},
		},
	}
	config := &mockMergerConfig{logger: slog.Default()}

	// Act
	summary, err := analyzer.MergeParserResults(results, config)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "TestParser", summary.ParserName)
	assert.Equal(t, timestamp.Unix(), summary.Timestamp)
	assert.Equal(t, []string{"/src/app"}, summary.SourceDirs)
	assert.Equal(t, 50, summary.LinesCovered)
	require.NotNil(t, summary.BranchesCovered)
	assert.Equal(t, 10, *summary.BranchesCovered)
	assert.Equal(t, 200, summary.TotalLines)
}

func TestMergeParserResults_WhenSingleResultWithoutBranches_ShouldNotIncludeBranchData(t *testing.T) {
	// Arrange
	results := []*parsers.ParserResult{
		{
			ParserName: "NoBranchParser",
			Assemblies: []model.Assembly{
				{
					Name:         "Assembly1",
					LinesCovered: 25,
					LinesValid:   50,
					// No branch data
				},
			},
		},
	}
	config := &mockMergerConfig{logger: slog.Default()}

	// Act
	summary, err := analyzer.MergeParserResults(results, config)

	// Assert
	require.NoError(t, err)
	assert.Nil(t, summary.BranchesCovered, "Expected no branch coverage data when source has none")
	assert.Nil(t, summary.BranchesValid, "Expected no branch validity data when source has none")
}

// =============================================================================
// MULTIPLE RESULTS TESTS
// =============================================================================

func TestMergeParserResults_WhenMultipleResultsWithSameName_ShouldUseSingleParserName(t *testing.T) {
	// Arrange
	results := []*parsers.ParserResult{
		{ParserName: "SameParser"},
		{ParserName: "SameParser"},
		{ParserName: "SameParser"},
	}
	config := &mockMergerConfig{logger: slog.Default()}

	// Act
	summary, err := analyzer.MergeParserResults(results, config)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "SameParser", summary.ParserName)
}

func TestMergeParserResults_WhenMultipleResultsWithDifferentNames_ShouldUseMultiReport(t *testing.T) {
	// Arrange
	results := []*parsers.ParserResult{
		{ParserName: "Parser1"},
		{ParserName: "Parser2"},
		{ParserName: "Parser3"},
	}
	config := &mockMergerConfig{logger: slog.Default()}

	// Act
	summary, err := analyzer.MergeParserResults(results, config)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "MultiReport", summary.ParserName)
}

func TestMergeParserResults_WhenMultipleResultsWithEmptyNames_ShouldUseUnknown(t *testing.T) {
	// Arrange
	results := []*parsers.ParserResult{
		{ParserName: ""},
		{ParserName: ""},
	}
	config := &mockMergerConfig{logger: slog.Default()}

	// Act
	summary, err := analyzer.MergeParserResults(results, config)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "Unknown", summary.ParserName)
}

// =============================================================================
// TIMESTAMP TESTS
// =============================================================================

func TestMergeParserResults_WhenMultipleTimestamps_ShouldUseEarliest(t *testing.T) {
	// Arrange
	early := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	middle := time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC)
	late := time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC)

	results := []*parsers.ParserResult{
		{ParserName: "Test", MinimumTimeStamp: &late},
		{ParserName: "Test", MinimumTimeStamp: &early},
		{ParserName: "Test", MinimumTimeStamp: &middle},
	}
	config := &mockMergerConfig{logger: slog.Default()}

	// Act
	summary, err := analyzer.MergeParserResults(results, config)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, early.Unix(), summary.Timestamp)
}

func TestMergeParserResults_WhenNoTimestamps_ShouldHaveZeroTimestamp(t *testing.T) {
	// Arrange
	results := []*parsers.ParserResult{
		{ParserName: "Test"},
		{ParserName: "Test"},
	}
	config := &mockMergerConfig{logger: slog.Default()}

	// Act
	summary, err := analyzer.MergeParserResults(results, config)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, int64(0), summary.Timestamp)
}

func TestMergeParserResults_WhenMixedTimestamps_ShouldIgnoreNilTimestamps(t *testing.T) {
	// Arrange
	validTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	results := []*parsers.ParserResult{
		{ParserName: "Test", MinimumTimeStamp: nil},
		{ParserName: "Test", MinimumTimeStamp: &validTime},
		{ParserName: "Test", MinimumTimeStamp: nil},
	}
	config := &mockMergerConfig{logger: slog.Default()}

	// Act
	summary, err := analyzer.MergeParserResults(results, config)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, validTime.Unix(), summary.Timestamp)
}

// =============================================================================
// SOURCE DIRECTORIES TESTS
// =============================================================================

func TestMergeParserResults_WhenDuplicateSourceDirs_ShouldDeduplicateDirectories(t *testing.T) {
	// Arrange
	results := []*parsers.ParserResult{
		{ParserName: "Test", SourceDirectories: []string{"/src/app", "/src/lib"}},
		{ParserName: "Test", SourceDirectories: []string{"/src/app", "/src/test"}},
		{ParserName: "Test", SourceDirectories: []string{"/src/lib", "/src/docs"}},
	}
	config := &mockMergerConfig{logger: slog.Default()}

	// Act
	summary, err := analyzer.MergeParserResults(results, config)

	// Assert
	require.NoError(t, err)
	expectedDirs := []string{"/src/app", "/src/lib", "/src/test", "/src/docs"}
	assert.Len(t, summary.SourceDirs, len(expectedDirs))

	// Check all expected directories are present (order doesn't matter)
	for _, expectedDir := range expectedDirs {
		assert.Contains(t, summary.SourceDirs, expectedDir)
	}
}

func TestMergeParserResults_WhenEmptySourceDirs_ShouldHandleGracefully(t *testing.T) {
	// Arrange
	results := []*parsers.ParserResult{
		{ParserName: "Test", SourceDirectories: []string{}},
		{ParserName: "Test", SourceDirectories: nil},
	}
	config := &mockMergerConfig{logger: slog.Default()}

	// Act
	summary, err := analyzer.MergeParserResults(results, config)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, summary.SourceDirs, "Expected non-nil source directories slice")
	assert.Empty(t, summary.SourceDirs, "Expected empty source directories")
}

// =============================================================================
// ASSEMBLY MERGING TESTS
// =============================================================================

func TestMergeParserResults_WhenDifferentAssemblies_ShouldCombineAllAssemblies(t *testing.T) {
	// Arrange
	results := []*parsers.ParserResult{
		{
			ParserName: "Test",
			Assemblies: []model.Assembly{
				{Name: "Assembly1", LinesCovered: 10, LinesValid: 20},
				{Name: "Assembly2", LinesCovered: 15, LinesValid: 25},
			},
		},
		{
			ParserName: "Test",
			Assemblies: []model.Assembly{
				{Name: "Assembly3", LinesCovered: 5, LinesValid: 10},
			},
		},
	}
	config := &mockMergerConfig{logger: slog.Default()}

	// Act
	summary, err := analyzer.MergeParserResults(results, config)

	// Assert
	require.NoError(t, err)
	assert.Len(t, summary.Assemblies, 3)

	// Check assemblies are sorted by name
	expectedNames := []string{"Assembly1", "Assembly2", "Assembly3"}
	for i, expected := range expectedNames {
		assert.Equal(t, expected, summary.Assemblies[i].Name)
	}
}

func TestMergeParserResults_WhenSameAssemblyInMultipleResults_ShouldMergeStatistics(t *testing.T) {
	// Arrange
	branches1 := 5
	branchesValid1 := 10
	branches2 := 8
	branchesValid2 := 12

	results := []*parsers.ParserResult{
		{
			ParserName: "Test",
			Assemblies: []model.Assembly{
				{
					Name:            "SharedAssembly",
					LinesCovered:    10,
					LinesValid:      20,
					BranchesCovered: &branches1,
					BranchesValid:   &branchesValid1,
					Classes: []model.Class{
						{Name: "Class1", LinesCovered: 5, LinesValid: 10},
					},
				},
			},
		},
		{
			ParserName: "Test",
			Assemblies: []model.Assembly{
				{
					Name:            "SharedAssembly",
					LinesCovered:    15,
					LinesValid:      25,
					BranchesCovered: &branches2,
					BranchesValid:   &branchesValid2,
					Classes: []model.Class{
						{Name: "Class2", LinesCovered: 8, LinesValid: 15},
					},
				},
			},
		},
	}
	config := &mockMergerConfig{logger: slog.Default()}

	// Act
	summary, err := analyzer.MergeParserResults(results, config)

	// Assert
	require.NoError(t, err)
	assert.Len(t, summary.Assemblies, 1, "Expected 1 merged assembly")

	asm := summary.Assemblies[0]
	assert.Equal(t, 25, asm.LinesCovered, "Expected merged lines covered")
	assert.Equal(t, 45, asm.LinesValid, "Expected merged lines valid")
	require.NotNil(t, asm.BranchesCovered)
	assert.Equal(t, 13, *asm.BranchesCovered, "Expected merged branches covered")
	require.NotNil(t, asm.BranchesValid)
	assert.Equal(t, 22, *asm.BranchesValid, "Expected merged branches valid")
	assert.Len(t, asm.Classes, 2, "Expected 2 classes after merge")
}

func TestMergeParserResults_WhenMixedBranchData_ShouldHandlePartialBranchInfo(t *testing.T) {
	// Arrange
	branches := 10
	branchesValid := 15

	results := []*parsers.ParserResult{
		{
			ParserName: "Test",
			Assemblies: []model.Assembly{
				{
					Name:            "TestAssembly",
					LinesCovered:    10,
					LinesValid:      20,
					BranchesCovered: &branches,
					BranchesValid:   &branchesValid,
				},
			},
		},
		{
			ParserName: "Test",
			Assemblies: []model.Assembly{
				{
					Name:         "TestAssembly",
					LinesCovered: 5,
					LinesValid:   10,
					// No branch data
				},
			},
		},
	}
	config := &mockMergerConfig{logger: slog.Default()}

	// Act
	summary, err := analyzer.MergeParserResults(results, config)

	// Assert
	require.NoError(t, err)

	asm := summary.Assemblies[0]
	require.NotNil(t, asm.BranchesCovered)
	assert.Equal(t, 10, *asm.BranchesCovered)
	require.NotNil(t, asm.BranchesValid)
	assert.Equal(t, 15, *asm.BranchesValid)
}

// =============================================================================
// GLOBAL STATISTICS TESTS
// =============================================================================

func TestMergeParserResults_WhenCalculatingGlobalStats_ShouldSumAllAssemblyStats(t *testing.T) {
	// Arrange
	branches1 := 5
	branchesValid1 := 10
	branches2 := 8
	branchesValid2 := 12

	results := []*parsers.ParserResult{
		{
			ParserName: "Test",
			Assemblies: []model.Assembly{
				{
					Name:            "Assembly1",
					LinesCovered:    10,
					LinesValid:      20,
					BranchesCovered: &branches1,
					BranchesValid:   &branchesValid1,
				},
				{
					Name:            "Assembly2",
					LinesCovered:    15,
					LinesValid:      25,
					BranchesCovered: &branches2,
					BranchesValid:   &branchesValid2,
				},
			},
		},
	}
	config := &mockMergerConfig{logger: slog.Default()}

	// Act
	summary, err := analyzer.MergeParserResults(results, config)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, 25, summary.LinesCovered)
	assert.Equal(t, 45, summary.LinesValid)
	require.NotNil(t, summary.BranchesCovered)
	assert.Equal(t, 13, *summary.BranchesCovered)
	require.NotNil(t, summary.BranchesValid)
	assert.Equal(t, 22, *summary.BranchesValid)
}

func TestMergeParserResults_WhenCalculatingTotalLines_ShouldCountUniqueFilesOnly(t *testing.T) {
	// Arrange
	results := []*parsers.ParserResult{
		{
			ParserName: "Test",
			Assemblies: []model.Assembly{
				{
					Name: "Assembly1",
					Classes: []model.Class{
						{
							Name: "Class1",
							Files: []model.CodeFile{
								{Path: "/app/file1.go", TotalLines: 100},
								{Path: "/app/file2.go", TotalLines: 200},
							},
						},
					},
				},
				{
					Name: "Assembly2",
					Classes: []model.Class{
						{
							Name: "Class2",
							Files: []model.CodeFile{
								{Path: "/app/file1.go", TotalLines: 100}, // Duplicate
								{Path: "/app/file3.go", TotalLines: 150},
							},
						},
					},
				},
			},
		},
	}
	config := &mockMergerConfig{logger: slog.Default()}

	// Act
	summary, err := analyzer.MergeParserResults(results, config)

	// Assert
	require.NoError(t, err)
	// Should be 100 + 200 + 150 = 450 (file1.go counted only once)
	assert.Equal(t, 450, summary.TotalLines)
}

func TestMergeParserResults_WhenFilesHaveZeroLines_ShouldIgnoreZeroLineFiles(t *testing.T) {
	// Arrange
	results := []*parsers.ParserResult{
		{
			ParserName: "Test",
			Assemblies: []model.Assembly{
				{
					Name: "Assembly1",
					Classes: []model.Class{
						{
							Name: "Class1",
							Files: []model.CodeFile{
								{Path: "/app/file1.go", TotalLines: 0}, // Should be ignored
								{Path: "/app/file2.go", TotalLines: 100},
								{Path: "/app/file3.go", TotalLines: 0}, // Should be ignored
							},
						},
					},
				},
			},
		},
	}
	config := &mockMergerConfig{logger: slog.Default()}

	// Act
	summary, err := analyzer.MergeParserResults(results, config)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, 100, summary.TotalLines)
}

// =============================================================================
// EDGE CASES AND BOUNDARY CONDITIONS
// =============================================================================

func TestMergeParserResults_WhenAssembliesHaveNoClasses_ShouldHandleGracefully(t *testing.T) {
	// Arrange
	results := []*parsers.ParserResult{
		{
			ParserName: "Test",
			Assemblies: []model.Assembly{
				{
					Name:         "EmptyAssembly",
					LinesCovered: 10,
					LinesValid:   20,
					Classes:      []model.Class{}, // Empty classes
				},
			},
		},
	}
	config := &mockMergerConfig{logger: slog.Default()}

	// Act
	summary, err := analyzer.MergeParserResults(results, config)

	// Assert
	require.NoError(t, err)
	assert.Len(t, summary.Assemblies, 1)
	assert.Equal(t, 0, summary.TotalLines, "Expected total lines 0 for empty classes")
}

func TestMergeParserResults_WhenClassesHaveNoFiles_ShouldHandleGracefully(t *testing.T) {
	// Arrange
	results := []*parsers.ParserResult{
		{
			ParserName: "Test",
			Assemblies: []model.Assembly{
				{
					Name: "Assembly1",
					Classes: []model.Class{
						{
							Name:  "EmptyClass",
							Files: []model.CodeFile{}, // Empty files
						},
					},
				},
			},
		},
	}
	config := &mockMergerConfig{logger: slog.Default()}

	// Act
	summary, err := analyzer.MergeParserResults(results, config)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, 0, summary.TotalLines, "Expected total lines 0 for empty files")
}

func TestMergeParserResults_WhenLargeNumbers_ShouldHandleIntegerOverflow(t *testing.T) {
	// Arrange
	largeNumber := 1000000000 // 1 billion
	hugeBranches := largeNumber

	results := []*parsers.ParserResult{
		{
			ParserName: "Test",
			Assemblies: []model.Assembly{
				{
					Name:            "LargeAssembly1",
					LinesCovered:    largeNumber,
					LinesValid:      largeNumber,
					BranchesCovered: &hugeBranches,
					BranchesValid:   &hugeBranches,
				},
				{
					Name:            "LargeAssembly2",
					LinesCovered:    largeNumber,
					LinesValid:      largeNumber,
					BranchesCovered: &hugeBranches,
					BranchesValid:   &hugeBranches,
				},
			},
		},
	}
	config := &mockMergerConfig{logger: slog.Default()}

	// Act
	summary, err := analyzer.MergeParserResults(results, config)

	// Assert
	require.NoError(t, err)
	expectedTotal := largeNumber * 2
	assert.Equal(t, expectedTotal, summary.LinesCovered)
	require.NotNil(t, summary.BranchesCovered)
	assert.Equal(t, expectedTotal, *summary.BranchesCovered)
}

// =============================================================================
// BEHAVIOR VALIDATION TESTS
// =============================================================================

func TestMergeParserResults_WhenComplexScenario_ShouldProduceConsistentResults(t *testing.T) {
	// Arrange - Complex scenario with multiple parsers, overlapping assemblies, mixed data
	time1 := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	time2 := time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC)

	branches1 := 10
	branchesValid1 := 20
	branches2 := 15
	branchesValid2 := 25

	results := []*parsers.ParserResult{
		{
			ParserName:        "Parser1",
			MinimumTimeStamp:  &time2,
			SourceDirectories: []string{"/src", "/test"},
			Assemblies: []model.Assembly{
				{
					Name:            "SharedAssembly",
					LinesCovered:    50,
					LinesValid:      100,
					BranchesCovered: &branches1,
					BranchesValid:   &branchesValid1,
					Classes: []model.Class{
						{
							Name: "Class1",
							Files: []model.CodeFile{
								{Path: "/src/shared.go", TotalLines: 200},
							},
						},
					},
				},
				{
					Name:         "UniqueAssembly1",
					LinesCovered: 25,
					LinesValid:   50,
					Classes: []model.Class{
						{
							Name: "UniqueClass1",
							Files: []model.CodeFile{
								{Path: "/src/unique1.go", TotalLines: 150},
							},
						},
					},
				},
			},
		},
		{
			ParserName:        "Parser2",
			MinimumTimeStamp:  &time1,                   // Earlier timestamp
			SourceDirectories: []string{"/src", "/lib"}, // Overlap with first
			Assemblies: []model.Assembly{
				{
					Name:            "SharedAssembly", // Same name as first
					LinesCovered:    30,
					LinesValid:      60,
					BranchesCovered: &branches2,
					BranchesValid:   &branchesValid2,
					Classes: []model.Class{
						{
							Name: "Class2",
							Files: []model.CodeFile{
								{Path: "/src/shared.go", TotalLines: 200}, // Duplicate file
								{Path: "/lib/helper.go", TotalLines: 100},
							},
						},
					},
				},
				{
					Name:         "UniqueAssembly2",
					LinesCovered: 40,
					LinesValid:   80,
					// No branch data
					Classes: []model.Class{
						{
							Name: "UniqueClass2",
							Files: []model.CodeFile{
								{Path: "/lib/utils.go", TotalLines: 75},
							},
						},
					},
				},
			},
		},
	}
	config := &mockMergerConfig{logger: slog.Default()}

	// Act
	summary, err := analyzer.MergeParserResults(results, config)

	// Assert
	require.NoError(t, err)

	// Verify parser name logic
	assert.Equal(t, "MultiReport", summary.ParserName)

	// Verify timestamp logic (should use earliest)
	assert.Equal(t, time1.Unix(), summary.Timestamp)

	// Verify source directories deduplication
	expectedSourceDirs := 3 // /src, /test, /lib
	assert.Len(t, summary.SourceDirs, expectedSourceDirs)

	// Verify assembly merging
	assert.Len(t, summary.Assemblies, 3, "Expected 3 assemblies (1 merged, 2 unique)")

	// Find the merged SharedAssembly
	var sharedAsm *model.Assembly
	for i := range summary.Assemblies {
		if summary.Assemblies[i].Name == "SharedAssembly" {
			sharedAsm = &summary.Assemblies[i]
			break
		}
	}
	require.NotNil(t, sharedAsm, "SharedAssembly not found in merged results")

	// Verify merged assembly statistics
	assert.Equal(t, 80, sharedAsm.LinesCovered, "Expected merged assembly lines covered") // 50 + 30
	assert.Equal(t, 160, sharedAsm.LinesValid, "Expected merged assembly lines valid")    // 100 + 60
	require.NotNil(t, sharedAsm.BranchesCovered)
	assert.Equal(t, 25, *sharedAsm.BranchesCovered, "Expected merged assembly branches covered") // 10 + 15
	assert.Len(t, sharedAsm.Classes, 2, "Expected merged assembly to have 2 classes")            // Class1 + Class2

	// Verify global statistics
	expectedLinesCovered := 145 // 80 (shared) + 25 (unique1) + 40 (unique2)
	assert.Equal(t, expectedLinesCovered, summary.LinesCovered)

	expectedLinesValid := 290 // 160 (shared) + 50 (unique1) + 80 (unique2)
	assert.Equal(t, expectedLinesValid, summary.LinesValid)

	// Verify total lines (unique files only)
	expectedTotalLines := 525 // 200 (shared.go) + 150 (unique1.go) + 100 (helper.go) + 75 (utils.go)
	assert.Equal(t, expectedTotalLines, summary.TotalLines)

	// Verify branch data presence
	assert.NotNil(t, summary.BranchesCovered, "Expected branch coverage data in final summary")
	assert.NotNil(t, summary.BranchesValid, "Expected branch validity data in final summary")
}

// =============================================================================
// PERFORMANCE AND MEMORY TESTS
// =============================================================================

func TestMergeParserResults_WhenManyAssemblies_ShouldPerformEfficiently(t *testing.T) {
	// Arrange - Create many assemblies to test performance
	results := []*parsers.ParserResult{
		{
			ParserName: "Test",
			Assemblies: make([]model.Assembly, 1000), // 1000 assemblies
		},
	}

	// Initialize assemblies
	for i := 0; i < 1000; i++ {
		results[0].Assemblies[i] = model.Assembly{
			Name:         fmt.Sprintf("Assembly%d", i),
			LinesCovered: i + 1,
			LinesValid:   (i + 1) * 2,
		}
	}

	config := &mockMergerConfig{logger: slog.Default()}

	// Act
	summary, err := analyzer.MergeParserResults(results, config)

	// Assert
	require.NoError(t, err)
	assert.Len(t, summary.Assemblies, 1000)

	// Verify assemblies are sorted by checking if the slice is sorted
	isSorted := sort.SliceIsSorted(summary.Assemblies, func(i, j int) bool {
		return summary.Assemblies[i].Name < summary.Assemblies[j].Name
	})
	assert.True(t, isSorted, "Assemblies should be sorted by name")
}

func TestMergeParserResults_WhenManySourceDirectories_ShouldDeduplicateEfficiently(t *testing.T) {
	// Arrange - Create many overlapping source directories
	results := make([]*parsers.ParserResult, 10)
	for i := 0; i < 10; i++ {
		dirs := make([]string, 100)
		for j := 0; j < 100; j++ {
			dirs[j] = fmt.Sprintf("/src%d", j%50) // Create overlap
		}
		results[i] = &parsers.ParserResult{
			ParserName:        "Test",
			SourceDirectories: dirs,
		}
	}

	config := &mockMergerConfig{logger: slog.Default()}

	// Act
	summary, err := analyzer.MergeParserResults(results, config)

	// Assert
	require.NoError(t, err)
	assert.Len(t, summary.SourceDirs, 50, "Expected 50 unique source directories")
}

// =============================================================================
// HELPER FUNCTION TESTS
// =============================================================================

func TestMergeParserResults_WhenAssembliesAreSorted_ShouldMaintainConsistentOrder(t *testing.T) {
	// Arrange - Create assemblies in random order
	results := []*parsers.ParserResult{
		{
			ParserName: "Test",
			Assemblies: []model.Assembly{
				{Name: "Zebra", LinesCovered: 1, LinesValid: 2},
				{Name: "Alpha", LinesCovered: 3, LinesValid: 4},
				{Name: "Bravo", LinesCovered: 5, LinesValid: 6},
			},
		},
	}
	config := &mockMergerConfig{logger: slog.Default()}

	// Act
	summary, err := analyzer.MergeParserResults(results, config)

	// Assert
	require.NoError(t, err)

	// Extract just the names to check the order
	var actualNames []string
	for _, asm := range summary.Assemblies {
		actualNames = append(actualNames, asm.Name)
	}

	expectedOrder := []string{"Alpha", "Bravo", "Zebra"}
	assert.Equal(t, expectedOrder, actualNames, "Assemblies should be sorted alphabetically by name")
}

func TestMergeParserResults_WhenMultipleCallsWithSameData_ShouldProduceIdenticalResults(t *testing.T) {
	// Arrange
	branches := 5
	branchesValid := 10

	results := []*parsers.ParserResult{
		{
			ParserName: "Test",
			Assemblies: []model.Assembly{
				{
					Name:            "TestAssembly",
					LinesCovered:    25,
					LinesValid:      50,
					BranchesCovered: &branches,
					BranchesValid:   &branchesValid,
					Classes: []model.Class{
						{
							Name: "TestClass",
							Files: []model.CodeFile{
								{Path: "/test/file.go", TotalLines: 100},
							},
						},
					},
				},
			},
		},
	}
	config := &mockMergerConfig{logger: slog.Default()}

	// Act - Call multiple times
	summary1, err1 := analyzer.MergeParserResults(results, config)
	summary2, err2 := analyzer.MergeParserResults(results, config)

	// Assert
	require.NoError(t, err1)
	require.NoError(t, err2)

	assert.Equal(t, summary1, summary2, "Multiple calls with the same data should produce identical results")
}

// =============================================================================
// DEEP MERGE TESTS
// =============================================================================

// TestMergeParserResults_WhenDeepMerging_ShouldAggregateCorrectly verifies the end-to-end process
// when assemblies and their nested classes are merged from multiple reports.
func TestMergeParserResults_WhenDeepMerging_ShouldAggregateCorrectly(t *testing.T) {
	// Arrange
	// Report 1: Contains a shared assembly with a shared class and a unique class.
	result1 := &parsers.ParserResult{
		ParserName: "Test",
		Assemblies: []model.Assembly{
			{
				Name:         "SharedAssembly",
				LinesCovered: 30, // 20 (SharedClass) + 10 (UniqueClass1)
				LinesValid:   60, // 40 (SharedClass) + 20 (UniqueClass1)
				Classes: []model.Class{
					{
						Name:         "SharedClass",
						LinesCovered: 20,
						LinesValid:   40,
						Files: []model.CodeFile{
							{Path: "/app/shared.go", TotalLines: 100},
							{Path: "/app/unique1.go", TotalLines: 50},
						},
					},
					{
						Name:         "UniqueClass1",
						LinesCovered: 10,
						LinesValid:   20,
					},
				},
			},
		},
	}

	// Report 2: Also contains the shared assembly with new data for the shared class and another unique class.
	result2 := &parsers.ParserResult{
		ParserName: "Test",
		Assemblies: []model.Assembly{
			{
				Name:         "SharedAssembly",
				LinesCovered: 45, // 30 (SharedClass) + 15 (UniqueClass2)
				LinesValid:   90, // 60 (SharedClass) + 30 (UniqueClass2)
				Classes: []model.Class{
					{
						Name:         "SharedClass",
						LinesCovered: 30,
						LinesValid:   60,
						Files: []model.CodeFile{
							{Path: "/app/shared.go", TotalLines: 100}, // Duplicate file path
							{Path: "/app/unique2.go", TotalLines: 75},
						},
					},
					{
						Name:         "UniqueClass2",
						LinesCovered: 15,
						LinesValid:   30,
					},
				},
			},
		},
	}

	config := &mockMergerConfig{logger: slog.Default()}

	// Act
	summary, err := analyzer.MergeParserResults([]*parsers.ParserResult{result1, result2}, config)

	// Assert
	require.NoError(t, err, "Merging should not produce an error")
	require.Len(t, summary.Assemblies, 1, "There should be only one final assembly after merging")

	// --- Verify Assembly Level ---
	mergedAsm := summary.Assemblies[0]
	assert.Equal(t, "SharedAssembly", mergedAsm.Name)
	// Verify top-level assembly stats are correctly summed from the original assemblies
	assert.Equal(t, 75, mergedAsm.LinesCovered, "Assembly lines covered should be summed (30 + 45)")
	assert.Equal(t, 150, mergedAsm.LinesValid, "Assembly lines valid should be summed (60 + 90)")
	require.Len(t, mergedAsm.Classes, 3, "Final assembly should contain the merged shared class and both unique classes")

	// --- Find and Verify Shared Class ---
	var sharedClass *model.Class
	for i := range mergedAsm.Classes {
		if mergedAsm.Classes[i].Name == "SharedClass" {
			sharedClass = &mergedAsm.Classes[i]
			break
		}
	}
	require.NotNil(t, sharedClass, "The merged 'SharedClass' should exist")
	assert.Equal(t, 50, sharedClass.LinesCovered, "SharedClass lines covered should be summed (20 + 30)")
	assert.Equal(t, 100, sharedClass.LinesValid, "SharedClass lines valid should be summed (40 + 60)")

	// --- Verify File De-duplication in Shared Class ---
	require.Len(t, sharedClass.Files, 3, "Files in SharedClass should be the union of both reports (shared.go, unique1.go, unique2.go)")

	// Check file paths to ensure correct de-duplication
	filePathMap := make(map[string]bool)
	for _, f := range sharedClass.Files {
		filePathMap[f.Path] = true
	}
	assert.Contains(t, filePathMap, "/app/shared.go")
	assert.Contains(t, filePathMap, "/app/unique1.go")
	assert.Contains(t, filePathMap, "/app/unique2.go")

	// --- Verify Global Statistics ---
	// TotalLines should be calculated from unique files only: 100 (shared.go) + 50 (unique1.go) + 75 (unique2.go)
	assert.Equal(t, 225, summary.TotalLines, "Grand total lines should be calculated from unique files")
	// Global Lines Covered should equal the final merged assembly's lines covered
	assert.Equal(t, mergedAsm.LinesCovered, summary.LinesCovered)
}

// TestMergeParserResults_WhenClassHasNoFiles_ShouldMergeCorrectly verifies that
// the deep merge logic handles classes without any associated files gracefully.
func TestMergeParserResults_WhenClassHasNoFiles_ShouldMergeCorrectly(t *testing.T) {
	// Arrange
	result1 := &parsers.ParserResult{
		Assemblies: []model.Assembly{
			{
				Name: "AssemblyA",
				Classes: []model.Class{
					{Name: "SharedClass", LinesCovered: 10, LinesValid: 20, Files: []model.CodeFile{}},
				},
			},
		},
	}
	result2 := &parsers.ParserResult{
		Assemblies: []model.Assembly{
			{
				Name: "AssemblyA",
				Classes: []model.Class{
					{Name: "SharedClass", LinesCovered: 15, LinesValid: 30, Files: nil}, // Use nil to test both cases
				},
			},
		},
	}

	config := &mockMergerConfig{logger: slog.Default()}

	// Act
	summary, err := analyzer.MergeParserResults([]*parsers.ParserResult{result1, result2}, config)

	// Assert
	require.NoError(t, err)
	require.Len(t, summary.Assemblies, 1)

	mergedAsm := summary.Assemblies[0]
	require.Len(t, mergedAsm.Classes, 1)

	mergedClass := mergedAsm.Classes[0]
	assert.Equal(t, "SharedClass", mergedClass.Name)
	assert.Equal(t, 25, mergedClass.LinesCovered, "Lines covered should be summed even with no files")
	assert.Equal(t, 50, mergedClass.LinesValid, "Lines valid should be summed even with no files")
	assert.NotNil(t, mergedClass.Files, "Files slice should not be nil even if empty")
	assert.Empty(t, mergedClass.Files, "Files slice should be empty after merging")
}
