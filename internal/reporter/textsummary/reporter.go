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
)

type TextReportBuilder struct {
	outputDir string
	logger    *slog.Logger
	config    *config.AppConfig
}

func NewTextReportBuilder(outputDir string, logger *slog.Logger, cfg *config.AppConfig) reporter.ReportBuilder {
	// Guarantee we always have a config and ActiveFileMetrics map, even in tests
	if cfg == nil {
		cfg = config.GetDefaultConfig()
		cfg.FileMetrics = config.DefaultFileMetrics
		cfg.ActiveFileMetrics = make(map[config.MetricKey]bool)
		for _, m := range cfg.FileMetrics {
			cfg.ActiveFileMetrics[m] = true
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
	fmt.Fprintf(f, "Generated on: %s\n", time.Now().Format("02/01/2006 - 15:04:05"))
	if tree.Timestamp > 0 {
		fmt.Fprintf(f, "Coverage date: %s\n", time.Unix(tree.Timestamp, 0).Format("02/01/2006 - 15:04:05"))
	}
	if len(tree.ParserNames) > 0 {
		fmt.Fprintf(f, "Parser: %s\n", strings.Join(tree.ParserNames, " | "))
	}

	fmt.Fprintf(f, "\n")

	// Print metric summaries in the order defined by cfg.FileMetrics.
	for _, key := range b.config.FileMetrics {
		if !b.config.ActiveFileMetrics[key] {
			continue
		}
		if calcData, exists := tree.Metrics.Calculated[key]; exists {
			prettyKey := formatKey(string(key))
			switch v := calcData.(type) {
			case model.CoverageDetail:
				fmt.Fprintf(f, "%s: %.1f%%\n", prettyKey, v.Percentage)
			case model.ScoreDetail:
				fmt.Fprintf(f, "%s: %.1f\n", prettyKey, v.Value)
			}
		}
	}

	// TODO: Create a better visualization of this table, maybe separate in different report
	// Print the hierarchical summary table.
	// tw := tabwriter.NewWriter(f, 0, 0, 2, ' ', 0)
	// defer tw.Flush()

	// fmt.Fprintln(tw) // Newline before the table
	// // Start the recursive walk from the root's children.
	// b.printNode(tw, tree.Root, 0)

	return nil
}

// buildNodeParts builds the metric summary parts for a node in FileMetrics order.
func (b *TextReportBuilder) buildNodeParts(m model.CoverageMetrics) []string {
	var parts []string
	for _, key := range b.config.FileMetrics {
		if !b.config.ActiveFileMetrics[key] {
			continue
		}
		if val, exists := m.Calculated[key]; exists {
			switch v := val.(type) {
			case model.CoverageDetail:
				parts = append(parts, fmt.Sprintf("%.0f%%", v.Percentage))
			case model.ScoreDetail:
				parts = append(parts, fmt.Sprintf("%.0f", v.Value))
			}
		}
	}
	return parts
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
		parts := b.buildNodeParts(sub.Metrics)
		fmt.Fprintf(tw, "%s%s/\t  %s\n", indent, sub.Name, strings.Join(parts, " | "))
		b.printNode(tw, sub, indentLevel+1)
	}
	// Then print files in the current directory.
	for _, file := range sortedFiles {
		parts := b.buildNodeParts(file.Metrics)
		fmt.Fprintf(tw, "%s%s\t  %s\n", indent, file.Name, strings.Join(parts, " | "))
	}
}

func formatKey(s string) string {
	parts := strings.Split(s, "_")
	for i := range parts {
		if len(parts[i]) > 0 {
			parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}
	return strings.Join(parts, " ")
}
