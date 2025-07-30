// Path: internal/reporter/htmlreport/adapter.go
package htmlreport

import (
	"log/slog"
	"path"
	"sort"
	"strings"

	"github.com/IgorBayerl/AdlerCov/internal/filereader"
	"github.com/IgorBayerl/AdlerCov/internal/model"
	"github.com/IgorBayerl/AdlerCov/internal/utils"
)

// ToLegacySummaryResult converts the new file-system-tree-based model into the
// old Assembly/Class-based model. This function acts as an anti-corruption layer,
// allowing the refactored core application to support the legacy HTML reporter
// without being coupled to its outdated data structures.
func ToLegacySummaryResult(tree *model.SummaryTree, fileReader filereader.Reader, sourceDirs []string, logger *slog.Logger) *SummaryResult {
	if tree == nil {
		return &SummaryResult{}
	}

	legacyResult := &SummaryResult{
		ParserName: tree.ParserName,
		Timestamp:  tree.Timestamp,
	}

	// Step 1: Collect all file nodes from the tree to analyze their paths.
	var fileNodes []*model.FileNode
	var collectFiles func(*model.DirNode)
	collectFiles = func(dir *model.DirNode) {
		for _, file := range dir.Files {
			fileNodes = append(fileNodes, file)
		}
		for _, subDir := range dir.Subdirs {
			collectFiles(subDir)
		}
	}
	collectFiles(tree.Root)

	// Step 2: Find the common parent directory of all source files.
	// This becomes the base path we trim from the full paths to create clean display names.
	var filePaths []string
	for _, fileNode := range fileNodes {
		filePaths = append(filePaths, fileNode.Path)
	}
	displayBasePath := findDisplayBasePath(filePaths)

	// Step 3: Group files by a synthetic "Assembly" name.
	filesByAssembly := make(map[string][]*model.FileNode)
	for _, fileNode := range fileNodes {
		parts := strings.Split(fileNode.Path, "/")
		assemblyName := "Default"
		if len(parts) > 1 {
			assemblyName = parts[0]
		}
		filesByAssembly[assemblyName] = append(filesByAssembly[assemblyName], fileNode)
	}

	// Step 4: For each assembly, group files by their directory to form "Classes".
	for assemblyName, filesInAssembly := range filesByAssembly {
		legacyAssembly := Assembly{Name: assemblyName}
		filesByClass := make(map[string][]*model.FileNode)

		for _, fileNode := range filesInAssembly {
			className := path.Dir(fileNode.Path)
			if className == "." {
				className = "(root)"
			}
			filesByClass[className] = append(filesByClass[className], fileNode)
		}

		for className, filesInClass := range filesByClass {
			legacyClass := Class{Name: className}

			displayName := strings.TrimPrefix(className, displayBasePath)
			displayName = strings.TrimPrefix(displayName, "/")
			if displayName == "" {
				if base := path.Base(className); base != "." && base != "/" {
					displayName = base
				} else {
					displayName = "(root)"
				}
			}
			legacyClass.DisplayName = displayName

			for _, fileNode := range filesInClass {
				// *** FIX: Pass dependencies to the builder function ***
				legacyFile := buildLegacyCodeFile(fileNode, fileReader, sourceDirs, logger)

				legacyClass.Files = append(legacyClass.Files, legacyFile)
				legacyClass.LinesCovered += legacyFile.CoveredLines
				legacyClass.LinesValid += legacyFile.CoverableLines
				legacyClass.TotalLines += legacyFile.TotalLines
			}
			legacyAssembly.Classes = append(legacyAssembly.Classes, legacyClass)
			legacyAssembly.LinesCovered += legacyClass.LinesCovered
			legacyAssembly.LinesValid += legacyClass.LinesValid
			legacyAssembly.TotalLines += legacyClass.TotalLines
		}
		legacyResult.Assemblies = append(legacyResult.Assemblies, legacyAssembly)
	}

	// Final Step: Copy root metrics to the legacy result.
	legacyResult.LinesCovered = tree.Metrics.LinesCovered
	legacyResult.LinesValid = tree.Metrics.LinesValid
	if tree.Metrics.BranchesValid > 0 || tree.Metrics.BranchesCovered > 0 {
		bc := tree.Metrics.BranchesCovered
		bv := tree.Metrics.BranchesValid
		legacyResult.BranchesCovered = &bc
		legacyResult.BranchesValid = &bv
	}
	legacyResult.TotalLines = calculateTotalLines(fileNodes)

	return legacyResult
}

// findDisplayBasePath finds the common parent directory of a list of file paths.
// This allows for creating shorter, relative display names in the report.
func findDisplayBasePath(paths []string) string {
	if len(paths) < 2 {
		// If there's only one file, its parent directory is the base.
		if len(paths) == 1 {
			dir := path.Dir(paths[0])
			if dir == "." {
				return ""
			}
			return dir
		}
		return ""
	}

	// Split all paths into their directory components.
	pathComponents := make([][]string, len(paths))
	for i, p := range paths {
		pathComponents[i] = strings.Split(path.Dir(p), "/")
	}

	// Find the point where the paths diverge.
	shortestPathLen := len(pathComponents[0])
	for _, components := range pathComponents[1:] {
		if len(components) < shortestPathLen {
			shortestPathLen = len(components)
		}
	}

	var commonPrefix []string
	for i := 0; i < shortestPathLen; i++ {
		firstPathComponent := pathComponents[0][i]
		isCommon := true
		for _, otherComponents := range pathComponents[1:] {
			if otherComponents[i] != firstPathComponent {
				isCommon = false
				break
			}
		}
		if !isCommon {
			break
		}
		commonPrefix = append(commonPrefix, firstPathComponent)
	}

	if len(commonPrefix) == 0 {
		return ""
	}
	return strings.Join(commonPrefix, "/")
}

// buildLegacyCodeFile and calculateTotalLines are unchanged from the previous version.
func buildLegacyCodeFile(node *model.FileNode, fileReader filereader.Reader, sourceDirs []string, logger *slog.Logger) CodeFile {
	legacyFile := CodeFile{
		Path:           node.Path,
		CoveredLines:   node.Metrics.LinesCovered,
		CoverableLines: node.Metrics.LinesValid,
	}

	// First, we must find the absolute path to the source file on disk so we can read its content.
	var sourceLines []string
	resolvedPath, err := utils.FindFileInSourceDirs(node.Path, sourceDirs, fileReader)
	if err != nil {
		// If the file can't be found, we log a warning but continue. The report will
		// still be generated, but this specific file will not have its source code visible.
		logger.Warn("Could not resolve source file for HTML adapter", "file", node.Path, "error", err)
	} else {
		// If the path was resolved, we attempt to read the file's contents line by line.
		sourceLines, err = fileReader.ReadFile(resolvedPath)
		if err != nil {
			logger.Warn("Could not read source file for HTML adapter", "file", resolvedPath, "error", err)
			// Ensure sourceLines is an empty slice on error to prevent panics later.
			sourceLines = []string{}
		}
	}
	legacyFile.TotalLines = len(sourceLines)

	var legacyLines []Line
	for lineNum, lineMetrics := range node.Lines {
		var content string
		// Fetch the line content, making sure to do a bounds check as the report
		// and the source file could theoretically be out of sync.
		if lineNum > 0 && lineNum <= len(sourceLines) {
			content = sourceLines[lineNum-1] // Slices are 0-indexed, line numbers are 1-indexed.
		}

		// Determine the legacy LineVisitStatus required by the report's CSS for color-coding.
		status := NotCoverable
		if lineMetrics.Hits >= 0 { // A non-negative hit count indicates an executable line.
			if lineMetrics.TotalBranches > 0 {
				// For branch points, the status depends on how many branches were taken.
				if lineMetrics.CoveredBranches == lineMetrics.TotalBranches {
					status = Covered
				} else if lineMetrics.CoveredBranches > 0 {
					status = PartiallyCovered
				} else {
					status = NotCovered
				}
			} else if lineMetrics.Hits > 0 {
				// For simple lines, a positive hit count means it's covered.
				status = Covered
			} else {
				// A zero hit count means it's an uncovered executable line.
				status = NotCovered
			}
		}

		legacyLines = append(legacyLines, Line{
			Number:          lineNum,
			Hits:            lineMetrics.Hits,
			Content:         content, // This is the crucial field for code visualization.
			IsBranchPoint:   lineMetrics.TotalBranches > 0,
			CoveredBranches: lineMetrics.CoveredBranches,
			TotalBranches:   lineMetrics.TotalBranches,
			LineVisitStatus: status,
		})
	}

	// The map of lines from the new model is unordered. Sorting is essential here
	// to ensure the lines are rendered in the correct order in the final HTML file.
	sort.Slice(legacyLines, func(i, j int) bool {
		return legacyLines[i].Number < legacyLines[j].Number
	})

	legacyFile.Lines = legacyLines
	return legacyFile
}

func calculateTotalLines(fileNodes []*model.FileNode) int {
	total := 0
	for _, file := range fileNodes {
		total += len(file.Lines)
	}
	return total
}
