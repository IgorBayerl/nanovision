package blob

import (
	"slices"
	"strings"
)

const diffVersion = 1

// DiffFile is the changed lines of one file of the run.
type DiffFile struct {
	Path     string   `json:"path"`
	Kind     DiffKind `json:"kind"`
	Added    []int    `json:"added,omitempty"`
	Modified []int    `json:"modified,omitempty"`
}

// Diff is the diff a run was measured with, limited to the files in its tree.
type Diff struct {
	Files []DiffFile `json:"files"`
}

// Find returns the file of path, or nil.
func (d Diff) Find(path string) *DiffFile {
	i, ok := slices.BinarySearchFunc(d.Files, path, func(f DiffFile, p string) int {
		return strings.Compare(f.Path, p)
	})
	if !ok {
		return nil
	}
	return &d.Files[i]
}

func EncodeDiff(diff Diff) []byte {
	files := slices.Clone(diff.Files)
	slices.SortStableFunc(files, func(a, b DiffFile) int { return strings.Compare(a.Path, b.Path) })
	files = slices.CompactFunc(files, func(a, b DiffFile) bool { return a.Path == b.Path })

	e := newEncoder(KindDiff, diffVersion, 16*len(files)+4)
	e.count(len(files))
	prev := ""
	for _, f := range files {
		e.path(prev, f.Path)
		e.uvarint(uint64(f.Kind & 3))
		encodeLineSet(e, f.Added)
		encodeLineSet(e, f.Modified)
		prev = f.Path
	}
	return e.buf
}

func DecodeDiff(data []byte) (Diff, error) {
	d, _, err := newDecoder(data, KindDiff, diffVersion)
	if err != nil {
		return Diff{}, err
	}
	n, err := d.length(4)
	if err != nil {
		return Diff{}, err
	}
	out := Diff{Files: make([]DiffFile, n)}
	prev := ""
	for i := range n {
		f := &out.Files[i]
		if f.Path, err = d.path(prev, i == 0); err != nil {
			return Diff{}, err
		}
		kind, err := d.uvarint()
		if err != nil {
			return Diff{}, err
		}
		if kind > uint64(DiffModified) {
			return Diff{}, corrupt("file %q has diff kind %d", f.Path, kind)
		}
		f.Kind = DiffKind(kind)
		if f.Added, err = decodeLineSet(d); err != nil {
			return Diff{}, err
		}
		if f.Modified, err = decodeLineSet(d); err != nil {
			return Diff{}, err
		}
		prev = f.Path
	}
	return out, d.done()
}

// encodeLineSet writes positive line numbers as runs of consecutive lines,
// because changed lines come in blocks.
func encodeLineSet(e *encoder, lines []int) {
	sorted := slices.Clone(lines)
	slices.Sort(sorted)
	sorted = slices.Compact(sorted)
	sorted = slices.DeleteFunc(sorted, func(l int) bool { return l < 1 || l > maxLineNumber })

	type run struct{ start, length int }
	var runs []run
	for _, l := range sorted {
		if n := len(runs); n > 0 && runs[n-1].start+runs[n-1].length == l {
			runs[n-1].length++
			continue
		}
		runs = append(runs, run{start: l, length: 1})
	}

	e.count(len(runs))
	end := 0 // one past the previous run
	for _, r := range runs {
		e.count(r.start - end)
		e.count(r.length - 1)
		end = r.start + r.length
	}
}

const maxLineNumber = 1<<31 - 1

func decodeLineSet(d *decoder) ([]int, error) {
	n, err := d.length(2)
	if err != nil {
		return nil, err
	}
	var lines []int
	end := 0
	for i := range n {
		gap, err := d.count()
		if err != nil {
			return nil, err
		}
		// runs are separated by at least one line, except the first which starts at 1 or later
		if i > 0 && gap == 0 {
			return nil, corrupt("adjacent line runs are not merged")
		}
		if i == 0 && gap == 0 {
			return nil, corrupt("line 0 in a line set")
		}
		length, err := d.count()
		if err != nil {
			return nil, err
		}
		start := end + gap
		stop := int64(start) + int64(length) + 1
		if stop-1 > maxLineNumber {
			return nil, corrupt("line %d out of range", stop-1)
		}
		if len(lines)+length+1 > 1<<24 {
			return nil, corrupt("line set too large")
		}
		for l := start; l < int(stop); l++ {
			lines = append(lines, l)
		}
		end = int(stop)
	}
	return lines, nil
}
