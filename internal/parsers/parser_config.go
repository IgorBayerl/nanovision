package parsers

import (
	"log/slog"
	"time"

	"github.com/IgorBayerl/AdlerCov/internal/filtering"
	"github.com/IgorBayerl/AdlerCov/internal/language"
	"github.com/IgorBayerl/AdlerCov/internal/model"
	"github.com/IgorBayerl/AdlerCov/internal/settings"
)

type ParserResult struct {
	FileCoverage          []FileCoverage
	ParserName            string
	UnresolvedSourceFiles []string
	SourceDirectory       string // <-- CORRECTED: Changed from plural to singular
	Timestamp             *time.Time
}

type FileCoverage struct {
	Path  string
	Lines map[int]model.LineMetrics
}

type ParserConfig interface {
	SourceDirectories() []string
	AssemblyFilters() filtering.IFilter
	ClassFilters() filtering.IFilter
	FileFilters() filtering.IFilter
	Settings() *settings.Settings
	Logger() *slog.Logger
	LanguageProcessorFactory() *language.ProcessorFactory
}

type IParser interface {
	Name() string
	SupportsFile(filePath string) bool
	Parse(filePath string, config ParserConfig) (*ParserResult, error)
}
