package golang

import (
	"fmt"
	"log/slog"
	"strings"

	sitter "github.com/IgorBayerl/nanovision/tree-sitter/go-tree-sitter"
	tsgo "github.com/IgorBayerl/nanovision/tree-sitter/tree-sitter-go/bindings/go"

	"github.com/IgorBayerl/nanovision/internal/analyzer"
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
	funcQueryString = `
    (function_declaration
      name: (identifier) @name)

    (method_declaration
      receiver: (parameter_list
                  (parameter_declaration
                    type: (_) @receiver))
      name: (field_identifier) @name)
  `

	complexityQueryString = `
    (if_statement)  @decision
    (for_statement) @decision
    (expression_case)     @case
    (type_case)           @case
    (communication_case)  @case
    (binary_expression operator: "&&") @op
    (binary_expression operator: "||") @op
  `

	// Precision Query for Atomic Statements
	statementQueryString = `
		; --- ATOMIC STATEMENTS ---
		; These are always statements, regardless of where they are.
		(expression_statement) @stmt
		(send_statement) @stmt
		(inc_statement) @stmt
		(dec_statement) @stmt
		(assignment_statement) @stmt
		(short_var_declaration) @stmt
		(return_statement) @stmt
		(go_statement) @stmt
		(defer_statement) @stmt
		(break_statement) @stmt
		(continue_statement) @stmt
		(fallthrough_statement) @stmt
		(goto_statement) @stmt

		; --- CONTROL FLOW LOGIC (No Field Names) ---
		
		; IF Conditions: Capture logic expressions inside 'if'
		; We explicitly list expression types to avoid capturing the 'block' or 'init' statement (which is already caught above).
		(if_statement 
			[
				(binary_expression)
				(unary_expression)
				(call_expression)
				(identifier)
				(parenthesized_expression)
				(selector_expression)
				(index_expression)
			] @stmt
		)

		; FOR Conditions: While-style loops
		(for_statement 
			[
				(binary_expression)
				(unary_expression)
				(call_expression)
				(identifier)
			] @stmt
		)

		; FOR Clauses: Standard (i=0; i<n; i++) and Range (i := range x)
		; These are specific nodes in the grammar, so we capture the node itself.
		(for_clause) @stmt
		(range_clause) @stmt

		; SWITCH Value: The expression being switched on
		(expression_switch_statement 
			[
				(binary_expression)
				(call_expression)
				(identifier)
				(selector_expression)
				(index_expression)
			] @stmt
		)

		; SELECT: Capture the keyword as the anchor
		(select_statement "select" @stmt)
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

	stmtQ, stmtQErr := sitter.NewQuery(lang, statementQueryString)
	if stmtQErr != nil {
		return analyzer.AnalysisResult{}, fmt.Errorf("compile statement query: %w", stmtQErr)
	}
	defer stmtQ.Close()

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

	stmtQc := sitter.NewQueryCursor()
	defer stmtQc.Close()
	stmtMatches := stmtQc.Matches(stmtQ, root, sourceCode)

	for m := stmtMatches.Next(); m != nil; m = stmtMatches.Next() {
		for _, capture := range m.Captures {
			node := capture.Node
			start := node.StartPosition().Row + 1
			end := node.EndPosition().Row + 1
			result.Statements = append(result.Statements, analyzer.StatementMetric{
				StartLine: int(start),
				EndLine:   int(end),
				Type:      node.Kind(),
			})
		}
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
