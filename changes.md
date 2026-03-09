# Changes Summary

This branch introduced dependency support between calculators and implemented several composite method-level risk metrics.

## Calculator Dependencies
- **`MetricCalculator` and `MethodMetricCalculator` updates**:
  - Added `DependsOn() []config.MetricKey` to declare dependencies.
  - Updated the `Calculate` signature to accept a `prior` map containing already computed metrics.
- **Topological Sorting Execution**:
  - Refactored `CalculateTree` in `internal/calculator/engine.go` to use a Depth-First Search (DFS) topological sort. This ensures that a calculator's dependencies are executed before the calculator itself.
- **Updated Existing Calculators**:
  - All existing calculators were refactored to conform to the new interface signatures by returning empty dependency arrays.

## New Derived Risk Metrics
- **`config.MetricKey` additions**:
  - Added `MethodCrapScore`, `MethodPatchCrapScore`, `MethodExposedRisk`, and `MethodDefectProbability`.
- **Calculators Implementation**:
  - **`MethodCrapScoreCalculator`**: Uses Cyclomatic Complexity and Method Line Coverage to compute the CRAP score.
  - **`MethodPatchCrapScoreCalculator` (PCRAP)**: Adapts the CRAP metric for patches using Patch Statement Coverage.
  - **`MethodExposedRiskCalculator`**: Calculates the absolute volume of unprotected complexity using Statement Coverage.
  - **`MethodDefectProbabilityCalculator` (DPI)**: Flags methods as High Risk if complexity is > 10, patch coverage < 50%, and statement coverage < 70%.
- **Registry Integration**:
  - Registered all new method metric calculators in `internal/calculator/registry.go`.

## Reporting & UI
- **`nanovision.yaml` Config**:
  - Added the new method metrics to the default `method_metrics` configuration to surface them in self-coverage reports.
- **UI Schema Definition (`internal/reporter/htmlreact/schema.go`)**:
  - Introduced corresponding UI constants (`MethodUICrapScore`, etc.) to enforce alphabetical sorting.
- **HTML Report Builder (`internal/reporter/htmlreact/builder.go`)**:
  - Wired up metric definitions for the JSON payload mapped to the new React HTML UI components.
- **Details Page Generation (`internal/reporter/htmlreact/details_generator.go`)**:
  - Dynamically extracted calculations and applied format strings (e.g., `%.2f`) before injecting them into the schema structure.
- **Rebuilt UI Assets (`assets/dist`)**:
  - Generated fresh bundles with `npm run build` using the updated UI structure logic.
