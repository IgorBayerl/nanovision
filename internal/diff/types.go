package diff

// DiffData represents a complete diff, potentially containing multiple files
type DiffData struct {
	Files []FileDiff
}

// FileDiff represents changes to a single file
type FileDiff struct {
	OldPath string // Original file path, "/dev/null" for new files
	NewPath string // New file path
	Kind    string // "added", "modified", or "renamed"
	Hunks   []Hunk
}

// Hunk represents a single diff hunk with line offsets
type Hunk struct {
	// Original file line information
	OldStart int // Starting line in old file
	OldLines int // Number of lines from old file

	// New file line information
	NewStart int // Starting line in new file
	NewLines int // Number of lines in new file

	// Line offsets from NewStart for added/modified lines
	AddedLineOffsets    []int // Lines that were added (no corresponding removal)
	ModifiedLineOffsets []int // Lines that replace removed lines

}
