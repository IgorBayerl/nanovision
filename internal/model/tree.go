package model

import "github.com/IgorBayerl/nanovision/internal/config"

// SummaryTree is the new root of the entire analyzed coverage result.
type SummaryTree struct {
	Root        *DirNode        // The root directory node of the project.
	Metrics     CoverageMetrics // Aggregated metrics for the entire project.
	Timestamp   int64           // The timestamp of the report generation.
	SourceFiles []string        // List of original source directories provided by the user.
	ReportFiles []string        // List of report files that were parsed.
	ParserNames []string        // Name of the parser(s) used.
	ReportNames []string        // Holds the list of reports, the index of an element needs to correspond to the index of LineMetrics.ReportHits

	// the diff files sorted into groups; nil without a diff
	Change *ChangeSet `json:",omitempty"`
	// the revisions being compared; nil when nothing is known about them
	Versions *Versions `json:",omitempty"`
	// the delta against the base run; nil when history is off or no base run was found
	Comparison *Comparison `json:",omitempty"`
}

// DirNode represents a directory in the file system tree.
type DirNode struct {
	Name     string                      `json:"name"`
	Path     string                      `json:"path"`
	Metrics  CoverageMetrics             `json:"metrics"`
	Subdirs  map[string]*DirNode         `json:"subdirs,omitempty"`
	Files    map[string]*FileNode        `json:"files,omitempty"`
	Parent   *DirNode                    `json:"-"` // Ignore this field during JSON serialization to prevent cycles
	Statuses map[config.MetricKey]string `json:"statuses,omitempty"`
}

type Statement struct {
	StartLine int    `json:"startLine"`
	EndLine   int    `json:"endLine"`
	Type      string `json:"type,omitempty"`
}

// FileNode represents a single source code file in the tree.
type FileNode struct {
	Name       string                      `json:"name"`
	Path       string                      `json:"path"`
	Metrics    CoverageMetrics             `json:"metrics"`
	Lines      map[int]LineMetrics         `json:"lines,omitempty"`
	Methods    []MethodMetrics             `json:"methods,omitempty"`
	Statements []Statement                 `json:"statements,omitempty"`
	Parent     *DirNode                    `json:"-"`
	TotalLines int                         `json:"totalLines"`
	SourceDir  string                      `json:"sourceDir"`
	Statuses   map[config.MetricKey]string `json:"statuses,omitempty"`
	Diff       *DiffInfo                   `json:"diff,omitempty"`
	// SHA-256 of the source file as the enricher read it; zero when it could not be read
	ContentHash [32]byte `json:"-"`
}
