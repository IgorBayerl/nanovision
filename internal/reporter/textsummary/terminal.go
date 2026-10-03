package textsummary

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/diagnostics"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/review"
	"github.com/IgorBayerl/nanovision/internal/status/evaluators"
)

// RunInfo is what the terminal summary says about the run beyond the tree.
type RunInfo struct {
	RunID    int64
	RunURL   string
	Revision string
	Stream   string
	Kind     string
	Duration time.Duration
	// one line each, e.g. why there is no delta
	Notes []string
}

// how many file metrics the summary lists
const terminalMetrics = 4

// WriteTerminal prints the short run summary for standard output: the
// headline numbers, the changed code, the delta, and where to look next.
// ASCII only, like the delta block.
func WriteTerminal(w io.Writer, tree *model.SummaryTree, cfg *config.AppConfig, info RunInfo) {
	header := []string{"nanovision"}
	if info.RunID > 0 {
		header = append(header, fmt.Sprintf("run %d", info.RunID))
	}
	if where := revisionLabel(info); where != "" {
		header = append(header, where)
	}
	if info.Duration > 0 {
		header = append(header, info.Duration.Round(100*time.Millisecond).String())
	}
	fmt.Fprintln(w)
	fmt.Fprintln(w, strings.Join(header, "  |  "))

	type row struct{ label, value string }
	var rows []row
	// the headline metric leads, as in the delta block
	keys := []config.MetricKey{config.LineCoverage}
	if tree.Metrics.StatementsValid > 0 {
		keys[0] = config.StatementCoverage
	}
	for _, key := range cfg.FileMetrics {
		if key != keys[0] {
			keys = append(keys, key)
		}
	}
	for _, key := range keys {
		if len(rows) == terminalMetrics {
			break
		}
		if strings.HasPrefix(string(key), "patch_") {
			continue
		}
		switch v := tree.Metrics.Calculated[key].(type) {
		case model.CoverageDetail:
			rows = append(rows, row{metricName(key), fmt.Sprintf("%6.2f%%  (%d of %d)", v.Percentage, v.Covered, v.Total)})
		case model.ScoreDetail:
			rows = append(rows, row{metricName(key), fmt.Sprintf("%6.0f", v.Value)})
		}
	}

	if tree.Change != nil {
		m := tree.Metrics
		switch {
		case m.PatchStatementsValid > 0:
			rows = append(rows, row{"Patch coverage", fmt.Sprintf("%6.2f%%  (%d of %d changed statements)",
				100*float64(m.PatchStatementsCovered)/float64(m.PatchStatementsValid), m.PatchStatementsCovered, m.PatchStatementsValid)})
		case m.PatchLinesValid > 0:
			rows = append(rows, row{"Patch coverage", fmt.Sprintf("%6.2f%%  (%d of %d changed lines)",
				100*float64(m.PatchLinesCovered)/float64(m.PatchLinesValid), m.PatchLinesCovered, m.PatchLinesValid)})
		default:
			rows = append(rows, row{"Patch coverage", "none: the change has no measured lines"})
		}

		res := review.Evaluate(tree, cfg)
		if len(res.Checks) > 0 {
			verdict := "passed"
			for _, c := range res.Checks {
				if !c.Passed {
					verdict = fmt.Sprintf("failed: %s is %.1f (limit %.1f)", c.Label, c.Value, c.Threshold)
					break
				}
			}
			rows = append(rows, row{"Review gate", verdict})
		}
		if res.Stats.UntestedChangedMethods > 0 {
			rows = append(rows, row{"Risk", fmt.Sprintf("%d changed %s with no covered change",
				res.Stats.UntestedChangedMethods, plural(res.Stats.UntestedChangedMethods, "method", "methods"))})
		}
		errs, warns := 0, 0
		for _, d := range diagnostics.OnlyChanged(diagnostics.Extract(tree, cfg, evaluators.Registry)) {
			switch d.Severity {
			case diagnostics.SeverityError:
				errs++
			case diagnostics.SeverityWarning:
				warns++
			}
		}
		if errs+warns > 0 {
			rows = append(rows, row{"Problems", fmt.Sprintf("%d %s, %d %s in changed code",
				errs, plural(errs, "error", "errors"), warns, plural(warns, "warning", "warnings"))})
		}
	}

	width := 0
	for _, r := range rows {
		width = max(width, len(r.label))
	}
	for _, r := range rows {
		fmt.Fprintf(w, "  %-*s   %s\n", width, r.label, r.value)
	}

	if tree.Comparison != nil {
		fmt.Fprintln(w)
		WriteComparison(w, tree.Comparison)
	}
	if len(info.Notes) > 0 {
		fmt.Fprintln(w)
		for _, n := range info.Notes {
			fmt.Fprintf(w, "  %s\n", n)
		}
	}

	if info.RunURL != "" {
		fmt.Fprintf(w, "\n  Open: %s\n", info.RunURL)
	}
	fmt.Fprintln(w)
}

func revisionLabel(info RunInfo) string {
	rev := info.Revision
	if len(rev) == 40 && !strings.Contains(rev, "@") {
		rev = rev[:10] // a git hash
	}
	switch {
	case info.Stream != "" && rev != "" && !strings.HasPrefix(rev, info.Stream):
		return info.Stream + " @ " + rev
	case rev != "":
		return rev
	default:
		return info.Stream
	}
}

func metricName(key config.MetricKey) string {
	if ev, ok := evaluators.Registry[key]; ok {
		return sentenceCase(ev.Name())
	}
	return sentenceCase(strings.ReplaceAll(string(key), "_", " "))
}
