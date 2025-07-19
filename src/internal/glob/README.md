# Package: glob

## Overview

This package provides an advanced file and directory path matching implementation, commonly known as "globbing". It was created to offer more powerful and flexible pattern matching than Go's standard `filepath.Glob` function.

While inspired by the capabilities of the original C# ReportGenerator, this implementation is a standalone Go package designed to be robust and easy to use. Its main purpose within this project is to allow users to easily specify one or multiple coverage report files using intuitive, powerful patterns.

## Features

This `glob` package supports the following features:

*   **Recursive Matching (`**`)**: Matches any number of subdirectories (including none). This is perfect for finding files deep within a project structure.
*   **Wildcard (`*`)**: Matches zero or more characters within a single file or directory name. It does not cross path separators (`/` or `\`).
*   **Single-Character Wildcard (`?`)**: Matches exactly one character in a file or directory name.
*   **Brace Expansion (`{a,b,...}`)**: Matches any of the comma-separated patterns provided inside the braces. This can be used for matching multiple names or extensions.
*   **Character Sets (`[...]`)**: Matches any single character within the set (e.g., `[abc]`) or range (e.g., `[0-9]`).
*   **Case-Insensitive Matching**: By default, all pattern matching is case-insensitive for a better user experience, especially on case-sensitive filesystems.
*   **Cross-Platform Path Separators**: The globber correctly handles both `/` and `\` as path separators in patterns, making it work seamlessly across different operating systems.

## Why a Custom `glob` Package?

Go's standard library provides `filepath.Glob`, but it is intentionally limited in its functionality. It lacks support for critical features like recursive directory matching (`**`) and brace expansion (`{...}`), which are standard in many modern development tools and build scripts.

To provide a powerful and user-friendly command-line interface where users can easily and intuitively specify files, a more advanced implementation was necessary. This custom package fills that gap by providing these advanced features, leading to a much better user experience.

## How It Works

The globber processes patterns by breaking them into path segments. It then recursively walks the filesystem, converting wildcard segments (`*`, `?`, `[]`) into cached regular expressions to efficiently match against file and directory names. This segment-by-segment approach allows it to handle complex patterns like `src/**/{cmd,internal}/*.go` effectively.

## Public API

The package exposes a `Glob` struct for fine-grained control and a convenient wrapper function for direct use.

#### `NewGlob(pattern string, fs filesystem.Filesystem, opts ...GlobOption) *Glob`

This is the main constructor for creating a `Glob` instance.

*   `pattern`: The glob pattern string to be matched.
*   `fs`: A `filesystem.Filesystem` interface. Pass `nil` to use the real operating system filesystem. This allows for mocking during tests.
*   `opts`: A variadic list of options to configure the globber's behavior.

#### `WithIgnoreCase(bool)`

An option to enable or disable case-insensitive matching. It is enabled by default.

```go
// Example: Create a case-sensitive globber
g := glob.NewGlob(pattern, nil, glob.WithIgnoreCase(false))
```

#### `g.ExpandNames() ([]string, error)`

This is the primary method on a `Glob` instance. It executes the matching logic and returns a slice of all absolute file paths that match the pattern.

#### `GetFiles(pattern string) ([]string, error)`

This is the main convenience function that most of the application will use. It creates a `Glob` instance with default settings, expands the pattern, and returns the results.

```go
// The most common way to use the package:
files, err := glob.GetFiles("src/**/*.go")
```

## Usage Examples

Assuming a directory structure like:

```
/app/
├── project-a/
│   └── bin/
│       ├── coverage.xml
│       └── debug.log
├── project-b/
│   └── bin/
│       ├── coverage.xml
│       └── TestResults/
│           └── coverage.cobertura.xml
└── shared/
    └── utils.cs
```

**Example 1: Basic Wildcard**
Find all `.xml` files in a specific directory.

```go
// Pattern: "/app/project-a/bin/*.xml"
files, err := glob.GetFiles("/app/project-a/bin/*.xml")
// Returns: ["/app/project-a/bin/coverage.xml"]
```

**Example 2: Recursive Directory Search (`**`)**
Find all `coverage.xml` files anywhere under the `/app` directory.

```go
// Pattern: "/app/**/coverage.xml"
files, err := glob.GetFiles("/app/**/coverage.xml")
// Returns: [
//   "/app/project-a/bin/coverage.xml",
//   "/app/project-b/bin/coverage.xml"
// ]
```

**Example 3: Brace Expansion (`{}`)**
Target specific projects' coverage files.

```go
// Pattern: "/app/{project-a,project-b}/bin/coverage.xml"
files, err := glob.GetFiles("/app/{project-a,project-b}/bin/coverage.xml")
// Returns: [
//   "/app/project-a/bin/coverage.xml",
//   "/app/project-b/bin/coverage.xml"
// ]
```

**Example 4: Combining Patterns**
Find any `.xml` or `.cobertura.xml` file within any `bin` or `TestResults` folder under `/app`.

```go
// Pattern: "/app/**/{bin,TestResults}/*.{xml,cobertura.xml}"
files, err := glob.GetFiles("/app/**/{bin,TestResults}/*.{xml,cobertura.xml}")
// Returns: [
//   "/app/project-a/bin/coverage.xml",
//   "/app/project-b/bin/coverage.xml",
//   "/app/project-b/bin/TestResults/coverage.cobertura.xml"
// ]
```

## Integration in ReportGenerator

In this project, the `glob` package is used in `cmd/main.go` to resolve the file paths provided by the user via the `-report` command-line flag. This enables users to provide flexible and powerful patterns to select one or more coverage reports for processing.