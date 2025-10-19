package parser_gocover

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/IgorBayerl/AdlerCov/filereader"
	"github.com/IgorBayerl/AdlerCov/internal/model"
	"github.com/IgorBayerl/AdlerCov/internal/parsers"
	"github.com/IgorBayerl/AdlerCov/internal/utils"
)

var (
	moduleName string
	moduleRoot string
	moduleErr  error
	moduleOnce sync.Once
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

// Helper function to find the go.mod file and parse the module name
func getGoModuleInfo() (string, string, error) {
	moduleOnce.Do(func() {
		currentDir, err := os.Getwd()
		if err != nil {
			moduleErr = fmt.Errorf("could not get current working directory: %w", err)
			return
		}

		dir := currentDir
		for {
			goModPath := filepath.Join(dir, "go.mod")
			if _, err := os.Stat(goModPath); err == nil {
				moduleRoot = dir
				break
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				moduleErr = fmt.Errorf("could not find go.mod in any parent directory of %s", currentDir)
				return
			}
			dir = parent
		}

		file, err := os.Open(filepath.Join(moduleRoot, "go.mod"))
		if err != nil {
			moduleErr = fmt.Errorf("could not open go.mod: %w", err)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		if scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "module ") {
				moduleName = strings.TrimSpace(strings.TrimPrefix(line, "module "))
			}
		}

		if moduleName == "" {
			moduleErr = fmt.Errorf("could not parse module name from go.mod")
		}
	})
	return moduleName, moduleRoot, moduleErr
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

	modName, modRoot, err := getGoModuleInfo()
	if err != nil {
		o.logger.Error("Could not determine Go module information, path normalization will fail.", "error", err)
		// Fallback to old behavior if we can't find module info
		modName = ""
	}

	for _, block := range blocks {
		normalizedPath := block.FileName

		// If the path is absolute and we have module info, normalize it.
		if modName != "" && filepath.IsAbs(normalizedPath) {
			relPath, err := filepath.Rel(modRoot, normalizedPath)
			if err == nil {
				// Join module name and relative path to get the canonical module path
				// e.g., "github.com/IgorBayerl/AdlerCov" + "cmd/main.go"
				normalizedPath = filepath.Join(modName, relPath)
			}
		}

		// Always convert to forward slashes for consistency.
		normalizedPath = filepath.ToSlash(normalizedPath)

		if !o.config.FileFilters().IsElementIncludedInReport(normalizedPath) {
			continue
		}

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
