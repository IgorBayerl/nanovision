package evaluators

import (
	"testing"

	"github.com/IgorBayerl/nanovision/internal/calculator"
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/status"
)

func TestEvaluate(t *testing.T) {
	band := &config.Band{Min: 60, Max: 75}
	coverage := func(pct float64) model.CoverageMetrics {
		return model.CoverageMetrics{Calculated: map[config.MetricKey]any{config.StatementCoverage: model.CoverageDetail{Percentage: pct}}}
	}
	score := func(v float64) model.CoverageMetrics {
		return model.CoverageMetrics{Calculated: map[config.MetricKey]any{config.MaxCyclomaticComplexity: model.ScoreDetail{Value: v}}}
	}
	for _, tc := range []struct {
		key     config.MetricKey
		metrics model.CoverageMetrics
		want    status.RiskLevel
	}{
		{config.StatementCoverage, coverage(59.9), status.RiskDanger},
		{config.StatementCoverage, coverage(60), status.RiskWarning},
		{config.StatementCoverage, coverage(75.1), status.RiskSafe},
		// lower is better: the same band reads the other way
		{config.MaxCyclomaticComplexity, score(59), status.RiskSafe},
		{config.MaxCyclomaticComplexity, score(75), status.RiskWarning},
		{config.MaxCyclomaticComplexity, score(76), status.RiskDanger},
	} {
		got, ok := Registry[tc.key].Evaluate(tc.metrics, band)
		if !ok || got != tc.want {
			t.Errorf("%s: got %q (%v), want %q", tc.key, got, ok, tc.want)
		}
	}

	if _, ok := Registry[config.StatementCoverage].Evaluate(model.CoverageMetrics{}, band); ok {
		t.Error("a metric that was not calculated has no status")
	}
	if _, ok := Registry[config.StatementCoverage].Evaluate(coverage(50), nil); ok {
		t.Error("a metric without a warning range has no status")
	}
	if Registry[config.MethodsHit].IsApplicable(status.Capabilities{}) {
		t.Error("methods_hit needs method data")
	}
}

// A metric is one table entry in config plus one calculator. This fails when
// one of the two is missing.
func TestEveryMetricHasACalculator(t *testing.T) {
	for _, def := range config.FileMetricDefs {
		if _, ok := calculator.FileRegistry[def.Key]; !ok {
			t.Errorf("file metric %s has no calculator in calculator.FileRegistry", def.Key)
		}
	}
	for _, def := range config.MethodMetricDefs {
		if _, ok := calculator.MethodRegistry[def.Key]; !ok {
			t.Errorf("method metric %s has no calculator in calculator.MethodRegistry", def.Key)
		}
	}
	if len(calculator.FileRegistry) != len(config.FileMetricDefs) || len(calculator.MethodRegistry) != len(config.MethodMetricDefs) {
		t.Error("a calculator has no entry in the metric tables of internal/config/metrics.go")
	}
}
