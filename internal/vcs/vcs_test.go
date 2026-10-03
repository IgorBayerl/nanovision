package vcs

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseRevision(t *testing.T) {
	assert.Equal(t, Revision{ID: "//game/main@118432", Branch: "//game/main", Change: 118432}, ParseRevision("//game/main@118432"))
	assert.Equal(t, Revision{ID: "@42", Change: 42}, ParseRevision(" @42 "))
	assert.Equal(t, Revision{ID: "3f2c9ab"}, ParseRevision("3f2c9ab"))
	assert.Equal(t, "//game/main@7", PerforceRevision("//game/main", 7))
}

func TestDetectNone(t *testing.T) {
	v, err := Detect("", t.TempDir(), Options{})
	assert.NoError(t, err)
	assert.Nil(t, v)
}

// --- git, against a real temporary repository ---

func gitRepo(t *testing.T) string {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not installed")
	}
	dir := t.TempDir()
	run := func(args ...string) {
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(), "GIT_AUTHOR_DATE=2026-01-01T00:00:00Z", "GIT_COMMITTER_DATE=2026-01-01T00:00:00Z")
		out, err := cmd.CombinedOutput()
		require.NoError(t, err, "git %s: %s", strings.Join(args, " "), out)
	}
	write := func(path, content string) {
		full := filepath.Join(dir, filepath.FromSlash(path))
		require.NoError(t, os.MkdirAll(filepath.Dir(full), 0o755))
		require.NoError(t, os.WriteFile(full, []byte(content), 0o644))
	}

	run("init", "--quiet", "--initial-branch=main")
	run("config", "user.email", "dev@example.com")
	run("config", "user.name", "Dev")
	run("config", "core.autocrlf", "false")
	write("game/src/inventory.cpp", "int a() {\n  return 1;\n}\n")
	write("tools/build.py", "print(1)\n")
	run("add", ".")
	run("commit", "--quiet", "-m", "first")
	write("game/src/inventory.cpp", "int a() {\n  return 2;\n}\n")
	run("commit", "--quiet", "-am", "second")

	run("checkout", "--quiet", "-b", "feature")
	write("game/src/item.cpp", "int b() {\n  return 3;\n}\n")
	write("tools/build.py", "print(2)\n")
	run("add", ".")
	run("commit", "--quiet", "-m", "feature work")
	write("game/src/inventory.cpp", "int a() {\n  return 2;\n}\nint c() { return 4; }\n") // uncommitted
	write("game/tests/new_test.cpp", "TEST(x) {}\n")                                      // untracked
	return dir
}

func TestGitAdapter(t *testing.T) {
	repo := gitRepo(t)
	project := filepath.Join(repo, "game")

	v, err := Detect("git", project, Options{BaseBranch: "main"})
	require.NoError(t, err)
	require.Equal(t, "git", v.Name())

	cur, err := v.Current()
	require.NoError(t, err)
	assert.Len(t, cur.ID, 40)
	assert.Equal(t, "feature", cur.Branch)

	edits, err := v.HasLocalEdits()
	require.NoError(t, err)
	assert.True(t, edits)

	base, err := v.Base()
	require.NoError(t, err)
	assert.NotEqual(t, cur.ID, base.ID)
	assert.Equal(t, "main", base.Branch)

	d, err := v.Diff(base)
	require.NoError(t, err)
	var paths []string
	for _, f := range d.Files {
		paths = append(paths, f.NewPath+":"+f.Kind)
	}
	sort.Strings(paths)
	assert.Equal(t, []string{"src/inventory.cpp:modified", "src/item.cpp:added", "tests/new_test.cpp:added"}, paths,
		"paths are relative to the project root, and tools/ outside it is left out")

	ancestry, ok := v.(Ancestry)
	require.True(t, ok)
	ancestors, err := ancestry.Ancestors(base, 5)
	require.NoError(t, err)
	require.Len(t, ancestors, 2)
	assert.Equal(t, base.ID, ancestors[0])

	author, err := v.Author(true)
	require.NoError(t, err)
	assert.Equal(t, "dev@example.com", author)
}

func TestGitWithoutBaseBranchComparesWithHead(t *testing.T) {
	repo := gitRepo(t)
	v, err := Detect("git", repo, Options{})
	require.NoError(t, err)

	base, err := v.Base()
	require.NoError(t, err)
	cur, _ := v.Current()
	assert.Equal(t, cur.ID, base.ID, "no remote: the change is the uncommitted edits")

	d, err := v.Diff(base)
	require.NoError(t, err)
	assert.Len(t, d.Files, 2, "the edit and the untracked file")
}

func TestGitDetectOutsideRepository(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not installed")
	}
	_, err := Detect("git", t.TempDir(), Options{})
	assert.Error(t, err)
}

// --- Perforce, against recorded command output ---

type fakeRunner struct {
	outputs map[string]string
	errors  map[string]string
	calls   []string
}

func (f *fakeRunner) Run(_ string, name string, args ...string) ([]byte, error) {
	key := name + " " + strings.Join(args, " ")
	f.calls = append(f.calls, key)
	if msg, ok := f.errors[key]; ok {
		return nil, &CommandError{Command: key, Message: msg}
	}
	out, ok := f.outputs[key]
	if !ok {
		return nil, fmt.Errorf("unexpected command %q", key)
	}
	return []byte(out), nil
}

func perforceFixture(t *testing.T) (*fakeRunner, string) {
	root := t.TempDir()
	project := filepath.Join(root, "game")
	require.NoError(t, os.MkdirAll(filepath.Join(project, "src"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(project, "src", "new.cpp"), []byte("int n() {\n  return 0;\n}\n"), 0o644))

	local := func(p string) string { return filepath.Join(root, filepath.FromSlash(p)) }
	runner := &fakeRunner{outputs: map[string]string{
		"p4 -ztag info":                          "... userName dev\n... clientName dev-ws\n... clientRoot " + root + "\n... serverAddress p4:1666\n",
		"p4 -ztag client -o":                     "... Client dev-ws\n... Root " + root + "\n... Stream //game/dev-virtual\n",
		"p4 -ztag stream -o //game/dev-virtual":  "... Stream //game/dev-virtual\n... Type virtual\n... Parent //game/main\n",
		"p4 -ztag changes -m1 //dev-ws/...#have": "... change 118420\n... status submitted\n... user build\n",
		"p4 -ztag opened": "... depotFile //game/main/game/src/inventory.cpp\n... clientFile //dev-ws/game/src/inventory.cpp\n... action edit\n\n" +
			"... depotFile //game/main/game/src/new.cpp\n... clientFile //dev-ws/game/src/new.cpp\n... action add\n\n" +
			"... depotFile //game/main/game/src/old.cpp\n... clientFile //dev-ws/game/src/old.cpp\n... action delete\n",
		"p4 diff -du": "--- //game/main/game/src/inventory.cpp\t2026/09/20 10:00:00\n" +
			"+++ " + local("game/src/inventory.cpp") + "\t2026/09/27 12:00:00\n" +
			"@@ -1,2 +1,3 @@\n int a() {\n+  log();\n   return 2;\n",
		"p4 -ztag changes -s submitted //dev-ws/...@118421,@118432": "... change 118432\n\n... change 118425\n",
	}}
	return runner, project
}

func TestPerforceAdapter(t *testing.T) {
	runner, project := perforceFixture(t)
	v, err := Detect("perforce", project, Options{Runner: runner})
	require.NoError(t, err)

	cur, err := v.Current()
	require.NoError(t, err)
	assert.Equal(t, Revision{ID: "//game/main@118420", Branch: "//game/main", Change: 118420}, cur, "a virtual stream uses its parent")

	edits, err := v.HasLocalEdits()
	require.NoError(t, err)
	assert.True(t, edits)

	d, err := v.Diff(cur)
	require.NoError(t, err)
	byPath := map[string]string{}
	for _, f := range d.Files {
		byPath[f.NewPath+"|"+f.OldPath] = f.Kind
	}
	assert.Equal(t, map[string]string{
		"src/inventory.cpp|src/inventory.cpp": "modified",
		"src/new.cpp|/dev/null":               "added",
		"/dev/null|src/old.cpp":               "modified",
	}, byPath)
	for _, f := range d.Files {
		if f.NewPath == "src/new.cpp" {
			assert.Len(t, f.Hunks[0].AddedLineOffsets, 3, "every line of an added file is added")
		}
	}

	counter, ok := v.(ChangeCounter)
	require.True(t, ok)
	n, err := counter.ChangesBetween(118420, 118432)
	require.NoError(t, err)
	assert.Equal(t, 2, n)

	author, err := v.Author(false)
	require.NoError(t, err)
	assert.Equal(t, "dev", author)
}

func TestPerforceNeedsAWorkspace(t *testing.T) {
	runner := &fakeRunner{outputs: map[string]string{"p4 -ztag info": "... userName dev\n... clientName *unknown*\n"}}
	_, err := Detect("perforce", t.TempDir(), Options{Runner: runner})
	assert.ErrorContains(t, err, "P4CLIENT")
}

func TestPerforceExpiredTicket(t *testing.T) {
	runner, project := perforceFixture(t)
	v, err := Detect("perforce", project, Options{Runner: runner})
	require.NoError(t, err)
	runner.errors = map[string]string{"p4 -ztag changes -m1 //dev-ws/...#have": "Your session has expired, please login again."}
	_, err = v.Current()
	assert.ErrorContains(t, err, "p4 login")
}

func TestParseZtag(t *testing.T) {
	recs := parseZtag("... change 3\n... user a\n... change 2\n... user b\n\n\n... change 1\r\n")
	require.Len(t, recs, 3)
	assert.Equal(t, "3", recs[0]["change"])
	assert.Equal(t, "b", recs[1]["user"])
	assert.Equal(t, "1", recs[2]["change"])
}
