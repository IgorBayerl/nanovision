package lang_go

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/IgorBayerl/AdlerCov/internal/language"
	"github.com/IgorBayerl/AdlerCov/internal/model"
	"github.com/fzipp/gocyclo"
)

type GoProcessor struct{}

func NewGoProcessor() language.Processor {
	return &GoProcessor{}
}

func (p *GoProcessor) Name() string {
	return "Go"
}

func (p *GoProcessor) Detect(filePath string) bool {
	return strings.HasSuffix(strings.ToLower(filePath), ".go")
}

func (p *GoProcessor) AnalyzeFile(filePath string, sourceLines []string) ([]model.MethodMetrics, error) {
	// First, parse the file to ensure it's syntactically correct.
	fset := token.NewFileSet()
	src := strings.Join(sourceLines, "\n")
	node, err := parser.ParseFile(fset, filePath, src, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	// Only if parsing succeeds, run the cyclomatic complexity analysis.
	cycloStats := gocyclo.Analyze([]string{filePath}, nil)
	complexityMap := make(map[string]int)
	for _, stat := range cycloStats {
		complexityMap[stat.FuncName] = stat.Complexity
	}

	methods := make([]model.MethodMetrics, 0)
	ast.Inspect(node, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok {
			return true // Continue traversing
		}

		startPos := fset.Position(fn.Pos())
		endPos := fset.Position(fn.End())

		funcName := fn.Name.Name
		if fn.Recv != nil && len(fn.Recv.List) > 0 {
			var typeName string
			switch T := fn.Recv.List[0].Type.(type) {
			case *ast.StarExpr: // Pointer receiver: func (s *MyStruct)
				if ident, ok := T.X.(*ast.Ident); ok {
					typeName = "(*" + ident.Name + ")"
				}
			case *ast.Ident: // Value receiver: func (s MyStruct)
				// FIX: Correctly handle *ast.Ident. T is the identifier itself.
				typeName = "(" + T.Name + ")"
			}
			if typeName != "" {
				funcName = typeName + "." + funcName
			}
		}

		methods = append(methods, model.MethodMetrics{
			Name:                 funcName,
			StartLine:            startPos.Line,
			EndLine:              endPos.Line,
			CyclomaticComplexity: complexityMap[funcName],
		})

		return false
	})

	return methods, nil
}
