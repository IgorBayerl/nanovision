package config

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/IgorBayerl/nanovision/internal/filtering"
	"github.com/IgorBayerl/nanovision/internal/logging"
	"go.yaml.in/yaml/v3"
)

// Every field with a yaml key needs a doc tag: Schema() turns the tags into
// the config reference of the CLI and the website, and a test fails when one
// is missing. A values tag lists the allowed values. A flag tag names the
// command-line flag that sets the same key; a test checks that the flag
// exists. Defaults are not tags: Schema() reads them from GetDefaultConfig().

// Band is an inclusive range [Min, Max], written "min..max" in the config.
type Band struct {
	Min float64
	Max float64
}

func (b *Band) UnmarshalYAML(n *yaml.Node) error {
	if n.Kind != yaml.ScalarNode {
		return fmt.Errorf("line %d: a warning range is written \"min..max\", for example \"60..75\"", n.Line)
	}
	parsed, err := parseBandString(n.Value)
	if err != nil {
		return fmt.Errorf("line %d: %w", n.Line, err)
	}
	*b = parsed
	return nil
}

func (b Band) String() string {
	return strconv.FormatFloat(b.Min, 'f', -1, 64) + ".." + strconv.FormatFloat(b.Max, 'f', -1, 64)
}

func (b Band) MarshalYAML() (any, error) { return b.String(), nil }

func parseBandString(s string) (Band, error) {
	parts := strings.Split(s, "..")
	if len(parts) != 2 {
		return Band{}, fmt.Errorf("invalid range %q (expected \"min..max\")", s)
	}
	min, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	max, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err1 != nil || err2 != nil {
		return Band{}, fmt.Errorf("invalid numbers in %q", s)
	}
	if min < 0 || min > max {
		return Band{}, fmt.Errorf("out of range in %q (require 0 <= min <= max)", s)
	}
	return Band{Min: min, Max: max}, nil
}

// StatusBands is the warning range of each metric that has one.
type StatusBands map[MetricKey]Band

// FolderBand is the warning ranges a folder sets for itself.
type FolderBand struct {
	Path  string // relative to the project root, "/" separated
	Bands StatusBands
	// the config file that sets them, for messages
	Source string
}

// ReportInput is one coverage report and the folder its paths are relative to.
type ReportInput struct {
	Path   string `yaml:"path" doc:"Report file or glob pattern."`
	Source string `yaml:"source" doc:"Folder of the source code for this report. Default: the project root."`
	Name   string `yaml:"name" doc:"Label in the report selector. Default: the file name."`
}

// UnmarshalYAML also accepts a plain string, which is the path.
func (r *ReportInput) UnmarshalYAML(n *yaml.Node) error {
	if n.Kind == yaml.ScalarNode {
		r.Path = n.Value
		return nil
	}
	if n.Kind != yaml.MappingNode {
		return fmt.Errorf("line %d: a report is a path, or a mapping with path, source and name", n.Line)
	}
	for i := 0; i+1 < len(n.Content); i += 2 {
		key, value := n.Content[i], n.Content[i+1]
		switch key.Value {
		case "path":
			r.Path = value.Value
		case "source":
			r.Source = value.Value
		case "name":
			r.Name = value.Value
		default:
			return fmt.Errorf("line %d: unknown key %q in a report (expected path, source or name)", key.Line, key.Value)
		}
	}
	return nil
}

// MetricSetting enables one metric and holds its settings.
type MetricSetting struct {
	Name    string `yaml:"name" doc:"Metric name."`
	Warning *Band  `yaml:"warning,omitempty" doc:"Warning range \"min..max\". Below is danger, above is safe. Reversed for lower-is-better metrics."`
}

type MetricsConfig struct {
	Files   []MetricSetting `yaml:"files,omitempty" doc:"File and folder metrics, in column order. Empty shows all."`
	Methods []MetricSetting `yaml:"methods,omitempty" doc:"Method metrics, in column order. Empty shows all."`
}

// ScopedConfig is the part of the configuration a folder can set for itself,
// in a nanovision.yaml inside the folder. Every other key is valid only in the
// root file.
type ScopedConfig struct {
	Reports     []ReportInput `yaml:"reports,omitempty" doc:"Coverage reports to read."`
	IgnoreFiles []string      `yaml:"ignore_files,omitempty" doc:"Source files to leave out. Glob patterns."`
	Metrics     MetricsConfig `yaml:"metrics,omitempty" doc:"Metrics to show, with their warning ranges."`
}

// FolderConfig is the nanovision.yaml of one sub-folder. It applies to the
// folder and everything below it.
type FolderConfig struct {
	Path string // the folder, relative to the project root, "/" separated
	ScopedConfig
	File string // the config file
}

// DiffPathMap defines a single path rewrite rule for diffs.
type DiffPathMap struct {
	From string `yaml:"from" doc:"Path prefix in the diff."`
	To   string `yaml:"to" doc:"Path prefix in the project."`
}

// DiffConfig holds all settings related to diff processing.
type DiffConfig struct {
	File         string        `yaml:"file" flag:"diff" doc:"Unified diff file. Empty: the VCS makes the diff."`
	Strip        string        `yaml:"strip" flag:"diff-strip" doc:"Leading path parts to remove: auto, or 0 to 6."`
	OnlyChanged  bool          `yaml:"only_changed" flag:"only-changed" doc:"The HTML report holds only the changed files."`
	RootOverride string        `yaml:"root_override" doc:"Folder the diff paths are relative to."`
	PathMaps     []DiffPathMap `yaml:"path_maps" doc:"Rewrite rules for diff paths."`
}

// ReviewGate defines pass/fail thresholds evaluated against the changelist.
// Nil pointers mean "check disabled".
type ReviewGate struct {
	PatchStatementCoverage     *float64 `yaml:"patch_statement_coverage" doc:"Minimum coverage of the changed statements, 0 to 100."`
	MaxChangedMethodComplexity *int     `yaml:"max_changed_method_complexity" doc:"Maximum complexity of a changed method."`
}

type ReviewConfig struct {
	Gate   ReviewGate `yaml:"gate" doc:"Limits for the changed code."`
	FailOn string     `yaml:"fail_on" flag:"fail-on" doc:"Fail the build on errors in changed code, or also on warnings." values:"never,error,warning"`
}

// HistoryConfig turns on the run store. An empty Store keeps history off.
type HistoryConfig struct {
	Store       string `yaml:"store" flag:"store" doc:"Server URL or local folder. Empty turns history off."`
	Project     string `yaml:"project" flag:"project" doc:"Project name. Default: the project folder name."`
	Profile     string `yaml:"profile" flag:"profile" doc:"Name of the test suite. Keeps its runs apart."`
	MaxDistance int    `yaml:"max_distance" doc:"Revisions to search back for a base run."`
	KeepLocal   int    `yaml:"keep_local" doc:"Local runs to keep for each stream."`
}

// IsServer reports whether the store is a team server rather than a local folder.
func (h HistoryConfig) IsServer() bool {
	s := strings.ToLower(h.Store)
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}

// VCSConfig selects the version control adapter. An empty Type turns it off.
type VCSConfig struct {
	Type       string `yaml:"type" flag:"vcs" doc:"Version control system. Empty turns it off." values:"git,perforce"`
	BaseBranch string `yaml:"base_branch" doc:"Git: branch to compare with. Default: the remote default branch."`
}

// RunFlags describe one run. They only come from flags, because they change
// with every run.
type RunFlags struct {
	Kind         string // submit, review or local (default)
	Revision     string // the revision this run measured
	BaseRevision string // the revision to compare with
	Stream       string // the stream or branch, when no VCS adapter knows it
	ReviewID     string // the shelved changelist or pull request of a review run
	Author       string
	CIURL        string
}

type RawConfigInput struct {
	// each "path[,source=DIR][,name=LABEL]"; several paths can share one
	// entry, separated by ";"
	Reports    []string
	SourceDirs string
	// each "key=value" for a config key, e.g. "history.keep_local=5"
	Set         []string
	ReportTypes string
	FileFilters string
	OutputDir   string
	Tag         string
	Title       string
	LogFile     string
	LogFormat   string
	Verbosity   string
	Verbose     bool
	DiffFile    string
	DiffStrip   string
	OnlyChanged bool
	StatusBands []string
	// each "name" or "name=min..max"
	FileMetrics    []string
	MethodMetrics  []string
	IgnoreCache    bool
	DefaultFilters string
	FailOn         string

	Store        string
	Project      string
	Profile      string
	VCS          string
	RunKind      string
	Revision     string
	BaseRevision string
	Stream       string
	ReviewID     string
	Author       string
	CIURL        string
}

type AppConfig struct {
	Title          string   `yaml:"title" flag:"title" doc:"Title of the HTML report."`
	OutputDir      string   `yaml:"output_dir" flag:"output" doc:"Folder for the output files."`
	ReportTypes    []string `yaml:"report_types" flag:"reporttypes" doc:"Output formats to write."`
	ScopedConfig   `yaml:",inline"`
	FileFilters    []string      `yaml:"file_filters" flag:"filefilters" doc:"Source files to include (+pattern) or exclude (-pattern)."`
	NestedConfigs  bool          `yaml:"nested_configs" doc:"Read the nanovision.yaml files in sub-folders."`
	Diff           DiffConfig    `yaml:"diff" doc:"The change to measure."`
	Review         ReviewConfig  `yaml:"review" doc:"Rules for the changed code."`
	VCS            VCSConfig     `yaml:"vcs" doc:"Version control."`
	History        HistoryConfig `yaml:"history" doc:"Run store. It gives the comparison with a base run."`
	DefaultFilters string        `yaml:"default_filters" flag:"default-filters" doc:"Filters applied when the report opens. Example: risk=danger."`
	Tag            string        `yaml:"tag" flag:"tag" doc:"Tag of the run. Example: a build number."`
	Verbosity      string        `yaml:"verbosity" flag:"verbosity" doc:"Log level." values:"Off,Error,Warning,Info,Verbose"`
	LogFile        string        `yaml:"log_file" flag:"logfile" doc:"Also write the log to this file."`
	LogFormat      string        `yaml:"log_format" flag:"logformat" doc:"Log format." values:"text,json"`
	IgnoreCache    bool          `yaml:"ignore_cache" flag:"ignore-cache" doc:"Analyze every file again. Do not use the cache."`

	// everything below is derived from the fields above

	// the nanovision.yaml files found in sub-folders, each for its folder
	Folders     []FolderConfig `yaml:"-"`
	ProjectRoot string         `yaml:"-"`
	ConfigFile  string         `yaml:"-"` // the root config file; empty when only flags were used
	Run         RunFlags       `yaml:"-"`

	FileMetrics         []MetricKey        `yaml:"-"` // in column order
	MethodMetrics       []MetricKey        `yaml:"-"`
	ActiveFileMetrics   map[MetricKey]bool `yaml:"-"`
	ActiveMethodMetrics map[MetricKey]bool `yaml:"-"`
	// the warning ranges of the root config; BandsFor adds the folder configs
	StatusBands StatusBands  `yaml:"-"`
	FolderBands []FolderBand `yaml:"-"` // shortest path first

	FileFilterInstance filtering.IFilter      `yaml:"-"`
	VerbosityLevel     logging.VerbosityLevel `yaml:"-"`
	InputPairs         []ReportInputPair      `yaml:"-"`

	cliThresholds []string
}

type ReportInputPair struct {
	ReportPattern string
	SourceDir     string
	Name          string
}

// a profile is part of URLs and store keys
var profileNameRE = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)

// keys of older versions, with what replaces them
var removedKeys = map[string]string{
	"overrides":      "put a nanovision.yaml inside the folder; it applies to that folder and everything below it",
	"source_dirs":    "put the source folder on each report:\n  reports:\n    - path: coverage.out\n      source: .",
	"status_bands":   "put the range on the metric:\n  metrics:\n    files:\n      - name: statement_coverage\n        warning: \"60..75\"",
	"file_metrics":   "list them under metrics.files:\n  metrics:\n    files:\n      - name: statement_coverage",
	"method_metrics": "list them under metrics.methods:\n  metrics:\n    methods:\n      - name: statement_coverage",
}

// Load loads the configuration from defaults, a YAML file, and CLI flags.
func Load(configPath string, cliInput RawConfigInput) (*AppConfig, error) {
	cfg := GetDefaultConfig()

	if configPath == "" {
		found, err := FindFile(".")
		if err != nil {
			return nil, err
		}
		configPath = found
	}

	if configPath != "" {
		yamlFile, err := os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
		}
		if err := decodeRoot(yamlFile, cfg); err != nil {
			return nil, fmt.Errorf("%s: %w", configPath, err)
		}
		cfg.ConfigFile = configPath
	}

	var err error
	if cfg.ProjectRoot, err = projectRoot(configPath); err != nil {
		return nil, err
	}

	if err := cfg.mergeCliOverrides(cliInput); err != nil {
		return nil, err
	}

	if cfg.ConfigFile != "" && cfg.NestedConfigs {
		nested, err := findNestedConfigs(cfg)
		if err != nil {
			return nil, err
		}
		cfg.Folders = append(cfg.Folders, nested...)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	if err := cfg.computeDerivedFields(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// projectRoot is the folder of the config file, or the working directory.
func projectRoot(configPath string) (string, error) {
	if configPath != "" {
		abs, err := filepath.Abs(configPath)
		if err != nil {
			return "", fmt.Errorf("could not determine absolute path for config file: %w", err)
		}
		return filepath.Dir(abs), nil
	}
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("could not get current working directory: %w", err)
	}
	return wd, nil
}

// decodeRoot reads the root config file. Unknown keys are errors, so a typo
// does not silently turn a setting off.
func decodeRoot(data []byte, cfg *AppConfig) error {
	keys, err := topLevelKeys(data)
	if err != nil {
		return err
	}
	for _, k := range keys {
		if hint, ok := removedKeys[k]; ok {
			return fmt.Errorf("%q is no longer a config key; %s", k, hint)
		}
	}
	return decodeStrict(data, cfg)
}

func decodeStrict(data []byte, into any) error {
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true)
	if err := dec.Decode(into); err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("failed to parse YAML config: %w", cleanYAMLError(err))
	}
	return nil
}

func topLevelKeys(data []byte) ([]string, error) {
	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("failed to parse YAML config: %w", err)
	}
	if len(doc.Content) == 0 || doc.Content[0].Kind != yaml.MappingNode {
		return nil, nil
	}
	var keys []string
	for i := 0; i < len(doc.Content[0].Content); i += 2 {
		keys = append(keys, doc.Content[0].Content[i].Value)
	}
	return keys, nil
}

// cleanYAMLError drops the Go type names from "field x not found in type config.AppConfig".
func cleanYAMLError(err error) error {
	msg := regexp.MustCompile(`field (\S+) not found in type \S+`).ReplaceAllString(err.Error(), `unknown key "$1"`)
	return errors.New(msg)
}

// GetDefaultConfig returns a new AppConfig with hard-coded default values.
func GetDefaultConfig() *AppConfig {
	return &AppConfig{
		OutputDir:     "coverage-report",
		ReportTypes:   []string{"TextSummary", "Html"},
		Title:         "Coverage Report",
		LogFormat:     "text",
		Verbosity:     "Info",
		NestedConfigs: true,
		Diff: DiffConfig{
			Strip: "auto",
		},
		Review: ReviewConfig{
			FailOn: "never",
		},
		History: HistoryConfig{
			Profile:     "default",
			MaxDistance: 50,
			KeepLocal:   20,
		},
	}
}

// mergeCliOverrides updates the config with values from the CLI.
func (c *AppConfig) mergeCliOverrides(cli RawConfigInput) error {
	if err := c.applySet(cli.Set); err != nil {
		return err
	}
	if err := c.mergeCliReports(cli.Reports, cli.SourceDirs); err != nil {
		return err
	}
	if cli.ReportTypes != "TextSummary,Html" && cli.ReportTypes != "" {
		c.ReportTypes = strings.Split(cli.ReportTypes, ",")
	}
	if cli.FileFilters != "" {
		c.FileFilters = strings.Split(cli.FileFilters, ";")
	}
	if cli.OutputDir != "coverage-report" && cli.OutputDir != "" {
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
	if cli.LogFormat != "text" && cli.LogFormat != "" {
		c.LogFormat = cli.LogFormat
	}
	if cli.Verbosity != "Info" && cli.Verbosity != "" {
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
	if cli.OnlyChanged {
		c.Diff.OnlyChanged = true
	}
	if cli.IgnoreCache {
		c.IgnoreCache = true
	}
	if cli.DefaultFilters != "" {
		c.DefaultFilters = cli.DefaultFilters
	}
	if cli.FailOn != "" {
		c.Review.FailOn = cli.FailOn
	}
	if cli.Store != "" {
		c.History.Store = cli.Store
	}
	if cli.Project != "" {
		c.History.Project = cli.Project
	}
	if cli.Profile != "" {
		c.History.Profile = cli.Profile
	}
	if cli.VCS != "" {
		c.VCS.Type = cli.VCS
	}
	c.Run = RunFlags{
		Kind:         strings.ToLower(strings.TrimSpace(cli.RunKind)),
		Revision:     strings.TrimSpace(cli.Revision),
		BaseRevision: strings.TrimSpace(cli.BaseRevision),
		Stream:       strings.TrimSpace(cli.Stream),
		ReviewID:     strings.TrimSpace(cli.ReviewID),
		Author:       strings.TrimSpace(cli.Author),
		CIURL:        strings.TrimSpace(cli.CIURL),
	}
	var err error
	if c.Metrics.Files, err = settingsFromFlags(cli.FileMetrics, c.Metrics.Files); err != nil {
		return fmt.Errorf("-file-metric: %w", err)
	}
	if c.Metrics.Methods, err = settingsFromFlags(cli.MethodMetrics, c.Metrics.Methods); err != nil {
		return fmt.Errorf("-method-metric: %w", err)
	}
	c.cliThresholds = cli.StatusBands
	return nil
}

// mergeCliReports applies -report, one report for each use of the flag:
// "path[,source=DIR][,name=LABEL]". The older forms still work: several paths
// in one value separated by ";", and -sourcedirs with one folder for each
// report or one folder for all.
func (c *AppConfig) mergeCliReports(entries []string, dirs string) error {
	if len(entries) > 0 {
		c.Reports = nil
	}
	for _, entry := range entries {
		for _, one := range splitList(entry, ";") {
			parts := splitList(one, ",")
			r := ReportInput{Path: parts[0]}
			for _, part := range parts[1:] {
				key, value, _ := strings.Cut(part, "=")
				switch strings.TrimSpace(key) {
				case "source":
					r.Source = strings.TrimSpace(value)
				case "name":
					r.Name = strings.TrimSpace(value)
				default:
					return fmt.Errorf("-report %q: unknown part %q (expected path[,source=DIR][,name=LABEL])", one, part)
				}
			}
			c.Reports = append(c.Reports, r)
		}
	}

	sources := splitList(dirs, ";")
	switch {
	case len(sources) == 0:
	case len(sources) == 1:
		for i := range c.Reports {
			c.Reports[i].Source = sources[0]
		}
	case len(sources) == len(c.Reports):
		for i := range c.Reports {
			c.Reports[i].Source = sources[i]
		}
	default:
		return fmt.Errorf("configuration error: mismatch between number of report patterns (%d) and source directories (%d)",
			len(c.Reports), len(sources))
	}
	return nil
}

func splitList(s, sep string) []string {
	var out []string
	for _, p := range strings.Split(s, sep) {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

// settingsFromFlags is the metric list of -file-metric or -method-metric:
// one "name" or "name=min..max" for each use. A metric without a range keeps
// the one the config file gave it. Without the flag the list does not change.
func settingsFromFlags(entries []string, old []MetricSetting) ([]MetricSetting, error) {
	if len(entries) == 0 {
		return old, nil
	}
	var out []MetricSetting
	for _, entry := range entries {
		for _, one := range splitList(entry, ",") {
			name, value, hasRange := strings.Cut(one, "=")
			setting := MetricSetting{Name: strings.TrimSpace(name)}
			if hasRange {
				band, err := parseBandString(strings.TrimSpace(value))
				if err != nil {
					return nil, fmt.Errorf("%s: %w", setting.Name, err)
				}
				setting.Warning = &band
			} else {
				for _, o := range old {
					if o.Name == setting.Name {
						setting.Warning = o.Warning
					}
				}
			}
			out = append(out, setting)
		}
	}
	return out, nil
}

// applySet applies -set key=value, which sets any plain config key by its
// name in the config file, e.g. -set history.keep_local=5. A key that holds a
// list takes one -set for each item.
func (c *AppConfig) applySet(entries []string) error {
	if len(entries) == 0 {
		return nil
	}
	types := make(map[string]string)
	for _, f := range Flatten(Schema().Fields) {
		types[f.Key] = f.Type
	}

	root := &yaml.Node{Kind: yaml.MappingNode}
	// child returns the value node of key in a mapping, adding it when new
	child := func(m *yaml.Node, key string, kind yaml.Kind) *yaml.Node {
		for i := 0; i+1 < len(m.Content); i += 2 {
			if m.Content[i].Value == key {
				return m.Content[i+1]
			}
		}
		value := &yaml.Node{Kind: kind}
		m.Content = append(m.Content, &yaml.Node{Kind: yaml.ScalarNode, Value: key}, value)
		return value
	}

	for _, entry := range entries {
		key, value, ok := strings.Cut(entry, "=")
		key = strings.TrimSpace(key)
		if !ok {
			return fmt.Errorf("invalid -set %q (expected key=value)", entry)
		}
		kind, known := types[key]
		switch {
		case !known || strings.Contains(key, "[]"):
			return fmt.Errorf("-set %s: unknown config key; 'nanovision config docs' lists the keys", key)
		case kind == "object" || kind == "list":
			return fmt.Errorf("-set %s: this key holds a list of entries; set it in nanovision.yaml", key)
		}

		parts := strings.Split(key, ".")
		node := root
		for _, part := range parts[:len(parts)-1] {
			node = child(node, part, yaml.MappingNode)
		}
		leaf := parts[len(parts)-1]
		scalar := &yaml.Node{Kind: yaml.ScalarNode, Value: value}
		if kind == "string" || kind == "list of strings" {
			scalar.Tag = "!!str"
		}
		if kind == "list of strings" {
			list := child(node, leaf, yaml.SequenceNode)
			list.Content = append(list.Content, scalar)
		} else {
			*child(node, leaf, yaml.ScalarNode) = *scalar
		}
	}
	if err := root.Decode(c); err != nil {
		return fmt.Errorf("-set: %w", cleanYAMLError(err))
	}
	return nil
}

// validate checks the final configuration for logical errors.
func (c *AppConfig) validate() error {
	reports := len(c.Reports)
	for _, o := range c.Folders {
		reports += len(o.Reports)
	}
	if reports == 0 {
		return errors.New("configuration error: at least one report pattern must be specified")
	}
	if _, err := logging.ParseVerbosity(c.Verbosity); err != nil {
		return fmt.Errorf("invalid verbosity level '%s'", c.Verbosity)
	}
	if err := checkReportTypes(c.ReportTypes); err != nil {
		return err
	}

	switch strings.ToLower(c.Review.FailOn) {
	case "", "never", "warning", "error":
	default:
		return fmt.Errorf("invalid review.fail_on value '%s' (expected 'error', 'warning' or 'never')", c.Review.FailOn)
	}

	if !profileNameRE.MatchString(c.History.Profile) {
		return fmt.Errorf("invalid history.profile %q: use letters, digits, '.', '_' and '-'", c.History.Profile)
	}
	if c.History.MaxDistance < 0 {
		return fmt.Errorf("history.max_distance must be 0 or more, got %d", c.History.MaxDistance)
	}
	switch strings.ToLower(c.VCS.Type) {
	case "", "git", "perforce":
	default:
		return fmt.Errorf("invalid vcs.type %q (expected 'git' or 'perforce'; remove the key to turn version control off)", c.VCS.Type)
	}
	switch c.Run.Kind {
	case "", "local", "submit", "review":
	default:
		return fmt.Errorf("invalid -run-kind %q (expected 'submit', 'review' or 'local')", c.Run.Kind)
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
	if err := c.resolveMetrics(); err != nil {
		return err
	}
	if err := c.resolveFolders(); err != nil {
		return err
	}

	allFilters := append([]string{}, c.FileFilters...)
	for _, pattern := range c.IgnoreFiles {
		allFilters = append(allFilters, "-"+pattern)
	}
	for _, o := range c.Folders {
		for _, pattern := range o.IgnoreFiles {
			allFilters = append(allFilters, "-"+o.Path+"/"+pattern)
		}
	}
	filter, err := filtering.NewDefaultFilter(allFilters, true)
	if err != nil {
		return fmt.Errorf("failed to initialize file filter: %w", err)
	}
	c.FileFilterInstance = filter

	c.VerbosityLevel, _ = logging.ParseVerbosity(c.Verbosity)

	c.InputPairs = nil
	for _, r := range c.Reports {
		c.addInputPair(r, "", ".")
	}
	for _, o := range c.Folders {
		// a folder config writes its paths relative to its own folder
		for _, r := range o.Reports {
			c.addInputPair(r, filepath.FromSlash(o.Path), o.Path)
		}
	}

	c.VCS.Type = strings.ToLower(c.VCS.Type)
	if c.Run.Kind == "" {
		c.Run.Kind = "local"
	}

	c.Review.FailOn = strings.ToLower(c.Review.FailOn)
	if c.Review.FailOn == "" {
		c.Review.FailOn = "never"
	}

	return nil
}

func (c *AppConfig) addInputPair(r ReportInput, base, defaultSource string) {
	pattern := strings.TrimSpace(r.Path)
	if pattern == "" {
		return
	}
	source := strings.TrimSpace(r.Source)
	if source == "" {
		source = defaultSource
	} else if base != "" && !filepath.IsAbs(source) {
		source = filepath.ToSlash(filepath.Join(base, source))
	}
	if base != "" && !filepath.IsAbs(pattern) {
		pattern = filepath.ToSlash(filepath.Join(base, pattern))
	}
	c.InputPairs = append(c.InputPairs, ReportInputPair{ReportPattern: pattern, SourceDir: source, Name: strings.TrimSpace(r.Name)})
}

// resolveMetrics turns metrics.files and metrics.methods into the metric
// lists and the warning ranges of the root config.
func (c *AppConfig) resolveMetrics() error {
	if len(c.Metrics.Files) == 0 {
		c.Metrics.Files = defaultSettings(FileMetricDefs)
	}
	if len(c.Metrics.Methods) == 0 {
		c.Metrics.Methods = defaultSettings(MethodMetricDefs)
	}
	if err := c.applyCliThresholds(); err != nil {
		return err
	}

	c.StatusBands = make(StatusBands)
	var err error
	if c.FileMetrics, err = resolveSettings(c.Metrics.Files, FileMetricDefs, "metrics.files", c.StatusBands); err != nil {
		return err
	}
	if c.MethodMetrics, err = resolveSettings(c.Metrics.Methods, MethodMetricDefs, "metrics.methods", c.StatusBands); err != nil {
		return err
	}
	c.ActiveFileMetrics = activeSet(c.FileMetrics)
	c.ActiveMethodMetrics = activeSet(c.MethodMetrics)
	return nil
}

func defaultSettings(defs []MetricDef) []MetricSetting {
	settings := make([]MetricSetting, len(defs))
	for i, d := range defs {
		settings[i] = MetricSetting{Name: d.Name}
	}
	return settings
}

func activeSet(keys []MetricKey) map[MetricKey]bool {
	set := make(map[MetricKey]bool, len(keys))
	for _, k := range keys {
		set[k] = true
	}
	return set
}

// resolveSettings checks the names of one metric list and adds its warning
// ranges to bands. It returns the metric keys in the listed order.
func resolveSettings(settings []MetricSetting, defs []MetricDef, where string, bands StatusBands) ([]MetricKey, error) {
	var keys []MetricKey
	seen := make(map[string]bool, len(settings))
	for _, s := range settings {
		def, ok := metricByName(defs, s.Name)
		if !ok {
			return nil, fmt.Errorf("%s: unknown metric %q (expected one of: %s)", where, s.Name, strings.Join(metricNames(defs), ", "))
		}
		if seen[s.Name] {
			return nil, fmt.Errorf("%s: metric %q is listed twice", where, s.Name)
		}
		seen[s.Name] = true
		keys = append(keys, def.Key)
		if s.Warning == nil {
			continue
		}
		if !def.Value && s.Warning.Max > 100 {
			return nil, fmt.Errorf("%s: warning range %q of %s is above 100%%", where, s.Warning, s.Name)
		}
		bands[def.Key] = *s.Warning
	}
	return keys, nil
}

// applyCliThresholds applies -threshold name=min..max. A bare name is a file
// metric; "methods.name" is a method metric. A metric that is not listed yet
// is added: asking for its warning range asks to see it.
func (c *AppConfig) applyCliThresholds() error {
	for _, t := range c.cliThresholds {
		name, value, ok := strings.Cut(t, "=")
		if !ok {
			return fmt.Errorf("invalid -threshold %q (expected name=min..max)", t)
		}
		band, err := parseBandString(strings.TrimSpace(value))
		if err != nil {
			return fmt.Errorf("-threshold %s: %w", name, err)
		}
		name = strings.TrimSpace(name)
		settings := &c.Metrics.Files
		if rest, isMethod := strings.CutPrefix(name, "methods."); isMethod {
			name, settings = rest, &c.Metrics.Methods
		} else {
			name = strings.TrimPrefix(name, "files.")
		}
		found := false
		for i := range *settings {
			if (*settings)[i].Name == name {
				(*settings)[i].Warning, found = &band, true
			}
		}
		if !found {
			*settings = append(*settings, MetricSetting{Name: name, Warning: &band})
		}
	}
	return nil
}
