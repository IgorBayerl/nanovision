# Roadmap: New HTML Reporter (SSR + React Islands)

This section outlines the comprehensive plan to re-architect the AdlerCov HTML reporter. The goal is to move from the current client-side rendered Angular application to a modern, high-performance solution that combines server-side rendering (SSR) with Go and targeted client-side interactivity using React "islands".

This initiative will significantly improve report load times, enhance the developer experience for contributors, and create a more maintainable and extensible frontend architecture for the future.

---

### Architectural Vision

For a high-level overview of the target stack, guiding principles, and expected benefits of this new approach, please see the main architectural vision document.

*   [**Feature Plan: Modern Frontend with Server-Side HTML and React Islands**](./feature-plan-ssr-react-islands.md)

---

### Detailed Implementation Plans

The refactoring process is broken down into two major components: the file detail pages and the main summary page. Each has a detailed implementation plan.

*   [**Feature Plan: File Page Refactor with `templ`**](./feature-plan-file-page-refactor.md) - Details the plan to rebuild the individual code file view, focusing on static HTML generation for performance.
*   [**Feature Plan: Summary Page Refactor with `templ`**](./feature-plan-summary-page-refactor.md) - Outlines the work to recreate the main `index.html`, including summary cards, the coverage tree, and interactive chart components.

---

### Contributor and Development Guide

For contributors interested in working on the new UI, this guide provides essential information on setting up the development environment, running the necessary tools (`templ`, Vite), and understanding the build and CI process.

*   [**UI Development Guide**](./ui-development-guide.md)