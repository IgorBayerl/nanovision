package htmlreact

import (
	"log/slog"
	"os"
	"testing"

	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
)

func TestHtmlReactReportBuilder_ActiveFileMetricsFilter(t *testing.T) {
	tempDir := t.TempDir()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg := &config.AppConfig{
		ActiveFileMetrics: map[config.MetricKey]bool{
			"line_coverage": true,
			// implicitly missing statement_coverage, branch_coverage
		},
	}

	builder := NewHtmlReactReportBuilder(tempDir, logger, false, cfg)

	tree := &model.SummaryTree{
		Metrics: model.CoverageMetrics{
			StatementsValid:   10,
			StatementsCovered: 5,
			LinesValid:        10,
			LinesCovered:      5,
			BranchesValid:     10,
			BranchesCovered:   5,
		},
		Root: &model.DirNode{
			Name: "root",
		},
	}

	b, ok := builder.(*HtmlReactReportBuilder)
	if !ok {
		t.Fatalf("Expected HtmlReactReportBuilder")
	}

	totalsData := b.buildTotals(tree, 1, 0)

	if totalsData.LineCoverage == nil {
		t.Errorf("Expected LineCoverage to be present, got nil")
	}
	if totalsData.StatementCoverage != nil {
		t.Errorf("Expected StatementCoverage to be nil, got %+v", totalsData.StatementCoverage)
	}
	if totalsData.BranchCoverage != nil {
		t.Errorf("Expected BranchCoverage to be nil, got %+v", totalsData.BranchCoverage)
	}
}
