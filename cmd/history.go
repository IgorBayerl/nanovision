package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/IgorBayerl/nanovision/internal/client"
	"github.com/IgorBayerl/nanovision/internal/compare"
	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/diff"
	"github.com/IgorBayerl/nanovision/internal/history"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/store"
	"github.com/IgorBayerl/nanovision/internal/vcs"
)

// runIdentity is what this run is: its revision, its base and its kind.
// It comes from the flags first and from the VCS adapter second.
type runIdentity struct {
	vcs     vcs.VCS
	current vcs.Revision
	base    vcs.Revision
	kind    store.RunKind
	clean   bool
	edits   bool
	stream  string
	author  string
	// true when the workspace is clean and at the base revision
	olderOnly bool
}

// resolveIdentity runs before the pipeline, because a diff from the VCS feeds
// patch coverage. It fails soft: without an answer from the VCS the run goes
// on with what the flags say.
func resolveIdentity(cfg *config.AppConfig, logger *slog.Logger) runIdentity {
	id := runIdentity{kind: store.RunKind(cfg.Run.Kind)}

	if cfg.VCS.Type != "" {
		v, err := vcs.Detect(cfg.VCS.Type, cfg.ProjectRoot, vcs.Options{BaseBranch: cfg.VCS.BaseBranch})
		switch {
		case errors.Is(err, vcs.ErrNotFound):
			logger.Info("No git repository or Perforce workspace found; the run has no revision")
		case err != nil:
			logger.Warn("VCS adapter is not available; continuing without it", "type", cfg.VCS.Type, "error", err)
		default:
			id.vcs = v
		}
	}

	if id.vcs != nil {
		var err error
		if id.current, err = id.vcs.Current(); err != nil {
			logger.Warn("Could not read the current revision", "vcs", id.vcs.Name(), "error", err)
		}
		if id.edits, err = id.vcs.HasLocalEdits(); err != nil {
			logger.Warn("Could not check for local edits", "vcs", id.vcs.Name(), "error", err)
			id.edits = true
		}
		if id.base, err = id.vcs.Base(); err != nil {
			logger.Warn("Could not find the base revision", "vcs", id.vcs.Name(), "error", err)
		}
		if author, err := id.vcs.Author(id.kind == store.KindLocal); err == nil {
			id.author = author
		}
	}

	// flags win over the VCS
	if cfg.Run.Revision != "" {
		id.current = vcs.ParseRevision(cfg.Run.Revision)
	}
	if cfg.Run.BaseRevision != "" {
		id.base = vcs.ParseRevision(cfg.Run.BaseRevision)
	}
	if cfg.Run.Author != "" {
		id.author = cfg.Run.Author
	}
	// a run matches a revision only without local edits; -revision promises
	// that when no adapter can check
	id.clean = id.current.ID != "" && !id.edits && (id.vcs != nil || cfg.Run.Revision != "")

	switch {
	case cfg.Run.Stream != "":
		id.stream = cfg.Run.Stream
	case id.current.Branch != "":
		id.stream = id.current.Branch
	case id.base.Branch != "":
		id.stream = id.base.Branch
	}

	// a run never compares with its own revision: a clean workspace at the
	// base revision compares with the revision before it
	if id.base.ID == "" && id.clean {
		id.base = id.current
	}
	id.olderOnly = id.clean && id.base.ID == id.current.ID
	return id
}

// versions is what the report says it compares; nil when nothing is known.
func (id runIdentity) versions(diffFile string, hasDiff bool) *model.Versions {
	v := &model.Versions{
		Base:       model.RunRef{Revision: id.base.ID, Stream: id.base.Branch},
		Current:    model.RunRef{Revision: id.current.ID, Stream: id.stream},
		LocalEdits: id.edits,
	}
	switch {
	case !hasDiff:
	case diffFile != "":
		v.Diff = filepath.Base(diffFile)
	case id.vcs != nil:
		v.Diff = id.vcs.Name()
	}
	if v.Base.Revision == "" && v.Current.Revision == "" && v.Diff == "" {
		return nil
	}
	return v
}

// vcsDiff makes the diff of the change when no diff file was given.
func vcsDiff(id runIdentity, logger *slog.Logger) *diff.DiffData {
	if id.vcs == nil || id.base.ID == "" || id.olderOnly {
		return nil
	}
	d, err := id.vcs.Diff(id.base)
	if err != nil {
		logger.Warn("Could not make the diff of the change; continuing without patch coverage", "vcs", id.vcs.Name(), "error", err)
		return nil
	}
	logger.Info("Diff made by the VCS adapter", "vcs", id.vcs.Name(), "base", id.base.ID, "files", len(d.Files))
	return d
}

// historyResult is what the run store step reports back to the terminal summary.
type historyResult struct {
	runID  int64
	runURL string
	notes  []string
}

// recordHistory compares the run with its base run and saves it. It never
// fails the build: a missing store, server or base run only costs the delta.
func recordHistory(cfg *config.AppConfig, tree *model.SummaryTree, id runIdentity, logger *slog.Logger) *historyResult {
	if cfg.History.Store == "" {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	res := &historyResult{}

	target, closeTarget, err := openTarget(cfg)
	if err != nil {
		logger.Warn("Run store is not available; continuing without history", "store", cfg.History.Store, "error", err)
		res.notes = append(res.notes, "History is off for this run: "+err.Error())
		return res
	}
	defer closeTarget()

	capture, err := history.NewCapture(tree, cfg, id.stream)
	if err != nil {
		logger.Warn("Could not prepare the run for the store", "error", err)
		return res
	}

	project := projectName(cfg)
	if base, err := history.FindBase(ctx, target, id.vcs, history.BaseQuery{
		Project: project, Stream: id.stream, Profile: cfg.History.Profile,
		Base: id.base, OlderOnly: id.olderOnly, MaxDistance: cfg.History.MaxDistance,
		Kinds: baseKinds(cfg),
	}); err != nil {
		logger.Warn("Could not search for a base run", "error", err)
		res.notes = append(res.notes, "No delta: the search for a base run failed.")
	} else if base == nil {
		res.notes = append(res.notes, noBaseNote(id))
	} else if baseManifest, err := target.Manifest(ctx, base.Run.Manifest); err != nil {
		logger.Warn("Could not read the base run", "run", base.Run.ID, "error", err)
	} else {
		baseMeta, _ := history.ParseMeta(base.Run.Meta)
		cmp := compare.Compare(baseManifest, capture.Manifest, compare.Options{
			Include:        cfg.FileFilterInstance.IsElementIncludedInReport,
			FileMetrics:    cfg.FileMetrics,
			HasDiff:        tree.Change != nil,
			Change:         tree.Change,
			BaseReports:    baseMeta.Reports,
			CurrentReports: tree.ReportNames,
			BaseTool:       base.Run.Tool,
			CurrentTool:    version,
		})
		cmp.Base = history.Ref(base.Run)
		cmp.Exact, cmp.Distance = base.Exact, base.Distance
		cmp.Current = model.RunRef{Revision: id.current.ID, Stream: id.stream, Profile: cfg.History.Profile, Kind: string(id.kind), Tool: version}
		tree.Comparison = cmp
		capture.Meta.Base = &history.BaseRef{RunRef: cmp.Base, Exact: base.Exact, Distance: base.Distance}
	}

	if reason := saveBlocked(cfg, id); reason != "" {
		res.notes = append(res.notes, reason)
		return res
	}
	run := store.Run{
		Project: project, Stream: id.stream, Profile: cfg.History.Profile, Kind: id.kind,
		Revision: id.current.ID, Change: id.current.Change, Clean: id.clean,
		ReviewID: cfg.Run.ReviewID, Author: id.author, CIURL: ciURL(cfg), Tool: version, CreatedAt: time.Now(),
	}
	runID, err := history.Save(ctx, target, capture, run, logger)
	if err != nil {
		logger.Warn("Could not save the run", "store", cfg.History.Store, "error", err)
		res.notes = append(res.notes, "The run was not saved: "+err.Error())
		return res
	}
	res.runID = runID
	if linker, ok := target.(interface{ RunURL(int64) string }); ok {
		res.runURL = linker.RunURL(runID)
	}
	if tree.Comparison != nil {
		tree.Comparison.Current.ID = runID
	}
	logger.Info("Run saved", "id", runID, "store", cfg.History.Store, "kind", id.kind, "revision", id.current.ID)

	if s, ok := target.(*store.Store); ok {
		cleanLocalStore(ctx, s, cfg, logger)
	}
	return res
}

// saveBlocked applies the rules of the run kinds; an empty string allows saving.
func saveBlocked(cfg *config.AppConfig, id runIdentity) string {
	switch {
	case id.kind == store.KindLocal && cfg.History.IsServer():
		return "Local runs read base runs from the team server but are not uploaded to it."
	case id.kind != store.KindLocal && cfg.History.IsServer() && os.Getenv(tokenEnv) == "":
		return "The run was not uploaded: set " + tokenEnv + " to the upload token of the team server."
	case id.kind == store.KindSubmit && id.current.ID == "":
		return "The run was not saved: a submit run needs -revision, or a VCS adapter (vcs.type)."
	case id.kind == store.KindSubmit && !id.clean:
		return "The run was not saved: a submit run needs a workspace without local edits."
	}
	return ""
}

// the environment variable that holds the upload token of a team server
const tokenEnv = "NANOVISION_TOKEN"

func openTarget(cfg *config.AppConfig) (history.Target, func(), error) {
	if cfg.History.IsServer() {
		c := client.New(cfg.History.Store, os.Getenv(tokenEnv))
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := c.Ping(ctx); err != nil {
			return nil, nil, fmt.Errorf("team server %s does not answer: %w", cfg.History.Store, err)
		}
		return c, func() {}, nil
	}
	s, err := store.Open(storeDir(cfg))
	if err != nil {
		return nil, nil, err
	}
	return s, func() { s.Close() }, nil
}

// storeDir resolves a local store folder against the project root.
func storeDir(cfg *config.AppConfig) string {
	dir := cfg.History.Store
	if !filepath.IsAbs(dir) {
		dir = filepath.Join(cfg.ProjectRoot, dir)
	}
	return dir
}

func projectName(cfg *config.AppConfig) string {
	if cfg.History.Project != "" {
		return cfg.History.Project
	}
	return filepath.Base(cfg.ProjectRoot)
}

// base run candidates: submit runs on a team server, and in a local store
// the developer's own clean runs too
func baseKinds(cfg *config.AppConfig) []store.RunKind {
	if cfg.History.IsServer() {
		return []store.RunKind{store.KindSubmit}
	}
	return []store.RunKind{store.KindSubmit, store.KindLocal}
}

func noBaseNote(id runIdentity) string {
	switch {
	case id.base.ID == "" && id.vcs == nil:
		return "No delta: set vcs.type, or pass -base-revision, so nanovision knows the base revision."
	case id.base.ID == "":
		return "No delta: the base revision is unknown."
	case id.olderOnly:
		return "No delta: the store has no earlier run to compare with."
	default:
		return fmt.Sprintf("No delta: the store has no base run at or before %s.", id.base.ID)
	}
}

func cleanLocalStore(ctx context.Context, s *store.Store, cfg *config.AppConfig, logger *slog.Logger) {
	policy := store.DefaultRetention()
	policy.LocalKeep = cfg.History.KeepLocal
	deleted, err := s.Cleanup(ctx, policy, time.Now())
	if err != nil {
		logger.Warn("Store cleanup failed", "error", err)
		return
	}
	if deleted == 0 {
		return
	}
	gc, err := s.GC(ctx)
	if err != nil {
		logger.Warn("Store cleanup failed", "error", err)
		return
	}
	logger.Info("Old runs removed from the local store", "runs", deleted, "blobs", gc.Blobs, "bytes", gc.Bytes)
}

// ciURL is the -ci-url flag, or the build link the CI system sets.
func ciURL(cfg *config.AppConfig) string {
	if cfg.Run.CIURL != "" {
		return cfg.Run.CIURL
	}
	for _, name := range []string{"BUILD_URL", "CI_JOB_URL"} { // Jenkins, GitLab
		if v := os.Getenv(name); v != "" {
			return v
		}
	}
	if server, repo, run := os.Getenv("GITHUB_SERVER_URL"), os.Getenv("GITHUB_REPOSITORY"), os.Getenv("GITHUB_RUN_ID"); server != "" && repo != "" && run != "" {
		return strings.Join([]string{server, repo, "actions", "runs", run}, "/")
	}
	return ""
}
