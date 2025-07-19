package parsers

import (
	"log/slog"
	"time"

	"github.com/IgorBayerl/AdlerCov/internal/filtering"
	"github.com/IgorBayerl/AdlerCov/internal/language"
	"github.com/IgorBayerl/AdlerCov/internal/model"
	"github.com/IgorBayerl/AdlerCov/internal/settings"
)

// holds the processed data from a single coverage report.
type ParserResult struct {
	Assemblies             []model.Assembly
	SourceDirectories      []string
	SupportsBranchCoverage bool
	ParserName             string
	MinimumTimeStamp       *time.Time
	MaximumTimeStamp       *time.Time
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
