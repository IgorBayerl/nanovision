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
	gitDiffRE = regexp.MustCompile(`^diff --git a/(.*) b/(.*)$`)
	// fileHeaderRE captures the prefix (--- or +++) and the rest of the line
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
func Parse(path string, logger *slog.Logger) (*DiffData, error) {
	if logger == nil {
		logger = slog.Default()
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening diff file: %w", err)
	}
	defer file.Close()

	logger.Debug("Starting diff parse", "path", path)

	reader := bufio.NewReader(file)
	data := &DiffData{}
	var currentFile *FileDiff
	var currentHunk *Hunk
	var pendingRemovals int
	var currentFileCreatedViaDiffGit bool

	for {
		lineBytes, err := reader.ReadBytes('\n')
		line := strings.TrimSuffix(string(lineBytes), "\n")
		line = strings.TrimSuffix(line, "\r")
		// Strip BOM if present (common in Windows files)
		line = strings.TrimPrefix(line, "\uFEFF")

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

			oldPath := strings.TrimPrefix(matches[1], "a/")
			newPath := strings.TrimPrefix(matches[2], "b/")

			logger.Debug("Diff found file", "raw_old", matches[1], "raw_new", matches[2])

			currentFile = &FileDiff{
				OldPath: oldPath,
				NewPath: newPath,
				Kind:    "modified",
			}
			currentHunk = nil
			currentFileCreatedViaDiffGit = true
			goto nextIteration
		}

		// Detect file additions via "new file mode" marker
		if newFileRE.MatchString(line) {
			if currentFile != nil {
				currentFile.Kind = "added"
				currentFile.OldPath = "/dev/null"
				logger.Debug("Diff file marked as added", "file", currentFile.NewPath)
			}
			goto nextIteration
		}

		// Parse "---" and "+++" file path headers
		if matches := fileHeaderRE.FindStringSubmatch(line); matches != nil {
			prefix, rawPath := matches[1], matches[2]

			// Clean the path:
			// Perforce/Unified diffs separate path and timestamp with a tab
			pathParts := strings.Split(rawPath, "\t")
			cleanPath := strings.TrimSpace(pathParts[0])

			// Strip Perforce version suffix (e.g., file.go#3) if present (common in OldPath)
			if strings.Contains(cleanPath, "#") {
				// Check if it looks like a version number #123
				if idx := strings.LastIndex(cleanPath, "#"); idx > 0 {
					// Simple heuristic: if it's #number, strip it.
					// This prevents stripping valid filenames that happen to contain #
					suffix := cleanPath[idx+1:]
					if _, err := strconv.Atoi(suffix); err == nil {
						cleanPath = cleanPath[:idx]
					}
				}
			}

			if prefix == "---" {
				// If we encounter "---", it might be the start of a new file if:
				// 1. We have a current file that has hunks (definitely finished).
				// 2. We have a current file that wasn't created by "diff --git" and already has a NewPath set (finished empty file?).
				if currentFile != nil {
					shouldClose := false
					if len(currentFile.Hunks) > 0 {
						shouldClose = true
					} else if !currentFileCreatedViaDiffGit && currentFile.NewPath != "" {
						shouldClose = true
					}

					if shouldClose {
						if currentHunk != nil {
							currentFile.Hunks = append(currentFile.Hunks, *currentHunk)
							currentHunk = nil
						}
						data.Files = append(data.Files, *currentFile)
						currentFile = nil
						currentFileCreatedViaDiffGit = false
					}
				}

				if currentFile == nil {
					currentFile = &FileDiff{Kind: "modified"}
					currentFileCreatedViaDiffGit = false
				}

				if cleanPath == "/dev/null" {
					currentFile.Kind = "added"
				}
				// Git style stripping (only if it starts with a/)
				currentFile.OldPath = strings.TrimPrefix(cleanPath, "a/")
			} else {
				// prefix == "+++"
				if currentFile == nil {
					// Should not happen in valid diffs, but handle gracefully
					currentFile = &FileDiff{Kind: "modified"}
				}
				// Git style stripping (only if it starts with b/)
				currentFile.NewPath = strings.TrimPrefix(cleanPath, "b/")
			}

			logger.Debug("Parsed diff header", "prefix", prefix, "raw", rawPath, "clean", cleanPath)
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

			pathForLog := "unknown"
			if currentFile != nil {
				pathForLog = currentFile.NewPath
			}
			logger.Debug("Diff found hunk", "file", pathForLog, "new_start", newStart, "new_lines", newLines)

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

				// Simplification: Always treat '+' lines as "Added" lines.
				// Previously, we tried to distinguish between "Modified" (replacing a -) and "Added" (pure insertion).
				// However, for coverage reporting on the *new* file, we want to highlight ALL changed lines uniformly.
				// This ensures the visual indicators in the report show up for all changes.

				// We still consume pendingRemovals to allow for future logic if needed, but we don't map to ModifiedLineOffsets.
				if pendingRemovals > 0 {
					pendingRemovals--
				}

				currentHunk.AddedLineOffsets = append(currentHunk.AddedLineOffsets, currentLineOffset)
				currentHunk.content += line + "\n"
			case strings.HasPrefix(line, " "):
				currentHunk.content += line + "\n"
				pendingRemovals = 0
			case strings.HasPrefix(line, `\`):
				// Skip "\ No newline at end of file" markers
				goto nextIteration
			default:
				logger.Warn("Ignoring invalid line inside a diff hunk", "line", line)
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
				logger.Warn("Ignoring overly long line in diff, skipping rest of hunk.", "file", fileName)

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

	logger.Debug("Diff parse completed", "files_count", len(data.Files))
	return data, nil
}
