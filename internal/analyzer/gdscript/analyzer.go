package gdscript

import (
	"fmt"
	"log/slog"
	"strings"

	sitter "github.com/IgorBayerl/nanovision/tree-sitter/go-tree-sitter"
	tsgdscript "github.com/IgorBayerl/nanovision/tree-sitter/tree-sitter-gdscript/bindings/go"

	"github.com/IgorBayerl/nanovision/internal/analyzer"
)

type GdScriptAnalyzer struct{}

func New() analyzer.Analyzer { return &GdScriptAnalyzer{} }

func (a *GdScriptAnalyzer) Name() string {
	return "gdscript"
}

func (a *GdScriptAnalyzer) SupportsFile(filePath string) bool {
	return strings.HasSuffix(strings.ToLower(filePath), ".gd")
}

const (
	funcQueryString = `
	(function_definition
		name: (name) @name)
	`

	complexityQueryString = `
	(if_statement)    @decision
	(elif_clause)     @decision
	(for_statement)   @decision
	(while_statement) @decision
	(match_statement) @decision

	(binary_operator "and") @op
	(binary_operator "or")  @op
	(binary_operator "&&")  @op
	(binary_operator "||")  @op
	`
)

// Analyze performs static analysis on the provided source code
func (a *GdScriptAnalyzer) Analyze(sourceCode []byte) (analyzer.AnalysisResult, error) {
	parser := sitter.NewParser()
	defer parser.Close()

	// Load the generated GDScript language binding
	lang := sitter.NewLanguage(tsgdscript.Language())
	if err := parser.SetLanguage(lang); err != nil {
		return analyzer.AnalysisResult{}, fmt.Errorf("set language: %w", err)
	}

	tree := parser.Parse(sourceCode, nil)
	if tree == nil {
		return analyzer.AnalysisResult{}, fmt.Errorf("parse returned nil tree")
	}
	defer tree.Close()

	root := tree.RootNode()

	q, qerr := sitter.NewQuery(lang, funcQueryString)
	if qerr != nil {
		return analyzer.AnalysisResult{}, fmt.Errorf("compile function query: %w", qerr)
	}
	defer q.Close()

	qc := sitter.NewQueryCursor()
	defer qc.Close()

	matches := qc.Matches(q, root, sourceCode)

	var result analyzer.AnalysisResult

	// Iterate over all function definitions found
	for m := matches.Next(); m != nil; m = matches.Next() {
		var funcNode *sitter.Node
		funcName := ""

		for _, capture := range m.Captures {
			captureName := q.CaptureNames()[capture.Index]
			if captureName == "name" {
				funcName = capture.Node.Utf8Text(sourceCode)
				funcNode = capture.Node.Parent()
			}
		}

		if funcNode == nil {
			continue
		}

		// GDScript function body is typically a child named "body"
		bodyNode := funcNode.ChildByFieldName("body")

		// Calculate complexity strictly within this function's body
		complexity := calculateComplexity(lang, sourceCode, bodyNode)

		// Tree-sitter rows are 0-indexed; convert to 1-indexed for reports
		start := funcNode.StartPosition().Row + 1
		end := funcNode.EndPosition().Row + 1

		result.Functions = append(result.Functions, analyzer.FunctionMetric{
			Name:                 funcName,
			Position:             analyzer.Position{StartLine: int(start), EndLine: int(end)},
			CyclomaticComplexity: &complexity,
		})
	}

	return result, nil
}

func calculateComplexity(lang *sitter.Language, src []byte, bodyNode *sitter.Node) int {
	if bodyNode == nil {
		return 1 // Default complexity for empty function
	}
	complexity := 1

	q, qerr := sitter.NewQuery(lang, complexityQueryString)
	if qerr != nil {
		slog.Warn("Error compiling GDScript complexity query", "error", qerr)
		return -1
	}
	defer q.Close()

	qc := sitter.NewQueryCursor()
	defer qc.Close()

	// Execute query only within the scope of the function body
	matches := qc.Matches(q, bodyNode, src)

	for m := matches.Next(); m != nil; m = matches.Next() {
		// In GDScript, we just count every match defined in the query string
		// (unlike Go where we might filter out 'default' cases)
		complexity++
	}

	return complexity
}
