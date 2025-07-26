package testutil

import (
	"io"
	"log/slog"

	"github.com/IgorBayerl/AdlerCov/internal/filtering"
	"github.com/IgorBayerl/AdlerCov/internal/language"
	"github.com/IgorBayerl/AdlerCov/internal/language/csharp"
	"github.com/IgorBayerl/AdlerCov/internal/language/defaultformatter"
	"github.com/IgorBayerl/AdlerCov/internal/language/gcc"
	"github.com/IgorBayerl/AdlerCov/internal/language/golang"
	"github.com/IgorBayerl/AdlerCov/internal/settings"
)

// MockParserConfig implements ParserConfig for providing test configuration.
type MockParserConfig struct {
	SrcDirs     []string
	AsmFilter   filtering.IFilter
	ClsFilter   filtering.IFilter
	FileFilter  filtering.IFilter
	SettingsObj *settings.Settings
	Log         *slog.Logger
	LangFactory *language.ProcessorFactory
}

func (m *MockParserConfig) SourceDirectories() []string        { return m.SrcDirs }
func (m *MockParserConfig) AssemblyFilters() filtering.IFilter { return m.AsmFilter }
func (m *MockParserConfig) ClassFilters() filtering.IFilter    { return m.ClsFilter }
func (m *MockParserConfig) FileFilters() filtering.IFilter     { return m.FileFilter }
func (m *MockParserConfig) Settings() *settings.Settings       { return m.SettingsObj }
func (m *MockParserConfig) Logger() *slog.Logger               { return m.Log }
func (m *MockParserConfig) LanguageProcessorFactory() *language.ProcessorFactory {
	return m.LangFactory
}

func NewTestConfig(sourceDirs []string) *MockParserConfig {
	noFilter, _ := filtering.NewDefaultFilter(nil)
	langFactory := language.NewProcessorFactory(
		defaultformatter.NewDefaultProcessor(),
		golang.NewGoProcessor(),
		csharp.NewCSharpProcessor(),
		gcc.NewGccProcessor(),
	)

	return &MockParserConfig{
		SrcDirs:     sourceDirs,
		AsmFilter:   noFilter,
		ClsFilter:   noFilter,
		FileFilter:  noFilter,
		SettingsObj: settings.NewSettings(),
		Log:         slog.New(slog.NewTextHandler(io.Discard, nil)),
		LangFactory: langFactory,
	}
}
