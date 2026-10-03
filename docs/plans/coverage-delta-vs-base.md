# Plan: coverage delta against the base revision

Status: draft, 2026-09-23. Updated 2026-09-24: the run store moved to [run-store-and-server.md](run-store-and-server.md), and this plan now covers test-only changes.

## 1. Goal

When you run nanovision on a change, it tells you how much the change increases or decreases the total coverage. You do not run the tests two times by hand to get this number.

Today, nanovision shows patch coverage. Patch coverage is the coverage of the lines that the change adds or modifies. It does not show how the change moves the total coverage of the project.

Patch coverage also misses a common case: a change that only adds or edits tests. Coverage filters usually ignore test files. So the change has no measured lines, and the review says that nothing changed. But the new tests can cover a large part of the code. Only a comparison with the base run shows this.

With this feature, nanovision does these steps for you:

1. It asks Perforce or git for the base revision of your change.
2. It gets the base run from the run store.
3. It compares the base run with the current run and prints the delta.

Example: a change that edits code.

```text
Coverage delta against the base run
  Base run   //game/main@118432  (exact)
  Profile    unit-win64
  Change     4 files: 3 measured, 1 ignored by the coverage filters

  Metric               Base      Current   Delta
  Statement coverage   78.45%    78.91%    +0.46
  Branch coverage      67.24%    66.10%    -1.14
  Methods hit          83.53%    83.60%    +0.07

  Statements           3033/3866 -> 3101/3930  (+68 covered, +64 valid)
  Patch coverage       84.0%  (42 of 50 changed statements)

  Files with the largest delta
    src/net/Replication.cpp   61.0% -> 74.2%   +13.2   modified
    src/core/Inventory.cpp    88.0% -> 71.5%   -16.5   not in diff
```

The last line is important. The coverage of `Inventory.cpp` decreased, but the change does not touch this file. This often means that the change removed or disabled a test, or that a test is not stable.

Example: a change that only adds tests.

```text
Coverage delta against the base run
  Base run   //game/main@118432  (exact)
  Profile    unit-win64
  Change     3 files, all ignored by the coverage filters (tests)

  Metric               Base      Current   Delta
  Statement coverage   78.45%    79.62%    +1.16
  Branch coverage      67.24%    70.69%    +3.45
  Methods hit          83.53%    84.33%    +0.79

  Statements           3033/3866 -> 3078/3866  (+45 covered, +0 valid)
  Patch coverage       none: the change has no measured lines

  Files with the largest delta
    src/core/Inventory.cpp    61.0% -> 88.1%   +27.1   not in diff
    src/core/ItemStack.cpp    70.2% -> 79.4%    +9.2   not in diff
```

The valid count did not change, so the change added no product code. The two files are the code that the new tests reach.

## 2. Terms

This plan uses one name for each item.

| Term | Meaning |
|---|---|
| Run | One execution of nanovision on a set of coverage reports. |
| Run store | The place that keeps runs: a local SQLite file or the team server. See [run-store-and-server.md](run-store-and-server.md). |
| Run kind | `submit`, `review`, or `local`. See section 4.5 of the store plan. |
| Revision | One state of the code. In git, a commit hash. In Perforce, a stream and a changelist number. |
| Change | The code that you want to measure. In Perforce, a pending or shelved changelist. In git, a branch or a pull request. |
| Local edits | Files in the workspace that are different from the current revision. In git, modified tracked files. In Perforce, opened files. |
| Clean workspace | A workspace with no local edits. |
| Base revision | The revision that the change starts from. |
| Base run | The stored run that nanovision compares with. Only a clean run can be a base run. |
| Profile | A name that keeps runs of different test suites or platforms apart. |
| Delta | The current value minus the base value, in percentage points. |
| Measured file | A file in the coverage tree. The coverage filters keep it, and a coverage report has it. |
| Test-only change | A change with no measured files. Usually, it adds or edits only tests. |
| VCS adapter | The code that asks git or Perforce for revisions and diffs. |
| Submit build | A CI build of a submitted changelist or a merged commit on the main line. |
| Review build | A CI build of a change before submit or merge. |

## 3. How it works

Two kinds of runs share one run store:

- A submit build saves a `submit` run for its revision.
- A review build or a developer run finds the base revision, loads the base run, and prints the delta.

```text
Submit build (clean workspace)          Review build or developer run (local edits)
------------------------------          -------------------------------------------
run the tests                           run the tests
nanovision -run-kind=submit             nanovision
  revision: //game/main@118432            base revision: //game/main@118432
  upload the run ---------+               diff: the change against the base revision
                          |               load the base run <-----+
                          v                                       |
                 run store (team server) -------------------------+
                                          print the delta
```

The design supports two modes:

- **Shared mode, for a team.** Submit builds upload `submit` runs to the team server. Review builds upload `review` runs. Developers read from the server. Nobody runs the tests on the base revision by hand.
- **Local mode, for one developer.** The run store is a local file. When you run the tests on a clean workspace, nanovision keeps that run as a base run. After that, each run with local edits compares with it. This mode needs one run on the clean base revision.

## 4. Design decisions

Each decision keeps the feature small.

**D1. Store compact facts, and share the parts that did not change.**
The store plan describes the format. The comparison needs only the manifests of the two runs, which have the counts for each file. The store also keeps the line data in a compact form, so a later version can show which lines lost coverage.

**D2. The team server owns the shared store.**
At studio scale, a network folder is too slow, and SQLite is not safe on a network share. The store plan explains this.

**D3. Only a clean run can be a base run.**
A run with local edits does not match a revision. If nanovision compared with it, the numbers could be wrong. The store keeps review runs and runs with local edits for viewing only.

**D4. The key is the profile and the revision.**
Two runs of the same revision with different tests are not comparable. The profile keeps them apart, for example `unit-win64` and `integration-linux`.

**D5. Compare counts, not stored percentages.**
nanovision calculates each percentage from the counts when it compares. Before it sums the base counts, it applies the current file filters to the base run. A new `ignore_files` entry then removes the same files from both sides.

**D6. Use the nearest older run when the exact run is missing.**
CI often skips revisions. For example, a build queue builds only the newest changelist. So the exact base run is often missing. nanovision then uses the nearest older run, up to `max_distance` revisions back. The output shows how far back it is.

**D7. Flags win over the VCS.**
CI systems already know the revision. Jenkins has `P4_CHANGELIST`, TeamCity has `BUILD_VCS_NUMBER`, and GitHub Actions has `GITHUB_SHA`. The flags `-revision` and `-base-revision` set the values directly. This makes Phase 1 useful before the VCS adapters exist.

**D8. Fail soft.**
If the run store, the team server, or the VCS is not available, nanovision logs a warning. It then makes the reports without the delta. The delta does not change the exit code. Only the optional gate in Phase 3 can do that.

**D9. The VCS integration is off by default.**
`vcs.type: none` keeps the current behavior. A Perforce server can be slow to answer, so you turn the integration on in the config.

## 5. What this feature needs from the run store

The store plan describes the run store. This feature uses these parts of it:

- The run kinds, and the rule that only clean runs can be base runs.
- A search for runs by project, stream, profile, and changelist number (`GET /api/v1/runs?...&near=`).
- The manifest of a run, with the counts for each file.
- The revision format: the commit hash for git, and `//stream@changelist` for Perforce.

The Perforce key includes the stream because changelist numbers are global to the server. The main stream and a release stream can both have runs near the same number. Without the stream, nanovision could compare with a run of the wrong stream. Git does not need this, because a commit hash is unique.

## 6. Find the base revision

### 6.1 Rules for git and Perforce

**The base revision:**

- Git: the merge-base of `HEAD` and the base branch. For a branch or a pull request, this is the commit where the branch starts, or the last merge from the base branch.
- Perforce: the synced changelist of the workspace. Pending and shelved changelists sit on top of it.

**Base run candidates:** only clean runs. On the team server, these are the `submit` runs. In a local store, these are the clean `local` runs.

**The base run:**

1. The search for the base run starts at the base revision.
2. Exception: if the workspace is clean and the base revision is the current revision, the search starts at the revision before it. A run never compares with its own revision. This is the normal case for a submit build. The delta then shows the effect of the changelists or commits since the previous run.
3. If the first revision of the search has a run, that run is the base run. The distance is 0.
4. If not, nanovision searches older revisions, up to `max_distance`. The output shows the distance, for example `3 changelists before the base revision`.
5. If nanovision finds no base run, it prints one line and continues without the delta.

Perforce adds one more exact case. Section 6.2 describes it.

**Flags:** `-revision` sets the current revision. `-base-revision` sets the base revision. A flag always wins over the VCS answer.

### 6.2 Perforce

| Question | Command | Use |
|---|---|---|
| Is this a Perforce workspace? | `p4 -ztag info` | `clientName` has a value, and `clientRoot` contains the project root. |
| Which stream? | `p4 -ztag client -o` | The `Stream` field. It is empty in a classic depot. |
| Is the stream virtual? | `p4 -ztag stream -o <stream>` | If `Type` is `virtual`, use the `Parent` stream in the key. |
| Which changelist does the workspace have? | `p4 -ztag changes -m1 //<client>/...#have` | The base revision. |
| Are there local edits? | `p4 -ztag opened` | Any opened file is a local edit. |
| What does the change contain? | `p4 diff -du` | The diff of all opened files against the synced revisions. |
| How many changelists are between two numbers? | `p4 -ztag changes -s submitted //<client>/...@<low+1>,@<high>` | Only the changelists that changed files in the workspace view. |

Notes:

- The diff includes all local edits of the workspace, not only one pending changelist. The test run also includes all of them, so the two numbers agree.
- Phase 0 checks if `p4 diff -du` includes files opened for add. If it does not, the adapter reads each added file and marks every line as added.
- A review build unshelves the change into its workspace with `p4 unshelve -s <changelist>`. It must do this to build the code anyway. Perforce then shows the files as opened, so the normal rules apply.
- `p4 diff -du` writes local absolute paths in the `+++` lines. The adapter makes these paths relative to the project root. The diff resolver then finds exact matches and does not need its path guesses.
- A virtual stream has no changelists of its own. Its content comes from the parent stream, so the key uses the parent.
- If the workspace has no stream (a classic depot), the key is only `@<changelist>`. Use one profile for each codeline in this case.
- The adapter uses the Perforce environment of the user: `P4PORT`, `P4USER`, `P4CLIENT`, `P4CONFIG`, and the login ticket. If the ticket has expired, the error message tells the user to run `p4 login`.
- Each command has a timeout of 30 seconds.

**The base run in Perforce:**

The `#have` query gives the newest changelist that changed a file in your workspace view. A submit build can save its run at a higher changelist that changed only files outside your view. That run has the same code as your base revision, so it is exact for you.

1. Take the stored changelists of the same stream.
2. Select the lowest stored changelist that is equal to or higher than the base changelist, if one exists.
3. Count the changelists in the view between the base changelist and this changelist.
4. If the count is 0, this run is the base run. The distance is 0.
5. Otherwise, select the highest stored changelist below the base changelist.
6. Count the changelists in the view between this changelist and the base changelist. This count is the distance.
7. If the distance is higher than `max_distance`, do not use the run.

In the exception of 6.1, skip steps 2 to 4. The search then finds only older runs.

Steps 2 to 4 are important in game studios. There, most changelists change art outside the code view of a programmer. Without these steps, the base run of a programmer is almost never exact.

### 6.3 Git

| Question | Command | Use |
|---|---|---|
| Is this a git repository? | `git rev-parse --show-toplevel` | The project root is inside the repository. |
| Which commit? | `git rev-parse HEAD` | The current revision. |
| Which branch? | `git rev-parse --abbrev-ref HEAD` | For the output only. |
| Are there local edits? | `git status --porcelain --untracked-files=no` | Any output is a local edit. |
| Which base? | `git merge-base HEAD <base_branch>` | The base revision. |
| What does the change contain? | `git diff --no-color --no-ext-diff --src-prefix=a/ --dst-prefix=b/ <base>` | The working tree against the base revision: the commits and the uncommitted edits. |
| Which new files does git not track yet? | `git ls-files --others --exclude-standard` | Every line of these files counts as added. |
| Which older commits can hold the base run? | `git rev-list --first-parent --max-count=<max_distance> <base>` | The first commit in the list with a run is the base run. |

Notes:

- The default base branch is the target of `origin/HEAD`. If it does not exist, the adapter tries `origin/main` and then `origin/master`. Set `vcs.base_branch` to use another branch.
- The diff flags make the output the same for all user git configs, for example `diff.noprefix` or an external diff tool.
- A shallow clone breaks `git merge-base`. In CI, fetch the full history. In GitHub Actions, set `fetch-depth: 0`.
- A pull request build on GitHub checks out a merge commit. The merge-base of that commit and the base branch is the tip of the base branch. So the delta is the delta of the pull request.
- The clean check ignores untracked files, because build outputs are often untracked. The limit: an untracked source file in a clean workspace makes the saved run a little wrong.

## 7. The comparison

### 7.1 What nanovision calculates

1. The headline metric is statement coverage when both runs have statements. Otherwise it is line coverage.2. For each coverage metric in `file_metrics` that both runs have, nanovision shows the base value, the current value, the delta, and the counts.
3. For each file, nanovision calculates the delta of the headline metric. The output lists the 5 files with the largest delta.
4. Each listed file has a label: `added`, `modified`, or `not in diff`. The label comes from the diff. Without a diff, the output has no label.
5. The block also shows the patch coverage. nanovision calculates these numbers today. The block only repeats them.

### 7.2 Changed files and test-only changes

nanovision puts each file of the diff into one of three groups:

1. **Measured.** The file is in the coverage tree. Patch coverage uses it.
2. **Ignored by the filters.** `ignore_files` or `file_filters` removes it. Test files are usually in this group.
3. **Not in the reports.** No coverage report has the file, for example a build script or a data file.

The `Change` line of the output shows these counts. Files in groups 2 and 3 no longer cause the warning `unable to map diff path`. They get a debug log only.

When group 1 is empty, the change is a test-only change:

- The `Change` line says so.
- Patch coverage shows `none`.
- The delta is the main result. Files with a higher coverage are the files that the new tests reach.
- The patch gate does not apply. The gate `max_coverage_decrease` still applies, so a change that deletes tests can fail it.

To put a diff file into group 2, nanovision applies the filters to the path of the file, relative to the project root. The VCS adapters give these relative paths. For a diff file from `-diff`, nanovision uses the path that the resolver finds.

### 7.3 Fairness warnings

The delta is correct only when both runs used the same tests. nanovision cannot prove this, but it can see the signs of a problem. It prints a warning when:

- The report patterns of the two runs are different.
- More than 10% of the coverable lines of the base run are in files that the current run does not have. This usually means that fewer tests ran.
- The major or minor version of nanovision is different between the two runs. Metric rules can change between versions.

A warning does not change the exit code. The distance to the base revision is not a warning, because the header line already shows it.

### 7.4 Where the delta shows

- Phase 1: the terminal and `Summary.txt`. The terminal block goes to standard output. It uses ASCII characters only, so old Windows consoles and CI logs show it correctly.
- Phase 3: the review comment, the `HtmlReview` report, and the compare page of the team server.

## 8. Configuration

```yaml
history:
  store: ""              # the team server URL, or a folder for a local store. Empty (the default) turns history off
  project: ""            # the project name on the team server
  profile: "default"     # keeps runs of different test suites apart
  max_distance: 50       # how many revisions back to search for a base run

vcs:
  type: "none"           # none (default), auto, git or perforce
  base_branch: ""        # git only. Empty means origin/HEAD, origin/main, origin/master
```

| Flag | Config key | Meaning |
|---|---|---|
| `-store` | `history.store` | The team server URL, or a local store folder. |
| `-project` | `history.project` | The project name. |
| `-profile` | `history.profile` | The profile name. |
| `-run-kind=submit\|review\|local` | (flag only) | The kind of this run. The default is `local`. |
| `-vcs` | `vcs.type` | `none`, `auto`, `git`, or `perforce`. |
| `-revision` | (flag only) | The revision of this run. Use it only when the workspace matches the revision exactly. |
| `-base-revision` | (flag only) | The revision to compare with. |

Rules:

- A `local` run never goes to the team server. It reads base runs from the server.
- A `submit` run needs a clean workspace. Without the VCS adapter, it also needs `-revision`.
- A `review` run goes to the team server. With a local store, it stays in the local store.
- The flags use the same revision format as the adapters: the commit hash for git, and `//stream@changelist` for Perforce.

**Example: local mode for one developer.**

```yaml
history:
  store: ".nanovision"
vcs:
  type: "auto"
```

Add `.nanovision/` to `.gitignore` or `.p4ignore`.

**Example: shared mode for a team.**

```yaml
history:
  store: "http://nanovision.studio.local:7070"
  project: "game"
  profile: "unit-win64"
vcs:
  type: "perforce"
```

## 9. Code changes

The store plan adds the store, the server, and the upload. This feature adds these parts.

### 9.1 New packages

| Package | Content |
|---|---|
| `internal/compare` | `Compare`, which calculates the delta between two manifests. |
| `internal/vcs` | The `VCS` interface, `Detect`, a command runner with a timeout, `git.go`, and `perforce.go`. |

### 9.2 Changed files

| File | Change |
|---|---|
| `internal/config/config.go` | `HistoryConfig`, `VCSConfig`, the new flags in `RawConfigInput`, the defaults, and a check of the profile name. |
| `internal/model/comparison.go` (new) | A plain data type for the result, and a `Comparison` field on `SummaryTree`. Reporters read it, as they read the diff data today. |
| `cmd/main.go` | The new flags. Split `executePipeline`, so the REPORT stage runs after the comparison. |
| `cmd/history.go` (new) | Before the pipeline: detect the VCS and make the diff. After the pipeline: find the base run, compare, save or upload the run, and print. |
| `internal/diff/parser.go` | `ParseText`, so the adapter can give diff text to the parser without a temporary file. `Parse` calls it. |
| `internal/diffapply/apply.go` | Put the diff files into the three groups of section 7.2. |
| `internal/review/review.go` | Add the group counts to the review statistics. |
| `internal/reporter/textsummary/reporter.go` | Write the comparison block into `Summary.txt`. |

### 9.3 Main types

This is a sketch. The names can change during the work.

```go
// internal/compare

func Compare(base, current *store.Manifest, include func(path string) bool,
	diffStatus map[string]string) *model.Comparison
```

```go
// internal/vcs

type Revision struct {
	ID     string // git: commit hash. Perforce: "//stream@changelist"
	Branch string // git branch or Perforce stream, for the output only
}

type VCS interface {
	Name() string
	Current() (Revision, error)
	HasLocalEdits() (bool, error)
	Base() (Revision, error)
	Diff(base Revision) (*diff.DiffData, error)
	// Nearest selects the stored revision to compare with, by the rules in section 6.
	// olderOnly is true in the exception of 6.1. A distance of 0 means an exact match.
	Nearest(rev Revision, stored []string, max int, olderOnly bool) (id string, distance int, ok bool)
}

func Detect(kind, projectRoot string) (VCS, error) // nil when kind is "none"
```

The adapters run commands through a small runner interface. Tests replace the runner with recorded output.

## 10. Build order

This work is milestone M4 of the store plan. Phase 0 is part of milestone M0. Phase 1 needs the local store from M2. Shared mode needs the team server from M3.

Each phase ends with a feature that works.

### Phase 0: check the Perforce commands

Goal: real command output for the tests, and no surprises in Phase 2.

1. On a real workspace, run each command in section 6.2.
2. Check if `p4 diff -du` includes files opened for add, move, and integrate.
3. Unshelve a shelved changelist, then run `p4 diff -du`.
4. Measure the time of `p4 changes -m1 //<client>/...#have` on the largest workspace that you have.
5. Save the outputs in `internal/vcs/testdata/p4/`.
6. Replace real server names, user names, and depot paths in these files. The repository is public.

Done when: the test files are in the repository, and this plan has the results.

### Phase 1: the delta, with revisions from flags

Goal: the full delta for CI systems that know their revisions. There are no VCS calls yet.

1. Add the config keys and the flags of section 8.
2. Create `internal/compare` and `internal/model/comparison.go`.
3. Split `executePipeline`, so the REPORT stage runs after the comparison.
4. Put the diff files into the three groups of section 7.2.
5. Add the group counts to the review statistics.
6. Print the comparison block at the end of the run, and write it into `Summary.txt`. Include the form for test-only changes.
7. Find the base run from the flags. Use an exact match first. For a Perforce revision, use the highest stored changelist of the same stream below it next. The output then says that the base run is older, without a distance.
8. Add the unit tests and the end-to-end cases.

Done when: a review build with `-base-revision` prints the delta, and a test-only change shows its delta.

### Phase 2: Perforce adapter

Goal: no revision flags and no diff file in a Perforce workspace.

1. Create `internal/vcs` with the `VCS` interface, `Detect`, and the command runner.
2. Write `perforce.go`: the stream, the synced changelist, the opened files, the diff, and the base run search of 6.2.
3. Add `ParseText` to `internal/diff`.
4. Make the diff paths relative to the project root.
5. Handle files opened for add, as the Phase 0 result shows.
6. Add `VCSConfig` and the `-vcs` flag.
7. If `diff.file` is empty and the adapter is active, make the diff with the adapter.
8. Add the tests with the Phase 0 files, and the optional test with a local `p4d`.

Done when: in a workspace with a pending changelist, `nanovision` prints the delta and the patch coverage. The user gives no diff file and no revision flags.

### Phase 3: review builds and the gate

Goal: the delta in the review, where people decide.

1. Upload review runs with `-run-kind=review`. The server answers with a link to the compare page.
2. Add the delta and the link to `Annotations.md`, so a CI job can post it to the review.
3. Show the delta in the `HtmlReview` report. For a test-only change, show the delta first.
4. Add `review.gate.max_coverage_decrease`. The gate fails when the headline metric decreases by more than this number of points.
5. Apply the gate only with an exact base run. With an older base run, print a warning and pass.

Done when: the review of a shelved changelist shows the delta and a link, and the gate works with `-fail-on`.

### Phase 4: git adapter

Goal: the same result for git branches and pull requests.

1. Write `git.go`: the commit, the branch, the status, the merge-base, the diff, the untracked files, and the base run search with `git rev-list`.
2. Add `vcs.base_branch` with the default search order.
3. Add the integration test with a temporary repository.
4. Use the feature in the CI of this repository, with the cache example in section 12.

Done when: a pull request build of this repository prints the delta against the merge-base.

### Phase 5: later

This plan does not describe these items in detail:

- Patch coverage for Perforce submit builds with `p4 diff2`.
- A base run from the parent stream for development streams and task streams.
- The lines that lost coverage, and a delta for each method. The store keeps the line data for this (store plan, M6).

## 11. Tests

Unit tests:

- `internal/compare`:
  - deltas for simple numbers, and valid counts of 0
  - a file only in the base run, and a file only in the current run
  - the current filter on the base run
  - the headline rule
  - each fairness warning
- `internal/diffapply`: the three groups of section 7.2, and no warning for ignored files.
- `internal/vcs`: each adapter with recorded command output from `testdata/`.

Integration tests:

- Git: the test makes a temporary repository with `git init`, adds commits and a branch, and checks `Current`, `Base`, `Diff`, and `Nearest`.
- Perforce: the test runs only when `p4` and `p4d` are on the `PATH`. It starts a temporary server with `P4PORT=rsh:p4d -r <temp dir> -i`, so no real server is necessary. Otherwise, the test calls `t.Skip`.

End-to-end tests in `scripts/e2e_test.py`:

- Run nanovision two times with a temporary local store. The first run is a `submit` run with `-revision=A`. The second run uses `-base-revision=A` and a different report. Check the delta in `Summary.txt`.
- A test-only change: the diff has only files that `ignore_files` removes. Check the `Change` line, `none` for patch coverage, and the delta.

## 12. CI examples

**Perforce with Jenkins or TeamCity (Phase 2).** The team config has `vcs.type: perforce` and the server URL in `history.store`.

Submit build, for each submitted changelist on `//game/main`:

```bat
rem The CI plugin syncs the workspace to the changelist.
run_tests.bat
nanovision -config nanovision.yaml -run-kind=submit
```

Review build, for each shelved changelist:

```bat
p4 sync //game/main/...
p4 unshelve -s %SHELVED_CL%
run_tests.bat
nanovision -config nanovision.yaml -run-kind=review
```

A developer runs the tests, then runs `nanovision`. No other step is necessary.

**Perforce before the adapter exists (Phase 1).** The build gives the revisions and the diff by hand.

```bat
rem Submit build
nanovision -config nanovision.yaml -run-kind=submit -revision=//game/main@%P4_CHANGELIST%

rem Review build. BASE_CL is the synced changelist of the workspace.
p4 diff -du > changes.diff
nanovision -config nanovision.yaml -run-kind=review -diff=changes.diff -base-revision=//game/main@%BASE_CL%
```

**Git with GitHub-hosted runners (Phase 4).** These runners have no team server, so a local store lives in the Actions cache. This works for small projects.

```yaml
- uses: actions/checkout@v4
  with:
    fetch-depth: 0
- uses: actions/cache@v4
  with:
    path: .nanovision
    key: nanovision-store-${{ github.sha }}
    restore-keys: nanovision-store-
- run: go test -coverprofile=coverage.out ./...
- run: nanovision -config nanovision.yaml -vcs=git -store=.nanovision -run-kind=${{ github.event_name == 'push' && 'submit' || 'review' }}
```

A pull request build restores the newest cache of the base branch. Two parallel builds on the main branch can each lose the run of the other. For a base run, this is acceptable.

## 13. Risks and limits

| Risk | Effect | Answer in this plan |
|---|---|---|
| The two runs used different tests. | The delta is wrong. | Fairness warnings (7.3). The profile keeps test suites apart. |
| Some tests are not stable. | Small deltas are noise. | Show the counts next to the percentages. Set the gate above 0. |
| A Perforce workspace has files from different changelists. | The base revision is not exact. | nanovision expects a full sync. The docs must say so. |
| The `#have` query is slow on a very large workspace. | The run takes more time. | `-base-revision` skips the query. CI already knows the number. |
| The team server is not reachable. | No base run, so no delta. | The run continues without the delta (D8). |
| CI does not build development streams or task streams. | These streams have no base run. | Phase 5: use the parent stream. |
| A diff path does not match a path in the tree. | A file can land in the wrong group of 7.2. | The VCS adapters give paths relative to the project root. |

## 14. Decisions to confirm

1. **Which adapter comes first?**
   Recommendation: Perforce in Phase 2, git in Phase 4. Perforce is the main use case. Until Phase 4, git users can use Phase 1 with `-revision`.
2. **What happens when only an older base run exists?**
   Recommendation: compare up to 50 revisions back, and show the distance. The gate uses only exact base runs.
3. **Which metric is the headline?**
   Recommendation: statement coverage when both runs have it, else line coverage.
4. **What does the review comment show first for a test-only change?**
   Recommendation: the delta, because patch coverage has no data.

## 15. Out of scope

- The store, the team server, and the history pages. The store plan covers them.
- Test results in runs.
- Everything in Phase 5, until a new plan describes it.
