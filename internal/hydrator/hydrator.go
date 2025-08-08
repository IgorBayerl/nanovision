package hydrator

import (
	"log/slog"
	"strings"

	"github.com/IgorBayerl/AdlerCov/internal/filereader"
	"github.com/IgorBayerl/AdlerCov/internal/language"
	"github.com/IgorBayerl/AdlerCov/internal/model"
	"github.com/IgorBayerl/AdlerCov/internal/utils"
)

// Hydrator enriches a raw SummaryTree with detailed information
// derived from static analysis of the source code.
type Hydrator struct {
	fileReader  filereader.Reader
	langFactory *language.ProcessorFactory
	logger      *slog.Logger
}

// NewHydrator creates a new Hydrator instance.
func NewHydrator(fileReader filereader.Reader, langFactory *language.ProcessorFactory, logger *slog.Logger) *Hydrator {
	return &Hydrator{
		fileReader:  fileReader,
		langFactory: langFactory,
		logger:      logger,
	}
}

// HydrateTree is the main entry point. It takes a raw tree, populates it with
// analytical data, and calculates the final metrics.
func (h *Hydrator) HydrateTree(tree *model.SummaryTree) error {
	var allFiles []*model.FileNode
	collectFiles(tree.Root, &allFiles)

	for _, fileNode := range allFiles {
		h.hydrateFile(fileNode)
	}

	// After all files are hydrated, perform the final bottom-up metric aggregation.
	h.aggregateTreeMetrics(tree.Root)
	tree.Metrics = tree.Root.Metrics

	return nil
}

// hydrateFile performs the analysis for a single file. It reads the source code,
// performs analysis, populates the node, and then discards the source code from memory.
func (h *Hydrator) hydrateFile(node *model.FileNode) {
	resolvedPath, err := utils.FindFileInSourceDirs(node.Path, []string{node.SourceDir}, h.fileReader, h.logger)
	if err != nil {
		h.logger.Warn("Could not resolve source file for hydration, skipping analysis.", "file", node.Path, "error", err)
		return
	}

	sourceLines, err := h.fileReader.ReadFile(resolvedPath)
	if err != nil {
		h.logger.Warn("Could not read source file for hydration, skipping analysis.", "file", resolvedPath, "error", err)
		return
	}

	processor := h.langFactory.FindProcessorForFile(node.Path)

	methods, err := processor.AnalyzeFile(resolvedPath, sourceLines)
	if err != nil {
		h.logger.Warn(
			"Static analysis failed for file, method metrics will be incomplete.",
			"file", resolvedPath,
			"processor", processor.Name(),
			"error", err,
		)
	}
	node.Methods = methods

	h.normalizeClosingBraces(node, sourceLines)
}

// normalizeClosingBraces re-classifies lines containing only a '}' as non-coverable.
func (h *Hydrator) normalizeClosingBraces(node *model.FileNode, sourceLines []string) {
	for lineNum, lineMetric := range node.Lines {
		if lineNum > 0 && lineNum <= len(sourceLines) {
			lineContent := strings.TrimSpace(sourceLines[lineNum-1])
			if lineContent == "}" {
				lineMetric.Hits = -1 // Re-classify as non-coverable.
				node.Lines[lineNum] = lineMetric
			}
		}
	}
}

// aggregateTreeMetrics performs a recursive, post-order traversal of the tree
// to calculate the aggregated CoverageMetrics for every node.
func (h *Hydrator) aggregateTreeMetrics(node *model.DirNode) {
	for _, subdir := range node.Subdirs {
		h.aggregateTreeMetrics(subdir)
	}

	node.Metrics = model.CoverageMetrics{}

	for _, file := range node.Files {
		calculateFileMetrics(file) // Calculate metrics for the file itself.
		node.Metrics.LinesCovered += file.Metrics.LinesCovered
		node.Metrics.LinesValid += file.Metrics.LinesValid
		node.Metrics.BranchesCovered += file.Metrics.BranchesCovered
		node.Metrics.BranchesValid += file.Metrics.BranchesValid
	}

	for _, subdir := range node.Subdirs {
		node.Metrics.LinesCovered += subdir.Metrics.LinesCovered
		node.Metrics.LinesValid += subdir.Metrics.LinesValid
		node.Metrics.BranchesCovered += subdir.Metrics.BranchesCovered
		node.Metrics.BranchesValid += subdir.Metrics.BranchesValid
	}
}

// calculateFileMetrics computes metrics for a file based on its (now normalized) line data.
func calculateFileMetrics(node *model.FileNode) {
	metrics := model.CoverageMetrics{}
	for _, line := range node.Lines {
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

// collectFiles is a helper to get a flat list of all file nodes from the tree.
func collectFiles(dir *model.DirNode, files *[]*model.FileNode) {
	for _, file := range dir.Files {
		*files = append(*files, file)
	}
	for _, subDir := range dir.Subdirs {
		collectFiles(subDir, files)
	}
}
