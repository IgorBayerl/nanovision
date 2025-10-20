# Feature Plan: Modern Frontend with Server-Side HTML and React Islands

This document outlines the architectural vision for the next generation of the AdlerCov HTML reporter. We will be moving from a purely client-side rendered (Angular SPA) model to a hybrid approach that leverages Go for server-side templating and React for targeted client-side interactivity.

## 1. Architectural Vision

**Target Stack:**
*   **Server-Side Rendering:** Go 1.22+ with `templ` for high-performance, type-safe HTML generation.
*   **Client-Side Interactivity:** Vite 5 + React 18 for creating isolated, interactive components (known as "islands").

**Guiding Principle:**
The core philosophy is to minimize the amount of client-side JavaScript.
*   **Static & Data-Heavy Content:** The main structure of the report, including the file-by-file code view, summary tables, and metrics, will be rendered into static HTML files by Go at generation time. This ensures the report is incredibly fast to load and is fully functional without JavaScript.
*   **Stateful & Interactive Widgets:** Small, specific parts of the UI that require user interaction and state management (like a search filter, a history chart, or coverage toggles) will be implemented as React components. These "islands" will be mounted on the static HTML pages to provide targeted interactivity.

## 2. Benefits of this Approach

*   **Performance:** Reports will load and display coverage data almost instantly, as the bulk of the rendering is done ahead of time by Go. The JavaScript bundle size will be significantly smaller, containing only the code for the interactive islands.
*   **Maintainability:** By separating the static structure from the interactive logic, we create a cleaner codebase. Go is excellent at data handling and HTML generation, while React is purpose-built for stateful UI.
*   **Developer Experience:** `templ` provides compile-time type checking for HTML templates, eliminating a whole class of bugs. Vite offers a fast and modern development environment for building the React islands.
*   **Simplicity:** The final report remains a set of static HTML, CSS, and JS files that can be hosted anywhere, with no need for a live server.

## 3. Implementation Roadmap

The transition will be managed in a series of phased milestones:

1.  **M1 - Skeleton & Proof of Concept:** Establish the build tooling and render a single file page using `templ`, with one working React island to prove the architecture.
2.  **M2 - Core Parity:** Rebuild the file detail pages and the main summary page using the new `templ` + React stack, achieving functional parity with the current Angular report.
3.  **M3 - Polish & Enhancements:** Add remaining UI features like the metrics table, sidebar navigation, and improved tooltips.
4.  **M4 - Release:** Finalize the CI/CD pipeline, add development documentation, and clean up all obsolete code from the previous architecture.

Detailed plans for each phase are available in the following documents.