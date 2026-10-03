package vcs

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/IgorBayerl/nanovision/internal/diff"
)

// The Perforce adapter follows section 6.2 of docs/plans/coverage-delta-vs-base.md.
// It uses the environment of the user: P4PORT, P4USER, P4CLIENT, P4CONFIG and
// the login ticket.
type perforceVCS struct {
	root       string // the project root
	run        Runner
	client     string
	clientRoot string
	user       string
	stream     string // the stream in run keys: the parent of a virtual stream
}

func newPerforce(projectRoot string, opts Options) (*perforceVCS, error) {
	p := &perforceVCS{root: projectRoot, run: opts.Runner}
	info, err := p.record("info")
	if err != nil {
		return nil, err
	}
	p.client, p.clientRoot, p.user = info["clientName"], info["clientRoot"], info["userName"]
	if p.client == "" || p.client == "*unknown*" {
		return nil, fmt.Errorf("no Perforce client workspace is set (P4CLIENT)")
	}
	if !within(projectRoot, p.clientRoot) {
		return nil, fmt.Errorf("project root %s is outside the Perforce workspace root %s", projectRoot, p.clientRoot)
	}

	client, err := p.record("client", "-o")
	if err != nil {
		return nil, err
	}
	p.stream = client["Stream"]
	if p.stream != "" {
		// a virtual stream has no changelists of its own
		if spec, err := p.record("stream", "-o", p.stream); err == nil && strings.EqualFold(spec["Type"], "virtual") && spec["Parent"] != "" {
			p.stream = spec["Parent"]
		}
	}
	return p, nil
}

func (p *perforceVCS) Name() string { return "perforce" }

func (p *perforceVCS) p4(args ...string) (string, error) {
	out, err := p.run.Run(p.root, "p4", append([]string{"-ztag"}, args...)...)
	if err != nil {
		if strings.Contains(err.Error(), "expired") || strings.Contains(err.Error(), "P4PASSWD") {
			return "", fmt.Errorf("%w (run `p4 login`)", err)
		}
		return "", err
	}
	return string(out), nil
}

func (p *perforceVCS) records(args ...string) ([]map[string]string, error) {
	out, err := p.p4(args...)
	if err != nil {
		return nil, err
	}
	return parseZtag(out), nil
}

func (p *perforceVCS) record(args ...string) (map[string]string, error) {
	recs, err := p.records(args...)
	if err != nil {
		return nil, err
	}
	if len(recs) == 0 {
		return map[string]string{}, nil
	}
	return recs[0], nil
}

// Current is the newest changelist that changed a file the workspace has.
func (p *perforceVCS) Current() (Revision, error) {
	rec, err := p.record("changes", "-m1", "//"+p.client+"/...#have")
	if err != nil {
		return Revision{}, err
	}
	change, err := strconv.ParseInt(rec["change"], 10, 64)
	if err != nil || change <= 0 {
		return Revision{}, fmt.Errorf("workspace %s has no synced changelist", p.client)
	}
	return Revision{ID: PerforceRevision(p.stream, change), Branch: p.stream, Change: change}, nil
}

// pending and shelved changelists sit on top of the synced changelist
func (p *perforceVCS) Base() (Revision, error) { return p.Current() }

func (p *perforceVCS) HasLocalEdits() (bool, error) {
	recs, err := p.opened()
	return len(recs) > 0, err
}

func (p *perforceVCS) opened() ([]map[string]string, error) {
	recs, err := p.records("opened")
	if err != nil && strings.Contains(err.Error(), "not opened") {
		return nil, nil
	}
	var files []map[string]string
	for _, r := range recs {
		if r["depotFile"] != "" {
			files = append(files, r)
		}
	}
	return files, err
}

// Diff is every opened file against the synced revisions. The diff covers
// the whole workspace, not one pending changelist, like the test run does.
func (p *perforceVCS) Diff(Revision) (*diff.DiffData, error) {
	out, err := p.run.Run(p.root, "p4", "diff", "-du")
	if err != nil {
		return nil, err
	}
	data, err := diff.ParseText(string(out), nil)
	if err != nil {
		return nil, err
	}
	var files []diff.FileDiff
	have := make(map[string]bool)
	for _, f := range data.Files {
		// the +++ line holds a local absolute path
		rel, ok := p.relative(f.NewPath)
		if !ok {
			continue
		}
		f.NewPath = rel
		f.OldPath = rel
		have[rel] = true
		files = append(files, f)
	}

	// files opened for add, branch or a move target are not in `p4 diff`
	opened, err := p.opened()
	if err != nil {
		return nil, err
	}
	for _, r := range opened {
		switch r["action"] {
		case "add", "branch", "move/add", "import":
		case "delete", "move/delete":
			if rel, ok := p.relative(p.localPath(r["clientFile"])); ok && !have[rel] {
				files = append(files, diff.FileDiff{OldPath: rel, NewPath: "/dev/null", Kind: "modified"})
			}
			continue
		default:
			continue
		}
		local := p.localPath(r["clientFile"])
		rel, ok := p.relative(local)
		if !ok || have[rel] {
			continue
		}
		if f, ok := addedFile(local, rel); ok {
			files = append(files, f)
		}
	}
	data.Files = files
	return data, nil
}

// ChangesBetween counts the submitted changelists in (low, high] that touch
// the workspace view. Art changelists outside the view do not count.
func (p *perforceVCS) ChangesBetween(low, high int64) (int, error) {
	if high <= low {
		return 0, nil
	}
	recs, err := p.records("changes", "-s", "submitted", fmt.Sprintf("//%s/...@%d,@%d", p.client, low+1, high))
	if err != nil {
		return 0, err
	}
	n := 0
	for _, r := range recs {
		if r["change"] != "" {
			n++
		}
	}
	return n, nil
}

func (p *perforceVCS) Author(local bool) (string, error) { return p.user, nil }

// localPath turns "//client/src/a.cpp" into a path under the client root.
func (p *perforceVCS) localPath(clientFile string) string {
	rest, ok := strings.CutPrefix(clientFile, "//"+p.client+"/")
	if !ok {
		return clientFile
	}
	return filepath.Join(p.clientRoot, filepath.FromSlash(rest))
}

func (p *perforceVCS) relative(local string) (string, bool) {
	if !filepath.IsAbs(local) {
		return "", false
	}
	rel, err := filepath.Rel(p.root, local)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", false
	}
	return filepath.ToSlash(rel), true
}

// parseZtag reads `p4 -ztag` output: "... key value" lines, one record per
// block. Blank lines separate records.
func parseZtag(out string) []map[string]string {
	var recs []map[string]string
	cur := map[string]string{}
	flush := func() {
		if len(cur) > 0 {
			recs = append(recs, cur)
			cur = map[string]string{}
		}
	}
	for _, line := range strings.Split(strings.ReplaceAll(out, "\r\n", "\n"), "\n") {
		rest, ok := strings.CutPrefix(line, "... ")
		if !ok {
			if strings.TrimSpace(line) == "" {
				flush()
			}
			continue
		}
		key, value, _ := strings.Cut(rest, " ")
		if _, dup := cur[key]; dup {
			flush() // records without a blank line between them
		}
		cur[key] = value
	}
	flush()
	return recs
}
