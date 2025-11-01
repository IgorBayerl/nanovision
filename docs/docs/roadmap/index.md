# Project Roadmap

This section outlines the planned features and architectural enhancements for nanovision. Our goal is to build upon the solid foundation of ReportGenerator while introducing modern, Go-native features that improve developer workflows.

The following documents provide detailed implementation plans for major upcoming features.

---

### 1. Hierarchical and Decentralized Components
**Status: Planned**

Move away from a single, repository-wide coverage number and introduce a system of **components**. Teams will be able to define logical code boundaries (e.g., a microservice, a library) in a decentralized way, enabling granular insights, clear ownership, and autonomous configuration.

[**Read the detailed feature plan...**](./components.md)



---

### 2. Live Reloading with Watch Mode
**Status: Planned**

Introduce a `--watch` flag to provide a live-reloading development workflow. The tool will monitor input coverage reports for changes and automatically regenerate the output. This creates an instant feedback loop, perfect for integrating with IDEs and background test runners.

[**Read the detailed feature plan...**](./watch-mode.md)