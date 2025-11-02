# Contributing to **nanovision**

Thanks for taking the time to contribute! 

This project is a **Go CLI** with a **React** subproject for Html reports. We provide cross‑platform scripts, a devcontainer, and end‑to‑end tests to keep contributions smooth and consistent.


## Code of Conduct

This project follows a [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you agree to abide by it.


## Project Layout

```
/                  
├─ .devcontainer/ # devcontainer config
├─ cmd/           # CLI entry points
├─ internal/      # private Go packages
├─ ui/            # React subproject (used for Html reports)
├─ scripts/       # helper scripts (build, test, e2etest, ...)
├─ vendor/        
├─ README.md
└─ ...
```

A simplified diagram that shows the main parts of the system.

1. **Configuration loading:**
  - CLI arguments will overwrite the settings in the config file.
2. **Report parsing:**
  - Collects:
    - Line coverage.
    - Branch coverage.
    - List of files touched by the reports.
3. **Enricher:**
  - Uses **tree-sitter** for static analysis of each covered file from the source code, and collects some extra information.
    - Function/method start and finish lines.
    - Cyclomatic complexity.
    - Other metrics as applicable.
  - Each supported programming language has its own analyzer with potentially different features.

4. **Reporter:**
  - Uses the collected information from the previous steps to generate the requested reports.


![diagram](docs/docs/imgs/nanovision-flow.png)


## Before You Start

* Search existing **issues** and **PRs** to avoid duplicates.
* For non‑trivial changes, open an **issue** or **discussion** first to align on direction.
* Make sure you can **build and test** locally (see below) or use the **devcontainer**.

## Development Setup

> Native setup works on macOS/Linux/Windows. Use the devcontainer if you prefer an pre-configured environment.

* **Go:** ≥ **1.25**
* **Python:** ≥ **3.x**
* **Node.js:** ≥ **20.x**

  * Only for the React project, not necessary if you are not touching the UI side.
* **Package manager (web):** `pnpm` (preferred) or `npm`/`yarn`
* **C/C++ toolchain (for Windows build only):** MinGW-w64 on `PATH` providing `x86_64-w64-mingw32-gcc`/`g++`.

  * On **Windows**, use MSYS2 MinGW-w64.
  * On **macOS/Linux** (cross-compiling to Windows), use the MinGW-w64 cross-compiler.
  * **Linux build (`--linux`)**: no C/C++ toolchain needed unless you enable CGO yourself.

```bash
# Clone
git clone https://github.com/IgorBayerl/nanovision.git
cd nanovision

# Build
python ./scripts/build.py --windows 
python ./scripts/build.py --linux 

# Generates many reports based on the example projects
python ./scripts/e2e_test.py 
```

---

## Using the Devcontainer

We ship a **.devcontainer** configuration (VS Code compatible) so contributors can get a ready‑to‑run environment.

1. Install **Docker** and **VS Code** with the **Dev Containers** extension.
2. Open the repository in VS Code. 
  - VSCode should prompt you with **Reopen in Container**.
  - Or press `ctrl + shift + p` and run **Dev Containers: Build and Reopen in Containers**
3. The container comes with Go, Node, Python, C Compiler and tooling preinstalled (see `.devcontainer/devcontainer.json`).

> CI uses a similar toolchain to the devcontainer to reduce "works on my machine" issues.

---

## Scripts

Project automation lives under `scripts/`:

- **Build (Linux):** `scripts/build.py --linux`
- **Build (Windows):** `scripts/build.py --windows`
- **Unit tests:** `scripts/test.py`
- **E2E tests:** `scripts/e2d_test.py` - runs many example reports through the tool
- **E2E tests:** `scripts/e2d_test.py -sc` - collects code coverage during the run and generate a self coverage report based on the e2e tests.

## Testing

### Unit tests

```bash
python ./scripts/test.py
```

### e2e tests
```bash
python ./scripts/e2e_test.py -sc
  ```

## Continuous Integration (GitHub Actions)

**GitHub Actions** is used for CI/CD. There are two workflows in `.github/workflows/`:

### [PR Check](https://github.com/IgorBayerl/nanovision/actions/workflows/pr_checks.yml)

**Trigger:** Pull Requests

**What it does:**

- Runs a **try build** of the **docs** and the **self report** (no publishing)
- Executes relevant checks/tests (Go unit, web tests, E2E if configured)
- Uploads build **artifacts** for PR reviewers where applicable

**Intent:** Fast feedback to maintainers and contributors before merge.

### [Publish Docs and Self Report](https://github.com/IgorBayerl/nanovision/actions/workflows/pages.yml)

**Trigger:** Every **push to `main`**

**What it does:**

- Performs the same **build** as `pr_check`
- **Publishes** the built **docs** and **self report** to **GitHub Pages**


## Security

Report vulnerabilities privately to `dev.igorbayerl@gmail.com`. Do not open public issues for security problems.

<!-- 
TODO: Create release process
## Release Process

- Merge PRs into `main` when CI is green
- Update `CHANGELOG.md` (use Keep a Changelog or auto‑generated notes)
- Tag release: `vX.Y.Z`
- Build release artifacts for Linux/Windows (consider `goreleaser`) and upload
- Announce in Releases and/or community channels -->


## License

By contributing, you agree your contributions are licensed under the [Apache 2.0 Liscense](LICENSE) 

## Contact

- Discussions: `open an issue for open discussions`
- Email: `dev.igorbayerl@gmail.com`



### Quick Start (TLDR)

1. Fork and clone the repo
2. Setup locally or open in the **devcontainer**
3. Make your changes
4. Build and Test: `./scripts/e2e_test.py -sc`
5. Open a PR with a clear description
