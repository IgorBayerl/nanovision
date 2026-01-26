package parser_lcov

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/IgorBayerl/nanovision/internal/filereader"
	"github.com/IgorBayerl/nanovision/internal/parsers"
)

// LcovParser implements the parsers.IParser interface for LCOV (.info) reports.
type LcovParser struct {
	fileReader filereader.Reader
}

// NewLcovParser creates a new parser instance.
func NewLcovParser(fileReader filereader.Reader) parsers.IParser {
	return &LcovParser{
		fileReader: fileReader,
	}
}

// Name returns the unique, human-readable name of the parser.
func (p *LcovParser) Name() string {
	return "LCOV"
}

// SupportsFile performs a fast check to see if this parser can handle the file.
// It checks for the .info extension or the characteristic "SF:" or "TN:" markers.
func (p *LcovParser) SupportsFile(filePath string) bool {
	// Fast path: extension check
	if strings.HasSuffix(strings.ToLower(filePath), ".info") {
		return true
	}

	f, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer f.Close()

	// Content check: Look for LCOV headers in the first few lines
	scanner := bufio.NewScanner(f)
	for i := 0; i < 5 && scanner.Scan(); i++ {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "TN:") || strings.HasPrefix(line, "SF:") {
			return true
		}
	}

	return false
}

// Parse reads the LCOV report and transforms it into a list of FileCoverage objects.
func (p *LcovParser) Parse(filePath string, config parsers.ParserConfig) (*parsers.ParserResult, error) {
	logger := config.Logger().With(slog.String("parser", p.Name()), slog.String("file", filePath))

	lines, err := filereader.ReadLinesInFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read LCOV file %s: %w", filePath, err)
	}

	orchestrator := newProcessingOrchestrator(p.fileReader, config, logger)

	fileCoverage, unresolvedFiles := orchestrator.processLines(lines)

	return &parsers.ParserResult{
		FileCoverage:          fileCoverage,
		ParserName:            p.Name(),
		UnresolvedSourceFiles: unresolvedFiles,
	}, nil
}
