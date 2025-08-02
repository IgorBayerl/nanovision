package htmlreport

import (
	"log/slog"
	"path"
	"strings"

	"github.com/IgorBayerl/AdlerCov/internal/filereader"
	"github.com/IgorBayerl/AdlerCov/internal/model"
	"github.com/IgorBayerl/AdlerCov/internal/utils"
)

// ToLegacySummaryResult converts the new file-system-tree-based model into the
// old Assembly/Class-based model for the legacy HTML reporter.
func ToLegacySummaryResult(tree *model.SummaryTree, fileReader filereader.Reader, logger *slog.Logger) *SummaryResult {
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
			legacyFile := buildLegacyCodeFile(fileNode, fileReader, logger)

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

		// Aggregate assembly metrics from the newly created legacy classes
		for _, cls := range legacyAssembly.Classes {
			legacyAssembly.LinesCovered += cls.LinesCovered
			legacyAssembly.LinesValid += cls.LinesValid
			legacyAssembly.TotalLines += cls.TotalLines
			if cls.BranchesCovered != nil {
				if legacyAssembly.BranchesCovered == nil {
					legacyAssembly.BranchesCovered = new(int)
				}
				*legacyAssembly.BranchesCovered += *cls.BranchesCovered
			}
			if cls.BranchesValid != nil {
				if legacyAssembly.BranchesValid == nil {
					legacyAssembly.BranchesValid = new(int)
				}
				*legacyAssembly.BranchesValid += *cls.BranchesValid
			}
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
	legacyResult.TotalLines = calculateTotalLines(fileNodes, fileReader)

	return legacyResult
}

// buildLegacyCodeFile reads a source file and maps coverage data to each line.
// This is the core fix: it iterates over source lines, not coverage lines.
func buildLegacyCodeFile(node *model.FileNode, fileReader filereader.Reader, logger *slog.Logger) CodeFile {
	legacyFile := CodeFile{
		Path:           node.Path,
		CoveredLines:   node.Metrics.LinesCovered,
		CoverableLines: node.Metrics.LinesValid,
	}

	var sourceLines []string
	resolvedPath, err := utils.FindFileInSourceDirs(node.Path, []string{node.SourceDir}, fileReader, logger)
	if err != nil {
		logger.Warn("Could not resolve source file for HTML adapter", "file", node.Path, "error", err)
	} else {
		lines, readErr := fileReader.ReadFile(resolvedPath)
		if readErr != nil {
			logger.Warn("Could not read source file for HTML adapter", "file", resolvedPath, "error", readErr)
		} else {
			sourceLines = lines
		}
	}
	legacyFile.TotalLines = len(sourceLines)

	var legacyLines []Line
	for i, lineContent := range sourceLines {
		lineNumber := i + 1
		lineMetrics, hasCoverageData := node.Lines[lineNumber]

		status := NotCoverable
		hits := 0
		isBranch := false
		coveredBranches := 0
		totalBranches := 0

		// Only apply coverage status if the line exists in the report and is coverable.
		if hasCoverageData && lineMetrics.Hits >= 0 {
			hits = lineMetrics.Hits
			isBranch = lineMetrics.TotalBranches > 0
			coveredBranches = lineMetrics.CoveredBranches
			totalBranches = lineMetrics.TotalBranches

			if totalBranches > 0 {
				if coveredBranches == totalBranches {
					status = Covered
				} else if coveredBranches > 0 {
					status = PartiallyCovered
				} else {
					status = NotCovered
				}
			} else if hits > 0 {
				status = Covered
			} else {
				status = NotCovered
			}
		}

		legacyLines = append(legacyLines, Line{
			Number:          lineNumber,
			Hits:            hits,
			Content:         lineContent,
			IsBranchPoint:   isBranch,
			CoveredBranches: coveredBranches,
			TotalBranches:   totalBranches,
			LineVisitStatus: status,
		})
	}

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

// calculateTotalLines must now read each unique file to get its true line count.
func calculateTotalLines(fileNodes []*model.FileNode, fileReader filereader.Reader) int {
	total := 0
	uniqueFiles := make(map[string]bool)
	for _, file := range fileNodes {
		if _, exists := uniqueFiles[file.Path]; !exists {
			resolvedPath, err := utils.FindFileInSourceDirs(file.Path, []string{file.SourceDir}, fileReader, nil)
			if err == nil {
				count, err := fileReader.CountLines(resolvedPath)
				if err == nil {
					total += count
				}
			}
			uniqueFiles[file.Path] = true
		}
	}
	return total
}
