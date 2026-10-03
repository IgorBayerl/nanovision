// Package review scores the changed part of an annotated coverage tree:
// gate thresholds and changelist stats.
package review

import (
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
)

// GateCheck is one evaluated threshold from config.ReviewGate.
type GateCheck struct {
	Key       string  `json:"key"`       // stable identifier, e.g. "patch_statement_coverage"
	Label     string  `json:"label"`     // human-readable name
	Value     float64 `json:"value"`     // measured value
	Threshold float64 `json:"threshold"` // configured limit
	Passed    bool    `json:"passed"`
}

// Stats summarizes the changelist.
type Stats struct {
	ChangedFiles           int `json:"changedFiles"`
	MethodsAdded           int `json:"methodsAdded"`
	MethodsModified        int `json:"methodsModified"`
	UntestedChangedMethods int `json:"untestedChangedMethods"` // changed methods with zero covered changed statements/lines
	PatchStatementsValid   int `json:"patchStatementsValid"`
	PatchStatementsCovered int `json:"patchStatementsCovered"`
	MaxChangedComplexity   int `json:"maxChangedComplexity"`
}

// Result is the full review evaluation.
type Result struct {
	// every check passed, or no gate was configured
	Passed bool        `json:"passed"`
	Checks []GateCheck `json:"checks,omitempty"`
	Stats  Stats       `json:"stats"`
}

// Evaluate never returns nil. Without diff data every stat is zero.
func Evaluate(tree *model.SummaryTree, cfg *config.AppConfig) *Result {
	res := &Result{Passed: true}
	if tree == nil || tree.Root == nil {
		return res
	}

	walk(tree.Root, func(file *model.FileNode) {
		if file.Diff != nil && file.Diff.Kind != model.ChangeKindNone {
			res.Stats.ChangedFiles++
		}
		for i := range file.Methods {
			m := &file.Methods[i]
			if m.DiffStatus == "" {
				continue
			}
			switch m.DiffStatus {
			case "added":
				res.Stats.MethodsAdded++
			case "modified":
				res.Stats.MethodsModified++
			}
			if isUntested(m) {
				res.Stats.UntestedChangedMethods++
			}
			if m.CyclomaticComplexity != nil && *m.CyclomaticComplexity > res.Stats.MaxChangedComplexity {
				res.Stats.MaxChangedComplexity = *m.CyclomaticComplexity
			}
		}
	})

	res.Stats.PatchStatementsValid = tree.Metrics.PatchStatementsValid
	res.Stats.PatchStatementsCovered = tree.Metrics.PatchStatementsCovered

	res.Checks = evaluateGate(res.Stats, tree, cfg.Review.Gate)
	for _, c := range res.Checks {
		if !c.Passed {
			res.Passed = false
		}
	}
	return res
}

func evaluateGate(stats Stats, tree *model.SummaryTree, gate config.ReviewGate) []GateCheck {
	var checks []GateCheck

	if gate.PatchStatementCoverage != nil && stats.PatchStatementsValid > 0 {
		pct := 100.0 * float64(stats.PatchStatementsCovered) / float64(stats.PatchStatementsValid)
		checks = append(checks, GateCheck{
			Key:       "patch_statement_coverage",
			Label:     "Patch statement coverage",
			Value:     pct,
			Threshold: *gate.PatchStatementCoverage,
			Passed:    pct >= *gate.PatchStatementCoverage,
		})
	}

	if gate.MaxChangedMethodComplexity != nil && (stats.MethodsAdded+stats.MethodsModified) > 0 {
		checks = append(checks, GateCheck{
			Key:       "max_changed_method_complexity",
			Label:     "Max changed-method complexity",
			Value:     float64(stats.MaxChangedComplexity),
			Threshold: float64(*gate.MaxChangedMethodComplexity),
			Passed:    stats.MaxChangedComplexity <= *gate.MaxChangedMethodComplexity,
		})
	}

	return checks
}

// true only when the method has coverable changed code and none of it runs
func isUntested(m *model.MethodMetrics) bool {
	if m.PatchStatementsValid > 0 {
		return m.PatchStatementsCovered == 0
	}
	if m.PatchLinesValid > 0 {
		return m.PatchLinesCovered == 0
	}
	return false
}

func walk(dir *model.DirNode, fn func(*model.FileNode)) {
	for _, sub := range dir.Subdirs {
		walk(sub, fn)
	}
	for _, f := range dir.Files {
		fn(f)
	}
}
