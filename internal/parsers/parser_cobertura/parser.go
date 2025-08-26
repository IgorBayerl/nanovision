package parser_cobertura

import (
	"encoding/xml"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/IgorBayerl/AdlerCov/filereader"
	"github.com/IgorBayerl/AdlerCov/internal/parsers"
	"github.com/IgorBayerl/AdlerCov/internal/utils"
)

// CoberturaParser implements the parsers.IParser interface for Cobertura XML reports.
type CoberturaParser struct {
	fileReader filereader.Reader
}

func NewCoberturaParser(fileReader filereader.Reader) parsers.IParser {
	return &CoberturaParser{
		fileReader: fileReader,
	}
}

func (p *CoberturaParser) Name() string {
	return "Cobertura"
}

// SupportsFile performs a fast check to see if this parser can handle the given file.
// It verifies the file has a ".xml" extension and that its root element is "<coverage>".
func (p *CoberturaParser) SupportsFile(filePath string) bool {
	if !strings.HasSuffix(strings.ToLower(filePath), ".xml") {
		return false
	}

	f, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer f.Close()

	// We only need to check the very first token to identify the root element.
	decoder := xml.NewDecoder(f)
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			return false // Reached end of file without finding any elements.
		}
		if err != nil {
			return false // Malformed XML.
		}

		if se, ok := token.(xml.StartElement); ok {
			return se.Name.Local == "coverage"
		}
	}
}

// Parse unmarshals the Cobertura XML report and delegates the conversion to a
// flat list of FileCoverage objects to the processingOrchestrator.
func (p *CoberturaParser) Parse(filePath string, config parsers.ParserConfig) (*parsers.ParserResult, error) {
	logger := config.Logger().With(slog.String("parser", p.Name()), slog.String("file", filePath))

	rawReport, err := p.loadAndUnmarshalCoberturaXML(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load/unmarshal Cobertura XML from %s: %w", filePath, err)
	}

	orchestrator := newProcessingOrchestrator(p.fileReader, config, logger)

	// The orchestrator now directly returns the flat list of file coverage data and any unresolved files.
	fileCoverage, unresolvedFiles := orchestrator.processPackages(rawReport.Packages.Package)

	timestamp := p.getReportTimestamp(rawReport.Timestamp, logger)

	return &parsers.ParserResult{
		FileCoverage:          fileCoverage,
		ParserName:            p.Name(),
		Timestamp:             timestamp,
		UnresolvedSourceFiles: unresolvedFiles,
	}, nil
}

// ------ Helper Functions ------

// getReportTimestamp parses the Cobertura timestamp string into a *time.Time object.
func (p *CoberturaParser) getReportTimestamp(rawTimestamp string, logger *slog.Logger) *time.Time {
	if rawTimestamp == "" {
		return nil
	}
	parsedTs, err := strconv.ParseInt(rawTimestamp, 10, 64)
	if err != nil {
		logger.Warn("Failed to parse Cobertura timestamp", "timestamp", rawTimestamp, "error", err)
		return nil
	}

	// Handle timestamps in milliseconds vs. seconds by checking if it's a valid seconds timestamp.
	if !utils.IsValidUnixSeconds(parsedTs) && utils.IsValidUnixSeconds(parsedTs/1000) {
		parsedTs /= 1000
	}

	if utils.IsValidUnixSeconds(parsedTs) {
		t := time.Unix(parsedTs, 0)
		return &t
	}

	logger.Warn("Cobertura timestamp is outside the valid range", "timestamp", rawTimestamp)
	return nil
}

// loadAndUnmarshalCoberturaXML reads and unmarshals the Cobertura XML file.
func (p *CoberturaParser) loadAndUnmarshalCoberturaXML(path string) (*CoberturaRoot, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	bytes, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	var rawReport CoberturaRoot
	if err := xml.Unmarshal(bytes, &rawReport); err != nil {
		return nil, fmt.Errorf("unmarshal xml: %w", err)
	}
	return &rawReport, nil
}
