package parser_lcov

import (
	"log/slog"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/IgorBayerl/nanovision/internal/filereader"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/parsers"
	"github.com/IgorBayerl/nanovision/internal/utils"
)

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

// processLines iterates through the lines of the LCOV file and builds coverage data.
func (o *processingOrchestrator) processLines(lines []string) ([]parsers.FileCoverage, []string) {
	var allFileCoverage []parsers.FileCoverage
	var allUnresolvedFiles []string

	// State variables for the current record
	var currentPath string
	currentLines := make(map[int]model.LineMetrics)
	inRecord := false

	finishRecord := func() {
		if currentPath == "" {
			return
		}

		// Resolve the file path against source directories
		sourceDirs := o.config.SourceDirectories()
		// Default to empty slice if not provided to avoid nil issues, though logic handles it
		if len(sourceDirs) == 0 {
			sourceDirs = []string{""}
		}

		// Use the utils helper to find the actual file on disk
		_, err := utils.FindFileInSourceDirs(currentPath, sourceDirs, o.fileReader, o.logger)
		if err != nil {
			o.logger.Warn("Source file not found, marking as unresolved", "file", currentPath)
			allUnresolvedFiles = append(allUnresolvedFiles, currentPath)
		}

		allFileCoverage = append(allFileCoverage, parsers.FileCoverage{
			Path:  filepath.ToSlash(currentPath),
			Lines: currentLines,
		})

		// Reset state
		currentPath = ""
		currentLines = make(map[int]model.LineMetrics)
		inRecord = false
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Handle Source File (Start of record)
		if strings.HasPrefix(line, "SF:") {
			// If we were already in a record without an end_of_record (malformed), finish it
			if inRecord {
				finishRecord()
			}
			currentPath = strings.TrimPrefix(line, "SF:")
			inRecord = true
			continue
		}

		if !inRecord {
			continue
		}

		// Handle Line Data: DA:<lineNumber>,<hits>[,<checksum>]
		if after, ok := strings.CutPrefix(line, "DA:"); ok {
			parts := strings.Split(after, ",")
			if len(parts) >= 2 {
				ln, err1 := strconv.Atoi(parts[0])
				hits, err2 := strconv.Atoi(parts[1])

				if err1 == nil && err2 == nil && ln > 0 {
					metric := currentLines[ln]
					metric.Hits += hits // Add to existing in case of duplicate entries
					currentLines[ln] = metric
				}
			}
			continue
		}

		// Handle End of Record
		if line == "end_of_record" {
			finishRecord()
		}
	}

	// Handle trailing record if file didn't end with newline/end_of_record
	if inRecord {
		finishRecord()
	}

	return allFileCoverage, allUnresolvedFiles
}
