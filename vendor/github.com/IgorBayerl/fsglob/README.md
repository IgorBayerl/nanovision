# fsglob

[![Go Reference](https://pkg.go.dev/badge/github.com/IgorBayerl/fsglob.svg)](https://pkg.go.dev/github.com/IgorBayerl/fsglob)
[![Go Report Card](https://goreportcard.com/badge/github.com/IgorBayerl/fsglob)](https://goreportcard.com/report/github.com/IgorBayerl/fsglob)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](https://opensource.org/licenses/MIT)

`fsglob` is a Go package for finding file paths that match glob patterns. It operates by walking a filesystem and is designed to be cross-platform and extensible.

## Features

-   **Standard Glob Patterns:** Supports common glob syntax for flexible matching.
-   **Cross-Platform Path Handling:** Automatically handles both Windows (`\`) and Unix (`/`) path separators in patterns.
-   **Extensible Filesystem:** Can operate on any filesystem that implements the provided `filesystem.Filesystem` interface, making it suitable for testing with mocks or for use with virtual filesystems.

## Installation

```bash
go get github.com/your-username/fsglob
```

## Quick Start

The main function in this package is `fsglob.GetFiles()`. It takes a pattern and returns a slice of all matching file and directory paths.

```go
package main

import (
	"fmt"
	"log"

	"github.com/your-username/fsglob"
)

func main() {
	// Example 1: Find all Go files recursively from the current directory.
	goFiles, err := fsglob.GetFiles("**/*.go")
	if err != nil {
		log.Fatalf("Failed to find Go files: %v", err)
	}

	fmt.Println("Found Go files:")
	for _, file := range goFiles {
		fmt.Println(file)
	}

	// Example 2: Find all Markdown or text files in a 'docs' directory.
	docFiles, err := fsglob.GetFiles("docs/*.{md,txt}")
	if err != nil {
		log.Fatalf("Failed to find doc files: %v", err)
	}

	fmt.Println("\nFound documentation files:")
	for _, file := range docFiles {
		fmt.Println(file)
	}
}
```

## Pattern Matching Details

`fsglob` supports the following pattern syntax:

| Pattern | Description                                                               | Example                  |
| :------ | :------------------------------------------------------------------------ | :----------------------- |
| `*`     | Matches any sequence of characters, except for path separators (`/` or `\`). | `*.log`                  |
| `?`     | Matches any single character.                                             | `file?.txt`              |
| `**`    | Matches zero or more directories, files, and subdirectories recursively.  | `reports/**/*.xml`       |
| `[]`    | Matches any single character within the brackets. Can be a set or a range. | `[abc].go`, `[0-9].txt`  |
| `{}`    | Matches any of the comma-separated patterns within the braces.            | `image.{jpg,png,gif}`    |

## Advanced Usage: Custom Filesystem

For testing or working with virtual filesystems (e.g., in-memory, TAR files), you can use `fsglob.NewGlob()` to create a globber instance that operates on a custom filesystem.

Your custom filesystem must implement the `filesystem.Filesystem` interface, which is exposed by `github.com/your-username/fsglob/filesystem`.

```go
package main

import (
	"fmt"
	"log"

	"github.com/your-username/fsglob"
	"github.com/your-username/fsglob/filesystem"
)

// InMemoryFS is a simple in-memory implementation of filesystem.Filesystem.
// (This is a conceptual example; a full implementation is required).
type InMemoryFS struct {
    // ... fields to store files and directories in memory
}

// Implement the methods of the filesystem.Filesystem interface for InMemoryFS...
func (fs *InMemoryFS) Stat(name string) (fs.FileInfo, error) { /* ... */ }
func (fs *InMemoryFS) ReadDir(name string) ([]fs.DirEntry, error) { /* ... */ }
func (fs *InMemoryFS) Getwd() (string, error) { /* ... */ }
func (fs *InMemoryFS) Abs(path string) (string, error) { /* ... */ }
// ... other methods

func main() {
	// 1. Create an instance of your custom filesystem.
	memFS := &InMemoryFS{}
	// ... code to populate memFS with files and directories ...

	// 2. Create a new Glob instance with the custom filesystem.
	g := fsglob.NewGlob("/**/*.log", memFS)

	// 3. Use the Expand() method to find matches.
	matches, err := g.Expand()
	if err != nil {
		log.Fatalf("Error during glob expansion: %v", err)
	}

	fmt.Println("Log files found in the in-memory filesystem:")
	for _, match := range matches {
		fmt.Println(match)
	}
}
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.