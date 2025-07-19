package gocover

// GoCoverProfileBlock represents a single parsed line from a Go coverage profile.
// The format is: path/to/file.go:startLine.startCol,endLine.endCol numStatements hitCount
type GoCoverProfileBlock struct {
	FileName      string
	StartLine     int
	StartCol      int
	EndLine       int
	EndCol        int
	NumStatements int
	HitCount      int
}
