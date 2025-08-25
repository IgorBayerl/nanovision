package parsers

import (
	"log/slog"
	"time"

	"github.com/IgorBayerl/AdlerCov/internal/filtering"
	"github.com/IgorBayerl/AdlerCov/internal/model"
)

type ParserResult struct {
	FileCoverage          []FileCoverage
	ParserName            string
	UnresolvedSourceFiles []string
	SourceDirectory       string
	Timestamp             *time.Time
}

type FileCoverage struct {
	Path  string
	Lines map[int]model.LineMetrics
}

type ParserConfig interface {
	SourceDirectories() []string
	FileFilters() filtering.IFilter
	Logger() *slog.Logger
}

type SimpleParserConfig struct {
	SrcDirs    []string
	FileFilter filtering.IFilter
	Log        *slog.Logger
}

func (sc *SimpleParserConfig) SourceDirectories() []string    { return sc.SrcDirs }
func (sc *SimpleParserConfig) FileFilters() filtering.IFilter { return sc.FileFilter }
func (sc *SimpleParserConfig) Logger() *slog.Logger           { return sc.Log }

type IParser interface {
	Name() string
	SupportsFile(filePath string) bool
	Parse(filePath string, config ParserConfig) (*ParserResult, error)
}
