package evaluators

import (
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/status"
)

// Registry maps each metric key to its Evaluator implementation.
// The annotation engine iterates over the active metric keys and looks
// up evaluators from this registry.
var Registry = map[config.MetricKey]status.Evaluator{
	config.LineCoverage:            LineCoverageEvaluator{},
	config.BranchCoverage:          BranchCoverageEvaluator{},
	config.StatementCoverage:       StatementCoverageEvaluator{},
	config.MaxCyclomaticComplexity: MaxComplexityEvaluator{},
}
