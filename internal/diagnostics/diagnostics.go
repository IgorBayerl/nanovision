// Package diagnostics is the single source of truth that converts the
// risk-annotated coverage tree (model.SummaryTree, already carrying per-node
// Statuses from the status.Annotate stage) into a flat list of editor-style
// diagnostics ("problems").
//
// Every consumer - the SARIF reporter, the Markdown reporter, the HTML
// problems panel, and future editor extensions - reads this same list instead
// of re-deriving rules. This keeps the "warning/error" logic centralized and
// consistent everywhere, which is the whole point: generate once, render many.
package diagnostics

import (
	"fmt"
	"sort"
	"strings"

	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/status"
)

// Severity is the editor-facing severity, aligned with what LSP / SARIF /
// the VSCode Problems tab understand.
type Severity string

const (
	SeverityError   Severity = "error"   // maps from RiskDanger
	SeverityWarning Severity = "warning" // maps from RiskWarning
	SeverityInfo    Severity = "info"    // reserved for future informational hints
)

// Scope tells whether the diagnostic came from a whole-file metric or a
// single method/function.
type Scope string

const (
	ScopeFile   Scope = "file"
	ScopeMethod Scope = "method"
)

// Diagnostic is one editor problem. Line numbers are 1-based. File is a
// project-root-relative, forward-slash path (matches the model tree), which is
// exactly what SARIF and editors expect for workspace-relative locations.
type Diagnostic struct {
	RuleID    string   `json:"ruleId"`    // metric key, e.g. "line_coverage"
	RuleName  string   `json:"ruleName"`  // human name, e.g. "Line Coverage"
	Severity  Severity `json:"severity"`  // error | warning | info
	File      string   `json:"file"`      // project-relative path, slash-separated
	StartLine int      `json:"startLine"` // 1-based
	EndLine   int      `json:"endLine"`   // 1-based, >= StartLine
	Message   string   `json:"message"`
	Scope     Scope    `json:"scope"`
	// "added" or "modified" when the anchor is in the diff, "" otherwise.
	// the anchor is the file for file scope, the method for method scope.
	DiffStatus string `json:"diffStatus,omitempty"`
}

// Extract walks the annotated tree and produces one diagnostic per file-level
// and method-level status that classifies as danger or warning. "safe"
// statuses are intentionally skipped - a passing metric is not a problem.
//
// registry is the evaluator registry (evaluators.Registry) passed in by the
// caller to avoid an import cycle; it supplies human-readable rule names.
func Extract(tree *model.SummaryTree, cfg *config.AppConfig, registry map[config.MetricKey]status.Evaluator) []Diagnostic {
	if tree == nil || tree.Root == nil {
		return nil
	}

	var out []Diagnostic
	walk(tree.Root, func(dir *model.DirNode) {
		for _, file := range dir.Files {
			fileDiff := ""
			if file.Diff != nil {
				fileDiff = file.Diff.Kind.String()
			}
			// File-level statuses -> anchored at the top of the file.
			for key, level := range file.Statuses {
				if d, ok := build(key, level, file.Metrics.Calculated, cfg, registry, ScopeFile, file.Path, 1, 1); ok {
					d.DiffStatus = fileDiff
					out = append(out, d)
				}
			}
			// Method-level statuses -> anchored at the method's line range.
			for i := range file.Methods {
				m := &file.Methods[i]
				start, end := lineRange(m)
				for key, level := range m.Statuses {
					if d, ok := build(key, level, m.Calculated, cfg, registry, ScopeMethod, file.Path, start, end); ok {
						d.Message = fmt.Sprintf("%s in %s", d.Message, m.Name)
						d.DiffStatus = m.DiffStatus
						out = append(out, d)
					}
				}
			}
		}
	})

	sortDiagnostics(out)
	return out
}

// OnlyChanged keeps the diagnostics anchored to diffed code, in order.
func OnlyChanged(ds []Diagnostic) []Diagnostic {
	var out []Diagnostic
	for _, d := range ds {
		if d.DiffStatus != "" {
			out = append(out, d)
		}
	}
	return out
}

// ok is false for a "safe" or unknown level, which is not a problem worth showing
func build(
	key config.MetricKey,
	level string,
	calc map[config.MetricKey]any,
	cfg *config.AppConfig,
	registry map[config.MetricKey]status.Evaluator,
	scope Scope,
	file string,
	start, end int,
) (Diagnostic, bool) {
	sev, ok := severityFor(level)
	if !ok {
		return Diagnostic{}, false
	}

	name := ruleName(key, registry)
	msg := name
	if v := valueString(calc, key); v != "" {
		msg += " " + v
	}
	if h := bandHint(cfg.StatusBands, key); h != "" {
		msg += fmt.Sprintf(" (threshold %s)", h)
	}

	return Diagnostic{
		RuleID:    string(key),
		RuleName:  name,
		Severity:  sev,
		File:      file,
		StartLine: start,
		EndLine:   end,
		Message:   msg,
		Scope:     scope,
	}, true
}

func severityFor(level string) (Severity, bool) {
	switch level {
	case string(status.RiskDanger):
		return SeverityError, true
	case string(status.RiskWarning):
		return SeverityWarning, true
	default:
		return "", false
	}
}

func ruleName(key config.MetricKey, registry map[config.MetricKey]status.Evaluator) string {
	if ev, ok := registry[key]; ok {
		return ev.Name()
	}
	return prettifyKey(string(key))
}

// valueString formats the actual measured value so the problem message shows
// *why* it tripped (e.g. "42.0%" or a raw complexity score).
func valueString(calc map[config.MetricKey]any, key config.MetricKey) string {
	if calc == nil {
		return ""
	}
	v, ok := calc[key]
	if !ok {
		return ""
	}
	switch d := v.(type) {
	case model.CoverageDetail:
		return fmt.Sprintf("%.1f%%", d.Percentage)
	case model.ScoreDetail:
		return fmt.Sprintf("%.0f", d.Value)
	default:
		return ""
	}
}

func bandHint(bands config.StatusBands, key config.MetricKey) string {
	if b, ok := bands[key]; ok {
		return fmt.Sprintf("%.0f..%.0f", b.Min, b.Max)
	}
	return ""
}

func lineRange(m *model.MethodMetrics) (int, int) {
	start := m.StartLine
	if start < 1 {
		start = 1
	}
	end := m.EndLine
	if end < start {
		end = start
	}
	return start, end
}

func prettifyKey(s string) string {
	parts := strings.Split(s, "_")
	for i := range parts {
		if len(parts[i]) > 0 {
			parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}
	return strings.Join(parts, " ")
}

// sortDiagnostics gives deterministic output: by file, then line, then rule.
// Errors are not floated above warnings here so that per-file ordering stays
// natural; consumers that want severity grouping can re-sort.
func sortDiagnostics(ds []Diagnostic) {
	sort.SliceStable(ds, func(i, j int) bool {
		if ds[i].File != ds[j].File {
			return ds[i].File < ds[j].File
		}
		if ds[i].StartLine != ds[j].StartLine {
			return ds[i].StartLine < ds[j].StartLine
		}
		return ds[i].RuleID < ds[j].RuleID
	})
}

// walk performs a recursive traversal of the directory tree (dirs first),
// matching the traversal used by the status annotator.
func walk(n *model.DirNode, fn func(*model.DirNode)) {
	for _, c := range n.Subdirs {
		walk(c, fn)
	}
	fn(n)
}
