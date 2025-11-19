package diff

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	gitDiffRE    = regexp.MustCompile(`^diff --git a/(.*) b/(.*)$`)
	fileHeaderRE = regexp.MustCompile(`^(\+\+\+|---)\s+(.*)$`)
	newFileRE    = regexp.MustCompile(`^new file mode \d+$`)
	hunkHeaderRE = regexp.MustCompile(`^@@ -(\d+)(?:,(\d+))? \+(\d+)(?:,(\d+))? @@.*$`)
)

// Parse reads a unified diff file from the given path and returns a structured
// representation of its contents. It handles standard git diff format including
// file additions, modifications, and the individual hunks within each file.
//
// The parser is resilient to malformed input - it will log warnings and continue
// processing when encountering oversized lines or unexpected content rather than
// failing immediately.
func Parse(path string) (*DiffData, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening diff file: %w", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	data := &DiffData{}
	var currentFile *FileDiff
	var currentHunk *Hunk
	var pendingRemovals int

	for {
		lineBytes, err := reader.ReadBytes('\n')
		line := strings.TrimSuffix(string(lineBytes), "\n")

		if line == "" {
			goto nextIteration
		}

		// Parse git diff header: "diff --git a/path b/path"
		if matches := gitDiffRE.FindStringSubmatch(line); matches != nil {
			if currentFile != nil {
				if currentHunk != nil {
					currentFile.Hunks = append(currentFile.Hunks, *currentHunk)
				}
				data.Files = append(data.Files, *currentFile)
			}
			currentFile = &FileDiff{
				OldPath: strings.TrimPrefix(matches[1], "a/"),
				NewPath: strings.TrimPrefix(matches[2], "b/"),
				Kind:    "modified",
			}
			currentHunk = nil
			goto nextIteration
		}

		// Detect file additions via "new file mode" marker
		if newFileRE.MatchString(line) {
			if currentFile != nil {
				currentFile.Kind = "added"
				currentFile.OldPath = "/dev/null"
			}
			goto nextIteration
		}

		// Parse "---" and "+++" file path headers
		if matches := fileHeaderRE.FindStringSubmatch(line); matches != nil {
			prefix, path := matches[1], matches[2]
			if currentFile == nil {
				currentFile = &FileDiff{Kind: "modified"}
			}
			if prefix == "---" {
				if path == "/dev/null" {
					currentFile.Kind = "added"
				}
				currentFile.OldPath = strings.TrimPrefix(path, "a/")
			} else {
				currentFile.NewPath = strings.TrimPrefix(path, "b/")
			}
			goto nextIteration
		}

		// Parse hunk header: "@@ -1,5 +1,6 @@"
		if strings.HasPrefix(line, "@@ ") && strings.Contains(line, " @@") {
			if !hunkHeaderRE.MatchString(line) {
				return nil, fmt.Errorf("malformed hunk header: %s", line)
			}
			matches := hunkHeaderRE.FindStringSubmatch(line)
			if currentFile == nil {
				return nil, fmt.Errorf("hunk header without file header")
			}
			if currentHunk != nil {
				currentFile.Hunks = append(currentFile.Hunks, *currentHunk)
			}
			oldStart, _ := strconv.Atoi(matches[1])
			oldLines, _ := strconv.Atoi(matches[2])
			if matches[2] == "" {
				oldLines = 1
			}
			newStart, _ := strconv.Atoi(matches[3])
			newLines, _ := strconv.Atoi(matches[4])
			if matches[4] == "" {
				newLines = 1
			}
			currentHunk = &Hunk{
				OldStart: oldStart,
				OldLines: oldLines,
				NewStart: newStart,
				NewLines: newLines,
				content:  "",
			}
			pendingRemovals = 0
			goto nextIteration
		}

		// Process hunk content lines (-, +, and context lines)
		if currentHunk != nil {
			switch {
			case strings.HasPrefix(line, "-"):
				pendingRemovals++
			case strings.HasPrefix(line, "+"):
				// Calculate the current line offset within the new file version
				currentLineOffset := 0
				for _, prevLine := range strings.Split(currentHunk.content, "\n") {
					if prevLine == "" {
						continue
					}
					if strings.HasPrefix(prevLine, " ") || strings.HasPrefix(prevLine, "+") {
						currentLineOffset++
					}
				}
				// Match additions with previous removals to identify modifications
				if pendingRemovals > 0 {
					currentHunk.ModifiedLineOffsets = append(currentHunk.ModifiedLineOffsets, currentLineOffset)
					pendingRemovals--
				} else {
					currentHunk.AddedLineOffsets = append(currentHunk.AddedLineOffsets, currentLineOffset)
				}
				currentHunk.content += line + "\n"
			case strings.HasPrefix(line, " "):
				currentHunk.content += line + "\n"
				pendingRemovals = 0
			case strings.HasPrefix(line, `\`):
				// Skip "\ No newline at end of file" markers
				goto nextIteration
			default:
				slog.Warn("Ignoring invalid line inside a diff hunk", "line", line)
			}
		}

	nextIteration:
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			if errors.Is(err, bufio.ErrBufferFull) {
				// Line exceeded buffer capacity - skip to next parseable line
				fileName := "unknown"
				if currentFile != nil {
					fileName = currentFile.NewPath
				}
				slog.Warn("Ignoring overly long line in diff, skipping rest of hunk.", "file", fileName)

				// Consume remainder of oversized line
				for errors.Is(err, bufio.ErrBufferFull) {
					_, err = reader.ReadBytes('\n')
				}

				// Discard corrupted hunk state
				currentHunk = nil
				pendingRemovals = 0
				continue
			}

			return nil, fmt.Errorf("reading diff file: %w", err)
		}
	}

	// Finalize any remaining hunk and file data
	if currentHunk != nil {
		currentFile.Hunks = append(currentFile.Hunks, *currentHunk)
	}
	if currentFile != nil {
		data.Files = append(data.Files, *currentFile)
	}

	return data, nil
}
