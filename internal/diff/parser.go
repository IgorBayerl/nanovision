package diff

import (
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"

	"github.com/IgorBayerl/nanovision/internal/filereader"
)

var (
	gitDiffRE    = regexp.MustCompile(`^diff --git a/(.*) b/(.*)$`)
	fileHeaderRE = regexp.MustCompile(`^(\+\+\+|---)\s+(.*)$`)
	newFileRE    = regexp.MustCompile(`^new file mode \d+$`)
	hunkHeaderRE = regexp.MustCompile(`^@@ -(\d+)(?:,(\d+))? \+(\d+)(?:,(\d+))? @@.*$`)
	ansiColorRE  = regexp.MustCompile(`\x1b\[[0-9;]*m`)
)

// Parse reads a unified diff file and returns a structured representation.
// It handles standard git diffs, file additions, modifications, and hunks.
// It relies on internal/filereader to handle character encoding (UTF-16/BOM).
func Parse(path string, logger *slog.Logger) (*DiffData, error) {
	if logger == nil {
		logger = slog.Default()
	}

	lines, err := filereader.ReadLinesInFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading diff file: %w", err)
	}

	logger.Debug("Starting diff parse", "path", path, "lines", len(lines))

	p := newParser(logger)
	return p.parseLines(lines)
}

// parser encapsulates the state required during the parsing process.
type parser struct {
	logger *slog.Logger
	data   *DiffData

	// State variables
	currFile        *FileDiff
	currHunk        *Hunk
	pendingRemovals int
	isGitDiff       bool // Tracks if the current file started with 'diff --git'
}

func newParser(logger *slog.Logger) *parser {
	return &parser{
		logger: logger,
		data:   &DiffData{},
	}
}

func (p *parser) parseLines(lines []string) (*DiffData, error) {
	for _, line := range lines {
		line = p.cleanLine(line)
		if line == "" {
			continue
		}

		// Try to match headers first
		if p.handleGitDiffHeader(line) {
			continue
		}
		if p.handleNewFileMarker(line) {
			continue
		}
		if p.handleFileHeader(line) {
			continue
		}
		if p.handleHunkHeader(line) {
			continue
		}

		// If inside a hunk, handle content lines (+, -, space)
		p.handleContentLine(line)
	}

	// Ensure the last file/hunk is saved
	p.finalizeHunk()
	p.finalizeFile()

	p.logger.Debug("Diff parse completed", "files_count", len(p.data.Files))
	return p.data, nil
}

// cleanLine handles whitespace and ANSI color stripping.
func (p *parser) cleanLine(line string) string {
	if strings.Contains(line, "\x1b") {
		line = ansiColorRE.ReplaceAllString(line, "")
	}
	return strings.TrimSuffix(line, "\r")
}

// handleGitDiffHeader handles "diff --git a/... b/..."
func (p *parser) handleGitDiffHeader(line string) bool {
	matches := gitDiffRE.FindStringSubmatch(line)
	if matches == nil {
		return false
	}

	p.finalizeHunk()
	p.finalizeFile()

	p.currFile = &FileDiff{
		OldPath: strings.TrimPrefix(matches[1], "a/"),
		NewPath: strings.TrimPrefix(matches[2], "b/"),
		Kind:    "modified",
	}
	p.isGitDiff = true
	p.logger.Debug("Diff found file", "old", p.currFile.OldPath, "new", p.currFile.NewPath)
	return true
}

// handleNewFileMarker handles "new file mode 123456"
func (p *parser) handleNewFileMarker(line string) bool {
	if !newFileRE.MatchString(line) {
		return false
	}
	if p.currFile != nil {
		p.currFile.Kind = "added"
		p.currFile.OldPath = "/dev/null"
		p.logger.Debug("Diff file marked as added", "file", p.currFile.NewPath)
	}
	return true
}

// handleFileHeader handles "--- a/foo" or "+++ b/foo"
func (p *parser) handleFileHeader(line string) bool {
	matches := fileHeaderRE.FindStringSubmatch(line)
	if matches == nil {
		// Log warning if it looks like a header but failed regex
		if strings.HasPrefix(line, "---") || strings.HasPrefix(line, "+++") {
			p.logger.Warn("Line looks like file header but failed regex", "line", line)
		}
		return false
	}

	prefix, rawPath := matches[1], matches[2]
	cleanPath := parsePathFromHeader(rawPath)

	if prefix == "---" {
		// Check if we need to close the previous file implicitly
		// (Common in non-git diffs where "diff --git" header is missing)
		if p.currFile != nil && (len(p.currFile.Hunks) > 0 || (!p.isGitDiff && p.currFile.NewPath != "")) {
			p.finalizeHunk()
			p.finalizeFile()
		}

		if p.currFile == nil {
			p.currFile = &FileDiff{Kind: "modified"}
			p.isGitDiff = false
		}

		if cleanPath == "/dev/null" {
			p.currFile.Kind = "added"
		}
		p.currFile.OldPath = strings.TrimPrefix(cleanPath, "a/")

	} else { // prefix == "+++"
		if p.currFile == nil {
			p.currFile = &FileDiff{Kind: "modified"}
		}
		p.currFile.NewPath = strings.TrimPrefix(cleanPath, "b/")
	}

	p.logger.Debug("Parsed diff header", "prefix", prefix, "path", cleanPath)
	return true
}

// handleHunkHeader handles "@@ -1,2 +3,4 @@"
func (p *parser) handleHunkHeader(line string) bool {
	if !strings.HasPrefix(line, "@@ ") {
		return false
	}

	matches := hunkHeaderRE.FindStringSubmatch(line)
	if matches == nil {
		p.logger.Warn("Malformed hunk header found", "line", line)
		return true // Return true to consume line even if invalid
	}

	if p.currFile == nil {
		p.logger.Error("Hunk header found without active file", "line", line)
		return true
	}

	p.finalizeHunk()

	oldStart, _ := strconv.Atoi(matches[1])
	oldLines := 1
	if matches[2] != "" {
		oldLines, _ = strconv.Atoi(matches[2])
	}

	newStart, _ := strconv.Atoi(matches[3])
	newLines := 1
	if matches[4] != "" {
		newLines, _ = strconv.Atoi(matches[4])
	}

	p.currHunk = &Hunk{
		OldStart: oldStart,
		OldLines: oldLines,
		NewStart: newStart,
		NewLines: newLines,
	}
	p.pendingRemovals = 0

	p.logger.Debug("Diff found hunk", "file", p.currFile.NewPath, "start", newStart)
	return true
}

// handleContentLine handles actual code lines (+, -, space)
func (p *parser) handleContentLine(line string) {
	if p.currHunk == nil {
		// Just a debug log if we encounter content outside a hunk (rare/garbage)
		if strings.TrimSpace(line) != "" {
			p.logger.Debug("Ignoring line outside hunk", "line", line)
		}
		return
	}

	switch {
	case strings.HasPrefix(line, "-"):
		p.pendingRemovals++

	case strings.HasPrefix(line, "+"):
		// Calculate offset relative to the new file version
		// We count lines in our accumulated content that are meant for the new file (starts with ' ' or '+')
		currentLineOffset := 0
		for _, prev := range strings.Split(p.currHunk.content, "\n") {
			if prev == "" {
				continue
			}
			if strings.HasPrefix(prev, " ") || strings.HasPrefix(prev, "+") {
				currentLineOffset++
			}
		}

		if p.pendingRemovals > 0 {
			p.pendingRemovals--
		}

		p.currHunk.AddedLineOffsets = append(p.currHunk.AddedLineOffsets, currentLineOffset)
		p.currHunk.content += line + "\n"

	case strings.HasPrefix(line, " "):
		p.currHunk.content += line + "\n"
		p.pendingRemovals = 0

	case strings.HasPrefix(line, `\`):
		// Skip "\ No newline at end of file"
		return

	default:
		// Only log if it's not just an empty line
		if strings.TrimSpace(line) != "" {
			p.logger.Debug("Ignoring unusual line inside hunk", "line", line)
		}
	}
}

// finalizeHunk appends the current hunk to the current file and resets state.
func (p *parser) finalizeHunk() {
	if p.currHunk != nil && p.currFile != nil {
		p.currFile.Hunks = append(p.currFile.Hunks, *p.currHunk)
	}
	p.currHunk = nil
	p.pendingRemovals = 0
}

// finalizeFile appends the current file to the data and resets state.
func (p *parser) finalizeFile() {
	if p.currFile != nil {
		p.data.Files = append(p.data.Files, *p.currFile)
	}
	p.currFile = nil
	p.isGitDiff = false
}

// parsePathFromHeader extracts clean paths from lines like "--- a/foo.txt\t2023..."
// It handles Perforce/Unified formats (tab separated) and #rev suffixes.
func parsePathFromHeader(rawPath string) string {
	parts := strings.Split(rawPath, "\t")
	path := strings.TrimSpace(parts[0])

	// Strip Perforce version suffix (e.g., file.go#3)
	if strings.Contains(path, "#") {
		if idx := strings.LastIndex(path, "#"); idx > 0 {
			suffix := path[idx+1:]
			if _, err := strconv.Atoi(suffix); err == nil {
				return path[:idx]
			}
		}
	}
	return path
}
