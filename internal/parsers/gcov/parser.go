package gcov

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/IgorBayerl/AdlerCov/internal/filereader"
	"github.com/IgorBayerl/AdlerCov/internal/parsers"
)

// GCovParser implements the parsers.IParser interface for gcov reports.
type GCovParser struct {
	fileReader filereader.Reader
}

// NewGCovParser creates a new parser instance.
func NewGCovParser(fileReader filereader.Reader) parsers.IParser {
	return &GCovParser{
		fileReader: fileReader,
	}
}

// Name returns the unique, human-readable name of the parser.
func (p *GCovParser) Name() string {
	return "GCov"
}

// SupportsFile checks if the file is a gcov report by reading its first line
// for the characteristic "0:Source:" marker.
func (p *GCovParser) SupportsFile(filePath string) bool {
	f, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	firstLine, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return false
	}

	return strings.Contains(firstLine, "0:Source:")
}

// Parse reads the gcov report from disk and delegates processing to the orchestrator.
// It then returns the resulting flat list of file coverage data.
func (p *GCovParser) Parse(filePath string, config parsers.ParserConfig) (*parsers.ParserResult, error) {
	logger := config.Logger().With(slog.String("parser", p.Name()), slog.String("file", filePath))

	lines, err := filereader.ReadLinesInFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read gcov file %s: %w", filePath, err)
	}

	orchestrator := newProcessingOrchestrator(p.fileReader, config, logger)

	// Since a gcov file maps to a single source file, the orchestrator returns
	// a single FileCoverage object or nil if the file is invalid/unresolved.
	fileCoverage, unresolvedFiles, err := orchestrator.processLines(lines)
	if err != nil {
		return nil, err
	}

	var coverageData []parsers.FileCoverage
	if fileCoverage != nil {
		coverageData = append(coverageData, *fileCoverage)
	}

	return &parsers.ParserResult{
		FileCoverage:          coverageData,
		ParserName:            p.Name(),
		UnresolvedSourceFiles: unresolvedFiles,
	}, nil
}
