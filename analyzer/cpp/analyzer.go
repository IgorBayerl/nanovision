package cpp

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"

	sitter "github.com/tree-sitter/go-tree-sitter"
	tscpp "github.com/tree-sitter/tree-sitter-cpp/bindings/go"

	"github.com/IgorBayerl/AdlerCov/analyzer"
)

type CppAnalyzer struct{}

func New() analyzer.Analyzer { return &CppAnalyzer{} }

func (a *CppAnalyzer) Name() string {
	return "C++"
}

func (a *CppAnalyzer) SupportsFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".c", ".cc", ".cpp", ".cxx", ".c++", ".hpp", ".hh", ".hxx", ".h":
		return true
	default:
		return false
	}
}

const (
	// -------- function & method names --------
	// Captures:
	//   @name       -> unqualified function/method name (identifier)
	//   @qualified  -> full qualified identifier (e.g., Namespace::Class::method)
	funcQueryString = `
    ; free function: int foo(...) { ... }
    (function_definition
      declarator: (function_declarator
        declarator: (identifier) @name))

    ; qualified member: Ret Class::method(...) { ... }
    (function_definition
      declarator: (function_declarator
        declarator: (qualified_identifier
          name: (identifier) @name) @qualified))
  `

	// -------- cyclomatic complexity drivers (C++) --------
	// Start at 1, then +1 for each capture here.
	// We capture *non-default* switch arms via case_statement.
	complexityQueryString = `
    ;; branches / loops
    (if_statement)     @decision
    (for_statement)    @decision
    (while_statement)  @decision
    (do_statement)     @decision

    ;; switch cases (non-default; default uses 'default_statement' and is not captured)
    (case_statement)   @case

    ;; ternary operator
    (conditional_expression) @condop

    ;; short-circuit boolean ops (anywhere in the body)
    (binary_expression operator: "&&") @boolop
    (binary_expression operator: "||") @boolop

    ;; each catch increases complexity (common convention)
    (catch_clause)     @decision
  `
)

func (a *CppAnalyzer) Analyze(sourceCode []byte) (analyzer.AnalysisResult, error) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tscpp.Language())
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
		var nameUnqualified string
		var nameQualified string
		var funcNode *sitter.Node

		for _, cap := range m.Captures {
			captureName := q.CaptureNames()[cap.Index]
			switch captureName {
			case "name":
				nameUnqualified = cap.Node.Utf8Text(sourceCode)
				if funcNode == nil {
					funcNode = ascendToFunctionDefinition(&cap.Node)
				}
			case "qualified":
				nameQualified = cap.Node.Utf8Text(sourceCode)
				if funcNode == nil {
					funcNode = ascendToFunctionDefinition(&cap.Node)
				}
			}
		}

		if funcNode == nil {
			continue
		}

		bodyNode := funcNode.ChildByFieldName("body")
		complexity := calculateComplexity(lang, sourceCode, bodyNode)

		start := funcNode.StartPosition().Row + 1
		end := funcNode.EndPosition().Row + 1

		// Choose qualified name when available.
		name := nameUnqualified
		if strings.TrimSpace(nameQualified) != "" {
			name = nameQualified
		}

		// Append the parameter list text to form a signature-like label.
		// function_definition -> declarator:function_declarator -> parameters:parameter_list
		if decl := funcNode.ChildByFieldName("declarator"); decl != nil && decl.Kind() == "function_declarator" {
			if params := decl.ChildByFieldName("parameters"); params != nil {
				name = name + params.Utf8Text(sourceCode) // e.g. "ns::sum(int a, int b)"
			}
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
		slog.Warn("Error compiling C++ complexity query", "error", qerr)
		return -1
	}
	defer q.Close()

	qc := sitter.NewQueryCursor()
	defer qc.Close()

	matches := qc.Matches(q, bodyNode, src)

	for m := matches.Next(); m != nil; m = matches.Next() {
		for range m.Captures {
			complexity++
		}
	}
	return complexity
}

func ascendToFunctionDefinition(n *sitter.Node) *sitter.Node {
	for node := n; node != nil; node = node.Parent() {
		if node.Kind() == "function_definition" {
			return node
		}
	}
	return nil
}
