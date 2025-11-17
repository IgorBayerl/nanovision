package diff

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	// Regular expressions for parsing diff headers and hunks
	gitDiffRE    = regexp.MustCompile(`^diff --git a/(.*) b/(.*)$`)
	fileHeaderRE = regexp.MustCompile(`^(\+\+\+|---)\s+(.*)$`)
	newFileRE    = regexp.MustCompile(`^new file mode \d+$`)
	hunkHeaderRE = regexp.MustCompile(`^@@ -(\d+)(?:,(\d+))? \+(\d+)(?:,(\d+))? @@.*$`)
)

// Parse reads a unified diff file and returns a structured DiffData
func Parse(path string) (*DiffData, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening diff file: %w", err)
	}
	defer file.Close()

	data := &DiffData{}
	scanner := bufio.NewScanner(file)
	var currentFile *FileDiff
	var currentHunk *Hunk
	var pendingRemovals int

	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines
		if line == "" {
			continue
		}

		// Check for file header
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
			continue
		}

		// Handle new file mode
		if newFileRE.MatchString(line) {
			if currentFile != nil {
				currentFile.Kind = "added"
				currentFile.OldPath = "/dev/null"
			}
			continue
		}

		// Handle file headers (---, +++)
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
			continue
		}

		// Handle hunk headers
		if strings.HasPrefix(line, "@@ ") && strings.Contains(line, " @@") {
			// Check if it's a valid hunk header format first
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

			oldStart, err := strconv.Atoi(matches[1])
			if err != nil {
				return nil, fmt.Errorf("invalid old start line number: %v", err)
			}
			oldLines := 1
			if matches[2] != "" {
				oldLines, err = strconv.Atoi(matches[2])
				if err != nil {
					return nil, fmt.Errorf("invalid old line count: %v", err)
				}
			}
			newStart, err := strconv.Atoi(matches[3])
			if err != nil {
				return nil, fmt.Errorf("invalid new start line number: %v", err)
			}
			newLines := 1
			if matches[4] != "" {
				newLines, err = strconv.Atoi(matches[4])
				if err != nil {
					return nil, fmt.Errorf("invalid new line count: %v", err)
				}
			}

			currentHunk = &Hunk{
				OldStart: oldStart,
				OldLines: oldLines,
				NewStart: newStart,
				NewLines: newLines,
				content:  "",
			}
			pendingRemovals = 0
			continue
		}

		// Handle content lines
		if currentHunk != nil {
			switch {
			case strings.HasPrefix(line, "-"):
				pendingRemovals++
			case strings.HasPrefix(line, "+"):
				// Calculate offset based on position within current hunk
				currentLineOffset := 0
				for _, prevLine := range strings.Split(currentHunk.content, "\n") {
					if prevLine == "" {
						continue
					}
					if strings.HasPrefix(prevLine, " ") || strings.HasPrefix(prevLine, "+") {
						currentLineOffset++
					}
				}

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
			case strings.HasPrefix(line, "\\"):
				// This handles the "\ No newline at end of file" line. Just ignore it.
				continue
			default:
				return nil, fmt.Errorf("invalid line in hunk: %s", line)
			}
		}
	}

	// Add the last hunk and file if present
	if currentHunk != nil {
		currentFile.Hunks = append(currentFile.Hunks, *currentHunk)
	}
	if currentFile != nil {
		data.Files = append(data.Files, *currentFile)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning diff file: %w", err)
	}

	return data, nil
}
