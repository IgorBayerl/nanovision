package gcov

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/IgorBayerl/AdlerCov/internal/filereader"
	"github.com/IgorBayerl/AdlerCov/internal/model"
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

// SupportsFile checks if the file is a gcov report by reading its first line.
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

// Parse reads the gcov report from disk and processes it.
func (p *GCovParser) Parse(filePath string, config parsers.ParserConfig) (*parsers.ParserResult, error) {
	logger := config.Logger().With(slog.String("parser", p.Name()), slog.String("file", filePath))

	// Read the actual report file from the filesystem, NOT from the mock filereader.
	lines, err := filereader.ReadLinesInFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read gcov file %s: %w", filePath, err)
	}

	orchestrator := newProcessingOrchestrator(p.fileReader, config, logger)

	assembly, err := orchestrator.processLines(lines)
	if err != nil {
		return nil, err
	}

	assemblies := []model.Assembly{}
	if assembly != nil {
		assemblies = append(assemblies, *assembly)
	}

	return &parsers.ParserResult{
		Assemblies:             assemblies,
		SupportsBranchCoverage: orchestrator.detectedBranchCoverage,
		ParserName:             p.Name(),
		UnresolvedSourceFiles:  orchestrator.unresolvedSourceFiles,
	}, nil
}
