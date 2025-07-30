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
	SourceDirectories     []string
	Timestamp             *time.Time
	// Other report-level metadata can go here.
}

type FileCoverage struct {
	Path  string                    // Project-relative path (e.g., "internal/analyzer/analyzer.go")
	Lines map[int]model.LineMetrics // Raw line data from the report
}
type ParserResultOld struct {
	Assemblies             []model.Assembly
	SourceDirectories      []string
	SupportsBranchCoverage bool
	ParserName             string
	MinimumTimeStamp       *time.Time
	MaximumTimeStamp       *time.Time
	UnresolvedSourceFiles  []string
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
