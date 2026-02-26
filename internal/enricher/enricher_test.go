package enricher

import (
	"testing"

	"github.com/IgorBayerl/nanovision/internal/analyzer"
	"github.com/IgorBayerl/nanovision/internal/model"
)

func TestCalculateStatementCoverage(t *testing.T) {
	fileNode := &model.FileNode{
		Lines: map[int]model.LineMetrics{
			10: {Hits: 1},
			11: {Hits: 5},
			12: {Hits: 0},
		},
	}

	analysis := analyzer.AnalysisResult{
		Functions: []analyzer.FunctionMetric{
			{
				Name: "TestFunction",
				Position: analyzer.Position{
					StartLine: 9,
					EndLine:   20,
				},
			},
		},
		Statements: []analyzer.StatementMetric{
			{
				StartLine: 10,
				EndLine:   11,
				Type:      "block",
			},
			{
				StartLine: 15,
				EndLine:   15,
				Type:      "call",
			},
		},
	}

	enricher := &Enricher{}
	enricher.applyAnalysisToFileNode(fileNode, analysis)

	if fileNode.Metrics.StatementsValid != 2 {
		t.Errorf("expected 2 valid statements in file, got %d", fileNode.Metrics.StatementsValid)
	}
	if fileNode.Metrics.StatementsCovered != 1 {
		t.Errorf("expected 1 covered statement in file, got %d", fileNode.Metrics.StatementsCovered)
	}

	if len(fileNode.Methods) != 1 {
		t.Fatalf("expected 1 method, got %d", len(fileNode.Methods))
	}

	method := fileNode.Methods[0]
	if method.StatementsValid != 2 {
		t.Errorf("expected 2 valid statements in method, got %d", method.StatementsValid)
	}
	if method.StatementsCovered != 1 {
		t.Errorf("expected 1 covered statement in method, got %d", method.StatementsCovered)
	}
}
