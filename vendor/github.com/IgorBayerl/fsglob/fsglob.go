// Package fsglob provides a powerful and extensible library for globbing, finding file
// paths that match standard Unix-style glob patterns.

// The main entrypoint is the GetFiles function, which provides a simple interface for
// common use cases. The package is designed to be cross-platform, automatically
// handling both forward and backslashes in patterns, and can be extended to work
// with any virtual filesystem.

// Features:
//   - Standard Glob Syntax: Supports `*`, `?`, `**`, `[]`, and `{}`.
//   - Cross-Platform: Correctly handles path separators on both Windows and Unix.
//   - Extensible: Operates on any filesystem that implements the filesystem.Filesystem
//     interface, perfect for testing or use with virtual filesystems like afero.
//   - Case-Insensitive Option: Provides an option for case-insensitive matching.

// Quick Start

// To find all Go files recursively from the current directory:
//     goFiles, err := fsglob.GetFiles("**/*.go")
//     if err != nil {
//         // Handle error
//     }
//     for _, file := range goFiles {
//         fmt.Println(file)
//     }

// Pattern Matching

// The following pattern syntax is supported:

// | Pattern | Description                                                               | Example                  |
// | :------ | :------------------------------------------------------------------------ | :----------------------- |
// | `*`     | Matches any sequence of characters, except for path separators.           | `*.log`                  |
// | `?`     | Matches any single character.                                             | `file?.txt`              |
// | `**`    | Matches zero or more directories recursively.                             | `reports/**/*.xml`       |
// | `[]`    | Matches any single character within the brackets (e.g., [abc], [a-z]).    | `[a-c].go`               |
// | `{}`    | Brace expansion matches any of the comma-separated patterns.              | `image.{jpg,png,gif}`    |

// Advanced Usage: Custom Filesystem and Options

// For more control, such as using a custom filesystem for testing or enabling
// case-insensitive matching, use NewGlob.

//     // Create a new globber with case-insensitivity enabled.
//     // On Windows, matching is case-insensitive by default for non-wildcard patterns.
//     g := fsglob.NewGlob("reports/*.{XML,json}", fs, fsglob.WithIgnoreCase(true))

//     matches, err := g.Expand()
//     if err != nil {
//         // Handle error
//     }
//     fmt.Println("Found matches:", matches)

package fsglob

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/IgorBayerl/fsglob/filesystem"
)

var (
	// globCharacters are special characters used in glob patterns.
	globCharacters = []rune{'*', '?', '[', ']', '{', '}'}

	// regexSpecialChars are characters that have special meaning in regular expressions.
	// Used when converting glob patterns to regex to know which characters to escape.
	regexSpecialChars = map[rune]bool{
		'[': true, '\\': true, '^': true, '$': true, '.': true, '|': true,
		'?': true, '*': true, '+': true, '(': true, ')': true, '{': true, '}': true,
	}

	// regexOrStringCache caches compiled regular expressions or literal strings
	// derived from glob pattern segments.
	// Key: pattern string + "|" + case-sensitivity_flag (e.g., "pat?tern*|true")
	regexOrStringCache = make(map[string]*RegexOrString)
	// cacheMutex protects concurrent access to regexOrStringCache.
	cacheMutex = &sync.Mutex{}
)

// RegexOrString is a helper struct that holds either a compiled regular expression
// or a literal string for efficient pattern matching. It is used internally to
// optimize matching by avoiding regex compilation for simple patterns.
type RegexOrString struct {
	// CompiledRegex is the compiled regular expression if the pattern segment contains wildcards.
	CompiledRegex *regexp.Regexp
	// IsRegex indicates if CompiledRegex is used for matching.
	IsRegex bool
	// LiteralPattern is the original glob pattern segment if it's treated as a literal.
	LiteralPattern string
	// IgnoreCase indicates if matching should be case-insensitive for literal patterns.
	IgnoreCase bool
	// OriginalRegexPattern stores the regex string pattern before compilation (for debugging or special checks).
	OriginalRegexPattern string
}

// IsMatch checks if the input string matches this RegexOrString.
// For regex patterns, it uses the compiled regex.
// For literal patterns, it performs a string comparison, respecting IgnoreCase.
func (ros *RegexOrString) IsMatch(input string) bool {
	if ros.IsRegex {
		return ros.CompiledRegex.MatchString(input)
	}
	if ros.IgnoreCase {
		return strings.EqualFold(ros.LiteralPattern, input)
	}
	return ros.LiteralPattern == input
}

// Glob is the main globbing engine. It holds the original pattern and configuration
// for a globbing operation.
//
// An instance is created with NewGlob, which allows for advanced configuration,
// such as setting a custom filesystem or enabling case-insensitive matching.
type Glob struct {
	OriginalPattern string
	IgnoreCase      bool
	FS              filesystem.Filesystem
	platform        string
}

func (g *Glob) joinPath(elem1, elem2 string) string {
	if g.platform == "windows" {
		e1s := strings.ReplaceAll(elem1, "\\", "/")
		e2s := strings.ReplaceAll(elem2, "\\", "/")
		joined := path.Join(e1s, e2s)
		return strings.ReplaceAll(joined, "/", "\\")
	}
	return path.Join(elem1, elem2)
}

func (g *Glob) parentDir(p string) string {
	if g.platform == "windows" {
		pWithSlashes := strings.ReplaceAll(p, "\\", "/")
		dirWithSlashes := path.Dir(pWithSlashes)
		return strings.ReplaceAll(dirWithSlashes, "/", "\\")
	}
	return path.Dir(p)
}

// absForPlatform turns a possibly-relative path into an absolute one
// without mixing host-OS separators.
func (g *Glob) absForPlatform(p string) (string, error) {
	if g.isAbsolutePath(p) {
		return g.normalizePathForFS(p), nil
	}

	cwd, err := g.FS.Getwd()
	if err != nil {
		return "", err
	}
	if g.platform == "windows" {
		return g.joinPath(cwd, g.normalizePathForFS(p)), nil
	}
	return path.Clean(path.Join(g.normalizePathForPattern(cwd), p)), nil
}

// A GlobOption configures a Glob instance. It's a functional option type
// used with NewGlob.
type GlobOption func(*Glob)

// WithIgnoreCase returns a GlobOption that sets the globber to perform
// case-insensitive matching. On Windows, literal (non-wildcard) path
// segments are case-insensitive by default. This option forces wildcards
// and all other patterns to also ignore case.
func WithIgnoreCase(v bool) GlobOption { return func(g *Glob) { g.IgnoreCase = v } }

// NewGlob creates and returns a new Glob engine for the given pattern.
// It uses the provided filesystem `fs` for all file operations. If `fs` is nil,
// it defaults to the host operating system's filesystem.
//
// Functional options (GlobOption) can be provided to customize the globber's
// behavior, such as enabling case-insensitive matching with WithIgnoreCase.
func NewGlob(pattern string, fs filesystem.Filesystem, opts ...GlobOption) *Glob {
	if fs == nil {
		fs = filesystem.DefaultFS{}
	}

	platform := "unix"
	if p, ok := fs.(filesystem.Platformer); ok {
		platform = p.Platform()
	} else if _, ok := fs.(filesystem.DefaultFS); ok {
		if filepath.Separator == '\\' {
			platform = "windows"
		}
	}

	g := &Glob{
		OriginalPattern: pattern,
		IgnoreCase:      false,
		FS:              fs,
		platform:        platform,
	}

	for _, opt := range opts {
		opt(g)
	}

	return g
}

// String returns the original, unmodified pattern that was used to create the Glob instance.
func (g *Glob) String() string { return g.OriginalPattern }

// Expand executes the glob pattern and returns all matching file paths.
// It is the primary method for retrieving results from a configured Glob instance.
// The function returns an empty, non-nil slice if no matches are found.
func (g *Glob) Expand() ([]string, error) {
	res, err := g.expandInternal(g.OriginalPattern, false)
	if err != nil && g.IgnoreCase {
		slog.Warn("Ignoring malformed glob pattern",
			"pattern", g.OriginalPattern, "error", err)
		return []string{}, nil
	}
	return res, err
}

// ExpandNames executes the glob pattern and returns all matching file paths.
//
// Deprecated: This function is an alias for Expand. Please use Expand instead.
func (g *Glob) ExpandNames() ([]string, error) {
	res, err := g.expandInternal(g.OriginalPattern, false)
	if err != nil && g.IgnoreCase { // tolerant mode
		slog.Warn("Ignoring malformed glob pattern",
			"pattern", g.OriginalPattern, "error", err)
		return []string{}, nil
	}
	return res, err
}

func (g *Glob) tryWindowsCaseFold(absPath string, dirOnly bool) ([]string, error) {
	// Normalize to forward slashes for processing with `path` package
	normalizedPath := g.normalizePathForPattern(absPath)
	cleanPath := path.Clean(normalizedPath)

	// Handle root case `C:\` which clean might change to `C:`
	if len(cleanPath) == 2 && cleanPath[1] == ':' {
		cleanPath += "/"
	}

	vol := ""
	rest := cleanPath
	if len(cleanPath) > 1 && cleanPath[1] == ':' {
		vol = cleanPath[:2]
		rest = cleanPath[2:]
	}

	// path.Clean might remove the leading slash. Put it back.
	if vol != "" && !strings.HasPrefix(rest, "/") {
		rest = "/" + rest
	}

	// Split path components
	parts := strings.Split(strings.Trim(rest, "/"), "/")

	cur := vol + `\`
	if vol == "" {
		// Handle UNC paths if necessary, for now assume drive letters
		if strings.HasPrefix(cleanPath, "//") {
			// very basic UNC handling
			parts = strings.Split(strings.TrimPrefix(cleanPath, "//"), "/")
			if len(parts) < 2 {
				return nil, nil
			}
			cur = `\\` + parts[0] + `\` + parts[1]
			parts = parts[2:]
		} else {
			cur = `\`
		}
	}

	// Walk every segment and pick the real-cased name.
	for _, p := range parts {
		if p == "" {
			continue
		}
		entries, err := g.FS.ReadDir(cur)
		if err != nil {
			return nil, nil // unreadable -> give up
		}
		var matched string
		for _, de := range entries {
			if strings.EqualFold(de.Name(), p) {
				matched = de.Name()
				break
			}
		}
		if matched == "" {
			return nil, nil // path element not found
		}
		cur = g.joinPath(cur, matched)
	}

	info, err := g.FS.Stat(cur)
	if err != nil || (dirOnly && !info.IsDir()) {
		return nil, nil
	}
	return []string{cur}, nil
}

// createRegexOrString compiles a glob pattern segment into a RegexOrString instance.
// It uses a cache to store and retrieve compiled regexes or literal patterns
// to avoid redundant compilations.
func (g *Glob) createRegexOrString(patternSegment string) (*RegexOrString, error) {
	hasWildcards := strings.ContainsAny(patternSegment, "*?[]")

	effectiveIC := g.IgnoreCase
	if !hasWildcards && g.platform == "windows" {
		effectiveIC = true
	}

	cacheKey := patternSegment + "|" + fmt.Sprintf("%t", effectiveIC)

	cacheMutex.Lock()
	if cached, ok := regexOrStringCache[cacheKey]; ok {
		cacheMutex.Unlock()
		return cached, nil
	}
	cacheMutex.Unlock()

	if !hasWildcards {
		ros := &RegexOrString{
			IsRegex:        false,
			LiteralPattern: patternSegment,
			IgnoreCase:     effectiveIC,
		}
		cacheMutex.Lock()
		regexOrStringCache[cacheKey] = ros
		cacheMutex.Unlock()
		return ros, nil
	}

	regexPatternStr, err := globToRegexPattern(patternSegment, g.IgnoreCase)
	if err != nil {
		return nil, fmt.Errorf("failed to convert glob segment %q: %w", patternSegment, err)
	}

	re, err := regexp.Compile(regexPatternStr)
	if err != nil {
		return nil, fmt.Errorf("failed to compile regex %q: %w", regexPatternStr, err)
	}

	ros := &RegexOrString{
		CompiledRegex:        re,
		IsRegex:              true,
		LiteralPattern:       patternSegment,
		IgnoreCase:           g.IgnoreCase,
		OriginalRegexPattern: regexPatternStr,
	}
	cacheMutex.Lock()
	regexOrStringCache[cacheKey] = ros
	cacheMutex.Unlock()
	return ros, nil
}

func (g *Glob) isAbsolutePath(path string) bool {
	if g.platform == "windows" {
		return (len(path) >= 3 && path[1] == ':' && (path[2] == '\\' || path[2] == '/')) ||
			strings.HasPrefix(path, "\\\\") ||
			strings.HasPrefix(path, "/")
	}
	return strings.HasPrefix(path, "/")
}

func (g *Glob) normalizePathForFS(p string) string {
	if g.platform == "windows" {
		return strings.ReplaceAll(p, "/", "\\")
	}
	return strings.ReplaceAll(p, "\\", "/")
}

// normalizePathForPattern converts a path to forward slashes for pattern matching
func (g *Glob) normalizePathForPattern(p string) string {
	return strings.ReplaceAll(p, "\\", "/")
}

// expandInternal is the core recursive matching function.
func (g *Glob) expandInternal(pattern string, dirOnly bool) ([]string, error) {
	if pattern == "" {
		return []string{}, nil
	}

	normalizedPattern := g.normalizePathForPattern(pattern)

	if !strings.ContainsAny(normalizedPattern, string(globCharacters)) {
		absPath, err := g.absForPlatform(pattern)
		if err != nil {
			return nil, err
		}

		if g.platform == "windows" {
			paths, _ := g.tryWindowsCaseFold(absPath, dirOnly)
			if paths == nil {
				return []string{}, nil
			}
			return paths, nil
		}

		info, err := g.FS.Stat(absPath)
		if err == nil && (!dirOnly || info.IsDir()) {
			return []string{absPath}, nil
		}
		return []string{}, nil
	}

	parent := path.Dir(normalizedPattern)
	child := path.Base(normalizedPattern)

	if parent == "." && !g.isAbsolutePath(normalizedPattern) {
		cwd, err := g.FS.Getwd()
		if err != nil {
			return []string{}, fmt.Errorf("failed to get working directory: %w", err)
		}
		parent = g.normalizePathForPattern(cwd)
	}

	if strings.Count(child, "}") > strings.Count(child, "{") {
		return g.handleCrossSeparatorBrace(normalizedPattern, dirOnly)
	}

	if child == "**" {
		parentDirs, err := g.expandInternal(g.normalizePathForFS(parent), true)
		if err != nil {
			return nil, err
		}
		var allResults []string
		seenPaths := make(map[string]bool)
		for _, pDir := range parentDirs {
			descendants, err := g.getRecursiveDirectoriesAndFiles(pDir, dirOnly)
			if err != nil {
				continue
			}
			if !seenPaths[pDir] {
				info, err := g.FS.Stat(pDir)
				if err == nil {
					if !dirOnly || info.IsDir() {
						allResults = append(allResults, pDir)
						seenPaths[pDir] = true
					}
				}
			}
			for _, d := range descendants {
				if !seenPaths[d] {
					allResults = append(allResults, d)
					seenPaths[d] = true
				}
			}
		}
		if allResults == nil {
			return []string{}, nil
		}
		return allResults, nil
	}

	return g.processPathSegment(parent, child, dirOnly)
}

// handleCrossSeparatorBrace handles patterns like "{a/b,c}/d.txt"
func (g *Glob) handleCrossSeparatorBrace(normalizedPattern string, dirOnly bool) ([]string, error) {
	groups, err := ungroup(normalizedPattern)
	if err != nil {
		return nil, fmt.Errorf("error ungrouping path '%s': %w", normalizedPattern, err)
	}
	var allResults []string
	seenPaths := make(map[string]bool)
	for _, groupPattern := range groups {
		expanded, err := g.expandInternal(g.normalizePathForFS(groupPattern), dirOnly)
		if err != nil {
			slog.Warn("Error expanding group pattern", "pattern", groupPattern, "error", err)
			continue
		}
		for _, p := range expanded {
			if !seenPaths[p] {
				allResults = append(allResults, p)
				seenPaths[p] = true
			}
		}
	}
	if allResults == nil {
		return []string{}, nil
	}
	return allResults, nil

}

// processPathSegment is the main workhorse for a standard `parent/child` pattern.
func (g *Glob) processPathSegment(parentPattern, childPattern string, dirOnly bool) ([]string, error) {
	parentForFS := g.normalizePathForFS(parentPattern)
	expandedParentDirs, err := g.expandInternal(parentForFS, true)
	if err != nil {
		return nil, err
	}

	ungroupedChildSegments, err := ungroup(childPattern)
	if err != nil {
		return nil, fmt.Errorf("error ungrouping child segment '%s': %w", childPattern, err)
	}

	var childRegexes []*RegexOrString
	for _, segment := range ungroupedChildSegments {
		ros, err := g.createRegexOrString(segment)
		if err != nil {
			return nil, err
		}
		childRegexes = append(childRegexes, ros)
	}

	var allMatches []string
	seenPaths := make(map[string]bool)

	for _, parentDir := range expandedParentDirs {
		entries, readDirErr := g.FS.ReadDir(parentDir)
		if readDirErr != nil {
			if os.IsNotExist(readDirErr) {
				continue
			}
			slog.Warn("Error reading directory", "directory", parentDir, "error", readDirErr)
			continue
		}

		for _, entry := range entries {
			isDir := entry.IsDir()
			if !dirOnly || isDir {
				for _, ros := range childRegexes {
					if ros.IsMatch(entry.Name()) {
						absEntryPath := g.joinPath(parentDir, entry.Name())
						if !seenPaths[absEntryPath] {
							allMatches = append(allMatches, absEntryPath)
							seenPaths[absEntryPath] = true
						}
						break
					}
				}
			}
		}

		// Handle '.' and '..' matching
		for _, ros := range childRegexes {
			if ros.LiteralPattern == "." {
				if !seenPaths[parentDir] {
					allMatches = append(allMatches, parentDir)
					seenPaths[parentDir] = true
				}
			} else if ros.LiteralPattern == ".." {
				grandParentDir := g.parentDir(parentDir)
				if grandParentDir != parentDir {
					if !seenPaths[grandParentDir] {
						allMatches = append(allMatches, grandParentDir)
						seenPaths[grandParentDir] = true
					}
				}
			}
		}
	}

	if allMatches == nil {
		return []string{}, nil
	}
	return allMatches, nil
}

// globToRegexPattern converts a glob pattern segment to a Go regular expression string.
func globToRegexPattern(globSegment string, ignoreCase bool) (string, error) {
	var regex strings.Builder
	if ignoreCase {
		regex.WriteString("(?i)")
	}
	regex.WriteRune('^')

	if globSegment == "**" {
		regex.WriteString(".*")
	} else {
		globSegment = strings.ReplaceAll(globSegment, "**", ".*")
		inCharClass := false
		for _, r := range globSegment {
			if inCharClass {
				if r == ']' {
					inCharClass = false
				}
				regex.WriteRune(r)
				continue
			}
			switch r {
			case '*':
				regex.WriteString("[^/\\\\]*")
			case '?':
				regex.WriteRune('.')
			case '[':
				inCharClass = true
				regex.WriteRune(r)
			default:
				if _, isSpecial := regexSpecialChars[r]; isSpecial {
					regex.WriteRune('\\')
				}
				regex.WriteRune(r)
			}
		}
		if inCharClass {
			return "", fmt.Errorf("unterminated character class: %s", globSegment)
		}
	}

	regex.WriteRune('$')
	return regex.String(), nil
}

// ungroup handles brace expansion, e.g., "{a,b}c" -> ["ac", "bc"].
func ungroup(path string) ([]string, error) {
	if !strings.Contains(path, "{") {
		return []string{path}, nil
	}

	var results []string
	level := 0
	firstOpenBrace := -1

	for i, char := range path {
		switch char {
		case '{':
			if level == 0 {
				firstOpenBrace = i
			}
			level++
		case '}':
			level--
			if level == 0 && firstOpenBrace != -1 {
				prefix := path[:firstOpenBrace]
				groupContent := path[firstOpenBrace+1 : i]
				suffix := path[i+1:]

				var groupParts []string
				partBuilder := strings.Builder{}
				subLevel := 0
				for _, gc := range groupContent {
					if gc == '{' {
						subLevel++
						partBuilder.WriteRune(gc)
					} else if gc == '}' {
						subLevel--
						partBuilder.WriteRune(gc)
					} else if gc == ',' && subLevel == 0 {
						groupParts = append(groupParts, partBuilder.String())
						partBuilder.Reset()
					} else {
						partBuilder.WriteRune(gc)
					}
				}
				groupParts = append(groupParts, partBuilder.String())

				expandedSuffixes, err := ungroup(suffix)
				if err != nil {
					return nil, err
				}

				for _, gp := range groupParts {
					currentCombinedPrefixPart := prefix + gp
					expandedPrefixParts, err := ungroup(currentCombinedPrefixPart)
					if err != nil {
						return nil, err
					}

					for _, epp := range expandedPrefixParts {
						for _, es := range expandedSuffixes {
							results = append(results, epp+es)
						}
					}
				}
				return results, nil
			}
		}
	}

	if level != 0 {
		return nil, fmt.Errorf("unbalanced braces in pattern: %s", path)
	}
	return []string{path}, nil
}

// getRecursiveDirectoriesAndFiles is a helper for `**`.
func (g *Glob) getRecursiveDirectoriesAndFiles(root string, dirOnly bool) ([]string, error) {
	var paths []string
	queue := []string{root}
	visited := make(map[string]struct{})

	for len(queue) > 0 {
		currentPath := queue[0]
		queue = queue[1:]

		if _, ok := visited[currentPath]; ok {
			continue
		}
		visited[currentPath] = struct{}{}

		entries, err := g.FS.ReadDir(currentPath)
		if err != nil {
			slog.Warn("Error reading directory", "path", currentPath, "error", err)
			continue
		}

		for _, entry := range entries {
			nextPath := g.joinPath(currentPath, entry.Name())
			entryInfo, err := entry.Info()
			if err != nil {
				slog.Warn("Could not get info for entry", "path", nextPath, "error", err)
				continue
			}
			if !dirOnly || entryInfo.IsDir() {
				paths = append(paths, nextPath)
			}
			if entryInfo.IsDir() {
				queue = append(queue, nextPath)
			}
		}
	}
	return paths, nil
}

// GetFiles is the public entry point for globbing.
func GetFiles(pattern string) ([]string, error) {
	if pattern == "" {
		return []string{}, nil
	}
	g := NewGlob(pattern, nil)
	return g.ExpandNames()
}
