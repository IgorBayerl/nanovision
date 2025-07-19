# Feature Plan: History System

### High-Level Architectural Approach

We will model the implementation closely on the original C# project, as its design is proven and separates concerns effectively. The main idea is to:

1.  **Store Historical Data**: After generating a report, save a lightweight snapshot of the coverage summary (per class) to a persistent location (the "history directory").
2.  **Enrich Current Data**: When generating a new report, read the historical snapshots and attach them to the corresponding classes being processed.
3.  **Calculate and Display Deltas**: In the report rendering step (specifically for the HTML report), compare the current coverage with a selected historical snapshot to calculate and display the delta.
4.  **Leverage Delta for Patch Coverage**: Use the same mechanism, but by comparing a feature branch against a historical snapshot from the main branch.

---

### Phase 1: The History Storage System (The Foundation)

**Goal:** Create the core components for saving and loading coverage snapshots. This decouples the analysis logic from the storage mechanism (e.g., local files), which is excellent for testability.

1.  **Create a New `history` Package:**
    This will encapsulate all logic related to history storage and parsing.
    *   Create directory: `go_report_generator/internal/history/`

2.  **Define the Storage Interface (`IHistoryStorage`):**
    This interface will abstract the storage mechanism.
    *   **Create file:** `internal/history/storage.go`
    *   **Content:**
        ```go
        package history

        import "io"

        // IHistoryStorage defines the contract for persisting and retrieving history files.
        type IHistoryStorage interface {
            // GetHistoryFilePaths returns the paths to all available history files.
            GetHistoryFilePaths() ([]string, error)
            // LoadFile opens a history file for reading.
            LoadFile(filePath string) (io.ReadCloser, error)
            // SaveFile saves a new history file.
            SaveFile(fileName string, content io.Reader) error
        }
        ```

3.  **Implement `FileHistoryStorage`:**
    This is the concrete implementation for storing files on the local disk.
    *   **Add to file:** `internal/history/storage.go`
    *   **Content:**
        ```go
        package history

        import (
            "io"
            "os"
            "path/filepath"
            "sort"
        )

        // FileHistoryStorage implements IHistoryStorage using the local file system.
        type FileHistoryStorage struct {
            historyDirectory string
        }

        func NewFileHistoryStorage(historyDirectory string) *FileHistoryStorage {
            return &FileHistoryStorage{historyDirectory: historyDirectory}
        }

        func (fhs *FileHistoryStorage) GetHistoryFilePaths() ([]string, error) {
            files, err := filepath.Glob(filepath.Join(fhs.historyDirectory, "*_CoverageHistory.xml"))
            if err != nil {
                return nil, err
            }
            sort.Strings(files) // Ensure consistent order
            return files, nil
        }

        func (fhs *FileHistoryStorage) LoadFile(filePath string) (io.ReadCloser, error) {
            return os.Open(filePath)
        }

        func (fhs *FileHistoryStorage) SaveFile(fileName string, content io.Reader) error {
            path := filepath.Join(fhs.historyDirectory, fileName)
            file, err := os.Create(path)
            if err != nil {
                return err
            }
            defer file.Close()

            _, err = io.Copy(file, content)
            return err
        }
        ```

4.  **Implement the History Writer:**
    This component will take the final `SummaryResult` and create an XML snapshot.
    *   **Create file:** `internal/history/writer.go`
    *   **Content:** The `HistoryWriter` will receive an `IHistoryStorage`, iterate through the assemblies and classes of the current report, and generate an XML file in a format similar to this, which it then saves using the storage interface.

        *Example XML Structure (`2025-07-15_10-30-00_CoverageHistory.xml`):*
        ```xml
        <coverage version="1.0" date="2025-07-15_10-30-00" tag="build_123">
          <assembly name="MyProject.Core">
            <class name="MyProject.Core.Calculator" coveredlines="10" coverablelines="12" totallines="20" coveredbranches="2" totalbranches="4" coveredcodeelements="2" totalcodeelements="2" />
            <class name="MyProject.Core.Utils" coveredlines="25" coverablelines="30" totallines="50" coveredbranches="8" totalbranches="10" coveredcodeelements="5" totalcodeelements="6" />
          </assembly>
        </coverage>
        ```

### Phase 2: Integrating History into the Analysis Pipeline

**Goal:** Modify the main application flow to use the new history system.

1.  **Update the Data Model:**
    The `Class` model needs a place to store the historical data that will be read.
    *   **Modify file:** `internal/model/analysis.go`
    *   **Action:** Add a field to the `Class` struct.
        ```go
        // In internal/model/analysis.go
        type Class struct {
            // ... existing fields ...
            HistoricCoverages   []HistoricCoverage // Add this line
        }
        ```
    *   The `HistoricCoverage` struct already exists in `internal/model/historiccoverage.go` and is suitable for this purpose.

2.  **Implement the History Parser:**
    This component will be responsible for reading history files and attaching the data to the `Class` models.
    *   **Create file:** `internal/history/parser.go`
    *   **Content:**
        ```go
        package history

        import (
            // ... imports ...
            "github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/model"
        )

        // HistoryParser reads history files and enriches the current analysis result.
        type HistoryParser struct {
            storage IHistoryStorage
            // ... other dependencies like max files to read
        }

        func NewHistoryParser(storage IHistoryStorage) *HistoryParser {
            return &HistoryParser{storage: storage}
        }

        // ApplyHistoricCoverage reads history and attaches it to the assemblies.
        func (hp *HistoryParser) ApplyHistoricCoverage(assemblies []model.Assembly) error {
            // 1. Get all history file paths from storage.
            // 2. Limit to the N most recent files.
            // 3. For each file, parse the XML.
            // 4. For each <class> element in the XML, find the matching model.Class
            //    in the provided assemblies slice.
            // 5. Create a model.HistoricCoverage object from the XML attributes.
            // 6. Add the HistoricCoverage object to the Class.HistoricCoverages slice.
            // This logic will be very similar to the C# HistoryParser.cs.
            return nil
        }
        ```

3.  **Update the Main Application Flow:**
    The entry point in `cmd/main.go` needs to orchestrate these new components.
    *   **Modify file:** `cmd/main.go`
    *   **Action Items:**
        1.  **Add a new flag:** Create a `-historydir` command-line flag to specify the history storage location.
        2.  **Modify `run()` function:**
            *   After parsing flags, check if `-historydir` was provided.
            *   If yes, create an instance of your new `history.FileHistoryStorage`.
            *   **During Analysis:** After `parseAndMergeReports` completes, create a `history.NewHistoryParser` and call its `ApplyHistoricCoverage` method, passing in the `summaryResult.Assemblies`. This enriches the data *before* reports are generated.
            *   **After Reporting:** After `generateReports` successfully completes, create a `history.NewHistoryWriter` and call a method like `CreateReport`, passing it the `summaryResult` and the current timestamp/tag. This saves the current run for future comparisons.

### Phase 3: Implementing Delta Coverage in the HTML Report

**Goal:** Use the newly attached historical data to calculate and display coverage changes.

1.  **Update the HTML View Model:**
    The data structure you pass to your HTML templates needs to include fields for the deltas.
    *   **Modify file:** `internal/reporter/htmlreport/viewmodels.go`
    *   **Action:** In `AngularClassViewModel`, you already have fields for historic data (`lch`, `bch`, etc.). Your Angular frontend seems prepared for this. The key is to ensure the data passed to it is correct.
    *   Your `go_report_generator/angular_frontend_spa/src/app/components/coverageinfo/class-row.component.ts` already contains the logic to display differences:
        ```typescript
        <div class="currenthistory {{getClassName(clazz.coveredLines, clazz.currentHistoricCoverage.cl)}}">
          {{clazz.coveredLines}}
        </div>
        <div [title]="clazz.currentHistoricCoverage.et">
          {{clazz.currentHistoricCoverage.cl}}
        </div>
        ```
        The `getClassName` function handles the `lightgreen`/`lightred` styling. This is perfect.

2.  **Enhance the Report Builder:**
    The logic for building the view models needs to populate the history fields.
    *   **Modify file:** `internal/reporter/htmlreport/summary_builder.go`
    *   **Action:**
        1.  In `buildAngularClassViewModelForSummary`, you are already iterating through `class.HistoricCoverages`. This is correct. This data, once populated by the `HistoryParser`, will flow directly into the JSON used by the Angular frontend.
        2.  You also need to collect all unique execution times from all classes to populate the "Compare with" dropdown. In `prepareGlobalJSONData`, create a helper function that iterates through `report.Assemblies`, collects all `HistoricCoverage` execution times into a `map[string]struct{}` to get unique values, and then marshals this list into `historicCoverageExecutionTimesJSON`.

The front-end seems largely ready to consume this data. The main work is in the Go backend to correctly parse, attach, and serialize the historical data.

### Phase 4: Using Delta Coverage for Patch Analysis

**Goal:** Document and explain the workflow for using the new system to analyze patches. This phase is primarily about process, not code.

1.  **Establish a Baseline:**
    *   Configure your CI/CD pipeline to run the report generator on every commit to your `main` or `develop` branch.
    *   This run **must** use the `-historydir` flag, pointing to a persistent, shared location (e.g., a mounted volume, a checked-in directory if small, or a custom S3 storage implementation of `IHistoryStorage`).
    *   This creates a history of coverage for your main branch.

2.  **Analyze a Feature Branch/Pull Request:**
    *   In the CI job for a pull request or feature branch, run the report generator on the test results for that branch.
    *   Crucially, configure it to use the **same `-historydir`** as the main branch job.
    *   The `HistoryParser` will now load the history from the `main` branch.
    *   The generated HTML report will automatically contain the historical data. In the UI, a user can select a recent `main` branch build from the "Compare with" dropdown.
    *   The report will then display the delta coverage, which effectively represents the **patch coverage**: how the changes in the pull request have affected the overall code coverage.

By following this plan, you will have a robust, testable, and maintainable history and delta coverage system that mirrors the capabilities of the original .NET tool.