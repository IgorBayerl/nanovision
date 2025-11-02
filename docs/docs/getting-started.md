# Getting Started

Welcome to Nanovision! This guide will walk you through the installation process and basic usage of the tool.

## Installation

Nanovision is a standalone Go binary, which makes installation simple.

1. **Download the Binary:**  
    Visit the [releases page](https://github.com/IgorBayerl/nanovision/releases) and download the appropriate binary for your operating system.

2. **Place it in Your PATH:**  
    For ease of use, place the downloaded binary in a directory that is part of your system's `PATH`.

## Basic Usage

Running Nanovision is straightforward. The most important flags are `--report` and `--output`.

```bash
nanovision --report="coverage.out" --sourcedir="." --reporttypes="Html"
```

Alternatively, if you have a `nanovision.yaml` configuration file in place:

```yaml
reports:
  - "coverage-unit.out"

source_dirs:
  - "."

report_types:
  - "Html"

output_dir: "reports"
```

You can simply run the command without any flags in the directory containing your configuration file:

```bash
nanovision
```
