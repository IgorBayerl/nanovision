// In: internal/model/tree.go
package model

// SummaryTree is the new root of the entire analyzed coverage result.
type SummaryTree struct {
	Root        *DirNode        // The root directory node of the project.
	Metrics     CoverageMetrics // Aggregated metrics for the entire project.
	Timestamp   int64           // The timestamp of the report generation.
	SourceFiles []string        // List of original source directories provided by the user.
	ReportFiles []string        // List of report files that were parsed.
	ParserName  string          // Name of the parser(s) used.
}

// DirNode represents a directory in the file system tree.
type DirNode struct {
	Name    string               // The name of the directory (e.g., "analyzer").
	Path    string               // The full path from the project root.
	Metrics CoverageMetrics      // Aggregated metrics for this directory and all its children.
	Subdirs map[string]*DirNode  // Child directories, keyed by name.
	Files   map[string]*FileNode // Child files, keyed by name.
	Parent  *DirNode             // Reference to the parent directory.
}

// FileNode represents a single source code file in the tree.
type FileNode struct {
	Name      string              // The name of the file (e.g., "analyzer.go").
	Path      string              // The full path from the project root.
	Metrics   CoverageMetrics     // Metrics for this specific file.
	Lines     map[int]LineMetrics // Coverage data per line number.
	Methods   []MethodMetrics     // Metrics for methods/functions within this file.
	Parent    *DirNode            // Reference to the parent directory.
	SourceDir string              // The original source directory for this file.
}
