# Feature Plan: Hierarchical and Decentralized Components

This document outlines the implementation plan for a component-based coverage analysis feature. This feature will allow teams to define logical code boundaries, gain granular insights, and establish clear ownership within the repository.

## 1. Feature Overview

The core of this feature is to move away from a single, repository-wide coverage number and introduce a system of **components**. A component is a logical group of source files, typically representing a microservice, a library, or a feature set.

This implementation will be founded on two key principles:

1.  **Decentralized Configuration**: Instead of a single root configuration file, teams will define components using a `component.yml` file placed directly within the directory that the component represents.
2.  **Hierarchical Structure**: The system will automatically build a nested component hierarchy based on the file system structure. A `component.yml` inside a directory managed by another `component.yml` becomes a sub-component (e.g., `backend/authentication`).

## 2. User Impact & Benefits

* **Team Autonomy**: The backend and frontend teams can manage their own component definitions in their respective directories (`/services/backend/`, `/packages/frontend/`) without creating configuration bottlenecks or merge conflicts. [cite: 556]
* **Intuitive Configuration**: The component structure directly mirrors the codebase's file structure, making it easy to understand and maintain.
* **Granular Analysis**: Users will be able to see coverage metrics not just for the whole project, but for individual components like `backend`, `frontend`, or even `backend/authentication`.
* **Actionable Reporting**: The HTML report will be enhanced to allow filtering by component, making it easy for a developer to see the impact of their changes on the parts of the code they own.
* **Foundation for Future Features**: This data structure will enable future enhancements like component-based status checks, ownership tracking via a code owners file, and per-component historical trend analysis.

## 3. High-Level Architectural Approach

We will model this implementation on the existing patterns within the `go-report-generator`, particularly the separation of concerns seen in the `parser` -> `analyzer` -> `reporter` pipeline. The plan is as follows:

1.  **Model (`internal/model`)**: The core data model will be extended to support component data. Specifically, `model.Class` will be tagged with the components it belongs to.
2.  **Discovery (`internal/components`)**: A new package will be created to handle the discovery and parsing of all `component.yml` files, building an in-memory tree of the component hierarchy.
3.  **Enrichment (`internal/analyzer`)**: The analysis step will be augmented with a new enrichment phase. After parser results are merged into a `SummaryResult`, this phase will use the discovered component tree to tag each `Class` with its corresponding component(s).
4.  **Reporting (`internal/reporter/htmlreport`)**: The HTML reporter will be updated to consume and display this new component data, primarily by adding component tags and filtering capabilities to the Angular frontend.

## 4. Detailed Implementation Plan

### Phase 1: Core Model and Configuration Parsing

**Goal:** Establish the data structures and the logic for discovering and parsing `component.yml` files.

1.  **Create a New `components` Package:**
    * **Action:** Create a new directory `internal/components/`. This will encapsulate all logic related to component definition and discovery, keeping it separate from other concerns.

2.  **Define the `component.yml` Structure:**
    * **Action:** Create `internal/components/config.go`.
    * **Content:** Define the struct that represents the `component.yml` file format.
        ```go
        // in: internal/components/config.go
        package components

        // ComponentConfig defines the structure of a component.yml file.
        type ComponentConfig struct {
            Name   string  `yaml:"name"`
            Target float64 `yaml:"target,omitempty"`
            Owner  string  `yaml:"owner,omitempty"`
        }
        ```

3.  **Implement the Component Discovery Logic:**
    * **Action:** Create `internal/components/discovery.go`.
    * **Content:** This will contain the logic to walk the filesystem and build the component tree.
        ```go
        // in: internal/components/discovery.go
        package components

        import "path/filepath"

        // ComponentNode represents a single component in the hierarchy.
        type ComponentNode struct {
            Path     string           // The absolute directory path
            FullName string           // The full, resolved component name (e.g., "backend/auth")
            Config   *ComponentConfig
            Parent   *ComponentNode
            Children map[string]*ComponentNode
        }

        // DiscoverComponents walks the filesystem from the root, finds all component.yml files,
        // and builds a hierarchical tree of ComponentNodes.
        func DiscoverComponents(rootDir string) (*ComponentNode, error) {
            // 1. Create a root node for the repository.
            // 2. Use filepath.Walk to scan for "component.yml".
            // 3. For each component.yml found:
            //    a. Parse the YAML into a ComponentConfig.
            //    b. Create a ComponentNode.
            //    c. Traverse up from the component's directory to find its parent node in the tree.
            //    d. Construct the FullName (e.g., parent.FullName + "/" + config.Name).
            //    e. Add the new node to its parent's Children map.
            // 4. Return the root node of the tree.
        }
        ```

4.  **Update the Core Data Model:**
    * **Goal:** Allow a class to be associated with one or more components.
    * **Action:** Modify `internal/model/analysis.go`.
    * **Change:** Add a `Components` field to the `model.Class` struct.
        ```go
        // in: internal/model/analysis.go
        type Class struct {
            // ... existing fields ...
            Components []string // Add this line. Holds resolved component names like "backend/auth".
        }
        ```

### Phase 2: Integrating Component Enrichment into the Pipeline

**Goal:** Use the discovered component tree to tag the parsed coverage data.

1.  **Implement the Enrichment Logic:**
    * **Action:** Create a new file `internal/analyzer/enrichment.go` or add to `analyzer.go`.
    * **Content:** This logic will apply the component data to the `SummaryResult`.
        ```go
        // in: internal/analyzer/enrichment.go (or similar)
        package analyzer

        // ApplyComponents iterates through the summary result and tags each class
        // with the components it belongs to based on file paths.
        func ApplyComponents(summary *model.SummaryResult, componentTree *components.ComponentNode) {
            // 1. Create a cache to store resolved component paths for file paths to avoid redundant lookups.
            // 2. Iterate through summary.Assemblies -> class -> file.
            // 3. For each file.Path:
            //    a. Look up its component path from the cache.
            //    b. If not in cache, walk up from file.Path's directory, checking against the componentTree
            //       to find the owning ComponentNode and its FullName. Store it in the cache.
            //    c. Add the resolved component name to the class.Components slice (ensure no duplicates).
        }
        ```

2.  **Update the Main Application Flow (`cmd/main.go`):**
    * **Goal:** Orchestrate the new discovery and enrichment steps.
    * **Action:** Modify the `run()` function in `cmd/main.go`.
        1.  **Add a new flag:** Introduce a `-components` flag to enable/disable this feature. E.g., `componentsEnabled := flag.Bool("components", true, "Enable component-based analysis")`.
        2.  **Modify `run()` function:** The new steps should be called after parsing/merging and before report generation.

            ```go
            // in: cmd/main.go's run() function

            // ... after resolveAndValidateInputs ...
            reportConfig, err := createReportConfiguration(...)
            if err != nil { return err }

            summaryResult, err := parseAndMergeReports(logger, reportConfig, parserFactory)
            if err != nil { return err }

            // --- NEW COMPONENT ENRICHMENT STEP ---
            if *flags.componentsEnabled { // Check if the feature is enabled
                logger.Info("Starting component discovery...")
                // Assume the project root can be determined (e.g., from CWD or a new flag).
                projectRoot := "." // Placeholder for project root discovery
                componentTree, err := components.DiscoverComponents(projectRoot)
                if err != nil {
                    logger.Warn("Component discovery failed, continuing without component data.", "error", err)
                } else {
                    logger.Info("Applying component data to analysis results...")
                    analyzer.ApplyComponents(summaryResult, componentTree)
                    logger.Info("Component enrichment complete.")
                }
            }
            // --- END OF NEW STEP ---

            reportCtx := reporting.NewReportContext(...)
            return generateReports(reportCtx, summaryResult)
            ```

### Phase 3: Updating the HTML Reporter

**Goal:** Expose component data in the final HTML report for user interaction.

1.  **Update HTML View Models:**
    * **Action:** Modify `internal/reporter/htmlreport/viewmodels.go`.
    * **Change:** Add a `Components` field to `AngularClassViewModel`.
        ```go
        // in: internal/reporter/htmlreport/viewmodels.go
        type AngularClassViewModel struct {
            // ... existing fields ...
            Components []string `json:"cmps,omitempty"` // Add this for component tags
        }
        ```

2.  **Update Summary Builder Logic:**
    * **Action:** Modify `buildAngularClassViewModelForSummary` in `internal/reporter/htmlreport/summary_builder.go`.
    * **Change:** Copy the component data from the model to the view model.
        ```go
        // in: internal/reporter/htmlreport/summary_builder.go
        func (b *HtmlReportBuilder) buildAngularClassViewModelForSummary(...) AngularClassViewModel {
            angularClass := AngularClassViewModel{
                // ... existing assignments ...
                Components: class.Components, // Add this line
            }
            // ... rest of the function ...
            return angularClass
        }
        ```

3.  **Enhance the Angular Frontend (Conceptual):**
    * **Goal:** Since the Angular source is not provided, this plan outlines the necessary changes.
    * **Actions:**
        1.  **Modify `coverage-info` component:** Update the component to read the new `cmps` array from each class object in `window.assemblies`.
        2.  **Display Component Tags:** In the class rows of the coverage table, render the component names as labels or tags next to the class name.
        3.  **Implement Component Filter:** Add a new dropdown filter to the report controls. This dropdown will be populated with a unique list of all component names discovered from `window.assemblies`. Selecting a component will filter the class list to show only classes belonging to that component.

## 5. Future Possibilities Enabled by This Feature

* **Component-Based Status Checks**: Introduce a mechanism in CI to fail a build if the coverage for a specific, critical component (e.g., `backend/payments`) drops below its defined `target` in `component.yml`.
* **Ownership and Notifications**: Use the `owner` field in `component.yml` to integrate with GitHub/GitLab APIs, automatically assigning reviewers or notifying the owning team of coverage changes in their components.
* **Component-Level History**: Extend the `HistoryCoverage.md` plan to store and display coverage trends on a per-component basis, not just per-class.
* **Targeted Reporting**: Add a command-line flag (`-report.components=backend/auth,frontend`) to generate a report containing only the specified components.