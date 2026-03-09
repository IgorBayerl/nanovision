package evaluators

import (
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/status"
)

// Registry maps each metric key to its Evaluator implementation.
// The annotation engine iterates over the active metric keys and looks
// up evaluators from this registry.
var Registry = map[config.MetricKey]status.Evaluator{
	config.LineCoverage:                 LineCoverageEvaluator{},
	config.BranchCoverage:               BranchCoverageEvaluator{},
	config.StatementCoverage:            StatementCoverageEvaluator{},
	config.MaxCyclomaticComplexity:      MaxComplexityEvaluator{},
	config.MethodsHit:                   MethodsHitEvaluator{},
	config.MethodsFullyCovered:          MethodsFullyCoveredEvaluator{},
	config.PatchLineCoverage:            PatchLineCoverageEvaluator{},
	config.PatchMethodsHit:              PatchMethodsHitEvaluator{},
	config.PatchStatementCoverage:       PatchStatementCoverageEvaluator{},
	config.StatementMethodsHit:          StatementMethodsHitEvaluator{},
	config.StatementMethodsFullyCovered: StatementMethodsFullyCoveredEvaluator{},
	config.PatchStatementMethodsHit:     PatchStatementMethodsHitEvaluator{},

	config.MethodLineCoverage:           MethodLineCoverageEvaluator{},
	config.MethodPatchLineCoverage:      MethodPatchLineCoverageEvaluator{},
	config.MethodStatementCoverage:      MethodStatementCoverageEvaluator{},
	config.MethodPatchStatementCoverage: MethodPatchStatementCoverageEvaluator{},
	config.CyclomaticComplexity:         CyclomaticComplexityEvaluator{},
	config.MethodCrapScore:              MethodCrapScoreEvaluator{},
	config.MethodPatchCrapScore:         MethodPatchCrapScoreEvaluator{},
	config.MethodExposedRisk:            MethodExposedRiskEvaluator{},
	config.MethodDefectProbability:      MethodDefectProbabilityEvaluator{},
}
