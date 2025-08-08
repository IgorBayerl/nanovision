package htmlreport

import (
	"encoding/json"
	"fmt"
	"html/template"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/IgorBayerl/AdlerCov/internal/utils"
)

const decimalPlaces int = 1

func (b *HtmlReportBuilder) generateClassDetailHTML(classModel *Class, classReportFilename string, tag string) error {
	classVM := b.buildClassViewModelForDetailServer(classModel, tag)

	angularClassDetailForJS, err := b.buildAngularClassDetailForJS(classModel)
	if err != nil {
		return fmt.Errorf("failed to build Angular class detail JSON for %s: %w", classModel.DisplayName, err)
	}
	classDetailJSONBytes, err := json.Marshal(angularClassDetailForJS)
	if err != nil {
		return fmt.Errorf("failed to marshal Angular class detail JSON for %s: %w", classModel.DisplayName, err)
	}

	templateData := b.buildClassDetailPageData(classVM, tag, template.JS(classDetailJSONBytes))
	return b.renderClassDetailPage(templateData, classReportFilename)
}

func (b *HtmlReportBuilder) buildClassViewModelForDetailServer(classModel *Class, tag string) ClassViewModelForDetail {
	cvm := ClassViewModelForDetail{
		Name:         classModel.DisplayName,
		AssemblyName: classModel.Name,
		IsMultiFile:  len(classModel.Files) > 1,
	}
	if dotIndex := strings.LastIndex(classModel.Name, "."); dotIndex > -1 && dotIndex < len(classModel.Name)-1 {
		cvm.AssemblyName = classModel.Name[:dotIndex]
	}

	cvm.CoveredLines = classModel.LinesCovered
	cvm.CoverableLines = classModel.LinesValid
	cvm.UncoveredLines = cvm.CoverableLines - cvm.CoveredLines
	cvm.TotalLines = classModel.TotalLines

	b.populateLineCoverageMetricsForClassVM(&cvm, classModel)
	b.populateBranchCoverageMetricsForClassVM(&cvm, classModel)
	b.populateMethodCoverageMetricsForClassVM(&cvm, classModel)
	b.populateHistoricCoveragesForClassVM(&cvm, classModel)
	b.populateAggregatedMetricsForClassVM(&cvm, classModel)

	metricsTable := b.buildMetricsTableForClassVM(classModel)
	cvm.MetricsTable = metricsTable
	cvm.FilesWithMetrics = len(metricsTable.Rows) > 0

	sortedFiles := make([]CodeFile, len(classModel.Files))
	copy(sortedFiles, classModel.Files)
	sort.Slice(sortedFiles, func(i, j int) bool { return sortedFiles[i].Path < sortedFiles[j].Path })

	for fileIdx, fileInClassValue := range sortedFiles {
		fileInClass := fileInClassValue
		fileVM := b.buildFileViewModelForServerRender(&fileInClass)
		cvm.Files = append(cvm.Files, fileVM)

		for i := range fileInClass.CodeElements {
			codeElem := &fileInClass.CodeElements[i]
			sidebarElem := b.buildSidebarElementViewModel(codeElem, fileVM.ShortPath, fileIdx+1, len(sortedFiles) > 1)
			cvm.SidebarElements = append(cvm.SidebarElements, sidebarElem)
		}
	}

	return cvm
}

func (b *HtmlReportBuilder) populateLineCoverageMetricsForClassVM(cvm *ClassViewModelForDetail, classModel *Class) {
	lineCoverage := utils.CalculatePercentage(cvm.CoveredLines, cvm.CoverableLines, b.ReportContext.Config().MaximumDecimalPlacesForCoverageQuotas)
	cvm.CoveragePercentageForDisplay = utils.FormatPercentage(lineCoverage, b.ReportContext.Config().MaximumDecimalPlacesForPercentageDisplay)

	if !math.IsNaN(lineCoverage) {
		cvm.CoveragePercentageBarValue = 100 - int(math.Round(lineCoverage))
		cvm.CoverageRatioTextForDisplay = fmt.Sprintf("%d of %d", cvm.CoveredLines, cvm.CoverableLines)
	} else {
		cvm.CoveragePercentageBarValue = 0
		cvm.CoverageRatioTextForDisplay = "-"
	}
}

func (b *HtmlReportBuilder) populateBranchCoverageMetricsForClassVM(cvm *ClassViewModelForDetail, classModel *Class) {
	if b.branchCoverageAvailable && classModel.BranchesValid != nil && *classModel.BranchesValid > 0 && classModel.BranchesCovered != nil {
		cvm.CoveredBranches = *classModel.BranchesCovered
		cvm.TotalBranches = *classModel.BranchesValid
		branchCoverage := utils.CalculatePercentage(*classModel.BranchesCovered, *classModel.BranchesValid, b.ReportContext.Config().MaximumDecimalPlacesForCoverageQuotas)
		cvm.BranchCoveragePercentageForDisplay = utils.FormatPercentage(branchCoverage, b.ReportContext.Config().MaximumDecimalPlacesForPercentageDisplay)
	} else {
		cvm.BranchCoveragePercentageForDisplay = "N/A"
	}
}

func (b *HtmlReportBuilder) populateMethodCoverageMetricsForClassVM(cvm *ClassViewModelForDetail, classModel *Class) {
	cvm.TotalMethods = classModel.TotalMethods
	cvm.CoveredMethods = classModel.CoveredMethods
	cvm.FullyCoveredMethods = classModel.FullyCoveredMethods

	if cvm.TotalMethods > 0 {
		// Calculate with configured precision
		methodCovVal := utils.CalculatePercentage(cvm.CoveredMethods, cvm.TotalMethods, decimalPlaces)
		fullMethodCovVal := utils.CalculatePercentage(cvm.FullyCoveredMethods, cvm.TotalMethods, decimalPlaces)

		// Format for display with 0 decimal places
		cvm.MethodCoveragePercentageForDisplay = utils.FormatPercentage(methodCovVal, decimalPlaces)
		cvm.FullMethodCoveragePercentageForDisplay = utils.FormatPercentage(fullMethodCovVal, decimalPlaces)

		cvm.MethodCoveragePercentageBarValue = 100 - int(math.Round(methodCovVal)) // Bar value should use the calculated value
		cvm.MethodCoverageRatioTextForDisplay = fmt.Sprintf("%d of %d", cvm.CoveredMethods, cvm.TotalMethods)
		cvm.FullMethodCoverageRatioTextForDisplay = fmt.Sprintf("%d of %d", cvm.FullyCoveredMethods, cvm.TotalMethods)
	} else {
		cvm.MethodCoveragePercentageForDisplay = "N/A"
		cvm.MethodCoveragePercentageBarValue = 0
		cvm.MethodCoverageRatioTextForDisplay = "-"
		cvm.FullMethodCoveragePercentageForDisplay = "N/A"
		cvm.FullMethodCoverageRatioTextForDisplay = "-"
	}
}

func (b *HtmlReportBuilder) populateHistoricCoveragesForClassVM(cvm *ClassViewModelForDetail, classModel *Class) {
	if classModel.HistoricCoverages == nil {
		return
	}
	for _, hist := range classModel.HistoricCoverages {
		angularHist := b.buildAngularHistoricCoverageViewModel(&hist)
		cvm.HistoricCoverages = append(cvm.HistoricCoverages, angularHist)
		if angularHist.LineCoverageQuota >= 0 {
			cvm.LineCoverageHistory = append(cvm.LineCoverageHistory, angularHist.LineCoverageQuota)
		}
		if angularHist.BranchCoverageQuota >= 0 {
			cvm.BranchCoverageHistory = append(cvm.BranchCoverageHistory, angularHist.BranchCoverageQuota)
		}
	}
}

func (b *HtmlReportBuilder) populateAggregatedMetricsForClassVM(cvm *ClassViewModelForDetail, classModel *Class) {
	cvm.Metrics = make(map[string]float64)
	for name, val := range classModel.Metrics {
		cvm.Metrics[name] = val
	}
}

func (b *HtmlReportBuilder) buildFileViewModelForServerRender(fileInClass *CodeFile) FileViewModelForDetail {
	fileVM := FileViewModelForDetail{
		Path:      fileInClass.Path,
		ShortPath: utils.ReplaceInvalidPathChars(filepath.Base(fileInClass.Path)),
	}

	// This function now correctly iterates through all lines from the legacy model.
	for _, modelCovLine := range fileInClass.Lines {
		lineVM := b.buildLineViewModelForServerRender(modelCovLine.Content, modelCovLine.Number, &modelCovLine)
		fileVM.Lines = append(fileVM.Lines, lineVM)
	}
	return fileVM
}

func (b *HtmlReportBuilder) buildLineViewModelForServerRender(lineContent string, actualLineNumber int, modelCovLine *Line) LineViewModelForDetail {
	lineVM := LineViewModelForDetail{
		LineNumber:      actualLineNumber,
		LineContent:     lineContent,
		LineVisitStatus: lineVisitStatusToString(modelCovLine.LineVisitStatus),
	}
	dataCoverageMap := map[string]map[string]string{"AllTestMethods": {"VC": "", "LVS": "gray"}}

	if modelCovLine.LineVisitStatus != NotCoverable {
		lineVM.Hits = fmt.Sprintf("%d", modelCovLine.Hits)
		dataCoverageMap["AllTestMethods"]["VC"] = lineVM.Hits
		dataCoverageMap["AllTestMethods"]["LVS"] = lineVM.LineVisitStatus

		if modelCovLine.IsBranchPoint && modelCovLine.TotalBranches > 0 {
			lineVM.IsBranch = true
			branchCoverageVal := (float64(modelCovLine.CoveredBranches) / float64(modelCovLine.TotalBranches)) * 100.0
			lineVM.BranchBarValue = 100 - int(math.Round(branchCoverageVal))
		}

		tooltipBranchRate := ""
		if lineVM.IsBranch {
			tooltipBranchRate = fmt.Sprintf(", %d of %d branches are covered", modelCovLine.CoveredBranches, modelCovLine.TotalBranches)
		}

		switch modelCovLine.LineVisitStatus {
		case Covered:
			lineVM.Tooltip = fmt.Sprintf("Covered (%d visits%s)", modelCovLine.Hits, tooltipBranchRate)
		case NotCovered:
			lineVM.Tooltip = fmt.Sprintf("Not covered (%d visits%s)", modelCovLine.Hits, tooltipBranchRate)
		case PartiallyCovered:
			lineVM.Tooltip = fmt.Sprintf("Partially covered (%d visits%s)", modelCovLine.Hits, tooltipBranchRate)
		}
	} else {
		lineVM.Hits = ""
		lineVM.Tooltip = "Not coverable"
	}

	dataCoverageBytes, _ := json.Marshal(dataCoverageMap)
	lineVM.DataCoverage = template.JS(dataCoverageBytes)
	return lineVM
}

func (b *HtmlReportBuilder) buildSidebarElementViewModel(codeElem *CodeElement, fileShortPath string, fileIndexPlus1 int, isMultiFile bool) SidebarElementViewModel {
	sidebarElem := SidebarElementViewModel{
		Name:          codeElem.Name,
		FullName:      codeElem.FullName,
		FileShortPath: fileShortPath,
		Line:          codeElem.FirstLine,
		Icon:          "cube",
	}
	if isMultiFile {
		sidebarElem.FileIndexPlus1 = fileIndexPlus1
	}
	if codeElem.Type == PropertyElementType {
		sidebarElem.Icon = "wrench"
	}

	var coverageTitleText string
	if codeElem.CoverageQuota != nil {
		sidebarElem.CoverageBarValue = getCoverageBarValue(*codeElem.CoverageQuota)
		coverageTitleText = fmt.Sprintf("Line coverage: %.1f%%", *codeElem.CoverageQuota)
	} else {
		sidebarElem.CoverageBarValue = -1
		coverageTitleText = "Line coverage: N/A"
	}
	sidebarElem.CoverageTitle = fmt.Sprintf("%s - %s", coverageTitleText, codeElem.FullName)
	return sidebarElem
}

func (b *HtmlReportBuilder) getStandardMetricHeaders() []AngularMetricDefinitionViewModel {
	standardMetricKeys := []string{
		"Branch coverage",
		"CrapScore",
		"Cyclomatic complexity",
		"Line coverage",
	}
	var headers []AngularMetricDefinitionViewModel
	for _, key := range standardMetricKeys {
		translatedName := b.translations[key]
		if translatedName == "" {
			translatedName = key
		}
		headers = append(headers, AngularMetricDefinitionViewModel{
			Name:           translatedName,
			ExplanationURL: b.getMetricExplanationURL(key),
		})
	}
	return headers
}

func (b *HtmlReportBuilder) buildSingleMetricRow(
	method *Method,
	correspondingCE *CodeElement,
	fileShortPath string,
	fileIndexPlus1 int,
	headers []AngularMetricDefinitionViewModel,
) AngularMethodMetricsViewModel {
	var fullNameForTitle string
	var lineToLink int
	var isProperty bool
	var coverageQuota *float64
	cleanedFullName := method.DisplayName
	fullNameForTitle = cleanedFullName
	if correspondingCE != nil {
		lineToLink = correspondingCE.FirstLine
		isProperty = (correspondingCE.Type == PropertyElementType)
		coverageQuota = correspondingCE.CoverageQuota
	} else {
		lineToLink = method.FirstLine
		isProperty = strings.HasPrefix(cleanedFullName, "get_") || strings.HasPrefix(cleanedFullName, "set_")
	}
	var shortDisplayNameForTable string
	if isProperty {
		shortDisplayNameForTable = cleanedFullName
	} else {
		shortDisplayNameForTable = utils.GetShortMethodName(cleanedFullName)
	}
	row := AngularMethodMetricsViewModel{
		Name:           shortDisplayNameForTable,
		FullName:       fullNameForTitle,
		FileIndexPlus1: fileIndexPlus1,
		Line:           lineToLink,
		FileShortPath:  fileShortPath,
		IsProperty:     isProperty,
		CoverageQuota:  coverageQuota,
		MetricValues:   make([]string, len(headers)),
	}

	// Create a map for easy lookup of existing metrics for the method
	methodMetricsMap := make(map[string]Metric)
	for _, mm := range method.MethodMetrics {
		for _, m := range mm.Metrics {
			methodMetricsMap[m.Name] = m
		}
	}

	for i, headerVM := range headers {
		var originalMetricKey string

		// This switch maps the display name of the header back to the metric's key.
		switch headerVM.Name {
		case b.translations["Branch coverage"]:
			originalMetricKey = "Branch coverage"
		case b.translations["CrapScore"]:
			originalMetricKey = "CrapScore"
		case b.translations["Cyclomatic complexity"]:
			originalMetricKey = "Cyclomatic complexity"
		case b.translations["Line coverage"]:
			originalMetricKey = "Line coverage"
		default:
			originalMetricKey = headerVM.Name
		}

		// Use new ratio format for line and branch coverage
		if originalMetricKey == "Line coverage" {
			if method.LinesValid > 0 {
				row.MetricValues[i] = fmt.Sprintf("%d/%d", method.LinesCovered, method.LinesValid)
			} else {
				row.MetricValues[i] = "-"
			}
		} else if originalMetricKey == "Branch coverage" {
			if b.branchCoverageAvailable {
				if method.BranchesValid > 0 {
					row.MetricValues[i] = fmt.Sprintf("%d/%d", method.BranchesCovered, method.BranchesValid)
				} else {
					row.MetricValues[i] = "-" // Method has no branches
				}
			} else {
				row.MetricValues[i] = "N/A"
			}
		} else {
			// Use existing logic for other metrics like Cyclomatic Complexity
			if metric, ok := methodMetricsMap[originalMetricKey]; ok {
				row.MetricValues[i] = b.formatMetricValue(metric)
			} else {
				row.MetricValues[i] = "N/A"
			}
		}
	}
	return row
}

// buildMetricsTableForClassVM constructs the view model for the metrics table.
// It collects all methods from all files within the class and sorts them
// primarily by file path, then by line number, then by short method name.
func (b *HtmlReportBuilder) buildMetricsTableForClassVM(classModel *Class) MetricsTableViewModel {
	metricsTable := MetricsTableViewModel{}
	metricsTable.Headers = b.getStandardMetricHeaders()

	if len(classModel.Methods) == 0 && len(classModel.Files) == 0 { // Check if there are any files to iterate
		return metricsTable
	}

	// Create a temporary struct to hold methods along with their file context for sorting
	type methodWithFileContext struct {
		method         *Method
		filePath       string // Full path for primary sort
		fileShortPath  string // For linking
		fileIndexPlus1 int    // For display in multi-file scenarios
	}
	var allMethodsWithContext []methodWithFileContext

	// Sort files first to ensure consistent file indexing and path usage
	sortedFiles := make([]CodeFile, len(classModel.Files))
	copy(sortedFiles, classModel.Files)
	sort.Slice(sortedFiles, func(i, j int) bool {
		return sortedFiles[i].Path < sortedFiles[j].Path
	})

	// Collect all methods from all files, associating them with their file context
	for fileIdx, file := range sortedFiles {
		// Methods within a CodeFile's MethodMetrics list might not be what we want directly.
		// We need to iterate through model.Method objects that are part of the classModel.Methods,
		// and then find which file they belong to for sorting.
		// A better way: iterate classModel.Methods, and for each method, find its file.
		// However, model.Method doesn't directly link back to a specific CodeFile.
		// We need to find the methods that are defined within this specific file.
		// The model.Method.FirstLine and model.Method.DisplayName are key.
		// CodeFile.CodeElements helps link method display names to file lines.

		// Let's find methods that are defined in *this* specific file (file)
		// by checking if their first line is within this file's scope
		// and if a corresponding code element exists.
		for methIdx := range classModel.Methods {
			method := &classModel.Methods[methIdx]
			// Check if this method is primarily defined in the current file
			// This is a bit heuristic: a method might span files in partial classes,
			// but for metrics, we usually associate it with its main definition file.
			// The `CodeElement` for this method within `file.CodeElements` will confirm.
			var foundInThisFile bool
			for ceIdx := range file.CodeElements {
				ce := &file.CodeElements[ceIdx]
				if ce.FirstLine == method.FirstLine && ce.FullName == method.DisplayName {
					foundInThisFile = true
					break
				}
			}

			if foundInThisFile {
				allMethodsWithContext = append(allMethodsWithContext, methodWithFileContext{
					method:         method,
					filePath:       file.Path, // Full path of the file
					fileShortPath:  utils.ReplaceInvalidPathChars(filepath.Base(file.Path)),
					fileIndexPlus1: fileIdx + 1,
				})
			}
		}
	}

	// Sort all collected methods:
	// 1. Primary: File Path
	// 2. Secondary: Method's First Line
	// 3. Tertiary: Method's Short Display Name
	sort.Slice(allMethodsWithContext, func(i, j int) bool {
		itemI := allMethodsWithContext[i]
		itemJ := allMethodsWithContext[j]

		if itemI.filePath != itemJ.filePath {
			return itemI.filePath < itemJ.filePath
		}
		if itemI.method.FirstLine != itemJ.method.FirstLine {
			return itemI.method.FirstLine < itemJ.method.FirstLine
		}
		return utils.GetShortMethodName(itemI.method.DisplayName) < utils.GetShortMethodName(itemJ.method.DisplayName)
	})

	// Now build the rows from the sorted list
	for _, mCtx := range allMethodsWithContext {
		// Find the CodeElement again, this time specifically for the method in its context
		// (or pass it if already available from a more direct link)
		var correspondingCE *CodeElement
		for _, f := range classModel.Files { // Iterate original files to find the CE
			if f.Path == mCtx.filePath {
				for ceIdx := range f.CodeElements {
					ce := &f.CodeElements[ceIdx]
					if ce.FirstLine == mCtx.method.FirstLine && ce.FullName == mCtx.method.DisplayName {
						correspondingCE = ce
						break
					}
				}
			}
			if correspondingCE != nil {
				break
			}
		}

		if correspondingCE == nil && len(mCtx.method.MethodMetrics) > 0 && mCtx.method.MethodMetrics[0].Line == mCtx.method.FirstLine {
			// Fallback: Create a temporary CodeElement if it's truly missing but metrics exist for the method at its first line.
			// This situation suggests an inconsistency or a method that exists in metrics but not explicitly in CodeElements.
			// For metrics table purposes, we primarily need FirstLine, FullName (as DisplayName), and Type (Method).
			// The CoverageQuota for the method itself might be derived if available.
			var methCovQuota *float64
			if !math.IsNaN(mCtx.method.LineRate) {
				lrq := mCtx.method.LineRate * 100.0
				methCovQuota = &lrq
			}
			correspondingCE = &CodeElement{
				Name:          mCtx.method.DisplayName, // Use display name as short name for this fallback
				FullName:      mCtx.method.DisplayName,
				Type:          MethodElementType,
				FirstLine:     mCtx.method.FirstLine,
				LastLine:      mCtx.method.LastLine, // Approx
				CoverageQuota: methCovQuota,
			}
			// This warning is now more specific to metric table generation for a method without a clear CE
			fmt.Fprintf(os.Stderr, "Metrics Table Warning: Could not find exact CodeElement for method %s (line %d) in class %s for file %s. Using method data as fallback for table row.\n", mCtx.method.DisplayName, mCtx.method.FirstLine, classModel.DisplayName, mCtx.filePath)
		}

		row := b.buildSingleMetricRow(mCtx.method, correspondingCE, mCtx.fileShortPath, mCtx.fileIndexPlus1, metricsTable.Headers)
		metricsTable.Rows = append(metricsTable.Rows, row)
	}

	return metricsTable
}

func (b *HtmlReportBuilder) getMetricExplanationURL(metricKey string) string {
	switch metricKey {
	case "Cyclomatic complexity", "Complexity":
		return "https://en.wikipedia.org/wiki/Cyclomatic_complexity"
	case "CrapScore":
		return "https://testing.googleblog.com/2011/02/this-code-is-crap.html"
	case "Line coverage", "Branch coverage":
		return "https://en.wikipedia.org/wiki/Code_coverage"
	default:
		return ""
	}
}

func (b *HtmlReportBuilder) formatMetricValue(metric Metric) string {
	if metric.Value == nil {
		return "-"
	}
	valFloat, isFloat := metric.Value.(float64)
	if !isFloat {
		if valInt, isInt := metric.Value.(int); isInt {
			return fmt.Sprintf("%d", valInt)
		}
		return fmt.Sprintf("%v", metric.Value)
	}
	if math.IsNaN(valFloat) {
		return "NaN"
	}
	if math.IsInf(valFloat, 0) {
		return "Inf"
	}
	switch metric.Name {
	case "Line coverage", "Branch coverage":
		return utils.FormatPercentage(valFloat, decimalPlaces)
	case "CrapScore":
		return fmt.Sprintf("%.2f", valFloat)
	case "Cyclomatic complexity", "Complexity":
		return fmt.Sprintf("%.0f", valFloat)
	default:
		return fmt.Sprintf(fmt.Sprintf("%%.%df", decimalPlaces), valFloat)
	}
}

func (b *HtmlReportBuilder) buildAngularClassDetailForJS(classModel *Class) (AngularClassDetailViewModel, error) {
	// This function now takes the complete legacy model as its source of truth.
	classVMServer := b.buildClassViewModelForDetailServer(classModel, "") // We only need a subset of fields

	angularClassVMForJS := AngularClassViewModel{
		Name:                  classModel.DisplayName,
		CoveredLines:          classModel.LinesCovered,
		UncoveredLines:        classModel.LinesValid - classModel.LinesCovered,
		CoverableLines:        classModel.LinesValid,
		TotalLines:            classModel.TotalLines,
		CoveredMethods:        classVMServer.CoveredMethods,
		FullyCoveredMethods:   classVMServer.FullyCoveredMethods,
		TotalMethods:          classVMServer.TotalMethods,
		HistoricCoverages:     classVMServer.HistoricCoverages,
		LineCoverageHistory:   classVMServer.LineCoverageHistory,
		BranchCoverageHistory: classVMServer.BranchCoverageHistory,
		Metrics:               classVMServer.Metrics,
	}
	if classModel.BranchesCovered != nil {
		angularClassVMForJS.CoveredBranches = *classModel.BranchesCovered
	}
	if classModel.BranchesValid != nil {
		angularClassVMForJS.TotalBranches = *classModel.BranchesValid
	}

	detailVM := AngularClassDetailViewModel{Class: angularClassVMForJS, Files: []AngularCodeFileViewModel{}}
	if classModel.Files == nil {
		return detailVM, nil
	}

	for _, fileInClass := range classModel.Files {
		// Pass the legacy file model, which now contains all necessary data.
		angularFileForJS := b.buildAngularFileViewModelForJS(&fileInClass)
		detailVM.Files = append(detailVM.Files, angularFileForJS)
	}
	return detailVM, nil
}

func (b *HtmlReportBuilder) buildAngularFileViewModelForJS(fileInClass *CodeFile) AngularCodeFileViewModel {
	// This function NO LONGER performs any file I/O.
	// It relies on the Content field already populated in fileInClass.Lines.
	angularFile := AngularCodeFileViewModel{
		Path:           fileInClass.Path,
		CoveredLines:   fileInClass.CoveredLines,
		CoverableLines: fileInClass.CoverableLines,
		TotalLines:     fileInClass.TotalLines,
		Lines:          []AngularLineAnalysisViewModel{},
	}

	if fileInClass.Lines != nil {
		for _, modelCovLine := range fileInClass.Lines {
			// Pass the legacy line model directly.
			angularLine := b.buildAngularLineViewModelForJS(&modelCovLine)
			angularFile.Lines = append(angularFile.Lines, angularLine)
		}
	}

	return angularFile
}

func (b *HtmlReportBuilder) buildAngularLineViewModelForJS(modelCovLine *Line) AngularLineAnalysisViewModel {
	// This function is now much simpler. It just translates the legacy model to the Angular view model.
	lineVM := AngularLineAnalysisViewModel{
		LineNumber:      modelCovLine.Number,
		LineContent:     modelCovLine.Content, // Use the content that was already read.
		Hits:            modelCovLine.Hits,
		CoveredBranches: modelCovLine.CoveredBranches,
		TotalBranches:   modelCovLine.TotalBranches,
		LineVisitStatus: lineVisitStatusToString(modelCovLine.LineVisitStatus),
	}
	return lineVM
}

func (b *HtmlReportBuilder) buildClassDetailPageData(classVM ClassViewModelForDetail, tag string, classDetailJS template.JS) ClassDetailData {
	appVersion := "0.0.1" // Simplified
	return ClassDetailData{
		ReportTitle:                           b.reportTitle,
		AppVersion:                            appVersion,
		CurrentDateTime:                       time.Now().Format("02/01/2006 - 15:04:05"),
		Class:                                 classVM,
		BranchCoverageAvailable:               b.branchCoverageAvailable,
		MethodCoverageAvailable:               b.methodCoverageAvailable,
		Tag:                                   tag,
		Translations:                          b.translations,
		MaximumDecimalPlacesForCoverageQuotas: b.ReportContext.Config().MaximumDecimalPlacesForCoverageQuotas,
		AngularCssFile:                        b.angularCssFile,
		CombinedAngularJsFile:                 b.combinedAngularJsFile,
		AssembliesJSON:                        b.assembliesJSON,
		RiskHotspotsJSON:                      b.riskHotspotsJSON,
		MetricsJSON:                           b.metricsJSON,
		RiskHotspotMetricsJSON:                b.riskHotspotMetricsJSON,
		HistoricCoverageExecutionTimesJSON:    b.historicCoverageExecutionTimesJSON,
		TranslationsJSON:                      b.translationsJSON,
		ClassDetailJSON:                       classDetailJS,
	}
}
