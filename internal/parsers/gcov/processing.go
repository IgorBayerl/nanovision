package gcov

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/IgorBayerl/AdlerCov/internal/filereader"
	"github.com/IgorBayerl/AdlerCov/internal/model"
	"github.com/IgorBayerl/AdlerCov/internal/parsers"
	"github.com/IgorBayerl/AdlerCov/internal/utils"
)

var (
	lineCoverageRegex   = regexp.MustCompile(`^\s*(?P<Visits>-|#####|=====|\d+):\s*(?P<LineNumber>[1-9]\d*):.*`)
	branchCoverageRegex = regexp.MustCompile(`^branch\s*(?P<Number>\d+)\s*(?:taken\s*(?P<Visits>\d+)|never executed)`)
	functionRegex       = regexp.MustCompile(`^function\s(?P<Name>.*?)\sline\s(?P<Line>\d+)`)
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

	o.parseCoverageData(lines, codeFile, sourceLines)

	class.Files = append(class.Files, *codeFile)
	o.aggregateClassMetrics(class)
	assembly.Classes = append(assembly.Classes, *class)
	o.aggregateAssemblyMetrics(assembly)

	return assembly, nil
}

func (o *processingOrchestrator) parseCoverageData(gcovLines []string, codeFile *model.CodeFile, sourceLines []string) {
	lineAnalyses := make(map[int]*model.Line)
	var functions []model.Method

	var lastLineNumber int

	for _, line := range gcovLines {
		if match := lineCoverageRegex.FindStringSubmatch(line); match != nil {
			visitsText := match[1]
			lineNumber, _ := strconv.Atoi(match[2])
			lastLineNumber = lineNumber

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
				// CORRECTED: Reassign the result of append to satisfy the linter.
				lineAnalysis.Branch = append(lineAnalysis.Branch, model.BranchCoverageDetail{
					Identifier: strconv.Itoa(branchNumber),
					Visits:     visits,
				})
			}
		} else if match := functionRegex.FindStringSubmatch(line); match != nil {
			functions = append(functions, model.Method{
				Name:      match[1],
				FirstLine: utils.ParseInt(match[2], 0),
			})
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

		// CORRECTED: Calculate file-level metrics here.
		if line.Hits >= 0 {
			coverableLines++
			if line.Hits > 0 {
				coveredLines++
			}
		}
	}

	// CORRECTED: Set the final file-level metrics.
	codeFile.CoveredLines = coveredLines
	codeFile.CoverableLines = coverableLines

	codeFile.CodeElements = []model.CodeElement{}
}

func (o *processingOrchestrator) aggregateClassMetrics(class *model.Class) {
	for _, f := range class.Files {
		// CORRECTED: Directly use the pre-calculated file metrics. This removes the "unused write" warnings.
		class.LinesCovered += f.CoveredLines
		class.LinesValid += f.CoverableLines

		// Branch aggregation remains the same as it's a sum of line-level data.
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
