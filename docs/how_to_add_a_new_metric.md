# How to add a metric

A metric is one table entry and one function. The config validation, the
report labels and tooltips, the status colors, `nanovision config docs` and
the website configurator all read the table, so you do not touch them.

## 1. Describe the metric

Add a key and an entry in `internal/config/metrics.go`. Use `FileMetricDefs`
for a metric of a file or folder, and `MethodMetricDefs` for a metric of one
method.

```go
const TodoDensity MetricKey = "todo_density"

var FileMetricDefs = []MetricDef{
	// ...
	{
		Key: TodoDensity, 
		Name: "todo_density", 
		Label: "TODO Density", 
		Value: true, 
		LowerIsBetter: true,
		Doc: "TODO comments for each 100 lines. Lower is better.",
	},
}
```

| Field                     | Use                                                                                               |
|---------------------------|---------------------------------------------------------------------------------------------------|
| `Key`                     | The id in the report data and in stored runs. Do not change it later.                             |
| `Name`                    | The name in `nanovision.yaml`, under `metrics.files` or `metrics.methods`.                        |
| `Label`, `Short`, `Title` | The label in the report, the column header, and the name in sentences. Only `Label` is necessary. |
| `Doc`                     | One sentence. It is the tooltip in the report and the text in the docs.                           |
| `Value`                   | Set it for a plain number such as complexity. Leave it out for a percentage.                      |
| `LowerIsBetter`           | Set it when a higher number is worse. Then danger is above the warning range.                     |
| `Needs`                   | Set it when the metric needs data that some report formats lack.                                  |

## 2. Calculate it

Add one function in `internal/calculator/registry.go`, under the same key.

```go
var FileRegistry = map[config.MetricKey]FileCalc{
	// ...
	config.TodoDensity: func(m model.CoverageMetrics) (any, bool) {
		if m.TotalLines == 0 {
			return nil, false // no data: the metric has no value
		}
		return model.ScoreDetail{Value: float64(m.Todos) * 100 / float64(m.TotalLines)}, true
	},
}
```

- Return `model.CoverageDetail` for a percentage. The helper `ratio(covered, valid)` builds it.
- Return `model.ScoreDetail` for a plain number.
- A method metric can use other method metrics: list them in `Deps` and read them from `prior`.

If the metric needs a new raw count (here `Todos`), add the field to
`model.CoverageMetrics` and fill it where the other counts are filled, in
`internal/aggregator`.

## 3. Generate and test

```bash
go generate ./...
go test ./...
```

`go generate` writes the config reference for the website
(`docs/src/generated/config.json`). The tests fail when:

- the metric has no `Label` or no `Doc`
- the metric has a table entry but no calculator, or the reverse
- the generated file is stale

## 4. Use it

```yaml
metrics:
  files:
    - name: todo_density
      warning: "2..5"
```

`nanovision -list-metrics` lists the metric with its description.

## Notes

- **The report selector.** When the user unticks a report, the page calculates
  percentages again from per-report data. A new percentage metric keeps its
  merged value there until you add it to `BuildFileReportIndex` in
  `internal/aggregator/report_index.go`. The page names such metrics in the
  selector, so the number is not mistaken for a filtered one.
- **The delta against a base run.** The comparison covers the metrics whose
  counts the run store keeps (`coverageMetrics` in `internal/compare`). A new
  metric shows no delta until its counts are in `blob.Counts`.
