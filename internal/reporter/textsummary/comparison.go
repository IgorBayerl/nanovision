package textsummary

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/IgorBayerl/nanovision/internal/model"
)

// how many files the delta block lists
const topFiles = 5

// WriteComparison prints the delta block. It uses ASCII only, so old Windows
// consoles and CI logs show it as it is.
func WriteComparison(w io.Writer, cmp *model.Comparison) {
	if cmp == nil {
		return
	}
	fmt.Fprintln(w, "Coverage delta against the base run")
	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
	fmt.Fprintf(tw, "  Base run\t%s\n", baseRunLine(cmp))
	if cmp.Base.Profile != "" {
		fmt.Fprintf(tw, "  Profile\t%s\n", cmp.Base.Profile)
	}
	if line := changeLine(cmp.Change); line != "" {
		fmt.Fprintf(tw, "  Change\t%s\n", line)
	}
	tw.Flush()
	fmt.Fprintln(w)

	// the metric table and the lines under it share one label column
	type row struct{ label, value string }
	var rows []row
	if h, ok := cmp.HeadlineDelta(); ok {
		rows = append(rows, row{sentenceCase(h.Unit), fmt.Sprintf("%d/%d -> %d/%d  (%s covered, %s valid)",
			h.BaseCovered, h.BaseTotal, h.CurrentCovered, h.CurrentTotal,
			signedInt(h.CurrentCovered-h.BaseCovered), signedInt(h.CurrentTotal-h.BaseTotal))})
	}
	if cmp.Patch != nil {
		rows = append(rows, row{"Patch coverage", patchLine(cmp.Patch)})
	}
	if cmp.AddedFiles > 0 || cmp.RemovedFiles > 0 {
		rows = append(rows, row{"Files", fmt.Sprintf("%d new, %d gone since the base run", cmp.AddedFiles, cmp.RemovedFiles)})
	}
	width := len("Metric")
	for _, m := range cmp.Metrics {
		width = max(width, len(m.Label))
	}
	for _, r := range rows {
		width = max(width, len(r.label))
	}

	fmt.Fprintf(w, "  %-*s   %-9s %-9s %s\n", width, "Metric", "Base", "Current", "Delta")
	for _, m := range cmp.Metrics {
		fmt.Fprintf(w, "  %-*s   %-9s %-9s %s\n", width, sentenceCase(m.Label),
			fmt.Sprintf("%.2f%%", m.Base), fmt.Sprintf("%.2f%%", m.Current), signed(m.Delta, 2))
	}
	if len(rows) > 0 {
		fmt.Fprintln(w)
	}
	for _, r := range rows {
		fmt.Fprintf(w, "  %-*s   %s\n", width, r.label, r.value)
	}

	if len(cmp.Files) > 0 {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "  Files with the largest delta")
		tw = tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
		for _, f := range cmp.Files[:min(topFiles, len(cmp.Files))] {
			fmt.Fprintf(tw, "    %s\t%.1f%% -> %.1f%%\t%s\t%s\n", f.Path, f.Base, f.Current, signed(f.Delta, 1), f.Change)
		}
		tw.Flush()
	}

	if len(cmp.Warnings) > 0 {
		fmt.Fprintln(w)
		for _, warning := range cmp.Warnings {
			fmt.Fprintf(w, "  Warning: %s\n", warning)
		}
	}
}

func baseRunLine(cmp *model.Comparison) string {
	name := cmp.Base.Revision
	if name == "" {
		name = fmt.Sprintf("run %d", cmp.Base.ID)
	}
	switch {
	case cmp.Exact:
		return name + "  (exact)"
	case cmp.Distance > 0:
		return fmt.Sprintf("%s  (%d %s before the base revision)", name, cmp.Distance, plural(cmp.Distance, "revision", "revisions"))
	default:
		return name + "  (older than the base revision)"
	}
}

func changeLine(c *model.ChangeSet) string {
	if c == nil || c.Total() == 0 {
		return ""
	}
	total := c.Total()
	if c.TestOnly() && len(c.NotInReports) == 0 {
		return fmt.Sprintf("%d %s, all ignored by the coverage filters (tests)", total, plural(total, "file", "files"))
	}
	var parts []string
	add := func(n int, what string) {
		if n > 0 {
			parts = append(parts, fmt.Sprintf("%d %s", n, what))
		}
	}
	add(len(c.Measured), "measured")
	add(len(c.Ignored), "ignored by the coverage filters")
	add(len(c.NotInReports), "not in the coverage reports")
	add(len(c.Deleted), "deleted")
	line := fmt.Sprintf("%d %s: %s", total, plural(total, "file", "files"), strings.Join(parts, ", "))
	if c.TestOnly() {
		line += " (no measured code)"
	}
	return line
}

func patchLine(p *model.PatchSummary) string {
	if p.Total == 0 {
		return "none: the change has no measured lines"
	}
	return fmt.Sprintf("%.1f%%  (%d of %d changed %s)", 100*float64(p.Covered)/float64(p.Total), p.Covered, p.Total, p.Unit)
}

func signed(v float64, decimals int) string {
	return fmt.Sprintf("%+.*f", decimals, v)
}

func signedInt(v int) string { return fmt.Sprintf("%+d", v) }

func plural(n int, one, many string) string {
	if n == 1 {
		return one
	}
	return many
}

// "Statement Coverage" -> "Statement coverage"
func sentenceCase(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}
