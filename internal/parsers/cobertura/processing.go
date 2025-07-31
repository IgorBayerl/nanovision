// Path: internal/parsers/cobertura/processing.go
package cobertura

import (
	"log/slog"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/IgorBayerl/AdlerCov/internal/filereader"
	"github.com/IgorBayerl/AdlerCov/internal/model"
	"github.com/IgorBayerl/AdlerCov/internal/parsers"
	"github.com/IgorBayerl/AdlerCov/internal/utils"
)

var (
	conditionCoverageRegexCobertura = regexp.MustCompile(`\((?P<NumberOfCoveredBranches>\d+)/(?P<NumberOfTotalBranches>\d+)\)$`)
)

// processingOrchestrator is responsible for converting the raw XML data into
// a flat list of per-file coverage metrics. It handles the logic of grouping
// line data by filename as it appears across different classes in the report.
type processingOrchestrator struct {
	fileReader filereader.Reader
	config     parsers.ParserConfig
	logger     *slog.Logger
}

func newProcessingOrchestrator(
	fileReader filereader.Reader,
	config parsers.ParserConfig,
	logger *slog.Logger,
) *processingOrchestrator {
	return &processingOrchestrator{
		fileReader: fileReader,
		config:     config,
		logger:     logger,
	}
}

// processPackages is the main entry point for the orchestrator. It iterates through
// the XML packages and classes to build a map of coverage data keyed by file path.
func (o *processingOrchestrator) processPackages(packages []PackageXML) ([]parsers.FileCoverage, []string) {
	fileData := make(map[string]map[int]model.LineMetrics)
	var unresolvedFiles []string

	for _, pkgXML := range packages {
		for _, classXML := range pkgXML.Classes.Class {
			filePath := filepath.ToSlash(classXML.Filename)
			if filePath == "" || !o.config.FileFilters().IsElementIncludedInReport(filePath) {
				continue
			}
			if _, ok := fileData[filePath]; !ok {
				fileData[filePath] = make(map[int]model.LineMetrics)
			}
			allLinesInClass := classXML.Lines.Line
			for _, methodXML := range classXML.Methods.Method {
				allLinesInClass = append(allLinesInClass, methodXML.Lines.Line...)
			}
			o.mergeLinesIntoFile(fileData[filePath], allLinesInClass)
		}
	}

	var finalFileCoverage []parsers.FileCoverage
	sourceDir := ""
	if len(o.config.SourceDirectories()) > 0 {
		sourceDir = o.config.SourceDirectories()[0]
	}

	for path, lines := range fileData {
		// Pass the logger from the orchestrator into the find utility
		if _, err := utils.FindFileInSourceDirs(path, []string{sourceDir}, o.fileReader, o.logger); err != nil {
			o.logger.Warn("Source file not found, it will be marked as unresolved.", "file", path, "error", err)
			unresolvedFiles = append(unresolvedFiles, path)
		}

		finalFileCoverage = append(finalFileCoverage, parsers.FileCoverage{
			Path:  path,
			Lines: lines,
		})
	}

	return finalFileCoverage, unresolvedFiles
}

// mergeLinesIntoFile processes a list of XML line elements and merges their
// data into a map of line metrics for a specific file. This handles cases
// where the same line might appear multiple times in a report.
func (o *processingOrchestrator) mergeLinesIntoFile(lineMetrics map[int]model.LineMetrics, linesXML []LineXML) {
	for _, lineXML := range linesXML {
		lineNumber, err := strconv.Atoi(lineXML.Number)
		if err != nil || lineNumber <= 0 {
			continue
		}

		hits, err := strconv.Atoi(lineXML.Hits)
		if err != nil {
			continue // Skip lines without a valid hit count.
		}

		existingMetric := lineMetrics[lineNumber]
		existingMetric.Hits += hits // Sum hits from different report sections.

		// Process branch coverage.
		if lineXML.Branch == "true" {
			covered, total := o.parseBranchData(lineXML)
			existingMetric.CoveredBranches += covered
			existingMetric.TotalBranches += total
		}

		lineMetrics[lineNumber] = existingMetric
	}
}

// parseBranchData extracts the number of covered and total branches from a single line's
// `condition-coverage` attribute.
func (o *processingOrchestrator) parseBranchData(lineXML LineXML) (covered, total int) {
	matches := conditionCoverageRegexCobertura.FindStringSubmatch(lineXML.ConditionCoverage)
	if len(matches) == 3 {
		covered, _ = strconv.Atoi(matches[1])
		total, _ = strconv.Atoi(matches[2])
	}
	return
}
