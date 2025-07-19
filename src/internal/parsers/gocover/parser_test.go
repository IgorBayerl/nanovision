package gocover

import (
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/filtering"
	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/language"
	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/language/defaultformatter"
	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/language/golang"
	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/settings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockFileInfo implements fs.FileInfo for testing.
type MockFileInfo struct {
	name  string
	isDir bool
}

func (m MockFileInfo) Name() string       { return m.name }
func (m MockFileInfo) Size() int64        { return 0 }
func (m MockFileInfo) Mode() fs.FileMode  { return 0 }
func (m MockFileInfo) ModTime() time.Time { return time.Now() }
func (m MockFileInfo) IsDir() bool        { return m.isDir }
func (m MockFileInfo) Sys() interface{}   { return nil }

// MockFileReader for testing without hitting the disk.
type MockFileReader struct {
	Files map[string]string
	Dirs  map[string]bool
}

func NewMockFileReader() *MockFileReader {
	return &MockFileReader{
		Files: make(map[string]string),
		Dirs:  make(map[string]bool),
	}
}

// normalize path to use forward slashes for cross-platform consistency in tests.
func (m *MockFileReader) normalize(path string) string {
	return filepath.ToSlash(path)
}

func (m *MockFileReader) ReadFile(path string) ([]string, error) {
	content, ok := m.Files[m.normalize(path)]
	if !ok {
		return nil, fmt.Errorf("file not found: %s", path)
	}
	return strings.Split(content, "\n"), nil
}

func (m *MockFileReader) CountLines(path string) (int, error) {
	content, ok := m.Files[m.normalize(path)]
	if !ok {
		return 0, fmt.Errorf("file not found: %s", path)
	}
	if content == "" {
		return 0, nil
	}
	count := strings.Count(content, "\n")
	if !strings.HasSuffix(content, "\n") {
		count++
	}
	return count, nil
}

func (m *MockFileReader) Stat(name string) (fs.FileInfo, error) {
	normalizedName := m.normalize(name)
	if _, ok := m.Files[normalizedName]; ok {
		return MockFileInfo{name: filepath.Base(normalizedName), isDir: false}, nil
	}
	if _, ok := m.Dirs[normalizedName]; ok {
		return MockFileInfo{name: filepath.Base(normalizedName), isDir: true}, nil
	}
	return nil, os.ErrNotExist
}

func (m *MockFileReader) AddFile(path, content string) {
	normalizedPath := m.normalize(path)
	m.Files[normalizedPath] = content
	dir := filepath.Dir(normalizedPath)
	for {
		if _, exists := m.Dirs[dir]; exists {
			break
		}
		m.Dirs[dir] = true
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
}

// mockParserConfig for providing test configuration.
type mockParserConfig struct {
	srcDirs        []string
	assemblyFilter filtering.IFilter
	classFilter    filtering.IFilter
	fileFilter     filtering.IFilter
	settings       *settings.Settings
	logger         *slog.Logger
	langFactory    *language.ProcessorFactory
}

func (m *mockParserConfig) SourceDirectories() []string        { return m.srcDirs }
func (m *mockParserConfig) AssemblyFilters() filtering.IFilter { return m.assemblyFilter }
func (m *mockParserConfig) ClassFilters() filtering.IFilter    { return m.classFilter }
func (m *mockParserConfig) FileFilters() filtering.IFilter     { return m.fileFilter }
func (m *mockParserConfig) Settings() *settings.Settings       { return m.settings }
func (m *mockParserConfig) Logger() *slog.Logger               { return m.logger }
func (m *mockParserConfig) LanguageProcessorFactory() *language.ProcessorFactory {
	return m.langFactory
}

func newTestConfig() *mockParserConfig {
	noFilter, _ := filtering.NewDefaultFilter(nil)

	// Create a language factory with the processors needed for this test.
	langFactory := language.NewProcessorFactory(
		defaultformatter.NewDefaultProcessor(),
		golang.NewGoProcessor(),
	)

	return &mockParserConfig{
		srcDirs:        []string{"/project/src"},
		assemblyFilter: noFilter,
		classFilter:    noFilter,
		fileFilter:     noFilter,
		settings:       settings.NewSettings(),
		logger:         slog.New(slog.NewTextHandler(io.Discard, nil)),
		langFactory:    langFactory,
	}
}

func TestGoCoverParser_Parse_Success(t *testing.T) {
	coverProfileContent := `mode: set
calculator/calculator.go:4.2,4.13 1 1
calculator/calculator.go:8.2,8.13 1 1
calculator/calculator.go:12.2,12.23 1 1
calculator/calculator.go:13.3,13.12 1 0
calculator/calculator.go:15.2,15.13 1 1
calculator/calculator.go:17.2,17.13 1 0`

	calculatorGoContent := `package calculator

func Add(a, b int) int {
	return a + b // Line 4
}

func Subtract(a, b int) int {
	return a - b // Line 8
}

func Multiply(a, b int) int {
	if a == 0 || b == 0 { // Line 12
		return 0 // Line 13
	}
	return a * b // Line 15
}
func Divide(a, b int) int {
	return a / b // Line 17
}`
	goModContent := `module example.com/calculator`

	reportFile, err := os.CreateTemp(t.TempDir(), "cover_*.out")
	require.NoError(t, err)
	_, err = reportFile.WriteString(coverProfileContent)
	require.NoError(t, err)
	reportFile.Close()

	mockFileReader := NewMockFileReader()
	mockFileReader.AddFile("/project/src/calculator/calculator.go", calculatorGoContent)
	mockFileReader.AddFile("/project/src/go.mod", goModContent)

	p := NewGoCoverParser(mockFileReader)
	config := newTestConfig()
	result, err := p.Parse(reportFile.Name(), config)
	require.NoError(t, err)
	require.NotNil(t, result)

	assemblies := result.Assemblies
	require.Len(t, assemblies, 1)
	assembly := assemblies[0]
	assert.Equal(t, "example.com/calculator", assembly.Name)

	require.Len(t, assembly.Classes, 1)
	class := assembly.Classes[0]

	require.Len(t, class.Methods, 4)
	methodCoverage := make(map[string]float64)
	for _, m := range class.Methods {
		methodCoverage[m.Name] = m.LineRate
	}

	assert.InDelta(t, 1.0, methodCoverage["Add"], 0.001)
	assert.InDelta(t, 1.0, methodCoverage["Subtract"], 0.001)
	assert.InDelta(t, 0.666, methodCoverage["Multiply"], 0.001, "2 of 3 statements covered (if + final return)")
	assert.InDelta(t, 0.0, methodCoverage["Divide"], 0.001)
}

func TestProcessingOrchestrator_findModuleNameFromGoMod(t *testing.T) {
	mockFileReader := NewMockFileReader()
	mockFileReader.AddFile("/project/src/go.mod", "module github.com/example/myproject\n")
	mockFileReader.AddFile("/project/src/pkg/math/add.go", "package math")
	mockFileReader.AddFile("/other/project/main.go", "package main")

	config := newTestConfig()
	orchestrator := newProcessingOrchestrator(mockFileReader, config, slog.Default())

	result, err := orchestrator.findModuleNameFromGoMod("/project/src/pkg/math/add.go")
	assert.NoError(t, err)
	assert.Equal(t, "github.com/example/myproject", result)

	_, err = orchestrator.findModuleNameFromGoMod("/other/project/main.go")
	assert.Error(t, err)
}

func TestDefaultFileReader_Integration(t *testing.T) {
	content := "line 1\nline 2\nline 3"
	file, err := os.CreateTemp("", "test_*.go")
	require.NoError(t, err)
	defer os.Remove(file.Name())

	_, err = file.WriteString(content)
	require.NoError(t, err)
	file.Close()

	reader := &DefaultFileReader{}
	lines, err := reader.ReadFile(file.Name())
	assert.NoError(t, err)
	assert.Len(t, lines, 3)

	count, err := reader.CountLines(file.Name())
	assert.NoError(t, err)
	assert.Equal(t, 3, count)
}
