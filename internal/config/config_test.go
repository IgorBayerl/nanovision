package config

import (
	"reflect"
	"strings"
	"testing"
)

func defaultCLIInput() RawConfigInput {
	return RawConfigInput{
		ReportTypes: "TextSummary,Html",
		OutputDir:   "coverage-report",
		LogFormat:   "text",
		Verbosity:   "Info",
	}
}

func TestConfig_DisplayMetrics_Default(t *testing.T) {
	// Happy Path: No display metrics provided in YAML or CLI -> defaults to DefaultDisplayMetrics
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

	if !reflect.DeepEqual(cfg.DisplayMetrics, DefaultDisplayMetrics) {
		t.Errorf("expected DisplayMetrics to be default, got %v", cfg.DisplayMetrics)
	}

	for _, dm := range DefaultDisplayMetrics {
		if !cfg.ActiveMetrics[dm] {
			t.Errorf("expected ActiveMetrics to contain %s", dm)
		}
	}
}

func TestConfig_DisplayMetrics_CLIOverride(t *testing.T) {
	// Happy Path: Valid comma-separated string provided via CLI -> successfully builds ActiveMetrics map
	cfg := GetDefaultConfig()
	cliInput := defaultCLIInput()
	cliInput.ReportPatterns = "report.xml"
	cliInput.SourceDirs = "."
	cliInput.DisplayMetrics = "line_coverage,branch_coverage"

	cfg.mergeCliOverrides(cliInput)
	if err := cfg.validate(); err != nil {
		t.Fatalf("expected no validation error, got: %v", err)
	}
	if err := cfg.computeDerivedFields(); err != nil {
		t.Fatalf("expected no compute error, got: %v", err)
	}

	expected := []MetricKey{LineCoverage, BranchCoverage}
	if !reflect.DeepEqual(cfg.DisplayMetrics, expected) {
		t.Errorf("expected DisplayMetrics to be %v, got %v", expected, cfg.DisplayMetrics)
	}

	if !cfg.ActiveMetrics[LineCoverage] || !cfg.ActiveMetrics[BranchCoverage] {
		t.Errorf("expected ActiveMetrics to contain LineCoverage and BranchCoverage")
	}
	if cfg.ActiveMetrics[MethodsHit] {
		t.Errorf("did not expect ActiveMetrics to contain MethodsHit")
	}
}

func TestConfig_DisplayMetrics_Invalid(t *testing.T) {
	// Unhappy Path: User provides an unknown/typo metric name -> validate() returns a clear error
	cfg := GetDefaultConfig()
	cliInput := defaultCLIInput()
	cliInput.ReportPatterns = "report.xml"
	cliInput.SourceDirs = "."
	cliInput.DisplayMetrics = "lines_coverage" // invalid, should be line_coverage

	cfg.mergeCliOverrides(cliInput)
	err := cfg.validate()
	if err == nil {
		t.Fatal("expected validation error for invalid metric, got nil")
	}

	expectedErrMsg := "unknown display metric 'lines_coverage'"
	if !strings.Contains(err.Error(), expectedErrMsg) {
		t.Errorf("expected error to contain %q, but got: %v", expectedErrMsg, err)
	}
}

func TestConfig_DisplayMetrics_YAML(t *testing.T) {
	// Happy Path: YAML provided (Testing via manual config assignment)
	cfg := GetDefaultConfig()
	cfg.ReportPatterns = []string{"report.xml"}
	cfg.SourceDirs = []string{"."}
	cfg.DisplayMetrics = []MetricKey{StatementCoverage, MethodsHit}

	cliInput := defaultCLIInput() // no override
	cfg.mergeCliOverrides(cliInput)
	if err := cfg.validate(); err != nil {
		t.Fatalf("expected no validation error, got: %v", err)
	}
	if err := cfg.computeDerivedFields(); err != nil {
		t.Fatalf("expected no compute error, got: %v", err)
	}

	expected := []MetricKey{StatementCoverage, MethodsHit}
	if !reflect.DeepEqual(cfg.DisplayMetrics, expected) {
		t.Errorf("expected DisplayMetrics to be %v, got %v", expected, cfg.DisplayMetrics)
	}
}
