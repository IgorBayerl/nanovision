package parsers

import (
	"log/slog"
	"time"

	"github.com/IgorBayerl/AdlerCov/internal/filtering"
	"github.com/IgorBayerl/AdlerCov/internal/language"
	"github.com/IgorBayerl/AdlerCov/internal/model"
)

// ParserResult holds the raw, translated data from a single report file.
type ParserResult struct {
	FileCoverage          []FileCoverage
	ParserName            string
	UnresolvedSourceFiles []string
	SourceDirectory       string
	Timestamp             *time.Time
}

// FileCoverage holds the raw line and branch metrics for a single source file.
type FileCoverage struct {
	Path  string
	Lines map[int]model.LineMetrics
}

// ParserConfig defines the lean contract for configuration needed by a parser during its operation.
type ParserConfig interface {
	SourceDirectories() []string
	FileFilters() filtering.IFilter
	Logger() *slog.Logger
	LanguageProcessorFactory() *language.ProcessorFactory
}

// SimpleParserConfig is a basic, concrete implementation of ParserConfig for individual parse tasks.
type SimpleParserConfig struct {
	SrcDirs     []string
	FileFilter  filtering.IFilter
	Log         *slog.Logger
	LangFactory *language.ProcessorFactory
}

func (sc *SimpleParserConfig) SourceDirectories() []string    { return sc.SrcDirs }
func (sc *SimpleParserConfig) FileFilters() filtering.IFilter { return sc.FileFilter }
func (sc *SimpleParserConfig) Logger() *slog.Logger           { return sc.Log }
func (sc *SimpleParserConfig) LanguageProcessorFactory() *language.ProcessorFactory {
	return sc.LangFactory
}

// IParser defines the contract that all report parsers must implement.
type IParser interface {
	Name() string
	SupportsFile(filePath string) bool
	Parse(filePath string, config ParserConfig) (*ParserResult, error)
}
