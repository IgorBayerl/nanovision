package calculator

import "github.com/IgorBayerl/nanovision/internal/config"

var FileRegistry = map[config.MetricKey]MetricCalculator{
	config.LineCoverage:            LineCoverageCalculator{},
	config.StatementCoverage:       StatementCoverageCalculator{},
	config.MethodsHit:              MethodsHitCalculator{},
	config.MethodsFullyCovered:     MethodsFullyCoveredCalculator{},
	config.PatchLineCoverage:       PatchLineCoverageCalculator{},
	config.PatchStatementCoverage:  PatchStatementCoverageCalculator{},
	config.PatchMethodsHit:         PatchMethodsHitCalculator{},
	config.MaxCyclomaticComplexity: MaxCyclomaticComplexityCalculator{},
}

var MethodRegistry = map[config.MetricKey]MethodMetricCalculator{
	config.MethodLineCoverage:           MethodLineCoverageCalculator{},
	config.MethodStatementCoverage:      MethodStatementCoverageCalculator{},
	config.MethodPatchLineCoverage:      MethodPatchLineCoverageCalculator{},
	config.MethodPatchStatementCoverage: MethodPatchStatementCoverageCalculator{},
	config.CyclomaticComplexity:         MethodCyclomaticComplexityCalculator{},
	config.MethodCrapScore:              MethodCrapScoreCalculator{},
	config.MethodExposedRisk:            MethodExposedRiskCalculator{},
}
