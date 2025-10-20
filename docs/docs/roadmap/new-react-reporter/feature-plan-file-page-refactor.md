# Feature Plan: File Page Refactor with `templ`

This document details the plan to refactor the file detail pages from the current Angular implementation to a server-side rendered approach using `templ`, enhanced with small React islands for interactivity. This corresponds to Phases 1 and 2 of the UI modernization project.

## 1. Feature Overview

The goal is to generate a static HTML page for each covered file in the report. The page will display the source code with line-by-line coverage information (hits, status, branch data). The rendering will be performed by Go at report generation time.

## 2. Implementation Plan

### Phase 1: Minimal Vertical Slice

This phase focuses on creating a single, working file page to validate the architecture.

1.  **Define Data Structures (`internal/web/types.go`):**
    *   Create new Go structs to pass data to the `templ` components.
        ```go
        package web

        // FilePageData holds all necessary data for rendering a file detail page.
        type FilePageData struct {
            FileName     string
            SourceLines  []LineRow
            // ... other metadata
        }

        // LineRow represents a single line of code with its coverage info.
        type LineRow struct {
            LineNumber  int
            Content     string
            Hits        int
            Status      string // e.g., "covered", "uncovered", "partial"
            HasBranch   bool
            // ... other line-specific data
        }
        ```

2.  **Create the `templ` Component (`templates/file_page.templ`):**
    *   Write a `templ` component that accepts `FilePageData` and renders the HTML structure for the file view, including the code table.
    *   Use `templ`'s features for loops (`for _, line := range data.SourceLines`) and conditionals (`if line.HasBranch`).

3.  **Integrate with the Builder:**
    *   Modify `HtmlReportBuilder.renderClassDetailPages` to use the new `templ` component.
    *   Initially, this will be hardcoded to process only the first file of the first class to create a proof-of-concept output file.

4.  **Create React "Island" for Interactivity:**
    *   In the `/ui` directory, create a simple React component, e.g., `CoverageToggle`.
    *   Use a custom element wrapper to allow it to be used in the static HTML like `<coverage-toggle>`.
    *   Build this island using `vite build`, which will produce a single JS file (e.g., `react-islands.js`).

5.  **Embed and Use Assets:**
    *   Update `internal/assets/assets.go` to embed the new `react-islands.js` and a minimal `report.css`.
    *   The `file_page.templ` component will include the `<script>` tag to load the React island.

### Phase 2: Full Feature Parity

This phase expands the minimal slice to include all features of the current file page.

1.  **Enhance Line Data:** The `LineRow` struct and the builder logic will be expanded to include detailed branch information, hit counts, and tooltip text.
2.  **CSS Styling:** Port the necessary CSS rules from the existing `custom.css` to `assets/report.css` to correctly style covered, uncovered, and partially covered lines, as well as the branch coverage bar.
3.  **Generate All Pages:** Modify the `HtmlReportBuilder` to loop through all assemblies, classes, and files, generating a unique HTML page for each file.

**Exit Criteria:**
Running the report generator produces a complete set of file detail pages. Each page displays the source code with accurate coverage highlighting, branch indicators, and tooltips. The interactive toggle (React island) works correctly on all pages.