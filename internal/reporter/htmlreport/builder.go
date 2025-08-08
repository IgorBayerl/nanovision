package htmlreport

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/IgorBayerl/AdlerCov/internal/filereader"
	"github.com/IgorBayerl/AdlerCov/internal/model"
	"github.com/IgorBayerl/AdlerCov/internal/reporter"
)

type HtmlReportBuilder struct {
	OutputDir     string
	ReportContext reporter.IBuilderContext
	fileReader    filereader.Reader

	// Cached data
	angularCssFile                     string
	combinedAngularJsFile              string
	assembliesJSON                     template.JS
	riskHotspotsJSON                   template.JS
	metricsJSON                        template.JS
	riskHotspotMetricsJSON             template.JS
	historicCoverageExecutionTimesJSON template.JS
	translationsJSON                   template.JS

	// Settings derived from context
	branchCoverageAvailable bool
	methodCoverageAvailable bool
	parserName              string
	reportTimestamp         int64
	reportTitle             string
	tag                     string
	translations            map[string]string
	onlySummary             bool

	classReportFilenames       map[string]string
	tempExistingLowerFilenames map[string]struct{}
}

func NewHtmlReportBuilder(outputDir string, reportCtx reporter.IBuilderContext, fileReader filereader.Reader) *HtmlReportBuilder {
	return &HtmlReportBuilder{
		OutputDir:                  outputDir,
		ReportContext:              reportCtx,
		fileReader:                 fileReader,
		classReportFilenames:       make(map[string]string),
		tempExistingLowerFilenames: make(map[string]struct{}),
	}
}

func (b *HtmlReportBuilder) ReportType() string {
	return "Html"
}

func (b *HtmlReportBuilder) CreateReport(tree *model.SummaryTree) error {
	report := ToLegacySummaryResult(tree, b.fileReader, b.ReportContext.Logger())

	if err := b.validateContext(); err != nil {
		return err
	}
	if err := b.prepareOutputDirectory(); err != nil {
		return err
	}
	if err := b.initializeAssets(); err != nil {
		return err
	}

	b.initializeBuilderProperties(report)
	if err := b.prepareGlobalJSONData(report); err != nil {
		return err
	}

	angularAssembliesForSummary, err := b.buildAngularAssemblyViewModelsForSummary(report)
	if err != nil {
		return fmt.Errorf("failed to build angular assembly view models for summary: %w", err)
	}

	var angularRiskHotspots []AngularRiskHotspotViewModel
	if err := b.setRiskHotspotsJSON(angularRiskHotspots); err != nil {
		return err
	}

	summaryData, err := b.buildSummaryPageData(report, angularAssembliesForSummary, angularRiskHotspots)
	if err != nil {
		return fmt.Errorf("failed to build summary page data: %w", err)
	}
	if err := b.renderSummaryPage(summaryData); err != nil {
		return fmt.Errorf("failed to render summary page: %w", err)
	}

	if !b.onlySummary {
		if err := b.renderClassDetailPages(report); err != nil {
			return fmt.Errorf("failed to render class detail pages: %w", err)
		}
	}
	return nil
}

func (b *HtmlReportBuilder) validateContext() error {
	if b.ReportContext == nil {
		return fmt.Errorf("HtmlReportBuilder.ReportContext is not set")
	}
	return nil
}

func (b *HtmlReportBuilder) prepareOutputDirectory() error {
	return os.MkdirAll(b.OutputDir, 0755)
}

func (b *HtmlReportBuilder) initializeBuilderProperties(report *SummaryResult) {
	appConfig := b.ReportContext.Config()

	b.reportTitle = appConfig.Title
	if b.reportTitle == "" {
		b.reportTitle = "Summary"
	}
	b.parserName = report.ParserName
	b.reportTimestamp = report.Timestamp
	b.tag = appConfig.Tag
	b.branchCoverageAvailable = report.BranchesValid != nil && *report.BranchesValid > 0
	b.methodCoverageAvailable = true
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

func (b *HtmlReportBuilder) renderClassDetailPages(report *SummaryResult) error {
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

func (b *HtmlReportBuilder) determineClassReportFilename(assemblyName string, className string, assemblyShortNameForFile string) string {
	classKey := assemblyName + "_" + className

	if filename, ok := b.classReportFilenames[classKey]; ok {
		return filename
	}

	newFilename := generateUniqueFilename(className, b.tempExistingLowerFilenames)
	b.classReportFilenames[classKey] = newFilename
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
