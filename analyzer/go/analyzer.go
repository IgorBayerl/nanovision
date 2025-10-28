package golang

import (
	"fmt"
	"log/slog"
	"strings"

	sitter "github.com/IgorBayerl/nanovision/tree-sitter/go-tree-sitter"
	tsgo "github.com/IgorBayerl/nanovision/tree-sitter/tree-sitter-go/bindings/go"

	"github.com/IgorBayerl/nanovision/analyzer"
)

type GoAnalyzer struct{}

func New() analyzer.Analyzer { return &GoAnalyzer{} }

func (a *GoAnalyzer) Name() string {
	return "Go"
}

func (a *GoAnalyzer) SupportsFile(filePath string) bool {
	return strings.HasSuffix(strings.ToLower(filePath), ".go")
}

const (
	// -------- function & method names --------
	// Captures:
	//   @name     -> function or method name
	//   @receiver -> receiver TYPE only (e.g. T or *pkg.T), not "(x T)"
	funcQueryString = `
    (function_declaration
      name: (identifier) @name)

    (method_declaration
      receiver: (parameter_list
                  (parameter_declaration
                    type: (_) @receiver))
      name: (field_identifier) @name)
  `

	// -------- cyclomatic complexity drivers --------
	// Start at 1, then +1 for each capture here.
	// We capture only non-default cases (default is its own node type).
	complexityQueryString = `
    ;; branches / loops
    (if_statement)  @decision
    (for_statement) @decision

    ;; switch/type-switch/select non-default arms
    (expression_case)     @case
    (type_case)           @case
    (communication_case)  @case

    ;; short-circuit boolean ops (anywhere in the body)
    (binary_expression operator: "&&") @op
    (binary_expression operator: "||") @op
  `
)

func (a *GoAnalyzer) Analyze(sourceCode []byte) (analyzer.AnalysisResult, error) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tsgo.Language())
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

	for m := matches.Next(); m != nil; m = matches.Next() {
		var funcNode *sitter.Node
		funcName := ""
		receiver := ""

		for _, capture := range m.Captures {
			captureName := q.CaptureNames()[capture.Index]
			switch captureName {
			case "name":
				funcName = capture.Node.Utf8Text(sourceCode)
				funcNode = capture.Node.Parent()
			case "receiver":
				receiver = capture.Node.Utf8Text(sourceCode)
			}
		}

		if funcNode == nil {
			continue
		}

		bodyNode := funcNode.ChildByFieldName("body")
		complexity := calculateComplexity(lang, sourceCode, bodyNode)

		start := funcNode.StartPosition().Row + 1
		end := funcNode.EndPosition().Row + 1

		name := funcName
		if strings.TrimSpace(receiver) != "" {
			name = fmt.Sprintf("(%s).%s", receiver, funcName) // e.g. (*MessageBuilder).Greet
		}

		result.Functions = append(result.Functions, analyzer.FunctionMetric{
			Name:                 name,
			Position:             analyzer.Position{StartLine: int(start), EndLine: int(end)},
			CyclomaticComplexity: &complexity,
		})
	}

	return result, nil
}

func calculateComplexity(lang *sitter.Language, src []byte, bodyNode *sitter.Node) int {
	if bodyNode == nil {
		return 1
	}
	complexity := 1

	q, qerr := sitter.NewQuery(lang, complexityQueryString)
	if qerr != nil {
		slog.Warn("Error compiling Go complexity query", "error", qerr)
		return -1
	}
	defer q.Close()

	captureNames := q.CaptureNames()

	qc := sitter.NewQueryCursor()
	defer qc.Close()

	matches := qc.Matches(q, bodyNode, src)

	for m := matches.Next(); m != nil; m = matches.Next() {
		for _, capture := range m.Captures {
			if captureNames[capture.Index] == "case" {
				firstChild := capture.Node.Child(0)
				if firstChild != nil && firstChild.Kind() == "default" {
					continue
				}
			}
			complexity++
		}
	}
	return complexity
}
