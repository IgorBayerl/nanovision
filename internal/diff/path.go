package diff

import (
	"path/filepath"
	"strings"
)

// Normalize converts a diff path into a normalized form:
// - Converts backslashes to forward slashes first
// - Strips "a/" and "b/" prefixes (git-style) - only once
// - Cleans the path (removes "." components, etc.)
func Normalize(p string) string {
	// Convert to forward slashes first.
	// We use strings.ReplaceAll instead of filepath.ToSlash because valid diffs
	// might contain Windows-style paths even when the tool is running on Linux.
	p = strings.ReplaceAll(p, "\\", "/")

	// Strip git-style prefixes (only check once at the start)
	if strings.HasPrefix(p, "a/") {
		p = strings.TrimPrefix(p, "a/")
	} else if strings.HasPrefix(p, "b/") {
		p = strings.TrimPrefix(p, "b/")
	}

	// Clean the path to remove "." and "./" components
	p = filepath.Clean(p)

	// filepath.Clean might add "./" for relative paths or convert back to backslashes
	// on Windows, so we force forward slashes again and remove leading "./".
	p = filepath.ToSlash(p)
	for strings.HasPrefix(p, "./") {
		p = strings.TrimPrefix(p, "./")
	}

	return p
}
