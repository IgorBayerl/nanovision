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
	cliInput.FileMetrics = "line_coverage,branch_coverage"

	cfg.mergeCliOverrides(cliInput)
	if err := cfg.validate(); err != nil {
		t.Fatalf("expected no validation error, got: %v", err)
	}
	if err := cfg.computeDerivedFields(); err != nil {
		t.Fatalf("expected no compute error, got: %v", err)
	}

	expected := []MetricKey{LineCoverage, BranchCoverage}
	if !reflect.DeepEqual(cfg.FileMetrics, expected) {
		t.Errorf("expected FileMetrics to be %v, got %v", expected, cfg.FileMetrics)
	}

	if !cfg.ActiveFileMetrics[LineCoverage] || !cfg.ActiveFileMetrics[BranchCoverage] {
		t.Errorf("expected ActiveFileMetrics to contain LineCoverage and BranchCoverage")
	}
	if cfg.ActiveFileMetrics[MethodsHit] {
		t.Errorf("did not expect ActiveFileMetrics to contain MethodsHit")
	}
}

func TestConfig_FileMetrics_Invalid(t *testing.T) {
	// Unhappy Path: User provides an unknown/typo metric name -> validate() returns a clear error
	cfg := GetDefaultConfig()
	cliInput := defaultCLIInput()
	cliInput.ReportPatterns = "report.xml"
	cliInput.SourceDirs = "."
	cliInput.FileMetrics = "lines_coverage" // invalid, should be line_coverage

	cfg.mergeCliOverrides(cliInput)
	err := cfg.validate()
	if err == nil {
		t.Fatal("expected validation error for invalid metric, got nil")
	}

	expectedErrMsg := "unknown file metric 'lines_coverage'"
	if !strings.Contains(err.Error(), expectedErrMsg) {
		t.Errorf("expected error to contain %q, but got: %v", expectedErrMsg, err)
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
	cfg.FileMetrics = []MetricKey{LineCoverage, BranchCoverage}
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
	if !cfg.ActiveFileMetrics[BranchCoverage] {
		t.Error("expected ActiveFileMetrics to contain BranchCoverage")
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

func TestConfig_MethodMetrics_Invalid(t *testing.T) {
	// Unhappy Path: Unknown key in MethodMetrics produces a clear error
	cfg := GetDefaultConfig()
	cliInput := defaultCLIInput()
	cliInput.ReportPatterns = "report.xml"
	cliInput.SourceDirs = "."
	cliInput.MethodMetrics = "nonexistent_metric"

	cfg.mergeCliOverrides(cliInput)
	err := cfg.validate()
	if err == nil {
		t.Fatal("expected validation error for invalid method metric, got nil")
	}

	expectedErrMsg := "unknown method metric 'nonexistent_metric'"
	if !strings.Contains(err.Error(), expectedErrMsg) {
		t.Errorf("expected error to contain %q, but got: %v", expectedErrMsg, err)
	}
}
