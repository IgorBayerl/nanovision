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
	"runtime"
	"strings"
	"time"

	"github.com/IgorBayerl/fsglob"
	"github.com/IgorBayerl/nanovision/internal/aggregator"
	"github.com/IgorBayerl/nanovision/internal/analyzer"
	cpp "github.com/IgorBayerl/nanovision/internal/analyzer/cpp"
	"github.com/IgorBayerl/nanovision/internal/analyzer/gdscript"
	golang "github.com/IgorBayerl/nanovision/internal/analyzer/go"
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/diff"
	"github.com/IgorBayerl/nanovision/internal/diffapply"
	"github.com/IgorBayerl/nanovision/internal/enricher"
	"github.com/IgorBayerl/nanovision/internal/filereader"
	"github.com/IgorBayerl/nanovision/internal/logging"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/parsers"
	"github.com/IgorBayerl/nanovision/internal/parsers/parser_cobertura"
	"github.com/IgorBayerl/nanovision/internal/parsers/parser_gcov"
	"github.com/IgorBayerl/nanovision/internal/parsers/parser_gocover"
	"github.com/IgorBayerl/nanovision/internal/parsers/parser_lcov"
	"github.com/IgorBayerl/nanovision/internal/reporter/htmlreact"
	"github.com/IgorBayerl/nanovision/internal/reporter/lcov"
	"github.com/IgorBayerl/nanovision/internal/reporter/reporter_rawjson"
	"github.com/IgorBayerl/nanovision/internal/reporter/textsummary"
	"github.com/IgorBayerl/nanovision/internal/status"
	"github.com/IgorBayerl/nanovision/internal/tree"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

type repeatedStringFlag []string

func (r *repeatedStringFlag) String() string {
	return strings.Join(*r, ", ")
}

func (r *repeatedStringFlag) Set(value string) error {
	*r = append(*r, value)
	return nil
}

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
	flag.StringVar(&rawInput.DiffFile, "diff", "", "Path to a unified diff file for patch coverage analysis")
	flag.StringVar(&rawInput.DiffStrip, "diff-strip", "", "Strip N leading components from diff paths ('auto' or 0-6)")
	flag.BoolVar(&rawInput.IgnoreCache, "ignore-cache", false, "Ignore existing cache and force re-analysis")
	flag.Var((*repeatedStringFlag)(&rawInput.StatusBands), "threshold", "Metric threshold (e.g. 'line_coverage=60..80'). Can be repeated.")
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

			result, err := parserInstance.Parse(absFile, parseTaskConfig)
			if err != nil {
				msg := fmt.Sprintf("error parsing file %s with %s: %v", reportFile, parserInstance.Name(), err)
				parserErrors = append(parserErrors, msg)
				logger.Error(msg)
				continue
			}

			result.SourceDirectory = pair.SourceDir
			result.ReportPattern = pair.ReportPattern
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
			err = htmlreact.NewHtmlReactReportBuilder(outputDir, logger, false).CreateReport(summaryTree)
		case "HtmlUnified":
			err = htmlreact.NewHtmlReactReportBuilder(outputDir, logger, true).CreateReport(summaryTree)
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

func deriveCapabilities(tree *model.SummaryTree) status.Capabilities {
	caps := status.Capabilities{}
	var walk func(n *model.DirNode)
	walk = func(n *model.DirNode) {
		if n.Metrics.BranchesValid > 0 {
			caps.HasBranchCoverage = true
		}
		if n.Metrics.MethodsValid > 0 {
			caps.HasMethodCoverage = true
		}
		for _, c := range n.Subdirs {
			walk(c)
		}
		for _, f := range n.Files {
			if f.Metrics.BranchesValid > 0 {
				caps.HasBranchCoverage = true
			}
			if f.Metrics.MethodsValid > 0 {
				caps.HasMethodCoverage = true
			}
		}
	}
	walk(tree.Root)
	return caps
}

// executePipeline orchestrates the report generation process from start to finish.
func executePipeline(appConfig *config.AppConfig, diffData *diff.DiffData) error {
	logger := slog.Default()
	logger.Info("Executing report generation pipeline...")

	prodFileReader := filereader.NewDefaultReader()
	parserFactory := createParserFactory(logger, prodFileReader)
	treeBuilder := tree.NewBuilder(appConfig.ProjectRoot, appConfig.FileFilterInstance)

	allAnalyzers := []analyzer.Analyzer{
		golang.New(),
		cpp.New(),
		gdscript.New(),
	}

	// Determine a safe cache location.
	// We try the user's standard cache dir (e.g. ~/.cache/nanovision)
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		logger.Warn("Could not determine user cache directory; caching will be disabled.")
		userCacheDir = ""
	}
	cachePath := ""
	if userCacheDir != "" {
		cachePath = filepath.Join(userCacheDir, "nanovision")
	}

	treeEnricher := enricher.New(allAnalyzers, prodFileReader, logger, cachePath, appConfig.IgnoreCache)

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

	if diffData != nil {
		logger.Info("Executing DIFF ANALYSIS stage...")
		diffapply.Apply(summaryTree, diffData, logger)
		logger.Info("DIFF ANALYSIS stage completed successfully.")
	}

	aggregator.AggregateMetricsAfterEnrichment(summaryTree)

	logger.Info("Executing ANNOTATE stage...")
	caps := deriveCapabilities(summaryTree)
	status.Annotate(summaryTree, appConfig.StatusBands, caps)
	logger.Info("ANNOTATE stage completed successfully.")

	logger.Info("Executing REPORT stage...")
	return generateReports(appConfig, summaryTree)
}

func determineProjectRoot(configPath string) (string, error) {
	if configPath != "" {
		absConfigPath, err := filepath.Abs(configPath)
		if err != nil {
			return "", fmt.Errorf("could not determine absolute path for config file: %w", err)
		}
		return filepath.Dir(absConfigPath), nil
	}

	// Fallback to current working directory if no config file is used
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("could not get current working directory: %w", err)
	}
	return wd, nil
}

func main() {
	start := time.Now()
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	configPath := flag.String("config", "", "Path to a nanovision.yaml configuration file.")
	watchFlag := flag.Bool("watch", false, "Enable watch mode to automatically regenerate reports on file changes")
	versionFlag := flag.Bool("version", false, "Print version information and exit")
	listParsersFlag := flag.Bool("list-parsers", false, "List supported coverage report parsers and exit")

	rawInput := parseAndBindFlags()
	flag.Parse()

	if *listParsersFlag {
		logger := slog.Default()
		prodFileReader := filereader.NewDefaultReader()
		factory := createParserFactory(logger, prodFileReader)

		fmt.Println("Supported Coverage Parsers:")
		for _, name := range factory.RegisteredParsers() {
			fmt.Printf(" - %s\n", name)
		}
		os.Exit(0)
	}

	// Handle version output
	if *versionFlag {
		fmt.Printf("nanovision version %s\n", version)
		fmt.Printf("commit: %s\n", commit)
		fmt.Printf("built at: %s\n", date)
		fmt.Printf("os/arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}

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

	appConfig.ProjectRoot, err = determineProjectRoot(*configPath)
	if err != nil {
		slog.Error("Failed to determine project root", "error", err)
		os.Exit(1)
	}
	slog.Info("Project root determined", "path", appConfig.ProjectRoot)

	closer, err := buildLogger(appConfig)
	if err != nil {
		fmt.Fprintln(os.Stderr, "logger init error:", err)
		os.Exit(1)
	}
	if closer != nil {
		defer closer.Close()
	}

	var diffData *diff.DiffData
	if appConfig.Diff.File != "" {
		var parseErr error
		diffLogger := slog.Default()
		diffLogger.Info("Parsing diff file...", "path", appConfig.Diff.File)
		diffData, parseErr = diff.Parse(appConfig.Diff.File, diffLogger)
		if parseErr != nil {
			diffLogger.Warn("Failed to parse diff file; ignoring diff analysis.", "error", parseErr)
			diffData = nil
		}
	}

	if err := executePipeline(appConfig, diffData); err != nil {
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

func createParserFactory(logger *slog.Logger, reader filereader.Reader) *parsers.ParserFactory {
	return parsers.NewParserFactory(
		logger,
		parser_cobertura.NewCoberturaParser(reader),
		parser_gocover.NewGoCoverParser(reader),
		parser_gcov.NewGCovParser(reader),
		parser_lcov.NewLcovParser(reader),
	)
}
