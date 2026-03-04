package textsummary

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
)

func TestTextReportBuilder_ActiveFileMetricsFilter(t *testing.T) {
	tempDir := t.TempDir()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg := &config.AppConfig{
		ActiveFileMetrics: map[config.MetricKey]bool{
			"statement_coverage": true,
			// explicitly missing line_coverage and branch_coverage
		},
	}

	builder := NewTextReportBuilder(tempDir, logger, cfg)

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
			Metrics: model.CoverageMetrics{
				StatementsValid:   10,
				StatementsCovered: 5,
				LinesValid:        10,
				LinesCovered:      5,
			},
			Files: map[string]*model.FileNode{
				"test.go": {
					Name: "test.go",
					Metrics: model.CoverageMetrics{
						StatementsValid:   10,
						StatementsCovered: 5,
						LinesValid:        10,
						LinesCovered:      5,
					},
				},
			},
		},
	}

	err := builder.CreateReport(tree)
	if err != nil {
		t.Fatalf("CreateReport failed: %v", err)
	}

	contentBytes, err := os.ReadFile(filepath.Join(tempDir, "Summary.txt"))
	if err != nil {
		t.Fatalf("Failed to read generated report: %v", err)
	}
	content := string(contentBytes)

	if !strings.Contains(content, "Statement coverage:") {
		t.Errorf("Expected report to contain 'Statement coverage:', but it did not.\nContent:\n%s", content)
	}
	if strings.Contains(content, "Line coverage:") {
		t.Errorf("Expected report NOT to contain 'Line coverage:', but it did.\nContent:\n%s", content)
	}
	if strings.Contains(content, "Branch coverage:") {
		t.Errorf("Expected report NOT to contain 'Branch coverage:', but it did.\nContent:\n%s", content)
	}
}
