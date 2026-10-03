package main

import (
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"time"

	"github.com/IgorBayerl/nanovision/internal/bootlog"
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/diagnostics"
	"github.com/IgorBayerl/nanovision/internal/diff"
	"github.com/IgorBayerl/nanovision/internal/filereader"
	"github.com/IgorBayerl/nanovision/internal/logging"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/pipeline"
	"github.com/IgorBayerl/nanovision/internal/reporter/annotations"
	"github.com/IgorBayerl/nanovision/internal/reporter/htmlreact"
	"github.com/IgorBayerl/nanovision/internal/reporter/lcov"
	"github.com/IgorBayerl/nanovision/internal/reporter/reporter_rawjson"
	"github.com/IgorBayerl/nanovision/internal/reporter/sarif"
	"github.com/IgorBayerl/nanovision/internal/reporter/textsummary"
	"github.com/IgorBayerl/nanovision/internal/review"
	"github.com/IgorBayerl/nanovision/internal/status/evaluators"
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

	flag.Var((*repeatedStringFlag)(&rawInput.Reports), "report", "Coverage report: 'path[,source=DIR][,name=LABEL]'. Repeat it for each report.")
	flag.StringVar(&rawInput.SourceDirs, "sourcedirs", "", "Source folder for the reports that name none (several: semicolon-separated, one for each report)")
	flag.Var((*repeatedStringFlag)(&rawInput.FileMetrics), "file-metric", "File metric to show: 'name' or 'name=min..max' with a warning range. Repeat it, in column order. See -list-metrics.")
	flag.Var((*repeatedStringFlag)(&rawInput.MethodMetrics), "method-metric", "Method metric to show: 'name' or 'name=min..max'. Repeat it, in column order. See -list-metrics.")
	flag.Var((*repeatedStringFlag)(&rawInput.StatusBands), "threshold", "Warning range of a metric: 'name=min..max', or 'methods.name=min..max' for a method metric. Repeat it.")
	flag.Var((*repeatedStringFlag)(&rawInput.Set), "set", "Set any config key: 'key=value', e.g. 'review.hotspots=5'. Repeat it. 'nanovision config docs' lists the keys.")
	flag.StringVar(&rawInput.VCS, "vcs", "", "Version control system: 'git' or 'perforce' (vcs.type)")
	flag.StringVar(&rawInput.OutputDir, "output", "coverage-report", "Output directory for generated reports")
	flag.StringVar(&rawInput.ReportTypes, "reporttypes", "TextSummary,Html", "Report types (comma-separated)")
	flag.StringVar(&rawInput.Tag, "tag", "", "Optional tag, e.g. build number")
	flag.StringVar(&rawInput.Title, "title", "", "Optional report title (default: 'Coverage Report')")
	flag.StringVar(&rawInput.FileFilters, "filefilters", "", "File path filters (+Include;-Exclude, semicolon-separated)")
	flag.StringVar(&rawInput.LogFile, "logfile", "", "Write logs to this file as well as the console")
	flag.StringVar(&rawInput.LogFormat, "logformat", "text", "Log output format: text (default) or json")
	flag.StringVar(&rawInput.Verbosity, "verbosity", "Info", "Logging level: Verbose, Info, Warning, Error, Off")
	flag.BoolVar(&rawInput.Verbose, "verbose", false, "Shortcut for Verbose logging (overridden by -verbosity)")
	flag.StringVar(&rawInput.DiffFile, "diff", "", "Path to a unified diff file for patch coverage analysis")
	flag.StringVar(&rawInput.DiffStrip, "diff-strip", "", "Strip N leading components from diff paths ('auto' or 0-6)")
	flag.BoolVar(&rawInput.OnlyChanged, "only-changed", false, "The HTML report holds only the changed files; needs a diff")
	flag.BoolVar(&rawInput.IgnoreCache, "ignore-cache", false, "Ignore existing cache and force re-analysis")
	flag.StringVar(&rawInput.DefaultFilters, "default-filters", "", "Raw URL query string of filters auto-applied when the report opens (e.g. 'diff=changed&risk=danger')")
	flag.StringVar(&rawInput.FailOn, "fail-on", "", "Exit non-zero when the review gate fails or changed code has problems: 'error', 'warning' or 'never' (default)")

	// run history and the coverage delta
	flag.StringVar(&rawInput.Store, "store", "", "Run store: a folder for a local store, or the team server URL (history.store)")
	flag.StringVar(&rawInput.Project, "project", "", "Project name in the run store (history.project, default: the project folder name)")
	flag.StringVar(&rawInput.Profile, "profile", "", "Keeps runs of different test suites or platforms apart (history.profile, default 'default')")
	flag.StringVar(&rawInput.RunKind, "run-kind", "", "Kind of this run: 'submit', 'review' or 'local' (default)")
	flag.StringVar(&rawInput.Revision, "revision", "", "Revision of this run: a commit hash, or '//stream@changelist'. Use it only when the workspace matches it exactly")
	flag.StringVar(&rawInput.BaseRevision, "base-revision", "", "Revision to compare with: a commit hash, or '//stream@changelist'")
	flag.StringVar(&rawInput.Stream, "stream", "", "Stream or branch of this run, when no VCS adapter knows it")
	flag.StringVar(&rawInput.ReviewID, "review-id", "", "Shelved changelist or pull request of a review run")
	flag.StringVar(&rawInput.Author, "author", "", "Author of the change")
	flag.StringVar(&rawInput.CIURL, "ci-url", "", "Link to the CI build (default: from BUILD_URL, CI_JOB_URL or GitHub Actions)")
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
		write, ok := reportWriters[trimmedType]
		if !ok {
			return fmt.Errorf("no writer for report type '%s'", trimmedType)
		}
		if err := write(appConfig, summaryTree, logger); err != nil {
			return fmt.Errorf("failed to generate '%s' report: %w", trimmedType, err)
		}
	}
	return nil
}

// reportWriters writes each output format of config.OutputFormats. A test
// fails when a format has no writer here, or a writer has no entry there.
var reportWriters = map[string]func(cfg *config.AppConfig, tree *model.SummaryTree, logger *slog.Logger) error{
	"TextSummary": func(cfg *config.AppConfig, tree *model.SummaryTree, logger *slog.Logger) error {
		return textsummary.NewTextReportBuilder(cfg.OutputDir, logger, cfg).CreateReport(tree)
	},
	"Html": func(cfg *config.AppConfig, tree *model.SummaryTree, logger *slog.Logger) error {
		return writeHTML(cfg, tree, logger, false)
	},
	"HtmlUnified": func(cfg *config.AppConfig, tree *model.SummaryTree, logger *slog.Logger) error {
		return writeHTML(cfg, tree, logger, true)
	},
	"Lcov": func(cfg *config.AppConfig, tree *model.SummaryTree, _ *slog.Logger) error {
		return lcov.NewLcovReportBuilder(cfg.OutputDir).CreateReport(tree)
	},
	"RawJson": func(cfg *config.AppConfig, tree *model.SummaryTree, _ *slog.Logger) error {
		return reporter_rawjson.NewRawJsonReportBuilder(cfg.OutputDir).CreateReport(tree)
	},
	"Sarif": func(cfg *config.AppConfig, tree *model.SummaryTree, _ *slog.Logger) error {
		return sarif.NewSarifReportBuilder(cfg.OutputDir, cfg, evaluators.Registry).CreateReport(tree)
	},
	"Annotations": func(cfg *config.AppConfig, tree *model.SummaryTree, _ *slog.Logger) error {
		return annotations.NewAnnotationsReportBuilder(cfg.OutputDir, cfg, evaluators.Registry).CreateReport(tree)
	},
}

func writeHTML(cfg *config.AppConfig, tree *model.SummaryTree, logger *slog.Logger, singleFile bool) error {
	if !cfg.Diff.OnlyChanged {
		warnLargeStaticReport(tree, logger)
	}
	return htmlreact.NewHtmlReactReportBuilder(cfg.OutputDir, logger, singleFile, cfg).CreateReport(tree)
}

// static HTML gets slow to write and to open above this many files
const largeStaticReport = 5000

func warnLargeStaticReport(tree *model.SummaryTree, logger *slog.Logger) {
	files := countFiles(tree.Root)
	if files <= largeStaticReport {
		return
	}
	logger.Warn(fmt.Sprintf("This static HTML report has %d files; above %d it gets slow to write and to open. "+
		"Set diff.only_changed to report only the changed files, or drop Html from report_types.",
		files, largeStaticReport))
}

func countFiles(dir *model.DirNode) int {
	n := len(dir.Files)
	for _, sub := range dir.Subdirs {
		n += countFiles(sub)
	}
	return n
}

func buildInfo() pipeline.BuildInfo {
	return pipeline.BuildInfo{Version: version, Commit: commit}
}

// applies review.fail_on and returns the reasons to fail the build, empty means pass
func evaluateReviewGate(appConfig *config.AppConfig, summaryTree *model.SummaryTree) []string {
	failOn := appConfig.Review.FailOn
	if failOn == "" || failOn == "never" {
		return nil
	}

	var reasons []string

	result := review.Evaluate(summaryTree, appConfig)
	for _, check := range result.Checks {
		if !check.Passed {
			reasons = append(reasons, fmt.Sprintf("%s is %.1f (limit %.1f)", check.Label, check.Value, check.Threshold))
		}
	}

	changed := diagnostics.OnlyChanged(diagnostics.Extract(summaryTree, appConfig, evaluators.Registry))
	var errorCount, warningCount int
	for _, d := range changed {
		switch d.Severity {
		case diagnostics.SeverityError:
			errorCount++
		case diagnostics.SeverityWarning:
			warningCount++
		}
	}
	if errorCount > 0 {
		reasons = append(reasons, fmt.Sprintf("%d error-severity problem(s) in changed code", errorCount))
	}
	if failOn == "warning" && warningCount > 0 {
		reasons = append(reasons, fmt.Sprintf("%d warning-severity problem(s) in changed code", warningCount))
	}

	return reasons
}

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "serve":
			os.Exit(runServe(os.Args[2:]))
		case "store":
			os.Exit(runStoreCommand(os.Args[2:]))
		case "config":
			os.Exit(runConfigCommand(os.Args[2:]))
		}
	}

	start := time.Now()
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  nanovision [flags]          make the reports of a coverage run\n")
		fmt.Fprintf(os.Stderr, "  nanovision serve [flags]    run the team server on a store (-h for its flags)\n")
		fmt.Fprintf(os.Stderr, "  nanovision store <command>  upkeep of a local store: gc, backup\n")
		fmt.Fprintf(os.Stderr, "  nanovision config <command> the config reference and checks: docs, check, show, schema\n\nFlags:\n")
		flag.PrintDefaults()
	}

	configPath := flag.String("config", "", "Path to the config file (default: nanovision.yaml in this folder; .yml and a leading dot also work).")
	watchFlag := flag.Bool("watch", false, "Enable watch mode to automatically regenerate reports on file changes")
	versionFlag := flag.Bool("version", false, "Print version information and exit")
	listParsersFlag := flag.Bool("list-parsers", false, "List supported coverage report parsers and exit")
	listMetricsFlag := flag.Bool("list-metrics", false, "List all configurable metrics and exit")

	rawInput := parseAndBindFlags()
	flag.Parse()

	if *listParsersFlag {
		factory := pipeline.NewParserFactory(slog.Default(), filereader.NewDefaultReader())

		fmt.Println("Supported Coverage Parsers:")
		for _, name := range factory.RegisteredParsers() {
			fmt.Printf(" - %s\n", name)
		}
		os.Exit(0)
	}

	if *listMetricsFlag {
		printMetrics(os.Stdout)
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

	closer, err := buildLogger(appConfig)
	if err != nil {
		fmt.Fprintln(os.Stderr, "logger init error:", err)
		os.Exit(1)
	}
	if closer != nil {
		defer closer.Close()
	}

	// Print the visual checklist of resolved configurations
	bootlog.PrintBootSummary(appConfig, evaluators.Registry)

	logger := slog.Default()
	identity := resolveIdentity(appConfig, logger)
	if appConfig.History.Store == "" && (appConfig.Run.Revision != "" || appConfig.Run.BaseRevision != "" || appConfig.Run.Kind != "local") {
		logger.Info("The run flags (-run-kind, -revision, -base-revision) have no effect without a run store (-store or history.store)")
	}

	var diffData *diff.DiffData
	if appConfig.Diff.File != "" {
		var parseErr error
		logger.Info("Parsing diff file...", "path", appConfig.Diff.File)
		diffData, parseErr = diff.Parse(appConfig.Diff.File, logger)
		if parseErr != nil {
			logger.Warn("Failed to parse diff file; ignoring diff analysis.", "error", parseErr)
			diffData = nil
		}
	} else if diffData = vcsDiff(identity, logger); diffData == nil {
		logger.Info("No diff file specified in configuration, skipping diff analysis")
	}

	summaryTree, err := pipeline.Run(appConfig, diffData, buildInfo(), logger)
	if err != nil {
		slog.Error("An error occurred during report generation", "error", err)
		os.Exit(1)
	}

	summaryTree.Versions = identity.versions(appConfig.Diff.File, diffData != nil)

	// the comparison lands on the tree before the reports are written, so
	// Summary.txt carries the delta
	hist := recordHistory(appConfig, summaryTree, identity, logger)

	logger.Info("Executing REPORT stage...")
	if err := generateReports(appConfig, summaryTree); err != nil {
		slog.Error("An error occurred during report generation", "error", err)
		os.Exit(1)
	}

	info := textsummary.RunInfo{Revision: identity.current.ID, Stream: identity.stream, Kind: string(identity.kind), Duration: time.Since(start)}
	if hist != nil {
		info.RunID, info.RunURL, info.Notes = hist.runID, hist.runURL, hist.notes
	}
	textsummary.WriteTerminal(os.Stdout, summaryTree, appConfig, info)

	if gateReasons := evaluateReviewGate(appConfig, summaryTree); len(gateReasons) > 0 {
		for _, reason := range gateReasons {
			slog.Error("Review gate failed", "reason", reason)
		}
		slog.Error("Build failed by review gate", "fail_on", appConfig.Review.FailOn)
		os.Exit(2)
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
