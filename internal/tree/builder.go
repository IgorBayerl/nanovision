package tree

import (
	"fmt"
	"path"
	"strings"

	"github.com/IgorBayerl/AdlerCov/internal/model"
	"github.com/IgorBayerl/AdlerCov/internal/parsers"
)

// Package tree is responsible for constructing the canonical data model for a coverage report.
// Its primary component, the Builder, acts as the central point for structuring and merging
// raw coverage data from various sources into a unified, hierarchical representation.
//
// # Architectural Role
//
// This package sits squarely between the 'Parsing' and 'Hydration' stages of the application pipeline.
//
//  1. Input: It receives a flat list of `parsers.ParserResult` objects, where each object
//     contains the raw line and branch coverage for the files found in a single report.
//
//  2. Responsibility: Its sole purpose is to build a `model.SummaryTree`. This tree mirrors the
//     project's filesystem structure and serves as the single source of truth for all subsequent
//     processing. A key function of the builder is to correctly merge coverage data when multiple
//     reports (e.g., from unit and integration tests) cover the same source file. It does this
//     by summing the hit counts and branch coverage for each line.
//
//  3. Output: It produces a single, raw `model.SummaryTree`. This tree is "raw" because it contains
//     only the structural information and the merged coverage metrics. It does not yet contain
//     rich details derived from static code analysis, such as method boundaries or cyclomatic
//     complexity. That subsequent enrichment is the responsibility of the 'Hydrator' package.
//
// By centralizing the tree construction and merging logic here, we decouple the parsers from the
// complexities of data aggregation and ensure that the Hydrator and all subsequent Reporter
// components work with a single, consistent, and well-structured data model.
type Builder struct {
	// Future dependencies like a logger could be added here if the building
	// process requires more detailed logging.
}

// NewBuilder creates a new instance of a tree Builder.
func NewBuilder() *Builder {
	return &Builder{}
}

// BuildTree orchestrates the entire process of merging multiple parser results
// into a single, aggregated SummaryTree. It ensures that coverage data from
// different reports for the same file is correctly merged.
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

	return tree, nil
}

// findOrCreateFileNode traverses the tree from the given startNode according to the
// components of filePath. It creates any missing directory nodes along the way.
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

// mergeLineMetrics additively combines new line coverage data into an existing file node.
// This is the core of the merging logic, ensuring that hits and branch counts are summed
// when multiple reports cover the same line.
func (b *Builder) mergeLineMetrics(node *model.FileNode, newLines map[int]model.LineMetrics) {
	for lineNum, newLineMetric := range newLines {
		if existingLineMetric, ok := node.Lines[lineNum]; ok {
			existingLineMetric.Hits += newLineMetric.Hits
			existingLineMetric.CoveredBranches += newLineMetric.CoveredBranches
			existingLineMetric.TotalBranches = newLineMetric.TotalBranches
			node.Lines[lineNum] = existingLineMetric
		} else {
			node.Lines[lineNum] = newLineMetric
		}
	}
}
