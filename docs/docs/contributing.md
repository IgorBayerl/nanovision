# How to Contribute

We actively welcome contributions to nanovision! Whether it's adding a new parser, improving the documentation, or fixing a bug, your help is valued.

## Project Architecture

nanovision is built with a modular design to enable easy expansion. The core pipeline is:

1.  **Parsers (`internal/parsers`):** Convert various input formats (e.g., Cobertura, GoCover) into a standardized intermediate `model`.
2.  **Analyzer (`internal/analyzer`):** Merges one or more parsed results into a single, unified summary. This is where filtering and data enrichment happens.
3.  **Reporters (`internal/reporter`):** Takes the final summary and generates output formats like HTML, Text, or LCOV.

## Setting Up Your Development Environment

1.  Clone the repository:
    ```bash
    git clone https://github.com/IgorBayerl/nanovision.git
    ```
2.  Install Go (version 1.23 or higher).
3.  Install dependencies:
    ```bash
    go mod tidy
    ```
4.  Run the tests:
    ```bash
    go test ./...
    ```

## How to Add a New Parser

*(This would be a great place to add a detailed guide on the parser interface and how to implement it, as you planned).*