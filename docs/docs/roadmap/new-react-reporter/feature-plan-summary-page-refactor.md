# Feature Plan: Summary Page Refactor with `templ`

This document details the plan to refactor the main `index.html` summary page to a server-side rendered approach using `templ`, with React islands for dynamic features like search and charting. This corresponds to Phase 3 of the UI modernization project.

## 1. Feature Overview

The goal is to generate a static `index.html` that serves as the main entry point for the report. It will contain summary cards, a navigable tree of all covered assemblies and classes, and a history chart. The page will be rendered by Go, with client-side JavaScript limited to interactive islands.

## 2. Implementation Plan

1.  **Define Data Structures (`internal/web/types.go`):**
    *   Create a `SummaryPageData` struct to hold all information needed by the summary page template. This includes summary statistics, a list of assemblies and classes for the navigation tree, and historical data for the chart.

2.  **Create the `templ` Component (`templates/summary.templ`):**
    *   Write the main `templ` component for the summary page.
    *   This component will render the summary cards with statistics.
    *   It will generate a static `<ul>` or similar structure containing all assemblies and classes, with links to their respective file pages. This structure will be targeted by the search island.

3.  **Integrate with the Builder:**
    *   Update `HtmlReportBuilder.renderSummaryPage` to populate `SummaryPageData` and call the new `summary.templ` component to generate `index.html`.

4.  **Create React Islands:**
    *   **`SearchFilter` Island:** A React component that takes a text input and filters the visibility of the statically rendered `<li>` elements in the coverage tree. This keeps the data static while making the filtering dynamic.
    *   **`HistoryChart` Island:** A React component that uses a library like `Chart.js`. It will read historical coverage data from a `window.__HISTORY__` object (injected as JSON into a `<script>` tag by `templ`) and render the chart on the client side.

5.  **Cleanup:**
    *   Once the new summary page is functional, all Angular-specific asset copying, data serialization (`window.assemblies`), and related logic can be removed from the `HtmlReportBuilder`, significantly simplifying the codebase.

**Exit Criteria:**
The generated `index.html` is rendered by `templ`. It displays accurate summary information and a full navigation tree. The search input correctly filters the tree, and the history chart renders with the correct data. All links to file pages are correct and functional. The dependency on the Angular framework is completely removed.