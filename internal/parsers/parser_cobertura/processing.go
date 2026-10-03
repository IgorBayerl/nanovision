package parser_cobertura

import (
	"log/slog"
	"path/filepath"
	"strconv"

	"github.com/IgorBayerl/nanovision/internal/filereader"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/parsers"
	"github.com/IgorBayerl/nanovision/internal/utils"
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

		lineMetrics[lineNumber] = existingMetric
	}
}
