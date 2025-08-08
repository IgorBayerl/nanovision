package testutil

import (
	"io"
	"log/slog"

	"github.com/IgorBayerl/AdlerCov/internal/filtering"
	"github.com/IgorBayerl/AdlerCov/internal/language"
	"github.com/IgorBayerl/AdlerCov/internal/language/lang_cpp"
	"github.com/IgorBayerl/AdlerCov/internal/language/lang_csharp"
	"github.com/IgorBayerl/AdlerCov/internal/language/lang_default"
	"github.com/IgorBayerl/AdlerCov/internal/language/lang_go"
	"github.com/IgorBayerl/AdlerCov/internal/parsers"
)

// MockParserConfig implements the lean parsers.ParserConfig interface for testing.
type MockParserConfig struct {
	SrcDirs     []string
	FileFilter  filtering.IFilter
	Log         *slog.Logger
	LangFactory *language.ProcessorFactory
}

func (m *MockParserConfig) SourceDirectories() []string    { return m.SrcDirs }
func (m *MockParserConfig) FileFilters() filtering.IFilter { return m.FileFilter }
func (m *MockParserConfig) Logger() *slog.Logger           { return m.Log }
func (m *MockParserConfig) LanguageProcessorFactory() *language.ProcessorFactory {
	return m.LangFactory
}

// NewTestConfig creates a default, permissive config suitable for most parser tests.
func NewTestConfig(sourceDirs []string) parsers.ParserConfig {
	noFilter, _ := filtering.NewDefaultFilter(nil, true)
	langFactory := language.NewProcessorFactory(
		lang_default.NewDefaultProcessor(),
		lang_go.NewGoProcessor(),
		lang_csharp.NewCSharpProcessor(),
		lang_cpp.NewCppProcessor(),
	)

	return &MockParserConfig{
		SrcDirs:     sourceDirs,
		FileFilter:  noFilter,
		Log:         slog.New(slog.NewTextHandler(io.Discard, nil)),
		LangFactory: langFactory,
	}
}
