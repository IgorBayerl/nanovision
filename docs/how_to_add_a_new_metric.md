# How to Add a New Metric to NanoVision

A step-by-step guide covering every file you need to touch when adding a new metric to the pipeline — from configuration constants all the way to reporters.

---

## Architecture Overview

```
  ┌─────────────────────────────────┐
  │  1. config.go                   │
  │     MetricKey constant          │
  └───────────────┬─────────────────┘
                  │
                  ▼
  ┌─────────────────────────────────┐
  │  2. model/metrics.go            │
  │     Data fields                 │
  └───────────────┬─────────────────┘
                  │
                  ▼
  ┌─────────────────────────────────┐
  │  3. aggregator                  │
  │     Calculation & aggregation   │
  └───────────────┬─────────────────┘
                  │
                  ▼
  ┌─────────────────────────────────┐
  │  4. status/evaluator            │
  │     Risk classification         │
  └───────────────┬─────────────────┘
                  │
                  ▼
  ┌─────────────────────────────────┐
  │  5. reporters                   │
  │     HTML + Text output          │
  └─────────────────────────────────┘
```

The metric pipeline has **5 layers**. Each new metric must be registered in all of them.

---

## Step 1 — Define the MetricKey Constant

**File:** [config.go](file:///c:/www/nanovision/internal/config/config.go)

### 1a. Add the MetricKey constant

Add a new `const` entry in the `MetricKey` block:

```diff
 const (
     LineCoverage                 MetricKey = "line_coverage"
     BranchCoverage               MetricKey = "branch_coverage"
     // ... existing keys ...
+    MyNewMetric                  MetricKey = "my_new_metric"
 )
```

### 1b. Optionally add to defaults

If the metric should be shown **by default**, add it to `DefaultFileMetrics` or `DefaultMethodMetrics`:

```diff
 var DefaultFileMetrics = []MetricKey{
     BranchCoverage,
     MethodsHit,
     // ...
+    MyNewMetric,
 }
```

> [!NOTE]
> There is **no** `isValidMetric` switch to update. Metric validation is handled automatically by the evaluator Registry. Once registered in Step 4b, the metric is valid.

---

## Step 2 — Add Data Fields to the Model

**File:** [metrics.go](file:///c:/www/nanovision/internal/model/metrics.go)

### 2a. `CoverageMetrics` (file/dir level)

Add fields to hold numerator and denominator values:

```diff
 type CoverageMetrics struct {
     // ... existing fields ...
+    MyNewMetricCovered int
+    MyNewMetricValid   int
 }
```

### 2b. `MethodMetrics` (method level) — if applicable

If the metric also applies at the method scope, add the corresponding fields:

```diff
 type MethodMetrics struct {
     // ... existing fields ...
+    MyNewMetricCovered int
+    MyNewMetricValid   int
 }
```

---

## Step 3 — Calculate & Aggregate the Metric

**File:** [aggrgator.go](file:///c:/www/nanovision/internal/aggregator/aggrgator.go)

### 3a. Calculate in `calculateFileMethodMetrics` or its helpers

If the metric is derived from method data, update `aggregateStandardMethodMetrics`:

```diff
 func aggregateStandardMethodMetrics(file *model.FileNode, method *model.MethodMetrics) {
     // ... existing logic ...
+    if method.MyNewMetricValid > 0 {
+        file.Metrics.MyNewMetricValid += method.MyNewMetricValid
+        file.Metrics.MyNewMetricCovered += method.MyNewMetricCovered
+    }
 }
```

If it's a patch-based metric, update `calculateMethodPatchMetrics` and/or `calculateFilePatchMetrics`.

### 3b. Propagate in `addMetrics`

This is the function that sums child metrics into parent directories:

```diff
 func addMetrics(dest *model.CoverageMetrics, src model.CoverageMetrics) {
     // ... existing additions ...
+    dest.MyNewMetricCovered += src.MyNewMetricCovered
+    dest.MyNewMetricValid += src.MyNewMetricValid
 }
```

### 3c. Reset in `resetFileMetrics`

Zero out aggregation-target fields before each recalculation:

```diff
 func resetFileMetrics(file *model.FileNode) {
     // ... existing resets ...
+    file.Metrics.MyNewMetricCovered = 0
+    file.Metrics.MyNewMetricValid = 0
 }
```

---

## Step 4 — Status Evaluator (Risk Classification)

This is where your metric gets a **danger / warning / safe** badge.

### 4a. Create the evaluator file

**New file:** `internal/status/evaluators/my_new_metric.go`

Use an existing evaluator as a template. Here's a complete example for a percentage-based "higher is better" metric:

```go
package evaluators

import (
    "github.com/IgorBayerl/nanovision/internal/config"
    "github.com/IgorBayerl/nanovision/internal/model"
    "github.com/IgorBayerl/nanovision/internal/status"
    "github.com/IgorBayerl/nanovision/internal/utils"
)

// MyNewMetricEvaluator evaluates my_new_metric (higher is better).
type MyNewMetricEvaluator struct{}

func (MyNewMetricEvaluator) Key() config.MetricKey { return config.MyNewMetric }
func (MyNewMetricEvaluator) Name() string           { return "My New Metric" }
func (MyNewMetricEvaluator) Description() string    { return "Description of what this metric measures." }
func (MyNewMetricEvaluator) SupportedScopes() status.MetricScope {
    return status.FileScope // or status.MethodScope
}

func (MyNewMetricEvaluator) IsApplicable(_ status.Capabilities) bool {
    // Return true if always applicable, or check caps for conditional metrics.
    // Example for conditional: return caps.HasBranchCoverage
    return true
}

func (MyNewMetricEvaluator) Evaluate(m model.CoverageMetrics, band *config.Band) (status.RiskLevel, bool) {
    // Zero-guard: skip if no data is present.
    if m.MyNewMetricValid == 0 {
        return "", false
    }
    pct := utils.CalculatePercentage(m.MyNewMetricCovered, m.MyNewMetricValid, 2)
    return status.ClassifyHigherIsBetter(pct, band)
}
```

> [!TIP]
> `Name()`, `Description()`, and `SupportedScopes()` power the `--list-metrics` CLI flag and the boot log.
> Adding these ensures your metric is **automatically documented** in the CLI.

### 4b. Register in the evaluator registry

**File:** [registry.go](file:///c:/www/nanovision/internal/status/evaluators/registry.go)

```diff
 var Registry = map[config.MetricKey]status.Evaluator{
     config.LineCoverage:            LineCoverageEvaluator{},
     // ... existing entries ...
+    config.MyNewMetric:             MyNewMetricEvaluator{},
 }
```

### 4c. Update `Capabilities` if conditional

If your evaluator's `IsApplicable` depends on a new capability flag, add it to:

**File:** [types.go](file:///c:/www/nanovision/internal/status/types.go)

```diff
 type Capabilities struct {
     HasBranchCoverage    bool
     HasMethodCoverage    bool
     HasStatementCoverage bool
+    HasMyNewMetric       bool
 }
```

Then update `deriveCapabilities` in [main.go](file:///c:/www/nanovision/cmd/main.go):

```diff
 func deriveCapabilities(tree *model.SummaryTree) status.Capabilities {
     // ... existing walk logic ...
+    if n.Metrics.MyNewMetricValid > 0 {
+        caps.HasMyNewMetric = true
+    }
 }
```

### 4d. Update `methodToCoverageMetrics` (for method-level metrics)

If your metric applies at method scope, map it in [annotate.go](file:///c:/www/nanovision/internal/status/annotate.go):

```diff
 func methodToCoverageMetrics(m *model.MethodMetrics) model.CoverageMetrics {
     cm := model.CoverageMetrics{
         // ... existing mappings ...
+        MyNewMetricCovered: m.MyNewMetricCovered,
+        MyNewMetricValid:   m.MyNewMetricValid,
     }
     return cm
 }
```

---

## Step 5 — Reporter: HTML (React)

Three sub-areas need updates for the HTML report.

### 5a. Metric extraction

**File:** [builder.go](file:///c:/www/nanovision/internal/reporter/htmlreact/builder.go)

Add your metric to `buildMetricsMap()` so it is mapped for the React UI:

```go
func (b *HtmlReactReportBuilder) buildMetricsMap(m model.CoverageMetrics) metricsMap {
    // ... inside the ActiveFileMetrics switch ...
    case config.MyNewMetric:
        if detail, ok := calcData.(model.CoverageDetail); ok {
            metrics[string(key)] = lineCoverageDetail{
                Covered:    detail.Covered,
                Uncovered:  detail.Uncovered,
                Coverable:  detail.Total,
                Total:      detail.Total,
                Percentage: detail.Percentage,
            }
        }
```

### 5b. Metric definitions

**File:** [builder.go](file:///c:/www/nanovision/internal/reporter/htmlreact/builder.go)

Add metric column definitions in `buildMetricDefinitions()`:

```diff
 func (b *HtmlReactReportBuilder) buildMetricDefinitions() metricDefinitions {
     // ... existing definitions ...

+    if b.config.ActiveFileMetrics[config.MyNewMetric] {
+        defs[string(config.MyNewMetric)] = metricDefinition{
+            Label:      "My New Metric",
+            ShortLabel: "New Metric",
+            SubMetrics: []subMetric{
+                {ID: "covered", Label: "Covered", Width: 100},
+                {ID: "total", Label: "Total", Width: 80},
+                {ID: "percentage", Label: "Percentage %", Width: 160},
+            },
+        }
+    }
 }
```

### 5c. Totals struct (if it needs a top-level hero card)

**File:** [schema.go](file:///c:/www/nanovision/internal/reporter/htmlreact/schema.go)

```diff
 type totals struct {
     // ... existing fields ...
+    MyNewMetric *lineCoverageDetail `json:"my_new_metric,omitempty"`
 }
```

Then extract it in `buildTotals()` in [builder.go](file:///c:/www/nanovision/internal/reporter/htmlreact/builder.go):

```diff
 func (b *HtmlReactReportBuilder) buildTotals(...) totals {
     // ... existing extractions ...
+    if mnm, ok := metrics[string(config.MyNewMetric)].(lineCoverageDetail); ok {
+        t.MyNewMetric = &mnm
+    }
 }
```

And the same assignment block in `buildFileTotals()` in [details_generator.go](file:///c:/www/nanovision/internal/reporter/htmlreact/details_generator.go):

```diff
+    assignLineMetric(&t.MyNewMetric, fileMetrics, string(config.MyNewMetric))
```

---

## Step 6 — Reporter: Text Summary

**File:** No changes required! 🎉

The text reporter iterates over configured metrics automatically using `b.config.ActiveFileMetrics` and formats them based on `tree.Metrics.Calculated`. As long as your metric yields a `model.CoverageDetail` or `model.ScoreDetail` and you registered the evaluator, it will just work!

---

## Step 7 — Configuration (YAML)

Users enable the metric in `nanovision.yaml`:

```yaml
file_metrics:
  - "my_new_metric"

# Optional: risk thresholds
status_bands:
  my_new_metric: "60..80"
```

Or via CLI flags:

```bash
nanovision -file-metrics "statement_coverage,my_new_metric" \
           -threshold "my_new_metric=60..80"
```

---

## Step 8 — Write Tests

### Evaluator tests

**File:** `internal/status/evaluators/evaluators_test.go`

Follow the existing pattern — test Key, IsApplicable, zero-guard, and all three threshold bands:

```go
func TestMyNewMetricEvaluator_Key(t *testing.T) {
    assert.Equal(t, config.MyNewMetric, MyNewMetricEvaluator{}.Key())
}

func TestMyNewMetricEvaluator_IsApplicable(t *testing.T) {
    assert.True(t, MyNewMetricEvaluator{}.IsApplicable(status.Capabilities{}))
}

func TestMyNewMetricEvaluator_ZeroGuard(t *testing.T) {
    m := model.CoverageMetrics{MyNewMetricValid: 0}
    band := &config.Band{Min: 60, Max: 80}
    lvl, show := MyNewMetricEvaluator{}.Evaluate(m, band)
    assert.Equal(t, status.RiskLevel(""), lvl)
    assert.False(t, show)
}

func TestMyNewMetricEvaluator_Thresholds(t *testing.T) {
    band := &config.Band{Min: 60, Max: 80}
    tests := []struct {
        name    string
        covered int
        valid   int
        want    status.RiskLevel
    }{
        {"danger", 50, 100, status.RiskDanger},
        {"warning", 70, 100, status.RiskWarning},
        {"safe", 90, 100, status.RiskSafe},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            m := model.CoverageMetrics{
                MyNewMetricCovered: tt.covered,
                MyNewMetricValid:   tt.valid,
            }
            lvl, show := MyNewMetricEvaluator{}.Evaluate(m, band)
            assert.True(t, show)
            assert.Equal(t, tt.want, lvl)
        })
    }
}
```

Also update `TestRegistryContainsAllEvaluators` to include the new key.

---

## Quick-Reference Checklist

| # | File | Action |
|---|------|--------|
| 1 | `internal/config/config.go` | Add `MetricKey` constant, optionally add to defaults |
| 2 | `internal/model/metrics.go` | Add fields to `CoverageMetrics` and/or `MethodMetrics` |
| 3 | `internal/aggregator/aggrgator.go` | Calculate, propagate in `addMetrics`, reset in `resetFileMetrics` |
| 4 | `internal/status/evaluators/<name>.go` | **[NEW]** Create evaluator struct |
| 5 | `internal/status/evaluators/registry.go` | Register evaluator |
| 6 | `internal/status/types.go` | *(Optional)* Add `Capabilities` flag |
| 7 | `internal/status/annotate.go` | *(Optional)* Map method fields in `methodToCoverageMetrics` |
| 8 | `cmd/main.go` | *(Optional)* Update `deriveCapabilities` |
| 9 | `internal/reporter/htmlreact/builder.go` | Map value in `buildMetricsMap` + add definition + extract in `buildTotals` |
| 10 | `internal/reporter/htmlreact/schema.go` | *(Optional)* Add field to `totals` struct |
| 11 | `internal/reporter/htmlreact/details_generator.go` | *(Optional)* Add assignment in `buildFileTotals` |
| 12 | `internal/status/evaluators/evaluators_test.go` | Add evaluator tests + update registry test |

---

## Design Patterns at Play

| Pattern | Where | Purpose |
|---------|-------|---------|
| **Strategy** | `Evaluator` interface | Each metric evaluates its own risk logically |
| **Registry** | `evaluators.Registry` | Decouple metric logic from annotation logic |
| **Open/Closed** | The `Annotate` and text reporting functions | Evaluators and Text UI dynamically handle new metrics with no code changes |

### Classifier Functions

Two classification strategies are available in [classifier.go](file:///c:/www/nanovision/internal/status/classifier.go):

| Function | Semantic | Example |
|----------|----------|---------|
| `ClassifyHigherIsBetter(val, band)` | `val < Min` → Danger, `Min ≤ val ≤ Max` → Warning, `val > Max` → Safe | Coverage % |
| `ClassifyLowerIsBetter(val, band)` | `val < Min` → Safe, `Min ≤ val ≤ Max` → Warning, `val > Max` → Danger | Complexity |
