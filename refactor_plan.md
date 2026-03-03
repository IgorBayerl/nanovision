# NanoVision — Metrics Decoupling & Strategy Refactoring
### Architecture Plan

---

## Table of Contents

1. [Context & Background](#1-context--background)
2. [The Problem](#2-the-problem)
3. [Final Architecture Decisions](#3-final-architecture-decisions)
4. [Full Pipeline Walkthrough](#4-full-pipeline-walkthrough)
5. [Files to Create and Modify](#5-files-to-create-and-modify)
6. [Implementation Tasks](#6-implementation-tasks)
7. [Testing Strategy](#7-testing-strategy)
8. [How the System Works After Refactoring](#8-how-the-system-works-after-refactoring)

---

## 1. Context & Background

NanoVision is a CLI code-coverage analysis engine. It parses coverage reports from multiple formats (Cobertura, GCov, etc.), builds an in-memory directory/file/method tree, annotates every node with risk statuses (safe / warning / danger), and renders the results as an HTML report and a text summary.

The system works today. However, it has reached a design inflection point driven by two colliding forces:

**Force 1 — Growing metric set.** The codebase started with line coverage. It now has branch and statement coverage, and cyclomatic complexity is next. Each metric added so far required surgical changes in at least four files simultaneously.

**Force 2 — Metric scope mismatch.** Users want to control metrics at two distinct scopes: what appears in the file/directory tree, and what appears inside the method detail panel. These are different concerns and should be configured independently. There is currently no mechanism to express this.

The goal of this refactoring is to resolve both forces permanently by introducing the Strategy Pattern uniformly across the pipeline and splitting metric configuration into two explicitly scoped lists.

---

## 2. The Problem

### 2.1 Tight Coupling Across the Pipeline

Every metric today is handled by hardcoded conditional logic scattered across three layers of the pipeline. The specific violations are:

**`internal/config/config.go` — Monolithic metric list.**
`DisplayMetrics` is a single flat list with no concept of scope. There is no distinction between "show this in the file tree" and "show this in the method table". Everything goes into one `ActiveMetrics map[MetricKey]bool`.

**`internal/status/annotate.go` — Hardcoded annotator.**
`annotateNode` and `annotateFile` contain large `if ActiveMetrics[X]` chains that perform percentage calculations inline and assign statuses directly. Adding any new metric requires opening this file and writing another branch.

**`internal/reporter/htmlreact/builder.go` — Monolithic builder.**
`buildMetricsMap` and `buildTotals` repeat the same `if b.config.ActiveMetrics[...]` pattern over and over. Math, capability checks, and JSON struct allocation are all mixed together inside the same function.

**No boundary between domain and presentation.**
There is no separation between "is this metric in danger?" (a domain question) and "how do I render this metric as a JSON object?" (a presentation question). UI-specific structures leak into status calculation logic.

### 2.2 Open/Closed Principle Violated

The system is currently *closed for extension, open for modification*. This is the exact inverse of the Open/Closed Principle. Every new metric requires modifying the core traversal engine. No amount of adding files fixes this — the engine itself must be changed each time.

### 2.3 Incompatible Metric Types

The current architecture implicitly assumes all metrics are percentages where higher is better. Cyclomatic Complexity breaks this assumption completely: it is a raw integer, it has no upper bound, and lower is better. The threshold logic is inverted. There is no clean way to express this difference in the current design without adding more special-case branches.

---

## 3. Final Architecture Decisions

These are the decisions reached after design iteration. They are not proposals.

### 3.1 Config Split

Replace the single `display_metrics` list in `nanovision.yaml` with two explicitly scoped lists:

```yaml
# nanovision.yaml

# Metrics shown in the file and directory tree
file_metrics:
  - "line_coverage"
  - "branch_coverage"
  - "statement_coverage"

# Metrics shown inside the method/function detail panel
method_metrics:
  - "line_coverage"
  - "statement_coverage"
  - "max_cyclomatic_complexity"
```

In `config.go`, replace `ActiveMetrics map[MetricKey]bool` with two maps:

```go
ActiveFileMetrics   map[MetricKey]bool
ActiveMethodMetrics map[MetricKey]bool
```

### 3.2 Domain Layer — Evaluator Interface

A new `Evaluator` interface lives in `internal/status/evaluator.go`. It encapsulates everything the annotator engine needs to know about a single metric. No UI logic is present here.

```go
// internal/status/evaluator.go

type Evaluator interface {
    Key() config.MetricKey

    // IsApplicable checks capabilities derived from the parser.
    // e.g. BranchCoverageEvaluator returns false if caps.HasBranchCoverage == false
    IsApplicable(caps Capabilities) bool

    // Evaluate computes the metric value and applies threshold logic.
    // Returns the RiskLevel and whether the metric has data to display.
    Evaluate(metrics model.CoverageMetrics, band *config.Band) (RiskLevel, bool)
}
```

`ExtractUIData` was explicitly **not** included in this interface. Putting it here would violate SRP by pulling presentation concerns into the domain layer.

### 3.3 Domain Layer — Classifiers

Two classifier helpers in `internal/status/classifier.go` eliminate duplicated threshold logic across evaluators:

```go
// internal/status/classifier.go

// For coverage metrics — higher is better (line%, branch%, statement%)
// val < Min → Danger | Min ≤ val ≤ Max → Warning | val > Max → Safe
func ClassifyHigherIsBetter(val float64, band *config.Band) (RiskLevel, bool)

// For complexity/churn metrics — lower is better
// val < Min → Safe | Min ≤ val ≤ Max → Warning | val > Max → Danger
func ClassifyLowerIsBetter(val float64, band *config.Band) (RiskLevel, bool)
```

### 3.4 Domain Layer — Evaluator Registry

Individual evaluators live in `internal/status/evaluators/` — one file per metric. A `Registry` map is the single source of truth for all known evaluators.

```go
// internal/status/evaluators/registry.go

var Registry = map[config.MetricKey]Evaluator{
    config.LineCoverage:            LineCoverageEvaluator{},
    config.BranchCoverage:          BranchCoverageEvaluator{},
    config.StatementCoverage:       StatementCoverageEvaluator{},
    config.MaxCyclomaticComplexity: MaxComplexityEvaluator{},
}
```

Example evaluator implementations:

```go
// internal/status/evaluators/line_coverage.go

type LineCoverageEvaluator struct{}

func (e LineCoverageEvaluator) Key() config.MetricKey { return config.LineCoverage }
func (e LineCoverageEvaluator) IsApplicable(c Capabilities) bool { return true }
func (e LineCoverageEvaluator) Evaluate(m model.CoverageMetrics, b *config.Band) (RiskLevel, bool) {
    if m.LinesValid == 0 { return "", false }
    pct := utils.CalculatePercentage(m.LinesCovered, m.LinesValid, 2)
    return ClassifyHigherIsBetter(pct, b)
}
```

```go
// internal/status/evaluators/complexity.go

type MaxComplexityEvaluator struct{}

func (e MaxComplexityEvaluator) Key() config.MetricKey { return config.MaxCyclomaticComplexity }
func (e MaxComplexityEvaluator) IsApplicable(c Capabilities) bool { return true }
func (e MaxComplexityEvaluator) Evaluate(m model.CoverageMetrics, b *config.Band) (RiskLevel, bool) {
    if m.MaxCyclomaticComplexity == 0 { return "", false }
    return ClassifyLowerIsBetter(float64(m.MaxCyclomaticComplexity), b)
}
```

### 3.5 Presentation Layer — Provider Interface

A `FileMetricProvider` interface lives strictly inside the `htmlreact` package. The existing `MethodMetricProvider` is refactored to match the same pattern.

```go
// internal/reporter/htmlreact/providers.go

type FileMetricProvider interface {
    Key() config.MetricKey
    // Apply takes raw model data and writes the UI-specific struct into the metrics map.
    Apply(metrics model.CoverageMetrics, ui metricsMap)
}
```

### 3.6 Refactored Annotator

`annotate.go` becomes a pure loop engine with zero metric-specific logic:

```go
// For file/directory nodes
for _, key := range cfg.ActiveFileMetrics {
    evaluator, ok := evaluators.Registry[key]
    if !ok || !evaluator.IsApplicable(caps) { continue }
    level, hasData := evaluator.Evaluate(node.Metrics, cfg.Bands[key])
    if hasData { node.Statuses[key] = level }
}

// For method nodes
for _, key := range cfg.ActiveMethodMetrics {
    evaluator, ok := evaluators.Registry[key]
    if !ok || !evaluator.IsApplicable(caps) { continue }
    level, hasData := evaluator.Evaluate(method.Metrics, cfg.Bands[key])
    if hasData { method.Statuses[key] = level }
}
```

### 3.7 Model Expansion

Two additions to the data model are required:

- Add `MaxCyclomaticComplexity int` to `model.CoverageMetrics` — so the evaluator can access it at the file and directory level.
- Add `Statuses map[config.MetricKey]string` to `model.MethodMetrics` — so per-method statuses can be stored and rendered.

The aggregator must propagate `MaxCyclomaticComplexity` upward through the tree by taking the **max value** at each level: methods → file → directory.

---

## 4. Full Pipeline Walkthrough

End-to-end lifecycle of a `nanovision` CLI run under the new architecture:

| Stage | Status | What Changes |
|---|---|---|
| 1. Config Load | **CHANGED** | YAML parsed into `AppConfig`. `display_metrics` replaced by `file_metrics` and `method_metrics`. Two O(1) lookup maps built: `ActiveFileMetrics` and `ActiveMethodMetrics`. |
| 2. Parse | UNCHANGED | Parsers (Cobertura, GCov, etc.) emit flat `ParserResult` arrays with line/branch hits. |
| 3. Tree Build | UNCHANGED | `internal/tree/builder.go` assembles the directory/file tree. |
| 4. Enrich | UNCHANGED | AST analysis in `internal/enricher` populates method/function nodes. |
| 5. Aggregate | **SLIGHTLY CHANGED** | Aggregator now also computes `MaxCyclomaticComplexity`, propagating the max value up from methods → files → directories. |
| 6. Diff Apply | UNCHANGED | Diff tagging is independent of metrics. |
| 7. Annotate | **HEAVILY CHANGED** | Engine loops over Evaluator registry using `ActiveFileMetrics` and `ActiveMethodMetrics`. Zero inline math. Statuses attached to nodes and methods via their `Statuses` maps. |
| 8. Report (HTML) | **HEAVILY CHANGED** | Builder loops over Provider registry using the split metric configs. Delegates JSON struct creation to each provider. No metric-specific logic remains in `builder.go`. |
| 9. Report (Text) | **CHANGED** | Text reporter updated to iterate over `ActiveFileMetrics` instead of the removed `ActiveMetrics` map. |

---

## 5. Files to Create and Modify

### 5.1 New Files

| File | Purpose |
|---|---|
| `internal/status/evaluator.go` | Defines the `Evaluator` interface. |
| `internal/status/classifier.go` | `ClassifyHigherIsBetter` and `ClassifyLowerIsBetter` helpers. |
| `internal/status/evaluators/registry.go` | Registry map of all known evaluators. |
| `internal/status/evaluators/line_coverage.go` | `LineCoverageEvaluator` implementation. |
| `internal/status/evaluators/branch_coverage.go` | `BranchCoverageEvaluator` implementation. |
| `internal/status/evaluators/statement_coverage.go` | `StatementCoverageEvaluator` implementation. |
| `internal/status/evaluators/complexity.go` | `MaxComplexityEvaluator` implementation (uses `ClassifyLowerIsBetter`). |

### 5.2 Modified Files

| File | Change Summary |
|---|---|
| `nanovision.yaml` | Replace `display_metrics` with `file_metrics` and `method_metrics`. |
| `internal/config/config.go` | Replace `DisplayMetrics` / `ActiveMetrics` with `FileMetrics`, `MethodMetrics`, `ActiveFileMetrics`, `ActiveMethodMetrics`. |
| `internal/model/metrics.go` | Add `MaxCyclomaticComplexity int` to `CoverageMetrics`. |
| `internal/model/tree.go` | Add `Statuses map[config.MetricKey]string` to `MethodMetrics`. |
| `internal/aggregator/aggregator.go` | Bubble up `MaxCyclomaticComplexity` (max value) through the tree. |
| `internal/status/annotate.go` | Rewrite to loop over Evaluator registry. Remove all inline metric logic. |
| `internal/reporter/htmlreact/providers.go` | Add `FileMetricProvider` interface. Refactor `MethodMetricProvider` to align. Implement file-level providers. |
| `internal/reporter/htmlreact/builder.go` | Replace `buildMetricsMap` if-chains with provider registry loops. |
| `internal/reporter/textsummary/reporter.go` | Update to iterate over `ActiveFileMetrics` instead of the removed `ActiveMetrics`. |

---

## 6. Implementation Tasks

Tasks are sequential. Do not start a task until the previous one passes its Definition of Done.

---

### Task 01 — Configuration Split

**Context.**
This is the foundation that every subsequent task depends on. The goal is to split the single `display_metrics` config into two scoped lists and update the data model to accommodate complexity as a first-class metric. At the end of this task, the project must still compile and all existing tests must still pass. No behavior changes at runtime — the new config keys will simply feed into the old `ActiveMetrics` map as a temporary shim until Task 02 removes it. The purpose of doing this as a standalone task is to keep the diff reviewable and to confirm the config parsing layer is correct in isolation before the engine is touched.

**Files.**
- `nanovision.yaml`
- `internal/config/config.go`
- `internal/config/config_test.go`
- `internal/model/metrics.go`
- `internal/model/tree.go`
- `internal/aggregator/aggregator.go`
- `internal/aggregator/aggregator_test.go`

**Definition of Done.**
- [ ] `nanovision.yaml` contains `file_metrics` and `method_metrics` lists. `display_metrics` is removed.
- [ ] `AppConfig` exposes `ActiveFileMetrics map[MetricKey]bool` and `ActiveMethodMetrics map[MetricKey]bool`. `ActiveMetrics` is removed — zero references remain in the codebase.
- [ ] `model.CoverageMetrics` has a `MaxCyclomaticComplexity int` field.
- [ ] `model.MethodMetrics` has a `Statuses map[config.MetricKey]string` field.
- [ ] The aggregator correctly propagates `MaxCyclomaticComplexity` (taking the max value) from methods up to files and from files up to directories.
- [ ] New config tests verify that both `ActiveFileMetrics` and `ActiveMethodMetrics` are populated correctly when parsing valid YAML.
- [ ] New config tests verify that an unknown metric key in the YAML produces a clear, informative error.
- [ ] New aggregator tests verify the `MaxCyclomaticComplexity` propagation with a fixture tree containing methods with differing complexity values.
- [ ] `go build ./...` succeeds with zero errors.
- [ ] `go test ./...` passes — no regressions.

---

### Task 02 — Domain Layer: Evaluators

**Context.**
This task introduces the Strategy Pattern into the core status annotation pipeline. The `Evaluator` interface and its implementations are the heart of the entire refactoring. The key insight driving this design is that the annotator engine should be completely ignorant of what a "metric" is — it should only know how to loop and assign results. All metric-specific knowledge (how to calculate the value, which direction is good, what capabilities are required) lives exclusively in the individual evaluator structs.

The classifier helpers (`ClassifyHigherIsBetter` / `ClassifyLowerIsBetter`) are extracted as shared utilities specifically to handle the cyclomatic complexity case, where lower is better. Without these helpers, every coverage evaluator would duplicate the same `if/else` threshold logic, and the complexity evaluator would need special-case handling in the engine itself.

After this task, `annotate.go` must contain zero hardcoded metric keys. Any reviewer should be able to read it and understand the loop without knowing what metrics exist.

**Files.**
- `internal/status/evaluator.go` *(new)*
- `internal/status/classifier.go` *(new)*
- `internal/status/evaluators/registry.go` *(new)*
- `internal/status/evaluators/line_coverage.go` *(new)*
- `internal/status/evaluators/branch_coverage.go` *(new)*
- `internal/status/evaluators/statement_coverage.go` *(new)*
- `internal/status/evaluators/complexity.go` *(new)*
- `internal/status/annotate.go` *(rewrite)*

**Definition of Done.**
- [ ] `Evaluator` interface is defined with `Key()`, `IsApplicable()`, and `Evaluate()` — and nothing else.
- [ ] `ClassifyHigherIsBetter` and `ClassifyLowerIsBetter` are implemented. Each is covered by table-driven tests for: nil band, below-min, at-min, at-max, above-max edge cases.
- [ ] All four evaluators (`LineCoverage`, `BranchCoverage`, `StatementCoverage`, `MaxComplexity`) implement the `Evaluator` interface and are registered in `Registry`.
- [ ] `MaxComplexityEvaluator` uses `ClassifyLowerIsBetter`. A test asserts that a complexity value above `band.Max` returns `RiskDanger`.
- [ ] `BranchCoverageEvaluator.IsApplicable()` returns `false` when `caps.HasBranchCoverage == false`.
- [ ] Each evaluator has unit tests covering: the not-applicable path, the zero-data guard (`ok = false`), and each of the three threshold outcomes (safe, warning, danger).
- [ ] `annotate.go` contains zero hardcoded `config.MetricKey` references and zero inline percentage calculations. Its entire metric logic is a loop over the registry.
- [ ] Integration test: build a minimal `model.SummaryTree` in memory, run `Annotate()`, and assert that file nodes and method nodes receive the correct statuses for each metric.
- [ ] `go test ./internal/status/...` passes with ≥ 90% coverage on all new files.
- [ ] `go build ./...` succeeds with zero errors.
- [ ] `go test ./...` passes — no regressions.

---

### Task 03 — Presentation Layer: HTML Providers

**Context.**
This task applies the same Strategy Pattern to the reporting layer. The reason this is a separate task from Task 02 is that the domain layer and the presentation layer must remain completely independent. An `Evaluator` must never import anything from `htmlreact`, and a `FileMetricProvider` must never import anything from `internal/status`.

The current `builder.go` mixes three concerns in `buildMetricsMap`: it traverses config, it knows the internal structure of each metric's data (how to calculate percentages, how to format them), and it builds JSON structs. After this task, `buildMetricsMap` is a loop that knows nothing about individual metrics.

The existing `MethodMetricProvider` pattern in `providers.go` is already the right idea — this task extends it consistently to the file level and ensures both interfaces follow the same contract.

**Files.**
- `internal/reporter/htmlreact/providers.go` *(add `FileMetricProvider` interface + implementations)*
- `internal/reporter/htmlreact/builder.go` *(refactor `buildMetricsMap` and `buildTotals`)*

**Definition of Done.**
- [ ] `FileMetricProvider` interface is defined with `Key()` and `Apply()` — mirroring the `MethodMetricProvider` pattern.
- [ ] A provider is implemented for each metric in the default `file_metrics` config: `LineCoverageProvider`, `BranchCoverageProvider`, `StatementCoverageProvider`.
- [ ] `builder.go`'s `buildMetricsMap` contains zero hardcoded metric keys. It is a loop over the provider registry filtered by `cfg.ActiveFileMetrics`.
- [ ] `buildTotals` follows the same pattern.
- [ ] The existing `MethodMetricProvider` implementations are refactored to align with the updated interface signature (no functional change, just consistency).
- [ ] A golden-file test confirms that the HTML JSON output for a fixed fixture tree is structurally equivalent to the pre-refactoring output.
- [ ] A test confirms that removing `branch_coverage` from `file_metrics` in the config removes branch data from the rendered file detail view.
- [ ] `go test ./internal/reporter/htmlreact/...` passes.
- [ ] `go build ./...` succeeds with zero errors.
- [ ] `go test ./...` passes — no regressions.

---

### Task 04 — Text Reporter Update

**Context.**
This is the final integration task. The text summary reporter is a simpler rendering surface than the HTML builder, but it still needs to be updated to consume the new split config. Because it is the last piece of the pipeline to be updated, this task also includes the final end-to-end smoke test against real fixture data.

The scope is intentionally small: the reporter should iterate over `ActiveFileMetrics` to decide which metric columns to print, rather than using the removed `ActiveMetrics` map. No new patterns are introduced here — it is a straightforward adaptation of existing logic to the new config structure.

**Files.**
- `internal/reporter/textsummary/reporter.go`
- `internal/reporter/textsummary/reporter_test.go`

**Definition of Done.**
- [ ] `textsummary` reporter references `ActiveFileMetrics` — `ActiveMetrics` is gone from this file.
- [ ] CLI output shows exactly the metrics listed in `file_metrics` — no more, no less.
- [ ] A test confirms that removing `statement_coverage` from `file_metrics` removes it from text output.
- [ ] A test confirms that a config with only `line_coverage` in `file_metrics` produces a single metric column in the output.
- [ ] `go test ./internal/reporter/textsummary/...` passes.
- [ ] End-to-end smoke test: run `nanovision` against the sample Cobertura fixture in the repository. Assert exit code 0, HTML output file exists and contains expected metric section IDs, and text output contains the configured metric names.
- [ ] `go build ./...` succeeds with zero errors.
- [ ] `go test ./...` passes — no regressions.

---

## 7. Testing Strategy

### 7.1 Unit Tests — Classifiers

`ClassifyHigherIsBetter` and `ClassifyLowerIsBetter` are pure functions with no dependencies. Use table-driven tests covering all boundary conditions:

```go
// Example: ClassifyHigherIsBetter, band = {Min: 50, Max: 80}
{ val: 49.9,  expected: RiskDanger  }
{ val: 50.0,  expected: RiskWarning }
{ val: 80.0,  expected: RiskWarning }
{ val: 80.1,  expected: RiskSafe    }
{ band: nil,  expected: ("", false) }
```

Invert the expectations for `ClassifyLowerIsBetter`. These tests are the contract for the entire classification system and should be exhaustive.

### 7.2 Unit Tests — Evaluators

Each evaluator is tested in complete isolation — no mocks, no tree traversal, just struct literals:

```go
func TestMaxComplexityEvaluator_Danger(t *testing.T) {
    e := MaxComplexityEvaluator{}

    // Not applicable when 0
    _, ok := e.Evaluate(model.CoverageMetrics{MaxCyclomaticComplexity: 0}, &config.Band{Min: 10, Max: 20})
    assert.False(t, ok)

    // Above max → danger (lower is better)
    level, ok := e.Evaluate(model.CoverageMetrics{MaxCyclomaticComplexity: 25}, &config.Band{Min: 10, Max: 20})
    assert.True(t, ok)
    assert.Equal(t, RiskDanger, level)

    // Below min → safe
    level, ok = e.Evaluate(model.CoverageMetrics{MaxCyclomaticComplexity: 5}, &config.Band{Min: 10, Max: 20})
    assert.True(t, ok)
    assert.Equal(t, RiskSafe, level)
}
```

Every evaluator must cover: `IsApplicable` false path, zero-data guard, and all three threshold outcomes.

### 7.3 Integration Tests — Annotator

Build a minimal `model.SummaryTree` in memory, run `Annotate()`, and assert the resulting `Statuses` maps on both file nodes and method nodes. This validates the full domain pipeline without involving any reporter.

### 7.4 Golden File Tests — HTML Builder

Run the HTML builder against a fixed, deterministic `SummaryTree` fixture and compare the JSON output to a committed golden file. Use an `--update-golden` flag to regenerate when intentional output changes are made. This test guards against accidental regressions in the presentation layer.

### 7.5 Config Tests

Parse test YAML files covering:
- Valid config with both `file_metrics` and `method_metrics` populated.
- Overlapping metrics (same key in both lists) — must be handled correctly, not deduplicated erroneously.
- Unknown metric key — must produce a clear error message.
- Empty `method_metrics` list — valid, means no method-level metrics are shown.

### 7.6 End-to-End Smoke Test

Run `nanovision` against the sample Cobertura fixture in the repository. Assert:
- Exit code is 0.
- The HTML output file exists.
- The HTML file contains expected metric section IDs (e.g. `data-metric="line_coverage"`).
- The text output contains the names of the configured metrics.
- Running with only `line_coverage` in `file_metrics` produces HTML and text output that contains no branch or statement coverage data.

---

## 8. How the System Works After Refactoring

### 8.1 Adding a New Metric

Adding `condition_coverage` after the refactoring is done requires touching exactly five locations — none of which are core engine files:

| Step | File | Action |
|---|---|---|
| 1. Constant | `internal/config/config.go` | Add `ConditionCoverage MetricKey = "condition_coverage"`. |
| 2. Model | `internal/model/metrics.go` | Add `ConditionsCovered int` and `ConditionsValid int` to `CoverageMetrics`. |
| 3. Evaluator | `internal/status/evaluators/condition_coverage.go` | Implement `Evaluator`. Call `ClassifyHigherIsBetter`. |
| 4. Provider | `internal/reporter/htmlreact/providers.go` | Implement `FileMetricProvider`. Format as percentage detail struct. |
| 5. Register | Both registry maps | Add entries for the new evaluator and provider. |

`annotate.go` is not touched. `builder.go` is not touched. The text reporter is not touched. The user enables the metric by adding `"condition_coverage"` to `file_metrics` in `nanovision.yaml`.

### 8.2 Configuration Independence

`file_metrics` and `method_metrics` are fully independent. Disabling `branch_coverage` in `file_metrics` has zero effect on method-level reporting. The engine handles the difference by using `ActiveFileMetrics` for the file/directory traversal loop and `ActiveMethodMetrics` for the method traversal loop — they are never merged or shared.

### 8.3 Complexity as a First-Class Metric

`max_cyclomatic_complexity` is a native metric at the file level after this refactoring. The aggregator computes it via max-value propagation. The `MaxComplexityEvaluator` inverts the threshold direction via `ClassifyLowerIsBetter`. The UI provider formats it as a raw integer rather than a percentage. There are no special cases in the engine. It behaves identically to any other metric from the engine's perspective.

### 8.4 SOLID Principles Satisfied

| Principle | How It Is Satisfied |
|---|---|
| **Open/Closed (OCP)** | Core engine files (`annotate.go`, `builder.go`) are never modified to add metrics. Register a new evaluator/provider and you are done. |
| **Single Responsibility (SRP)** | Evaluators own math and classification. Providers own JSON formatting. The annotator engine owns traversal only. The builder owns rendering only. |
| **Dependency Inversion (DIP)** | The annotator depends on the `Evaluator` abstraction, not on concrete metric structs. |
| **Testability** | Evaluators are pure structs with no side effects. Each threshold case is a 5-line test with zero mocks or tree setup. |