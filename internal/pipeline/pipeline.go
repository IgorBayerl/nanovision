// Package pipeline runs the stages that turn coverage reports into an
// annotated tree: PARSE, BUILD, ENRICH, DIFF, AGGREGATE, CALCULATE and
// ANNOTATE. Reports, the run store and the delta all start from its result.
package pipeline

import (
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/IgorBayerl/fsglob"
	"github.com/IgorBayerl/nanovision/internal/aggregator"
	"github.com/IgorBayerl/nanovision/internal/analyzer"
	cpp "github.com/IgorBayerl/nanovision/internal/analyzer/cpp"
	"github.com/IgorBayerl/nanovision/internal/analyzer/gdscript"
	golang "github.com/IgorBayerl/nanovision/internal/analyzer/go"
	"github.com/IgorBayerl/nanovision/internal/cache"
	"github.com/IgorBayerl/nanovision/internal/calculator"
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/diff"
	"github.com/IgorBayerl/nanovision/internal/diffapply"
	"github.com/IgorBayerl/nanovision/internal/enricher"
	"github.com/IgorBayerl/nanovision/internal/filereader"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/parsers"
	"github.com/IgorBayerl/nanovision/internal/parsers/parser_cobertura"
	"github.com/IgorBayerl/nanovision/internal/parsers/parser_gcov"
	"github.com/IgorBayerl/nanovision/internal/parsers/parser_gocover"
	"github.com/IgorBayerl/nanovision/internal/parsers/parser_lcov"
	"github.com/IgorBayerl/nanovision/internal/status"
	"github.com/IgorBayerl/nanovision/internal/status/evaluators"
	"github.com/IgorBayerl/nanovision/internal/tree"
)

// BuildInfo identifies the binary, for the analysis cache.
type BuildInfo struct {
	Version string
	Commit  string
}

// Dev is true for a build without a release commit.
func (b BuildInfo) Dev() bool { return b.Commit == "dev" || b.Commit == "none" || b.Commit == "" }

// Run executes every stage up to ANNOTATE. diffData may be nil.
func Run(cfg *config.AppConfig, diffData *diff.DiffData, build BuildInfo, logger *slog.Logger) (*model.SummaryTree, error) {
	logger.Info("Executing report generation pipeline...")

	reader := filereader.NewDefaultReader()
	parserFactory := NewParserFactory(logger, reader)
	treeBuilder := tree.NewBuilder(cfg.ProjectRoot, cfg.FileFilterInstance)

	allAnalyzers := []analyzer.Analyzer{
		golang.New(),
		cpp.New(),
		gdscript.New(),
	}
	buildMeta := cache.BuildMetadata{CommitHash: build.Commit, AnalyzerVersion: build.Version}
	cacheManager := setupCacheManager(cfg, logger, build, buildMeta)
	treeEnricher := enricher.New(allAnalyzers, reader, logger, cacheManager, buildMeta)

	if len(cfg.InputPairs) == 0 {
		return nil, fmt.Errorf("no valid report pattern and source directory pairs were provided")
	}

	logger.Info("Executing PARSE stage...")
	parserResults, err := parseReportFiles(logger, cfg, parserFactory)
	if err != nil {
		return nil, err
	}
	logger.Info("PARSE stage completed successfully.", "parsed_report_sets", len(parserResults))

	logger.Info("Executing BUILD stage...")
	summaryTree, err := treeBuilder.BuildTree(parserResults)
	if err != nil {
		return nil, fmt.Errorf("failed to build and aggregate coverage tree: %w", err)
	}
	logger.Info("BUILD stage completed successfully.")

	logger.Info("Executing ENRICH stage...")
	treeEnricher.EnrichTree(summaryTree)
	logger.Info("ENRICH stage completed successfully.")

	if diffData != nil {
		logger.Info("Executing DIFF ANALYSIS stage...")
		diffapply.Apply(summaryTree, diffData, cfg.FileFilterInstance, logger)
		logger.Info("DIFF ANALYSIS stage completed successfully.")
	}

	aggregator.AggregateMetricsAfterEnrichment(summaryTree)

	logger.Info("Executing CALCULATE stage...")
	calculator.CalculateTree(summaryTree, cfg.ActiveFileMetrics, cfg.ActiveMethodMetrics)
	logger.Info("CALCULATE stage completed successfully.")

	logger.Info("Executing ANNOTATE stage...")
	status.Annotate(summaryTree, cfg, status.DeriveCapabilities(summaryTree), evaluators.Registry)
	logger.Info("ANNOTATE stage completed successfully.")

	return summaryTree, nil
}

func NewParserFactory(logger *slog.Logger, reader filereader.Reader) *parsers.ParserFactory {
	return parsers.NewParserFactory(
		logger,
		parser_cobertura.NewCoberturaParser(reader),
		parser_gocover.NewGoCoverParser(reader),
		parser_gcov.NewGCovParser(reader),
		parser_lcov.NewLcovParser(reader),
	)
}

func parseReportFiles(logger *slog.Logger, cfg *config.AppConfig, parserFactory *parsers.ParserFactory) ([]*parsers.ParserResult, error) {
	var parserResults []*parsers.ParserResult
	var parserErrors []string
	var totalFilesParsed int

	for _, pair := range cfg.InputPairs {
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
				FileFilter: cfg.FileFilterInstance,
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

// Try to find a cache directory, if it cant will continue anyway without cache
// Fail soft, cache is not necessary for the application to work, its just good
// But we throw a warning
func setupCacheManager(cfg *config.AppConfig, logger *slog.Logger, build BuildInfo, buildMeta cache.BuildMetadata) *cache.Manager {
	if cfg.IgnoreCache {
		logger.Debug("Cache ignored by user configuration.")
		return nil
	}

	// Ask the cache package to find the best writable directory (3-path fallback)
	cacheDir, err := cache.DetermineCacheDir(logger)
	if err != nil {
		logger.Warn("Could not determine a writable cache directory; caching will be disabled.", "error", err)
		return nil
	}

	// SELECT VALIDATOR BASED ON BUILD TYPE
	var validator cache.CacheValidator
	if build.Dev() {
		validator = &cache.DevValidator{} // Always invalidate in dev
		logger.Info("Dev mode: cache will be invalidated on each run")
	} else {
		validator = &cache.StrictValidator{CurrentBuildMetadata: buildMeta}
		logger.Info("Production mode: cache validated by commit hash and analyzer version")
	}

	// Initialize the manager at that directory
	manager, err := cache.NewManager(cacheDir, logger, validator)
	if err != nil {
		logger.Warn("Failed to initialize cache manager; proceeding without cache.", "error", err)
		return nil
	}

	return manager
}
