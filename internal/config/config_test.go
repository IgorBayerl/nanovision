package config

import (
	"reflect"
	"testing"

	"go.yaml.in/yaml/v3"
)

func defaultCLIInput() RawConfigInput {
	return RawConfigInput{
		ReportTypes: "TextSummary,Html",
		OutputDir:   "coverage-report",
		LogFormat:   "text",
		Verbosity:   "Info",
	}
}

func TestConfig_FileMetrics_Default(t *testing.T) {
	// Happy Path: No file metrics provided in YAML or CLI -> defaults to DefaultFileMetrics
	cfg := GetDefaultConfig()
	cliInput := defaultCLIInput()
	cliInput.ReportPatterns = "report.xml"
	cliInput.SourceDirs = "."

	cfg.mergeCliOverrides(cliInput)
	if err := cfg.validate(); err != nil {
		t.Fatalf("expected no validation error, got: %v", err)
	}
	if err := cfg.computeDerivedFields(); err != nil {
		t.Fatalf("expected no compute error, got: %v", err)
	}

	if !reflect.DeepEqual(cfg.FileMetrics, DefaultFileMetrics) {
		t.Errorf("expected FileMetrics to be default, got %v", cfg.FileMetrics)
	}

	for _, dm := range DefaultFileMetrics {
		if !cfg.ActiveFileMetrics[dm] {
			t.Errorf("expected ActiveFileMetrics to contain %s", dm)
		}
	}
}

func TestConfig_FileMetrics_CLIOverride(t *testing.T) {
	// Happy Path: Valid comma-separated string provided via CLI -> successfully builds ActiveFileMetrics map
	cfg := GetDefaultConfig()
	cliInput := defaultCLIInput()
	cliInput.ReportPatterns = "report.xml"
	cliInput.SourceDirs = "."
	cliInput.FileMetrics = "line_coverage,statement_coverage"

	cfg.mergeCliOverrides(cliInput)
	if err := cfg.validate(); err != nil {
		t.Fatalf("expected no validation error, got: %v", err)
	}
	if err := cfg.computeDerivedFields(); err != nil {
		t.Fatalf("expected no compute error, got: %v", err)
	}

	expected := []MetricKey{LineCoverage, StatementCoverage}
	if !reflect.DeepEqual(cfg.FileMetrics, expected) {
		t.Errorf("expected FileMetrics to be %v, got %v", expected, cfg.FileMetrics)
	}

	if !cfg.ActiveFileMetrics[LineCoverage] || !cfg.ActiveFileMetrics[StatementCoverage] {
		t.Errorf("expected ActiveFileMetrics to contain LineCoverage and StatementCoverage")
	}
	if cfg.ActiveFileMetrics[MethodsHit] {
		t.Errorf("did not expect ActiveFileMetrics to contain MethodsHit")
	}
}

func TestConfig_FileMetrics_UnknownKeysPassConfig(t *testing.T) {
	// After refactoring, config.validate() no longer rejects unknown metric keys.
	// Validation is handled by main.go using the evaluator Registry.
	cfg := GetDefaultConfig()
	cliInput := defaultCLIInput()
	cliInput.ReportPatterns = "report.xml"
	cliInput.SourceDirs = "."
	cliInput.FileMetrics = "lines_coverage" // would be invalid, but config no longer validates this

	cfg.mergeCliOverrides(cliInput)
	err := cfg.validate()
	if err != nil {
		t.Fatalf("expected no validation error (validation moved to main.go), got: %v", err)
	}
}

func TestConfig_FileMetrics_YAML(t *testing.T) {
	// Happy Path: YAML provided (Testing via manual config assignment)
	cfg := GetDefaultConfig()
	cfg.ReportPatterns = []string{"report.xml"}
	cfg.SourceDirs = []string{"."}
	cfg.FileMetrics = []MetricKey{StatementCoverage, MethodsHit}

	cliInput := defaultCLIInput() // no override
	cfg.mergeCliOverrides(cliInput)
	if err := cfg.validate(); err != nil {
		t.Fatalf("expected no validation error, got: %v", err)
	}
	if err := cfg.computeDerivedFields(); err != nil {
		t.Fatalf("expected no compute error, got: %v", err)
	}

	expected := []MetricKey{StatementCoverage, MethodsHit}
	if !reflect.DeepEqual(cfg.FileMetrics, expected) {
		t.Errorf("expected FileMetrics to be %v, got %v", expected, cfg.FileMetrics)
	}
}

func TestConfig_FileAndMethodMetrics_YAML(t *testing.T) {
	// Verify both ActiveFileMetrics and ActiveMethodMetrics are populated correctly
	cfg := GetDefaultConfig()
	cfg.ReportPatterns = []string{"report.xml"}
	cfg.SourceDirs = []string{"."}
	cfg.FileMetrics = []MetricKey{LineCoverage, MethodsHit}
	cfg.MethodMetrics = []MetricKey{StatementCoverage, MaxCyclomaticComplexity}

	cliInput := defaultCLIInput()
	cfg.mergeCliOverrides(cliInput)
	if err := cfg.validate(); err != nil {
		t.Fatalf("expected no validation error, got: %v", err)
	}
	if err := cfg.computeDerivedFields(); err != nil {
		t.Fatalf("expected no compute error, got: %v", err)
	}

	// Check ActiveFileMetrics
	if !cfg.ActiveFileMetrics[LineCoverage] {
		t.Error("expected ActiveFileMetrics to contain LineCoverage")
	}
	if !cfg.ActiveFileMetrics[MethodsHit] {
		t.Error("expected ActiveFileMetrics to contain MethodsHit")
	}
	if cfg.ActiveFileMetrics[StatementCoverage] {
		t.Error("did not expect ActiveFileMetrics to contain StatementCoverage")
	}

	// Check ActiveMethodMetrics
	if !cfg.ActiveMethodMetrics[StatementCoverage] {
		t.Error("expected ActiveMethodMetrics to contain StatementCoverage")
	}
	if !cfg.ActiveMethodMetrics[MaxCyclomaticComplexity] {
		t.Error("expected ActiveMethodMetrics to contain MaxCyclomaticComplexity")
	}
	if cfg.ActiveMethodMetrics[LineCoverage] {
		t.Error("did not expect ActiveMethodMetrics to contain LineCoverage")
	}
}

func TestConfig_MethodMetrics_UnknownKeysPassConfig(t *testing.T) {
	// After refactoring, config.validate() no longer rejects unknown metric keys.
	// Validation is handled by main.go using the evaluator Registry.
	cfg := GetDefaultConfig()
	cliInput := defaultCLIInput()
	cliInput.ReportPatterns = "report.xml"
	cliInput.SourceDirs = "."
	cliInput.MethodMetrics = "nonexistent_metric"

	cfg.mergeCliOverrides(cliInput)
	err := cfg.validate()
	if err != nil {
		t.Fatalf("expected no validation error (validation moved to main.go), got: %v", err)
	}
}

func TestConfig_History(t *testing.T) {
	cfg := GetDefaultConfig()
	if err := yaml.Unmarshal([]byte("history:\n  store: .nanovision\nvcs:\n  type: Git\n"), cfg); err != nil {
		t.Fatal(err)
	}
	cli := defaultCLIInput()
	cli.ReportPatterns, cli.SourceDirs = "report.xml", "."
	cli.RunKind, cli.Revision, cli.Profile = "Submit", " //game/main@118432 ", "unit-win64"

	cfg.mergeCliOverrides(cli)
	if err := cfg.validate(); err != nil {
		t.Fatalf("expected no validation error, got: %v", err)
	}
	if err := cfg.computeDerivedFields(); err != nil {
		t.Fatal(err)
	}

	if cfg.History.Store != ".nanovision" || cfg.History.IsServer() {
		t.Errorf("store: got %+v", cfg.History)
	}
	if cfg.History.MaxDistance != 50 || cfg.History.KeepLocal != 20 {
		t.Errorf("YAML without these keys must keep the defaults, got %+v", cfg.History)
	}
	if cfg.History.Profile != "unit-win64" {
		t.Errorf("the flag wins over the default profile, got %q", cfg.History.Profile)
	}
	if cfg.VCS.Type != "git" {
		t.Errorf("vcs type is case-insensitive, got %q", cfg.VCS.Type)
	}
	if cfg.Run.Kind != "submit" || cfg.Run.Revision != "//game/main@118432" {
		t.Errorf("run flags: got %+v", cfg.Run)
	}
	if !(HistoryConfig{Store: "https://nv.studio.local:7070"}).IsServer() {
		t.Error("an http URL is a team server")
	}
}

func TestConfig_HistoryValidation(t *testing.T) {
	for name, mutate := range map[string]func(*AppConfig, *RawConfigInput){
		"bad profile":  func(c *AppConfig, _ *RawConfigInput) { c.History.Profile = "unit win64" },
		"bad vcs":      func(c *AppConfig, _ *RawConfigInput) { c.VCS.Type = "svn" },
		"bad run kind": func(_ *AppConfig, cli *RawConfigInput) { cli.RunKind = "nightly" },
		"bad distance": func(c *AppConfig, _ *RawConfigInput) { c.History.MaxDistance = -1 },
	} {
		t.Run(name, func(t *testing.T) {
			cfg := GetDefaultConfig()
			cli := defaultCLIInput()
			cli.ReportPatterns, cli.SourceDirs = "report.xml", "."
			mutate(cfg, &cli)
			cfg.mergeCliOverrides(cli)
			if err := cfg.validate(); err == nil {
				t.Error("expected a validation error")
			}
		})
	}
}
