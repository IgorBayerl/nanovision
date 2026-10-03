"""
End-to-end checks for run history: the run store, the coverage delta against
the base run, and the team server (docs/plans/run-store-and-server.md and
docs/plans/coverage-delta-vs-base.md).

Used by e2e_test.py with --history (and --serve to browse the result).

Everything happens in reports/history/:
  repo/    a git repository made from the Go demo project
  store/   the local run store the checks write to
  team/    the store of a temporary team server
"""
import json
import os
import shutil
import socket
import stat
import subprocess
import sys
import time
import urllib.request

SCRIPT_ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
DEMO = os.path.join(SCRIPT_ROOT, "demo_projects")
WORK = os.path.join(SCRIPT_ROOT, "reports", "history")
REPO = os.path.join(WORK, "repo")
STORE = os.path.join(WORK, "store")
TEAM_STORE = os.path.join(WORK, "team")
TOKEN = "e2e-upload-token"

CONFIG = """\
title: "Calculator (history e2e)"
reports: ["cov.out"]
source_dirs: ["project"]
report_types: ["TextSummary"]
output_dir: "out"
verbosity: "Warning"
ignore_files: ["**/*_test.go"]
status_bands:
  statement_coverage: "60..80"
history:
  project: "calculator"
  profile: "unit"
vcs:
  type: "git"
  base_branch: "main"
"""


class Checks:
    """Collects one result per check, in the shape e2e_test.py prints."""

    def __init__(self):
        self.results = []

    def add(self, name, ok, details):
        self.results.append({
            "name": f"History: {name}",
            "status": "✅ SUCCESS" if ok else "❌ FAILED",
            "details": details,
        })
        print(f"--- {'PASS' if ok else 'FAIL'}: {name} ---")
        if not ok:
            print(details, file=sys.stderr)

    def expect(self, name, output, *wanted):
        missing = [w for w in wanted if w not in output]
        if missing:
            self.add(name, False, f"missing {missing!r} in output:\n{output}")
        else:
            self.add(name, True, "output contains " + ", ".join(repr(w) for w in wanted))
        return not missing


def git(*args):
    subprocess.run(["git", *args], cwd=REPO, check=True, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)


def nanovision(binary, *args, cwd=REPO, env=None):
    """Runs the CLI and returns its exit code and its terminal summary (stdout)."""
    full_env = os.environ.copy()
    full_env.pop("NANOVISION_TOKEN", None)
    if env:
        full_env.update(env)
    print(f"--- Running Command: nanovision {' '.join(args)}")
    proc = subprocess.run([binary, *args], cwd=cwd, env=full_env, capture_output=True, text=True, encoding="utf-8")
    return proc.returncode, proc.stdout


def edit_coverage(replacements):
    path = os.path.join(REPO, "cov.out")
    with open(path, encoding="utf-8") as f:
        text = f.read()
    for old, new in replacements:
        text = text.replace(old, new)
    with open(path, "w", encoding="utf-8", newline="\n") as f:
        f.write(text)


def append(path, text):
    with open(os.path.join(REPO, path), "a", encoding="utf-8", newline="\n") as f:
        f.write(text)


def remove_tree(path):
    """shutil.rmtree that also removes read-only files, like git's objects on Windows."""
    def make_writable(func, target, _exc):
        os.chmod(target, stat.S_IWRITE)
        func(target)

    if not os.path.exists(path):
        return
    try:
        shutil.rmtree(path, onexc=make_writable)
    except PermissionError as e:
        sys.exit(f"--- Cannot clean {path}: {e.filename} is in use. "
                 "Stop the `nanovision serve` of an earlier --serve/--showcase run (Ctrl+C) and try again. ---")


def setup_repo():
    remove_tree(WORK)
    os.makedirs(REPO)
    shutil.copytree(os.path.join(DEMO, "go", "project"), os.path.join(REPO, "project"))
    shutil.copy(os.path.join(DEMO, "go", "report", "gocover", "coverage.out"), os.path.join(REPO, "cov.out"))
    with open(os.path.join(REPO, "nanovision.yaml"), "w", encoding="utf-8", newline="\n") as f:
        f.write(CONFIG)
    # coverage and outputs are build products, not local edits
    with open(os.path.join(REPO, ".gitignore"), "w", encoding="utf-8", newline="\n") as f:
        f.write("cov.out\nout/\n")
    git("init", "--quiet", "--initial-branch=main")
    git("config", "user.email", "ada@example.com")
    git("config", "user.name", "Ada")
    git("config", "core.autocrlf", "false")
    git("add", ".")
    git("commit", "--quiet", "-m", "base")


def local_store_checks(binary, checks):
    store = ["-config", "nanovision.yaml", "-store", STORE]

    code, out = nanovision(binary, *store, "-run-kind=submit")
    checks.expect("the first submit run is saved, with no base run yet", out if code == 0 else f"exit {code}\n{out}",
                  "run 1", "No delta: the store has no earlier run")

    # a commit that makes the tests reach one more statement
    edit_coverage([("calculator.go:19.2,19.14 1 0", "calculator.go:19.2,19.14 1 1")])
    append("project/calculator/entities.go", "// tweak\n")
    git("commit", "--quiet", "-am", "B")
    code, out = nanovision(binary, *store, "-run-kind=submit")
    checks.expect("a submit run compares with the previous revision", out,
                  "run 2", "Coverage delta against the base run", "(exact)")

    # the plan's key case: a change that only edits tests
    append("project/calculator/calculator_test.go", "// more tests\n")
    edit_coverage([("calculator.go:25.12,27.3 1 0", "calculator.go:25.12,27.3 1 1"),
                   ("calculator.go:32.36,33.12 1 0", "calculator.go:32.36,33.12 1 1")])
    code, out = nanovision(binary, *store)
    checks.expect("a test-only change shows the delta", out,
                  "run 3", "all ignored by the coverage filters (tests)",
                  "none: the change has no measured lines", "+4.61")
    summary = os.path.join(REPO, "out", "Summary.txt")
    with open(summary, encoding="utf-8") as f:
        checks.expect("Summary.txt carries the delta block", f.read(), "Coverage delta against the base run")
    git("checkout", "--quiet", "--", "project/calculator/calculator_test.go")

    # a review run on a branch that adds a tested function
    git("checkout", "--quiet", "-b", "feature/rounding")
    append("project/calculator/calculator.go",
           "\n// Round rounds to n digits.\nfunc Round(v float64, n int) float64 {\n\tif n < 0 {\n\t\treturn v\n\t}\n"
           "\tp := 1.0\n\tfor i := 0; i < n; i++ {\n\t\tp *= 10\n\t}\n\treturn float64(int(v*p+0.5)) / p\n}\n")
    git("commit", "--quiet", "-am", "feature: Round")
    with open(os.path.join(REPO, "project", "calculator", "calculator.go"), encoding="utf-8") as f:
        n = len(f.read().splitlines())
    with open(os.path.join(REPO, "cov.out"), "a", encoding="utf-8", newline="\n") as f:
        path = "test_project_go/calculator/calculator.go"
        f.write(f"{path}:{n - 9}.42,{n - 8}.12 1 1\n{path}:{n - 8}.12,{n - 6}.3 1 0\n"
                f"{path}:{n - 5}.2,{n - 4}.18 2 1\n{path}:{n - 4}.18,{n - 2}.3 1 1\n{path}:{n - 1}.2,{n - 1}.40 1 1\n")
    code, out = nanovision(binary, *store, "-run-kind=review", "-review-id=42")
    checks.expect("a review run gets patch coverage and the delta", out,
                  "run 4", "(exact)", "Patch coverage", "modified")
    git("checkout", "--quiet", "main")

    # Perforce-style revisions from flags only, for CI without a VCS adapter
    cpp = os.path.join(DEMO, "cpp")
    gcov = os.path.join(cpp, "report", "gcov", "branch-probabilities", "*.gcov")
    flags = ["-report=" + gcov, "-sourcedirs=" + os.path.join(cpp, "project"), "-reporttypes=TextSummary",
             "-output=" + os.path.join(WORK, "cpp-out"), "-verbosity=Warning", "-store", STORE,
             "-project", "game", "-profile", "unit-win64"]
    nanovision(binary, *flags, "-run-kind=submit", "-revision=//game/main@100", cwd=WORK)
    code, out = nanovision(binary, *flags, "-run-kind=review", "-revision=//game/main@112",
                           "-base-revision=//game/main@107", "-review-id=113", cwd=WORK)
    checks.expect("a Perforce base revision without a run finds the older run", out,
                  "//game/main@100", "older than the base revision")

    backup = os.path.join(WORK, "backup.db")
    code, out = nanovision(binary, "store", "backup", "-store", STORE, backup)
    checks.add("nanovision store backup writes a copy of the store", code == 0 and os.path.getsize(backup) > 0, out.strip())


def free_port():
    with socket.socket() as s:
        s.bind(("127.0.0.1", 0))
        return s.getsockname()[1]


def start_server(binary, store, port, token=None):
    args = [binary, "serve", "-store", store, "-addr", f"127.0.0.1:{port}"]
    env = os.environ.copy()
    env.pop("NANOVISION_TOKEN", None)
    if token:
        env["NANOVISION_TOKEN"] = token
    proc = subprocess.Popen(args, env=env, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
    url = f"http://127.0.0.1:{port}"
    for _ in range(100):
        try:
            urllib.request.urlopen(url + "/api/v1/info", timeout=1).read()
            return proc, url
        except OSError:
            time.sleep(0.1)
    proc.kill()
    raise RuntimeError("nanovision serve did not start")


def get_json(url):
    with urllib.request.urlopen(url, timeout=30) as resp:
        return resp.status, json.loads(resp.read().decode("utf-8"))


def team_server_checks(binary, checks):
    proc, url = start_server(binary, TEAM_STORE, free_port(), token=TOKEN)
    try:
        store = ["-config", "nanovision.yaml", "-store", url]
        git("checkout", "--quiet", "main")

        code, out = nanovision(binary, *store, "-run-kind=submit")
        checks.expect("a CI run without the token is not uploaded", out, "set NANOVISION_TOKEN")

        code, out = nanovision(binary, *store, "-run-kind=submit", env={"NANOVISION_TOKEN": TOKEN})
        checks.expect("a CI run with the token is uploaded", out, "Open: " + url + "/runs/1")

        git("checkout", "--quiet", "feature/rounding")
        code, out = nanovision(binary, *store, "-run-kind=review", env={"NANOVISION_TOKEN": TOKEN})
        checks.expect("a review run finds its base run on the server", out, "(exact)", "Open: " + url + "/runs/2")

        code, out = nanovision(binary, *store)
        checks.expect("a local run reads the base but does not upload", out, "(exact)", "are not uploaded")
        git("checkout", "--quiet", "main")

        status, summary = get_json(url + "/api/v1/runs/2/summary")
        checks.add("the server builds the run summary", status == 200 and len(summary.get("nodes", [])) > 0,
                   f"{len(summary.get('nodes', []))} nodes")
        cmp = summary.get("comparison") or {}
        checks.add("the summary carries the delta against the base run",
                   cmp.get("headline") == "statement_coverage" and cmp.get("base", {}).get("id") == 1,
                   f"headline {cmp.get('headline')}, base {cmp.get('base', {}).get('id')}, "
                   f"{len(cmp.get('files') or [])} files changed")
        checks.add("the summary carries the verdict on the changed code", bool(summary.get("review")), "review block")
    finally:
        proc.terminate()
        proc.wait(timeout=10)

    code, out = nanovision(binary, "-config", "nanovision.yaml", "-store", url, "-run-kind=submit")
    checks.expect("a run fails soft when the server is down", out if code == 0 else f"exit {code}", "History is off")


def run_history_workflow(binary):
    print("\n" + "=" * 80)
    print("--- Starting run history checks ---")
    print("=" * 80)
    checks = Checks()
    try:
        setup_repo()
        local_store_checks(binary, checks)
        team_server_checks(binary, checks)
    except Exception as e:  # a crash in the setup is a failed check, not a traceback
        checks.add("the history workflow ran to the end", False, repr(e))
    return checks.results


def record_self_delta(binary, diff_file, checks=None):
    """
    Records nanovision's own unit coverage twice into the store: once at the
    merge-base with main (the base run, measured in a git worktree) and once
    for the working tree with its diff (a review run). The store then shows
    the delta, the compare page and the review block on real code.
    """
    checks = checks or Checks()
    worktree = os.path.join(WORK, "base-worktree")
    out_dir = os.path.join(WORK, "self-out")

    def repo_git(*args):
        return subprocess.run(["git", *args], cwd=SCRIPT_ROOT, check=True, capture_output=True, text=True).stdout.strip()

    try:
        base = repo_git("merge-base", "HEAD", "main")
        branch = repo_git("rev-parse", "--abbrev-ref", "HEAD")
        repo_git("worktree", "prune")
        remove_tree(worktree)
        repo_git("worktree", "add", "--detach", worktree, base)
        try:
            base_cov = os.path.join(WORK, "base-unit.out")
            print(f"--- Running the unit tests of the base revision {base[:10]} (this takes a while) ---")
            proc = subprocess.run(["go", "test", f"-coverprofile={base_cov}", "./..."], cwd=worktree,
                                  stdout=subprocess.DEVNULL, stderr=subprocess.PIPE, text=True)
            if not os.path.exists(base_cov):
                checks.add("the base revision's tests ran", False, proc.stderr[-2000:])
                return checks.results

            code, out = nanovision(binary, "-config", os.path.join(worktree, "nanovision.yaml"),
                                   "-report=" + base_cov, "-sourcedirs=.", "-reporttypes=TextSummary",
                                   "-output=" + out_dir, "-verbosity=Warning", "-store", STORE,
                                   "-project", "nanovision", "-profile", "unit", "-run-kind=submit",
                                   "-revision=" + base, "-stream", "main", cwd=worktree)
            checks.add("nanovision's base revision is saved as a submit run", code == 0 and "run " in out, out)
        finally:
            subprocess.run(["git", "worktree", "remove", "--force", worktree], cwd=SCRIPT_ROOT,
                           stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)

        current_cov = os.path.join(SCRIPT_ROOT, "reports", "nanovision_self_coverage", "coverage-unit.out")
        args = ["-config", os.path.join(SCRIPT_ROOT, "nanovision.yaml"), "-report=" + current_cov, "-sourcedirs=.",
                "-reporttypes=TextSummary", "-output=" + out_dir, "-verbosity=Warning", "-store", STORE,
                "-project", "nanovision", "-profile", "unit", "-run-kind=review", "-base-revision=" + base,
                "-stream", branch]
        if diff_file:
            args.append("-diff=" + diff_file)
        code, out = nanovision(binary, *args, cwd=SCRIPT_ROOT)
        checks.expect("nanovision's working tree is compared with main", out,
                      "Coverage delta against the base run", "(exact)")
        print(out)
    except subprocess.CalledProcessError as e:
        checks.add("nanovision's delta against main was recorded", False, f"{e}: {e.stderr}")
    return checks.results


def serve_store(binary, port=7070):
    """Serves the store of the checks until Ctrl+C, to look at it in a browser."""
    if not os.path.isdir(STORE):
        print(f"--- No store at {STORE}; run with --history first ---", file=sys.stderr)
        return
    print(f"\n--- Serving {STORE} at http://localhost:{port} (Ctrl+C to stop) ---")
    try:
        subprocess.run([binary, "serve", "-store", STORE, "-addr", f"localhost:{port}"],
                       env={**os.environ, "NANOVISION_TOKEN": TOKEN})
    except KeyboardInterrupt:
        pass
