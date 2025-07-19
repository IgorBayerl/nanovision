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

// Attempts to locate a file using a Stater interface.
func FindFileInSourceDirs(relativePath string, sourceDirs []string, stater Stater) (string, error) {
	if filepath.IsAbs(relativePath) {
		if _, err := stater.Stat(relativePath); err == nil {
			return relativePath, nil
		}
	}

	cleanedRelativePath := filepath.Clean(relativePath)

	for _, dir := range sourceDirs {
		absPath := filepath.Join(filepath.Clean(dir), cleanedRelativePath)
		if _, err := stater.Stat(absPath); err == nil {
			return absPath, nil
		}

		pathParts := strings.Split(cleanedRelativePath, string(os.PathSeparator))
		for i := 0; i < len(pathParts); i++ {
			suffixToTry := filepath.Join(pathParts[i:]...)
			potentialPath := filepath.Join(filepath.Clean(dir), suffixToTry)
			if _, err := stater.Stat(potentialPath); err == nil {
				return potentialPath, nil
			}
		}
	}
	return "", fmt.Errorf("file %q not found in any source directory (%v) or as absolute path", relativePath, sourceDirs)
}
