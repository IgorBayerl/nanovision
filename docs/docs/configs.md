# Configuration & CLI

This page explains **how to configure nanovision** using a `nanovision.yaml` file **or** the **command‑line flags**. 
Both do the same thing. Use whichever you prefer. We'll show side‑by‑side mappings and plenty of copy‑pasteable examples.

---

## TLDR

* Put a `nanovision.yaml` next to your project and run `nanovision`.
* Or skip the file and run `nanovision -report=... -sourcedirs=... -reporttypes=Html -output=...`.
* **CLI flags override** settings from the config file.
* Reports supported as input: **Cobertura**, **Go Cover**, **GCov**.
* Output formats: **Html**, **TextSummary**, **Lcov**, **RawJson**.

---

## Config file vs CLI

nanovision reads settings from these two places (in this order):

1. **Config file** `nanovision.yaml` (optional)
2. **CLI flags** (optional) - these **override** the config file

This lets you keep a shared base config in the repo and tweak a few values on the command line when needed (e.g., change the output folder or add an extra report type in CI).

---

## Key concepts (2 minutes)

- A **report** is an input coverage file (or glob pattern) you already generated with your test tool (Cobertura/XML, Go's `coverage.out`, GCov `.gcov`, etc.).
- A **source dir** is the matching source code directory for that report. **Order matters** when you pass multiple reports: the *nth* report matches the *nth* source dir.
- **Report types** are the output formats you want nanovision to generate.

---

## CLI reference (flags)

Short list of the most useful flags. All flags can be combined. Paths support globs where noted.

- `-report=PATH[;PATH2;...]` - One or more coverage inputs (Cobertura XML, Go `coverage.out`, GCov `*.gcov`). Use `;` to separate.
- `-sourcedirs=DIR[;DIR2;...]` - One source directory per report, in the same order.
- `-output=DIR` - Where to write reports (folder is created if missing).
- `-reporttypes=Html,TextSummary,Lcov,RawJson` - Comma‑separated output formats.
- `-title=STRING` - Title for the HTML report.
- `-filefilters="RULE[;RULE2;...]"` - Include/exclude patterns. Each rule starts with `+` (include) or `-` (exclude); wildcards `*`/`?` allowed.
- `-verbosity=Verbose|Info|Warning|Error|Off` - Log level. (`-verbose` is a shortcut for chatty output.)
- `-logfile=FILE` - Also write logs to a file.
- `-logformat=text|json` - Console log format.
- `-config=path/to/nanovision.yaml` - Point to a specific config file if not in CWD.

---

## YAML config example (full, commented)

Below is a single, comprehensive example (mirrors the one in the repo root). You can run `nanovision` with no flags and it will use this file; any CLI flags you pass will override these values.

```yaml
# ------------------------------------------------------------------
#  nanovision Configuration for the Self-Coverage Report
# ------------------------------------------------------------------
# This file defines the default settings for generating nanovision's own
# coverage report. CLI flags provided by the python script will override
# the settings in this file.

# A list of coverage report files to parse for the merged report.
reports:
  - "reports/nanovision_self_coverage/coverage-unit.out"
  - "reports/nanovision_self_coverage/coverage-integration.out"
  - "demo_projects/cpp/report/gcov/branch-probabilities/*.gcov"
  - "demo_projects/csharp/report/cobertura/cobertura.xml"
  - "demo_projects/go/report/gocover/coverage.out"

# A list of source code directories. The order must match the `reports` list.
source_dirs:
  - "."
  - "."
  - "demo_projects/cpp/project"
  - "demo_projects/csharp/project"
  - "demo_projects/go/project"

# The directory where the final self-coverage report will be saved.
output_dir: "reports/nanovision_self_coverage_full"

# The types of reports to generate.
report_types:
  - "Html"
  - "TextSummary"
  - "Lcov"
  - "RawJson"

# The title for the generated HTML report.
title: "nanovision Self-Coverage (Full Merged)"

# Logging verbosity for the self-coverage run.
verbosity: "Verbose"

# Optional settings for patch coverage analysis.
diff:
  file: ""            # Path to the .diff file (can also be set via -diff flag)
  strip: "auto"       # "auto" or an integer (0-6) to strip leading path components

# A list of glob patterns for files and directories to exclude from the report.
# This is crucial for ignoring generated files and the vendored tree-sitter grammars,
# which are not part of the core tool's logic.
ignore_files:
  - "tree-sitter/**"       # Exclude all downloaded tree-sitter grammars
  - "**/*_test.go"         # Exclude test files themselves from coverage metrics
  - "tools/**"             # Exclude helper tools
  - "vendor/**"            # Exclude vendored dependencies
```

---

Use filters to include/exclude files from the final report.

- Start each rule with `+` (include) or `-` (exclude)
- Wildcards supported: `*` (many chars), `?` (single char)
- Matching is **case‑insensitive**
- Example list:

```yaml
ignore_files:
  - "-vendor/**"        # exclude
  - "-**/*_test.go"     # exclude test files themselves
  - "+src/**"           # include all under src
```

On the CLI, use `;` between rules:

```bash
-filefilters="-vendor/**;-**/*_test.go;+src/**"
```

> **Rule of thumb:** if no `+` rules are provided, everything is considered included **unless** excluded by a `-` rule.

---

### 1) Single Go report -> HTML only

**CLI**

```bash
nanovision -report=coverage.out -sourcedirs=. -reporttypes=Html -output=reports/go_html
```

**Config file** (`nanovision.yaml`)

```yaml
reports:
  - "coverage.out"
source_dirs:
  - "."
report_types:
  - "Html"
output_dir: "reports/go_html"
```

Run:

```bash
nanovision
```

### 2) Single Cobertura report -> HTML + Text summary

**CLI**

```bash
nanovision -report=report/cobertura/cobertura.xml -sourcedirs=. -reporttypes=Html,TextSummary -output=reports/csharp
```

**Config file**

```yaml
reports:
  - "report/cobertura/cobertura.xml"
source_dirs:
  - "."
report_types: [Html, TextSummary]
output_dir: "reports/csharp"
```

### 3) Single C++ GCov (glob) -> HTML

**CLI**

```bash
nanovision -report="cpp/report/gcov/branch-probabilities/*.gcov" -sourcedirs=cpp/project -reporttypes=Html -output=reports/cpp
```

**Config file**

```yaml
reports:
  - "cpp/report/gcov/branch-probabilities/*.gcov"
source_dirs:
  - "cpp/project"
report_types: [Html]
output_dir: "reports/cpp"
```

---

### 4) Merge Cobertura (C#) + Cobertura (C++) -> HTML + Text

**CLI**

```bash
nanovision -report="csharp/.../cobertura.xml;cpp/.../cobertura.xml" -sourcedirs="csharp/project;cpp/project" -reporttypes=Html,TextSummary -output=reports/merged_all_cobertura
```

### 5) Merge C++ GCov + C++ Cobertura -> all outputs

**CLI**

```bash
nanovision -report="cpp/report/gcov/branch-probabilities/*.gcov;cpp/.../cobertura.xml" -sourcedirs="cpp/project;cpp/project" -reporttypes=Html,TextSummary,Lcov,RawJson -output=reports/merged_all_cpp
```

### 6) Merge mixed projects (C#, Go, C++) -> HTML + Text + Lcov + RawJson

**CLI**

```bash
nanovision \
  -report="csharp/report/cobertura/cobertura.xml;go/report/gocover/coverage.out;cpp/report/gcov/branch-probabilities/*.gcov" \
  -sourcedirs="csharp/project;go/project;cpp/project" \
  -reporttypes=Html,TextSummary,Lcov,RawJson \
  -output=reports/merged_all_projects_mixed
```

**Config file**

```yaml
reports:
  - "csharp/report/cobertura/cobertura.xml"
  - "go/report/gocover/coverage.out"
  - "cpp/report/gcov/branch-probabilities/*.gcov"
source_dirs:
  - "csharp/project"
  - "go/project"
  - "cpp/project"
report_types: [Html, TextSummary, Lcov, RawJson]
output_dir: "reports/merged_all_projects_mixed"
```

---

### Custom HTML title and tag

```bash
nanovision -title="Release 1.2 Coverage" -tag=build-1042 ...
```

### Logging

```bash
# Very chatty logs in the console and a file
nanovision -verbosity=Verbose -logformat=json -logfile=nanovision.log ...

# Quick shortcut for verbose (overridden by -verbosity if both given)
nanovision -verbose ...
```

### Point to a specific config file

```bash
nanovision -config=path/to/nanovision.yaml
```

---

### Single project (Go)

```yaml
reports:
  - "coverage.out"
source_dirs:
  - "."
report_types: [Html, TextSummary]
output_dir: "reports/go"
verbosity: "Info"
ignore_files:
  - "-vendor/**"
  - "-**/*_test.go"
```

### Mixed, merged projects

```yaml
reports:
  - "csharp/report/cobertura/cobertura.xml"
  - "go/report/gocover/coverage.out"
  - "cpp/report/gcov/branch-probabilities/*.gcov"
source_dirs:
  - "csharp/project"
  - "go/project"
  - "cpp/project"
report_types: [Html, TextSummary, Lcov, RawJson]
output_dir: "reports/all"
Title: "All Projects Coverage"
verbosity: "Verbose"
ignore_files:
  - "-vendor/**"
  - "-tree-sitter/**"
  - "-tools/**"
```

---

## FAQ

**Do I have to pass `-sourcedirs`?**
Yes. nanovision needs source roots to map coverage to files and to analyze code (complexity, function ranges). For merged runs, pass one source dir per report.

**What output formats should I pick?**
For humans: `Html` + `TextSummary`. For CI badges or downstream tools: also add `Lcov` or `RawJson`.

**Can I generate multiple outputs at once?**
Yes - list them in `report_types` (config) or comma‑separate in `-reporttypes` (CLI).

**Can I run with only a config file?**
Yes. Put `nanovision.yaml` in your repo and run just `nanovision`.

---
