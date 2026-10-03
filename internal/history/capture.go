package history

import (
	"fmt"
	"sort"

	"github.com/IgorBayerl/nanovision/internal/analyzer"
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/store"
	"github.com/IgorBayerl/nanovision/internal/store/blob"
)

// Capture is a finished pipeline tree in store form. Source texts are not
// part of it: the version control system already keeps them.
type Capture struct {
	Manifest     blob.Manifest
	ManifestHash blob.Hash
	Meta         Meta
	// the manifest, coverage, analysis and diff blobs, one of each hash
	Blobs []store.Blob
}

// NewCapture encodes a tree after the ANNOTATE stage. The tree is not changed.
// branch is the git branch or Perforce stream, for display.
func NewCapture(tree *model.SummaryTree, cfg *config.AppConfig, branch string) (*Capture, error) {
	if tree == nil || tree.Root == nil {
		return nil, fmt.Errorf("capture: empty tree")
	}
	files := make(map[string]*model.FileNode)
	collectFiles(tree.Root, files)
	paths := make([]string, 0, len(files))
	for p := range files {
		paths = append(paths, p)
	}
	sort.Strings(paths)

	c := &Capture{}
	seen := make(map[blob.Hash]bool)
	add := func(data []byte) blob.Hash {
		h := blob.Sum(data)
		if !seen[h] {
			seen[h] = true
			c.Blobs = append(c.Blobs, store.Blob{Hash: h, Data: data})
		}
		return h
	}

	var diff blob.Diff
	reports := len(tree.ReportNames)
	for _, p := range paths {
		f := files[p]
		entry := blob.Entry{
			Path:     p,
			Coverage: add(blob.EncodeCoverage(coverageOf(f, reports))),
			Counts:   countsOf(f.Metrics),
		}
		if f.ContentHash != ([32]byte{}) {
			copy(entry.Content[:], f.ContentHash[:blob.HashSize])
			entry.Analysis = add(blob.EncodeAnalysis(analysisOf(f)))
		}
		if f.Diff != nil && f.Diff.Kind != model.ChangeKindNone {
			entry.Diff = diffKind(f.Diff.Kind)
			patch := patchCountsOf(f.Metrics)
			entry.Patch = &patch
			diff.Files = append(diff.Files, blob.DiffFile{
				Path:     p,
				Kind:     entry.Diff,
				Added:    lineSet(f.Diff.AddedLines),
				Modified: lineSet(f.Diff.ModifiedLines),
			})
		}
		c.Manifest.Entries = append(c.Manifest.Entries, entry)
	}

	manifestData := blob.EncodeManifest(c.Manifest)
	c.ManifestHash = add(manifestData)

	totals, patch := c.Manifest.Totals()
	c.Meta = Meta{
		Title:    cfg.Title,
		Reports:  tree.ReportNames,
		Parsers:  tree.ParserNames,
		Branch:   branch,
		Files:    len(paths),
		Totals:   totals,
		Change:   tree.Change,
		Versions: tree.Versions,
		Config:   SnapshotConfig(cfg),
	}
	if tree.Change != nil || len(diff.Files) > 0 {
		h := add(blob.EncodeDiff(diff))
		c.Meta.Diff = &h
		c.Meta.Patch = &patch
	}
	return c, nil
}

func coverageOf(f *model.FileNode, reports int) blob.Coverage {
	c := blob.Coverage{Reports: reports, Lines: make([]blob.Line, 0, len(f.Lines))}
	for n, l := range f.Lines {
		c.Lines = append(c.Lines, blob.Line{
			Number:     n,
			Hits:       l.Hits,
			ReportHits: l.ReportHits,
		})
	}
	return c
}

func analysisOf(f *model.FileNode) blob.Analysis {
	a := blob.Analysis{TotalLines: f.TotalLines}
	for _, m := range f.Methods {
		a.Functions = append(a.Functions, blob.Function{
			Name:       m.Name,
			StartLine:  m.StartLine,
			EndLine:    m.EndLine,
			Complexity: m.CyclomaticComplexity,
		})
	}
	for _, s := range f.Statements {
		a.Statements = append(a.Statements, blob.Statement{StartLine: s.StartLine, EndLine: s.EndLine, Type: s.Type})
	}
	return a
}

// analysisResult is the inverse of analysisOf, in the shape the enricher applies.
func analysisResult(a blob.Analysis) analyzer.AnalysisResult {
	var r analyzer.AnalysisResult
	for _, f := range a.Functions {
		r.Functions = append(r.Functions, analyzer.FunctionMetric{
			Name:                 f.Name,
			Position:             analyzer.Position{StartLine: f.StartLine, EndLine: f.EndLine},
			CyclomaticComplexity: f.Complexity,
		})
	}
	for _, s := range a.Statements {
		r.Statements = append(r.Statements, analyzer.StatementMetric{StartLine: s.StartLine, EndLine: s.EndLine, Type: s.Type})
	}
	return r
}

func countsOf(m model.CoverageMetrics) blob.Counts {
	return blob.Counts{
		TotalLines:                   m.TotalLines,
		LinesCovered:                 m.LinesCovered,
		LinesValid:                   m.LinesValid,
		StatementsCovered:            m.StatementsCovered,
		StatementsValid:              m.StatementsValid,
		MethodsHit:                   m.MethodsHit,
		MethodsFullyCovered:          m.MethodsFullyCovered,
		MethodsValid:                 m.MethodsValid,
		StatementMethodsHit:          m.StatementMethodsHit,
		StatementMethodsFullyCovered: m.StatementMethodsFullyCovered,
		StatementMethodsValid:        m.StatementMethodsValid,
		MaxComplexity:                m.MaxCyclomaticComplexity,
	}
}

func patchCountsOf(m model.CoverageMetrics) blob.PatchCounts {
	return blob.PatchCounts{
		LinesTotal:            m.PatchLinesTotal,
		LinesCovered:          m.PatchLinesCovered,
		LinesValid:            m.PatchLinesValid,
		StatementsCovered:     m.PatchStatementsCovered,
		StatementsValid:       m.PatchStatementsValid,
		MethodsHit:            m.PatchMethodsHit,
		MethodsValid:          m.PatchMethodsValid,
		StatementMethodsHit:   m.PatchStatementMethodsHit,
		StatementMethodsValid: m.PatchStatementMethodsValid,
	}
}

func diffKind(k model.ChangeKind) blob.DiffKind {
	switch k {
	case model.ChangeKindAdded:
		return blob.DiffAdded
	case model.ChangeKindModified:
		return blob.DiffModified
	default:
		return blob.DiffNone
	}
}

func lineSet(lines map[int]bool) []int {
	out := make([]int, 0, len(lines))
	for l, set := range lines {
		if set {
			out = append(out, l)
		}
	}
	sort.Ints(out)
	return out
}

func collectFiles(dir *model.DirNode, out map[string]*model.FileNode) {
	for _, f := range dir.Files {
		out[f.Path] = f
	}
	for _, sub := range dir.Subdirs {
		collectFiles(sub, out)
	}
}
