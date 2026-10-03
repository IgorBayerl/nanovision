// Package compare calculates the coverage delta between two runs from their
// manifests. It needs no line data: the manifest counts are enough.
package compare

import (
	"math"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/status/evaluators"
	"github.com/IgorBayerl/nanovision/internal/store/blob"
	"github.com/IgorBayerl/nanovision/internal/utils"
)

// Options carry what the manifests do not.
type Options struct {
	// the current run's file filters; base files it rejects are left out
	Include func(path string) bool
	// file_metrics of the current run, in order
	FileMetrics []config.MetricKey
	// the current run was measured with a diff
	HasDiff bool
	Change  *model.ChangeSet
	// report patterns and nanovision versions, for the fairness warnings
	BaseReports, CurrentReports []string
	BaseTool, CurrentTool       string
}

// coverage metrics that the manifest counts can answer
var coverageMetrics = map[config.MetricKey]struct {
	unit  string
	count func(c blob.Counts) (covered, total int)
}{
	config.StatementCoverage: {"statements", func(c blob.Counts) (int, int) { return c.StatementsCovered, c.StatementsValid }},
	config.LineCoverage:      {"lines", func(c blob.Counts) (int, int) { return c.LinesCovered, c.LinesValid }},
	config.MethodsHit:        {"methods", func(c blob.Counts) (int, int) { return c.MethodsHit, c.MethodsValid }},
	config.MethodsFullyCovered: {"methods", func(c blob.Counts) (int, int) {
		return c.MethodsFullyCovered, c.MethodsValid
	}},
	config.StatementMethodsHit: {"methods", func(c blob.Counts) (int, int) {
		return c.StatementMethodsHit, c.StatementMethodsValid
	}},
	config.StatementMethodsFullyCovered: {"methods", func(c blob.Counts) (int, int) {
		return c.StatementMethodsFullyCovered, c.StatementMethodsValid
	}},
}

// share of the base run's coverable lines that may disappear before the delta
// looks like fewer tests ran
const missingLinesWarning = 0.10

// Compare calculates the delta of current against base. Run references,
// distance and exactness are left for the caller, who found the base run.
func Compare(base, current blob.Manifest, opts Options) *model.Comparison {
	cmp := &model.Comparison{Distance: -1, Change: opts.Change}

	var baseTotals, curTotals blob.Counts
	baseByPath := make(map[string]*blob.Entry, len(base.Entries))
	for i := range base.Entries {
		e := &base.Entries[i]
		if opts.Include != nil && !opts.Include(e.Path) {
			cmp.FilteredBaseFiles++
			continue
		}
		baseByPath[e.Path] = e
		baseTotals.Add(e.Counts)
	}
	var curPatch blob.PatchCounts
	for _, e := range current.Entries {
		curTotals.Add(e.Counts)
		if e.Patch != nil {
			curPatch.Add(*e.Patch)
		}
	}

	cmp.Headline = config.LineCoverage
	if baseTotals.StatementsValid > 0 && curTotals.StatementsValid > 0 {
		cmp.Headline = config.StatementCoverage
	}

	keys := []config.MetricKey{cmp.Headline}
	for _, k := range opts.FileMetrics {
		if _, ok := coverageMetrics[k]; ok && k != cmp.Headline {
			keys = append(keys, k)
		}
	}
	for _, k := range keys {
		m := coverageMetrics[k]
		bc, bt := m.count(baseTotals)
		cc, ct := m.count(curTotals)
		if k != cmp.Headline && (bt == 0 || ct == 0) {
			continue
		}
		cmp.Metrics = append(cmp.Metrics, delta(k, m.unit, bc, bt, cc, ct))
	}

	cmp.Deltas = nodeDeltas(baseByPath, current.Entries, keys)

	headline := coverageMetrics[cmp.Headline].count
	var missingLines int
	seen := make(map[string]bool, len(current.Entries))
	for _, e := range current.Entries {
		seen[e.Path] = true
		b, ok := baseByPath[e.Path]
		if !ok {
			cmp.AddedFiles++
			continue
		}
		bc, bt := headline(b.Counts)
		cc, ct := headline(e.Counts)
		if bc == cc && bt == ct {
			continue
		}
		fd := model.FileDelta{
			Path: e.Path, BaseCovered: bc, BaseTotal: bt, CurrentCovered: cc, CurrentTotal: ct,
			Base: percent(bc, bt), Current: percent(cc, ct),
		}
		fd.Delta = round2(fd.Current - fd.Base)
		if opts.HasDiff {
			fd.Change = "not in diff"
			if k := e.Diff.String(); k != "" {
				fd.Change = k
			}
		}
		cmp.Files = append(cmp.Files, fd)
	}
	for path, b := range baseByPath {
		if !seen[path] {
			cmp.RemovedFiles++
			missingLines += b.Counts.LinesValid
		}
	}
	sort.SliceStable(cmp.Files, func(i, j int) bool {
		a, b := math.Abs(cmp.Files[i].Delta), math.Abs(cmp.Files[j].Delta)
		if a != b {
			return a > b
		}
		return cmp.Files[i].Path < cmp.Files[j].Path
	})

	if opts.HasDiff {
		cmp.Patch = &model.PatchSummary{Unit: "statements", Covered: curPatch.StatementsCovered, Total: curPatch.StatementsValid}
		if curPatch.StatementsValid == 0 && curPatch.LinesValid > 0 {
			cmp.Patch = &model.PatchSummary{Unit: "lines", Covered: curPatch.LinesCovered, Total: curPatch.LinesValid}
		}
	}

	cmp.Warnings = fairnessWarnings(opts, missingLines, baseTotals.LinesValid)
	return cmp
}

// nodeDeltas gives the change of each metric, in percentage points, for every
// file and folder whose coverage changed. A file that only one run has gets no
// delta of its own, but it moves the folders above it.
func nodeDeltas(base map[string]*blob.Entry, current []blob.Entry, keys []config.MetricKey) map[string]map[config.MetricKey]float64 {
	type counts struct {
		base, current     blob.Counts
		inBase, inCurrent bool
	}
	nodes := make(map[string]*counts, len(current))
	// add counts a file to its own path and to every folder above it
	add := func(file string, c blob.Counts, isBase bool) {
		for p := file; p != "." && p != "/"; p = path.Dir(p) {
			n := nodes[p]
			if n == nil {
				n = &counts{}
				nodes[p] = n
			}
			if isBase {
				n.base.Add(c)
				n.inBase = true
			} else {
				n.current.Add(c)
				n.inCurrent = true
			}
		}
	}
	for p, e := range base {
		add(p, e.Counts, true)
	}
	for _, e := range current {
		add(e.Path, e.Counts, false)
	}

	deltas := make(map[string]map[config.MetricKey]float64)
	for p, n := range nodes {
		if !n.inBase || !n.inCurrent {
			continue
		}
		for _, k := range keys {
			count := coverageMetrics[k].count
			bc, bt := count(n.base)
			cc, ct := count(n.current)
			if bt == 0 || ct == 0 {
				continue
			}
			if d := round2(percent(cc, ct) - percent(bc, bt)); d != 0 {
				if deltas[p] == nil {
					deltas[p] = make(map[config.MetricKey]float64)
				}
				deltas[p][k] = d
			}
		}
	}
	if len(deltas) == 0 {
		return nil
	}
	return deltas
}

func delta(key config.MetricKey, unit string, bc, bt, cc, ct int) model.MetricDelta {
	d := model.MetricDelta{
		Key: key, Label: Label(key), Unit: unit,
		BaseCovered: bc, BaseTotal: bt, CurrentCovered: cc, CurrentTotal: ct,
		Base: percent(bc, bt), Current: percent(cc, ct),
	}
	d.Delta = round2(d.Current - d.Base)
	return d
}

// Label is the evaluator's name of a metric, for example "Statement Coverage".
func Label(key config.MetricKey) string {
	if ev, ok := evaluators.Registry[key]; ok {
		return ev.Name()
	}
	parts := strings.Split(string(key), "_")
	for i, p := range parts {
		if p != "" {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, " ")
}

func fairnessWarnings(opts Options, missingLines, baseLines int) []string {
	var warnings []string
	if len(opts.BaseReports) > 0 && len(opts.CurrentReports) > 0 && !sameSet(opts.BaseReports, opts.CurrentReports) {
		warnings = append(warnings, "The report patterns differ from the base run, so the runs may not have run the same tests.")
	}
	if baseLines > 0 && float64(missingLines) > missingLinesWarning*float64(baseLines) {
		warnings = append(warnings, strconv.Itoa(int(math.Round(100*float64(missingLines)/float64(baseLines))))+
			"% of the base run's coverable lines are in files this run does not have. Fewer tests may have run.")
	}
	if b, ok := majorMinor(opts.BaseTool); ok {
		if c, ok := majorMinor(opts.CurrentTool); ok && b != c {
			warnings = append(warnings, "The base run was made by nanovision "+opts.BaseTool+" and this run by "+
				opts.CurrentTool+". Metric rules can change between versions.")
		}
	}
	return warnings
}

// majorMinor reads "v1.4.2" or "1.4" as "1.4"; development builds have none.
func majorMinor(version string) (string, bool) {
	parts := strings.SplitN(strings.TrimPrefix(version, "v"), ".", 3)
	if len(parts) < 2 {
		return "", false
	}
	for _, p := range parts[:2] {
		if _, err := strconv.Atoi(p); err != nil {
			return "", false
		}
	}
	return parts[0] + "." + parts[1], true
}

func sameSet(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	x, y := append([]string{}, a...), append([]string{}, b...)
	sort.Strings(x)
	sort.Strings(y)
	for i := range x {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}

// percent truncates like the rest of the reports, so 99.999% never shows as
// 100%. Unlike them, nothing to measure is 0%, not 100%.
func percent(covered, total int) float64 {
	if total == 0 {
		return 0
	}
	return utils.CalculatePercentage(covered, total, 2)
}

func round2(v float64) float64 { return math.Round(v*100) / 100 }
