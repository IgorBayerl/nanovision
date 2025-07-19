package gocover

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/filereader"
	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/parsers"
)

var (
	// Regex to parse a Go coverage line, e.g., "file.go:1.2,3.4 5 6"
	goCoverLineRegex = regexp.MustCompile(`^(.+):(\d+)\.(\d+),(\d+)\.(\d+)\s(\d+)\s(\d+)$`)
)

// GoCoverParser implements the parsers.IParserinterface for Go coverage reports.
type GoCoverParser struct {
	fileReader filereader.Reader // Injected dependency
}

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

// NewGoCoverParser creates a new parser instance.
func NewGoCoverParser(fileReader filereader.Reader) parsers.IParser {
	return &GoCoverParser{
		fileReader: fileReader,
	}
}

// Name returns the unique, human-readable name of the parsers.
func (p *GoCoverParser) Name() string {
	return "GoCover"
}

// SupportsFile performs a fast check to see if this parser can handle the file.
func (p *GoCoverParser) SupportsFile(filePath string) bool {
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

	return strings.HasPrefix(firstLine, "mode:")
}

// Parse reads the entire Go coverage report, transforms it into `GoCoverProfileBlock`s,
// and then delegates the complex processing to the processingOrchestrator.
func (p *GoCoverParser) Parse(filePath string, config parsers.ParserConfig) (*parsers.ParserResult, error) {
	logger := config.Logger().With(slog.String("parser", p.Name()), slog.String("file", filePath))

	profileBlocks, err := p.loadAndParseGoCoverFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load/parse Go coverage file from %s: %w", filePath, err)
	}

	orchestrator := newProcessingOrchestrator(p.fileReader, config, logger)

	assemblies, err := orchestrator.processBlocks(profileBlocks)
	if err != nil {
		return nil, fmt.Errorf("failed to process Go coverage blocks: %w", err)
	}

	return &parsers.ParserResult{
		Assemblies:             assemblies,
		SourceDirectories:      []string{}, // Go cover files don't list source directories
		SupportsBranchCoverage: false,
		ParserName:             p.Name(),
		MinimumTimeStamp:       nil,
		MaximumTimeStamp:       nil,
	}, nil
}

// loadAndParseGoCoverFile reads the specified file line-by-line and parses each
// valid coverage data line into a GoCoverProfileBlock.
func (p *GoCoverParser) loadAndParseGoCoverFile(path string) ([]GoCoverProfileBlock, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	var blocks []GoCoverProfileBlock
	scanner := bufio.NewScanner(file)

	// Skip the first line ("mode: ...")
	if !scanner.Scan() {
		return nil, fmt.Errorf("file is empty or could not be read")
	}

	for scanner.Scan() {
		line := scanner.Text()
		match := goCoverLineRegex.FindStringSubmatch(line)

		if len(match) == 8 {
			startLine, _ := strconv.Atoi(match[2])
			startCol, _ := strconv.Atoi(match[3])
			endLine, _ := strconv.Atoi(match[4])
			endCol, _ := strconv.Atoi(match[5])
			numStatements, _ := strconv.Atoi(match[6])
			hitCount, _ := strconv.Atoi(match[7])

			blocks = append(blocks, GoCoverProfileBlock{
				FileName:      match[1],
				StartLine:     startLine,
				StartCol:      startCol,
				EndLine:       endLine,
				EndCol:        endCol,
				NumStatements: numStatements,
				HitCount:      hitCount,
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return blocks, nil
}
