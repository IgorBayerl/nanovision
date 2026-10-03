package config

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
)

// FileNames are the names a config file can have. The root config and the
// config of a sub-folder use the same names. The first one is the name the
// documentation uses.
var FileNames = []string{"nanovision.yaml", "nanovision.yml", ".nanovision.yaml", ".nanovision.yml"}

// IsFileName reports whether name is one of FileNames.
func IsFileName(name string) bool {
	return slices.Contains(FileNames, name)
}

// FindFile returns the path of the config file in dir, or "" when dir has
// none. Two config files in one folder are an error: the tool cannot know
// which one the user means.
func FindFile(dir string) (string, error) {
	found := ""
	for _, name := range FileNames {
		p := filepath.Join(dir, name)
		if info, err := os.Stat(p); err != nil || info.IsDir() {
			continue
		}
		if found != "" {
			return "", sameFolderError(found, p)
		}
		found = p
	}
	return found, nil
}

func sameFolderError(a, b string) error {
	return fmt.Errorf("%s and %s are both config files of the same folder; keep one", a, b)
}
