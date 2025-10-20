package tree

import (
	"fmt"
	"path"
	"strings"

	"github.com/IgorBayerl/AdlerCov/internal/model"
	"github.com/IgorBayerl/AdlerCov/internal/parsers"
)

type Builder struct{}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) BuildTree(results []*parsers.ParserResult) (*model.SummaryTree, error) {
	if len(results) == 0 {
		return nil, fmt.Errorf("cannot build tree with no parser results")
	}

	tree := &model.SummaryTree{
		Root: &model.DirNode{
			Name:    "Root",
			Path:    ".",
			Subdirs: make(map[string]*model.DirNode),
			Files:   make(map[string]*model.FileNode),
		},
		ParserName: results[0].ParserName,
	}

	for _, result := range results {
		for _, fileCov := range result.FileCoverage {
			fileNode := b.findOrCreateFileNode(tree.Root, fileCov.Path, result.SourceDirectory)
			b.mergeLineMetrics(fileNode, fileCov.Lines)
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

func (b *Builder) mergeLineMetrics(node *model.FileNode, newLines map[int]model.LineMetrics) {
	for lineNum, newLineMetric := range newLines {
		existing := node.Lines[lineNum]
		existing.Hits += newLineMetric.Hits
		existing.CoveredBranches += newLineMetric.CoveredBranches
		// Total branches should not be summed - the last value from a report is authoritative.
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
