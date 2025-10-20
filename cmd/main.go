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

	"github.com/IgorBayerl/AdlerCov/analyzer"
	cpp "github.com/IgorBayerl/AdlerCov/analyzer/cpp"
	golang "github.com/IgorBayerl/AdlerCov/analyzer/go"
	"github.com/IgorBayerl/AdlerCov/filereader"
	"github.com/IgorBayerl/AdlerCov/internal/aggregator"
	"github.com/IgorBayerl/AdlerCov/internal/config"
	"github.com/IgorBayerl/AdlerCov/internal/enricher"
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
	"github.com/IgorBayerl/AdlerCov/logging"
)

func parseAndBindFlags() *config.RawConfigInput {
	rawInput := &config.RawConfigInput{}

	flag.StringVar(&rawInput.ReportPatterns, "report", "", "Coverage report file paths or patterns (semicolon-separated)")
	flag.StringVar(&rawInput.OutputDir, "output", "coverage-report", "Output directory for generated reports")
	flag.StringVar(&rawInput.ReportTypes, "reporttypes", "TextSummary,Html", "Report types (comma-separated)")
	flag.StringVar(&rawInput.SourceDirs, "sourcedirs", "", "Source directories (semicolon-separated, one per report pattern)")
	flag.StringVar(&rawInput.Tag, "tag", "", "Optional tag, e.g. build number")
	flag.StringVar(&rawInput.Title, "title", "", "Optional report title (default: 'Coverage Report')")
	flag.StringVar(&rawInput.FileFilters, "filefilters", "", "File path filters (+Include;-Exclude, semicolon-separated)")
	flag.StringVar(&rawInput.LogFile, "logfile", "", "Write logs to this file as well as the console")
	flag.StringVar(&rawInput.LogFormat, "logformat", "text", "Log output format: text (default) or json")
	flag.StringVar(&rawInput.Verbosity, "verbosity", "Info", "Logging level: Verbose, Info, Warning, Error, Off")
	flag.BoolVar(&rawInput.Verbose, "verbose", false, "Shortcut for Verbose logging (overridden by -verbosity)")
	return rawInput
}

func buildLogger(appConfig *config.AppConfig) (io.Closer, error) {
	cfg := logging.Config{
		Verbosity: appConfig.VerbosityLevel,
		File:      appConfig.LogFile,
		Format:    appConfig.LogFormat,
	}
	return logging.Init(&cfg)
}

func parseReportFiles(logger *slog.Logger, appConfig *config.AppConfig, inputPairs []config.ReportInputPair, parserFactory *parsers.ParserFactory) ([]*parsers.ParserResult, error) {
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
				SrcDirs:    []string{pair.SourceDir},
				FileFilter: appConfig.FileFilterInstance,
				Log:        logger,
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
	if len(parserErrors) > 0 {
		return parserResults, fmt.Errorf("encountered errors during parsing: %s", strings.Join(parserErrors, "; "))
	}
	return parserResults, nil
}

func generateReports(appConfig *config.AppConfig, summaryTree *model.SummaryTree) error {
	logger := slog.Default()
	outputDir := appConfig.OutputDir

	logger.Info("Generating reports", "directory", outputDir)
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	for _, reportType := range appConfig.ReportTypes {
		trimmedType := strings.TrimSpace(reportType)
		logger.Info("Generating report", "type", trimmedType)
		var err error
		switch trimmedType {
		case "TextSummary":
			err = textsummary.NewTextReportBuilder(outputDir, logger).CreateReport(summaryTree)
		case "Html":
			err = htmlreact.NewHtmlReactReportBuilder(outputDir, logger).CreateReport(summaryTree)
		case "Lcov":
			err = lcov.NewLcovReportBuilder(outputDir).CreateReport(summaryTree)
		case "RawJson":
			err = reporter_rawjson.NewRawJsonReportBuilder(outputDir).CreateReport(summaryTree)
		}
		if err != nil {
			return fmt.Errorf("failed to generate '%s' report: %w", trimmedType, err)
		}
	}
	return nil
}

// executePipeline orchestrates the report generation process from start to finish.
//
// The function attempts to parse each report file provided by the user. If an
// individual file cannot be processed, the error is logged, and the pipeline
// continues to the next file.
//
// To ensure accuracy, the entire process will halt if any parsing errors
// occurred. This prevents the creation of a final report from incomplete data.
//
// The process follows these stages:
//   - Component Initialization: Prepares all the tools needed for the pipeline,
//     such as the parsers and the tree builder.
//   - Parse: Reads the different coverage report formats into a standard structure.
//   - Build: Combines data from all parsed reports into a single project tree.
//   - Enrich: Gathers extra details from the source code, like method complexity.
//   - Report: Generates the final output files, such as the HTML and text summaries.
func executePipeline(appConfig *config.AppConfig) error {
	logger := slog.Default()
	logger.Info("Executing report generation pipeline...")

	prodFileReader := filereader.NewDefaultReader()
	parserFactory := parsers.NewParserFactory(
		parser_cobertura.NewCoberturaParser(prodFileReader),
		parser_gocover.NewGoCoverParser(prodFileReader),
		parser_gcov.NewGCovParser(prodFileReader),
	)
	treeBuilder := tree.NewBuilder()

	allAnalyzers := []analyzer.Analyzer{
		golang.New(),
		cpp.New(),
	}
	treeEnricher := enricher.New(allAnalyzers, prodFileReader, logger)

	if len(appConfig.InputPairs) == 0 {
		return fmt.Errorf("no valid report pattern and source directory pairs were provided")
	}

	logger.Info("Executing PARSE stage...")
	parserResults, err := parseReportFiles(logger, appConfig, appConfig.InputPairs, parserFactory)
	if err != nil {
		return err
	}
	logger.Info("PARSE stage completed successfully.", "parsed_report_sets", len(parserResults))

	logger.Info("Executing BUILD stage...")
	summaryTree, err := treeBuilder.BuildTree(parserResults)
	if err != nil {
		return fmt.Errorf("failed to build and aggregate coverage tree: %w", err)
	}
	logger.Info("BUILD stage completed successfully.")

	logger.Info("Executing ENRICH stage...")
	treeEnricher.EnrichTree(summaryTree)
	logger.Info("ENRICH stage completed successfully.")

	aggregator.AggregateMetricsAfterEnrichment(summaryTree)

	logger.Info("Executing REPORT stage...")
	return generateReports(appConfig, summaryTree)
}

func main() {
	start := time.Now()
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	configPath := flag.String("config", "", "Path to a adlercov.yaml configuration file.")
	watchFlag := flag.Bool("watch", false, "Enable watch mode to automatically regenerate reports on file changes")

	rawInput := parseAndBindFlags()
	flag.Parse()

	if _, err := logging.ParseVerbosity(rawInput.Verbosity); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: %v. Defaulting to 'Info' level.\n", err)
	}

	appConfig, err := config.Load(*configPath, *rawInput)
	if err != nil {
		slog.Error("Configuration error", "error", err)
		if strings.Contains(err.Error(), "must be specified") {
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

	if *watchFlag {
		slog.Info("Watch mode enabled. Press Ctrl+C to exit.")
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		<-quit
		slog.Info("Shutdown signal received, exiting.")
	}

	slog.Info("Report generation completed successfully", "duration", time.Since(start).Round(time.Millisecond))
}
