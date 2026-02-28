package textsummary

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/reporter"
	"github.com/IgorBayerl/nanovision/internal/utils"
)

type TextReportBuilder struct {
	outputDir string
	logger    *slog.Logger
	config    *config.AppConfig
}

func NewTextReportBuilder(outputDir string, logger *slog.Logger, cfg *config.AppConfig) reporter.ReportBuilder {
	// Guarantee we always have a config and ActiveMetrics map, even in tests
	if cfg == nil {
		cfg = config.GetDefaultConfig()
		cfg.DisplayMetrics = config.DefaultDisplayMetrics
		cfg.ActiveMetrics = make(map[config.MetricKey]bool)
		for _, m := range cfg.DisplayMetrics {
			cfg.ActiveMetrics[m] = true
		}
	}

	return &TextReportBuilder{
		outputDir: outputDir,
		logger:    logger,
		config:    cfg,
	}
}

func (b *TextReportBuilder) ReportType() string {
	return "TextSummary"
}

// CreateReport now accepts the new model.SummaryTree.
func (b *TextReportBuilder) CreateReport(tree *model.SummaryTree) error {
	outputPath := filepath.Join(b.outputDir, "Summary.txt")
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create report file: %w", err)
	}
	defer f.Close()

	b.logger.Info("Writing text summary to file", "path", outputPath)

	// Print top-level summary information.
	fmt.Fprintf(f, "Summary\n")
	fmt.Fprintf(f, "  Generated on: %s\n", time.Now().Format("02/01/2006 - 15:04:05"))
	if tree.Timestamp > 0 {
		fmt.Fprintf(f, "  Coverage date: %s\n", time.Unix(tree.Timestamp, 0).Format("02/01/2006 - 15:04:05"))
	}
	if len(tree.ParserNames) > 0 {
		fmt.Fprintf(f, "  Parser: %s\n", strings.Join(tree.ParserNames, " | "))
	}

	if b.config.ActiveMetrics[config.StatementCoverage] {
		statementCoverage := utils.CalculatePercentage(tree.Metrics.StatementsCovered, tree.Metrics.StatementsValid, 1)
		fmt.Fprintf(f, "  Statement coverage: %s\n", utils.FormatPercentage(statementCoverage, 0))
		fmt.Fprintf(f, "  Covered statements: %d\n", tree.Metrics.StatementsCovered)
		fmt.Fprintf(f, "  Uncovered statements: %d\n", tree.Metrics.StatementsValid-tree.Metrics.StatementsCovered)
		fmt.Fprintf(f, "  Total statements: %d\n", tree.Metrics.StatementsValid)
	}

	if b.config.ActiveMetrics[config.LineCoverage] {
		lineCoverage := utils.CalculatePercentage(tree.Metrics.LinesCovered, tree.Metrics.LinesValid, 1)
		fmt.Fprintf(f, "  Line coverage: %s\n", utils.FormatPercentage(lineCoverage, 0))
		fmt.Fprintf(f, "  Covered lines: %d\n", tree.Metrics.LinesCovered)
		fmt.Fprintf(f, "  Uncovered lines: %d\n", tree.Metrics.LinesValid-tree.Metrics.LinesCovered)
		fmt.Fprintf(f, "  Coverable lines: %d\n", tree.Metrics.LinesValid)
	}

	if tree.Metrics.BranchesValid > 0 && b.config.ActiveMetrics[config.BranchCoverage] {
		branchCoverage := utils.CalculatePercentage(tree.Metrics.BranchesCovered, tree.Metrics.BranchesValid, 1)
		fmt.Fprintf(f, "  Branch coverage: %s (%d of %d)\n", utils.FormatPercentage(branchCoverage, 0), tree.Metrics.BranchesCovered, tree.Metrics.BranchesValid)
	}

	// Print the hierarchical summary table.
	tw := tabwriter.NewWriter(f, 0, 0, 2, ' ', 0)
	defer tw.Flush()

	fmt.Fprintln(tw) // Newline before the table
	// Start the recursive walk from the root's children.
	b.printNode(tw, tree.Root, 0)

	return nil
}

// printNode is a recursive helper to print the tree hierarchy.
func (b *TextReportBuilder) printNode(tw *tabwriter.Writer, dir *model.DirNode, indentLevel int) {
	indent := strings.Repeat("  ", indentLevel)

	// Sort subdirectories by name for consistent output.
	sortedSubdirs := make([]*model.DirNode, 0, len(dir.Subdirs))
	for _, sub := range dir.Subdirs {
		sortedSubdirs = append(sortedSubdirs, sub)
	}
	sort.Slice(sortedSubdirs, func(i, j int) bool {
		return sortedSubdirs[i].Name < sortedSubdirs[j].Name
	})

	// Sort files by name for consistent output.
	sortedFiles := make([]*model.FileNode, 0, len(dir.Files))
	for _, file := range dir.Files {
		sortedFiles = append(sortedFiles, file)
	}
	sort.Slice(sortedFiles, func(i, j int) bool {
		return sortedFiles[i].Name < sortedFiles[j].Name
	})

	// Print subdirectories first.
	for _, sub := range sortedSubdirs {
		stmtCov := utils.CalculatePercentage(sub.Metrics.StatementsCovered, sub.Metrics.StatementsValid, 1)
		lineCov := utils.CalculatePercentage(sub.Metrics.LinesCovered, sub.Metrics.LinesValid, 1)

		var parts []string
		if b.config.ActiveMetrics[config.StatementCoverage] {
			parts = append(parts, fmt.Sprintf("%s (Stmt)", utils.FormatPercentage(stmtCov, 0)))
		}
		if b.config.ActiveMetrics[config.LineCoverage] {
			parts = append(parts, fmt.Sprintf("%s (Line)", utils.FormatPercentage(lineCov, 0)))
		}

		fmt.Fprintf(tw, "%s%s/\t  %s\n", indent, sub.Name, strings.Join(parts, " | "))
		b.printNode(tw, sub, indentLevel+1)
	}

	// Then print files in the current directory.
	for _, file := range sortedFiles {
		stmtCov := utils.CalculatePercentage(file.Metrics.StatementsCovered, file.Metrics.StatementsValid, 1)
		lineCov := utils.CalculatePercentage(file.Metrics.LinesCovered, file.Metrics.LinesValid, 1)

		var parts []string
		if b.config.ActiveMetrics[config.StatementCoverage] {
			parts = append(parts, fmt.Sprintf("%s (Stmt)", utils.FormatPercentage(stmtCov, 0)))
		}
		if b.config.ActiveMetrics[config.LineCoverage] {
			parts = append(parts, fmt.Sprintf("%s (Line)", utils.FormatPercentage(lineCov, 0)))
		}

		fmt.Fprintf(tw, "%s%s\t  %s\n", indent, file.Name, strings.Join(parts, " | "))
	}
}
