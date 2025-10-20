# Build Configuration

You may have noticed our project contains multiple Vite configuration files. This setup is intentional and necessary to meet two distinct deployment requirements for the generated report.

## The Core Requirement: Web vs. Local

We must produce a report that works in two environments:

2.  **Local Filesystem (`file://`)**: To embed in the CLI tool, we will be able to open this without a web server, double click the .html file and it just works.
1.  **Web Server (`http://`)**: The standard environment for hosting on a website (Not Used Yet).

Vite's default build uses modern JavaScript Modules (`<script type="module">`). For security (CORS), browsers block these modules when an HTML file is opened locally via `file://`.

To support the local filesystem, we must generate a "portable" version of the report that uses classic scripts (`<script defer>`), which browsers allow on `file://`.

| Script Type        | Loaded As                | Behavior on `file://` |
|:-------------------|:-------------------------|:----------------------|
| **Modern Module**  | `<script type="module">` | **Blocked**           |
| **Classic Script** | `<script defer>`         | **Allowed**           |

## Build Targets

We have two build commands to generate these different versions.

### 1. Modern Build

*   **Command:** `pnpm build`
*   **Output:** `dist-modern/`
*   **Use Case:** Deploying the report to any standard web server.
*   **Technology:** Uses `<script type="module">` and efficient code-splitting. **Will not work on `file://`**.

### 2. Portable Build

*   **Command:** `pnpm build:portable`
*   **Output:** `dist/`
*   **Use Case:** Creating a downloadable report where every page works locally via `file://`.
*   **Technology:** Uses classic `<script defer>`. To bypass bundler limitations with this format, it runs a separate build for each HTML entry point.

## Vite Configuration Files

This two-target approach explains our file structure.

#### Modern Build

*   `vite.config.modern.ts`: A single config that builds all pages at once using modern defaults.

#### Portable Build

*   `vite.config.portable.base.ts`: Contains shared logic for the portable build, including the custom plugin to convert scripts to the classic format.
*   `vite.config.portable.index.ts`: Config for the first step, building *only* `index.html`.
*   `vite.config.portable.details.ts`: Config for the second step, building *only* `details.html`.

## Summary of Commands

| To Achieve This...                      | Run This Command...     | Output Directory      |
|:----------------------------------------|:------------------------|:----------------------|
| Create the report for a **web server**. | `pnpm build`            | `dist-modern/`        |
| Create the report for **local use**.    | `pnpm build:portable`   | `dist/`               |
| Preview the **modern build** locally.   | `pnpm preview`          | Serves `dist-modern/` |
| Preview the **portable build** locally. | `pnpm preview:portable` | Serves `dist/`        |

> **Important**: The `preview` command always uses a web server (`http://`). To test the `file://` functionality of the portable build, you must navigate to the `dist` folder in your file explorer and **double-click the HTML files** directly.