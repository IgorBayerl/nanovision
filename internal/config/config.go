package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/IgorBayerl/nanovision/filtering"
	"github.com/IgorBayerl/nanovision/logging"
	"gopkg.in/yaml.v3"
)

type ReportInputPair struct {
	ReportPattern string
	SourceDir     string
}

type RawConfigInput struct {
	ReportPatterns string
	SourceDirs     string
	ReportTypes    string
	FileFilters    string
	OutputDir      string
	Tag            string
	Title          string
	LogFile        string
	LogFormat      string
	Verbosity      string
	Verbose        bool
}

type AppConfig struct {
	ReportPatterns []string `yaml:"reports"`
	SourceDirs     []string `yaml:"source_dirs"`
	ReportTypes    []string `yaml:"report_types"`
	FileFilters    []string `yaml:"file_filters"`
	OutputDir      string   `yaml:"output_dir"`
	Tag            string   `yaml:"tag"`
	Title          string   `yaml:"title"`
	LogFile        string   `yaml:"log_file"`
	LogFormat      string   `yaml:"log_format"`
	Verbosity      string   `yaml:"verbosity"`
	IgnoreFiles    []string `yaml:"ignore_files"`
	ProjectRoot    string   `yaml:"-"`

	FileFilterInstance filtering.IFilter
	VerbosityLevel     logging.VerbosityLevel
	InputPairs         []ReportInputPair
}

// resolveInputPairs matches slices of report patterns and source directories into structured pairs.
func resolveInputPairs(patterns []string, dirs []string) []ReportInputPair {
	var pairs []ReportInputPair
	for i := range patterns {
		trimmedPattern := strings.TrimSpace(patterns[i])
		var trimmedSourceDir string
		if i < len(dirs) {
			trimmedSourceDir = strings.TrimSpace(dirs[i])
		}

		if trimmedPattern != "" && trimmedSourceDir != "" {
			pairs = append(pairs, ReportInputPair{
				ReportPattern: trimmedPattern,
				SourceDir:     trimmedSourceDir,
			})
		}
	}
	return pairs
}

// Load loads the configuration from defaults, a YAML file, and CLI flags.
func Load(configPath string, cliInput RawConfigInput) (*AppConfig, error) {
	cfg := GetDefaultConfig()

	if configPath == "" {
		if _, err := os.Stat("nanovision.yaml"); err == nil {
			configPath = "nanovision.yaml"
		}
	}

	if configPath != "" {
		yamlFile, err := os.ReadFile(configPath)
		if err != nil {
			if !os.IsNotExist(err) || (configPath != "nanovision.yaml" && configPath != "") {
				return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
			}
		} else {
			if err := yaml.Unmarshal(yamlFile, cfg); err != nil {
				return nil, fmt.Errorf("failed to parse YAML config %s: %w", configPath, err)
			}
		}
	}

	cfg.mergeCliOverrides(cliInput)

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	if err := cfg.computeDerivedFields(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// GetDefaultConfig returns a new AppConfig with hard-coded default values.
func GetDefaultConfig() *AppConfig {
	return &AppConfig{
		OutputDir:   "coverage-report",
		ReportTypes: []string{"TextSummary", "Html"},
		Title:       "Coverage Report",
		LogFormat:   "text",
		Verbosity:   "Info",
	}
}

// mergeCliOverrides updates the config with values from the CLI.
func (c *AppConfig) mergeCliOverrides(cli RawConfigInput) {
	if cli.ReportPatterns != "" {
		c.ReportPatterns = strings.Split(cli.ReportPatterns, ";")
	}
	if cli.SourceDirs != "" {
		c.SourceDirs = strings.Split(cli.SourceDirs, ";")
	}
	if cli.ReportTypes != "TextSummary,Html" {
		c.ReportTypes = strings.Split(cli.ReportTypes, ",")
	}
	if cli.FileFilters != "" {
		c.FileFilters = strings.Split(cli.FileFilters, ";")
	}
	if cli.OutputDir != "coverage-report" {
		c.OutputDir = cli.OutputDir
	}
	if cli.Tag != "" {
		c.Tag = cli.Tag
	}
	if cli.Title != "" {
		c.Title = cli.Title
	}
	if cli.LogFile != "" {
		c.LogFile = cli.LogFile
	}
	if cli.LogFormat != "text" {
		c.LogFormat = cli.LogFormat
	}
	if cli.Verbosity != "Info" {
		c.Verbosity = cli.Verbosity
	}
	if cli.Verbose {
		c.Verbosity = "Verbose"
	}
}

// validate checks the final configuration for logical errors.
func (c *AppConfig) validate() error {
	if len(c.ReportPatterns) == 0 {
		return errors.New("configuration error: at least one report pattern must be specified")
	}
	if len(c.SourceDirs) == 0 {
		return errors.New("configuration error: at least one source directory must be specified")
	}
	if len(c.ReportPatterns) != len(c.SourceDirs) {
		return fmt.Errorf(
			"configuration error: mismatch between number of report patterns (%d) and source directories (%d)",
			len(c.ReportPatterns),
			len(c.SourceDirs),
		)
	}
	if _, err := logging.ParseVerbosity(c.Verbosity); err != nil {
		return fmt.Errorf("invalid verbosity level '%s'", c.Verbosity)
	}
	return nil
}

// computeDerivedFields processes raw config values into usable internal fields.
func (c *AppConfig) computeDerivedFields() error {
	allFilters := c.FileFilters
	for _, pattern := range c.IgnoreFiles {
		allFilters = append(allFilters, "-"+pattern)
	}

	filter, err := filtering.NewDefaultFilter(allFilters, true)
	if err != nil {
		return fmt.Errorf("failed to initialize file filter: %w", err)
	}
	c.FileFilterInstance = filter

	c.VerbosityLevel, _ = logging.ParseVerbosity(c.Verbosity)

	c.InputPairs = resolveInputPairs(c.ReportPatterns, c.SourceDirs)

	return nil
}
