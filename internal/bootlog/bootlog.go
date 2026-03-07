package bootlog

import (
	"fmt"
	"sort"
	"strings"

	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/status"
)

// PrintBootSummary logs the final parsed configuration and the active metric checklist.
func PrintBootSummary(cfg *config.AppConfig, registry map[config.MetricKey]status.Evaluator) {
	var sb strings.Builder

	sb.WriteString("\n======================================================\n")
	sb.WriteString(" NanoVision Execution Summary\n")
	sb.WriteString("======================================================\n")
	sb.WriteString(fmt.Sprintf(" Output Directory : %s\n", cfg.OutputDir))
	sb.WriteString(fmt.Sprintf(" Report Types     : %s\n", strings.Join(cfg.ReportTypes, ", ")))

	if cfg.Diff.File != "" {
		sb.WriteString(fmt.Sprintf(" Diff File        : %s (Strip: %s)\n", cfg.Diff.File, cfg.Diff.Strip))
	} else {
		sb.WriteString(" Diff File        : Disabled\n")
	}

	// Sort evaluators alphabetically by Name for consistent display
	var evals []status.Evaluator
	for _, ev := range registry {
		evals = append(evals, ev)
	}
	sort.Slice(evals, func(i, j int) bool { return evals[i].Name() < evals[j].Name() })

	sb.WriteString("\n File Metrics:\n")
	for _, ev := range evals {
		if HasScope(ev, status.FileScope) {
			mark := "[ ]"
			if cfg.ActiveFileMetrics[ev.Key()] {
				mark = "[x]"
			}
			sb.WriteString(fmt.Sprintf("   %s %-30s (%s)\n", mark, ev.Name(), ev.Key()))
		}
	}

	sb.WriteString("\n Method Metrics:\n")
	for _, ev := range evals {
		if HasScope(ev, status.MethodScope) {
			mark := "[ ]"
			if cfg.ActiveMethodMetrics[ev.Key()] {
				mark = "[x]"
			}
			sb.WriteString(fmt.Sprintf("   %s %-30s (%s)\n", mark, ev.Name(), ev.Key()))
		}
	}
	sb.WriteString("======================================================\n")

	// Print directly to standard output for immediate CLI visibility
	fmt.Print(sb.String())
}

// HasScope returns true if the evaluator supports the given scope.
func HasScope(ev status.Evaluator, scope status.MetricScope) bool {
	return ev.SupportedScopes() == scope
}
