// Package annotations emits coverage diagnostics as a human-readable Markdown
// file - a flat "Problems" list grouped by file, mirroring what an editor's
// problems panel shows. It reads the same diagnostics.Extract output as the
// SARIF reporter, so the two never drift apart.
package annotations

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/diagnostics"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/reporter"
	"github.com/IgorBayerl/nanovision/internal/status"
)

const outputFile = "Annotations.md"

type builder struct {
	outputDir string
	config    *config.AppConfig
	registry  map[config.MetricKey]status.Evaluator
}

// NewAnnotationsReportBuilder creates a Markdown annotations report builder.
func NewAnnotationsReportBuilder(outputDir string, cfg *config.AppConfig, registry map[config.MetricKey]status.Evaluator) reporter.ReportBuilder {
	return &builder{outputDir: outputDir, config: cfg, registry: registry}
}

func (b *builder) ReportType() string { return "Annotations" }

func (b *builder) CreateReport(tree *model.SummaryTree) error {
	diags := diagnostics.Extract(tree, b.config, b.registry)

	outputPath := filepath.Join(b.outputDir, outputFile)
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create annotations file: %w", err)
	}
	defer f.Close()

	errors, warnings := counts(diags)

	fmt.Fprintf(f, "# Coverage Annotations\n\n")
	fmt.Fprintf(f, "%d error(s), %d warning(s)\n\n", errors, warnings)

	if len(diags) == 0 {
		fmt.Fprintf(f, "No problems found. \n")
		return nil
	}

	// Group by file, preserving the already-sorted order.
	var currentFile string
	for _, d := range diags {
		if d.File != currentFile {
			if currentFile != "" {
				fmt.Fprintf(f, "\n")
			}
			currentFile = d.File
			fmt.Fprintf(f, "## %s\n\n", d.File)
		}
		fmt.Fprintf(f, "- %s **L%d**: %s `%s`\n", icon(d.Severity), d.StartLine, d.Message, d.RuleID)
	}
	fmt.Fprintf(f, "\n")
	return nil
}

func counts(diags []diagnostics.Diagnostic) (errors, warnings int) {
	for _, d := range diags {
		switch d.Severity {
		case diagnostics.SeverityError:
			errors++
		case diagnostics.SeverityWarning:
			warnings++
		}
	}
	return errors, warnings
}

func icon(s diagnostics.Severity) string {
	switch s {
	case diagnostics.SeverityError:
		return "🔴"
	case diagnostics.SeverityWarning:
		return "🟡"
	default:
		return "🔵"
	}
}
