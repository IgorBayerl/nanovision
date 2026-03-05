package htmlreact

import (
	"testing"

	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestPatchStatementMethodProvider_Apply(t *testing.T) {
	provider := PatchStatementMethodProvider{}
	method := &model.MethodMetrics{
		PatchStatementsValid:   10,
		PatchStatementsCovered: 5,
		StatementsValid:        1,
		DiffStatus:             "modified",
	}
	ui := &methodDetail{Metrics: make(map[string]methodMetric)}

	provider.Apply(method, ui)

	assert.Equal(t, "5 / 10", ui.Metrics[MethodUIPatchStmtCoverage].Value)
	assert.NotNil(t, ui.NewStatementCoverage)
	assert.Equal(t, 5, ui.NewStatementCoverage.Covered)
	assert.Equal(t, 10, ui.NewStatementCoverage.Total)
}

func TestPatchLineCoverageMethodProvider_Apply(t *testing.T) {
	provider := PatchLineCoverageMethodProvider{}
	method := &model.MethodMetrics{
		PatchLinesValid:   20,
		PatchLinesCovered: 15,
		DiffStatus:        "added",
	}
	ui := &methodDetail{Metrics: make(map[string]methodMetric)}

	provider.Apply(method, ui)

	assert.Equal(t, "15 / 20", ui.Metrics[MethodUIPatchLineCoverage].Value)
	assert.NotNil(t, ui.NewLinesCoverage)
	assert.Equal(t, 15, ui.NewLinesCoverage.Covered)
	assert.Equal(t, 20, ui.NewLinesCoverage.Total)
}
