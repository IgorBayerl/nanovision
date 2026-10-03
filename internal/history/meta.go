// Package history turns a finished pipeline tree into a stored run, and a
// stored run back into a tree the existing reporters can render.
package history

import (
	"encoding/json"
	"fmt"

	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/filtering"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/store/blob"
)

// Meta is the JSON a run keeps next to its manifest. Derived values such as
// percentages and statuses are not stored; the thresholds are, so an old run
// shows the rules of its time.
type Meta struct {
	Title string `json:"title,omitempty"`
	// report patterns; per-report hits in the coverage blobs use this order
	Reports []string `json:"reports,omitempty"`
	Parsers []string `json:"parsers,omitempty"`
	// git branch or Perforce stream, for display
	Branch string `json:"branch,omitempty"`

	Files  int               `json:"files"`
	Totals blob.Counts       `json:"totals"`
	Patch  *blob.PatchCounts `json:"patch,omitempty"`

	// the changed lines of the run, when it was measured with a diff
	Diff   *blob.Hash       `json:"diff,omitempty"`
	Change *model.ChangeSet `json:"change,omitempty"`
	// the revisions the run compared
	Versions *model.Versions `json:"versions,omitempty"`

	// the base run the delta was calculated against when the run was saved
	Base *BaseRef `json:"base,omitempty"`

	Config ConfigSnapshot `json:"config"`
}

// BaseRef is the base run of a saved run, and how well it matches the base revision.
type BaseRef struct {
	model.RunRef
	Exact bool `json:"exact"`
	// revisions between the base revision and the base run, -1 when unknown
	Distance int `json:"distance"`
}

// ConfigSnapshot is the part of the configuration that decides how a run looks.
type ConfigSnapshot struct {
	FileMetrics    []config.MetricKey        `json:"fileMetrics,omitempty"`
	MethodMetrics  []config.MetricKey        `json:"methodMetrics,omitempty"`
	StatusBands    map[config.MetricKey]Band `json:"statusBands,omitempty"`
	Gate           Gate                      `json:"gate"`
	Hotspots       int                       `json:"hotspots,omitempty"`
	DefaultFilters string                    `json:"defaultFilters,omitempty"`
	// file_filters and ignore_files, as filter rules
	Filters []string `json:"filters,omitempty"`
}

type Band struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

type Gate struct {
	PatchStatementCoverage     *float64 `json:"patchStatementCoverage,omitempty"`
	MaxChangedMethodComplexity *int     `json:"maxChangedMethodComplexity,omitempty"`
}

func SnapshotConfig(cfg *config.AppConfig) ConfigSnapshot {
	snap := ConfigSnapshot{
		FileMetrics:    cfg.FileMetrics,
		MethodMetrics:  cfg.MethodMetrics,
		Gate:           Gate{PatchStatementCoverage: cfg.Review.Gate.PatchStatementCoverage, MaxChangedMethodComplexity: cfg.Review.Gate.MaxChangedMethodComplexity},
		Hotspots:       cfg.Review.Hotspots,
		DefaultFilters: cfg.DefaultFilters,
	}
	if len(cfg.StatusBands) > 0 {
		snap.StatusBands = make(map[config.MetricKey]Band, len(cfg.StatusBands))
		for k, b := range cfg.StatusBands {
			snap.StatusBands[k] = Band{Min: b.Min, Max: b.Max}
		}
	}
	snap.Filters = append(snap.Filters, cfg.FileFilters...)
	for _, pattern := range cfg.IgnoreFiles {
		snap.Filters = append(snap.Filters, "-"+pattern)
	}
	return snap
}

// AppConfig rebuilds the configuration a stored run was made with, enough for
// the calculators, the annotator and the report builders.
func (m Meta) AppConfig() (*config.AppConfig, error) {
	snap := m.Config
	cfg := config.GetDefaultConfig()
	cfg.Title = m.Title
	if cfg.Title == "" {
		cfg.Title = "Coverage Report"
	}
	cfg.FileMetrics = snap.FileMetrics
	if len(cfg.FileMetrics) == 0 {
		cfg.FileMetrics = config.DefaultFileMetrics
	}
	cfg.MethodMetrics = snap.MethodMetrics
	if len(cfg.MethodMetrics) == 0 {
		cfg.MethodMetrics = config.DefaultMethodMetrics
	}
	cfg.ActiveFileMetrics = make(map[config.MetricKey]bool, len(cfg.FileMetrics))
	for _, k := range cfg.FileMetrics {
		cfg.ActiveFileMetrics[k] = true
	}
	cfg.ActiveMethodMetrics = make(map[config.MetricKey]bool, len(cfg.MethodMetrics))
	for _, k := range cfg.MethodMetrics {
		cfg.ActiveMethodMetrics[k] = true
	}
	cfg.StatusBands = make(config.StatusBands, len(snap.StatusBands))
	for k, b := range snap.StatusBands {
		cfg.StatusBands[k] = config.Band{Min: b.Min, Max: b.Max}
	}
	cfg.Review.Gate = config.ReviewGate{PatchStatementCoverage: snap.Gate.PatchStatementCoverage, MaxChangedMethodComplexity: snap.Gate.MaxChangedMethodComplexity}
	cfg.Review.Hotspots = snap.Hotspots
	if cfg.Review.Hotspots <= 0 {
		cfg.Review.Hotspots = 10
	}
	cfg.Review.FailOn = "never"
	cfg.DefaultFilters = snap.DefaultFilters

	filter, err := filtering.NewDefaultFilter(snap.Filters, true)
	if err != nil {
		return nil, fmt.Errorf("stored file filters: %w", err)
	}
	cfg.FileFilterInstance = filter
	return cfg, nil
}

func ParseMeta(raw json.RawMessage) (Meta, error) {
	var m Meta
	if len(raw) == 0 {
		return m, nil
	}
	if err := json.Unmarshal(raw, &m); err != nil {
		return m, fmt.Errorf("run meta: %w", err)
	}
	return m, nil
}
