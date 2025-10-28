package parser_cobertura

import (
	"log/slog"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/IgorBayerl/nanovision/filereader"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/parsers"
	"github.com/IgorBayerl/nanovision/internal/utils"
)

var (
	conditionCoverageRegexCobertura = regexp.MustCompile(`\((?P<NumberOfCoveredBranches>\d+)/(?P<NumberOfTotalBranches>\d+)\)$`)
)

// processingOrchestrator is responsible for converting the raw XML data into
// a flat list of per-file coverage metrics.
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

// processPackages is the main entry point for the orchestrator.
func (o *processingOrchestrator) processPackages(packages []PackageXML) ([]parsers.FileCoverage, []string) {
	fileData := make(map[string]map[int]model.LineMetrics)
	var unresolvedFiles []string

	for _, pkgXML := range packages {
		for _, classXML := range pkgXML.Classes.Class {
			filePath := filepath.ToSlash(classXML.Filename)
			if filePath == "" {
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
// data into a map of line metrics for a specific file.
func (o *processingOrchestrator) mergeLinesIntoFile(lineMetrics map[int]model.LineMetrics, linesXML []LineXML) {
	for _, lineXML := range linesXML {
		lineNumber, err := strconv.Atoi(lineXML.Number)
		if err != nil || lineNumber <= 0 {
			continue
		}

		hits, err := strconv.Atoi(lineXML.Hits)
		if err != nil {
			continue
		}

		existingMetric := lineMetrics[lineNumber]
		existingMetric.Hits += hits

		if lineXML.Branch == "true" || strings.ToLower(lineXML.Branch) == "true" {
			covered, total := o.parseBranchData(lineXML)
			existingMetric.CoveredBranches += covered
			existingMetric.TotalBranches += total
		}

		lineMetrics[lineNumber] = existingMetric
	}
}

// parseBranchData is the updated function that handles both Cobertura styles.
// It prioritizes the explicit 'condition-coverage' attribute and falls back
// to the nested '<conditions>' block only if necessary.
func (o *processingOrchestrator) parseBranchData(lineXML LineXML) (covered, total int) {
	// STRATEGY 1: Prioritize the 'condition-coverage' attribute.
	// This format is explicit (e.g., "50% (1/2)") and used by both gcovr (C++) and coverlet (C#).
	// It is the most reliable source for the line's overall branch statistics.
	matches := conditionCoverageRegexCobertura.FindStringSubmatch(lineXML.ConditionCoverage)
	if len(matches) == 3 {
		covered, _ = strconv.Atoi(matches[1])
		total, _ = strconv.Atoi(matches[2])
		// If we successfully parsed this attribute, we trust it and are done.
		return covered, total
	}

	// STRATEGY 2: Fallback to counting <condition> elements.
	// Some tools might omit the 'condition-coverage' attribute and only provide the detailed block.
	if len(lineXML.Conditions.Condition) > 0 {
		total = len(lineXML.Conditions.Condition)
		covered = 0
		for _, condition := range lineXML.Conditions.Condition {
			// A condition is considered covered if its coverage is 100%.
			if strings.HasPrefix(condition.Coverage, "100%") {
				covered++
			}
		}
		return covered, total
	}

	// If neither method yields data, return 0, 0.
	return 0, 0
}
