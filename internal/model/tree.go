package model

// SummaryTree is the new root of the entire analyzed coverage result.
type SummaryTree struct {
	Root        *DirNode        // The root directory node of the project.
	Metrics     CoverageMetrics // Aggregated metrics for the entire project.
	Timestamp   int64           // The timestamp of the report generation.
	SourceFiles []string        // List of original source directories provided by the user.
	ReportFiles []string        // List of report files that were parsed.
	ParserNames []string        // Name of the parser(s) used.
}

// DirNode represents a directory in the file system tree.
type DirNode struct {
	Name    string               `json:"name"`
	Path    string               `json:"path"`
	Metrics CoverageMetrics      `json:"metrics"`
	Subdirs map[string]*DirNode  `json:"subdirs,omitempty"`
	Files   map[string]*FileNode `json:"files,omitempty"`
	Parent  *DirNode             `json:"-"` // Ignore this field during JSON serialization to prevent cycles
}

// FileNode represents a single source code file in the tree.
type FileNode struct {
	Name       string              `json:"name"`
	Path       string              `json:"path"`
	Metrics    CoverageMetrics     `json:"metrics"`
	Lines      map[int]LineMetrics `json:"lines,omitempty"`
	Methods    []MethodMetrics     `json:"methods,omitempty"`
	Parent     *DirNode            `json:"-"`
	TotalLines int                 `json:"totalLines"`
	SourceDir  string              `json:"sourceDir"`
}
