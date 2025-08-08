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

	"github.com/IgorBayerl/AdlerCov/internal/model"
	"github.com/IgorBayerl/AdlerCov/internal/reporter"
	"github.com/IgorBayerl/AdlerCov/internal/utils"
)

type TextReportBuilder struct {
	outputDir string
	logger    *slog.Logger
}

func NewTextReportBuilder(outputDir string, logger *slog.Logger) reporter.ReportBuilder {
	return &TextReportBuilder{
		outputDir: outputDir,
		logger:    logger,
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
	fmt.Fprintf(f, "  Parser: %s\n", tree.ParserName)

	lineCoverage := utils.CalculatePercentage(tree.Metrics.LinesCovered, tree.Metrics.LinesValid, 1)
	fmt.Fprintf(f, "  Line coverage: %s\n", utils.FormatPercentage(lineCoverage, 0))
	fmt.Fprintf(f, "  Covered lines: %d\n", tree.Metrics.LinesCovered)
	fmt.Fprintf(f, "  Uncovered lines: %d\n", tree.Metrics.LinesValid-tree.Metrics.LinesCovered)
	fmt.Fprintf(f, "  Coverable lines: %d\n", tree.Metrics.LinesValid)

	if tree.Metrics.BranchesValid > 0 {
		branchCoverage := utils.CalculatePercentage(tree.Metrics.BranchesCovered, tree.Metrics.BranchesValid, 1)
		fmt.Fprintf(f, "  Branch coverage: %s (%d of %d)\n", utils.FormatPercentage(branchCoverage, 0), tree.Metrics.BranchesCovered, tree.Metrics.BranchesValid)
	}

	// Print the hierarchical summary table.
	tw := tabwriter.NewWriter(f, 0, 0, 2, ' ', 0)
	defer tw.Flush()

	fmt.Fprintln(tw) // Newline before the table
	// Start the recursive walk from the root's children.
	printNode(tw, tree.Root, 0)

	return nil
}

// printNode is a recursive helper to print the tree hierarchy.
func printNode(tw *tabwriter.Writer, dir *model.DirNode, indentLevel int) {
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
		lineCov := utils.CalculatePercentage(sub.Metrics.LinesCovered, sub.Metrics.LinesValid, 1)
		fmt.Fprintf(tw, "%s%s/\t  %s\n", indent, sub.Name, utils.FormatPercentage(lineCov, 0))
		printNode(tw, sub, indentLevel+1)
	}

	// Then print files in the current directory.
	for _, file := range sortedFiles {
		lineCov := utils.CalculatePercentage(file.Metrics.LinesCovered, file.Metrics.LinesValid, 1)
		fmt.Fprintf(tw, "%s%s\t  %s\n", indent, file.Name, utils.FormatPercentage(lineCov, 0))
	}
}
