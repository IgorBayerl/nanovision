// Path: internal/analyzer/tree_builder.go
package analyzer

import (
	"fmt"
	"path"
	"strings"

	"github.com/IgorBayerl/AdlerCov/internal/model"
	"github.com/IgorBayerl/AdlerCov/internal/parsers"
)

// TreeBuilder is responsible for constructing the canonical SummaryTree data model
// from the raw, per-file coverage data provided by one or more parsers.
// Its primary role is to create the hierarchical structure of directories and files
// and to perform the bottom-up aggregation of coverage metrics.
type TreeBuilder struct {
	// In the future, dependencies like a logger or configuration settings
	// could be added here to influence the tree-building process.
}

// NewTreeBuilder creates a new instance of a TreeBuilder.
func NewTreeBuilder() *TreeBuilder {
	return &TreeBuilder{}
}

// BuildTree orchestrates the entire process of merging multiple parser results
// into a single, aggregated SummaryTree. It ensures that coverage data from
// different reports (e.g., unit and integration tests) for the same file
// is correctly merged and that all directory-level metrics are accurately calculated.
func (b *TreeBuilder) BuildTree(results []*parsers.ParserResult) (*model.SummaryTree, error) {
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
		ParserName: results[0].ParserName, // Use the first parser's name as a representative.
	}

	for _, result := range results {
		for _, fileCov := range result.FileCoverage {
			// Find or create the node for this file in the tree. This builds out the
			// directory structure as a side effect.
			fileNode := b.findOrCreateFileNode(tree.Root, fileCov.Path)

			// Merge the line-level metrics from the current report into the file node.
			// This additively combines hits and branch data if a file is covered
			// by multiple reports.
			b.mergeLineMetrics(fileNode, fileCov.Lines)
		}
	}

	// After all files from all reports have been processed and their data merged,
	// perform a final recursive aggregation. This calculates all file-level,
	// directory-level, and root-level metrics in a single pass.
	b.aggregateTreeMetrics(tree.Root)

	// The root node's metrics now represent the entire project's coverage.
	// We copy them to the top level of the SummaryTree for easy access.
	tree.Metrics = tree.Root.Metrics

	return tree, nil
}

// findOrCreateFileNode traverses the tree from the given startNode according to the
// components of filePath. It creates any missing directory nodes along the way.
// This design ensures that the file system hierarchy is accurately represented
// in the tree structure without needing to know the structure in advance.
func (b *TreeBuilder) findOrCreateFileNode(startNode *model.DirNode, filePath string) *model.FileNode {
	// We assume paths use forward slashes as a normalized separator.
	// The parsers are responsible for this normalization.
	parts := strings.Split(filePath, "/")
	currentNode := startNode

	// Iterate through the directory parts of the path.
	for _, part := range parts[:len(parts)-1] {
		// If a subdirectory for the current path part does not exist, create it.
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
	// If a file node for the final path component does not exist, create it.
	if _, ok := currentNode.Files[fileName]; !ok {
		newFile := &model.FileNode{
			Name:   fileName,
			Path:   filePath,
			Lines:  make(map[int]model.LineMetrics),
			Parent: currentNode,
		}
		currentNode.Files[fileName] = newFile
	}

	return currentNode.Files[fileName]
}

// mergeLineMetrics additively combines new line coverage data into an existing file node.
// This is crucial for correctly calculating total coverage when multiple reports
// cover the same file. Hits and branch counts are summed.
func (b *TreeBuilder) mergeLineMetrics(node *model.FileNode, newLines map[int]model.LineMetrics) {
	for lineNum, newLineMetric := range newLines {
		if existingLineMetric, ok := node.Lines[lineNum]; ok {
			// The line already exists, so we sum the metrics.
			existingLineMetric.Hits += newLineMetric.Hits
			existingLineMetric.CoveredBranches += newLineMetric.CoveredBranches
			// TotalBranches should be the same across reports for the same line,
			// but we take the new value just in case of inconsistencies.
			existingLineMetric.TotalBranches = newLineMetric.TotalBranches
			node.Lines[lineNum] = existingLineMetric
		} else {
			// This is the first time we've seen coverage for this line.
			node.Lines[lineNum] = newLineMetric
		}
	}
}

// aggregateTreeMetrics performs a recursive, post-order traversal of the tree
// to calculate the aggregated CoverageMetrics for every node.
// It starts by calculating metrics for leaf files, then aggregates those into
// their parent directories, and so on, up to the root.
func (b *TreeBuilder) aggregateTreeMetrics(node *model.DirNode) {
	// First, recurse down to the deepest subdirectories.
	for _, subdir := range node.Subdirs {
		b.aggregateTreeMetrics(subdir)
	}

	// Reset the current node's metrics before recalculating.
	node.Metrics = model.CoverageMetrics{}

	// Aggregate metrics from all files within this directory.
	for _, file := range node.Files {
		// Ensure the file's own metrics are calculated from its lines first.
		b.calculateFileMetrics(file)
		node.Metrics.LinesCovered += file.Metrics.LinesCovered
		node.Metrics.LinesValid += file.Metrics.LinesValid
		node.Metrics.BranchesCovered += file.Metrics.BranchesCovered
		node.Metrics.BranchesValid += file.Metrics.BranchesValid
	}

	// Aggregate metrics from all subdirectories (which have already been aggregated).
	for _, subdir := range node.Subdirs {
		node.Metrics.LinesCovered += subdir.Metrics.LinesCovered
		node.Metrics.LinesValid += subdir.Metrics.LinesValid
		node.Metrics.BranchesCovered += subdir.Metrics.BranchesCovered
		node.Metrics.BranchesValid += subdir.Metrics.BranchesValid
	}
}

// calculateFileMetrics computes the total coverage metrics for a single file
// by iterating over all its line-level metrics. This is the lowest level of aggregation.
func (b *TreeBuilder) calculateFileMetrics(node *model.FileNode) {
	metrics := model.CoverageMetrics{}
	for _, line := range node.Lines {
		// A line is considered "coverable" or "valid" if it has a non-negative hit count.
		// A hit count of -1 is a sentinel for non-executable code (e.g., comments, braces).
		if line.Hits >= 0 {
			metrics.LinesValid++
			if line.Hits > 0 {
				metrics.LinesCovered++
			}
		}
		metrics.BranchesCovered += line.CoveredBranches
		metrics.BranchesValid += line.TotalBranches
	}
	node.Metrics = metrics
}
