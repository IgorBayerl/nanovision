// Path: internal/parsers/gcov/processing.go
package gcov

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/IgorBayerl/AdlerCov/internal/filereader"
	"github.com/IgorBayerl/AdlerCov/internal/model"
	"github.com/IgorBayerl/AdlerCov/internal/parsers"
	"github.com/IgorBayerl/AdlerCov/internal/utils"
)

var (
	lineCoverageRegex   = regexp.MustCompile(`^\s*(?P<Visits>-|#####|=====|\d+):\s*(?P<LineNumber>[1-9]\d*):.*`)
	branchCoverageRegex = regexp.MustCompile(`^branch\s*\d+\s*(taken\s*(?P<Visits>\d+)%?|never executed)`)
)

// processingOrchestrator is now responsible for processing the lines of a single
// gcov report into a single FileCoverage object.
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

// processLines reads the content of a gcov file, extracts the source file path,
// and parses all line and branch coverage data into a single FileCoverage struct.
func (o *processingOrchestrator) processLines(lines []string) (*parsers.FileCoverage, []string, error) {
	if len(lines) == 0 {
		return nil, nil, fmt.Errorf("gcov file is empty")
	}

	firstLine := lines[0]
	if !strings.Contains(firstLine, "0:Source:") {
		return nil, nil, fmt.Errorf("invalid gcov format: first line does not contain '0:Source:'")
	}
	// Extract the source file path and normalize it.
	sourceFilePath := filepath.ToSlash(strings.TrimSpace(strings.SplitN(firstLine, "0:Source:", 2)[1]))

	if !o.config.FileFilters().IsElementIncludedInReport(sourceFilePath) {
		return nil, nil, nil // Return nil to indicate the file was skipped by filters.
	}

	// Try to resolve the source file. If it fails, we still proceed but will return
	// it in the unresolved files list.
	var unresolvedFiles []string
	if _, err := utils.FindFileInSourceDirs(sourceFilePath, o.config.SourceDirectories(), o.fileReader); err != nil {
		o.logger.Warn("Source file not found, it will be marked as unresolved.", "file", sourceFilePath, "error", err)
		unresolvedFiles = append(unresolvedFiles, sourceFilePath)
	}

	lineMetrics := make(map[int]model.LineMetrics)
	var lastCoverableLineNumber int

	for _, line := range lines {
		if match := lineCoverageRegex.FindStringSubmatch(line); match != nil {
			visitsText := match[1]
			lineNumber, _ := strconv.Atoi(match[2])
			lastCoverableLineNumber = lineNumber

			// A '-' visit count means the line is not executable code.
			if visitsText == "-" {
				continue
			}

			metric := model.LineMetrics{}
			// '#####' or '=====' means the line is executable but was not covered.
			if visitsText != "#####" && visitsText != "=====" {
				metric.Hits, _ = strconv.Atoi(visitsText)
			}
			lineMetrics[lineNumber] = metric

		} else if match := branchCoverageRegex.FindStringSubmatch(line); match != nil {
			if metric, ok := lineMetrics[lastCoverableLineNumber]; ok {
				metric.TotalBranches++
				// The branch was taken if the "Visits" group (index 2) is not empty.
				// "taken 0%" still means the branch was not taken.
				if len(match) > 2 && match[2] != "" && match[2] != "0" {
					metric.CoveredBranches++
				}
				lineMetrics[lastCoverableLineNumber] = metric
			}
		}
	}

	fileCoverage := &parsers.FileCoverage{
		Path:  sourceFilePath,
		Lines: lineMetrics,
	}

	return fileCoverage, unresolvedFiles, nil
}
