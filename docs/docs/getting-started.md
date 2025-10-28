# Getting Started

Welcome to nanovision! This guide will walk you through the installation and basic usage of the tool.

## Installation

nanovision is a standalone Go binary, which makes installation simple.

1.  **Download the Binary:**
    Head over to the [releases page](https://github.com/IgorBayerl/nanovision/releases) and download the appropriate binary for your operating system.

2.  **Place it in your PATH:**
    For ease of use, place the downloaded binary in a directory that is part of your system's `PATH`.

## Basic Usage

Running nanovision is straightforward. The most important flags are `--report` and `--output`.

```bash
nanovision --report="coverage.out" --output="coverage-report"
```