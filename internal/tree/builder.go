package tree

import (
	"fmt"
	"log/slog"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/IgorBayerl/nanovision/filereader"
	"github.com/IgorBayerl/nanovision/filtering"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/parsers"
	"github.com/IgorBayerl/nanovision/internal/utils"
)

type Builder struct {
	projectRoot string
	fileReader  filereader.Reader
	fileFilter  filtering.IFilter
}

func NewBuilder(projectRoot string, fileFilter filtering.IFilter) *Builder {
	return &Builder{
		projectRoot: projectRoot,
		fileReader:  filereader.NewDefaultReader(),
		fileFilter:  fileFilter,
	}
}

func (b *Builder) BuildTree(results []*parsers.ParserResult) (*model.SummaryTree, error) {
	if len(results) == 0 {
		return nil, fmt.Errorf("cannot build tree with no parser results")
	}

	logger := slog.Default()

	tree := &model.SummaryTree{
		Root: &model.DirNode{
			Name:    "Root",
			Path:    ".",
			Subdirs: make(map[string]*model.DirNode),
			Files:   make(map[string]*model.FileNode),
		},
	}

	// Create a stable, ordered list of report sources to use as indices
	// Using the report pattern as a unique key for the report group
	// This prevents reports that cover the same part of the project to merge
	reportNameMap := make(map[string]int)
	for _, result := range results {
		reportKey := result.ReportPattern
		if _, exists := reportNameMap[reportKey]; !exists {
			reportNameMap[reportKey] = len(tree.ReportNames)
			tree.ReportNames = append(tree.ReportNames, reportKey)
		}
	}
	numReports := len(tree.ReportNames)

	uniqueParsers := make(map[string]struct{})
	for _, result := range results {
		if result.ParserName != "" {
			uniqueParsers[result.ParserName] = struct{}{}
		}
	}

	var parserNames []string
	for name := range uniqueParsers {
		parserNames = append(parserNames, name)
	}
	sort.Strings(parserNames) // Sort for consistent output
	tree.ParserNames = parserNames

	for _, result := range results {
		for _, fileCov := range result.FileCoverage {
			// Find the canonical absolute path of the file.
			absoluteFilePath, err := utils.FindFileInSourceDirs(fileCov.Path, []string{result.SourceDirectory}, b.fileReader, logger)
			if err != nil {
				logger.Warn("Could not resolve file path, skipping file.", "path", fileCov.Path, "sourceDir", result.SourceDirectory, "error", err)
				continue
			}

			// Make it relative to our project root.
			relativeToProjectRoot, err := filepath.Rel(b.projectRoot, absoluteFilePath)
			if err != nil {
				logger.Warn("Could not make path relative to project root, using absolute.", "path", absoluteFilePath, "projectRoot", b.projectRoot, "error", err)
				relativeToProjectRoot = absoluteFilePath // Fallback
			}

			// Ensure consistent separators and use this path to build the tree.
			finalPath := filepath.ToSlash(relativeToProjectRoot)

			// Apply filtering on project relative path.
			if !b.fileFilter.IsElementIncludedInReport(finalPath) {
				logger.Debug("File excluded by filter", "path", finalPath)
				continue
			}

			reportKey := result.ReportPattern
			reportIndex := reportNameMap[reportKey]
			fileNode := b.findOrCreateFileNode(tree.Root, finalPath, result.SourceDirectory)
			b.mergeLineMetrics(fileNode, fileCov.Lines, reportIndex, numReports)
		}
	}

	// After all files are added and merged, perform a final aggregation pass.
	tree.Metrics = b.aggregateMetrics(tree.Root)

	return tree, nil
}

func (b *Builder) findOrCreateFileNode(startNode *model.DirNode, filePath string, sourceDir string) *model.FileNode {
	parts := strings.Split(filePath, "/")
	currentNode := startNode

	for _, part := range parts[:len(parts)-1] {
		if _, ok := currentNode.Subdirs[part]; !ok {
			newDir := &model.DirNode{
				Name:    part,
				Path:    path.Join(currentNode.Path, part),
				Subdirs: make(map[string]*model.DirNode),
				Files:   make(map[string]*model.FileNode),
				Parent:  currentNode,
			}
			currentNode.Subdirs[part] = newDir
		}
		currentNode = currentNode.Subdirs[part]
	}

	fileName := parts[len(parts)-1]
	if _, ok := currentNode.Files[fileName]; !ok {
		newFile := &model.FileNode{
			Name:      fileName,
			Path:      filePath,
			Lines:     make(map[int]model.LineMetrics),
			Parent:    currentNode,
			SourceDir: sourceDir,
		}
		currentNode.Files[fileName] = newFile
	}

	return currentNode.Files[fileName]
}

func (b *Builder) mergeLineMetrics(node *model.FileNode, newLines map[int]model.LineMetrics, reportIndex int, numReports int) {
	for lineNum, newLineMetric := range newLines {
		existing := node.Lines[lineNum]

		if existing.ReportHits == nil {
			existing.ReportHits = make([]int, numReports)
		}

		existing.Hits += newLineMetric.Hits

		existing.ReportHits[reportIndex] = newLineMetric.Hits

		existing.CoveredBranches += newLineMetric.CoveredBranches
		if newLineMetric.TotalBranches > 0 {
			existing.TotalBranches = newLineMetric.TotalBranches
		}
		node.Lines[lineNum] = existing
	}
}

// aggregateMetrics performs a bottom-up aggregation of metrics throughout the tree.
func (b *Builder) aggregateMetrics(dir *model.DirNode) model.CoverageMetrics {
	var dirMetrics model.CoverageMetrics

	// Aggregate from subdirectories first
	for _, subDir := range dir.Subdirs {
		subDirMetrics := b.aggregateMetrics(subDir)
		subDir.Metrics = subDirMetrics
		dirMetrics.LinesCovered += subDirMetrics.LinesCovered
		dirMetrics.LinesValid += subDirMetrics.LinesValid
		dirMetrics.BranchesCovered += subDirMetrics.BranchesCovered
		dirMetrics.BranchesValid += subDirMetrics.BranchesValid
		dirMetrics.TotalLines += subDirMetrics.TotalLines
	}

	// Aggregate from files in the current directory
	for _, file := range dir.Files {
		fileMetrics := b.calculateFileMetrics(file)
		file.Metrics = fileMetrics
		dirMetrics.LinesCovered += fileMetrics.LinesCovered
		dirMetrics.LinesValid += fileMetrics.LinesValid
		dirMetrics.BranchesCovered += fileMetrics.BranchesCovered
		dirMetrics.BranchesValid += fileMetrics.BranchesValid
		dirMetrics.TotalLines += fileMetrics.TotalLines
	}

	return dirMetrics
}

func (b *Builder) calculateFileMetrics(file *model.FileNode) model.CoverageMetrics {
	metrics := model.CoverageMetrics{TotalLines: file.TotalLines}
	for _, line := range file.Lines {
		if line.Hits >= 0 { // Is a coverable line
			metrics.LinesValid++
			if line.Hits > 0 {
				metrics.LinesCovered++
			}
		}
		metrics.BranchesValid += line.TotalBranches
		metrics.BranchesCovered += line.CoveredBranches
	}
	return metrics
}
