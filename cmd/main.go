// Path: cmd/main.go
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/IgorBayerl/fsglob"

	// New core components
	"github.com/IgorBayerl/AdlerCov/internal/analyzer"
	"github.com/IgorBayerl/AdlerCov/internal/model"

	// Foundational packages
	"github.com/IgorBayerl/AdlerCov/internal/filereader"
	"github.com/IgorBayerl/AdlerCov/internal/language"
	"github.com/IgorBayerl/AdlerCov/internal/logging"
	"github.com/IgorBayerl/AdlerCov/internal/parsers"
	"github.com/IgorBayerl/AdlerCov/internal/reportconfig"
	"github.com/IgorBayerl/AdlerCov/internal/reporter"
	"github.com/IgorBayerl/AdlerCov/internal/settings"
	"github.com/IgorBayerl/AdlerCov/internal/utils"

	// Language specific processors
	"github.com/IgorBayerl/AdlerCov/internal/language/cpp"
	"github.com/IgorBayerl/AdlerCov/internal/language/csharp"
	"github.com/IgorBayerl/AdlerCov/internal/language/defaultformatter"
	"github.com/IgorBayerl/AdlerCov/internal/language/golang"

	// Parsers
	"github.com/IgorBayerl/AdlerCov/internal/parsers/cobertura"
	"github.com/IgorBayerl/AdlerCov/internal/parsers/gcov"
	"github.com/IgorBayerl/AdlerCov/internal/parsers/gocover"

	// Reporters
	"github.com/IgorBayerl/AdlerCov/internal/reporter/htmlreport"
	"github.com/IgorBayerl/AdlerCov/internal/reporter/lcov"
	"github.com/IgorBayerl/AdlerCov/internal/reporter/textsummary"
)

var ErrMissingReportFlag = errors.New("missing required -report flag")

// cliFlags defines the command-line flags for the application.
// Legacy flags like -assemblyfilters and -classfilters have been removed
// in favor of a simpler, path-based filtering model.
type cliFlags struct {
	// domain
	reportsPatterns *string
	outputDir       *string
	reportTypes     *string
	sourceDirs      *string
	tag             *string
	title           *string
	fileFilters     *string

	// logging
	verbose   *bool
	verbosity *string
	logFile   *string
	logFormat *string
}

func parseFlags() (*cliFlags, error) {
	f := &cliFlags{
		// domain flags
		reportsPatterns: flag.String("report", "", "Coverage report file paths or patterns (semicolon-separated)"),
		outputDir:       flag.String("output", "coverage-report", "Output directory for generated reports"),
		reportTypes:     flag.String("reporttypes", "TextSummary,Html", "Report types (comma-separated)"),
		sourceDirs:      flag.String("sourcedirs", "", "Source directories (comma-separated)"),
		tag:             flag.String("tag", "", "Optional tag, e.g. build number"),
		title:           flag.String("title", "", "Optional report title (default: 'Coverage Report')"),
		fileFilters:     flag.String("filefilters", "", "File path filters (+Include;-Exclude)"),

		// logging flags
		verbose:   flag.Bool("verbose", false, "Shortcut for Verbose logging (overridden by -verbosity)"),
		verbosity: flag.String("verbosity", "Error", "Logging level: Verbose, Info, Warning, Error, Off"),
		logFile:   flag.String("logfile", "", "Write logs to this file as well as the console"),
		logFormat: flag.String("logformat", "text", "Log output format: text (default) or json"),
	}

	flag.Parse()
	return f, nil
}

func buildLogger(f *cliFlags) (logging.VerbosityLevel, io.Closer, error) {
	verbosityStr := strings.TrimSpace(*f.verbosity)
	level, err := logging.ParseVerbosity(verbosityStr)
	if err != nil && verbosityStr != "" {
		return 0, nil, err
	}

	switch {
	case verbosityStr != "" && verbosityStr != "Error":
	case *f.verbose:
		level = logging.Verbose
	default:
		level = logging.Error
	}

	cfg := logging.Config{
		Verbosity: level,
		File:      *f.logFile,
		Format:    *f.logFormat,
	}
	closer, err := logging.Init(&cfg)
	return level, closer, err
}

func resolveAndValidateInputs(logger *slog.Logger, flags *cliFlags) ([]string, []string, error) {
	if *flags.reportsPatterns == "" {
		return nil, nil, ErrMissingReportFlag
	}

	reportFilePatterns := strings.Split(*flags.reportsPatterns, ";")
	var actualReportFiles []string
	var invalidPatterns []string
	seenFiles := make(map[string]struct{})

	for _, pattern := range reportFilePatterns {
		trimmedPattern := strings.TrimSpace(pattern)
		if trimmedPattern == "" {
			continue
		}
		expandedFiles, err := fsglob.GetFiles(trimmedPattern)
		if err != nil {
			logger.Warn("Error expanding report file pattern", "pattern", trimmedPattern, "error", err)
			invalidPatterns = append(invalidPatterns, trimmedPattern)
			continue
		}
		if len(expandedFiles) == 0 {
			logger.Warn("No files found for report pattern", "pattern", trimmedPattern)
			invalidPatterns = append(invalidPatterns, trimmedPattern)
		}
		for _, file := range expandedFiles {
			absFile, _ := filepath.Abs(file)
			if _, exists := seenFiles[absFile]; !exists {
				if stat, err := os.Stat(absFile); err == nil && !stat.IsDir() {
					actualReportFiles = append(actualReportFiles, absFile)
					seenFiles[absFile] = struct{}{}
				} else if err != nil {
					logger.Warn("Could not stat file from pattern", "pattern", trimmedPattern, "file", absFile, "error", err)
					invalidPatterns = append(invalidPatterns, file)
				}
			}
		}
	}

	if len(actualReportFiles) == 0 {
		return nil, invalidPatterns, fmt.Errorf("no valid report files found after expanding patterns")
	}

	logger.Info("Found report files", "count", len(actualReportFiles))
	logger.Debug("Report file list", "files", strings.Join(actualReportFiles, ", "))
	return actualReportFiles, invalidPatterns, nil
}

func createReportConfiguration(flags *cliFlags, verbosity logging.VerbosityLevel, actualReportFiles, invalidPatterns []string, langFactory *language.ProcessorFactory, logger *slog.Logger) (*reportconfig.ReportConfiguration, error) {
	reportTypes := strings.Split(*flags.reportTypes, ",")
	sourceDirsList := strings.Split(*flags.sourceDirs, ",")
	fileFilterStrings := strings.Split(*flags.fileFilters, ";")

	opts := []reportconfig.Option{
		reportconfig.WithLogger(logger),
		reportconfig.WithVerbosity(verbosity),
		reportconfig.WithInvalidPatterns(invalidPatterns),
		reportconfig.WithTitle(*flags.title),
		reportconfig.WithTag(*flags.tag),
		reportconfig.WithSourceDirectories(sourceDirsList),
		reportconfig.WithReportTypes(reportTypes),
		// Simplified WithFilters call, only passing file filters now.
		reportconfig.WithFilters(
			[]string{}, // assembly filters (deprecated)
			[]string{}, // class filters (deprecated)
			fileFilterStrings,
			[]string{}, // risk hotspot assembly filters (deprecated)
			[]string{}, // risk hotspot class filters (deprecated)
		),
		reportconfig.WithLanguageProcessorFactory(langFactory),
	}

	return reportconfig.NewReportConfiguration(
		actualReportFiles,
		*flags.outputDir,
		opts...,
	)
}

// parseReportFiles iterates through the report file patterns, finds the correct
// parser for each, and returns a collection of parser results.
func parseReportFiles(logger *slog.Logger, reportConfig *reportconfig.ReportConfiguration, parserFactory *parsers.ParserFactory) ([]*parsers.ParserResult, error) {
	var parserResults []*parsers.ParserResult
	var parserErrors []string
	var allUnresolvedFiles []string

	for _, reportFile := range reportConfig.ReportFiles() {
		logger.Info("Attempting to parse report file", "file", reportFile)
		parserInstance, err := parserFactory.FindParserForFile(reportFile)
		if err != nil {
			msg := fmt.Sprintf("no suitable parser found for file %s: %v", reportFile, err)
			parserErrors = append(parserErrors, msg)
			logger.Warn(msg)
			continue
		}

		logger.Info("Using parser for file", "parser", parserInstance.Name(), "file", reportFile)

		result, err := parserInstance.Parse(reportFile, reportConfig)
		if err != nil {
			msg := fmt.Sprintf("error parsing file %s with %s: %v", reportFile, parserInstance.Name(), err)
			parserErrors = append(parserErrors, msg)
			logger.Error(msg)
			continue
		}

		if len(result.UnresolvedSourceFiles) > 0 {
			allUnresolvedFiles = append(allUnresolvedFiles, result.UnresolvedSourceFiles...)
		}

		parserResults = append(parserResults, result)
		logger.Info("Successfully parsed file", "file", reportFile)
	}

	// Handle unresolved source files as a fatal error after attempting all parsers.
	if len(allUnresolvedFiles) > 0 {
		uniqueUnresolvedFiles := utils.DistinctBy(allUnresolvedFiles, func(s string) string { return s })
		logger.Error("Failed to find source files referenced in coverage report", "count", len(uniqueUnresolvedFiles))
		// ... (error logging as before) ...
		return nil, errors.New("failed to find source files referenced in coverage report")
	}

	// Handle case where no reports could be successfully parsed.
	if len(parserResults) == 0 {
		errMsg := "no coverage reports could be parsed successfully"
		if len(parserErrors) > 0 {
			errMsg = fmt.Sprintf("%s. Errors:\r\n- %s", errMsg, strings.Join(parserErrors, "\r\n- "))
		}
		return nil, errors.New(errMsg)
	}

	return parserResults, nil
}

// generateReports orchestrates the creation of all requested report types.
func generateReports(reportCtx reporter.IBuilderContext, summaryTree *model.SummaryTree, fileReader filereader.Reader) error {
	logger := reportCtx.Logger()
	reportConfig := reportCtx.ReportConfiguration()
	outputDir := reportConfig.TargetDirectory()

	logger.Info("Generating reports", "directory", outputDir)
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	for _, reportType := range reportConfig.ReportTypes() {
		trimmedType := strings.TrimSpace(reportType)
		logger.Info("Generating report", "type", trimmedType)

		switch trimmedType {
		case "TextSummary":
			if err := textsummary.NewTextReportBuilder(outputDir, logger).CreateReport(summaryTree); err != nil {
				return fmt.Errorf("failed to generate text report: %w", err)
			}
		case "Html":
			// *** FIX: Pass the fileReader to the HTML report builder ***
			builder := htmlreport.NewHtmlReportBuilder(outputDir, reportCtx, fileReader)
			if err := builder.CreateReport(summaryTree); err != nil {
				return fmt.Errorf("failed to generate HTML report: %w", err)
			}
		case "Lcov":
			if err := lcov.NewLcovReportBuilder(outputDir).CreateReport(summaryTree); err != nil {
				return fmt.Errorf("failed to generate lcov report: %w", err)
			}
		}
	}
	return nil
}

// run is the main application logic flow.
func run(flags *cliFlags) error {
	logger := slog.Default()
	logger.Info("Starting report generation with new architecture.")

	verbosityStr := strings.TrimSpace(*flags.verbosity)
	verbosity, _ := logging.ParseVerbosity(verbosityStr)
	if *flags.verbose {
		verbosity = logging.Verbose
	}

	langFactory := language.NewProcessorFactory(
		defaultformatter.NewDefaultProcessor(),
		csharp.NewCSharpProcessor(),
		golang.NewGoProcessor(),
		cpp.NewCppProcessor(),
	)
	prodFileReader := filereader.NewDefaultReader()
	parserFactory := parsers.NewParserFactory(
		cobertura.NewCoberturaParser(prodFileReader),
		gocover.NewGoCoverParser(prodFileReader),
		gcov.NewGCovParser(prodFileReader),
	)

	actualReportFiles, invalidPatterns, err := resolveAndValidateInputs(logger, flags)
	if err != nil {
		return err
	}
	reportConfig, err := createReportConfiguration(flags, verbosity, actualReportFiles, invalidPatterns, langFactory, logger)
	if err != nil {
		return err
	}

	logger.Info("Executing PARSE stage...")
	parserResults, err := parseReportFiles(logger, reportConfig, parserFactory)
	if err != nil {
		return err
	}
	logger.Info("PARSE stage completed successfully.", "parsed_report_sets", len(parserResults))

	logger.Info("Executing ANALYZE stage...")
	treeBuilder := analyzer.NewTreeBuilder()
	summaryTree, err := treeBuilder.BuildTree(parserResults)
	if err != nil {
		return fmt.Errorf("failed to analyze and build coverage tree: %w", err)
	}
	logger.Info("ANALYZE stage completed successfully. Coverage tree built.")

	logger.Info("Executing REPORT stage...")
	reportCtx := reporter.NewBuilderContext(reportConfig, settings.NewSettings(), logger)

	// *** FIX: Pass prodFileReader to the generateReports function ***
	return generateReports(reportCtx, summaryTree, prodFileReader)
}

func main() {
	start := time.Now()

	flags, err := parseFlags()
	if err != nil {
		fmt.Fprintln(os.Stderr, "flag error:", err)
		os.Exit(1)
	}

	_, closer, err := buildLogger(flags)
	if err != nil {
		fmt.Fprintln(os.Stderr, "logger init error:", err)
		os.Exit(1)
	}
	if closer != nil {
		defer closer.Close()
	}

	if err := run(flags); err != nil {
		slog.Error("An error occurred during report generation", "error", err)
		if errors.Is(err, ErrMissingReportFlag) {
			fmt.Fprintln(os.Stderr, "")
			flag.Usage()
		}
		os.Exit(1)
	}

	slog.Info("Report generation completed successfully", "duration", time.Since(start).Round(time.Millisecond))
}
