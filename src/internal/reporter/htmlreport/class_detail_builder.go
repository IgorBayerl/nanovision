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

	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/filereader"
	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/model"
	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/utils"
)

func (b *HtmlReportBuilder) generateClassDetailHTML(classModel *model.Class, classReportFilename string, tag string) error {
	// 1. Build the main ClassViewModelForDetail (server-side rendering focus)
	classVM := b.buildClassViewModelForDetailServer(classModel, tag)

	// 2. Build the AngularClassDetailViewModel (for client-side window.classDetails JSON)
	angularClassDetailForJS, err := b.buildAngularClassDetailForJS(classModel, &classVM)
	if err != nil {
		return fmt.Errorf("failed to build Angular class detail JSON for %s: %w", classModel.DisplayName, err)
	}
	classDetailJSONBytes, err := json.Marshal(angularClassDetailForJS)
	if err != nil {
		return fmt.Errorf("failed to marshal Angular class detail JSON for %s: %w", classModel.DisplayName, err)
	}

	// 3. Prepare overall data for the template
	templateData := b.buildClassDetailPageData(classVM, tag, template.JS(classDetailJSONBytes))

	// 4. Render the template
	return b.renderClassDetailPage(templateData, classReportFilename)

}

func (b *HtmlReportBuilder) buildClassViewModelForDetailServer(classModel *model.Class, tag string) ClassViewModelForDetail {
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

	var allMethodMetricsForClass []*model.MethodMetric

	sortedFiles := make([]model.CodeFile, len(classModel.Files))
	copy(sortedFiles, classModel.Files)
	sort.Slice(sortedFiles, func(i, j int) bool {
		return sortedFiles[i].Path < sortedFiles[j].Path
	})

	for fileIdx, fileInClassValue := range sortedFiles {
		fileInClass := fileInClassValue
		fileVM, _, err := b.buildFileViewModelForServerRender(&fileInClass)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not build file view model for %s: %v\n", fileInClass.Path, err)
			continue
		}
		cvm.Files = append(cvm.Files, fileVM)

		for i := range fileInClass.MethodMetrics {
			allMethodMetricsForClass = append(allMethodMetricsForClass, &fileInClass.MethodMetrics[i])
		}

		for i := range fileInClass.CodeElements {
			codeElem := &fileInClass.CodeElements[i]
			sidebarElem := b.buildSidebarElementViewModel(codeElem, fileVM.ShortPath, fileIdx+1, len(sortedFiles) > 1)
			cvm.SidebarElements = append(cvm.SidebarElements, sidebarElem)
		}
	}

	if len(allMethodMetricsForClass) > 0 {
		cvm.FilesWithMetrics = true
		cvm.MetricsTable = b.buildMetricsTableForClassVM(classModel) // Pass the original classModel
	}

	return cvm
}

func (b *HtmlReportBuilder) populateLineCoverageMetricsForClassVM(cvm *ClassViewModelForDetail, classModel *model.Class) {
	lineCoverage := utils.CalculatePercentage(cvm.CoveredLines, cvm.CoverableLines, b.maximumDecimalPlacesForCoverageQuotas)
	cvm.CoveragePercentageForDisplay = utils.FormatPercentage(lineCoverage, b.maximumDecimalPlacesForPercentageDisplay)

	if !math.IsNaN(lineCoverage) {

		cvm.CoveragePercentageBarValue = 100 - int(math.Round(lineCoverage))
		cvm.CoverageRatioTextForDisplay = fmt.Sprintf("%d of %d", cvm.CoveredLines, cvm.CoverableLines)
	} else {
		cvm.CoveragePercentageBarValue = 0
		cvm.CoverageRatioTextForDisplay = "-"
	}
}

func (b *HtmlReportBuilder) populateBranchCoverageMetricsForClassVM(cvm *ClassViewModelForDetail, classModel *model.Class) {
	if b.branchCoverageAvailable && classModel.BranchesValid != nil && *classModel.BranchesValid > 0 && classModel.BranchesCovered != nil {
		cvm.CoveredBranches = *classModel.BranchesCovered
		cvm.TotalBranches = *classModel.BranchesValid
		branchCoverage := utils.CalculatePercentage(*classModel.BranchesCovered, *classModel.BranchesValid, b.maximumDecimalPlacesForCoverageQuotas)
		cvm.BranchCoveragePercentageForDisplay = utils.FormatPercentage(branchCoverage, b.maximumDecimalPlacesForPercentageDisplay)

		if !math.IsNaN(branchCoverage) {

			cvm.BranchCoveragePercentageBarValue = 100 - int(math.Round(branchCoverage))
			cvm.BranchCoverageRatioTextForDisplay = fmt.Sprintf("%d of %d", cvm.CoveredBranches, cvm.TotalBranches)
		} else {
			cvm.BranchCoveragePercentageBarValue = 0
			cvm.BranchCoverageRatioTextForDisplay = "-"
		}
	} else {
		cvm.BranchCoveragePercentageForDisplay = "N/A"
		cvm.BranchCoveragePercentageBarValue = 0
		cvm.BranchCoverageRatioTextForDisplay = "-"
	}
}

func (b *HtmlReportBuilder) populateMethodCoverageMetricsForClassVM(cvm *ClassViewModelForDetail, classModel *model.Class) {
	cvm.TotalMethods = classModel.TotalMethods
	cvm.CoveredMethods = classModel.CoveredMethods
	cvm.FullyCoveredMethods = classModel.FullyCoveredMethods

	if cvm.TotalMethods > 0 {
		// Calculate with configured precision
		methodCovVal := utils.CalculatePercentage(cvm.CoveredMethods, cvm.TotalMethods, b.maximumDecimalPlacesForCoverageQuotas)
		fullMethodCovVal := utils.CalculatePercentage(cvm.FullyCoveredMethods, cvm.TotalMethods, b.maximumDecimalPlacesForCoverageQuotas)

		// Format for display with 0 decimal places
		cvm.MethodCoveragePercentageForDisplay = utils.FormatPercentage(methodCovVal, b.maximumDecimalPlacesForPercentageDisplay)
		cvm.FullMethodCoveragePercentageForDisplay = utils.FormatPercentage(fullMethodCovVal, b.maximumDecimalPlacesForPercentageDisplay)

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

func (b *HtmlReportBuilder) populateHistoricCoveragesForClassVM(cvm *ClassViewModelForDetail, classModel *model.Class) {
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

func (b *HtmlReportBuilder) populateAggregatedMetricsForClassVM(cvm *ClassViewModelForDetail, classModel *model.Class) {
	cvm.Metrics = make(map[string]float64)
	for name, val := range classModel.Metrics {
		cvm.Metrics[name] = val
	}
}

func (b *HtmlReportBuilder) buildFileViewModelForServerRender(fileInClass *model.CodeFile) (FileViewModelForDetail, []string, error) {
	fileVM := FileViewModelForDetail{
		Path:      fileInClass.Path,
		ShortPath: utils.ReplaceInvalidPathChars(filepath.Base(fileInClass.Path)),
	}
	sourceLines, err := filereader.ReadLinesInFile(fileInClass.Path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not read source file %s: %v\n", fileInClass.Path, err)
		sourceLines = []string{}
	}

	coverageLinesMap := make(map[int]*model.Line)
	for i := range fileInClass.Lines {
		covLine := &fileInClass.Lines[i]
		coverageLinesMap[covLine.Number] = covLine
	}

	for lineNumIdx, lineContent := range sourceLines {
		actualLineNumber := lineNumIdx + 1
		modelCovLine, hasCoverageData := coverageLinesMap[actualLineNumber]
		lineVM := b.buildLineViewModelForServerRender(lineContent, actualLineNumber, modelCovLine, hasCoverageData)
		fileVM.Lines = append(fileVM.Lines, lineVM)
	}
	return fileVM, sourceLines, nil
}

func (b *HtmlReportBuilder) buildLineViewModelForServerRender(lineContent string, actualLineNumber int, modelCovLine *model.Line, hasCoverageData bool) LineViewModelForDetail {
	lineVM := LineViewModelForDetail{LineNumber: actualLineNumber, LineContent: lineContent}
	dataCoverageMap := map[string]map[string]string{"AllTestMethods": {"VC": "", "LVS": "gray"}}

	if hasCoverageData {
		lineVM.Hits = fmt.Sprintf("%d", modelCovLine.Hits)
		status := determineLineVisitStatus(modelCovLine.Hits, modelCovLine.IsBranchPoint, modelCovLine.CoveredBranches, modelCovLine.TotalBranches)
		lineVM.LineVisitStatus = lineVisitStatusToString(status)
		if modelCovLine.IsBranchPoint && modelCovLine.TotalBranches > 0 {
			lineVM.IsBranch = true
			branchCoverageVal := (float64(modelCovLine.CoveredBranches) / float64(modelCovLine.TotalBranches)) * 100.0
			lineVM.BranchBarValue = 100 - int(math.Round(branchCoverageVal))
		}
		dataCoverageMap["AllTestMethods"]["VC"] = fmt.Sprintf("%d", modelCovLine.Hits)
		dataCoverageMap["AllTestMethods"]["LVS"] = lineVM.LineVisitStatus
		tooltipBranchRate := ""
		if lineVM.IsBranch {
			tooltipBranchRate = fmt.Sprintf(", %d of %d branches are covered", modelCovLine.CoveredBranches, modelCovLine.TotalBranches)
		}
		switch status {
		case lineVisitStatusCovered:
			lineVM.Tooltip = fmt.Sprintf("Covered (%d visits%s)", modelCovLine.Hits, tooltipBranchRate)
		case lineVisitStatusNotCovered:
			lineVM.Tooltip = fmt.Sprintf("Not covered (%d visits%s)", modelCovLine.Hits, tooltipBranchRate)
		case lineVisitStatusPartiallyCovered:
			lineVM.Tooltip = fmt.Sprintf("Partially covered (%d visits%s)", modelCovLine.Hits, tooltipBranchRate)
		default:
			lineVM.Tooltip = "Not coverable"
		}
	} else {
		lineVM.LineVisitStatus = lineVisitStatusToString(lineVisitStatusNotCoverable)
		lineVM.Hits = ""
		lineVM.Tooltip = "Not coverable"
	}
	dataCoverageBytes, _ := json.Marshal(dataCoverageMap)
	lineVM.DataCoverage = template.JS(dataCoverageBytes)
	return lineVM

}

func (b *HtmlReportBuilder) buildSidebarElementViewModel(codeElem *model.CodeElement, fileShortPath string, fileIndexPlus1 int, isMultiFile bool) SidebarElementViewModel {
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
	if codeElem.Type == model.PropertyElementType {
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
	method *model.Method,
	correspondingCE *model.CodeElement,
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
		isProperty = (correspondingCE.Type == model.PropertyElementType)
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
	methodMetricsMap := make(map[string]model.Metric)
	for _, mm := range method.MethodMetrics {
		for _, m := range mm.Metrics {
			methodMetricsMap[m.Name] = m
		}
	}

	// Manually add Line Coverage and Branch Coverage from the method model
	// to ensure they are available for formatting.
	methodMetricsMap["Line coverage"] = model.Metric{Value: method.LineRate * 100.0}
	if method.BranchRate != nil {
		methodMetricsMap["Branch coverage"] = model.Metric{Value: *method.BranchRate * 100.0}
	}
	// Note: Complexity and CrapScore are already in method.MethodMetrics, so they'll be in the map.

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

		if metric, ok := methodMetricsMap[originalMetricKey]; ok {
			row.MetricValues[i] = b.formatMetricValue(metric)
		} else {
			// If the metric is not in the map (e.g., Branch coverage for a method where BranchRate was nil),
			// explicitly set it to "N/A".
			row.MetricValues[i] = "N/A"
		}
	}
	return row
}

// buildMetricsTableForClassVM constructs the view model for the metrics table.
// It collects all methods from all files within the class and sorts them
// primarily by file path, then by line number, then by short method name.
func (b *HtmlReportBuilder) buildMetricsTableForClassVM(classModel *model.Class) MetricsTableViewModel {
	metricsTable := MetricsTableViewModel{}
	metricsTable.Headers = b.getStandardMetricHeaders()

	if len(classModel.Methods) == 0 && len(classModel.Files) == 0 { // Check if there are any files to iterate
		return metricsTable
	}

	// Create a temporary struct to hold methods along with their file context for sorting
	type methodWithFileContext struct {
		method         *model.Method
		filePath       string // Full path for primary sort
		fileShortPath  string // For linking
		fileIndexPlus1 int    // For display in multi-file scenarios
	}
	var allMethodsWithContext []methodWithFileContext

	// Sort files first to ensure consistent file indexing and path usage
	sortedFiles := make([]model.CodeFile, len(classModel.Files))
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
		// However, model.Method doesn't directly link back to a specific model.CodeFile.
		// We need to find the methods that are defined within this specific file.
		// The model.Method.FirstLine and model.Method.DisplayName are key.
		// model.CodeFile.CodeElements helps link method display names to file lines.

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
		var correspondingCE *model.CodeElement
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
			correspondingCE = &model.CodeElement{
				Name:          mCtx.method.DisplayName, // Use display name as short name for this fallback
				FullName:      mCtx.method.DisplayName,
				Type:          model.MethodElementType,
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

func (b *HtmlReportBuilder) formatMetricValue(metric model.Metric) string {
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
		return utils.FormatPercentage(valFloat, b.maximumDecimalPlacesForPercentageDisplay)
	case "CrapScore":
		return fmt.Sprintf("%.2f", valFloat)
	case "Cyclomatic complexity", "Complexity":
		return fmt.Sprintf("%.0f", valFloat)
	default:
		return fmt.Sprintf(fmt.Sprintf("%%.%df", b.maximumDecimalPlacesForCoverageQuotas), valFloat)
	}
}

func (b *HtmlReportBuilder) buildAngularClassDetailForJS(classModel *model.Class, classVMServer *ClassViewModelForDetail) (AngularClassDetailViewModel, error) {
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
		angularFileForJS, err := b.buildAngularFileViewModelForJS(&fileInClass)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error building Angular file view model for JS (%s): %v\n", fileInClass.Path, err)
			continue
		}
		detailVM.Files = append(detailVM.Files, angularFileForJS)
	}
	return detailVM, nil
}

func (b *HtmlReportBuilder) buildAngularFileViewModelForJS(fileInClass *model.CodeFile) (AngularCodeFileViewModel, error) {
	angularFile := AngularCodeFileViewModel{
		Path:           fileInClass.Path,
		CoveredLines:   fileInClass.CoveredLines,
		CoverableLines: fileInClass.CoverableLines,
		TotalLines:     fileInClass.TotalLines,
		Lines:          []AngularLineAnalysisViewModel{},
	}
	sourceLines, err := filereader.ReadLinesInFile(fileInClass.Path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not read source file %s for JS Angular VM: %v\n", fileInClass.Path, err)
		return angularFile, nil
	}
	coverageLinesMap := make(map[int]*model.Line)
	if fileInClass.Lines != nil {
		for i := range fileInClass.Lines {
			covLine := &fileInClass.Lines[i]
			coverageLinesMap[covLine.Number] = covLine
		}
	}
	for i, content := range sourceLines {
		actualLineNumber := i + 1
		modelCovLine, hasCoverageData := coverageLinesMap[actualLineNumber]
		angularLine := b.buildAngularLineViewModelForJS(content, actualLineNumber, modelCovLine, hasCoverageData)
		angularFile.Lines = append(angularFile.Lines, angularLine)
	}
	return angularFile, nil
}

func (b *HtmlReportBuilder) buildAngularLineViewModelForJS(content string, actualLineNumber int, modelCovLine *model.Line, hasCoverageData bool) AngularLineAnalysisViewModel {
	lineVM := AngularLineAnalysisViewModel{
		LineNumber:  actualLineNumber,
		LineContent: content,
	}
	if hasCoverageData {
		lineVM.Hits = modelCovLine.Hits
		lineVM.CoveredBranches = modelCovLine.CoveredBranches
		lineVM.TotalBranches = modelCovLine.TotalBranches
		lineVM.LineVisitStatus = lineVisitStatusToString(modelCovLine.LineVisitStatus) // Use the field here
	} else {
		lineVM.LineVisitStatus = lineVisitStatusToString(model.NotCoverable) // Use model.NotCoverable
	}
	return lineVM
}

func (b *HtmlReportBuilder) buildClassDetailPageData(classVM ClassViewModelForDetail, tag string, classDetailJS template.JS) ClassDetailData {
	appVersion := "0.0.1"
	if b.ReportContext.ReportConfiguration() != nil {
		appVersion = "0.0.1"
	}
	return ClassDetailData{
		ReportTitle:                           b.reportTitle,
		AppVersion:                            appVersion,
		CurrentDateTime:                       time.Now().Format("02/01/2006 - 15:04:05"),
		Class:                                 classVM,
		BranchCoverageAvailable:               b.branchCoverageAvailable,
		MethodCoverageAvailable:               b.methodCoverageAvailable,
		Tag:                                   tag,
		Translations:                          b.translations,
		MaximumDecimalPlacesForCoverageQuotas: b.maximumDecimalPlacesForCoverageQuotas,
		AngularCssFile:                        b.angularCssFile,
		CombinedAngularJsFile:                 b.combinedAngularJsFile,
		AngularRuntimeJsFile:                  b.angularRuntimeJsFile,
		AngularPolyfillsJsFile:                b.angularPolyfillsJsFile,
		AngularMainJsFile:                     b.angularMainJsFile,
		AssembliesJSON:                        b.assembliesJSON,
		RiskHotspotsJSON:                      b.riskHotspotsJSON,
		MetricsJSON:                           b.metricsJSON,
		RiskHotspotMetricsJSON:                b.riskHotspotMetricsJSON,
		HistoricCoverageExecutionTimesJSON:    b.historicCoverageExecutionTimesJSON,
		TranslationsJSON:                      b.translationsJSON,
		ClassDetailJSON:                       classDetailJS,
	}
}
