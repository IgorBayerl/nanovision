// Package analyzer provides tools for processing and merging parsed coverage data.
// It takes results from one or more parsers and combines them into a single,
// unified summary model, which can then be used by reporters.
package analyzer

import (
	"fmt"
	"log/slog"
	"sort"
	"time"

	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/filtering"
	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/model"
	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/parsers"
)

// MergerConfig defines the necessary configuration for the merging process.
// It provides access to source directories, filters, and a logger.
type MergerConfig interface {
	SourceDirectories() []string
	AssemblyFilters() filtering.IFilter
	Logger() *slog.Logger
}

// MergeParserResults orchestrates the process of merging multiple ParserResult objects
// into a single, unified model.SummaryResult.
func MergeParserResults(results []*parsers.ParserResult, config MergerConfig) (*model.SummaryResult, error) {
	if len(results) == 0 {
		return nil, fmt.Errorf("no parser results to merge")
	}

	logger := config.Logger()
	logger.Info("Starting merge process for parser results", "count", len(results))

	parserName := pickParserName(results)
	logger.Debug("Picked parser name", "name", parserName)

	minTimestamp := earliestTimestamp(results)

	sourceDirs := unionSourceDirs(results)

	mergedAssembliesMap := mergeAssemblies(results, logger)
	logger.Info("Assemblies merged", "count", len(mergedAssembliesMap))

	finalAssemblies := make([]model.Assembly, 0, len(mergedAssembliesMap))
	for _, asm := range mergedAssembliesMap {
		finalAssemblies = append(finalAssemblies, *asm)
	}
	sort.Slice(finalAssemblies, func(i, j int) bool {
		return finalAssemblies[i].Name < finalAssemblies[j].Name
	})

	linesCovered, linesValid, totalLines, branchesCovered, branchesValid, hasBranchData := computeGlobalStats(mergedAssembliesMap)
	logger.Debug("Computed global stats", "linesCovered", linesCovered, "linesValid", linesValid, "hasBranchData", hasBranchData)

	finalSummary := &model.SummaryResult{
		ParserName:   parserName,
		SourceDirs:   sourceDirs,
		Assemblies:   finalAssemblies,
		LinesCovered: linesCovered,
		LinesValid:   linesValid,
		TotalLines:   totalLines,
	}

	if minTimestamp != nil {
		finalSummary.Timestamp = minTimestamp.Unix()
	}

	if hasBranchData {
		finalSummary.BranchesCovered = &branchesCovered
		finalSummary.BranchesValid = &branchesValid
	}

	logger.Info("Merge process completed successfully")
	return finalSummary, nil
}

// --- Helpers ---

// pickParserName inspects the parser results and returns a single representative name.
// It returns "Unknown", the single unique name, or "MultiReport" if multiple parsers were used.
func pickParserName(results []*parsers.ParserResult) string {
	parserNames := make(map[string]struct{})
	for _, res := range results {
		if res.ParserName != "" {
			parserNames[res.ParserName] = struct{}{}
		}
	}
	if len(parserNames) == 1 {
		for name := range parserNames {
			return name
		}
	}
	if len(parserNames) > 1 {
		return "MultiReport"
	}
	return "Unknown"
}

// earliestTimestamp finds and returns the minimum non-nil MinimumTimeStamp from all parser results.
func earliestTimestamp(results []*parsers.ParserResult) *time.Time {
	var minTs *time.Time
	for _, res := range results {
		if res.MinimumTimeStamp != nil {
			if minTs == nil || res.MinimumTimeStamp.Before(*minTs) {
				minTs = res.MinimumTimeStamp
			}
		}
	}
	return minTs
}

// builds and returns a de-duplicated slice of all SourceDirectories from the parser results.
func unionSourceDirs(results []*parsers.ParserResult) []string {
	allSourceDirsSet := make(map[string]struct{})
	for _, res := range results {
		for _, dir := range res.SourceDirectories {
			allSourceDirsSet[dir] = struct{}{}
		}
	}
	sourceDirs := make([]string, 0, len(allSourceDirsSet))
	for dir := range allSourceDirsSet {
		sourceDirs = append(sourceDirs, dir)
	}
	return sourceDirs
}

// combines assemblies from all parser results into a single map using a deep merge strategy.
// If an assembly is found in multiple results, its statistics are summed.
// Its classes are also merged by name, summing their individual statistics and creating a union of their file lists.
func mergeAssemblies(results []*parsers.ParserResult, logger *slog.Logger) map[string]*model.Assembly {
	// Pre-allocate map capacity, guessing an average of 2 assemblies per result.
	mergedAssembliesMap := make(map[string]*model.Assembly, len(results)*2)

	for _, res := range results {
		for _, asmFromParser := range res.Assemblies {
			// Work with a copy to avoid modifying the original parser result data.
			asmCopy := asmFromParser

			if existingAsm, ok := mergedAssembliesMap[asmCopy.Name]; ok {
				logger.Debug("Merging existing assembly", "name", asmCopy.Name)

				// Merge top-level assembly statistics
				existingAsm.LinesCovered += asmCopy.LinesCovered
				existingAsm.LinesValid += asmCopy.LinesValid

				// Merge branch coverage data
				if asmCopy.BranchesCovered != nil {
					if existingAsm.BranchesCovered == nil {
						bc := *asmCopy.BranchesCovered
						existingAsm.BranchesCovered = &bc
					} else {
						*existingAsm.BranchesCovered += *asmCopy.BranchesCovered
					}
				}
				if asmCopy.BranchesValid != nil {
					if existingAsm.BranchesValid == nil {
						bv := *asmCopy.BranchesValid
						existingAsm.BranchesValid = &bv
					} else {
						*existingAsm.BranchesValid += *asmCopy.BranchesValid
					}
				}

				// Deep merge the classes within the assembly
				// Create a map of the existing classes for efficient lookup.
				classMap := make(map[string]*model.Class, len(existingAsm.Classes))
				for i := range existingAsm.Classes {
					classMap[existingAsm.Classes[i].Name] = &existingAsm.Classes[i]
				}

				// Iterate through the new classes from the current parser result
				for _, classFromParser := range asmCopy.Classes {
					if existingClass, found := classMap[classFromParser.Name]; found {
						// Class exists: merge its statistics and files
						existingClass.LinesCovered += classFromParser.LinesCovered
						existingClass.LinesValid += classFromParser.LinesValid

						// Merge the file list to avoid duplicates
						// Create a set of existing file paths for quick lookups.
						filePaths := make(map[string]struct{}, len(existingClass.Files))
						for _, f := range existingClass.Files {
							filePaths[f.Path] = struct{}{}
						}

						// Append only the files that have not been seen before in this class.
						for _, fileFromParser := range classFromParser.Files {
							if _, fileExists := filePaths[fileFromParser.Path]; !fileExists {
								existingClass.Files = append(existingClass.Files, fileFromParser)
								filePaths[fileFromParser.Path] = struct{}{}
							}
						}
					} else {
						// Class is new: append it to the existing assembly's class slice
						existingAsm.Classes = append(existingAsm.Classes, classFromParser)
						classMap[classFromParser.Name] = &existingAsm.Classes[len(existingAsm.Classes)-1]
					}
				}
			} else {
				// Assembly is new, so add a copy of it to the map.
				logger.Debug("Adding new assembly", "name", asmCopy.Name)
				mergedAssembliesMap[asmCopy.Name] = &asmCopy
			}
		}
	}
	return mergedAssembliesMap
}

// computeGlobalStats iterates through the merged assemblies and calculates the final summary statistics in a single pass.
func computeGlobalStats(mergedAssemblies map[string]*model.Assembly) (linesCovered, linesValid, totalLines, branchesCovered, branchesValid int, hasBranchData bool) {
	uniqueFilesForGrandTotal := make(map[string]int)

	for _, asm := range mergedAssemblies {
		linesCovered += asm.LinesCovered
		linesValid += asm.LinesValid

		// Calculate total lines from unique files across all assemblies.
		for _, cls := range asm.Classes {
			for _, f := range cls.Files {
				if _, exists := uniqueFilesForGrandTotal[f.Path]; !exists && f.TotalLines > 0 {
					uniqueFilesForGrandTotal[f.Path] = f.TotalLines
					totalLines += f.TotalLines
				}
			}
		}

		if asm.BranchesCovered != nil && asm.BranchesValid != nil {
			hasBranchData = true
			branchesCovered += *asm.BranchesCovered
			branchesValid += *asm.BranchesValid
		}
	}
	return
}
