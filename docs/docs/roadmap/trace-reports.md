# Feature Plan: Report Source Visualization

This document outlines the implementation plan for visualizing the source of code coverage. The feature will allow users to identify exactly which input report file is responsible for the coverage of any given line of code, directly within the final HTML report.

## 1. Feature Overview

When multiple coverage reports (e.g., from unit tests and integration tests) are merged, it becomes difficult to trace the origin of the coverage data. This feature addresses that by tagging each line of code with its source report during the analysis phase.

The final HTML report will then use this information to provide enhanced tooltips and enable new filtering capabilities, allowing a user to see not just *if* a line is covered, but *by which test suite*.

## 2. User Impact & Benefits

*   **Traceability & Debugging**: Developers can instantly see which report (and by extension, which test suite) covered a specific line. This is invaluable for debugging why a line is, or is not, being covered by a particular set of tests.
*   **Granular Analysis**: Enables a clear understanding of coverage contributions from different test strategies (e.g., backend vs. frontend, unit vs. E2E). This helps teams assess the effectiveness of each part of their testing pipeline.
*   **Confidence in Merged Reports**: Removes the "black box" nature of merged reports, giving users confidence that coverage data from all sources is being correctly processed and attributed.
*   **Actionable Reporting**: The new data provides a foundation for more advanced filtering in the UI, such as "show me only lines covered by `integration_tests.xml`".

## 3. High-Level Architectural Approach

The implementation will be integrated into the existing `Parse -> Analyze -> Report` pipeline with minimal disruption to the overall architecture.

1.  **Model (`internal/model`)**: The core `model.Line` struct will be extended to store a map of source report file paths to their corresponding hit counts for that line.
2.  **Parser (`internal/parser`)**: The `ParserResult` struct will be tagged with the path of the file it was parsed from. This happens in the application entrypoint (`cmd/main.go`).
3.  **Analyzer (`internal/analyzer`)**: The merging logic in the `analyzer` will be enhanced to perform a "deep merge". Instead of simply combining file lists, it will now merge line-level data for files that appear in multiple reports, aggregating hit counts and populating the new source report map on each line.
4.  **Reporter (`internal/reporter/htmlreport`)**: The HTML reporter's view models will be updated to carry this new source data. The builder will populate the view models, which will then be serialized to JSON and embedded in the final HTML file for use by the Angular frontend.

## 4. Detailed Implementation Plan

### Phase 1: Core Model and Parser Result Enhancement

**Goal**: Update the data model to support per-line report attribution.

1.  **Update the `Line` Model:**
    *   **File:** `internal/model/analysis.go`
    *   **Action:** Add a map to the `Line` struct to store hits per source report.

    ```go
    // in: internal/model/analysis.go
    type Line struct {
        // ... existing fields ...
        LineVisitStatus          LineVisitStatus
        CoverageByReport         map[string]int // Key: report file path, Value: hits from that report
    }
    ```

2.  **Update the `ParserResult` Struct:**
    *   **File:** `internal/parser/parser_config.go`
    *   **Action:** Add a field to `ParserResult` to track its origin file path. This makes the path available to the analyzer during merging.

    ```go
    // in: internal/parser/parser_config.go
    type ParserResult struct {
        ReportFilePath         string // Add this field
        Assemblies             []model.Assembly
        // ... existing fields ...
    }
    ```

3.  **Tag the `ParserResult` at Creation:**
    *   **File:** `cmd/main.go`
    *   **Action:** In the `parseAndMergeReports` function, assign the `reportFile` path to the new `ReportFilePath` field immediately after a successful parse.

    ```go
    // in: cmd/main.go (inside parseAndMergeReports loop)
    result, err := parserInstance.Parse(reportFile, reportConfig)
    if err != nil {
        // ... error handling ...
        continue
    }
    result.ReportFilePath = reportFile // Tag the result with its source file path
    parserResults = append(parserResults, result)
    ```

### Phase 2: Refactor the Analyzer for Deep Merging

**Goal**: Enhance the analyzer to merge coverage data at the line level and populate the `CoverageByReport` map.

1.  **Modify `mergeAssemblies` Logic:**
    *   **File:** `internal/analyzer/analyzer.go`
    *   **Action:** The core change is within the `mergeAssemblies` function. The current logic for merging classes is too shallow. It must be modified to handle deep merging of files and their lines.

2.  **Implement `mergeFiles` and `mergeLines` (Conceptual Representation):**
    *   **File:** `internal/analyzer/analyzer.go`
    *   **Action:** The logic inside the `if existingClass, found := classMap[classFromParser.Name]; found` block needs to be replaced with a more robust file-merging strategy.

    ```go
    // Conceptual change in internal/analyzer/analyzer.go inside mergeAssemblies

    // ... existing loop over classFromParser.Classes ...
    if existingClass, found := classMap[classFromParser.Name]; found {
        // ... merge class-level statistics ...

        // --- DEEP MERGE LOGIC ---
        // This replaces the simple file list append.
        existingFilesMap := make(map[string]*model.CodeFile)
        for i := range existingClass.Files {
            existingFilesMap[existingClass.Files[i].Path] = &existingClass.Files[i]
        }

        for _, fileFromParser := range classFromParser.Files {
            if existingFile, fileFound := existingFilesMap[fileFromParser.Path]; fileFound {
                // File already exists, merge lines
                mergeLines(existingFile, &fileFromParser, res.ReportFilePath) // res is the current ParserResult
            } else {
                // New file for this class, tag all its lines and append
                tagLinesOfNewFile(&fileFromParser, res.ReportFilePath)
                existingClass.Files = append(existingClass.Files, fileFromParser)
                existingFilesMap[fileFromParser.Path] = &existingClass.Files[len(existingClass.Files)-1]
            }
        }
        // --- END OF DEEP MERGE ---

    } else {
        // New class, just tag all lines in all its files before adding
        for i := range asmCopy.Classes {
            for j := range asmCopy.Classes[i].Files {
                tagLinesOfNewFile(&asmCopy.Classes[i].Files[j], res.ReportFilePath)
            }
        }
        mergedAssembliesMap[asmCopy.Name] = &asmCopy
    }
    ```

3.  **Define Helper Functions for Merging and Tagging:**
    *   **File:** `internal/analyzer/analyzer.go`
    *   **Action:** Add new helper functions to support the logic above.

    ```go
    // in: internal/analyzer/analyzer.go

    // mergeLines merges line data from newFile into existingFile.
    func mergeLines(existingFile, newFile *model.CodeFile, reportPath string) {
        // Create a map for quick lookup of existing lines by number
        linesMap := make(map[int]*model.Line, len(existingFile.Lines))
        for i := range existingFile.Lines {
            linesMap[existingFile.Lines[i].Number] = &existingFile.Lines[i]
        }

        for _, newLine := range newFile.Lines {
            if existingLine, found := linesMap[newLine.Number]; found {
                // Line exists, merge hits and report data
                existingLine.Hits += newLine.Hits
                if existingLine.CoverageByReport == nil {
                    existingLine.CoverageByReport = make(map[string]int)
                }
                existingLine.CoverageByReport[reportPath] = newLine.Hits
            } else {
                // New line for this file, tag and append
                if newLine.CoverageByReport == nil {
                    newLine.CoverageByReport = make(map[string]int)
                }
                newLine.CoverageByReport[reportPath] = newLine.Hits
                existingFile.Lines = append(existingFile.Lines, newLine)
                linesMap[newLine.Number] = &existingFile.Lines[len(existingFile.Lines)-1]
            }
        }
    }

    // tagLinesOfNewFile initializes the CoverageByReport map for all lines in a file.
    func tagLinesOfNewFile(file *model.CodeFile, reportPath string) {
        for i := range file.Lines {
            line := &file.Lines[i]
            if line.CoverageByReport == nil {
                line.CoverageByReport = make(map[string]int)
            }
            line.CoverageByReport[reportPath] = line.Hits
        }
    }
    ```

### Phase 3: HTML Reporter Integration

**Goal**: Pass the new report source data to the frontend via the `window` object.

1.  **Update View Models:**
    *   **File:** `internal/reporter/htmlreport/viewmodels.go`
    *   **Action:** Add the `CoverageByReport` map to the `AngularLineAnalysisViewModel`.

    ```go
    // in: internal/reporter/htmlreport/viewmodels.go
    type AngularLineAnalysisViewModel struct {
        // ... existing fields ...
        CoverageByReport map[string]int `json:"cbr,omitempty"` // Add this field
    }
    ```

2.  **Populate View Models in Builder:**
    *   **File:** `internal/reporter/htmlreport/class_detail_builder.go`
    *   **Action:** In `buildAngularLineViewModelForJS`, copy the map from the model to the view model.

    ```go
    // in: internal/reporter/htmlreport/class_detail_builder.go
    func (b *HtmlReportBuilder) buildAngularLineViewModelForJS(...) AngularLineAnalysisViewModel {
        lineVM := AngularLineAnalysisViewModel{
            // ... existing assignments ...
        }
        if hasCoverageData {
            // ... existing assignments ...
            if len(modelCovLine.CoverageByReport) > 0 { // Only include the map if it's not empty
                lineVM.CoverageByReport = modelCovLine.CoverageByReport
            }
        }
        // ...
        return lineVM
    }
    ```

### Phase 4: Frontend Enhancement (Conceptual)

**Goal**: Use the new data in the Angular frontend to provide a better user experience.

*   **Tooltip Enhancement**:
    *   **Action**: Modify the Angular component responsible for rendering code lines. When generating the tooltip for a line, check for the presence of the new `cbr` (`CoverageByReport`) map.
    *   **Logic**: If the map exists, iterate through its key-value pairs to build a detailed tooltip string.
    *   **Example Tooltip**: "Covered (5 visits) by: report_unit.xml (3), report_integ.xml (2)".

*   **Source Report Filtering (Stretch Goal)**:
    *   **Action**: Add a new dropdown or multi-select filter to the report controls. This control will be populated with a unique list of all report file paths found across all lines.
    *   **Logic**: When a user selects one or more reports from the filter, JavaScript will iterate through the line elements, adding a "highlight" CSS class to lines whose `cbr` map contains a key matching the selected report(s).

## 5. Testing Strategy

1.  **Unit Tests (Analyzer)**: Create a new test in `internal/analyzer/analyzer_test.go` that specifically validates the deep merge logic.
    *   Define two `parser.ParserResult` objects originating from different `ReportFilePath` values.
    *   Ensure they cover the same assembly, class, and file.
    *   Have them cover some of the same lines and some unique lines.
    *   Assert that the final merged `model.Line` objects have correctly summed `Hits` and a `CoverageByReport` map containing entries for both source reports with the correct hit counts.

2.  **Unit Tests (Reporter)**: In `internal/reporter/htmlreport/`, add a test to verify that a `model.Line` with a populated `CoverageByReport` map results in an `AngularLineAnalysisViewModel` with the `cbr` field correctly populated.

3.  **End-to-End Test**: Create a small, temporary test that runs the main application `run()` function with two simple, mock coverage files. Inspect the generated `index.html` and verify that the `window.classDetails` JSON object contains the expected `cbr` data for the relevant lines.

## 6. Definition of Done

*   [ ] The `model.Line` and `parser.ParserResult` structs are updated as specified.
*   [ ] The `analyzer` package correctly performs a deep merge of line-level data, populating the `CoverageByReport` map.
*   [ ] Unit tests for the analyzer's deep merge functionality are implemented and passing.
*   [ ] The `htmlreport` builder correctly populates the new `cbr` field in the `AngularLineAnalysisViewModel`.
*   [ ] The final HTML report contains the `cbr` data in the `window.classDetails` JSON object.
*   [ ] The line-level tooltips in the HTML report are enhanced to display the source report(s) and their respective hit counts.
*   [ ] All existing tests continue to pass.