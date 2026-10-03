package status_test

import (
	"testing"

	"github.com/IgorBayerl/nanovision/internal/calculator"
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/status"
	"github.com/IgorBayerl/nanovision/internal/status/evaluators"
	"github.com/stretchr/testify/assert"
)

// A folder with its own warning range classifies its files with it; the rest
// of the project keeps the root range.
func TestAnnotate_FolderWarningRanges(t *testing.T) {
	seventy := model.CoverageMetrics{StatementsCovered: 70, StatementsValid: 100}
	file := func(path string) *model.FileNode { return &model.FileNode{Path: path, Metrics: seventy} }

	strict := &model.DirNode{Path: "services/strict", Metrics: seventy, Files: map[string]*model.FileNode{"a.go": file("services/strict/a.go")}}
	other := &model.DirNode{Path: "services/other", Metrics: seventy, Files: map[string]*model.FileNode{"b.go": file("services/other/b.go")}}
	services := &model.DirNode{Path: "services", Metrics: seventy, Subdirs: map[string]*model.DirNode{"strict": strict, "other": other}}
	tree := &model.SummaryTree{Root: &model.DirNode{Path: ".", Metrics: seventy, Subdirs: map[string]*model.DirNode{"services": services}}}

	cfg := &config.AppConfig{
		ActiveFileMetrics: map[config.MetricKey]bool{config.StatementCoverage: true},
		StatusBands:       config.StatusBands{config.StatementCoverage: {Min: 50, Max: 60}},
		FolderBands: []config.FolderBand{
			{Path: "services/strict", Bands: config.StatusBands{config.StatementCoverage: {Min: 80, Max: 90}}},
		},
	}
	calculator.CalculateTree(tree, cfg.ActiveFileMetrics, nil)
	status.Annotate(tree, cfg, status.Capabilities{HasStatementCoverage: true}, evaluators.Registry)

	assert.Equal(t, "safe", tree.Root.Statuses[config.StatementCoverage])
	assert.Equal(t, "safe", other.Files["b.go"].Statuses[config.StatementCoverage])
	assert.Equal(t, "danger", strict.Statuses[config.StatementCoverage], "the folder itself uses its range")
	assert.Equal(t, "danger", strict.Files["a.go"].Statuses[config.StatementCoverage])
}
