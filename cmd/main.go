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
// It holds the raw string values directly from the command line.
type cliFlags struct {
	reportsPatterns *string
	outputDir       *string
	reportTypes     *string
	sourceDirs      *string
	tag             *string
	title           *string
	fileFilters     *string
	verbose         *bool
	verbosity       *string
	logFile         *string
	logFormat       *string
}

// AppConfig holds the parsed and validated configuration for the application.
type AppConfig struct {
	ReportPatterns []string
	SourceDirs     []string
	ReportTypes    []string
	FileFilters    []string
	OutputDir      string
	Tag            string
	Title          string
	LogFile        string
	LogFormat      string
	Verbosity      logging.VerbosityLevel
}

func parseFlags() *cliFlags {
	return &cliFlags{
		reportsPatterns: flag.String("report", "", "Coverage report file paths or patterns (semicolon-separated)"),
		outputDir:       flag.String("output", "coverage-report", "Output directory for generated reports"),
		reportTypes:     flag.String("reporttypes", "TextSummary,Html", "Report types (comma-separated)"),
		sourceDirs:      flag.String("sourcedirs", "", "Source directories (semicolon-separated, one per report pattern)"),
		tag:             flag.String("tag", "", "Optional tag, e.g. build number"),
		title:           flag.String("title", "", "Optional report title (default: 'Coverage Report')"),
		fileFilters:     flag.String("filefilters", "", "File path filters (+Include;-Exclude, semicolon-separated)"),
		verbose:         flag.Bool("verbose", false, "Shortcut for Verbose logging (overridden by -verbosity)"),
		verbosity:       flag.String("verbosity", "Info", "Logging level: Verbose, Info, Warning, Error, Off"),
		logFile:         flag.String("logfile", "", "Write logs to this file as well as the console"),
		logFormat:       flag.String("logformat", "text", "Log output format: text (default) or json"),
	}
}

// buildAppConfig is the single point of truth for parsing and validating CLI flags.
func buildAppConfig(flags *cliFlags) (*AppConfig, error) {
	if *flags.reportsPatterns == "" {
		return nil, ErrMissingReportFlag
	}

	verbosityStr := strings.TrimSpace(*flags.verbosity)
	level, err := logging.ParseVerbosity(verbosityStr)
	if err != nil && verbosityStr != "" {
		return nil, fmt.Errorf("invalid verbosity level %q", verbosityStr)
	}

	if *flags.verbose {
		level = logging.Verbose
	}

	reportPatterns := strings.Split(*flags.reportsPatterns, ";")
	sourceDirs := strings.Split(*flags.sourceDirs, ";")

	if len(reportPatterns) != len(sourceDirs) {
		return nil, fmt.Errorf(
			"mismatch between number of report patterns (%d) and source directories (%d). You must provide one source directory for each report pattern",
			len(reportPatterns),
			len(sourceDirs),
		)
	}

	return &AppConfig{
		ReportPatterns: reportPatterns,
		SourceDirs:     sourceDirs,
		ReportTypes:    strings.Split(*flags.reportTypes, ","),
		FileFilters:    strings.Split(*flags.fileFilters, ";"),
		OutputDir:      *flags.outputDir,
		Tag:            *flags.tag,
		Title:          *flags.title,
		LogFile:        *flags.logFile,
		LogFormat:      *flags.logFormat,
		Verbosity:      level,
	}, nil
}

func buildLogger(appConfig *AppConfig) (io.Closer, error) {
	cfg := logging.Config{
		Verbosity: appConfig.Verbosity,
		File:      appConfig.LogFile,
		Format:    appConfig.LogFormat,
	}
	return logging.Init(&cfg)
}

type reportInputPair struct {
	ReportPattern string
	SourceDir     string
}

// resolveInputPairs now takes the already-parsed slices from AppConfig.
func resolveInputPairs(appConfig *AppConfig) []reportInputPair {
	var pairs []reportInputPair
	for i := 0; i < len(appConfig.ReportPatterns); i++ {
		trimmedPattern := strings.TrimSpace(appConfig.ReportPatterns[i])
		trimmedSourceDir := strings.TrimSpace(appConfig.SourceDirs[i])
		if trimmedPattern != "" && trimmedSourceDir != "" {
			pairs = append(pairs, reportInputPair{
				ReportPattern: trimmedPattern,
				SourceDir:     trimmedSourceDir,
			})
		}
	}
	return pairs
}

// createReportConfiguration is simplified to accept the parsed AppConfig.
func createReportConfiguration(appConfig *AppConfig, actualReportFiles []string, langFactory *language.ProcessorFactory) (*reportconfig.ReportConfiguration, error) {
	opts := []reportconfig.Option{
		reportconfig.WithLogger(slog.Default()),
		reportconfig.WithVerbosity(appConfig.Verbosity),
		reportconfig.WithTitle(appConfig.Title),
		reportconfig.WithTag(appConfig.Tag),
		reportconfig.WithSourceDirectories(appConfig.SourceDirs), // Global list for context
		reportconfig.WithReportTypes(appConfig.ReportTypes),
		reportconfig.WithFilters(
			[]string{}, // assembly filters (deprecated)
			[]string{}, // class filters (deprecated)
			appConfig.FileFilters,
			[]string{}, // risk hotspot assembly filters (deprecated)
			[]string{}, // risk hotspot class filters (deprecated)
		),
		reportconfig.WithLanguageProcessorFactory(langFactory),
	}

	return reportconfig.NewReportConfiguration(
		actualReportFiles,
		appConfig.OutputDir,
		opts...,
	)
}

func parseReportFiles(logger *slog.Logger, appConfig *AppConfig, inputPairs []reportInputPair, parserFactory *parsers.ParserFactory, langFactory *language.ProcessorFactory) ([]*parsers.ParserResult, error) {
	var parserResults []*parsers.ParserResult
	var parserErrors []string
	var allUnresolvedFiles []string
	var totalFilesParsed int

	for _, pair := range inputPairs {
		expandedFiles, err := fsglob.GetFiles(pair.ReportPattern)
		if err != nil {
			logger.Warn("Error expanding report file pattern", "pattern", pair.ReportPattern, "error", err)
			continue
		}
		if len(expandedFiles) == 0 {
			logger.Warn("No files found for report pattern", "pattern", pair.ReportPattern)
		}

		for _, reportFile := range expandedFiles {
			absFile, _ := filepath.Abs(reportFile)
			logger.Info("Attempting to parse report file", "file", absFile, "sourcedir", pair.SourceDir)

			parseTaskConfig, err := reportconfig.NewReportConfiguration(
				[]string{absFile},
				appConfig.OutputDir,
				reportconfig.WithLogger(logger),
				reportconfig.WithVerbosity(appConfig.Verbosity),
				reportconfig.WithSourceDirectories([]string{pair.SourceDir}), // Use the specific paired source dir
				reportconfig.WithFilters([]string{}, []string{}, appConfig.FileFilters, []string{}, []string{}),
				reportconfig.WithLanguageProcessorFactory(langFactory),
			)
			if err != nil {
				msg := fmt.Sprintf("failed to create specific config for %s: %v", absFile, err)
				parserErrors = append(parserErrors, msg)
				logger.Error(msg)
				continue
			}

			parserInstance, err := parserFactory.FindParserForFile(absFile)
			if err != nil {
				msg := fmt.Sprintf("no suitable parser found for file %s: %v", absFile, err)
				parserErrors = append(parserErrors, msg)
				logger.Warn(msg)
				continue
			}

			logger.Info("Using parser for file", "parser", parserInstance.Name(), "file", absFile)

			result, err := parserInstance.Parse(absFile, parseTaskConfig)
			if err != nil {
				msg := fmt.Sprintf("error parsing file %s with %s: %v", reportFile, parserInstance.Name(), err)
				parserErrors = append(parserErrors, msg)
				logger.Error(msg)
				continue
			}

			result.SourceDirectory = pair.SourceDir

			if len(result.UnresolvedSourceFiles) > 0 {
				allUnresolvedFiles = append(allUnresolvedFiles, result.UnresolvedSourceFiles...)
			}

			parserResults = append(parserResults, result)
			totalFilesParsed++
			logger.Info("Successfully parsed file", "file", absFile)
		}
	}

	if len(allUnresolvedFiles) > 0 {
		uniqueUnresolvedFiles := utils.DistinctBy(allUnresolvedFiles, func(s string) string { return s })
		logger.Error("Failed to find source files referenced in coverage report", "count", len(uniqueUnresolvedFiles))
		return nil, errors.New("failed to find source files referenced in coverage report")
	}

	if totalFilesParsed == 0 {
		errMsg := "no coverage reports could be found or parsed successfully from the provided patterns"
		if len(parserErrors) > 0 {
			errMsg = fmt.Sprintf("%s. Errors:\r\n- %s", errMsg, strings.Join(parserErrors, "\r\n- "))
		}
		return nil, errors.New(errMsg)
	}

	return parserResults, nil
}

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

// run is the main application logic, now driven by the AppConfig.
func run(appConfig *AppConfig) error {
	logger := slog.Default()
	logger.Info("Starting report generation with new architecture.")

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

	// Step 1: Resolve the (report, sourcedir) pairs from the pre-validated config.
	inputPairs := resolveInputPairs(appConfig)
	if len(inputPairs) == 0 {
		return fmt.Errorf("no valid report pattern and source directory pairs were provided")
	}

	// Step 2: Execute the parsing stage.
	logger.Info("Executing PARSE stage...")
	parserResults, err := parseReportFiles(logger, appConfig, inputPairs, parserFactory, langFactory)
	if err != nil {
		return err
	}
	logger.Info("PARSE stage completed successfully.", "parsed_report_sets", len(parserResults))

	// Step 3: Analyze the results.
	logger.Info("Executing ANALYZE stage...")
	treeBuilder := analyzer.NewTreeBuilder()
	summaryTree, err := treeBuilder.BuildTree(parserResults)
	if err != nil {
		return fmt.Errorf("failed to analyze and build coverage tree: %w", err)
	}
	logger.Info("ANALYZE stage completed successfully. Coverage tree built.")

	// Step 4: Create the final report configuration for the reporting stage.
	globalReportConfig, err := createReportConfiguration(appConfig, []string{}, langFactory)
	if err != nil {
		return err
	}

	// Step 5: Generate reports.
	logger.Info("Executing REPORT stage...")
	reportCtx := reporter.NewBuilderContext(globalReportConfig, settings.NewSettings(), logger)

	return generateReports(reportCtx, summaryTree, prodFileReader)
}

func main() {
	start := time.Now()
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	// 1. Get raw flags
	flags := parseFlags()
	flag.Parse()

	// 2. Build unified and validated AppConfig
	appConfig, err := buildAppConfig(flags)
	if err != nil {
		slog.Error("Configuration error", "error", err)
		if errors.Is(err, ErrMissingReportFlag) {
			fmt.Fprintln(os.Stderr, "")
			flag.Usage()
		}
		os.Exit(1)
	}

	// 3. Initialize logger from AppConfig
	closer, err := buildLogger(appConfig)
	if err != nil {
		fmt.Fprintln(os.Stderr, "logger init error:", err)
		os.Exit(1)
	}
	if closer != nil {
		defer closer.Close()
	}

	// 4. Run the application logic with the clean AppConfig
	if err := run(appConfig); err != nil {
		slog.Error("An error occurred during report generation", "error", err)
		os.Exit(1)
	}

	slog.Info("Report generation completed successfully", "duration", time.Since(start).Round(time.Millisecond))
}
