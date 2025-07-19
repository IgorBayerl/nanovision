package reportconfig

import (
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/filtering"
	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/language"
	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/logging"
	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/settings"
)

var supportedReportTypes = map[string]bool{
	"TextSummary": true,
	"Html":        true,
	"Lcov":        true,
}

// ReportConfiguration struct remains the same.
type ReportConfiguration struct {
	RFiles                        []string
	TDirectory                    string
	SDirectories                  []string
	HDirectory                    string
	RTypes                        []string
	PluginsList                   []string
	AssemblyFilterInstance        filtering.IFilter
	ClassFilterInstance           filtering.IFilter
	FileFilterInstance            filtering.IFilter
	RiskHotspotAssemblyFilterInst filtering.IFilter
	RiskHotspotClassFilterInst    filtering.IFilter
	VLevel                        logging.VerbosityLevel
	CfgTag                        string
	CfgTitle                      string
	CfgLicense                    string
	InvalidPatterns               []string
	VLevelValid                   bool
	App                           *settings.Settings
	logr                          *slog.Logger
	LangFactory                   *language.ProcessorFactory
}

// All accessor methods remain the same.
func (rc *ReportConfiguration) ReportFiles() []string              { return rc.RFiles }
func (rc *ReportConfiguration) TargetDirectory() string            { return rc.TDirectory }
func (rc *ReportConfiguration) SourceDirectories() []string        { return rc.SDirectories }
func (rc *ReportConfiguration) HistoryDirectory() string           { return rc.HDirectory }
func (rc *ReportConfiguration) ReportTypes() []string              { return rc.RTypes }
func (rc *ReportConfiguration) Plugins() []string                  { return rc.PluginsList }
func (rc *ReportConfiguration) AssemblyFilters() filtering.IFilter { return rc.AssemblyFilterInstance }
func (rc *ReportConfiguration) ClassFilters() filtering.IFilter    { return rc.ClassFilterInstance }
func (rc *ReportConfiguration) FileFilters() filtering.IFilter     { return rc.FileFilterInstance }
func (rc *ReportConfiguration) RiskHotspotAssemblyFilters() filtering.IFilter {
	return rc.RiskHotspotAssemblyFilterInst
}
func (rc *ReportConfiguration) RiskHotspotClassFilters() filtering.IFilter {
	return rc.RiskHotspotClassFilterInst
}
func (rc *ReportConfiguration) VerbosityLevel() logging.VerbosityLevel { return rc.VLevel }
func (rc *ReportConfiguration) Tag() string                            { return rc.CfgTag }
func (rc *ReportConfiguration) Title() string                          { return rc.CfgTitle }
func (rc *ReportConfiguration) License() string                        { return rc.CfgLicense }
func (rc *ReportConfiguration) InvalidReportFilePatterns() []string    { return rc.InvalidPatterns }
func (rc *ReportConfiguration) IsVerbosityLevelValid() bool            { return rc.VLevelValid }
func (rc *ReportConfiguration) Settings() *settings.Settings           { return rc.App }

func (rc *ReportConfiguration) Logger() *slog.Logger { return rc.logr }
func (rc *ReportConfiguration) LanguageProcessorFactory() *language.ProcessorFactory {
	return rc.LangFactory
}

// --- Functional Options Pattern Implementation ---

// Option is a function that configures a ReportConfiguration.
type Option func(*ReportConfiguration) error

// NewReportConfiguration is the new, cleaner constructor.
func NewReportConfiguration(
	reportFiles []string,
	targetDir string,
	opts ...Option,
) (*ReportConfiguration, error) {
	defaultAssemblyFilter, _ := filtering.NewDefaultFilter(nil)
	defaultClassFilter, _ := filtering.NewDefaultFilter(nil)
	defaultFileFilter, _ := filtering.NewDefaultFilter(nil, true)

	cfg := &ReportConfiguration{
		RFiles:                        reportFiles,
		TDirectory:                    targetDir,
		RTypes:                        []string{"Html"},
		VLevel:                        logging.Info,
		VLevelValid:                   true,
		CfgTitle:                      "Coverage Report",
		App:                           settings.NewSettings(),
		AssemblyFilterInstance:        defaultAssemblyFilter,
		ClassFilterInstance:           defaultClassFilter,
		FileFilterInstance:            defaultFileFilter,
		RiskHotspotAssemblyFilterInst: defaultAssemblyFilter,
		RiskHotspotClassFilterInst:    defaultClassFilter,
		PluginsList:                   []string{},
		InvalidPatterns:               []string{},
		logr:                          slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
	}

	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

func WithLogger(logger *slog.Logger) Option {
	return func(c *ReportConfiguration) error {
		if logger != nil {
			c.logr = logger
		}
		return nil
	}
}

func WithSourceDirectories(dirs []string) Option {
	return func(c *ReportConfiguration) error {
		c.SDirectories = dirs
		return nil
	}
}

func WithHistoryDirectory(dir string) Option {
	return func(c *ReportConfiguration) error {
		c.HDirectory = dir
		return nil
	}
}

func WithReportTypes(types []string) Option {
	return func(c *ReportConfiguration) error {
		if len(types) == 0 {
			return nil
		}

		var validatedTypes []string
		for _, t := range types {
			trimmedType := strings.TrimSpace(t)
			if trimmedType == "" {
				continue
			}
			if !supportedReportTypes[trimmedType] {
				return fmt.Errorf("unsupported report type: %s", trimmedType)
			}
			validatedTypes = append(validatedTypes, trimmedType)
		}

		if len(validatedTypes) > 0 {
			c.RTypes = validatedTypes
		}

		return nil
	}
}

func WithTag(tag string) Option {
	return func(c *ReportConfiguration) error {
		c.CfgTag = tag
		return nil
	}
}

func WithTitle(title string) Option {
	return func(c *ReportConfiguration) error {
		if title != "" {
			c.CfgTitle = title
		}
		return nil
	}
}

func WithVerbosity(verbosity logging.VerbosityLevel) Option {
	return func(c *ReportConfiguration) error {
		c.VLevel = verbosity
		return nil
	}
}

func WithInvalidPatterns(patterns []string) Option {
	return func(c *ReportConfiguration) error {
		c.InvalidPatterns = patterns
		return nil
	}
}

func WithSettings(s *settings.Settings) Option {
	return func(c *ReportConfiguration) error {
		if s != nil {
			c.App = s
		}
		return nil
	}
}

func WithFilters(
	assemblyFilters []string,
	classFilters []string,
	fileFilters []string,
	rhAssemblyFilters []string,
	rhClassFilters []string,
) Option {
	return func(c *ReportConfiguration) error {
		var err error
		c.AssemblyFilterInstance, err = filtering.NewDefaultFilter(assemblyFilters)
		if err != nil {
			return fmt.Errorf("failed to create assembly filter: %w", err)
		}

		c.ClassFilterInstance, err = filtering.NewDefaultFilter(classFilters)
		if err != nil {
			return fmt.Errorf("failed to create class filter: %w", err)
		}

		c.FileFilterInstance, err = filtering.NewDefaultFilter(fileFilters, true)
		if err != nil {
			return fmt.Errorf("failed to create file filter: %w", err)
		}

		c.RiskHotspotAssemblyFilterInst, err = filtering.NewDefaultFilter(rhAssemblyFilters)
		if err != nil {
			return fmt.Errorf("failed to create risk hotspot assembly filter: %w", err)
		}

		c.RiskHotspotClassFilterInst, err = filtering.NewDefaultFilter(rhClassFilters)
		if err != nil {
			return fmt.Errorf("failed to create risk hotspot class filter: %w", err)
		}

		return nil
	}
}

func WithLanguageProcessorFactory(factory *language.ProcessorFactory) Option {
	return func(c *ReportConfiguration) error {
		if factory != nil {
			c.LangFactory = factory
		}
		return nil
	}
}
