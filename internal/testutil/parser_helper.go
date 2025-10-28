package testutil

import (
	"io"
	"log/slog"

	"github.com/IgorBayerl/nanovision/filtering"
	"github.com/IgorBayerl/nanovision/internal/parsers"
)

// MockParserConfig implements the lean parsers.ParserConfig interface for testing.
type MockParserConfig struct {
	SrcDirs    []string
	FileFilter filtering.IFilter
	Log        *slog.Logger
}

func (m *MockParserConfig) SourceDirectories() []string    { return m.SrcDirs }
func (m *MockParserConfig) FileFilters() filtering.IFilter { return m.FileFilter }
func (m *MockParserConfig) Logger() *slog.Logger           { return m.Log }

// NewTestConfig creates a default, permissive config suitable for most parser tests.
func NewTestConfig(sourceDirs []string) parsers.ParserConfig {
	noFilter, _ := filtering.NewDefaultFilter(nil, true)

	return &MockParserConfig{
		SrcDirs:    sourceDirs,
		FileFilter: noFilter,
		Log:        slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
}
