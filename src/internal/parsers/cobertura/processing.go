package cobertura

import (
	"errors"
	"fmt"
	"log/slog"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/filereader"
	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/language"
	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/model"
	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/parsers"
	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/utils"
)

var (
	conditionCoverageRegexCobertura = regexp.MustCompile(`\((?P<NumberOfCoveredBranches>\d+)/(?P<NumberOfTotalBranches>\d+)\)$`)
)

type fileProcessingMetrics struct {
	linesCovered    int
	linesValid      int
	branchesCovered int
	branchesValid   int
}

type processingOrchestrator struct {
	fileReader                        filereader.Reader
	config                            parsers.ParserConfig
	sourceDirs                        []string
	uniqueFilePathsForGrandTotalLines map[string]int
	processedAssemblyFiles            map[string]struct{}
	detectedBranchCoverage            bool
	logger                            *slog.Logger
}

func newProcessingOrchestrator(
	fileReader filereader.Reader,
	config parsers.ParserConfig,
	sourceDirs []string,
	logger *slog.Logger,
) *processingOrchestrator {
	return &processingOrchestrator{
		fileReader:                        fileReader,
		config:                            config,
		sourceDirs:                        sourceDirs,
		uniqueFilePathsForGrandTotalLines: make(map[string]int),
		detectedBranchCoverage:            false,
		logger:                            logger,
	}
}

func (o *processingOrchestrator) processPackages(packages []PackageXML) ([]model.Assembly, bool, error) {
	var parsedAssemblies []model.Assembly
	for _, pkgXML := range packages {
		assembly, err := o.processPackage(pkgXML)
		if err != nil {
			o.logger.Warn("Could not process Cobertura package, skipping.", "package", pkgXML.Name, "error", err)
			continue
		}
		if assembly != nil {
			parsedAssemblies = append(parsedAssemblies, *assembly)
		}
	}
	return parsedAssemblies, o.detectedBranchCoverage, nil
}

func (o *processingOrchestrator) processPackage(pkgXML PackageXML) (*model.Assembly, error) {
	if !o.config.AssemblyFilters().IsElementIncludedInReport(pkgXML.Name) {
		o.logger.Debug("Skipping assembly excluded by filter", "assembly", pkgXML.Name)
		return nil, nil
	}

	assembly := &model.Assembly{
		Name:    pkgXML.Name,
		Classes: []model.Class{},
	}
	o.processedAssemblyFiles = make(map[string]struct{})

	classesXMLGrouped := o.groupClassesByLogicalName(pkgXML.Classes.Class)

	for logicalName, classXMLGroup := range classesXMLGrouped {
		classModel, err := o.processClassGroup(logicalName, classXMLGroup)
		if err != nil {
			o.logger.Debug("Skipping class group.", "class", logicalName, "reason", err)
			continue
		}
		if classModel != nil {
			assembly.Classes = append(assembly.Classes, *classModel)
		}
	}

	o.aggregateAssemblyMetrics(assembly)
	return assembly, nil
}

func (o *processingOrchestrator) groupClassesByLogicalName(classes []ClassXML) map[string][]ClassXML {
	grouped := make(map[string][]ClassXML)
	for _, classXML := range classes {
		formatter := o.config.LanguageProcessorFactory().FindProcessorForFile(classXML.Filename)
		logicalName := formatter.GetLogicalClassName(classXML.Name)
		grouped[logicalName] = append(grouped[logicalName], classXML)
	}
	return grouped
}

func (o *processingOrchestrator) processClassGroup(logicalClassName string, classXMLs []ClassXML) (*model.Class, error) {
	if len(classXMLs) == 0 {
		return nil, nil
	}
	primaryFormatter := o.config.LanguageProcessorFactory().FindProcessorForFile(classXMLs[0].Filename)

	if !o.config.ClassFilters().IsElementIncludedInReport(logicalClassName) {
		return nil, fmt.Errorf("class '%s' is excluded by filters", logicalClassName)
	}

	classModel := &model.Class{
		Name:    logicalClassName,
		Files:   []model.CodeFile{},
		Methods: []model.Method{},
		Metrics: make(map[string]float64),
	}

	if primaryFormatter.IsCompilerGeneratedClass(classModel) {
		return nil, fmt.Errorf("class '%s' is a compiler-generated type and was filtered out", logicalClassName)
	}

	classModel.DisplayName = primaryFormatter.FormatClassName(classModel)

	classProcessedFilePaths := make(map[string]struct{})
	xmlFragmentsByFile := o.groupClassFragmentsByFile(classXMLs)

	for filePath, fragmentsForFile := range xmlFragmentsByFile {
		fileFormatter := o.config.LanguageProcessorFactory().FindProcessorForFile(filePath)
		codeFile, methodsInFile, err := o.processFileForClass(filePath, classModel, fragmentsForFile, fileFormatter)
		if err != nil {
			o.logger.Warn("Failed to process file for class, skipping file.", "file", filePath, "class", classModel.DisplayName, "error", err)
			continue
		}

		classModel.Files = append(classModel.Files, *codeFile)
		classModel.Methods = append(classModel.Methods, methodsInFile...)

		o.processedAssemblyFiles[codeFile.Path] = struct{}{}
		classProcessedFilePaths[codeFile.Path] = struct{}{}
	}

	o.aggregateClassMetrics(classModel, classProcessedFilePaths)
	return classModel, nil
}

func (o *processingOrchestrator) processFileForClass(filePath string, classModel *model.Class, fragments []ClassXML, fileFormatter language.Processor) (*model.CodeFile, []model.Method, error) {
	resolvedPath, err := utils.FindFileInSourceDirs(filePath, o.sourceDirs, o.fileReader)
	if err != nil {
		o.logger.Warn("Source file not found, line content will be missing.", "file", filePath, "class", classModel.DisplayName)
		resolvedPath = filePath
	}

	complexityMetrics, err := fileFormatter.CalculateCyclomaticComplexity(resolvedPath)
	if err != nil && !errors.Is(err, language.ErrNotSupported) {
		o.logger.Warn("Failed to calculate cyclomatic complexity", "file", resolvedPath, "error", err)
	}
	complexityMap := make(map[string]model.MethodMetric)
	for _, m := range complexityMetrics {
		complexityMap[m.Name] = m
	}
	// =================================================================

	sourceLines, _ := o.fileReader.ReadFile(resolvedPath)
	totalLines := o.getTotalLines(resolvedPath, sourceLines)
	maxLineNumInFile := getMaxLineNumber(fragments)
	mergedLineHits, mergedBranches := o.mergeLineAndBranchData(fragments)

	// Pass the complexity map down to the method processor
	methodsInFile, codeElementsInFile, err := o.processMethodsForFile(fragments, classModel, fileFormatter, complexityMap)
	if err != nil {
		return nil, nil, fmt.Errorf("processing methods for file %s: %w", filePath, err)
	}

	finalLinesForFile, fileMetrics := o.assembleLinesForFile(maxLineNumInFile, sourceLines, mergedLineHits, mergedBranches)

	codeFile := &model.CodeFile{
		Path:           resolvedPath,
		Lines:          finalLinesForFile,
		CoveredLines:   fileMetrics.linesCovered,
		CoverableLines: fileMetrics.linesValid,
		TotalLines:     totalLines,
		CodeElements:   codeElementsInFile,
	}

	for _, method := range methodsInFile {
		if method.MethodMetrics != nil {
			codeFile.MethodMetrics = append(codeFile.MethodMetrics, method.MethodMetrics...)
		}
	}
	codeFile.MethodMetrics = utils.DistinctBy(codeFile.MethodMetrics, func(mm model.MethodMetric) string { return mm.Name + fmt.Sprintf("_%d", mm.Line) })

	return codeFile, methodsInFile, nil
}

func (o *processingOrchestrator) processMethodsForFile(fragments []ClassXML, classModel *model.Class, fileFormatter language.Processor, complexityMap map[string]model.MethodMetric) ([]model.Method, []model.CodeElement, error) {
	var allMethods []model.Method

	for _, fragment := range fragments {
		for _, methodXML := range fragment.Methods.Method {
			methodModel := o.processMethodXML(methodXML, classModel, fileFormatter, complexityMap)
			allMethods = append(allMethods, *methodModel)
		}
	}

	distinctMethods := utils.DistinctBy(allMethods, func(m model.Method) string {
		return m.Name + m.Signature
	})

	var allCodeElements []model.CodeElement
	for i := range distinctMethods {
		allCodeElements = append(allCodeElements, o.createCodeElementFromMethod(&distinctMethods[i], fileFormatter))
	}

	utils.SortByLineAndName(distinctMethods)
	utils.SortByLineAndName(allCodeElements)

	return distinctMethods, allCodeElements, nil
}

func (o *processingOrchestrator) processMethodXML(methodXML MethodXML, classModel *model.Class, fileFormatter language.Processor, complexityMap map[string]model.MethodMetric) *model.Method {
	method := &model.Method{
		Name:       methodXML.Name,
		Signature:  methodXML.Signature,
		Complexity: parseFloat(methodXML.Complexity),
	}

	method.DisplayName = fileFormatter.FormatMethodName(method, classModel)

	if metric, ok := complexityMap[method.DisplayName]; ok {
		if len(metric.Metrics) > 0 {
			// Override the complexity from the Cobertura file with our more accurate one.
			method.Complexity = metric.Metrics[0].Value.(float64)
		}
	}

	o.processMethodLines(methodXML, method)
	o.populateStandardMethodMetrics(method)

	return method
}

func (o *processingOrchestrator) createCodeElementFromMethod(method *model.Method, fileFormatter language.Processor) model.CodeElement {
	elementType := fileFormatter.CategorizeCodeElement(method)

	var coverageQuota *float64
	if len(method.Lines) > 0 && !math.IsNaN(method.LineRate) {
		cq := method.LineRate * 100.0
		coverageQuota = &cq
	}

	shortName := method.DisplayName
	if elementType == model.MethodElementType {
		shortName = utils.GetShortMethodName(method.DisplayName)
	}

	return model.CodeElement{
		Name:          shortName,
		FullName:      method.DisplayName,
		Type:          elementType,
		FirstLine:     method.FirstLine,
		LastLine:      method.LastLine,
		CoverageQuota: coverageQuota,
	}
}

// The rest of this file (helper functions) remains unchanged as they are not
// directly involved in calculating or setting the cyclomatic complexity.

func (o *processingOrchestrator) processMethodLines(methodXML MethodXML, method *model.Method) {
	minLine, maxLine := math.MaxInt32, 0
	var methodLinesCovered, methodLinesValid int
	var methodBranchesCovered, methodBranchesValid int

	for _, lineXML := range methodXML.Lines.Line {
		currentLineNum, _ := strconv.Atoi(lineXML.Number)
		if currentLineNum < minLine {
			minLine = currentLineNum
		}
		if currentLineNum > maxLine {
			maxLine = currentLineNum
		}

		lineModel, lineMetricsStats := o.processLineXML(lineXML)
		method.Lines = append(method.Lines, lineModel)

		if lineModel.Hits >= 0 {
			methodLinesValid++
			if lineModel.Hits > 0 {
				methodLinesCovered++
			}
		}
		methodBranchesCovered += lineMetricsStats.branchesCovered
		methodBranchesValid += lineMetricsStats.branchesValid
	}

	method.FirstLine = 0
	if minLine != math.MaxInt32 {
		method.FirstLine = minLine
	}
	method.LastLine = maxLine

	if methodLinesValid > 0 {
		method.LineRate = float64(methodLinesCovered) / float64(methodLinesValid)
	} else {
		method.LineRate = 0.0
	}

	if !o.detectedBranchCoverage {
		method.BranchRate = nil
	} else if methodBranchesValid > 0 {
		rate := float64(methodBranchesCovered) / float64(methodBranchesValid)
		method.BranchRate = &rate
	} else {
		rate := 1.0
		method.BranchRate = &rate
	}
}

func (o *processingOrchestrator) processLineXML(lineXML LineXML) (model.Line, fileProcessingMetrics) {
	metrics := fileProcessingMetrics{}
	lineNumber, _ := strconv.Atoi(lineXML.Number)
	isBranchPoint := strings.EqualFold(lineXML.Branch, "true")

	if isBranchPoint && !o.detectedBranchCoverage {
		o.detectedBranchCoverage = true
	}

	line := model.Line{
		Number:        lineNumber,
		Hits:          parseInt(lineXML.Hits),
		IsBranchPoint: isBranchPoint,
		Branch:        make([]model.BranchCoverageDetail, 0),
	}

	if line.IsBranchPoint {
		matches := conditionCoverageRegexCobertura.FindStringSubmatch(lineXML.ConditionCoverage)
		if len(matches) > 0 {
			coveredStr := findNamedGroup(conditionCoverageRegexCobertura, matches, "NumberOfCoveredBranches")
			totalStr := findNamedGroup(conditionCoverageRegexCobertura, matches, "NumberOfTotalBranches")

			if coveredStr != "" && totalStr != "" {
				numberOfCoveredBranches, _ := strconv.Atoi(coveredStr)
				numberOfTotalBranches, _ := strconv.Atoi(totalStr)

				if numberOfTotalBranches > 0 {
					line.CoveredBranches = numberOfCoveredBranches
					line.TotalBranches = numberOfTotalBranches
					for i := 0; i < line.TotalBranches; i++ {
						visits := 0
						if i < line.CoveredBranches {
							visits = 1
						}
						line.Branch = append(line.Branch, model.BranchCoverageDetail{Identifier: fmt.Sprintf("%d_%d", lineNumber, i), Visits: visits})
					}
				}
			}
		} else if len(lineXML.Conditions.Condition) > 0 {
			for _, condition := range lineXML.Conditions.Condition {
				visits := 0
				if strings.HasPrefix(condition.Coverage, "100") {
					visits = 1
					line.CoveredBranches++
				}
				line.Branch = append(line.Branch, model.BranchCoverageDetail{Identifier: condition.Number, Visits: visits})
				line.TotalBranches++
			}
		} else {
			o.setFallbackBranchData(&line)
		}
	}

	metrics.branchesCovered = line.CoveredBranches
	metrics.branchesValid = line.TotalBranches
	return line, metrics
}

func (o *processingOrchestrator) setFallbackBranchData(line *model.Line) {
	if line.Hits > 0 {
		line.CoveredBranches = 1
	} else {
		line.CoveredBranches = 0
	}
	line.TotalBranches = 1
	line.Branch = append(line.Branch, model.BranchCoverageDetail{
		Identifier: fmt.Sprintf("%d_0", line.Number),
		Visits:     line.CoveredBranches,
	})
}

func (o *processingOrchestrator) populateStandardMethodMetrics(method *model.Method) {
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

	if method.BranchRate != nil {
		branchCoveragePercentage := *method.BranchRate * 100.0
		if !math.IsNaN(branchCoveragePercentage) {
			method.MethodMetrics = append(method.MethodMetrics, model.MethodMetric{
				Name: shortMetricName, Line: method.FirstLine,
				Metrics: []model.Metric{{Name: "Branch coverage", Value: branchCoveragePercentage, Status: model.StatusOk}},
			})
		}
	}

	var coverageForCrapScore float64
	if method.BranchRate != nil {
		coverageForCrapScore = *method.BranchRate
	} else {
		coverageForCrapScore = method.LineRate
	}

	crapScoreValue := o.calculateCrapScore(coverageForCrapScore, method.Complexity)
	if !math.IsNaN(crapScoreValue) {
		method.MethodMetrics = append(method.MethodMetrics, model.MethodMetric{
			Name: shortMetricName, Line: method.FirstLine,
			Metrics: []model.Metric{{Name: "CrapScore", Value: crapScoreValue, Status: model.StatusOk}},
		})
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

func (o *processingOrchestrator) groupClassFragmentsByFile(classXMLs []ClassXML) map[string][]ClassXML {
	grouped := make(map[string][]ClassXML)
	for _, classXML := range classXMLs {
		if classXML.Filename == "" || !o.config.FileFilters().IsElementIncludedInReport(classXML.Filename) {
			continue
		}
		grouped[classXML.Filename] = append(grouped[classXML.Filename], classXML)
	}
	return grouped
}

func (o *processingOrchestrator) aggregateAssemblyMetrics(assembly *model.Assembly) {
	var linesCovered, linesValid, branchesCovered, branchesValid, totalLines int
	hasBranchData := false

	for _, cls := range assembly.Classes {
		linesCovered += cls.LinesCovered
		linesValid += cls.LinesValid
		if cls.BranchesCovered != nil && cls.BranchesValid != nil {
			hasBranchData = true
			branchesCovered += *cls.BranchesCovered
			branchesValid += *cls.BranchesValid
		}
	}
	for path := range o.processedAssemblyFiles {
		if lineCount, ok := o.uniqueFilePathsForGrandTotalLines[path]; ok {
			totalLines += lineCount
		}
	}
	assembly.LinesCovered = linesCovered
	assembly.LinesValid = linesValid
	assembly.TotalLines = totalLines
	if hasBranchData {
		assembly.BranchesCovered = &branchesCovered
		assembly.BranchesValid = &branchesValid
	}
}

func (o *processingOrchestrator) aggregateClassMetrics(class *model.Class, processedFiles map[string]struct{}) {
	var totalClassLines, totalClassBranchesCovered, totalClassBranchesValid int
	var coveredM, fullyCoveredM, totalM int
	hasClassBranchData := false

	for _, f := range class.Files {
		class.LinesCovered += f.CoveredLines
		class.LinesValid += f.CoverableLines
		for _, line := range f.Lines {
			if line.IsBranchPoint {
				hasClassBranchData = true
				totalClassBranchesValid += line.TotalBranches
				totalClassBranchesCovered += line.CoveredBranches
			}
		}
	}
	if hasClassBranchData {
		class.BranchesCovered = &totalClassBranchesCovered
		class.BranchesValid = &totalClassBranchesValid
	}

	for path := range processedFiles {
		if lineCount, ok := o.uniqueFilePathsForGrandTotalLines[path]; ok {
			totalClassLines += lineCount
		}
	}
	class.TotalLines = totalClassLines

	if len(class.Methods) > 0 {
		totalM = len(class.Methods)
		for _, method := range class.Methods {
			atLeastOneLineCoveredInMethod := false
			methodIsFullyCovered := true
			methodHasCoverableLines := false
			for _, line := range method.Lines {
				if line.Hits >= 0 {
					methodHasCoverableLines = true
					if line.Hits > 0 {
						atLeastOneLineCoveredInMethod = true
					} else {
						methodIsFullyCovered = false
					}
				}
			}
			if atLeastOneLineCoveredInMethod {
				coveredM++
			}
			if methodHasCoverableLines && methodIsFullyCovered {
				fullyCoveredM++
			} else if !methodHasCoverableLines && len(method.Lines) == 0 {
				fullyCoveredM++
			}
		}
	}
	class.CoveredMethods = coveredM
	class.FullyCoveredMethods = fullyCoveredM
	class.TotalMethods = totalM

	for _, method := range class.Methods {
		if !math.IsNaN(method.Complexity) {
			class.Metrics["Cyclomatic complexity"] += method.Complexity
		}
	}
}

func (o *processingOrchestrator) getTotalLines(path string, sourceLines []string) int {
	if count, ok := o.uniqueFilePathsForGrandTotalLines[path]; ok {
		return count
	}
	if lineCount, err := o.fileReader.CountLines(path); err == nil {
		o.uniqueFilePathsForGrandTotalLines[path] = lineCount
		return lineCount
	}
	if sourceLines != nil {
		o.uniqueFilePathsForGrandTotalLines[path] = len(sourceLines)
		return len(sourceLines)
	}
	return 0
}

func getMaxLineNumber(fragments []ClassXML) int {
	maxLine := 0
	for _, fragment := range fragments {
		allLines := fragment.Lines.Line
		for _, method := range fragment.Methods.Method {
			allLines = append(allLines, method.Lines.Line...)
		}
		for _, lineXML := range allLines {
			if ln, _ := strconv.Atoi(lineXML.Number); ln > maxLine {
				maxLine = ln
			}
		}
	}
	return maxLine
}

func (o *processingOrchestrator) mergeBranches(existing, new []model.BranchCoverageDetail) []model.BranchCoverageDetail {
	if existing == nil {
		return new
	}
	for _, newBranch := range new {
		found := false
		for i, existingBranch := range existing {
			if existingBranch.Identifier == newBranch.Identifier {
				existing[i].Visits += newBranch.Visits
				found = true
				break
			}
		}
		if !found {
			existing = append(existing, newBranch)
		}
	}
	return existing
}

func (o *processingOrchestrator) mergeLineAndBranchData(fragments []ClassXML) (map[int]int, map[int][]model.BranchCoverageDetail) {
	lineHits := make(map[int]int)
	branchDetails := make(map[int][]model.BranchCoverageDetail)

	for _, fragment := range fragments {
		allLines := make([]LineXML, len(fragment.Lines.Line))
		copy(allLines, fragment.Lines.Line)
		for _, method := range fragment.Methods.Method {
			allLines = append(allLines, method.Lines.Line...)
		}

		for _, lineXML := range allLines {
			lineNumber, err := strconv.Atoi(lineXML.Number)
			if err != nil || lineNumber <= 0 {
				continue
			}

			if hits, err := strconv.Atoi(lineXML.Hits); err == nil {
				lineHits[lineNumber] += hits
			}

			if strings.EqualFold(lineXML.Branch, "true") {
				lineModel, _ := o.processLineXML(lineXML)
				if lineModel.IsBranchPoint {
					currentBranches := branchDetails[lineNumber]
					branchDetails[lineNumber] = o.mergeBranches(currentBranches, lineModel.Branch)
				}
			}
		}
	}
	return lineHits, branchDetails
}

func (o *processingOrchestrator) assembleLinesForFile(maxLineNum int, sourceLines []string, lineHits map[int]int, branches map[int][]model.BranchCoverageDetail) ([]model.Line, fileProcessingMetrics) {
	var finalLines []model.Line
	metrics := fileProcessingMetrics{}

	for lineNum := 1; lineNum <= maxLineNum; lineNum++ {
		lineContent := ""
		if lineNum > 0 && lineNum <= len(sourceLines) {
			lineContent = sourceLines[lineNum-1]
		}

		hits, hasHits := lineHits[lineNum]
		currentLine := model.Line{
			Number:  lineNum,
			Content: lineContent,
			Hits:    -1,
		}
		if hasHits {
			currentLine.Hits = hits
		}

		if branchData, ok := branches[lineNum]; ok {
			currentLine.IsBranchPoint = true
			currentLine.Branch = branchData
			for _, b := range branchData {
				if b.Visits > 0 {
					currentLine.CoveredBranches++
				}
				currentLine.TotalBranches++
			}
		}

		currentLine.LineVisitStatus = determineLineVisitStatus(currentLine.Hits, currentLine.IsBranchPoint, currentLine.CoveredBranches, currentLine.TotalBranches)

		if currentLine.Hits >= 0 {
			metrics.linesValid++
			if currentLine.Hits > 0 {
				metrics.linesCovered++
			}
		}
		metrics.branchesCovered += currentLine.CoveredBranches
		metrics.branchesValid += currentLine.TotalBranches
		finalLines = append(finalLines, currentLine)
	}

	return finalLines, metrics
}

func findNamedGroup(re *regexp.Regexp, match []string, groupName string) string {
	for i, name := range re.SubexpNames() {
		if i > 0 && i < len(match) && name == groupName {
			return match[i]
		}
	}
	return ""
}

func parseInt(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}

func parseFloat(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

func determineLineVisitStatus(hits int, isBranchPoint bool, coveredBranches int, totalBranches int) model.LineVisitStatus {
	if hits < 0 {
		return model.NotCoverable
	}
	if isBranchPoint {
		if totalBranches == 0 {
			return model.NotCoverable
		}
		if coveredBranches == totalBranches {
			return model.Covered
		}
		if coveredBranches > 0 {
			return model.PartiallyCovered
		}
		return model.NotCovered
	}
	if hits > 0 {
		return model.Covered
	}
	return model.NotCovered
}
