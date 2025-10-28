package parser_gocover

import (
	"log/slog"
	"path/filepath"

	"github.com/IgorBayerl/nanovision/filereader"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/parsers"
	"github.com/IgorBayerl/nanovision/internal/utils"
)

// processingOrchestrator now holds state for converting raw blocks into a flat
// list of per-file coverage data. It no longer performs aggregation.
type processingOrchestrator struct {
	fileReader filereader.Reader
	config     parsers.ParserConfig
	logger     *slog.Logger
}

func newProcessingOrchestrator(fileReader filereader.Reader, config parsers.ParserConfig, logger *slog.Logger) *processingOrchestrator {
	return &processingOrchestrator{
		fileReader: fileReader,
		config:     config,
		logger:     logger,
	}
}

// processBlocks is the main entry point for the orchestrator. It groups the raw
// profile blocks by their source file and processes each file independently.
func (o *processingOrchestrator) processBlocks(blocks []GoCoverProfileBlock) ([]parsers.FileCoverage, []string) {
	if len(blocks) == 0 {
		return []parsers.FileCoverage{}, []string{}
	}

	blocksByFile := o.groupBlocksByFile(blocks)
	var allFileCoverage []parsers.FileCoverage
	var allUnresolvedFiles []string

	sourceDir := ""
	if len(o.config.SourceDirectories()) > 0 {
		sourceDir = o.config.SourceDirectories()[0]
	}

	for filePath, fileBlocks := range blocksByFile {
		// Pass the logger from the orchestrator into the find utility
		_, err := utils.FindFileInSourceDirs(filePath, []string{sourceDir}, o.fileReader, o.logger)
		if err != nil {
			o.logger.Warn("Source file not found, it will be marked as unresolved.", "file", filePath, "error", err)
			allUnresolvedFiles = append(allUnresolvedFiles, filePath)
		}

		fileCoverage := o.processFile(filePath, fileBlocks)
		allFileCoverage = append(allFileCoverage, fileCoverage)
	}

	return allFileCoverage, allUnresolvedFiles
}

// groupBlocksByFile creates a map where keys are file paths and values are slices
// of all coverage blocks belonging to that file.
func (o *processingOrchestrator) groupBlocksByFile(blocks []GoCoverProfileBlock) map[string][]GoCoverProfileBlock {
	blocksByFile := make(map[string][]GoCoverProfileBlock)

	for _, block := range blocks {
		// The path is used exactly as it is in the report.
		// The tree builder will be responsible for finding its true location and filtering.
		normalizedPath := filepath.ToSlash(block.FileName)
		blocksByFile[normalizedPath] = append(blocksByFile[normalizedPath], block)
	}
	return blocksByFile
}

// processFile converts all coverage blocks for a single file into a single
// FileCoverage struct, which contains a map of line numbers to their metrics.
func (o *processingOrchestrator) processFile(filePath string, blocks []GoCoverProfileBlock) parsers.FileCoverage {
	lineMetrics := make(map[int]model.LineMetrics)

	for _, block := range blocks {
		// A single block can span multiple lines. We apply its hit count to
		// every line within its range.
		for l := block.StartLine; l <= block.EndLine; l++ {
			// If multiple blocks cover the same line, the Go tool's convention is
			// that the hit count of one of them applies. We take the highest hit count
			// as the most representative value for that line's execution status.
			if existing, ok := lineMetrics[l]; !ok || block.HitCount > existing.Hits {
				lineMetrics[l] = model.LineMetrics{
					Hits: block.HitCount,
					// Go coverage profiles do not support branch coverage.
					TotalBranches:   0,
					CoveredBranches: 0,
				}
			}
		}
	}

	return parsers.FileCoverage{
		Path:  filePath,
		Lines: lineMetrics,
	}
}
