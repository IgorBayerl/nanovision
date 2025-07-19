# Project Roadmap

This section outlines the planned features and architectural enhancements for AdlerCov. Our goal is to build upon the solid foundation of ReportGenerator while introducing modern, Go-native features that improve developer workflows.

The following documents provide detailed implementation plans for major upcoming features.

---

### 1. Hierarchical and Decentralized Components
**Status: Planned**

Move away from a single, repository-wide coverage number and introduce a system of **components**. Teams will be able to define logical code boundaries (e.g., a microservice, a library) in a decentralized way, enabling granular insights, clear ownership, and autonomous configuration.

[**Read the detailed feature plan...**](./components.md)

---

### 2. History and Delta Coverage Analysis
**Status: Planned**

Implement a history system to store lightweight snapshots of coverage over time. This will enable trend analysis and, most importantly, **delta coverage**. Developers will be able to compare a feature branch's coverage against the main branch to see the direct impact of their changes, effectively enabling patch coverage analysis.

[**Read the detailed feature plan...**](./history-coverage.md)

---

### 3. Traceable Reports and Source Visualization
**Status: Planned**

When merging multiple reports (e.g., from unit and integration tests), it can be hard to know which tests cover which lines. This feature will tag each line of code with its source report. The final HTML report will allow users to see exactly which test suite is responsible for the coverage of any given line.

[**Read the detailed feature plan...**](./trace-reports.md)

---

### 4. Live Reloading with Watch Mode
**Status: Planned**

Introduce a `--watch` flag to provide a live-reloading development workflow. The tool will monitor input coverage reports for changes and automatically regenerate the output. This creates an instant feedback loop, perfect for integrating with IDEs and background test runners.

[**Read the detailed feature plan...**](./watch-mode.md)