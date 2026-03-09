package diffapply

import (
	"bytes"
	"strings"

	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/sourcegraph/go-diff/diff"
)

// Translates a raw diff byte slice into the domain model map.
func Parse(rawDiff []byte) (map[string]*model.DiffInfo, error) {
	fileDiffs, err := diff.ParseMultiFileDiff(rawDiff)
	if err != nil {
		return nil, err
	}

	result := make(map[string]*model.DiffInfo)

	for _, fd := range fileDiffs {
		name := cleanFilePath(fd.NewName)
		if name == "" {
			continue
		}

		result[name] = buildDiffInfo(fd)
	}

	return result, nil
}

// Extracts modified lines from hunks.
func buildDiffInfo(fd *diff.FileDiff) *model.DiffInfo {
	info := &model.DiffInfo{
		AddedLines:    make(map[int]bool),
		ModifiedLines: make(map[int]bool),
	}

	// Determine kind
	info.Kind = model.ChangeKindModified
	if fd.OrigName == "/dev/null" || strings.HasPrefix(fd.OrigName, "a/dev/null") {
		info.Kind = model.ChangeKindAdded
	} else {
		for _, ext := range fd.Extended {
			if strings.HasPrefix(ext, "new file mode") {
				info.Kind = model.ChangeKindAdded
				break
			}
		}
	}

	for _, hunk := range fd.Hunks {
		processHunk(hunk, info)
	}

	return info
}

// Tracks line numbers for inserts and modifications within a single hunk.
func processHunk(hunk *diff.Hunk, info *model.DiffInfo) {
	lines := bytes.Split(hunk.Body, []byte("\n"))
	currentLine := int(hunk.NewStartLine)

	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		if bytes.HasPrefix(line, []byte("+")) {
			info.AddedLines[currentLine] = true
			info.ModifiedLines[currentLine] = true
			currentLine++
			continue
		}

		if bytes.HasPrefix(line, []byte(" ")) {
			currentLine++
		}
	}
}

// Strips a/ or b/ prefixes standard in git diffs.
func cleanFilePath(path string) string {
	path = strings.TrimPrefix(path, "b/")
	path = strings.TrimPrefix(path, "a/")
	return path
}
