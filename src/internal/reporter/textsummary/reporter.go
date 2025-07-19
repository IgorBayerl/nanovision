// In: internal/reporter/textsummary/reporter.go
package textsummary

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"text/tabwriter"
	"time"

	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/model"
	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/reporter"
	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/utils"
)

// TextReportBuilder generates a text summary report.
type TextReportBuilder struct {
	outputDir string
	logger    *slog.Logger
}

// NewTextReportBuilder creates a new TextReportBuilder.
func NewTextReportBuilder(outputDir string, logger *slog.Logger) reporter.ReportBuilder {
	return &TextReportBuilder{
		outputDir: outputDir,
		logger:    logger,
	}
}

// ReportType returns the type of report this builder generates.
func (b *TextReportBuilder) ReportType() string {
	return "TextSummary"
}

type summaryFileWriter struct {
	f *os.File
}

func (sfw *summaryFileWriter) writeLine(format string, args ...interface{}) {
	fmt.Fprintf(sfw.f, format+"\n", args...)
}

// CreateReport generates the text summary report using the analyzed model.SummaryResult.
func (b *TextReportBuilder) CreateReport(summary *model.SummaryResult) error {
	if err := os.MkdirAll(b.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	outputPath := filepath.Join(b.outputDir, "Summary.txt")
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create report file: %w", err)
	}
	defer f.Close()

	b.logger.Info("Writing text summary to file", "path", outputPath)

	sfw := &summaryFileWriter{f: f}

	decimalPlaces := 1                     // Placeholder, should be from settings
	decimalPlacesForPercentageDisplay := 0 // Placeholder, should be from settings

	sfw.writeLine("Summary")
	sfw.writeLine("  Generated on: %s", time.Now().Format("02/01/2006 - 15:04:05"))

	if summary.Timestamp > 0 {
		sfw.writeLine("  Coverage date: %s", time.Unix(summary.Timestamp, 0).Format("02/01/2006 - 15:04:05"))
	}

	sfw.writeLine("  Parser: %s", summary.ParserName)

	totalClasses := 0
	totalFiles := 0
	processedFilePaths := make(map[string]bool)
	for _, assembly := range summary.Assemblies {
		totalClasses += len(assembly.Classes)
		for _, class := range assembly.Classes {
			for _, codeFile := range class.Files {
				if !processedFilePaths[codeFile.Path] {
					processedFilePaths[codeFile.Path] = true
					totalFiles++
				}
			}
		}
	}

	sfw.writeLine("  Assemblies: %d", len(summary.Assemblies))
	sfw.writeLine("  Classes: %d", totalClasses)
	sfw.writeLine("  Files: %d", totalFiles)

	overallLineCoverage := utils.CalculatePercentage(summary.LinesCovered, summary.LinesValid, decimalPlaces)
	sfw.writeLine("  Line coverage: %s", utils.FormatPercentage(overallLineCoverage, decimalPlacesForPercentageDisplay))
	sfw.writeLine("  Covered lines: %d", summary.LinesCovered)
	sfw.writeLine("  Uncovered lines: %d", summary.LinesValid-summary.LinesCovered)
	sfw.writeLine("  Coverable lines: %d", summary.LinesValid)
	if summary.TotalLines > 0 {
		sfw.writeLine("  Total lines: %d", summary.TotalLines)
	} else {
		sfw.writeLine("  Total lines: N/A")
	}

	if summary.BranchesValid != nil && summary.BranchesCovered != nil {
		overallBranchCoverage := utils.CalculatePercentage(*summary.BranchesCovered, *summary.BranchesValid, decimalPlaces)
		// Only print percentage if there are valid branches (CalculatePercentage returns NaN if total is 0)
		if *summary.BranchesValid > 0 {
			sfw.writeLine("  Branch coverage: %s (%d of %d)", utils.FormatPercentage(overallBranchCoverage, decimalPlacesForPercentageDisplay), *summary.BranchesCovered, *summary.BranchesValid)
		} else { // No valid branches, just print counts or N/A for percentage
			sfw.writeLine("  Branch coverage: N/A (%d of %d)", *summary.BranchesCovered, *summary.BranchesValid)
		}
		sfw.writeLine("  Covered branches: %d", *summary.BranchesCovered)
		sfw.writeLine("  Total branches: %d", *summary.BranchesValid)
	}

	totalMethodsAgg, coveredMethodsAgg, fullyCoveredMethodsAgg := 0, 0, 0
	for _, assembly := range summary.Assemblies {
		for _, class := range assembly.Classes {
			// Assuming model.Class now has these pre-aggregated from analyzer/class.go
			totalMethodsAgg += class.TotalMethods
			coveredMethodsAgg += class.CoveredMethods
			fullyCoveredMethodsAgg += class.FullyCoveredMethods
		}
	}
	methodCoverage := utils.CalculatePercentage(coveredMethodsAgg, totalMethodsAgg, decimalPlaces)
	fullMethodCoverage := utils.CalculatePercentage(fullyCoveredMethodsAgg, totalMethodsAgg, decimalPlaces)

	sfw.writeLine("  Method coverage: %s (%d of %d)", utils.FormatPercentage(methodCoverage, decimalPlacesForPercentageDisplay), coveredMethodsAgg, totalMethodsAgg)
	sfw.writeLine("  Full method coverage: %s (%d of %d)", utils.FormatPercentage(fullMethodCoverage, decimalPlacesForPercentageDisplay), fullyCoveredMethodsAgg, totalMethodsAgg)
	sfw.writeLine("  Covered methods: %d", coveredMethodsAgg)
	sfw.writeLine("  Fully covered methods: %d", fullyCoveredMethodsAgg)
	sfw.writeLine("  Total methods: %d", totalMethodsAgg)

	tw := tabwriter.NewWriter(f, 0, 0, 2, ' ', 0)
	defer tw.Flush()
	for _, assembly := range summary.Assemblies {
		fmt.Fprintln(tw)
		assemblyLineCoverage := utils.CalculatePercentage(assembly.LinesCovered, assembly.LinesValid, decimalPlaces)
		fmt.Fprintf(tw, "%s\t  %s\n", assembly.Name, utils.FormatPercentage(assemblyLineCoverage, decimalPlacesForPercentageDisplay))

		sortedClasses := make([]model.Class, len(assembly.Classes))
		copy(sortedClasses, assembly.Classes)
		sort.Slice(sortedClasses, func(i, j int) bool {
			return sortedClasses[i].DisplayName < sortedClasses[j].DisplayName
		})
		for _, class := range sortedClasses {
			classLineCoverage := utils.CalculatePercentage(class.LinesCovered, class.LinesValid, decimalPlaces)
			fmt.Fprintf(tw, "  %s\t  %s\n", class.DisplayName, utils.FormatPercentage(classLineCoverage, decimalPlacesForPercentageDisplay))
		}
	}
	return nil
}
