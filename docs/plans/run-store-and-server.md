# Plan: run history at studio scale

Status: draft, 2026-09-24. Related plan: [coverage-delta-vs-base.md](coverage-delta-vs-base.md).

## 1. The problem

Today, each run writes finished reports: an HTML folder and, if you ask for it, `RawJson.json`. These are the sizes for the self-coverage run of this repository (88 files, 5,634 coverable lines):

| Output | Size | Size for each file |
|---|---|---|
| HTML folder | 8.2 MB | about 19 KB for each file page |
| `RawJson.json` | 3.7 MB | about 42 KB, or 687 bytes for each coverable line |

For a studio project with 30,000 files, the same sizes give about 0.6 GB of HTML and 1.3 GB of `RawJson.json` for each run. The files in this repository are small Go files. C++ files are larger, so the real numbers are probably higher.

This does not scale:

- With hundreds of runs each day, stored reports grow by hundreds of GB each day.
- A static HTML report with 30,000 pages is slow to write and slow to open.
- A team cannot share a folder of HTML files by link. Browsers do not open `file://` links from web pages.

History is still necessary. A change that only adds tests has no measured lines, so patch coverage cannot measure it. Only a comparison with the base run shows what the tests add. The delta plan describes this case.

## 2. The answer in short

1. **Store facts, not reports.** A report is a view. The server makes views when somebody asks for them.
2. **Split the data by how often it changes.** Source text and analysis change only when a file changes. Coverage changes in each run, but only for some files. nanovision stores each part once, and all runs that use it share it.
3. **Use a small binary format for the bulk data.** For the same coverage, it is about 200 times smaller than `RawJson.json` (measured in section 4.2).
4. **One process owns the store: `nanovision serve`.** It is the same binary. On a laptop, it serves one developer. On one machine in the studio, it serves the whole team. It does not sync, and it needs no accounts and no cloud.
5. **Make views on demand.** The browser gets small JSON answers from the server. Large projects do not write a static HTML report for each run.

## 3. The strategy: what changes and what stays

`docs/strategy/product-vision.md` plans a folder of runs first, then a shared folder, and a sync engine only when a shared folder becomes slow. Each run in that plan stores the model tree (`coverage.json`).

The first studio rollout gives real numbers, and they change two parts of that plan:

| Part of the strategy | Change | Reason |
|---|---|---|
| Each run stores the model tree. | Each run stores compact facts, and runs share the parts that did not change. | The model tree is about 1.3 GB for each run at 30,000 files. |
| Phase 3 uses one shared folder (an SMB share or a depot path). | One `nanovision serve` process owns the store. | Hundreds of writers and readers each day make a network folder slow. SQLite is not safe on a network share. A folder gives no links. |

These parts stay the same:

- Local first. The CLI works without a network, and a developer can keep a local store.
- Nothing leaves the studio. The team server runs inside the studio network.
- No sync engine, no CRDTs, and no accounts in the first version.
- One binary.
- A saved run never changes. Its key is the project, the stream, the profile, and the revision.

The strategy has one rule: `product-vision.md` and `plan.md` change only after `evidence.md` changes. So first add a dated entry to `evidence.md` with the studio facts: 30,000 files, C++ and Go, Perforce, hundreds of changelists each day, and a 30 MB report for a part of the project. Then change Phase 3 in `product-vision.md`.

## 4. The data

### 4.1 Five kinds of data

| Data | Changes when | Key | Shared between runs |
|---|---|---|---|
| Source text | the file changes | hash of the file content | yes |
| Analysis: methods, complexity, statements | the file or the analyzer changes | content hash and analyzer version | yes |
| File coverage: coverable lines, hits, branches, report masks | the tests or the code change | hash of the encoded coverage | yes, if the coverage of the file did not change |
| Manifest: path, content hash, coverage hash, and counts for each file | every run | hash of the manifest | no |
| Run metadata: revision, stream, profile, kind, author, time, CI link, totals, thresholds | every run | run ID | no |

nanovision does not store derived data: percentages, statuses, problems, review results, and HTML. The server calculates them when somebody asks, and keeps the result in a cache. The run metadata keeps a copy of the thresholds, so an old run shows the rules of its time.

The analysis is the same data that `internal/cache` keeps today. That cache already uses the content hash as its key.

### 4.2 Measured sizes

The same coverage data (88 files, 5,634 coverable lines, 5 reports) in three encodings:

| Encoding | Raw | Gzip for each file | Gzip for all files together |
|---|---|---|---|
| `RawJson.json` (today) | 3.7 MB | – | – |
| Compact JSON arrays | 36 KB | 8.8 KB | 3.6 KB |
| Binary varints | 17 KB | 5.9 KB | 3.1 KB |

The analysis of the same files is 28 KB raw in the binary format, or 8 KB after gzip. nanovision stores it once for each version of a file.

The structure gives the big gain, not the binary format. Compact JSON is already about 100 times smaller than `RawJson.json`. The binary format makes it 1.5 to 2 times smaller again.

### 4.3 Format choice

Use the binary format for the two bulk kinds: the manifest and the file coverage.

- The manifest has two hashes for each file. JSON must write them as text, which doubles their size. Gzip cannot make random hashes smaller.
- The server reads manifests for each summary and each comparison. A binary manifest decodes much faster than JSON. (Estimate: a few milliseconds for 30,000 entries.)
- The format is ours, so it is deterministic. The same data always gives the same bytes and the same hash. Deduplication needs this.

Do not use these formats:

- `encoding/gob`. It writes maps in random order, so the same data can give different bytes and a different hash. Also, only Go can read it.
- Protocol Buffers. It needs the `protoc` tool chain, and it does not promise the same bytes across versions. The gain over simple varints is small.

Rules for the binary format:

- Each blob starts with a kind byte and a version byte.
- Numbers are unsigned varints. A line number is the difference from the previous line number.
- A string has a length prefix.
- A debug command, `nanovision store dump <hash>`, prints a blob as JSON.

JSON stays for the run metadata, the threshold copy, and the API for the browser.

### 4.4 The manifest

The manifest has one entry for each file, sorted by path:

- the path, relative to the project root, with `/`
- the content hash (16 bytes)
- the coverage hash (16 bytes)
- the counts: covered and valid lines, statements, branches, and methods

A manifest is about 100 bytes for each file, so about 3 MB at 30,000 files. The upload always sends the full manifest. The server stores it in a smaller form (section 5.3).

### 4.5 Kinds of runs

| Kind | Made by | Workspace | Can be a base run | Kept |
|---|---|---|---|---|
| `submit` | CI, after a submit or a merge | clean | yes | always |
| `review` | CI, for a shelved changelist or a pull request | local edits | no | 30 days |
| `local` | a developer | clean or local edits | only in a local store, and only when clean | the last 20 |

Only a clean run can be a base run. A run with local edits does not match a revision.

### 4.6 Estimates at studio scale

These assumptions come from the numbers above. Milestone M0 checks them on real studio data.

- 30,000 files with 150 coverable lines each on average: 4.5 million coverable lines.
- 300 runs each day.
- In each run, the coverage changes for 1% of the files (300 files), and the content changes for 20 files.

| Item | First run | Each later run |
|---|---|---|
| Source text | about 110 MB after gzip | about 100 KB |
| Analysis | about 7 MB after gzip | about 15 KB |
| File coverage | about 14 MB | about 135 KB |
| Manifest, full | about 3 MB | about 3 MB |
| Manifest, stored as a difference | – | about 30 KB |

With full manifests, the store grows by about 1 GB each day. With manifest differences, it grows by about 90 MB each day, or about 33 GB each year. The server deletes review runs after 30 days, so the real growth is lower.

The largest unknown is the hit count. Hit counts can change between runs in files whose code did not change, for example because of threads. Then the file coverage part grows. M0 measures this. If it is a problem, split each coverage blob into two blobs. One blob keeps the covered and not-covered state forever. The other blob keeps the exact hit counts for 30 days.

## 5. The store

### 5.1 One SQLite file

The store is one SQLite file. It holds the run table and all blobs.

- SQLite reads small blobs faster than separate files.
- It has indexes and transactions.
- One file is simple to back up and to move.
- The build already uses cgo for tree-sitter, so a cgo SQLite driver adds no new tool.

```sql
create table runs (
  id          integer primary key,
  project     text    not null,
  stream      text    not null,  -- "//game/main", or a git branch
  profile     text    not null,
  revision    text    not null,  -- "//game/main@118432", or a commit hash
  change      integer,           -- 118432, or null for git
  kind        text    not null,  -- submit, review, or local
  clean       integer not null,  -- 1 when the run can be a base run
  review_id   text,              -- the shelved changelist or the pull request
  author      text,
  ci_url      text,
  created_at  integer not null,
  tool        text    not null,  -- the nanovision version
  manifest    blob    not null,  -- the hash of the manifest blob
  meta        text    not null   -- JSON: totals, thresholds, report patterns
);
create index runs_by_change on runs (project, stream, profile, change);

create table blobs (
  hash  blob    primary key,     -- the first 16 bytes of SHA-256
  kind  integer not null,        -- source, analysis, coverage, or manifest
  data  blob    not null
) without rowid;
```

### 5.2 One owner

Only one process writes to a store: the CLI for a local store, or `nanovision serve` for a team store. Do not put the SQLite file on a network share. SQLite locks are not reliable there, and the file can become corrupt.

### 5.3 Manifest differences

1. The server stores a full manifest from time to time. This is a snapshot.
2. For every other run, the server stores only the entries that are different from the last snapshot.
3. To read a run, the server reads two blobs: the snapshot and the difference.
4. When a difference grows over 10% of the snapshot size, the server writes a new snapshot.

The upload format does not change, so clients do not see this.

### 5.4 Cleanup

- Delete review runs after 30 days.
- Keep the last 20 local runs in a local store.
- Keep all submit runs.
- Once each week, delete the blobs that no run uses.

All numbers are settings.

### 5.5 Backups

Each night, `nanovision serve` writes a copy of the store with `VACUUM INTO`. Keep the copies on another disk.

## 6. The server: `nanovision serve`

### 6.1 Two modes

- **Local.** A developer starts `nanovision serve` on `localhost`. It shows the runs of the local store. It also keeps a cache of the team server. Blobs never change, so the cache never becomes wrong.
- **Team.** One studio machine runs `nanovision serve --store D:\nanovision --addr :7070`. CI uploads to it. Developers, the CLI, and review comments link to it.

### 6.2 API

| Method | Path | Use |
|---|---|---|
| `POST` | `/api/v1/blobs/missing` | The client sends hashes. The server answers with the hashes that it does not have. |
| `POST` | `/api/v1/blobs` | Upload the missing blobs as one gzip stream. |
| `POST` | `/api/v1/runs` | Register a run: the metadata and the manifest hash. The answer has the URL of the run. |
| `GET` | `/api/v1/runs?project=&stream=&profile=&near=` | Find runs, for example the candidates for a base run. |
| `GET` | `/api/v1/runs/{id}/summary` | The summary for the UI, in the current `summaryV1` shape. |
| `GET` | `/api/v1/runs/{id}/file?path=` | One file for the UI, in the current `detailsV1` shape. |
| `GET` | `/api/v1/compare?base=&head=` | The delta between two runs. |
| `GET` | `/api/v1/blobs/{hash}` | One blob, for the CLI and the local cache. |

Saved runs never change. So the server marks these answers as immutable, and the browser keeps them in its cache. A second visit needs no network.

### 6.3 Upload

1. The CLI makes the manifest and the blobs.
2. The CLI asks the server which blobs are missing.
3. The CLI uploads only the missing blobs. After the first run, these are mostly the blobs of the changed files.
4. The CLI registers the run.
5. The server answers with the URL of the run page, and the CLI prints it.

### 6.4 Web UI

The UI stays the React app in `ui/`. It gets a second data source: the API instead of `data.js`. The server builds the JSON shapes that the static report uses today (`summaryV1` and `detailsV1`) from the store. It uses the existing pipeline stages and builders for this. So most UI components do not change.

URLs are short and stable, for example:

- `/game/main`: the stream page, with the latest run, the trend, and recent changelists
- `/game/main/118432`: one run
- `/game/main/118432/src/core/Inventory.cpp`: one file in one run
- `/game/main/compare/118420/118432`: the delta between two runs

### 6.5 Access

- The first version has no accounts. Make the server reachable only inside the studio network.
- An upload needs a token. Only CI has the token.
- Perforce protections do not apply on the server. Some code, for example console platform code, can be under an NDA. Use the setting `source_exclude` to stop the upload of source text for these paths. The server then shows the coverage without the text.
- Add single sign-on later, if the studio needs it.

### 6.6 Code layout

| Package | Content |
|---|---|
| `internal/store/blob` | The binary formats: encode, decode, and dump. |
| `internal/store` | The SQLite store, cleanup, backups, and the local cache of a team server. |
| `internal/server` | The HTTP API and the UI files. |
| `cmd/main.go` | The `serve` subcommand. It is the first subcommand, because the CLI has only flags today. |
| `ui/src/lib/` | A data source for the API, next to the current data source for `data.js`. |

## 7. How it must feel

Linear won against Jira with speed, strong defaults, and a keyboard-first design. In this market, the heavy tools are SonarQube and a self-hosted Codecov: several services, a database server, and slow pages. nanovision must be the opposite: one binary, instant pages, and no configuration for developers.

| Moment | Target |
|---|---|
| Full nanovision run on the project, with a warm cache | less than 60 seconds (M0 checks this target) |
| The delta in the terminal, after the coverage run | less than 1 second more |
| Upload from CI, after the first run | less than 10 seconds |
| First open of a run with 30,000 files | less than 1 second on the studio network |
| Open the same run again | less than 100 ms |
| Open a file | less than 150 ms |
| Search for a file by name | results while you type |

Rules for the UI:

- Every view has a short, stable URL that people can paste into a review or a chat.
- No page waits behind a spinner. Show the layout at once, and load data on hover, before the click.
- Use the keyboard: `Ctrl+K` jumps to a file, a run, or a changelist. `j` and `k` move in lists.
- The home page has an opinion: your recent changes first, then the health of your streams.
- Developers get the server URL from the project config in the depot. They configure nothing.

## 8. Rollout in the first studio

The sizes are rough, for one developer.

### M0: measure (about 1 week)

1. Run the current nanovision on the full project. Record the time of each stage, the peak memory, and the output size.
2. Get 5 to 10 consecutive CI coverage outputs of one stream.
3. Run a prototype encoder on them. Record the blob sizes.
4. Record how many coverage blobs change between two runs. Do this once with exact hit counts and once with covered bits only.
5. Check the Perforce commands (Phase 0 of the delta plan).
6. Write the numbers into `evidence.md`. Keep studio data at the studio. Only numbers go into the repository.

Done when: the numbers are in `evidence.md` and in this plan.

### M1: a fast single run (1 to 2 weeks)

1. Fix the slowest stages that M0 found.
2. Print the terminal summary to standard output.
3. Warn when a static HTML report has more than 5,000 files, and point to `nanovision serve`.

Done when: a full-project run with a warm cache meets the target from M0.

### M2: the store and the local server (2 to 3 weeks)

1. Write the blob formats. Add fuzz tests for the decoders.
2. Write the SQLite store.
3. Save each run into the local store.
4. Write `nanovision serve` for the local mode: the API and the UI.
5. Add the API data source to the UI.
6. Measure the summary load for 30,000 files. If it misses the target, send the summary in columns (one array for each field), not as one object for each file.

Done when: a developer opens a 30,000-file run in the browser in less than 1 second.

### M3: the team server, pilot on one stream (2 to 3 weeks)

1. Add the upload API and the token.
2. Add the run kinds and the `-run-kind` flag.
3. Add manifest differences, cleanup, and backups.
4. Add the stream page: the latest run, the trend, and recent changelists.
5. Write the install guide: a Windows service, the store folder, and the backup folder.
6. Let the submit builds of one stream upload their runs.

Done when: the pilot stream uploads for two weeks, and its pages meet the targets in section 7.

### M4: the delta in reviews (2 to 3 weeks)

This milestone is the delta plan: the Perforce adapter, the base run search, the comparison, test-only changes, and the review comment with a link.

Done when: a programmer sees the delta of a shelved changelist in the review, and in the terminal before the submit.

### M5: all projects of the studio

1. Write a project config template and CI templates for the build system of the studio.
2. Write one page of documentation for developers.
3. Set the source exclusions for protected paths.
4. Watch the store size and the page times.

Done when: teams use nanovision without the help of the maintainer. These teams are the real usage signal for the strategy.

### M6: history views

- Compare any two runs.
- The coverage history of one file.
- Trends for each stream and each profile.
- The lines that lost coverage in a change.
- Lines whose coverage changes between runs without a code change (unstable tests).
- Test results in runs (the JUnit parser from the strategy).

Later: the git adapter, a shared analysis cache on the server, single sign-on, and the MCP server.

## 9. Risks

| Risk | Effect | Answer |
|---|---|---|
| Hit counts change in many unchanged files. | The store grows faster. | M0 measures it. Split the hit counts into a separate blob with a 30-day limit. |
| The browser is slow with 30,000 files. | Pages miss the targets. | M2 measures it. Send the summary in columns. |
| The server stops. | No links and no delta. | The CLI works without the server. The server starts again in seconds. Backups run each night. |
| Protected source (NDA) is on the server. | People see code that they must not see. | Source exclusions, access only inside the network, and single sign-on later. |
| Studio data goes into this public repository. | A leak. | Only numbers. No studio name, paths, or code in plans, test files, or issues. |
| The maintainer builds features as part of a job. | The employer can own that code (for example §69b UrhG in Germany). | Agree on ownership in writing before the work starts. |
| Coverage builds at the studio are slow. | Few runs, so base runs are older. | The delta plan uses the nearest run and shows the distance. |

## 10. Decisions to confirm

1. **A team server instead of a shared folder.** Recommendation: yes.
2. **Binary blobs for the manifest and the file coverage, JSON for the rest.** Recommendation: yes.
3. **SQLite as the store.** Recommendation: yes.
4. **Exact hit counts in the history.** Recommendation: decide after M0.
5. **Source text on the server, with exclusions.** Recommendation: yes. The other option is to read the text from Perforce for each request. That needs a Perforce service account and adds load on the Perforce server.
6. **The pilot machine.** Recommendation: a build-farm machine with a local SSD and a nightly backup to another disk.
