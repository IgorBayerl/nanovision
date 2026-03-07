package textsummary

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/IgorBayerl/nanovision/internal/calculator"
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
)

// helper builds a SummaryTree with non-zero values for all metric categories.
func buildTestTree() *model.SummaryTree {
	tree := &model.SummaryTree{
		Metrics: model.CoverageMetrics{
			StatementsValid:   10,
			StatementsCovered: 5,
			LinesValid:        20,
			LinesCovered:      10,
			BranchesValid:     8,
			BranchesCovered:   4,
		},
		Root: &model.DirNode{
			Name: "root",
			Metrics: model.CoverageMetrics{
				StatementsValid:   10,
				StatementsCovered: 5,
				LinesValid:        20,
				LinesCovered:      10,
				BranchesValid:     8,
				BranchesCovered:   4,
			},
			Files: map[string]*model.FileNode{
				"test.go": {
					Name: "test.go",
					Metrics: model.CoverageMetrics{
						StatementsValid:   10,
						StatementsCovered: 5,
						LinesValid:        20,
						LinesCovered:      10,
						BranchesValid:     8,
						BranchesCovered:   4,
					},
				},
			},
		},
	}
	activeActiveFiles := map[config.MetricKey]bool{
		config.LineCoverage:      true,
		config.StatementCoverage: true,
		config.BranchCoverage:    true,
	}
	calculator.CalculateTree(tree, activeActiveFiles, nil)
	return tree
}

func runReport(t *testing.T, fileMetrics []config.MetricKey) string {
	t.Helper()
	tempDir := t.TempDir()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	activeMap := make(map[config.MetricKey]bool, len(fileMetrics))
	for _, m := range fileMetrics {
		activeMap[m] = true
	}

	cfg := &config.AppConfig{
		FileMetrics:       fileMetrics,
		ActiveFileMetrics: activeMap,
	}

	builder := NewTextReportBuilder(tempDir, logger, cfg)
	tree := buildTestTree()

	if err := builder.CreateReport(tree); err != nil {
		t.Fatalf("CreateReport failed: %v", err)
	}

	contentBytes, err := os.ReadFile(filepath.Join(tempDir, "Summary.txt"))
	if err != nil {
		t.Fatalf("Failed to read generated report: %v", err)
	}
	return string(contentBytes)
}

func TestTextReportBuilder_AllThreeMetrics(t *testing.T) {
	content := runReport(t, []config.MetricKey{
		config.LineCoverage,
		config.BranchCoverage,
		config.StatementCoverage,
	})

	for _, want := range []string{"Line Coverage:", "Branch Coverage:", "Statement Coverage:"} {
		if !strings.Contains(content, want) {
			t.Errorf("Expected report to contain %q, but it did not.\nContent:\n%s", want, content)
		}
	}
}

func TestTextReportBuilder_NoStatementCoverage(t *testing.T) {
	content := runReport(t, []config.MetricKey{
		config.LineCoverage,
		config.BranchCoverage,
	})

	for _, want := range []string{"Line Coverage:", "Branch Coverage:"} {
		if !strings.Contains(content, want) {
			t.Errorf("Expected report to contain %q, but it did not.\nContent:\n%s", want, content)
		}
	}
	for _, notWant := range []string{"Statement Coverage:"} {
		if strings.Contains(content, notWant) {
			t.Errorf("Expected report NOT to contain %q, but it did.\nContent:\n%s", notWant, content)
		}
	}
}

func TestTextReportBuilder_SingleMetricOnly(t *testing.T) {
	content := runReport(t, []config.MetricKey{
		config.LineCoverage,
	})

	if !strings.Contains(content, "Line Coverage:") {
		t.Errorf("Expected report to contain 'Line Coverage:', but it did not.\nContent:\n%s", content)
	}
	for _, notWant := range []string{"Statement Coverage:", "Branch Coverage:"} {
		if strings.Contains(content, notWant) {
			t.Errorf("Expected report NOT to contain %q, but it did.\nContent:\n%s", notWant, content)
		}
	}
}

func TestTextReportBuilder_ColumnOrderRespectsFileMetrics(t *testing.T) {
	content := runReport(t, []config.MetricKey{
		config.StatementCoverage,
		config.LineCoverage,
	})

	stmtIdx := strings.Index(content, "Statement Coverage")
	lineIdx := strings.Index(content, "Line Coverage")

	if stmtIdx < 0 || lineIdx < 0 {
		t.Fatalf("Expected both Statement Coverage and Line Coverage in output.\nContent:\n%s", content)
	}
	if stmtIdx > lineIdx {
		t.Errorf("Expected Statement Coverage before Line Coverage, but Statement Coverage at %d, Line Coverage at %d.\nContent:\n%s", stmtIdx, lineIdx, content)
	}
}
