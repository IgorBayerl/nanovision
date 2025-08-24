package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/IgorBayerl/fsglob"

	"github.com/IgorBayerl/AdlerCov/internal/config"
	"github.com/IgorBayerl/AdlerCov/internal/filereader"
	"github.com/IgorBayerl/AdlerCov/internal/hydrator"
	"github.com/IgorBayerl/AdlerCov/internal/language"
	"github.com/IgorBayerl/AdlerCov/internal/language/lang_cpp"
	"github.com/IgorBayerl/AdlerCov/internal/language/lang_csharp"
	"github.com/IgorBayerl/AdlerCov/internal/language/lang_default"
	"github.com/IgorBayerl/AdlerCov/internal/language/lang_go"
	"github.com/IgorBayerl/AdlerCov/internal/logging"
	"github.com/IgorBayerl/AdlerCov/internal/model"
	"github.com/IgorBayerl/AdlerCov/internal/parsers"
	"github.com/IgorBayerl/AdlerCov/internal/parsers/parser_cobertura"
	"github.com/IgorBayerl/AdlerCov/internal/parsers/parser_gcov"
	"github.com/IgorBayerl/AdlerCov/internal/parsers/parser_gocover"
	"github.com/IgorBayerl/AdlerCov/internal/reporter/htmlreact"
	"github.com/IgorBayerl/AdlerCov/internal/reporter/lcov"
	"github.com/IgorBayerl/AdlerCov/internal/reporter/reporter_rawjson"
	"github.com/IgorBayerl/AdlerCov/internal/reporter/textsummary"
	"github.com/IgorBayerl/AdlerCov/internal/tree"
)

var ErrMissingReportFlag = errors.New("missing required -report flag")

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
	watch           *bool
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
		watch:           flag.Bool("watch", false, "Enable watch mode to automatically regenerate reports on file changes"),
	}
}

func buildLogger(appConfig *config.AppConfig) (io.Closer, error) {
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

func resolveInputPairs(appConfig *config.AppConfig) []reportInputPair {
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

func parseReportFiles(logger *slog.Logger, appConfig *config.AppConfig, inputPairs []reportInputPair, parserFactory *parsers.ParserFactory) ([]*parsers.ParserResult, error) {
	var parserResults []*parsers.ParserResult
	var parserErrors []string
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

			parseTaskConfig := &parsers.SimpleParserConfig{
				SrcDirs:     []string{pair.SourceDir},
				FileFilter:  appConfig.FileFilterInstance,
				Log:         logger,
				LangFactory: appConfig.LangFactory,
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
			parserResults = append(parserResults, result)
			totalFilesParsed++
			logger.Info("Successfully parsed file", "file", absFile)
		}
	}

	if totalFilesParsed == 0 {
		return nil, errors.New("no coverage reports could be found or parsed successfully")
	}
	return parserResults, nil
}

func generateReports(appConfig *config.AppConfig, summaryTree *model.SummaryTree, fileReader filereader.Reader) error {
	logger := slog.Default()
	outputDir := appConfig.OutputDir

	logger.Info("Generating reports", "directory", outputDir)
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// reportCtx := reporter.NewBuilderContext(appConfig, logger)

	for _, reportType := range appConfig.ReportTypes {
		trimmedType := strings.TrimSpace(reportType)
		logger.Info("Generating report", "type", trimmedType)

		switch trimmedType {
		case "TextSummary":
			if err := textsummary.NewTextReportBuilder(outputDir, logger).CreateReport(summaryTree); err != nil {
				return fmt.Errorf("failed to generate text report: %w", err)
			}
		// case "Html":
		// 	if err := htmlreport.NewHtmlReportBuilder(outputDir, reportCtx, fileReader).CreateReport(summaryTree); err != nil {
		// 		return fmt.Errorf("failed to generate HTML report: %w", err)
		// 	}
		case "Html":
			if err := htmlreact.NewHtmlReactReportBuilder(outputDir, logger).CreateReport(summaryTree); err != nil {
				return fmt.Errorf("failed to generate HTML report: %w", err)
			}
		case "Lcov":
			if err := lcov.NewLcovReportBuilder(outputDir).CreateReport(summaryTree); err != nil {
				return fmt.Errorf("failed to generate lcov report: %w", err)
			}
		case "RawJson":
			if err := reporter_rawjson.NewRawJsonReportBuilder(outputDir).CreateReport(summaryTree); err != nil {
				return fmt.Errorf("failed to generate json summary report: %w", err)
			}
		}

	}
	return nil
}

func executePipeline(appConfig *config.AppConfig) error {
	logger := slog.Default()
	logger.Info("Executing report generation pipeline...")

	// --- Component Initialization ---
	prodFileReader := filereader.NewDefaultReader()
	parserFactory := parsers.NewParserFactory(
		parser_cobertura.NewCoberturaParser(prodFileReader),
		parser_gocover.NewGoCoverParser(prodFileReader),
		parser_gcov.NewGCovParser(prodFileReader),
	)
	treeBuilder := tree.NewBuilder()
	hydrator := hydrator.NewHydrator(prodFileReader, appConfig.LangFactory, logger)

	// --- Pipeline Execution ---

	// 1. Resolve input pairs from the validated config.
	inputPairs := resolveInputPairs(appConfig)
	if len(inputPairs) == 0 {
		return fmt.Errorf("no valid report pattern and source directory pairs were provided")
	}

	// 2. PARSE Stage
	logger.Info("Executing PARSE stage...")
	parserResults, err := parseReportFiles(logger, appConfig, inputPairs, parserFactory)
	if err != nil {
		return err
	}
	logger.Info("PARSE stage completed successfully.", "parsed_report_sets", len(parserResults))

	// 3. TREE BUILDER Stage
	logger.Info("Executing TREE BUILDER stage...")
	rawTree, err := treeBuilder.BuildTree(parserResults)
	if err != nil {
		return fmt.Errorf("failed to build coverage tree: %w", err)
	}
	logger.Info("TREE BUILDER stage completed successfully.")

	// 4. HYDRATOR Stage
	logger.Info("Executing HYDRATOR stage...")
	if err := hydrator.HydrateTree(rawTree); err != nil {
		return fmt.Errorf("failed to hydrate coverage tree: %w", err)
	}
	logger.Info("HYDRATOR stage completed successfully.")

	// 5. REPORT Stage
	logger.Info("Executing REPORT stage...")
	return generateReports(appConfig, rawTree, prodFileReader)
}

func main() {
	start := time.Now()
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	rawFlags := parseFlags()
	flag.Parse()

	langFactory := language.NewProcessorFactory(
		lang_default.NewDefaultProcessor(),
		lang_csharp.NewCSharpProcessor(),
		lang_go.NewGoProcessor(),
		lang_cpp.NewCppProcessor(),
	)

	verbosity, _ := logging.ParseVerbosity(*rawFlags.verbosity)
	if *rawFlags.verbose {
		verbosity = logging.Verbose
	}

	appConfig, err := config.BuildAppConfig(
		*rawFlags.reportsPatterns, *rawFlags.sourceDirs, *rawFlags.reportTypes, *rawFlags.fileFilters,
		*rawFlags.outputDir, *rawFlags.tag, *rawFlags.title, *rawFlags.logFile, *rawFlags.logFormat,
		verbosity, langFactory,
	)
	if err != nil {
		slog.Error("Configuration error", "error", err)
		if errors.Is(err, ErrMissingReportFlag) {
			fmt.Fprintln(os.Stderr, "")
			flag.Usage()
		}
		os.Exit(1)
	}

	closer, err := buildLogger(appConfig)
	if err != nil {
		fmt.Fprintln(os.Stderr, "logger init error:", err)
		os.Exit(1)
	}
	if closer != nil {
		defer closer.Close()
	}

	if err := executePipeline(appConfig); err != nil {
		slog.Error("An error occurred during report generation", "error", err)
		os.Exit(1)
	}

	if *rawFlags.watch {
		slog.Info("Watch mode enabled. Press Ctrl+C to exit.")
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		<-quit
		slog.Info("Shutdown signal received, exiting.")
	}

	slog.Info("Report generation completed successfully", "duration", time.Since(start).Round(time.Millisecond))
}
