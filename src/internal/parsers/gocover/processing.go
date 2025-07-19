package gocover

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log/slog"
	"math"
	"path/filepath"
	"strings"

	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/filereader"
	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/language"
	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/model"
	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/parsers"
	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/utils"
)

// processingOrchestrator holds dependencies and state for a single parsing operation.
type processingOrchestrator struct {
	fileReader   filereader.Reader
	config       parsers.ParserConfig
	assemblyName string
	logger       *slog.Logger
}

// parsedMethod is a temporary struct to hold data from AST (Abstract System Tree) parsing.
type parsedMethod struct {
	DisplayName string
	FuncName    string
	StartLine   int
	EndLine     int
}

func newProcessingOrchestrator(fileReader filereader.Reader, config parsers.ParserConfig, logger *slog.Logger) *processingOrchestrator {
	return &processingOrchestrator{
		fileReader: fileReader,
		config:     config,
		logger:     logger,
	}
}

func (o *processingOrchestrator) processBlocks(blocks []GoCoverProfileBlock) ([]model.Assembly, error) {
	if len(blocks) == 0 {
		return []model.Assembly{}, nil
	}

	var foundAssemblyName string
	if len(blocks) > 0 {
		startPath := blocks[0].FileName
		resolvedStartPath, err := utils.FindFileInSourceDirs(startPath, o.config.SourceDirectories(), o.fileReader)
		if err == nil {
			startPath = resolvedStartPath
		}

		modName, err := o.findModuleNameFromGoMod(startPath)
		if err == nil {
			foundAssemblyName = modName
			o.logger.Info("Discovered Go module name for assembly", "name", foundAssemblyName)
		} else {
			o.logger.Warn("Could not discover Go module name, falling back to default.", "error", err)
		}
	}
	if foundAssemblyName != "" {
		o.assemblyName = foundAssemblyName
	} else {
		o.assemblyName = o.config.Settings().DefaultAssemblyName
	}

	filesByPackage := o.groupFilesByPackage(blocks)
	if len(filesByPackage) == 0 {
		return []model.Assembly{}, nil
	}

	assembly := &model.Assembly{
		Name:    o.assemblyName,
		Classes: []model.Class{},
	}

	for pkgPath, fileBlocks := range filesByPackage {
		class := o.processPackage(pkgPath, fileBlocks)
		if class != nil {
			assembly.Classes = append(assembly.Classes, *class)
		}
	}

	o.aggregateAssemblyMetrics(assembly)
	return []model.Assembly{*assembly}, nil
}

func (o *processingOrchestrator) findModuleNameFromGoMod(startPath string) (string, error) {
	dir := filepath.Dir(startPath)
	var goModPath string

	for {
		potentialPath := filepath.Join(dir, "go.mod")
		if _, err := o.fileReader.Stat(potentialPath); err == nil {
			goModPath = potentialPath
			break
		}

		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			return "", fmt.Errorf("go.mod not found in parent directories of %s", startPath)
		}
		dir = parentDir
	}

	lines, err := o.fileReader.ReadFile(goModPath)
	if err != nil || len(lines) == 0 {
		return "", fmt.Errorf("could not read or empty go.mod at %s: %w", goModPath, err)
	}

	for _, line := range lines {
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module ")), nil
		}
	}

	return "", fmt.Errorf("'module' directive not found in %s", goModPath)
}

func (o *processingOrchestrator) groupFilesByPackage(blocks []GoCoverProfileBlock) map[string]map[string][]GoCoverProfileBlock {
	filesByPackage := make(map[string]map[string][]GoCoverProfileBlock)
	for _, block := range blocks {
		if !o.config.FileFilters().IsElementIncludedInReport(block.FileName) {
			continue
		}
		pkgPath := filepath.ToSlash(filepath.Dir(block.FileName))
		if pkgPath == "." {
			pkgPath = o.assemblyName
		}
		if _, ok := filesByPackage[pkgPath]; !ok {
			filesByPackage[pkgPath] = make(map[string][]GoCoverProfileBlock)
		}
		filesByPackage[pkgPath][block.FileName] = append(filesByPackage[pkgPath][block.FileName], block)
	}
	return filesByPackage
}

func (o *processingOrchestrator) processPackage(pkgPath string, fileBlocks map[string][]GoCoverProfileBlock) *model.Class {
	if !o.config.ClassFilters().IsElementIncludedInReport(pkgPath) {
		return nil
	}
	displayName := pkgPath
	prefix := o.assemblyName + "/"
	if strings.HasPrefix(pkgPath, prefix) {
		displayName = strings.TrimPrefix(pkgPath, prefix)
	} else if pkgPath == o.assemblyName {
		displayName = "(root)"
	}

	packageClass := &model.Class{
		Name:        pkgPath,
		DisplayName: displayName,
		Files:       []model.CodeFile{},
		Methods:     []model.Method{},
		Metrics:     make(map[string]float64),
	}

	for filePath, blocksForFile := range fileBlocks {
		codeFile, methods := o.processFile(filePath, blocksForFile)
		if codeFile == nil {
			continue
		}
		packageClass.Files = append(packageClass.Files, *codeFile)
		packageClass.Methods = append(packageClass.Methods, methods...)
	}

	if len(packageClass.Files) > 0 {
		o.aggregateClassMetrics(packageClass)
		return packageClass
	}
	return nil
}

func (o *processingOrchestrator) processFile(filePath string, blocks []GoCoverProfileBlock) (*model.CodeFile, []model.Method) {
	resolvedPath, err := utils.FindFileInSourceDirs(filePath, o.config.SourceDirectories(), o.fileReader)
	if err != nil {
		o.logger.Warn("Source file not found, line content will be missing.", "file", filePath, "error", err)
		resolvedPath = filePath
	}

	sourceLines, _ := o.fileReader.ReadFile(resolvedPath)
	totalLines, _ := o.fileReader.CountLines(resolvedPath)
	if len(sourceLines) == 0 {
		o.logger.Warn("Source file is empty or could not be read.", "file", resolvedPath)
		return nil, nil
	}

	parsedMethods, err := parseGoSourceForFunctions(resolvedPath, sourceLines)
	if err != nil {
		o.logger.Warn("Failed to parse Go source for functions, method metrics will be unavailable.", "file", resolvedPath, "error", err)
	}

	langProcessor := o.config.LanguageProcessorFactory().FindProcessorForFile(filePath)
	complexityMetrics, err := langProcessor.CalculateCyclomaticComplexity(resolvedPath)
	if err != nil && !errors.Is(err, language.ErrNotSupported) {
		o.logger.Warn("Failed to calculate cyclomatic complexity", "file", resolvedPath, "error", err)
	}

	complexityMap := make(map[string]model.MethodMetric)
	for _, m := range complexityMetrics {
		complexityMap[m.Name] = m
	}

	blocksByMethod := make(map[string][]GoCoverProfileBlock)
	for _, block := range blocks {
		for _, pMethod := range parsedMethods {
			if block.StartLine >= pMethod.StartLine && block.StartLine <= pMethod.EndLine {
				blocksByMethod[pMethod.DisplayName] = append(blocksByMethod[pMethod.DisplayName], block)
				break
			}
		}
	}

	var methods []model.Method
	var codeElements []model.CodeElement
	for _, pMethod := range parsedMethods {
		methodBlocks := blocksByMethod[pMethod.DisplayName]
		totalStatements := 0
		coveredStatements := 0

		for _, block := range methodBlocks {
			totalStatements += block.NumStatements
			if block.HitCount > 0 {
				coveredStatements += block.NumStatements
			}
		}

		var lineRate float64
		if totalStatements > 0 {
			lineRate = float64(coveredStatements) / float64(totalStatements)
		}

		method := model.Method{
			Name:        pMethod.FuncName,
			DisplayName: pMethod.DisplayName,
			FirstLine:   pMethod.StartLine,
			LastLine:    pMethod.EndLine,
			LineRate:    lineRate,
			Complexity:  math.NaN(),
		}

		if metric, ok := complexityMap[method.DisplayName]; ok {
			if len(metric.Metrics) > 0 {
				method.Complexity = metric.Metrics[0].Value.(float64)
			}
		}

		o.populateStandardGoMethodMetrics(&method)
		methods = append(methods, method)

		lineRateForQuota := method.LineRate * 100
		codeElements = append(codeElements, model.CodeElement{
			Name:          utils.GetShortMethodName(method.DisplayName),
			FullName:      method.DisplayName,
			Type:          langProcessor.CategorizeCodeElement(&method),
			FirstLine:     method.FirstLine,
			LastLine:      method.LastLine,
			CoverageQuota: &lineRateForQuota,
		})
	}

	lineData := make(map[int]int)
	for _, block := range blocks {
		for l := block.StartLine; l <= block.EndLine; l++ {
			if _, ok := lineData[l]; !ok || block.HitCount > lineData[l] {
				lineData[l] = block.HitCount
			}
		}
	}

	var finalLines []model.Line
	fileTotalCoveredStatements := 0
	fileTotalCoverableStatements := 0
	for _, block := range blocks {
		fileTotalCoverableStatements += block.NumStatements
		if block.HitCount > 0 {
			fileTotalCoveredStatements += block.NumStatements
		}
	}

	for i, lineContent := range sourceLines {
		lineNumber := i + 1
		hits, isBlockMember := lineData[lineNumber]
		line := model.Line{Number: lineNumber, Content: lineContent, Hits: -1}

		if isBlockMember {
			line.Hits = hits
			line.LineVisitStatus = model.NotCovered
			if hits > 0 {
				line.LineVisitStatus = model.Covered
			}
		} else {
			line.LineVisitStatus = model.NotCoverable
		}

		if strings.TrimSpace(lineContent) == "}" {
			line.LineVisitStatus = model.NotCoverable
			line.Hits = -1
		}

		finalLines = append(finalLines, line)
	}

	var methodMetricsForFile []model.MethodMetric
	for _, method := range methods {
		if method.MethodMetrics != nil {
			methodMetricsForFile = append(methodMetricsForFile, method.MethodMetrics...)
		}
	}

	codeFile := &model.CodeFile{
		Path:           resolvedPath,
		Lines:          finalLines,
		CoveredLines:   fileTotalCoveredStatements,
		CoverableLines: fileTotalCoverableStatements,
		TotalLines:     totalLines,
		CodeElements:   codeElements,
		MethodMetrics:  methodMetricsForFile,
	}

	return codeFile, methods
}

func (o *processingOrchestrator) populateStandardGoMethodMetrics(method *model.Method) {
	method.MethodMetrics = []model.MethodMetric{}
	shortMetricName := utils.GetShortMethodName(method.DisplayName)

	if !math.IsNaN(method.Complexity) {
		method.MethodMetrics = append(method.MethodMetrics, model.MethodMetric{
			Name: shortMetricName, Line: method.FirstLine,
			Metrics: []model.Metric{{Name: "Cyclomatic complexity", Value: method.Complexity, Status: model.StatusOk}},
		})
	}
	lineCoveragePercentage := method.LineRate * 100.0
	if !math.IsNaN(lineCoveragePercentage) {
		method.MethodMetrics = append(method.MethodMetrics, model.MethodMetric{
			Name: shortMetricName, Line: method.FirstLine,
			Metrics: []model.Metric{{Name: "Line coverage", Value: lineCoveragePercentage, Status: model.StatusOk}},
		})
	}
	if !math.IsNaN(method.Complexity) {
		crapScoreValue := o.calculateCrapScore(method.LineRate, method.Complexity)
		if !math.IsNaN(crapScoreValue) {
			method.MethodMetrics = append(method.MethodMetrics, model.MethodMetric{
				Name: shortMetricName, Line: method.FirstLine,
				Metrics: []model.Metric{{Name: "CrapScore", Value: crapScoreValue, Status: model.StatusOk}},
			})
		}
	}
}

func (o *processingOrchestrator) calculateCrapScore(coverage float64, complexity float64) float64 {
	if math.IsNaN(coverage) || math.IsInf(coverage, 0) || coverage < 0 || coverage > 1 {
		coverage = 0
	}
	if math.IsNaN(complexity) || math.IsInf(complexity, 0) || complexity < 0 {
		return math.NaN()
	}
	uncoveredRatio := 1.0 - coverage
	return (math.Pow(complexity, 2) * math.Pow(uncoveredRatio, 3)) + complexity
}

func (o *processingOrchestrator) aggregateClassMetrics(class *model.Class) {
	for _, f := range class.Files {
		class.LinesCovered += f.CoveredLines
		class.LinesValid += f.CoverableLines
		class.TotalLines += f.TotalLines
	}
	class.TotalMethods = len(class.Methods)
	for _, method := range class.Methods {
		if !math.IsNaN(method.LineRate) {
			if method.LineRate > 0 {
				class.CoveredMethods++
			}
			if method.LineRate >= 1.0 {
				class.FullyCoveredMethods++
			}
		}
		if !math.IsNaN(method.Complexity) {
			class.Metrics["Cyclomatic complexity"] += method.Complexity
		}
	}
}

func (o *processingOrchestrator) aggregateAssemblyMetrics(assembly *model.Assembly) {
	for _, cls := range assembly.Classes {
		assembly.LinesCovered += cls.LinesCovered
		assembly.LinesValid += cls.LinesValid
		assembly.TotalLines += cls.TotalLines
	}
}

func parseGoSourceForFunctions(filePath string, sourceLines []string) ([]parsedMethod, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, strings.Join(sourceLines, "\n"), 0)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Go source: %w", err)
	}

	var methods []parsedMethod
	ast.Inspect(f, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok {
			funcName := fn.Name.Name
			displayName := funcName

			if fn.Recv != nil && len(fn.Recv.List) > 0 {
				typeExpr := fn.Recv.List[0].Type
				var receiverTypeNameBuilder strings.Builder
				var extractTypeName func(ast.Expr)
				extractTypeName = func(e ast.Expr) {
					switch t := e.(type) {
					case *ast.StarExpr:
						receiverTypeNameBuilder.WriteString("*")
						extractTypeName(t.X)
					case *ast.Ident:
						receiverTypeNameBuilder.WriteString(t.Name)
					}
				}
				extractTypeName(typeExpr)

				receiverTypeName := receiverTypeNameBuilder.String()
				if receiverTypeName != "" {
					displayName = fmt.Sprintf("(%s).%s", receiverTypeName, funcName)
				}
			}

			startPosition := fset.Position(fn.Pos())
			endPosition := fset.Position(fn.End())
			methods = append(methods, parsedMethod{
				DisplayName: displayName,
				FuncName:    funcName,
				StartLine:   startPosition.Line,
				EndLine:     endPosition.Line,
			})
		}
		return true
	})

	return methods, nil
}
