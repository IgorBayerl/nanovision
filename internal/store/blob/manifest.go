package blob

import (
	"fmt"
	"slices"
	"strings"
)

const manifestVersion = 1

// Counts are the raw per-file numbers every metric is calculated from. They
// let a summary or a comparison run on the manifest alone.
type Counts struct {
	TotalLines                   int `json:"totalLines,omitempty"`
	LinesCovered                 int `json:"linesCovered,omitempty"`
	LinesValid                   int `json:"linesValid,omitempty"`
	StatementsCovered            int `json:"statementsCovered,omitempty"`
	StatementsValid              int `json:"statementsValid,omitempty"`
	MethodsHit                   int `json:"methodsHit,omitempty"`
	MethodsFullyCovered          int `json:"methodsFullyCovered,omitempty"`
	MethodsValid                 int `json:"methodsValid,omitempty"`
	StatementMethodsHit          int `json:"statementMethodsHit,omitempty"`
	StatementMethodsFullyCovered int `json:"statementMethodsFullyCovered,omitempty"`
	StatementMethodsValid        int `json:"statementMethodsValid,omitempty"`
	MaxComplexity                int `json:"maxComplexity,omitempty"`
}

func (c *Counts) fields() []*int {
	return []*int{
		&c.TotalLines,
		&c.LinesCovered, &c.LinesValid,
		&c.StatementsCovered, &c.StatementsValid,
		&c.MethodsHit, &c.MethodsFullyCovered, &c.MethodsValid,
		&c.StatementMethodsHit, &c.StatementMethodsFullyCovered, &c.StatementMethodsValid,
		&c.MaxComplexity,
	}
}

// Add sums o into c; the complexity keeps the maximum.
func (c *Counts) Add(o Counts) {
	maxComplexity := max(c.MaxComplexity, o.MaxComplexity)
	dst, src := c.fields(), o.fields()
	for i := range dst {
		*dst[i] += *src[i]
	}
	c.MaxComplexity = maxComplexity
}

// PatchCounts are the counts of the changed lines. Only runs with a diff have them.
type PatchCounts struct {
	LinesTotal            int `json:"linesTotal,omitempty"`
	LinesCovered          int `json:"linesCovered,omitempty"`
	LinesValid            int `json:"linesValid,omitempty"`
	StatementsCovered     int `json:"statementsCovered,omitempty"`
	StatementsValid       int `json:"statementsValid,omitempty"`
	MethodsHit            int `json:"methodsHit,omitempty"`
	MethodsValid          int `json:"methodsValid,omitempty"`
	StatementMethodsHit   int `json:"statementMethodsHit,omitempty"`
	StatementMethodsValid int `json:"statementMethodsValid,omitempty"`
}

func (p *PatchCounts) fields() []*int {
	return []*int{
		&p.LinesTotal,
		&p.LinesCovered, &p.LinesValid,
		&p.StatementsCovered, &p.StatementsValid,
		&p.MethodsHit, &p.MethodsValid,
		&p.StatementMethodsHit, &p.StatementMethodsValid,
	}
}

func (p *PatchCounts) Add(o PatchCounts) {
	dst, src := p.fields(), o.fields()
	for i := range dst {
		*dst[i] += *src[i]
	}
}

// DiffKind is how the diff of a run changed a file.
type DiffKind byte

const (
	DiffNone     DiffKind = 0
	DiffAdded    DiffKind = 1
	DiffModified DiffKind = 2
)

func (k DiffKind) String() string {
	switch k {
	case DiffAdded:
		return "added"
	case DiffModified:
		return "modified"
	default:
		return ""
	}
}

func (k DiffKind) MarshalText() ([]byte, error) { return []byte(k.String()), nil }

func (k *DiffKind) UnmarshalText(text []byte) error {
	switch string(text) {
	case "added":
		*k = DiffAdded
	case "modified":
		*k = DiffModified
	case "":
		*k = DiffNone
	default:
		return fmt.Errorf("unknown diff kind %q", text)
	}
	return nil
}

// Entry is one file of a run. A zero hash means the run has no such data.
type Entry struct {
	Path string `json:"path"`
	// the hash of the source text, which the store does not keep
	Content  Hash         `json:"content"`
	Analysis Hash         `json:"analysis"`
	Coverage Hash         `json:"coverage"`
	Diff     DiffKind     `json:"diff,omitempty"`
	Counts   Counts       `json:"counts"`
	Patch    *PatchCounts `json:"patch,omitempty"`
}

// Manifest lists the files of a run, sorted by path.
type Manifest struct {
	Entries []Entry `json:"entries"`
}

// Find returns the entry of path, or nil.
func (m Manifest) Find(path string) *Entry {
	i, ok := slices.BinarySearchFunc(m.Entries, path, func(e Entry, p string) int {
		return strings.Compare(e.Path, p)
	})
	if !ok {
		return nil
	}
	return &m.Entries[i]
}

// Totals sums the counts of every entry.
func (m Manifest) Totals() (Counts, PatchCounts) {
	var c Counts
	var p PatchCounts
	for _, e := range m.Entries {
		c.Add(e.Counts)
		if e.Patch != nil {
			p.Add(*e.Patch)
		}
	}
	return c, p
}

const (
	entryHasContent = 1 << iota
	entryHasAnalysis
	entryHasCoverage
	entryHasPatch
	entryDiffShift = 4 // two bits of DiffKind
	entryKnownBits = entryHasContent | entryHasAnalysis | entryHasCoverage | entryHasPatch | 3<<entryDiffShift
)

// EncodeManifest sorts the entries by path and drops repeated paths.
func EncodeManifest(m Manifest) []byte {
	entries := sortedEntries(m.Entries)
	e := newEncoder(KindManifest, manifestVersion, 80*len(entries)+4)
	e.count(len(entries))
	prev := ""
	for i := range entries {
		encodeEntry(e, prev, &entries[i])
		prev = entries[i].Path
	}
	return e.buf
}

func DecodeManifest(data []byte) (Manifest, error) {
	d, _, err := newDecoder(data, KindManifest, manifestVersion)
	if err != nil {
		return Manifest{}, err
	}
	entries, err := decodeEntries(d)
	if err != nil {
		return Manifest{}, err
	}
	return Manifest{Entries: entries}, d.done()
}

func sortedEntries(in []Entry) []Entry {
	entries := slices.Clone(in)
	slices.SortStableFunc(entries, func(a, b Entry) int { return strings.Compare(a.Path, b.Path) })
	return slices.CompactFunc(entries, func(a, b Entry) bool { return a.Path == b.Path })
}

func encodeEntry(e *encoder, prev string, en *Entry) {
	e.path(prev, en.Path)

	flags := uint64(en.Diff&3) << entryDiffShift
	if !en.Content.IsZero() {
		flags |= entryHasContent
	}
	if !en.Analysis.IsZero() {
		flags |= entryHasAnalysis
	}
	if !en.Coverage.IsZero() {
		flags |= entryHasCoverage
	}
	if en.Patch != nil {
		flags |= entryHasPatch
	}
	e.uvarint(flags)

	if flags&entryHasContent != 0 {
		e.hash(en.Content)
	}
	if flags&entryHasAnalysis != 0 {
		e.hash(en.Analysis)
	}
	if flags&entryHasCoverage != 0 {
		e.hash(en.Coverage)
	}
	counts := en.Counts
	for _, v := range counts.fields() {
		e.count(*v)
	}
	if en.Patch != nil {
		patch := *en.Patch
		for _, v := range patch.fields() {
			e.count(*v)
		}
	}
}

func decodeEntries(d *decoder) ([]Entry, error) {
	n, err := d.length(3)
	if err != nil {
		return nil, err
	}
	entries := make([]Entry, n)
	prev := ""
	for i := range n {
		if err := decodeEntry(d, prev, i == 0, &entries[i]); err != nil {
			return nil, err
		}
		prev = entries[i].Path
	}
	return entries, nil
}

func decodeEntry(d *decoder, prev string, first bool, en *Entry) error {
	var err error
	if en.Path, err = d.path(prev, first); err != nil {
		return err
	}
	flags, err := d.uvarint()
	if err != nil {
		return err
	}
	if flags&^entryKnownBits != 0 {
		return corrupt("entry %q has unknown flags %#x", en.Path, flags)
	}
	en.Diff = DiffKind(flags >> entryDiffShift & 3)
	if en.Diff > DiffModified {
		return corrupt("entry %q has diff kind %d", en.Path, en.Diff)
	}

	for _, h := range []struct {
		bit  uint64
		dest *Hash
	}{{entryHasContent, &en.Content}, {entryHasAnalysis, &en.Analysis}, {entryHasCoverage, &en.Coverage}} {
		if flags&h.bit == 0 {
			continue
		}
		if *h.dest, err = d.hash(); err != nil {
			return err
		}
		if h.dest.IsZero() {
			return corrupt("entry %q stores a zero hash", en.Path)
		}
	}

	for _, v := range en.Counts.fields() {
		if *v, err = d.count(); err != nil {
			return err
		}
	}
	if flags&entryHasPatch != 0 {
		en.Patch = &PatchCounts{}
		for _, v := range en.Patch.fields() {
			if *v, err = d.count(); err != nil {
				return err
			}
		}
	}
	return nil
}
