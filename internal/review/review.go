// Package review scores the changed part of an annotated coverage tree:
// gate thresholds, changelist stats, and risk hotspots.
package review

import (
	"sort"

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

// Hotspot is a changed method ranked by exposed risk (CC * uncovered ratio).
type Hotspot struct {
	File       string   `json:"file"`
	Method     string   `json:"method"`
	StartLine  int      `json:"startLine"`
	DiffStatus string   `json:"diffStatus"`
	Complexity *float64 `json:"complexity,omitempty"`
	// coverage of the changed lines only, nil when none are coverable
	PatchCoverage *float64 `json:"patchCoverage,omitempty"`
	// CC * (1 - coverage). higher is riskier
	Risk float64 `json:"risk"`
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
	Passed   bool        `json:"passed"`
	Checks   []GateCheck `json:"checks,omitempty"`
	Stats    Stats       `json:"stats"`
	Hotspots []Hotspot   `json:"hotspots,omitempty"`
}

// Evaluate never returns nil. Without diff data every stat is zero.
func Evaluate(tree *model.SummaryTree, cfg *config.AppConfig) *Result {
	res := &Result{Passed: true}
	if tree == nil || tree.Root == nil {
		return res
	}

	var hotspots []Hotspot
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
			if h, ok := hotspot(file, m); ok {
				hotspots = append(hotspots, h)
			}
		}
	})

	res.Stats.PatchStatementsValid = tree.Metrics.PatchStatementsValid
	res.Stats.PatchStatementsCovered = tree.Metrics.PatchStatementsCovered

	sort.SliceStable(hotspots, func(i, j int) bool {
		if hotspots[i].Risk != hotspots[j].Risk {
			return hotspots[i].Risk > hotspots[j].Risk
		}
		if hotspots[i].File != hotspots[j].File {
			return hotspots[i].File < hotspots[j].File
		}
		return hotspots[i].StartLine < hotspots[j].StartLine
	})
	limit := cfg.Review.Hotspots
	if limit <= 0 {
		limit = 10
	}
	if len(hotspots) > limit {
		hotspots = hotspots[:limit]
	}
	res.Hotspots = hotspots

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

// a method with no complexity data ranks on uncovered ratio alone, CC counts as 1.
func hotspot(file *model.FileNode, m *model.MethodMetrics) (Hotspot, bool) {
	cov, hasCov := methodCoverage(m)
	if !hasCov {
		return Hotspot{}, false
	}

	cc := 1.0
	var ccPtr *float64
	if m.CyclomaticComplexity != nil {
		cc = float64(*m.CyclomaticComplexity)
		ccPtr = &cc
	}

	h := Hotspot{
		File:       file.Path,
		Method:     m.Name,
		StartLine:  m.StartLine,
		DiffStatus: m.DiffStatus,
		Complexity: ccPtr,
		Risk:       cc * (1.0 - cov/100.0),
	}

	if m.PatchStatementsValid > 0 {
		p := 100.0 * float64(m.PatchStatementsCovered) / float64(m.PatchStatementsValid)
		h.PatchCoverage = &p
	} else if m.PatchLinesValid > 0 {
		p := 100.0 * float64(m.PatchLinesCovered) / float64(m.PatchLinesValid)
		h.PatchCoverage = &p
	}

	return h, true
}

// falls back to line coverage when the parser found no statements
func methodCoverage(m *model.MethodMetrics) (float64, bool) {
	if m.StatementsValid > 0 {
		return 100.0 * float64(m.StatementsCovered) / float64(m.StatementsValid), true
	}
	if m.LinesValid > 0 {
		return 100.0 * float64(m.LinesCovered) / float64(m.LinesValid), true
	}
	return 0, false
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
