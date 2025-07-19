package htmlreport

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/model"
	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/utils"
)

func (b *HtmlReportBuilder) prepareGlobalJSONData(report *model.SummaryResult) error {
	translationsJSONBytes, err := json.Marshal(b.translations)
	if err != nil {
		b.translationsJSON = template.JS("({})") // Fallback
	} else {
		b.translationsJSON = template.JS(string(translationsJSONBytes)) // Ensure it's string(bytes)
	}

	availableMetrics := []AngularMetricViewModel{
		{Name: "NPath complexity", Abbreviation: "npath", ExplanationURL: "https://modess.io/npath-complexity-cyclomatic-complexity-explained/"},
		{Name: "CrapScore", Abbreviation: "crap", ExplanationURL: "https://testing.googleblog.com/2011/02/this-code-is-crap.html"},
	}
	metricsJSONBytes, err := json.Marshal(availableMetrics)
	if err != nil {
		b.metricsJSON = template.JS("([])")
	} else {
		b.metricsJSON = template.JS(string(metricsJSONBytes))
	}

	riskHotspotMetricHeaders := []AngularRiskHotspotMetricHeaderViewModel{
		{Name: "Cyclomatic complexity", Abbreviation: "cyclomatic", ExplanationURL: "https://www.ndepend.com/docs/code-metrics#CC"},
		{Name: "CrapScore", Abbreviation: "crap", ExplanationURL: "https://testing.googleblog.com/2011/02/this-code-is-crap.html"},
		{Name: "NPath complexity", Abbreviation: "npath", ExplanationURL: "https://modess.io/npath-complexity-cyclomatic-complexity-explained/"},
	}
	riskHotspotMetricsJSONBytes, err := json.Marshal(riskHotspotMetricHeaders)
	if err != nil {
		b.riskHotspotMetricsJSON = template.JS("([])")
	} else {
		b.riskHotspotMetricsJSON = template.JS(string(riskHotspotMetricsJSONBytes))
	}

	executionTimes := b.collectHistoricExecutionTimes(report)
	historicExecTimesJSONBytes, err := json.Marshal(executionTimes)
	if err != nil {
		b.historicCoverageExecutionTimesJSON = template.JS("([])")
	} else {
		b.historicCoverageExecutionTimesJSON = template.JS(string(historicExecTimesJSONBytes))
	}
	return nil
}

func (b *HtmlReportBuilder) collectHistoricExecutionTimes(report *model.SummaryResult) []string {
	var allHistoricCoverages []model.HistoricCoverage
	if report.Assemblies != nil {
		for _, assembly := range report.Assemblies {
			for _, class := range assembly.Classes {
				allHistoricCoverages = append(allHistoricCoverages, class.HistoricCoverages...)
			}
		}
	}

	if len(allHistoricCoverages) == 0 {
		return []string{}
	}

	distinctHistoricCoverages := utils.DistinctBy(allHistoricCoverages, func(hc model.HistoricCoverage) int64 {
		return hc.ExecutionTime
	})

	executionTimes := make([]string, len(distinctHistoricCoverages))
	for i, hc := range distinctHistoricCoverages {
		executionTimes[i] = time.Unix(hc.ExecutionTime, 0).Format("2006-01-02 15:04:05")
	}
	sort.Strings(executionTimes) // Ensure consistent order
	return executionTimes
}

func (b *HtmlReportBuilder) buildAngularAssemblyViewModelsForSummary(report *model.SummaryResult) ([]AngularAssemblyViewModel, error) {
	var angularAssemblies []AngularAssemblyViewModel
	if len(report.Assemblies) == 0 {
		log.Println("buildAngularAssemblyViewModelsForSummary: No assemblies in the report or report.Assemblies is nil.")
		b.assembliesJSON = template.JS("[]") // Default to a valid empty JS array literal
		return angularAssemblies, nil
	}

	log.Printf("buildAngularAssemblyViewModelsForSummary: Processing %d assemblies.\n", len(report.Assemblies))

	for _, assembly := range report.Assemblies {
		var assemblyShortNameForFile string
		if lastSlash := strings.LastIndexAny(assembly.Name, "/\\"); lastSlash != -1 {
			assemblyShortNameForFile = assembly.Name[lastSlash+1:]
		} else {
			assemblyShortNameForFile = assembly.Name
		}

		angularAssembly := AngularAssemblyViewModel{Name: assembly.Name, Classes: []AngularClassViewModel{}}
		log.Printf("  Processing Assembly: %s (ShortName for file: %s)\n", assembly.Name, assemblyShortNameForFile)

		if len(assembly.Classes) == 0 {
			log.Printf("    Assembly %s has no classes.\n", assembly.Name)
		}

		for _, class := range assembly.Classes {
			classReportFilename := b.determineClassReportFilename(assembly.Name, class.Name, assemblyShortNameForFile)
			log.Printf("    Processing Class: %s, ReportPath: %s\n", class.DisplayName, classReportFilename)

			angularClass := b.buildAngularClassViewModelForSummary(&class, classReportFilename)
			log.Printf("      AngularClass Built: Name=%s, RP=%s, CL=%d, CAL=%d\n", angularClass.Name, angularClass.ReportPath, angularClass.CoveredLines, angularClass.CoverableLines)
			angularAssembly.Classes = append(angularAssembly.Classes, angularClass)
		}
		angularAssemblies = append(angularAssemblies, angularAssembly)
	}

	if len(angularAssemblies) == 0 {
		log.Println("buildAngularAssemblyViewModelsForSummary: angularAssemblies slice is empty after processing (this shouldn't happen if report.Assemblies was not empty).")
		b.assembliesJSON = template.JS("[]")
		return angularAssemblies, nil
	}

	assembliesJSONBytes, err := json.Marshal(angularAssemblies)
	if err != nil {
		log.Printf("buildAngularAssemblyViewModelsForSummary: ERROR marshaling angularAssemblies to JSON: %v\n", err)
		b.assembliesJSON = template.JS("[]") // Fallback
		return nil, fmt.Errorf("failed to marshal angular assemblies for summary: %w", err)
	}

	jsonString := string(assembliesJSONBytes)
	log.Printf("buildAngularAssemblyViewModelsForSummary: Marshaled assembliesJSON (length: %d): %s\n", len(jsonString), jsonString)
	if jsonString == "null" { // Safeguard, though Marshal on non-empty slice shouldn't give "null"
		log.Println("buildAngularAssemblyViewModelsForSummary: Marshaled JSON is 'null', ensuring it's an empty array '[]' for JS.")
		b.assembliesJSON = template.JS("[]")
	} else {
		b.assembliesJSON = template.JS(jsonString) // Key: assign the string to template.JS
	}
	return angularAssemblies, nil
}

func (b *HtmlReportBuilder) buildAngularClassViewModelForSummary(class *model.Class, reportPath string) AngularClassViewModel {
	angularClass := AngularClassViewModel{
		Name:                      class.DisplayName,
		ReportPath:                reportPath,
		CoveredLines:              class.LinesCovered,
		UncoveredLines:            class.LinesValid - class.LinesCovered,
		CoverableLines:            class.LinesValid,
		TotalLines:                class.TotalLines,
		Metrics:                   make(map[string]float64),
		HistoricCoverages:         []AngularHistoricCoverageViewModel{},
		LineCoverageHistory:       []float64{},
		BranchCoverageHistory:     []float64{},
		MethodCoverageHistory:     []float64{},
		FullMethodCoverageHistory: []float64{},
	}

	angularClass.TotalMethods = class.TotalMethods
	angularClass.CoveredMethods = class.CoveredMethods
	angularClass.FullyCoveredMethods = class.FullyCoveredMethods

	if class.BranchesCovered != nil {
		angularClass.CoveredBranches = *class.BranchesCovered
	} else {
		angularClass.CoveredBranches = 0
	}
	if class.BranchesValid != nil {
		angularClass.TotalBranches = *class.BranchesValid
	} else {
		angularClass.TotalBranches = 0
	}

	for _, hist := range class.HistoricCoverages {
		angularHist := b.buildAngularHistoricCoverageViewModel(&hist)
		angularClass.HistoricCoverages = append(angularClass.HistoricCoverages, angularHist)

		if angularHist.LineCoverageQuota >= 0 {
			angularClass.LineCoverageHistory = append(angularClass.LineCoverageHistory, angularHist.LineCoverageQuota)
		}
		if angularHist.BranchCoverageQuota >= 0 {
			angularClass.BranchCoverageHistory = append(angularClass.BranchCoverageHistory, angularHist.BranchCoverageQuota)
		}
		if angularHist.MethodCoverageQuota >= 0 {
			angularClass.MethodCoverageHistory = append(angularClass.MethodCoverageHistory, angularHist.MethodCoverageQuota)
		}
		if angularHist.FullMethodCoverageQuota >= 0 {
			angularClass.FullMethodCoverageHistory = append(angularClass.FullMethodCoverageHistory, angularHist.FullMethodCoverageQuota)
		}
	}

	for name, val := range class.Metrics {
		angularClass.Metrics[name] = val
	}

	return angularClass
}

func (b *HtmlReportBuilder) buildAngularHistoricCoverageViewModel(hist *model.HistoricCoverage) AngularHistoricCoverageViewModel {
	angularHist := AngularHistoricCoverageViewModel{
		ExecutionTime:   time.Unix(hist.ExecutionTime, 0).Format("2006-01-02"), // Simplified
		CoveredLines:    hist.CoveredLines,
		CoverableLines:  hist.CoverableLines,
		TotalLines:      hist.TotalLines,
		CoveredBranches: hist.CoveredBranches,
		TotalBranches:   hist.TotalBranches,
		// FIXME: Populate these from your model.HistoricCoverage if it has method coverage history
		// CoveredMethods: hist.CoveredMethods,
		// FullyCoveredMethods: hist.FullyCoveredMethods,
		// TotalMethods: hist.TotalMethods,
	}

	angularHist.LineCoverageQuota = -1.0
	if hist.CoverableLines > 0 {
		angularHist.LineCoverageQuota = (float64(hist.CoveredLines) / float64(hist.CoverableLines)) * 100.0
	}

	angularHist.BranchCoverageQuota = -1.0
	if hist.TotalBranches > 0 {
		angularHist.BranchCoverageQuota = (float64(hist.CoveredBranches) / float64(hist.TotalBranches)) * 100.0
	}

	// FIXME: Update these based on fields in your model.HistoricCoverage
	angularHist.MethodCoverageQuota = -1.0
	angularHist.FullMethodCoverageQuota = -1.0
	// if hist.TotalMethods > 0 {
	// 	angularHist.MethodCoverageQuota = (float64(hist.CoveredMethods) / float64(hist.TotalMethods)) * 100.0
	// 	angularHist.FullMethodCoverageQuota = (float64(hist.FullyCoveredMethods) / float64(hist.TotalMethods)) * 100.0
	// }

	return angularHist
}

func (b *HtmlReportBuilder) setRiskHotspotsJSON(angularRiskHotspots []AngularRiskHotspotViewModel) error {
	log.Printf("setRiskHotspotsJSON: Received %d risk hotspots to marshal.", len(angularRiskHotspots))

	// json.Marshal on a nil slice results in "null" string.
	// json.Marshal on an empty non-nil slice (e.g., make([]Type, 0)) results in "[]".
	// Angular expects an array.
	if angularRiskHotspots == nil { // Explicitly handle nil slice
		b.riskHotspotsJSON = template.JS("[]")
		log.Println("setRiskHotspotsJSON: angularRiskHotspots was nil, b.riskHotspotsJSON set to '[]'")
		return nil
	}

	riskHotspotsJSONBytes, err := json.Marshal(angularRiskHotspots)
	if err != nil {
		b.riskHotspotsJSON = template.JS("[]") // Fallback
		log.Printf("setRiskHotspotsJSON: Error marshaling risk hotspots: %v. b.riskHotspotsJSON set to '[]'", err)
		return fmt.Errorf("failed to marshal angular risk hotspots: %w", err)
	}

	jsonString := string(riskHotspotsJSONBytes)
	// If angularRiskHotspots was an empty (but not nil) slice, jsonString would be "[]".
	// If it was nil, jsonString would be "null". The 'if angularRiskHotspots == nil' above handles this.
	if jsonString == "null" {
		log.Println("setRiskHotspotsJSON: Marshaled riskHotspotsJSON is 'null', changing to '[]' for JS.")
		b.riskHotspotsJSON = template.JS("[]")
	} else {
		b.riskHotspotsJSON = template.JS(jsonString)
	}
	log.Printf("setRiskHotspotsJSON: b.riskHotspotsJSON set to: %s", b.riskHotspotsJSON)
	return nil
}

func (b *HtmlReportBuilder) buildSummaryPageData(report *model.SummaryResult, angularAssembliesForSummary []AngularAssemblyViewModel, angularRiskHotspots []AngularRiskHotspotViewModel) (SummaryPageData, error) {
	log.Printf("buildSummaryPageData: b.assembliesJSON before assigning to SummaryPageData: %s\n", b.assembliesJSON) // Log

	data := SummaryPageData{
		ReportTitle:                        b.reportTitle,
		AppVersion:                         "0.0.1",
		CurrentDateTime:                    time.Now().Format("02/01/2006 - 15:04:05"),
		Translations:                       b.translations,
		HasRiskHotspots:                    len(angularRiskHotspots) > 0,
		HasAssemblies:                      len(report.Assemblies) > 0,
		AssembliesJSON:                     b.assembliesJSON,
		RiskHotspotsJSON:                   b.riskHotspotsJSON,
		MetricsJSON:                        b.metricsJSON,
		RiskHotspotMetricsJSON:             b.riskHotspotMetricsJSON,
		HistoricCoverageExecutionTimesJSON: b.historicCoverageExecutionTimesJSON,
		TranslationsJSON:                   b.translationsJSON,
		AngularCssFile:                     b.angularCssFile,
		CombinedAngularJsFile:              b.combinedAngularJsFile,
		AngularRuntimeJsFile:               b.angularRuntimeJsFile,
		AngularPolyfillsJsFile:             b.angularPolyfillsJsFile,
		AngularMainJsFile:                  b.angularMainJsFile,

		BranchCoverageAvailable:               b.branchCoverageAvailable,
		MethodCoverageAvailable:               b.methodCoverageAvailable,
		MaximumDecimalPlacesForCoverageQuotas: b.maximumDecimalPlacesForCoverageQuotas,
		SummaryCards:                          b.buildSummaryCards(report),
		OverallHistoryChartData:               HistoryChartDataViewModel{Series: false},
	}
	return data, nil
}

func (b *HtmlReportBuilder) buildSummaryCards(report *model.SummaryResult) []CardViewModel {
	var cards []CardViewModel
	decimalPlaces := b.maximumDecimalPlacesForCoverageQuotas
	decimalPlacesForPercentageDisplay := b.maximumDecimalPlacesForPercentageDisplay

	// Information Card
	infoCardRows := []CardRowViewModel{
		{Header: b.translations["Parser"], Text: report.ParserName},
		{Header: b.translations["Assemblies2"], Text: fmt.Sprintf("%d", len(report.Assemblies)), Alignment: "right"},
		{Header: b.translations["Classes"], Text: fmt.Sprintf("%d", countTotalClasses(report.Assemblies)), Alignment: "right"},
		{Header: b.translations["Files2"], Text: fmt.Sprintf("%d", countUniqueFiles(report.Assemblies)), Alignment: "right"},
	}
	if report.Timestamp > 0 {
		infoCardRows = append(infoCardRows, CardRowViewModel{Header: b.translations["CoverageDate"], Text: time.Unix(report.Timestamp, 0).Format("02/01/2006 - 15:04:05")})
	}
	if b.tag != "" {
		infoCardRows = append(infoCardRows, CardRowViewModel{Header: b.translations["Tag"], Text: b.tag})
	}
	cards = append(cards, CardViewModel{Title: b.translations["Information"], Rows: infoCardRows})

	// Line Coverage Card
	lineCovQuota := utils.CalculatePercentage(report.LinesCovered, report.LinesValid, decimalPlaces)
	lineCovText := utils.FormatPercentage(lineCovQuota, decimalPlacesForPercentageDisplay)
	lineCovTooltip := "-"
	if !math.IsNaN(lineCovQuota) {
		lineCovTooltip = fmt.Sprintf("%d of %d", report.LinesCovered, report.LinesValid)
	}
	lineCovBar := 0
	if !math.IsNaN(lineCovQuota) {
		lineCovBar = 100 - int(math.Round(lineCovQuota))
	}

	cards = append(cards, CardViewModel{Title: b.translations["LineCoverage"], SubTitle: lineCovText, SubTitlePercentageBarValue: lineCovBar, Rows: []CardRowViewModel{
		{Header: b.translations["CoveredLines"], Text: fmt.Sprintf("%d", report.LinesCovered), Alignment: "right"},
		{Header: b.translations["UncoveredLines"], Text: fmt.Sprintf("%d", report.LinesValid-report.LinesCovered), Alignment: "right"},
		{Header: b.translations["CoverableLines"], Text: fmt.Sprintf("%d", report.LinesValid), Alignment: "right"},
		{Header: b.translations["TotalLines"], Text: fmt.Sprintf("%d", report.TotalLines), Alignment: "right"},
		{Header: b.translations["LineCoverage"], Text: lineCovText, Tooltip: lineCovTooltip, Alignment: "right"},
	}})

	// Branch Coverage Card (Conditional)
	if b.branchCoverageAvailable && report.BranchesCovered != nil && report.BranchesValid != nil {
		branchCovQuota := utils.CalculatePercentage(*report.BranchesCovered, *report.BranchesValid, decimalPlaces)
		branchCovText := utils.FormatPercentage(branchCovQuota, decimalPlacesForPercentageDisplay)
		branchCovTooltip := "-"
		if !math.IsNaN(branchCovQuota) {
			branchCovTooltip = fmt.Sprintf("%d of %d", *report.BranchesCovered, *report.BranchesValid)
		}
		branchCovBar := 0
		if !math.IsNaN(branchCovQuota) {
			branchCovBar = 100 - int(math.Round(branchCovQuota))
		}

		cards = append(cards, CardViewModel{Title: b.translations["BranchCoverage"], SubTitle: branchCovText, SubTitlePercentageBarValue: branchCovBar, Rows: []CardRowViewModel{
			{Header: b.translations["CoveredBranches2"], Text: fmt.Sprintf("%d", *report.BranchesCovered), Alignment: "right"},
			{Header: b.translations["TotalBranches"], Text: fmt.Sprintf("%d", *report.BranchesValid), Alignment: "right"},
			{Header: b.translations["BranchCoverage"], Text: branchCovText, Tooltip: branchCovTooltip, Alignment: "right"},
		}})
	}

	// Method Coverage Card
	var totalMethods, coveredMethods, fullyCoveredMethods int
	for _, asm := range report.Assemblies {
		for _, cls := range asm.Classes {
			totalMethods += cls.TotalMethods
			coveredMethods += cls.CoveredMethods
			fullyCoveredMethods += cls.FullyCoveredMethods
		}
	}
	methodCovQuota := utils.CalculatePercentage(coveredMethods, totalMethods, decimalPlaces)
	methodCovText := utils.FormatPercentage(methodCovQuota, decimalPlacesForPercentageDisplay)
	methodCovTooltip := "-"
	if !math.IsNaN(methodCovQuota) {
		methodCovTooltip = fmt.Sprintf("%d of %d", coveredMethods, totalMethods)
	}
	methodCovBar := 0
	if !math.IsNaN(methodCovQuota) {
		methodCovBar = 100 - int(math.Round(methodCovQuota))
	}

	fullMethodCovQuota := utils.CalculatePercentage(fullyCoveredMethods, totalMethods, decimalPlaces)
	fullMethodCovText := utils.FormatPercentage(fullMethodCovQuota, decimalPlacesForPercentageDisplay)
	fullMethodCovTooltip := "-"
	if !math.IsNaN(fullMethodCovQuota) {
		fullMethodCovTooltip = fmt.Sprintf("%d of %d", fullyCoveredMethods, totalMethods)
	}

	cards = append(cards, CardViewModel{
		Title: b.translations["MethodCoverage"], ProRequired: !b.methodCoverageAvailable, SubTitle: methodCovText, SubTitlePercentageBarValue: methodCovBar,
		Rows: []CardRowViewModel{
			{Header: b.translations["CoveredCodeElements"], Text: fmt.Sprintf("%d", coveredMethods), Alignment: "right"},
			{Header: b.translations["FullCoveredCodeElements"], Text: fmt.Sprintf("%d", fullyCoveredMethods), Alignment: "right"},
			{Header: b.translations["TotalCodeElements"], Text: fmt.Sprintf("%d", totalMethods), Alignment: "right"},
			{Header: b.translations["CodeElementCoverageQuota2"], Text: methodCovText, Tooltip: methodCovTooltip, Alignment: "right"},
			{Header: b.translations["FullCodeElementCoverageQuota2"], Text: fullMethodCovText, Tooltip: fullMethodCovTooltip, Alignment: "right"},
		},
	})
	return cards
}
