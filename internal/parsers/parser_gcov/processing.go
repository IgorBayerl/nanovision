package parser_gcov

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/IgorBayerl/AdlerCov/filereader"
	"github.com/IgorBayerl/AdlerCov/internal/model"
	"github.com/IgorBayerl/AdlerCov/internal/parsers"
	"github.com/IgorBayerl/AdlerCov/internal/utils"
)

var (
	lineCoverageRegex   = regexp.MustCompile(`^\s*(?P<Visits>-|#####|=====|\d+):\s*(?P<LineNumber>[1-9]\d*):.*`)
	branchCoverageRegex = regexp.MustCompile(`^branch\s*\d+\s*(taken\s*(?P<Visits>\d+)%?|never executed)`)
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

func (o *processingOrchestrator) processLines(lines []string) (*parsers.FileCoverage, []string, error) {
	if len(lines) == 0 {
		return nil, nil, fmt.Errorf("gcov file is empty")
	}

	firstLine := lines[0]
	if !strings.Contains(firstLine, "0:Source:") {
		return nil, nil, fmt.Errorf("invalid gcov format: first line does not contain '0:Source:'")
	}

	sourceFilePathFromReport := filepath.ToSlash(strings.TrimSpace(strings.SplitN(firstLine, "0:Source:", 2)[1]))
	sourceDirs := o.config.SourceDirectories()

	// We will pass the original, potentially absolute path to the builder.
	// The builder will resolve it against the source_dir and project_root.
	displayPath := sourceFilePathFromReport

	var unresolvedFiles []string
	// Pass the logger from the orchestrator into the find utility
	if _, err := utils.FindFileInSourceDirs(sourceFilePathFromReport, sourceDirs, o.fileReader, o.logger); err != nil {
		o.logger.Warn("Source file not found, it will be marked as unresolved.", "file", sourceFilePathFromReport, "error", err)
		unresolvedFiles = append(unresolvedFiles, sourceFilePathFromReport)
	}

	lineMetrics := make(map[int]model.LineMetrics)
	var lastCoverableLineNumber int

	for _, line := range lines {
		if match := lineCoverageRegex.FindStringSubmatch(line); match != nil {
			visitsText := match[1]
			lineNumber, _ := strconv.Atoi(match[2])
			lastCoverableLineNumber = lineNumber
			if visitsText == "-" {
				continue
			}
			metric := model.LineMetrics{}
			if visitsText != "#####" && visitsText != "=====" {
				metric.Hits, _ = strconv.Atoi(visitsText)
			}
			lineMetrics[lineNumber] = metric
		} else if match := branchCoverageRegex.FindStringSubmatch(line); match != nil {
			if metric, ok := lineMetrics[lastCoverableLineNumber]; ok {
				metric.TotalBranches++
				if len(match) > 2 && match[2] != "" && match[2] != "0" {
					metric.CoveredBranches++
				}
				lineMetrics[lastCoverableLineNumber] = metric
			}
		}
	}

	fileCoverage := &parsers.FileCoverage{
		Path:  displayPath,
		Lines: lineMetrics,
	}

	return fileCoverage, unresolvedFiles, nil
}
