# Go Report Generator

This project is a Go port of the excellent [ReportGenerator](https://github.com/danielpalme/ReportGenerator) tool by Daniel Palme. It aims to provide a fast, dependency-free, and cross-platform code coverage reporting tool with a focus on modern development workflows.

## Motivation

The original ReportGenerator is a mature and feature-rich tool that has been a staple in the .NET ecosystem for over 14 years. The motivation for this Go port is threefold:
1.  **For Study:** To delve into the complexities of code coverage formats and report generation, providing a significant learning opportunity.
2.  **Dependency-Free:** To create a single, native executable that runs on any platform without requiring external runtimes like .NET, making it ideal for containerized CI/CD pipelines.
3.  **Extensibility:** To build a modern foundation in Go that can be easily extended to support new formats and features tailored to current needs, such as first-class support for Go's own coverage tooling.

## System Design

The Go Report Generator is designed with modularity and extensibility in mind, drawing inspiration from the original's architecture while leveraging Go's strengths.

*   **Entrypoint (`cmd/main.go`):** The main application entry point handles command-line argument parsing and orchestrates the report generation workflow. It is also responsible for composing the application's dependencies, such as creating parser and language processor instances and injecting them into their respective factories.
*   **Configuration (`internal/reportconfig`, `internal/settings`):** A dedicated package for managing all configuration, from input files and output directories to filters and report types.
*   **Parser Factory (`internal/parsers`):** At the core of the input processing is a parser factory. It is explicitly given a list of all available parser implementations (e.g., for Cobertura, GoCover). For a given input file, the factory iterates through these parsers and uses the first one that reports it can handle the format. This makes it easy to add support for new formats by simply creating a new parser and adding it to the factory's constructor in `cmd/main.go`.
*   **Language-Specific Processing (`internal/language`):** To handle language-specific nuances like formatting C#/.NET identifiers or calculating Go-specific metrics, a language processor factory is used. Similar to the parser factory, it holds a list of available language processors. When processing a source file, it selects the appropriate processor by checking which one supports the file's extension (e.g., `.cs`, `.go`). This allows for clean separation of concerns, where parsers focus on format and processors focus on language semantics.
*   **Intermediate Model (`internal/model`):** Once a report is parsed, its data is translated into a standardized intermediate model. This decouples the input formats from the output reporters, meaning any supported input format can be used to generate any supported output format.
*   **Analysis Engine (`internal/analyzer`):** The analyzer takes the results from one or more parsers and merges them into a single, unified `SummaryResult`. This is what enables the powerful feature of combining coverage reports from different test runs or even different languages (e.g., C# and Go) into one consolidated report.
*   **Report Builders (`internal/reporter`):** Report builders are responsible for generating the final output. The `htmlreport` package, for instance, generates a sophisticated single-page application (SPA).
*   **Angular SPA Frontend (`angular_frontend_spa`):** The HTML report is not just a static file. It's a full-fledged Angular application that provides rich, interactive features like real-time filtering, sorting, and collapsible views, offering a much more dynamic user experience than traditional reports.

## Feature Status

This project is a work in progress. The following table tracks the implementation status of features found in the original ReportGenerator and new features specific to this Go port.

| Feature Category | Feature | C# Status | Go Status | Notes |
| :--- | :--- | :---: | :---: | :--- |
| **Input Formats** | **Cobertura** | ✅ | ✅ | Core format, fully supported. |
| | **Go Cover** | ❌ | ✅ | **Go-native feature.** Direct parsing of `coverage.out`. |
| | OpenCover | ✅ | ❌ | |
| | dotCover | ✅ | ❌ | |
| | Visual Studio | ✅ | ❌ | |
| | JaCoCo | ✅ | ❌ | |
| | Clover | ✅ | ❌ | |
| | NCover | ✅ | ❌ | |
| | Multiple/Merged Reports | ✅ | ✅ | Merging is a core feature of the `analyzer`. |
| **Output Formats** | **HTML (SPA)** | ✅ | ✅ | Go version generates a modern Angular-based SPA. |
| | **TextSummary** | ✅ | ✅ | |
| | **lcov** | ✅ | ✅ | |
| | Badge | ✅ | ❌ | |
| | CodeClimate | ✅ | ❌ | |
| | Cobertura | ✅ | ❌ | |
| | CsvSummary | ✅ | ❌ | |
| | HtmlChart | ✅ | ❌ | |
| | HtmlInline | ✅ | ❌ | |
| | HtmlSummary | ✅ | ❌ | |
| | JsonSummary | ✅ | ❌ | |
| | Latex | ✅ | ❌ | |
| | MHtml | ✅ | ❌ | |
| | PngChart | ✅ | ❌ | |
| | SvgChart | ✅ | ❌ | |
| | TeamCitySummary | ✅ | ❌ | |
| | Xml | ✅ | ❌ | |
| | XmlSummary | ✅ | ❌ | |
| **Core Features** | **Filtering** (Assembly, Class, File) | ✅ | ✅ | Filtering logic is implemented. |
| | **Branch Coverage** | ✅ | ✅ | Supported for formats that provide it (e.g., Cobertura). |
| | **Method Coverage** | ✅ | ✅ | |
| | **Cyclomatic Complexity** | ✅ | ✅ | **Go-native support added.** C# support not ported yet. |
| | History / Trend Charts | ✅ | ❌ | Historic coverage tracking is not yet implemented. |
| | Risk Hotspots | ✅ | ❌ | Risk hotspot analysis based on metrics is not yet implemented. |
| | Raw Mode (No class merging) | ✅ | ❌ | |

## Command Line Arguments

The command-line interface aims to be compatible with the original ReportGenerator.

| Argument | C# Status | Go Status | Go Flag Name | Description |
| :--- | :---: | :---: | :--- | :--- |
| `reports` | ✅ | ✅ | `report` | The coverage reports that should be parsed. |
| `targetdir` | ✅ | ✅ | `output` | The directory where the generated report should be saved. |
| `sourcedirs` | ✅ | ✅ | `sourcedirs` | Optional directories which contain the source code. |
| `reporttypes` | ✅ | ✅ | `reporttypes` | The output formats to generate. |
| `assemblyfilters` | ✅ | ✅ | `assemblyfilters` | Filters for assemblies to include or exclude. |
| `classfilters` | ✅ | ✅ | `classfilters` | Filters for classes to include or exclude. |
| `filefilters` | ✅ | ✅ | `filefilters` | Filters for files to include or exclude. |
| `verbosity` | ✅ | ✅ | `verbosity` | The verbosity level of the log messages. |
| `tag` | ✅ | ✅ | `tag` | Optional tag or build version. |
| `title` | ✅ | ✅ | `title` | Optional report title. |
| `historydir` | ✅ | ❌ | `-` | Directory for storing persistent coverage information. |
| `plugins` | ✅ | ❌ | `-` | Plugin files for custom reports or history storage. |
| `riskhotspotassemblyfilters`| ✅ | ✅ | `riskhotspotassemblyfilters` | Assembly filters for risk hotspots. |
| `riskhotspotclassfilters`| ✅ | ✅ | `riskhotspotclassfilters` | Class filters for risk hotspots. |
| `license`| ✅ | ❌ | `-` | License for PRO version features. |

## How to Contribute

This project is in its early stages, and contributions are welcome! Whether it's porting a feature, adding a new parser, or improving documentation, your help is appreciated.

### A Note on Feature Parity

ReportGenerator is a 14-year-old project with a rich feature set. This Go port is only a few months old. If you need a feature from the original that has not yet been ported, please **open an issue on GitHub**.

To help accelerate the process, please include the following in your feature request:

1.  A clear description of the feature.
2.  A link to the original feature's documentation or related issue if possible.
3.  An example of the command-line usage.
4.  Sample input (coverage files) and expected output, if applicable.

This will make it much easier to understand the requirements and implement the feature correctly.

Thank you for your interest and support