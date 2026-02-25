package status

import (
	"math"

	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/utils"
)

// Annotate traverses the entire `model.SummaryTree` and attaches a risk status
// to every file and directory node based on their calculated coverage metrics.
// This function is the main entry point for the status annotation pipeline stage.
//
// What it does:
//
// It performs a recursive walk of the entire project tree. For each node (both
// files and directories), it performs the following steps for each supported metric
// (e.g., line coverage, branch coverage):
//  1. It calculates the metric's percentage value (e.g., 85.5% line coverage).
//  2. It calls `Classify` with this percentage and the corresponding thresholds
//     loaded from the user's configuration (`status_bands`).
//  3. If `Classify` indicates a status should be shown (`show: true`), it adds the
//     resulting `RiskLevel` ("danger", "warning", or "safe") to the node's
//     `Statuses` map.
//
// Contribution to Reports:
//
// This annotation step is what enables consistent risk visualization across all
// types of reports. By pre-calculating and embedding the status directly into the
// data model, the report generators themselves are kept simple.
//   - The HTML report reads the `Statuses` map to display the correct risk icon
//     (e.g., red circle, yellow triangle, blue shield) next to each metric, both
//     in the summary cards and in the file explorer tree, without needing to
//     contain any threshold logic itself.
//   - A text-based report could use this information to highlight files that are
//     in a 'danger' state.
//   - Future reporters (e.g., a "Risk Hotspots" report) can easily filter for
//     nodes with a specific status without re-implementing any logic.
//
// The `caps` parameter ensures that statuses are only applied for metrics that
// are actually supported by the input coverage formats (e.g., it prevents a
// "branch_coverage" status from being added if the report came from Go's native
// coverage tool, which lacks branch data).
func Annotate(tree *model.SummaryTree, bands config.StatusBands, caps Capabilities) {
	if tree == nil || tree.Root == nil {
		return
	}
	walk(tree.Root, func(dir *model.DirNode) {
		annotateNode(dir, bands, caps)
		for _, file := range dir.Files {
			annotateFile(file, bands, caps)
		}
	})
}

// annotateNode is a helper that calculates and attaches statuses to a single directory node.
func annotateNode(node *model.DirNode, bands config.StatusBands, caps Capabilities) {
	node.Statuses = make(map[config.MetricKey]string)
	metrics := node.Metrics

	// Line Coverage
	pct := utils.CalculatePercentage(metrics.LinesCovered, metrics.LinesValid, 2)
	if math.IsNaN(pct) {
		pct = 0.0
	}
	if lvl, show := Classify(pct, bandPtr(bands, LineCoverage)); show {
		node.Statuses[LineCoverage] = string(lvl)
	}

	// Branch Coverage (only if applicable)
	if caps.HasBranchCoverage && metrics.BranchesValid > 0 {
		pct := utils.CalculatePercentage(metrics.BranchesCovered, metrics.BranchesValid, 2)
		if lvl, show := Classify(pct, bandPtr(bands, BranchCoverage)); show {
			node.Statuses[BranchCoverage] = string(lvl)
		}
	}

	// Method Coverage (only if applicable)
	if caps.HasMethodCoverage && metrics.MethodsValid > 0 {
		pct := utils.CalculatePercentage(metrics.MethodsCovered, metrics.MethodsValid, 2)
		if lvl, show := Classify(pct, bandPtr(bands, MethodsCovered)); show {
			node.Statuses[MethodsCovered] = string(lvl)
		}

		pct = utils.CalculatePercentage(metrics.MethodsFullyCovered, metrics.MethodsValid, 2)
		if lvl, show := Classify(pct, bandPtr(bands, MethodsFullyCovered)); show {
			node.Statuses[MethodsFullyCovered] = string(lvl)
		}
	}

	// Patch Line Coverage
	// We check if PatchLinesTotal (or Valid) > 0 to ensure this file/folder was actually involved in the diff.
	if metrics.PatchLinesValid > 0 {
		pct := utils.CalculatePercentage(metrics.PatchLinesCovered, metrics.PatchLinesValid, 2)
		// Use the constant PatchLineCoverage which should resolve to "patch_line_coverage"
		if lvl, show := Classify(pct, bandPtr(bands, PatchLineCoverage)); show {
			node.Statuses[PatchLineCoverage] = string(lvl)
		}
	}

	// Patch Method Coverage
	if metrics.PatchMethodsValid > 0 {
		pct := utils.CalculatePercentage(metrics.PatchMethodsCovered, metrics.PatchMethodsValid, 2)
		// We check config for "patch_methods_covered" band
		if lvl, show := Classify(pct, bandPtr(bands, PatchMethodsCovered)); show {
			node.Statuses[PatchMethodsCovered] = string(lvl)
		}
	}
}

// annotateFile is a helper that calculates and attaches statuses to a single file node.
func annotateFile(node *model.FileNode, bands config.StatusBands, caps Capabilities) {
	node.Statuses = make(map[config.MetricKey]string)
	metrics := node.Metrics

	// Line Coverage
	pct := utils.CalculatePercentage(metrics.LinesCovered, metrics.LinesValid, 2)
	if lvl, show := Classify(pct, bandPtr(bands, LineCoverage)); show {
		node.Statuses[LineCoverage] = string(lvl)
	}

	// Branch Coverage (only if applicable)
	if caps.HasBranchCoverage && metrics.BranchesValid > 0 {
		pct := utils.CalculatePercentage(metrics.BranchesCovered, metrics.BranchesValid, 2)
		if lvl, show := Classify(pct, bandPtr(bands, BranchCoverage)); show {
			node.Statuses[BranchCoverage] = string(lvl)
		}
	}

	// Method Coverage (only if applicable)
	if caps.HasMethodCoverage && metrics.MethodsValid > 0 {
		pct := utils.CalculatePercentage(metrics.MethodsCovered, metrics.MethodsValid, 2)
		if lvl, show := Classify(pct, bandPtr(bands, MethodsCovered)); show {
			node.Statuses[MethodsCovered] = string(lvl)
		}

		pct = utils.CalculatePercentage(metrics.MethodsFullyCovered, metrics.MethodsValid, 2)
		if lvl, show := Classify(pct, bandPtr(bands, MethodsFullyCovered)); show {
			node.Statuses[MethodsFullyCovered] = string(lvl)
		}
	}

	// Patch Line Coverage
	if metrics.PatchLinesValid > 0 {
		pct := utils.CalculatePercentage(metrics.PatchLinesCovered, metrics.PatchLinesValid, 2)
		if lvl, show := Classify(pct, bandPtr(bands, PatchLineCoverage)); show {
			node.Statuses[PatchLineCoverage] = string(lvl)
		}
	}

	// Patch Method Coverage
	if metrics.PatchMethodsValid > 0 {
		pct := utils.CalculatePercentage(metrics.PatchMethodsCovered, metrics.PatchMethodsValid, 2)
		if lvl, show := Classify(pct, bandPtr(bands, PatchMethodsCovered)); show {
			node.Statuses[PatchMethodsCovered] = string(lvl)
		}
	}
}

// bandPtr is a small helper to get a pointer to a band from the map, which is
// needed for the `Classify` function's nil check.
func bandPtr(bands config.StatusBands, k MetricKey) *config.Band {
	if b, ok := bands[k]; ok {
		return &b
	}
	return nil
}

// walk performs a recursive traversal of the directory tree.
func walk(n *model.DirNode, fn func(*model.DirNode)) {
	for _, c := range n.Subdirs {
		walk(c, fn)
	}
	fn(n)
}
