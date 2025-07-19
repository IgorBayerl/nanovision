package htmlreport

import (
	"fmt" // fmt is still needed for fmt.Errorf
	"html/template"
	"os"
	"path/filepath"

	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/model"
	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/reporter"
)

type HtmlReportBuilder struct {
	OutputDir     string
	ReportContext reporter.IBuilderContext

	// Cached data for reuse across page generations
	angularCssFile                     string
	angularRuntimeJsFile               string
	angularPolyfillsJsFile             string
	angularMainJsFile                  string
	assembliesJSON                     template.JS
	riskHotspotsJSON                   template.JS
	metricsJSON                        template.JS
	riskHotspotMetricsJSON             template.JS
	historicCoverageExecutionTimesJSON template.JS
	translationsJSON                   template.JS

	// Settings derived from context
	branchCoverageAvailable                  bool
	methodCoverageAvailable                  bool
	maximumDecimalPlacesForCoverageQuotas    int
	maximumDecimalPlacesForPercentageDisplay int
	parserName                               string
	reportTimestamp                          int64
	reportTitle                              string
	tag                                      string
	translations                             map[string]string
	onlySummary                              bool

	classReportFilenames       map[string]string
	tempExistingLowerFilenames map[string]struct{}

	combinedAngularJsFile string // To store "reportgenerator.combined.js"
}

func NewHtmlReportBuilder(outputDir string, reportCtx reporter.IBuilderContext) *HtmlReportBuilder {
	return &HtmlReportBuilder{
		OutputDir:                  outputDir,
		ReportContext:              reportCtx,
		classReportFilenames:       make(map[string]string),
		tempExistingLowerFilenames: make(map[string]struct{}),
	}
}

func (b *HtmlReportBuilder) ReportType() string {
	return "Html"
}

func (b *HtmlReportBuilder) CreateReport(report *model.SummaryResult) error {
	if err := b.validateContext(); err != nil {
		return err
	}
	if err := b.prepareOutputDirectory(); err != nil {
		return err
	}
	if err := b.initializeAssets(); err != nil { // Copies static assets and parses Angular index.html
		return err
	}

	b.initializeBuilderProperties(report)                   // Sets up common properties like title, translations etc.
	if err := b.prepareGlobalJSONData(report); err != nil { // Prepares metricsJSON, riskHotspotMetricsJSON etc.
		return err
	}

	// This call will populate b.classReportFilenames and b.assembliesJSON
	// The returned 'angularAssemblies' is not strictly needed here if all subsequent
	// operations use the builder's stored JSON or maps.
	// However, buildSummaryPageData might still conceptually want the processed view models.
	// Let's keep it for buildSummaryPageData if that function's logic benefits from it.
	angularAssembliesForSummary, err := b.buildAngularAssemblyViewModelsForSummary(report)
	if err != nil {
		return fmt.Errorf("failed to build angular assembly view models for summary: %w", err)
	}

	var angularRiskHotspots []AngularRiskHotspotViewModel              // Placeholder
	if err := b.setRiskHotspotsJSON(angularRiskHotspots); err != nil { // Prepares b.riskHotspotsJSON
		return err
	}

	// Pass angularAssembliesForSummary to buildSummaryPageData
	summaryData, err := b.buildSummaryPageData(report, angularAssembliesForSummary, angularRiskHotspots)
	if err != nil {
		return fmt.Errorf("failed to build summary page data: %w", err)
	}
	if err := b.renderSummaryPage(summaryData); err != nil {
		return fmt.Errorf("failed to render summary page: %w", err)
	}

	if !b.onlySummary {
		// renderClassDetailPages uses b.classReportFilenames, so it doesn't need angularAssembliesForSummary
		if err := b.renderClassDetailPages(report); err != nil {
			return fmt.Errorf("failed to render class detail pages: %w", err)
		}
	}
	return nil
}

// --- CreateReport helper methods ---

func (b *HtmlReportBuilder) validateContext() error {
	if b.ReportContext == nil {
		return fmt.Errorf("HtmlReportBuilder.ReportContext is not set; it's required for configuration and settings")
	}
	return nil
}

func (b *HtmlReportBuilder) prepareOutputDirectory() error {
	return os.MkdirAll(b.OutputDir, 0755)
}

func (b *HtmlReportBuilder) initializeBuilderProperties(report *model.SummaryResult) {
	reportConfig := b.ReportContext.ReportConfiguration()
	settings := b.ReportContext.Settings()

	b.reportTitle = reportConfig.Title()
	if b.reportTitle == "" {
		b.reportTitle = "Summary" // Default for summary page
	}
	b.parserName = report.ParserName
	b.reportTimestamp = report.Timestamp
	b.tag = reportConfig.Tag()
	b.branchCoverageAvailable = report.BranchesValid != nil && *report.BranchesValid > 0
	b.methodCoverageAvailable = true
	b.maximumDecimalPlacesForCoverageQuotas = settings.MaximumDecimalPlacesForCoverageQuotas
	b.maximumDecimalPlacesForPercentageDisplay = settings.MaximumDecimalPlacesForPercentageDisplay
	b.translations = GetTranslations()
}

func (b *HtmlReportBuilder) renderSummaryPage(data SummaryPageData) error {
	outputIndexPath := filepath.Join(b.OutputDir, "index.html")
	summaryFile, err := os.Create(outputIndexPath)
	if err != nil {
		return fmt.Errorf("failed to create index.html: %w", err)
	}
	defer summaryFile.Close()
	return summaryPageTpl.Execute(summaryFile, data)
}

func (b *HtmlReportBuilder) renderClassDetailPages(report *model.SummaryResult) error { // Removed angularAssembliesForSummary
	if b.onlySummary {
		return nil
	}

	for _, assemblyModel := range report.Assemblies {
		for _, classModel := range assemblyModel.Classes {
			classKey := assemblyModel.Name + "_" + classModel.Name

			classReportFilename, ok := b.classReportFilenames[classKey]

			if !ok || classReportFilename == "" {
				b.ReportContext.Logger().Error(
					"Class report filename not found, skipping detail page generation",
					"class", classModel.DisplayName,
					"assembly", assemblyModel.Name,
				)
				continue
			}

			err := b.generateClassDetailHTML(&classModel, classReportFilename, b.tag)
			if err != nil {
				b.ReportContext.Logger().Error(
					"Failed to generate detail page for class",
					"class", classModel.DisplayName,
					"file", classReportFilename,
					"error", err,
				)
			}
		}
	}
	return nil
}

// determineClassReportFilename gets or generates a unique HTML filename for a class report.
// It uses and updates the builder's internal maps for filename tracking.
// assemblyName is the full assembly name, className is the model's raw/unique name.
func (b *HtmlReportBuilder) determineClassReportFilename(assemblyName string, className string, assemblyShortNameForFile string) string {
	// Create a unique key for the class within its assembly.
	// Using the full assembly name and raw class name for the key ensures uniqueness.
	classKey := assemblyName + "_" + className

	if filename, ok := b.classReportFilenames[classKey]; ok {
		return filename // Return already generated filename
	}

	// Filename not yet generated for this class. Generate a new one.
	// generateUniqueFilename expects a map of existing *lowercase* filenames.
	// b.tempExistingLowerFilenames serves this purpose.
	// assemblyShortNameForFile is used for constructing the base of the filename.
	newFilename := generateUniqueFilename(assemblyShortNameForFile, className, b.tempExistingLowerFilenames)

	// Store the actual generated filename (preserving case) in classReportFilenames.
	b.classReportFilenames[classKey] = newFilename
	// Also, add its lowercase version to tempExistingLowerFilenames for future uniqueness checks by generateUniqueFilename.
	// Note: generateUniqueFilename itself adds to the map it's passed, so this is already handled if it modifies its input map.
	// The current generateUniqueFilename modifies the map passed to it.

	return newFilename
}

func (b *HtmlReportBuilder) renderClassDetailPage(data ClassDetailData, classReportFilename string) error {
	outputFilePath := filepath.Join(b.OutputDir, classReportFilename)
	fileWriter, err := os.Create(outputFilePath)
	if err != nil {
		return fmt.Errorf("failed to create class report file %s: %w", outputFilePath, err)
	}
	defer fileWriter.Close()
	return classDetailTpl.Execute(fileWriter, data)
}
