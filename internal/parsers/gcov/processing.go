package gcov

import (
	"fmt"
	"log/slog"
	"math"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/IgorBayerl/AdlerCov/internal/filereader"
	"github.com/IgorBayerl/AdlerCov/internal/language"
	"github.com/IgorBayerl/AdlerCov/internal/language/cpp"
	"github.com/IgorBayerl/AdlerCov/internal/model"
	"github.com/IgorBayerl/AdlerCov/internal/parsers"
	"github.com/IgorBayerl/AdlerCov/internal/utils"
)

var (
	lineCoverageRegex   = regexp.MustCompile(`^\s*(?P<Visits>-|#####|=====|\d+):\s*(?P<LineNumber>[1-9]\d*):.*`)
	branchCoverageRegex = regexp.MustCompile(`^branch\s*(?P<Number>\d+)\s*(?:taken\s*(?P<Visits>\d+)|never executed)`)
	functionRegex       = regexp.MustCompile(`^function\s(?P<Name>.*?)\scalled.*`)
)

type processingOrchestrator struct {
	fileReader             filereader.Reader
	config                 parsers.ParserConfig
	logger                 *slog.Logger
	detectedBranchCoverage bool
	unresolvedSourceFiles  []string
}

func newProcessingOrchestrator(fileReader filereader.Reader, config parsers.ParserConfig, logger *slog.Logger) *processingOrchestrator {
	return &processingOrchestrator{
		fileReader:            fileReader,
		config:                config,
		logger:                logger,
		unresolvedSourceFiles: make([]string, 0),
	}
}

func (o *processingOrchestrator) processLines(lines []string) (*model.Assembly, error) {
	if len(lines) == 0 {
		return nil, fmt.Errorf("gcov file is empty")
	}

	firstLine := lines[0]
	if !strings.Contains(firstLine, "0:Source:") {
		return nil, fmt.Errorf("invalid gcov format: first line does not contain '0:Source:'")
	}
	sourceFilePath := strings.TrimSpace(strings.SplitN(firstLine, "0:Source:", 2)[1])

	if !o.config.FileFilters().IsElementIncludedInReport(sourceFilePath) {
		return nil, nil
	}

	resolvedPath, err := utils.FindFileInSourceDirs(sourceFilePath, o.config.SourceDirectories(), o.fileReader)
	if err != nil {
		o.logger.Error("Source file not found for gcov report.", "file", sourceFilePath)
		o.unresolvedSourceFiles = append(o.unresolvedSourceFiles, sourceFilePath)
		return nil, nil
	}

	sourceLines, err := o.fileReader.ReadFile(resolvedPath)
	if err != nil {
		o.logger.Error("Failed to read source file.", "file", resolvedPath, "error", err)
		o.unresolvedSourceFiles = append(o.unresolvedSourceFiles, sourceFilePath)
		return nil, nil
	}

	className := filepath.Base(sourceFilePath)
	if !o.config.ClassFilters().IsElementIncludedInReport(className) {
		return nil, nil
	}

	assembly := &model.Assembly{Name: o.config.Settings().DefaultAssemblyName}
	class := &model.Class{Name: className, DisplayName: className}
	codeFile := &model.CodeFile{Path: resolvedPath, TotalLines: len(sourceLines)}

	langProcessor := o.config.LanguageProcessorFactory().FindProcessorForFile(resolvedPath)
	o.parseCoverageData(lines, class, codeFile, sourceLines, langProcessor)

	class.Files = append(class.Files, *codeFile)
	o.aggregateClassMetrics(class)
	assembly.Classes = append(assembly.Classes, *class)
	o.aggregateAssemblyMetrics(assembly)

	return assembly, nil
}

func (o *processingOrchestrator) parseCoverageData(gcovLines []string, class *model.Class, codeFile *model.CodeFile, sourceLines []string, langProcessor language.Processor) {
	lineAnalyses := make(map[int]*model.Line)
	var functions []model.Method
	var lastLineNumber int
	var pendingFunctionName string

	for _, line := range gcovLines {
		if match := lineCoverageRegex.FindStringSubmatch(line); match != nil {
			visitsText := match[1]
			lineNumber, _ := strconv.Atoi(match[2])
			lastLineNumber = lineNumber

			if pendingFunctionName != "" {
				functions = append(functions, model.Method{
					Name:      pendingFunctionName,
					FirstLine: lineNumber,
				})
				pendingFunctionName = ""
			}

			if visitsText == "-" {
				continue
			}

			lineAnalysis := &model.Line{
				Number:  lineNumber,
				Content: sourceLines[lineNumber-1],
				Hits:    0,
			}
			if visitsText != "#####" && visitsText != "=====" {
				lineAnalysis.Hits, _ = strconv.Atoi(visitsText)
			}
			lineAnalyses[lineNumber] = lineAnalysis

		} else if match := branchCoverageRegex.FindStringSubmatch(line); match != nil {
			o.detectedBranchCoverage = true
			branchNumber, _ := strconv.Atoi(match[1])
			visits := 0
			if len(match) > 2 && match[2] != "" {
				visits, _ = strconv.Atoi(match[2])
			}

			if lineAnalysis, ok := lineAnalyses[lastLineNumber]; ok {
				lineAnalysis.IsBranchPoint = true
				lineAnalysis.Branch = append(lineAnalysis.Branch, model.BranchCoverageDetail{
					Identifier: strconv.Itoa(branchNumber),
					Visits:     visits,
				})
			}
		} else if match := functionRegex.FindStringSubmatch(line); match != nil {
			pendingFunctionName = match[1]
		}
	}

	var coveredLines, coverableLines int
	for _, line := range lineAnalyses {
		if line.IsBranchPoint {
			line.TotalBranches = len(line.Branch)
			for _, b := range line.Branch {
				if b.Visits > 0 {
					line.CoveredBranches++
				}
			}
		}
		line.LineVisitStatus = determineLineVisitStatus(line.Hits, line.IsBranchPoint, line.CoveredBranches, line.TotalBranches)
		codeFile.Lines = append(codeFile.Lines, *line)

		if line.Hits >= 0 {
			coverableLines++
			if line.Hits > 0 {
				coveredLines++
			}
		}
	}
	codeFile.CoveredLines = coveredLines
	codeFile.CoverableLines = coverableLines

	if cppProcessor, ok := langProcessor.(*cpp.CppProcessor); ok {
		cppProcessor.MapCsharpMangledNames(functions, sourceLines)
	}

	if len(functions) > 0 {
		sort.Slice(functions, func(i, j int) bool { return functions[i].FirstLine < functions[j].FirstLine })

		for i := 0; i < len(functions)-1; i++ {
			functions[i].LastLine = functions[i+1].FirstLine - 1
		}
		functions[len(functions)-1].LastLine = codeFile.TotalLines

		for i := range functions {
			method := &functions[i]
			var methodLinesCovered, methodLinesValid, methodBranchesCovered, methodBranchesValid int

			for _, line := range codeFile.Lines {
				if line.Number >= method.FirstLine && line.Number <= method.LastLine {
					method.Lines = append(method.Lines, line)
					if line.Hits >= 0 {
						methodLinesValid++
						if line.Hits > 0 {
							methodLinesCovered++
						}
					}
					if line.IsBranchPoint {
						methodBranchesCovered += line.CoveredBranches
						methodBranchesValid += line.TotalBranches
					}
				}
			}

			method.LineRate = utils.CalculatePercentage(methodLinesCovered, methodLinesValid, 1) / 100
			if o.detectedBranchCoverage && methodBranchesValid > 0 {
				branchRate := utils.CalculatePercentage(methodBranchesCovered, methodBranchesValid, 1) / 100
				method.BranchRate = &branchRate
			}
			o.populateStandardMethodMetrics(method)
		}

		for _, method := range functions {
			coverageQuota := method.LineRate * 100
			codeElement := model.CodeElement{
				Name:          utils.GetShortMethodName(method.DisplayName),
				FullName:      method.DisplayName,
				Type:          model.MethodElementType,
				FirstLine:     method.FirstLine,
				LastLine:      method.LastLine,
				CoverageQuota: &coverageQuota,
			}
			codeFile.CodeElements = append(codeFile.CodeElements, codeElement)
		}
		utils.SortByLineAndName(codeFile.CodeElements)
	}

	class.Methods = functions

	for _, method := range class.Methods {
		if method.MethodMetrics != nil {
			codeFile.MethodMetrics = append(codeFile.MethodMetrics, method.MethodMetrics...)
		}
	}
}

func (o *processingOrchestrator) populateStandardMethodMetrics(method *model.Method) {
	method.MethodMetrics = []model.MethodMetric{}
	shortMetricName := utils.GetShortMethodName(method.DisplayName)

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
}

func (o *processingOrchestrator) aggregateClassMetrics(class *model.Class) {
	for _, f := range class.Files {
		class.LinesCovered += f.CoveredLines
		class.LinesValid += f.CoverableLines

		var bcovered, bvalid int
		for _, line := range f.Lines {
			if line.IsBranchPoint {
				bcovered += line.CoveredBranches
				bvalid += line.TotalBranches
			}
		}

		if o.detectedBranchCoverage && bvalid > 0 {
			if class.BranchesCovered == nil {
				class.BranchesCovered = new(int)
				class.BranchesValid = new(int)
			}
			*class.BranchesCovered += bcovered
			*class.BranchesValid += bvalid
		}
	}

	class.TotalMethods = len(class.Methods)
	for _, method := range class.Methods {
		if method.LineRate > 0 {
			class.CoveredMethods++
		}
		if method.LineRate >= 1.0 {
			class.FullyCoveredMethods++
		}
	}
}

func (o *processingOrchestrator) aggregateAssemblyMetrics(assembly *model.Assembly) {
	for _, class := range assembly.Classes {
		assembly.LinesCovered += class.LinesCovered
		assembly.LinesValid += class.LinesValid
		if class.BranchesCovered != nil && class.BranchesValid != nil {
			if assembly.BranchesCovered == nil {
				assembly.BranchesCovered = new(int)
				assembly.BranchesValid = new(int)
			}
			*assembly.BranchesCovered += *class.BranchesCovered
			*assembly.BranchesValid += *class.BranchesValid
		}
	}
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
