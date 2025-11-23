package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/IgorBayerl/nanovision/internal/filtering"
	"github.com/IgorBayerl/nanovision/internal/logging"
	"gopkg.in/yaml.v3"
)

// -------- Status bands parsing ----------
// Band is an inclusive range [Min, Max].
type Band struct {
	Min float64 `yaml:"min"`
	Max float64 `yaml:"max"`
}

// MetricKey is the canonical snake_case key used in config & model.
type MetricKey string

const (
	LineCoverage         MetricKey = "line_coverage"
	BranchCoverage       MetricKey = "branch_coverage"
	MethodsCovered       MetricKey = "methods_covered"
	MethodsFullyCovered  MetricKey = "methods_fully_covered"
	PatchLineCoverage    MetricKey = "patch_line_coverage"
	PatchMethodsCovered  MetricKey = "patch_methods_covered"
	MethodBranchCoverage MetricKey = "method_branch_coverage"
)

// StatusBands supports either:
// 1) map form:    { "line_coverage": "60..75", ... }
// 2) list form:   [ { "line_coverage": "60..75" }, ... ]
type StatusBands map[MetricKey]Band

func (sb *StatusBands) UnmarshalYAML(unmarshal func(any) error) error {
	// Try map[string]string first
	var m1 map[string]string
	if err := unmarshal(&m1); err == nil && len(m1) > 0 {
		bands := make(StatusBands, len(m1))
		for k, v := range m1 {
			b, err := parseBandString(v)
			if err != nil {
				return fmt.Errorf("status_bands.%s: %w", k, err)
			}
			bands[MetricKey(k)] = b
		}
		*sb = bands
		return nil
	}
	// Try []map[string]string (list form)
	var m2 []map[string]string
	if err := unmarshal(&m2); err == nil && len(m2) > 0 {
		bands := make(StatusBands)
		for _, item := range m2 {
			for k, v := range item {
				b, err := parseBandString(v)
				if err != nil {
					return fmt.Errorf("status_bands.%s: %w", k, err)
				}
				bands[MetricKey(k)] = b
			}
		}
		*sb = bands
		return nil
	}
	// Empty or null is allowed (no bands configured)
	*sb = make(StatusBands)
	return nil
}

func parseBandString(s string) (Band, error) {
	parts := strings.Split(s, "..")
	if len(parts) != 2 {
		return Band{}, fmt.Errorf("invalid band %q (expected \"min..max\")", s)
	}
	min, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	max, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err1 != nil || err2 != nil {
		return Band{}, fmt.Errorf("invalid numbers in %q", s)
	}
	if min < 0 || max > 100 || min > max {
		return Band{}, fmt.Errorf("out of range in %q (require 0 <= min <= max <= 100)", s)
	}
	return Band{Min: min, Max: max}, nil
}

type ReportInputPair struct {
	ReportPattern string
	SourceDir     string
}

// DiffPathMap defines a single path rewrite rule for diffs.
type DiffPathMap struct {
	From string `yaml:"from"`
	To   string `yaml:"to"`
}

// DiffConfig holds all settings related to diff processing.
type DiffConfig struct {
	File         string        `yaml:"file"`
	Strip        string        `yaml:"strip"`
	RootOverride string        `yaml:"root_override"`
	PathMaps     []DiffPathMap `yaml:"path_maps"`
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
	DiffFile       string
	DiffStrip      string
	StatusBands    []string
}

type AppConfig struct {
	ReportPatterns []string    `yaml:"reports"`
	SourceDirs     []string    `yaml:"source_dirs"`
	ReportTypes    []string    `yaml:"report_types"`
	FileFilters    []string    `yaml:"file_filters"`
	OutputDir      string      `yaml:"output_dir"`
	Tag            string      `yaml:"tag"`
	Title          string      `yaml:"title"`
	LogFile        string      `yaml:"log_file"`
	LogFormat      string      `yaml:"log_format"`
	Verbosity      string      `yaml:"verbosity"`
	IgnoreFiles    []string    `yaml:"ignore_files"`
	ProjectRoot    string      `yaml:"-"`
	StatusBands    StatusBands `yaml:"status_bands"`
	Diff           DiffConfig  `yaml:"diff"`

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
		Diff: DiffConfig{
			Strip: "auto",
		},
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
	if cli.DiffFile != "" {
		c.Diff.File = cli.DiffFile
	}
	if cli.DiffStrip != "" {
		c.Diff.Strip = cli.DiffStrip
	}
	if len(cli.StatusBands) > 0 {
		if c.StatusBands == nil {
			c.StatusBands = make(StatusBands)
		}
		for _, p := range cli.StatusBands {
			kv := strings.Split(p, "=")
			if len(kv) == 2 {
				k := strings.TrimSpace(kv[0])
				v := strings.TrimSpace(kv[1])
				if band, err := parseBandString(v); err == nil {
					c.StatusBands[MetricKey(k)] = band
				}
			}
		}
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

	if c.Diff.File != "" {
		stripVal := strings.ToLower(c.Diff.Strip)
		if stripVal != "auto" {
			if n, err := strconv.Atoi(stripVal); err != nil || n < 0 || n > 6 {
				// The calling code should log a warning. We revert to the safe default.
				c.Diff.Strip = "auto"
			}
		}
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
