package vcs

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/IgorBayerl/nanovision/internal/diff"
)

type gitVCS struct {
	root       string // the project root; commands run here
	baseBranch string
	run        Runner
}

func newGit(projectRoot string, opts Options) (*gitVCS, error) {
	g := &gitVCS{root: projectRoot, baseBranch: opts.BaseBranch, run: opts.Runner}
	top, err := g.git("rev-parse", "--show-toplevel")
	if err != nil {
		return nil, fmt.Errorf("not a git repository: %w", err)
	}
	if !within(projectRoot, top) {
		return nil, fmt.Errorf("project root %s is outside the git repository %s", projectRoot, top)
	}
	return g, nil
}

func (g *gitVCS) Name() string { return "git" }

func (g *gitVCS) git(args ...string) (string, error) {
	out, err := g.run.Run(g.root, "git", args...)
	return strings.TrimSpace(string(out)), err
}

func (g *gitVCS) Current() (Revision, error) {
	hash, err := g.git("rev-parse", "HEAD")
	if err != nil {
		return Revision{}, err
	}
	rev := Revision{ID: hash}
	// a detached checkout, common in CI, has no branch; -stream can name it
	if branch, err := g.git("rev-parse", "--abbrev-ref", "HEAD"); err == nil && branch != "HEAD" {
		rev.Branch = branch
	}
	return rev, nil
}

// untracked files do not count: build outputs are often untracked
func (g *gitVCS) HasLocalEdits() (bool, error) {
	out, err := g.git("status", "--porcelain", "--untracked-files=no")
	if err != nil {
		return false, err
	}
	return out != "", nil
}

// Base is the merge-base of HEAD and the base branch. Without any base branch,
// for example in a repository with no remote, it is HEAD itself, so the change
// is the uncommitted edits.
func (g *gitVCS) Base() (Revision, error) {
	branch := g.findBaseBranch()
	if branch == "" {
		return g.Current()
	}
	hash, err := g.git("merge-base", "HEAD", branch)
	if err != nil {
		return Revision{}, fmt.Errorf("merge-base with %s (a shallow clone has no merge-base; fetch the full history): %w", branch, err)
	}
	return Revision{ID: hash, Branch: strings.TrimPrefix(branch, "origin/")}, nil
}

func (g *gitVCS) findBaseBranch() string {
	if g.baseBranch != "" {
		return g.baseBranch
	}
	if ref, err := g.git("symbolic-ref", "--quiet", "refs/remotes/origin/HEAD"); err == nil && ref != "" {
		return strings.TrimPrefix(ref, "refs/remotes/")
	}
	for _, candidate := range []string{"origin/main", "origin/master"} {
		if _, err := g.git("rev-parse", "--verify", "--quiet", candidate+"^{commit}"); err == nil {
			return candidate
		}
	}
	return ""
}

// Diff is the working tree against base: the commits since base and the
// uncommitted edits, plus untracked files with every line added. --relative
// keeps only the files under the project root, with paths relative to it.
func (g *gitVCS) Diff(base Revision) (*diff.DiffData, error) {
	out, err := g.run.Run(g.root, "git", "diff", "--no-color", "--no-ext-diff", "--relative",
		"--src-prefix=a/", "--dst-prefix=b/", base.ID)
	if err != nil {
		return nil, err
	}
	data, err := diff.ParseText(string(out), nil)
	if err != nil {
		return nil, err
	}

	untracked, err := g.git("ls-files", "--others", "--exclude-standard")
	if err != nil {
		return nil, err
	}
	for _, path := range strings.Split(untracked, "\n") {
		if path = strings.TrimSpace(path); path == "" {
			continue
		}
		if f, ok := addedFile(filepath.Join(g.root, filepath.FromSlash(path)), path); ok {
			data.Files = append(data.Files, f)
		}
	}
	return data, nil
}

// Ancestors lists base and its first parents, max+1 commits in all.
func (g *gitVCS) Ancestors(from Revision, max int) ([]string, error) {
	out, err := g.git("rev-list", "--first-parent", "--max-count="+strconv.Itoa(max+1), from.ID)
	if err != nil {
		return nil, err
	}
	if out == "" {
		return nil, nil
	}
	return strings.Split(out, "\n"), nil
}

func (g *gitVCS) Author(local bool) (string, error) {
	if local {
		return g.git("config", "user.email")
	}
	return g.git("log", "-1", "--format=%ae")
}

// addedFile marks every line of a new file as added.
func addedFile(abs, rel string) (diff.FileDiff, bool) {
	content, err := os.ReadFile(abs)
	if err != nil || bytes.IndexByte(content, 0) >= 0 {
		return diff.FileDiff{}, false // unreadable, or binary
	}
	lines := bytes.Count(content, []byte{'\n'})
	if len(content) > 0 && content[len(content)-1] != '\n' {
		lines++
	}
	hunk := diff.Hunk{NewStart: 1, NewLines: lines, AddedLineOffsets: make([]int, lines)}
	for i := range lines {
		hunk.AddedLineOffsets[i] = i
	}
	return diff.FileDiff{OldPath: "/dev/null", NewPath: filepath.ToSlash(rel), Kind: "added", Hunks: []diff.Hunk{hunk}}, true
}

// within reports whether path is dir or below it. Both are cleaned and
// compared case-insensitively on Windows.
func within(path, dir string) bool {
	rel, err := filepath.Rel(filepath.Clean(filepath.FromSlash(dir)), filepath.Clean(filepath.FromSlash(path)))
	if err != nil {
		return false
	}
	return rel == "." || (rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)))
}
