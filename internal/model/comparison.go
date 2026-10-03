package model

import "github.com/IgorBayerl/nanovision/internal/config"

// ChangeSet sorts the files of a diff into the groups the delta needs. Paths
// are relative to the project root.
type ChangeSet struct {
	// in the coverage tree; patch coverage uses them
	Measured []string `json:"measured,omitempty"`
	// removed by ignore_files or file_filters, usually tests
	Ignored []string `json:"ignored,omitempty"`
	// no coverage report has them, for example build scripts
	NotInReports []string `json:"notInReports,omitempty"`
	// deleted by the change and not removed by the filters
	Deleted []string `json:"deleted,omitempty"`
}

func (c *ChangeSet) Total() int {
	if c == nil {
		return 0
	}
	return len(c.Measured) + len(c.Ignored) + len(c.NotInReports) + len(c.Deleted)
}

// TestOnly is true for a change that touches no measured code, which usually
// means it adds or edits only tests. Only the delta shows what it is worth.
func (c *ChangeSet) TestOnly() bool {
	return c != nil && c.Total() > 0 && len(c.Measured) == 0 && len(c.Deleted) == 0
}

// RunRef names a stored run in a comparison.
type RunRef struct {
	ID       int64  `json:"id,omitempty"`
	Revision string `json:"revision,omitempty"`
	Stream   string `json:"stream,omitempty"`
	Profile  string `json:"profile,omitempty"`
	Kind     string `json:"kind,omitempty"`
	// unix milliseconds
	CreatedAt int64  `json:"createdAt,omitempty"`
	Tool      string `json:"tool,omitempty"`
}

// Versions names the two states of the code a report compares. Only the
// revision and stream of each RunRef are set.
type Versions struct {
	// the revision the change starts from; the diff is against it
	Base RunRef `json:"base"`
	// the revision the run measured
	Current RunRef `json:"current"`
	// the workspace had edits on top of Current
	LocalEdits bool `json:"localEdits,omitempty"`
	// where the changed lines come from: the diff file, or "git" or "perforce"
	Diff string `json:"diff,omitempty"`
}

// MetricDelta is one coverage metric of both runs. Percentages come from the
// counts, never from stored values.
type MetricDelta struct {
	Key            config.MetricKey `json:"key"`
	Label          string           `json:"label"`
	Unit           string           `json:"unit"` // what the counts count, e.g. "statements"
	BaseCovered    int              `json:"baseCovered"`
	BaseTotal      int              `json:"baseTotal"`
	CurrentCovered int              `json:"currentCovered"`
	CurrentTotal   int              `json:"currentTotal"`
	Base           float64          `json:"base"`
	Current        float64          `json:"current"`
	// Current minus Base, in percentage points
	Delta float64 `json:"delta"`
}

// FileDelta is the headline metric of one file in both runs.
type FileDelta struct {
	Path           string  `json:"path"`
	BaseCovered    int     `json:"baseCovered"`
	BaseTotal      int     `json:"baseTotal"`
	CurrentCovered int     `json:"currentCovered"`
	CurrentTotal   int     `json:"currentTotal"`
	Base           float64 `json:"base"`
	Current        float64 `json:"current"`
	Delta          float64 `json:"delta"`
	// "added" or "modified" from the diff, "not in diff" when the run has a
	// diff that does not touch the file, "" without a diff
	Change string `json:"change,omitempty"`
}

// PatchSummary repeats the patch coverage of the current run.
type PatchSummary struct {
	Unit    string `json:"unit"` // "statements" or "lines"
	Covered int    `json:"covered"`
	Total   int    `json:"total"`
}

// Comparison is the delta of a run against its base run.
type Comparison struct {
	Base    RunRef `json:"base"`
	Current RunRef `json:"current"`
	// the base run is at the base revision itself
	Exact bool `json:"exact"`
	// revisions between the base revision and the base run, -1 when unknown
	Distance int `json:"distance"`
	// statement coverage when both runs have statements, else line coverage
	Headline config.MetricKey `json:"headline"`
	Metrics  []MetricDelta    `json:"metrics"`
	// files whose headline changed, largest change first
	Files []FileDelta `json:"files,omitempty"`
	// for each file and folder whose coverage changed: the change of each
	// metric in percentage points, by path as in the tree
	Deltas map[string]map[config.MetricKey]float64 `json:"deltas,omitempty"`
	// files only one of the runs has
	AddedFiles   int `json:"addedFiles,omitempty"`
	RemovedFiles int `json:"removedFiles,omitempty"`
	// base files the current filters remove; they are left out of every number
	FilteredBaseFiles int           `json:"filteredBaseFiles,omitempty"`
	Change            *ChangeSet    `json:"change,omitempty"`
	Patch             *PatchSummary `json:"patch,omitempty"`
	// signs that the runs used different tests
	Warnings []string `json:"warnings,omitempty"`
}

// HeadlineDelta returns the delta of the headline metric, if the comparison has it.
func (c *Comparison) HeadlineDelta() (MetricDelta, bool) {
	if c == nil {
		return MetricDelta{}, false
	}
	for _, m := range c.Metrics {
		if m.Key == c.Headline {
			return m, true
		}
	}
	return MetricDelta{}, false
}
