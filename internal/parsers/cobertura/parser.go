package cobertura

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/IgorBayerl/AdlerCov/internal/filereader"
	"github.com/IgorBayerl/AdlerCov/internal/parsers"
	"github.com/IgorBayerl/AdlerCov/internal/utils"
)

// CoberturaParser implements the parsers.IParser interface for Cobertura XML reports.
type CoberturaParser struct {
	fileReader filereader.Reader
}

// DefaultFileReader is the production implementation of the FileReader interface.
type DefaultFileReader struct{}

func (dfr *DefaultFileReader) ReadFile(path string) ([]string, error) {
	return filereader.ReadLinesInFile(path)
}

func (dfr *DefaultFileReader) CountLines(path string) (int, error) {
	return filereader.CountLinesInFile(path)
}

func (dfr *DefaultFileReader) Stat(name string) (fs.FileInfo, error) {
	return os.Stat(name)
}

func NewCoberturaParser(fileReader filereader.Reader) parsers.IParser {
	return &CoberturaParser{
		fileReader: fileReader,
	}
}

func (cp *CoberturaParser) Name() string {
	return "Cobertura"
}

func (cp *CoberturaParser) SupportsFile(filePath string) bool {
	if !strings.HasSuffix(strings.ToLower(filePath), ".xml") {
		return false
	}

	f, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer f.Close()

	decoder := xml.NewDecoder(f)
	firstElementChecked := false
	isCoverageFile := false

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			// Return true only if we found a <coverage> root and no errors occurred.
			return isCoverageFile
		}
		if err != nil {
			// Any XML parsing error (like a malformed tag) means the file is not supported.
			return false
		}

		if se, ok := token.(xml.StartElement); ok {
			if !firstElementChecked {
				// This is the first element tag we have encountered.
				firstElementChecked = true
				if se.Name.Local == "coverage" {
					isCoverageFile = true
				} else {
					// The root element is not <coverage>, so this is not a supported file
					return false
				}
			}
		}
	}
}

// Parse is the main entry point for the Cobertura parsers. It unmarshals the XML
// and delegates the complex processing logic to the processingOrchestrator, which
// handles per-file language detection and formatting.
func (cp *CoberturaParser) Parse(filePath string, config parsers.ParserConfig) (*parsers.ParserResult, error) {
	logger := config.Logger().With(slog.String("parser", cp.Name()), slog.String("file", filePath))

	rawReport, sourceDirsFromXML, err := cp.loadAndUnmarshalCoberturaXML(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load/unmarshal Cobertura XML from %s: %w", filePath, err)
	}

	effectiveSourceDirs := cp.getEffectiveSourceDirs(config, sourceDirsFromXML)

	// The orchestrator is now simpler and does not take a pre-determined formatter.
	// It will determine the formatter for each file internally.
	orchestrator := newProcessingOrchestrator(cp.fileReader, config, effectiveSourceDirs, logger)

	assemblies, detectedBranchSupport, err := orchestrator.processPackages(rawReport.Packages.Package)
	if err != nil {
		return nil, fmt.Errorf("failed to process Cobertura packages: %w", err)
	}

	timestamp := cp.getReportTimestamp(rawReport.Timestamp, logger)

	return &parsers.ParserResult{
		Assemblies:             assemblies,
		SourceDirectories:      sourceDirsFromXML,
		SupportsBranchCoverage: detectedBranchSupport,
		ParserName:             cp.Name(),
		MinimumTimeStamp:       timestamp,
		MaximumTimeStamp:       timestamp,
		UnresolvedSourceFiles:  orchestrator.unresolvedSourceFiles,
	}, nil
}

// ------ Helper Functions ------

// getEffectiveSourceDirs combines source directories from the configuration (CLI)
// and from the XML file's <sources> tag to create a comprehensive list of search paths.
func (cp *CoberturaParser) getEffectiveSourceDirs(config parsers.ParserConfig, sourceDirsFromXML []string) []string {
	sourceDirsSet := make(map[string]struct{})

	for _, dir := range config.SourceDirectories() {
		if dir != "" {
			sourceDirsSet[dir] = struct{}{}
		}
	}

	for _, dir := range sourceDirsFromXML {
		if dir != "" {
			sourceDirsSet[dir] = struct{}{}
		}
	}

	var effectiveSourceDirs []string
	for dir := range sourceDirsSet {
		effectiveSourceDirs = append(effectiveSourceDirs, dir)
	}

	return effectiveSourceDirs
}

// getReportTimestamp parses the Cobertura timestamp string into a *time.Time object.
func (cp *CoberturaParser) getReportTimestamp(rawTimestamp string, logger *slog.Logger) *time.Time {
	if rawTimestamp == "" {
		return nil
	}
	parsedTs, err := strconv.ParseInt(rawTimestamp, 10, 64)
	if err != nil {
		logger.Warn("Failed to parse Cobertura timestamp", "timestamp", rawTimestamp, "error", err)
		return nil
	}

	// Handle timestamps in milliseconds vs. seconds
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
func (cp *CoberturaParser) loadAndUnmarshalCoberturaXML(path string) (*CoberturaRoot, []string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	bytes, err := io.ReadAll(f)
	if err != nil {
		return nil, nil, fmt.Errorf("read file: %w", err)
	}

	var rawReport CoberturaRoot
	if err := xml.Unmarshal(bytes, &rawReport); err != nil {
		return nil, nil, fmt.Errorf("unmarshal xml: %w", err)
	}
	return &rawReport, rawReport.Sources.Source, nil
}
