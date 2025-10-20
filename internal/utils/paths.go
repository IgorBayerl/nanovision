package utils

import (
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/IgorBayerl/AdlerCov/filereader"
)

// walk is a local, recursive helper that mimics filepath.WalkDir but uses our mockable filereader.Reader interface.
func walk(reader filereader.Reader, root string, fileNameToFind string, logger *slog.Logger) (string, error) {
	logger.Debug("Recursive search: walking directory", "dir", root)
	entries, err := reader.ReadDir(root)
	if err != nil {
		logger.Debug("Recursive search: cannot read directory, stopping this path.", "dir", root, "error", err)
		return "", nil // Not a fatal error, just a path that can't be explored.
	}

	for _, entry := range entries {
		if !entry.IsDir() && entry.Name() == fileNameToFind {
			foundPath := filepath.Join(root, entry.Name())
			logger.Debug("Recursive search: MATCH FOUND!", "path", foundPath)
			return foundPath, nil
		}
	}

	for _, entry := range entries {
		if entry.IsDir() {
			path := filepath.Join(root, entry.Name())
			foundPath, err := walk(reader, path, fileNameToFind, logger)
			if err != nil {
				return "", err // Propagate a fatal error.
			}
			if foundPath != "" {
				return foundPath, nil // Propagate the successful result immediately.
			}
		}
	}

	return "", nil
}

// FindFileInSourceDirs resolves a file path from a report against a list of source directories.
func FindFileInSourceDirs(relativePath string, sourceDirs []string, reader filereader.Reader, logger *slog.Logger) (string, error) {
	// If a nil logger is passed, default to a discarded one to prevent panics.
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}

	// Normalize path separators: replace backslashes with forward slashes
	// This handles Windows-style paths in coverage reports even when running on Unix
	// (coverage tools sometimes emit Windows paths like "src\app.go" even in cross-platform reports)
	normalizedRelativePath := strings.ReplaceAll(relativePath, "\\", "/")
	fileNameToFind := filepath.Base(normalizedRelativePath)
	logger.Debug("FindFileInSourceDirs starting", "relativePath", relativePath, "normalizedPath", normalizedRelativePath, "sourceDirs", sourceDirs)

	// Strategy 1: Absolute Path Check
	// Check both the original path and normalized path in case the report contains an absolute path
	for _, pathToCheck := range []string{relativePath, normalizedRelativePath} {
		if filepath.IsAbs(pathToCheck) {
			logger.Debug("Strategy 1: Path is absolute, checking existence.", "path", pathToCheck)
			if _, err := reader.Stat(pathToCheck); err == nil {
				logger.Debug("Strategy 1: Success.", "foundPath", pathToCheck)
				return pathToCheck, nil
			}
		}
	}

	// Strategy 2: Direct Join
	osRelativePath := filepath.FromSlash(normalizedRelativePath)
	for _, dir := range sourceDirs {
		potentialPath := filepath.Join(filepath.Clean(dir), osRelativePath)
		logger.Debug("Strategy 2: Trying direct join.", "path", potentialPath)
		if _, err := reader.Stat(potentialPath); err == nil {
			logger.Debug("Strategy 2: Success.", "foundPath", potentialPath)
			return potentialPath, nil
		}
	}

	// Strategy 3: Suffix Matching
	pathParts := strings.Split(normalizedRelativePath, "/")
	if len(pathParts) > 1 {
		for i := 1; i < len(pathParts); i++ {
			suffix := strings.Join(pathParts[i:], "/")
			osSuffix := filepath.FromSlash(suffix)
			for _, dir := range sourceDirs {
				potentialPath := filepath.Join(filepath.Clean(dir), osSuffix)
				logger.Debug("Strategy 3: Trying suffix join.", "path", potentialPath)
				if _, err := reader.Stat(potentialPath); err == nil {
					logger.Debug("Strategy 3: Success.", "foundPath", potentialPath)
					return potentialPath, nil
				}
			}
		}
	}

	// Strategy 4: Recursive Fallback Search
	logger.Debug("Strategy 4: Starting recursive search.", "filename", fileNameToFind)
	for _, dir := range sourceDirs {
		foundPath, err := walk(reader, filepath.Clean(dir), fileNameToFind, logger)
		if err != nil {
			return "", fmt.Errorf("error during recursive search in '%s': %w", dir, err)
		}
		if foundPath != "" {
			logger.Debug("Strategy 4: Success.", "foundPath", foundPath)
			return foundPath, nil
		}
	}

	logger.Warn("All strategies failed to find file.", "relativePath", relativePath)
	return "", fmt.Errorf("file %q not found in any of the provided source directories", relativePath)
}