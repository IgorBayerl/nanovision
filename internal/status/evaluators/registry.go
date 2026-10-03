// Package evaluators builds the status rule of every metric from the metric
// tables in internal/config. A metric needs no code here: its direction
// (higher or lower is better) is a field of its definition.
package evaluators

import (
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/status"
)

// Registry maps each metric key to its status rule.
var Registry = func() map[config.MetricKey]status.Evaluator {
	registry := make(map[config.MetricKey]status.Evaluator)
	for _, def := range config.FileMetricDefs {
		registry[def.Key] = evaluator{def: def, scope: status.FileScope}
	}
	for _, def := range config.MethodMetricDefs {
		registry[def.Key] = evaluator{def: def, scope: status.MethodScope}
	}
	return registry
}()

type evaluator struct {
	def   config.MetricDef
	scope status.MetricScope
}

func (e evaluator) Key() config.MetricKey               { return e.def.Key }
func (e evaluator) Name() string                        { return e.def.TitleLabel() }
func (e evaluator) Description() string                 { return e.def.Doc }
func (e evaluator) SupportedScopes() status.MetricScope { return e.scope }

func (e evaluator) IsApplicable(caps status.Capabilities) bool {
	switch e.def.Needs {
	case config.NeedsMethods:
		return caps.HasMethodCoverage
	case config.NeedsStatements:
		return caps.HasStatementCoverage
	}
	return true
}

// Evaluate classifies the calculated value of the metric against the band.
func (e evaluator) Evaluate(m model.CoverageMetrics, band *config.Band) (status.RiskLevel, bool) {
	var value float64
	switch detail := m.Calculated[e.def.Key].(type) {
	case model.CoverageDetail:
		value = detail.Percentage
	case model.ScoreDetail:
		value = detail.Value
	default:
		return "", false // not calculated: the data is absent
	}
	if e.def.LowerIsBetter {
		return status.ClassifyLowerIsBetter(value, band)
	}
	return status.ClassifyHigherIsBetter(value, band)
}
