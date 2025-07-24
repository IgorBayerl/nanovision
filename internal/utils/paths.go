package utils

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Stater interface {
	Stat(name string) (fs.FileInfo, error)
}

type DefaultStater struct{}

func (ds DefaultStater) Stat(name string) (fs.FileInfo, error) {
	return os.Stat(name)
}

// FindFileInSourceDirs resolves the often-inconsistent file paths found in
// coverage reports against a list of local source code directories.
//
// # Why this function is necessary
//
// Coverage reports are frequently generated in one environment (e.g., a CI/CD pipeline)
// but viewed in another (e.g., a developer's local machine). This leads to several
// common problems that this function is designed to solve:
//
//	Absolute Path Mismatch: A report from a Linux CI server might contain the path
//	`/home/runner/work/my-project/src/app.go`. This path is meaningless on a developers
//	Windows machine where the code lives at `C:\Users\dev\my-project`.
//
//	Partial Path Suffixes: Some coverage tools only record paths relative to the
//	project root (e.g., `services/user.go`), not the full repository path. The function
//	must be able to find the file even if the report path is just a suffix of the real path.
//
// By providing the local `sourceDirs`, a user tells the tool where to start searching.
// This function then intelligently combines the reports path with the local directories
// to robustly locate the correct source file on the current filesystem.
//
// # Examples
//
// Example: Resolving a CI path on a local machine
//
//	// Path from a report generated on a Linux CI server
//	relativePath := "/home/runner/work/adler-cov/internal/analyzer/analyzer.go"
//
//	// The user's local source code directory on Windows
//	sourceDirs := []string{"C:\\Users\\igor\\dev\\adler-cov"}
//
//	// The function will correctly find the file at:
//	// "C:\\Users\\igor\\dev\\adler-cov\\internal\\analyzer\\analyzer.go"
//	// by trying various suffixes of the relativePath within the source directory.
//
// Example : Resolving a partial path from a .NET report
//
//	// Path from a Cobertura report, relative to the .csproj file
//	relativePath := "Services/UserService.cs"
//
//	// The local source directory is the root of the API project
//	sourceDirs := []string{"C:\\dev\\my-api\\src\\MyProject.Api"}
//
//	// The function will find the file by directly joining the paths:
//	// "C:\\dev\\my-api\\src\\MyProject.Api\\Services\\UserService.cs"
// In: internal/utils/paths.go

func FindFileInSourceDirs(relativePath string, sourceDirs []string, stater Stater) (string, error) {
	// Clean the incoming relative path using the host OS's rules.
	// This will correctly handle both '/' and '\' on Windows.
	cleanedRelativePath := filepath.Clean(relativePath)

	// First, check if the cleaned path is absolute and exists.
	if filepath.IsAbs(cleanedRelativePath) {
		if _, err := stater.Stat(cleanedRelativePath); err == nil {
			return cleanedRelativePath, nil
		}
	}

	// Iterate through the user-provided source directories.
	for _, dir := range sourceDirs {
		// Clean the source directory path.
		cleanedDir := filepath.Clean(dir)

		// Attempt to join the source directory with the full relative path.
		// filepath.Join will use the correct OS-specific separator.
		potentialPath := filepath.Join(cleanedDir, cleanedRelativePath)
		if _, err := stater.Stat(potentialPath); err == nil {
			return potentialPath, nil
		}

		// If that fails, try to find a match by using suffixes of the relative path.
		// This handles cases where the report path includes extra parent directories
		// (e.g., /build/src/my/project/file.go) and the source dir is my/project.
		pathParts := strings.Split(cleanedRelativePath, string(filepath.Separator))
		for i := 1; i < len(pathParts); i++ { // Start at 1 to skip the first part
			suffixToTry := filepath.Join(pathParts[i:]...)
			potentialPathWithSuffix := filepath.Join(cleanedDir, suffixToTry)
			if _, err := stater.Stat(potentialPathWithSuffix); err == nil {
				return potentialPathWithSuffix, nil
			}
		}
	}

	return "", fmt.Errorf("file %q not found in any source directory", relativePath)
}
