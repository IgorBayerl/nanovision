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
// allowing the refactored core application to support the legacy HTML reporter.
//
// STRATEGY: To enable the frontend's "Group by" feature and provide a flat-file view,
// this adapter treats EACH source file as its OWN "Class". The frontend then receives
// a granular list that it can group by namespace (directory path).
// The 'allSourceDirs' parameter is now ignored, but kept for signature stability if needed elsewhere.
// The real source of truth is the fileNode.SourceDir.
func ToLegacySummaryResult(tree *model.SummaryTree, fileReader filereader.Reader, allSourceDirs []string, logger *slog.Logger) *SummaryResult {
	if tree == nil {
		return &SummaryResult{}
	}

	legacyResult := &SummaryResult{
		ParserName: tree.ParserName,
		Timestamp:  tree.Timestamp,
	}

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

	var filePaths []string
	for _, fileNode := range fileNodes {
		filePaths = append(filePaths, fileNode.Path)
	}
	displayBasePath := findDisplayBasePath(filePaths)

	filesByAssembly := make(map[string][]*model.FileNode)
	for _, fileNode := range fileNodes {
		parts := strings.Split(fileNode.Path, "/")
		assemblyName := "Default"
		if len(parts) > 1 {
			assemblyName = parts[0]
		}
		filesByAssembly[assemblyName] = append(filesByAssembly[assemblyName], fileNode)
	}

	for assemblyName, filesInAssembly := range filesByAssembly {
		legacyAssembly := Assembly{Name: assemblyName}

		for _, fileNode := range filesInAssembly {
			// THIS IS THE KEY CHANGE:
			// We pass a list containing ONLY the file's true source directory, which is stored on the node itself.
			// This completely removes the ambiguity of searching in multiple directories.
			legacyFile := buildLegacyCodeFile(fileNode, fileReader, []string{fileNode.SourceDir}, logger)

			legacyClass := Class{
				Name:         fileNode.Path,
				DisplayName:  strings.TrimPrefix(strings.TrimPrefix(fileNode.Path, displayBasePath), "/"),
				Files:        []CodeFile{legacyFile},
				LinesCovered: legacyFile.CoveredLines,
				LinesValid:   legacyFile.CoverableLines,
				TotalLines:   legacyFile.TotalLines,
				TotalMethods: len(fileNode.Methods),
			}
			if fileNode.Metrics.BranchesValid > 0 {
				bc := fileNode.Metrics.BranchesCovered
				bv := fileNode.Metrics.BranchesValid
				legacyClass.BranchesCovered = &bc
				legacyClass.BranchesValid = &bv
			}

			legacyAssembly.Classes = append(legacyAssembly.Classes, legacyClass)
		}

		for _, cls := range legacyAssembly.Classes {
			legacyAssembly.LinesCovered += cls.LinesCovered
			legacyAssembly.LinesValid += cls.LinesValid
			legacyAssembly.TotalLines += cls.TotalLines
		}

		legacyResult.Assemblies = append(legacyResult.Assemblies, legacyAssembly)
	}

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

func buildLegacyCodeFile(node *model.FileNode, fileReader filereader.Reader, sourceDirs []string, logger *slog.Logger) CodeFile {
	legacyFile := CodeFile{
		Path:           node.Path,
		CoveredLines:   node.Metrics.LinesCovered,
		CoverableLines: node.Metrics.LinesValid,
	}

	var sourceLines []string
	// The sourceDirs slice is now guaranteed to contain exactly one, correct directory for this file.
	resolvedPath, err := utils.FindFileInSourceDirs(node.Path, sourceDirs, fileReader, logger)
	if err != nil {
		logger.Warn("Could not resolve source file for HTML adapter", "file", node.Path, "error", err)
	} else {
		logger.Info("Successfully resolved source file for HTML report", "file", node.Path, "resolved_path", resolvedPath)
		sourceLines, err = fileReader.ReadFile(resolvedPath)
		if err != nil {
			logger.Warn("Could not read source file for HTML adapter", "file", resolvedPath, "error", err)
			sourceLines = []string{}
		}
	}
	legacyFile.TotalLines = len(sourceLines)

	var legacyLines []Line
	for lineNum, lineMetrics := range node.Lines {
		var content string
		if lineNum > 0 && lineNum <= len(sourceLines) {
			content = sourceLines[lineNum-1]
		}

		status := NotCoverable
		if lineMetrics.Hits >= 0 {
			if lineMetrics.TotalBranches > 0 {
				if lineMetrics.CoveredBranches == lineMetrics.TotalBranches {
					status = Covered
				} else if lineMetrics.CoveredBranches > 0 {
					status = PartiallyCovered
				} else {
					status = NotCovered
				}
			} else if lineMetrics.Hits > 0 {
				status = Covered
			} else {
				status = NotCovered
			}
		}

		legacyLines = append(legacyLines, Line{
			Number:          lineNum,
			Hits:            lineMetrics.Hits,
			Content:         content,
			IsBranchPoint:   lineMetrics.TotalBranches > 0,
			CoveredBranches: lineMetrics.CoveredBranches,
			TotalBranches:   lineMetrics.TotalBranches,
			LineVisitStatus: status,
		})
	}

	sort.Slice(legacyLines, func(i, j int) bool {
		return legacyLines[i].Number < legacyLines[j].Number
	})

	legacyFile.Lines = legacyLines
	return legacyFile
}

// findDisplayBasePath finds the common parent directory of a list of file paths.
func findDisplayBasePath(paths []string) string {
	if len(paths) < 2 {
		if len(paths) == 1 {
			dir := path.Dir(paths[0])
			if dir == "." {
				return ""
			}
			return dir
		}
		return ""
	}

	pathComponents := make([][]string, len(paths))
	for i, p := range paths {
		pathComponents[i] = strings.Split(path.Dir(p), "/")
	}

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

// calculateTotalLines sums the actual line counts from all unique source files.
func calculateTotalLines(fileNodes []*model.FileNode) int {
	// This function remains an estimation based on available data.
	// A more precise count would require passing the fully constructed legacy files,
	// but this is sufficient for the summary card.
	total := 0
	uniqueFiles := make(map[string]bool)
	for _, file := range fileNodes {
		if _, exists := uniqueFiles[file.Path]; !exists {
			// Using file.Metrics.LinesValid as a proxy for lines of code,
			// though the actual file line count might differ.
			total += file.Metrics.LinesValid
			uniqueFiles[file.Path] = true
		}
	}
	return total
}
