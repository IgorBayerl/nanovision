// Package vcs asks git or Perforce the questions the coverage delta needs:
// which revision this is, what it is based on, and what changed.
package vcs

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/IgorBayerl/nanovision/internal/diff"
)

// Revision is one state of the code.
type Revision struct {
	// git: a commit hash. Perforce: "//stream@changelist", or "@changelist"
	// in a classic depot
	ID string
	// git branch or Perforce stream, for display and for grouping runs
	Branch string
	// the Perforce changelist number; 0 for git
	Change int64
}

// VCS is a version control adapter.
type VCS interface {
	Name() string
	// the revision the workspace is at
	Current() (Revision, error)
	// true when files differ from the current revision
	HasLocalEdits() (bool, error)
	// the revision the change starts from
	Base() (Revision, error)
	// the change against base, with paths relative to the project root
	Diff(base Revision) (*diff.DiffData, error)
	// who made the change: the commit author, or the user of the workspace
	Author(local bool) (string, error)
}

// Ancestry is implemented by git: the first-parent history below a revision,
// newest first and starting with the revision itself.
type Ancestry interface {
	Ancestors(from Revision, max int) ([]string, error)
}

// ChangeCounter is implemented by Perforce: how many submitted changelists
// in the workspace view lie in (low, high].
type ChangeCounter interface {
	ChangesBetween(low, high int64) (int, error)
}

// Runner runs one command. Tests replace it with recorded output.
type Runner interface {
	Run(dir, name string, args ...string) ([]byte, error)
}

// CommandTimeout bounds every command; a Perforce server can be slow to answer.
const CommandTimeout = 30 * time.Second

type execRunner struct{}

func (execRunner) Run(dir, name string, args ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CommandTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return nil, fmt.Errorf("%s %s: no answer after %s", name, strings.Join(args, " "), CommandTimeout)
	}
	if err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return stdout.Bytes(), &CommandError{Command: name + " " + strings.Join(args, " "), Message: msg}
	}
	return stdout.Bytes(), nil
}

// CommandError is a command that ran and failed.
type CommandError struct {
	Command string
	Message string
}

func (e *CommandError) Error() string { return e.Command + ": " + e.Message }

// DefaultRunner runs real commands.
var DefaultRunner Runner = execRunner{}

// Options configure Detect.
type Options struct {
	// git only; empty tries origin/HEAD, origin/main, origin/master
	BaseBranch string
	Runner     Runner
}

// ErrNotFound means auto detection found neither git nor Perforce.
var ErrNotFound = errors.New("no git repository or Perforce workspace contains the project root")

// Detect returns the adapter for kind: "git" or "perforce". An empty kind
// gives nil: version control is off.
func Detect(kind, projectRoot string, opts Options) (VCS, error) {
	if opts.Runner == nil {
		opts.Runner = DefaultRunner
	}
	switch strings.ToLower(kind) {
	case "":
		return nil, nil
	case "git":
		return newGit(projectRoot, opts)
	case "perforce", "p4":
		return newPerforce(projectRoot, opts)
	default:
		return nil, fmt.Errorf("unknown vcs type %q", kind)
	}
}

var perforceRevisionRE = regexp.MustCompile(`^(//[^@]+)?@(\d+)$`)

// ParseRevision reads a revision flag: "//stream@changelist" or "@changelist"
// for Perforce, anything else is a git commit.
func ParseRevision(s string) Revision {
	s = strings.TrimSpace(s)
	if m := perforceRevisionRE.FindStringSubmatch(s); m != nil {
		change, _ := strconv.ParseInt(m[2], 10, 64)
		return Revision{ID: s, Branch: m[1], Change: change}
	}
	return Revision{ID: s}
}

// PerforceRevision formats a stream and changelist the way runs store them.
func PerforceRevision(stream string, change int64) string {
	return stream + "@" + strconv.FormatInt(change, 10)
}
