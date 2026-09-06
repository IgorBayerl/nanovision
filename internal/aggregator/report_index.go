package aggregator

import (
	"sort"
	"strings"

	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
)

// MaxIndexedReports is the number of reports a bitmask can address. Runs with
// more reports than this get no index and stay pinned to the merged numbers.
const MaxIndexedReports = 32

// ReportBucket is a group of entities (lines, statements or methods) that share
// the same coverage requirement.
//
// Masks are report bitsets: bit i is set when report i covers that part of the
// entity. The entities count as covered for a selection S when every mask
// intersects S. "Hit"-style metrics need one mask ("any of these reports
// touched it"); "fully covered" metrics carry one mask per part, so all parts
// must be reached.
type ReportBucket struct {
	Masks []uint32 `json:"m"`
	Count int      `json:"n"`
}

// FileReportIndex maps a metric to the buckets covering one file. The total for
// a metric is the sum of all bucket counts and does not depend on the
// selection: which lines and methods exist is decided by the union of all
// reports, exactly as the per-line view on the details page already does.
type FileReportIndex map[config.MetricKey][]ReportBucket

// BuildFileReportIndex compresses a file's per-report hit data into buckets for
// every active metric that varies with the report selection.
//
// Branch coverage is absent on purpose: the model merges branch hits across
// reports without keeping which report covered which branch, so it cannot be
// recomputed for a subset. Cyclomatic complexity is absent because it does not
// depend on coverage at all.
func BuildFileReportIndex(file *model.FileNode, numReports int, active map[config.MetricKey]bool) FileReportIndex {
	if numReports <= 0 || numReports > MaxIndexedReports {
		return nil
	}

	full := fullMask(numReports)
	idx := FileReportIndex{}

	add := func(key config.MetricKey, build func(*bucketSet)) {
		if !active[key] {
			return
		}
		set := &bucketSet{}
		build(set)
		if buckets := set.buckets(); len(buckets) > 0 {
			idx[key] = buckets
		}
	}

	add(config.LineCoverage, func(s *bucketSet) { indexLines(s, file, full, false) })
	add(config.PatchLineCoverage, func(s *bucketSet) { indexLines(s, file, full, true) })
	add(config.StatementCoverage, func(s *bucketSet) { indexStatements(s, file, full, false) })
	add(config.PatchStatementCoverage, func(s *bucketSet) { indexStatements(s, file, full, true) })
	add(config.MethodsHit, func(s *bucketSet) { indexMethodLines(s, file, full, false) })
	add(config.MethodsFullyCovered, func(s *bucketSet) { indexMethodLines(s, file, full, true) })
	add(config.PatchMethodsHit, func(s *bucketSet) { indexPatchMethodLines(s, file, full) })
	add(config.StatementMethodsHit, func(s *bucketSet) { indexMethodStatements(s, file, full, false) })
	add(config.StatementMethodsFullyCovered, func(s *bucketSet) { indexMethodStatements(s, file, full, true) })
	add(config.PatchStatementMethodsHit, func(s *bucketSet) { indexPatchMethodStatements(s, file, full) })

	if len(idx) == 0 {
		return nil
	}
	return idx
}

// -----------------------------------------------------------------------------
// Entity indexers, one per metric family
// -----------------------------------------------------------------------------

func indexLines(s *bucketSet, file *model.FileNode, full uint32, patchOnly bool) {
	if patchOnly && file.Diff == nil {
		return
	}
	wholeFileIsPatch := patchOnly && file.Diff.Kind == model.ChangeKindAdded

	for ln, lm := range file.Lines {
		if lm.Hits < 0 {
			continue
		}
		if patchOnly && !wholeFileIsPatch && !isLineInPatch(ln, file.Diff) {
			continue
		}
		s.add(lineMask(lm, full))
	}
}

func indexStatements(s *bucketSet, file *model.FileNode, full uint32, patchOnly bool) {
	if patchOnly && file.Diff == nil {
		return
	}

	for _, stmt := range file.Statements {
		if patchOnly {
			if inPatch, _ := evaluateStatementPatchStatus(stmt, file); !inPatch {
				continue
			}
		}
		s.add(statementMask(file, stmt, full))
	}
}

func indexMethodLines(s *bucketSet, file *model.FileNode, full uint32, fully bool) {
	for i := range file.Methods {
		method := &file.Methods[i]
		masks := coverableLineMasks(file, method.StartLine, method.EndLine, full)
		if len(masks) == 0 {
			continue // LinesValid == 0, the method is not counted at all
		}
		if fully {
			s.add(masks...)
		} else {
			s.add(orMasks(masks))
		}
	}
}

// A method joins the patch when it owns at least one changed coverable line, or
// when the whole file is new. It is hit when any of those lines was executed.
func indexPatchMethodLines(s *bucketSet, file *model.FileNode, full uint32) {
	if file.Diff == nil {
		return
	}

	for i := range file.Methods {
		method := &file.Methods[i]
		if file.Diff.Kind == model.ChangeKindAdded {
			masks := coverableLineMasks(file, method.StartLine, method.EndLine, full)
			if len(masks) > 0 {
				s.add(orMasks(masks))
			}
			continue
		}

		var patchMasks []uint32
		for ln := method.StartLine; ln <= method.EndLine; ln++ {
			if !isLineInPatch(ln, file.Diff) {
				continue
			}
			if lm, ok := file.Lines[ln]; ok && lm.Hits >= 0 {
				patchMasks = append(patchMasks, lineMask(lm, full))
			}
		}
		if len(patchMasks) > 0 {
			s.add(orMasks(patchMasks))
		}
	}
}

func indexMethodStatements(s *bucketSet, file *model.FileNode, full uint32, fully bool) {
	for i := range file.Methods {
		masks := methodStatementMasks(file, &file.Methods[i], full)
		if len(masks) == 0 {
			continue
		}
		if fully {
			s.add(masks...)
		} else {
			s.add(orMasks(masks))
		}
	}
}

// Mirrors aggregatePatchMethodMetrics: the gate is the method being in the
// patch, but the hit itself is measured over all of its statements.
func indexPatchMethodStatements(s *bucketSet, file *model.FileNode, full uint32) {
	if file.Diff == nil {
		return
	}
	wholeFileIsPatch := file.Diff.Kind == model.ChangeKindAdded

	for i := range file.Methods {
		method := &file.Methods[i]
		if !wholeFileIsPatch && !methodHasPatchLines(file, method) {
			continue
		}
		masks := methodStatementMasks(file, method, full)
		if len(masks) == 0 {
			continue
		}
		s.add(orMasks(masks))
	}
}

// -----------------------------------------------------------------------------
// Mask helpers
// -----------------------------------------------------------------------------

func fullMask(numReports int) uint32 {
	if numReports >= MaxIndexedReports {
		return ^uint32(0)
	}
	return uint32(1)<<uint(numReports) - 1
}

// lineMask is the set of reports that executed this line. Line data merged
// before per-report hits existed has no breakdown, so an executed line falls
// back to "every report", which keeps the all-selected view exact.
func lineMask(lm model.LineMetrics, full uint32) uint32 {
	if len(lm.ReportHits) == 0 {
		if lm.Hits > 0 {
			return full
		}
		return 0
	}

	var mask uint32
	for i, hits := range lm.ReportHits {
		if i >= MaxIndexedReports {
			break
		}
		if hits > 0 {
			mask |= 1 << uint(i)
		}
	}
	return mask
}

// A statement counts as covered when any line inside it ran, so its mask is the
// union of its lines' masks.
func statementMask(file *model.FileNode, stmt model.Statement, full uint32) uint32 {
	var mask uint32
	for ln := stmt.StartLine; ln <= stmt.EndLine; ln++ {
		if lm, ok := file.Lines[ln]; ok {
			mask |= lineMask(lm, full)
		}
	}
	return mask
}

func coverableLineMasks(file *model.FileNode, startLine, endLine int, full uint32) []uint32 {
	var masks []uint32
	for ln := startLine; ln <= endLine; ln++ {
		if lm, ok := file.Lines[ln]; ok && lm.Hits >= 0 {
			masks = append(masks, lineMask(lm, full))
		}
	}
	return masks
}

func methodStatementMasks(file *model.FileNode, method *model.MethodMetrics, full uint32) []uint32 {
	var masks []uint32
	for _, stmt := range file.Statements {
		if stmt.StartLine >= method.StartLine && stmt.EndLine <= method.EndLine {
			masks = append(masks, statementMask(file, stmt, full))
		}
	}
	return masks
}

func methodHasPatchLines(file *model.FileNode, method *model.MethodMetrics) bool {
	for ln := method.StartLine; ln <= method.EndLine; ln++ {
		if !isLineInPatch(ln, file.Diff) {
			continue
		}
		if lm, ok := file.Lines[ln]; ok && lm.Hits >= 0 {
			return true
		}
	}
	return false
}

func orMasks(masks []uint32) uint32 {
	var out uint32
	for _, m := range masks {
		out |= m
	}
	return out
}

// -----------------------------------------------------------------------------
// Bucket accumulation
// -----------------------------------------------------------------------------

// bucketSet tallies entities by their normalised mask set, which is what keeps
// the index small: a file's lines usually fall into a handful of patterns.
type bucketSet struct {
	counts map[string]int
	masks  map[string][]uint32
	order  []string
}

func (s *bucketSet) add(masks ...uint32) {
	norm := normaliseMasks(masks)
	key := maskKey(norm)

	if s.counts == nil {
		s.counts = map[string]int{}
		s.masks = map[string][]uint32{}
	}
	if _, seen := s.counts[key]; !seen {
		s.masks[key] = norm
		s.order = append(s.order, key)
	}
	s.counts[key]++
}

func (s *bucketSet) buckets() []ReportBucket {
	out := make([]ReportBucket, 0, len(s.order))
	for _, key := range s.order {
		out = append(out, ReportBucket{Masks: s.masks[key], Count: s.counts[key]})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Count != out[j].Count {
			return out[i].Count > out[j].Count
		}
		return maskKey(out[i].Masks) < maskKey(out[j].Masks)
	})
	return out
}

// normaliseMasks sorts and dedupes so equivalent requirements share a bucket. A
// zero mask means a part no report reached, which no selection can satisfy, so
// the whole set collapses to the canonical "never covered" requirement.
func normaliseMasks(masks []uint32) []uint32 {
	if len(masks) == 0 {
		return []uint32{0}
	}

	sorted := make([]uint32, len(masks))
	copy(sorted, masks)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })

	if sorted[0] == 0 {
		return []uint32{0}
	}

	out := sorted[:1]
	for _, m := range sorted[1:] {
		if m != out[len(out)-1] {
			out = append(out, m)
		}
	}
	return out
}

func maskKey(masks []uint32) string {
	var sb strings.Builder
	for _, m := range masks {
		sb.WriteByte(byte(m))
		sb.WriteByte(byte(m >> 8))
		sb.WriteByte(byte(m >> 16))
		sb.WriteByte(byte(m >> 24))
	}
	return sb.String()
}
