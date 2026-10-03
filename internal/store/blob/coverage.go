package blob

import (
	"math"
	"slices"
)

const coverageVersion = 1

// Line is the coverage of one line. Hits is -1 for a line a report marks as
// not coverable. ReportHits has one entry per report of the run.
type Line struct {
	Number     int   `json:"line"`
	Hits       int   `json:"hits"`
	ReportHits []int `json:"reportHits,omitempty"`
}

// Coverage is the coverage of one file, lines sorted by number.
type Coverage struct {
	Reports int    `json:"reports"`
	Lines   []Line `json:"lines"`
}

const (
	lineHasReportHits = 1 << 0
	lineFlagBits      = 1
	// hits above this are clamped, so the flag bits never overflow the varint
	maxHits = math.MaxInt64 >> (lineFlagBits + 1)
)

// EncodeCoverage sorts the lines by number. Two lines with the same number
// keep the first.
func EncodeCoverage(c Coverage) []byte {
	lines := slices.Clone(c.Lines)
	slices.SortStableFunc(lines, func(a, b Line) int { return a.Number - b.Number })

	e := newEncoder(KindCoverage, coverageVersion, 3*len(lines)+4)
	e.count(c.Reports)

	unique := make([]Line, 0, len(lines))
	for _, line := range lines {
		if len(unique) > 0 && line.Number == unique[len(unique)-1].Number {
			continue
		}
		if line.Number < math.MinInt32 || line.Number > math.MaxInt32 {
			continue
		}
		unique = append(unique, line)
	}
	e.count(len(unique))

	prev := 0
	for i, line := range unique {
		if i == 0 {
			e.varint(int64(line.Number))
		} else {
			e.uvarint(uint64(line.Number - prev - 1))
		}
		prev = line.Number

		hits := clampHits(line.Hits)
		reportHits := normalizeReportHits(line.ReportHits, c.Reports)
		flags := uint64(0)
		if !slices.Equal(reportHits, defaultReportHits(hits, c.Reports)) {
			flags |= lineHasReportHits
		}
		e.uvarint(uint64(hits+1)<<lineFlagBits | flags)

		if flags&lineHasReportHits != 0 {
			encodeReportHits(e, reportHits)
		}
	}
	return e.buf
}

func DecodeCoverage(data []byte) (Coverage, error) {
	d, _, err := newDecoder(data, KindCoverage, coverageVersion)
	if err != nil {
		return Coverage{}, err
	}
	reports, err := d.count()
	if err != nil {
		return Coverage{}, err
	}
	if reports > maxReports {
		return Coverage{}, corrupt("%d reports", reports)
	}
	n, err := d.length(2)
	if err != nil {
		return Coverage{}, err
	}

	c := Coverage{Reports: reports, Lines: make([]Line, n)}
	prev := 0
	for i := range n {
		var number int
		if i == 0 {
			v, err := d.varint()
			if err != nil {
				return Coverage{}, err
			}
			if v < math.MinInt32 || v > math.MaxInt32 {
				return Coverage{}, corrupt("line %d out of range", v)
			}
			number = int(v)
		} else {
			delta, err := d.count()
			if err != nil {
				return Coverage{}, err
			}
			number = prev + delta + 1
			if number > math.MaxInt32 {
				return Coverage{}, corrupt("line %d out of range", number)
			}
		}
		prev = number

		packed, err := d.uvarint()
		if err != nil {
			return Coverage{}, err
		}
		flags := packed & (1<<lineFlagBits - 1)
		if packed>>lineFlagBits > maxHits+1 {
			return Coverage{}, corrupt("line %d: hits out of range", number)
		}
		line := Line{Number: number, Hits: int(packed>>lineFlagBits) - 1}

		if flags&lineHasReportHits != 0 {
			if line.ReportHits, err = decodeReportHits(d, reports); err != nil {
				return Coverage{}, err
			}
			if slices.Equal(line.ReportHits, defaultReportHits(line.Hits, reports)) {
				return Coverage{}, corrupt("line %d stores its default report hits", number)
			}
		} else {
			line.ReportHits = defaultReportHits(line.Hits, reports)
		}
		c.Lines[i] = line
	}
	return c, d.done()
}

// the most reports one run can carry; far above anything real, it only
// bounds the decoder
const maxReports = 1 << 16

func clampHits(h int) int {
	if h < -1 {
		return -1
	}
	return min(h, maxHits)
}

// normalizeReportHits gives the slice exactly one entry per report.
func normalizeReportHits(hits []int, reports int) []int {
	if reports == 0 {
		return nil
	}
	out := make([]int, reports)
	for i := range min(len(hits), reports) {
		out[i] = clampHits(hits[i])
	}
	return out
}

// defaultReportHits is what a line without stored report hits decodes to: a
// single report owns every hit, several reports own none.
func defaultReportHits(hits, reports int) []int {
	if reports == 0 {
		return nil
	}
	out := make([]int, reports)
	if reports == 1 {
		out[0] = hits
	}
	return out
}

// only the non-zero entries are written, as (index gap, hits+1) pairs
func encodeReportHits(e *encoder, hits []int) {
	nonZero := 0
	for _, h := range hits {
		if h != 0 {
			nonZero++
		}
	}
	e.count(nonZero)
	prev := -1
	for i, h := range hits {
		if h == 0 {
			continue
		}
		e.uvarint(uint64(i - prev - 1))
		e.uvarint(uint64(h + 1))
		prev = i
	}
}

func decodeReportHits(d *decoder, reports int) ([]int, error) {
	n, err := d.length(2)
	if err != nil {
		return nil, err
	}
	if n > reports {
		return nil, corrupt("%d report hits for %d reports", n, reports)
	}
	out := make([]int, reports)
	prev := -1
	for range n {
		gap, err := d.count()
		if err != nil {
			return nil, err
		}
		idx := prev + gap + 1
		if idx >= reports {
			return nil, corrupt("report index %d of %d", idx, reports)
		}
		v, err := d.uvarint()
		if err != nil {
			return nil, err
		}
		if v == 1 || v > maxHits+1 {
			return nil, corrupt("report hit value %d", v)
		}
		out[idx] = int(v) - 1
		prev = idx
	}
	return out, nil
}
