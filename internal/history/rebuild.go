package history

import (
	"context"
	"fmt"
	"runtime"
	"sync"

	"github.com/IgorBayerl/nanovision/internal/aggregator"
	"github.com/IgorBayerl/nanovision/internal/calculator"
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/enricher"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/status"
	"github.com/IgorBayerl/nanovision/internal/status/evaluators"
	"github.com/IgorBayerl/nanovision/internal/store"
	"github.com/IgorBayerl/nanovision/internal/store/blob"
	"github.com/IgorBayerl/nanovision/internal/tree"
)

// Rebuild makes the tree of a stored run and evaluates it with the run's own
// configuration, so the reporters render it as they rendered the run. With
// only set, the tree holds just the files it accepts.
func Rebuild(ctx context.Context, src *store.Store, meta Meta, m blob.Manifest, only func(path string) bool) (*model.SummaryTree, *config.AppConfig, error) {
	cfg, err := meta.AppConfig()
	if err != nil {
		return nil, nil, err
	}

	entries := m.Entries
	if only != nil {
		entries = nil
		for _, e := range m.Entries {
			if only(e.Path) {
				entries = append(entries, e)
			}
		}
	}

	var hashes []blob.Hash
	for _, e := range entries {
		hashes = append(hashes, e.Coverage)
		if !e.Analysis.IsZero() {
			hashes = append(hashes, e.Analysis)
		}
	}
	if meta.Diff != nil {
		hashes = append(hashes, *meta.Diff)
	}
	blobs, err := src.Blobs(ctx, hashes)
	if err != nil {
		return nil, nil, err
	}

	var diff blob.Diff
	if meta.Diff != nil {
		data, ok := blobs[*meta.Diff]
		if !ok {
			return nil, nil, fmt.Errorf("diff blob %s is missing", meta.Diff)
		}
		if diff, err = blob.DecodeDiff(data); err != nil {
			return nil, nil, err
		}
	}

	nodes := make([]*model.FileNode, len(entries))
	errs := make([]error, len(entries))
	var wg sync.WaitGroup
	next := make(chan int)
	for range runtime.GOMAXPROCS(0) {
		wg.Go(func() {
			for i := range next {
				nodes[i], errs[i] = rebuildFile(entries[i], blobs, diff.Find(entries[i].Path), len(meta.Reports))
			}
		})
	}
	for i := range entries {
		next <- i
	}
	close(next)
	wg.Wait()
	for _, err := range errs {
		if err != nil {
			return nil, nil, err
		}
	}

	t := &model.SummaryTree{
		Root: &model.DirNode{
			Name:    "Root",
			Path:    ".",
			Subdirs: make(map[string]*model.DirNode),
			Files:   make(map[string]*model.FileNode),
		},
		ParserNames: meta.Parsers,
		ReportNames: meta.Reports,
		Change:      meta.Change,
		Versions:    meta.Versions,
	}
	for _, n := range nodes {
		placed := tree.AddFile(t.Root, n.Path, "")
		parent := placed.Parent
		n.Name, n.Parent = placed.Name, parent
		parent.Files[n.Name] = n
	}

	Evaluate(t, cfg)
	return t, cfg, nil
}

// Evaluate runs the stages after ENRICH and DIFF: aggregate, calculate, annotate.
func Evaluate(t *model.SummaryTree, cfg *config.AppConfig) {
	aggregator.AggregateMetricsAfterEnrichment(t)
	calculator.CalculateTree(t, cfg.ActiveFileMetrics, cfg.ActiveMethodMetrics)
	status.Annotate(t, cfg, status.DeriveCapabilities(t), evaluators.Registry)
}

// rebuildFile repeats the BUILD, ENRICH and DIFF stages for one file.
func rebuildFile(e blob.Entry, blobs map[blob.Hash][]byte, diff *blob.DiffFile, reports int) (*model.FileNode, error) {
	data, ok := blobs[e.Coverage]
	if !ok {
		return nil, fmt.Errorf("%s: coverage blob %s is missing", e.Path, e.Coverage)
	}
	cov, err := blob.DecodeCoverage(data)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", e.Path, err)
	}

	f := &model.FileNode{Path: e.Path, Lines: make(map[int]model.LineMetrics, len(cov.Lines))}
	for _, l := range cov.Lines {
		hits := l.ReportHits
		if len(hits) < reports {
			hits = append(hits, make([]int, reports-len(hits))...)
		}
		f.Lines[l.Number] = model.LineMetrics{Hits: l.Hits, ReportHits: hits}
	}

	var analysis *blob.Analysis
	if !e.Analysis.IsZero() {
		data, ok := blobs[e.Analysis]
		if !ok {
			return nil, fmt.Errorf("%s: analysis blob %s is missing", e.Path, e.Analysis)
		}
		a, err := blob.DecodeAnalysis(data)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", e.Path, err)
		}
		analysis = &a
		f.TotalLines = a.TotalLines
	}
	f.Metrics = tree.CalculateFileMetrics(f)
	if analysis != nil {
		enricher.ApplyAnalysis(f, analysisResult(*analysis))
	}

	if diff != nil && diff.Kind != blob.DiffNone {
		f.Diff = &model.DiffInfo{
			Kind:          model.ChangeKindModified,
			AddedLines:    make(map[int]bool, len(diff.Added)),
			ModifiedLines: make(map[int]bool, len(diff.Modified)),
		}
		if diff.Kind == blob.DiffAdded {
			f.Diff.Kind = model.ChangeKindAdded
		}
		for _, l := range diff.Added {
			f.Diff.AddedLines[l] = true
		}
		for _, l := range diff.Modified {
			f.Diff.ModifiedLines[l] = true
		}
	}
	return f, nil
}
