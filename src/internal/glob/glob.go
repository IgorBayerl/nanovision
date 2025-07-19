// Package glob provides functionality for finding files and directories
// by matching their path names against a pattern.
// Aim to support:
//   - `?`: Matches any single character in a file or directory name.
//   - `*`: Matches zero or more characters in a file or directory name.
//   - `**`: Matches zero or more recursive directories.
//   - `[...]`: Matches a set of characters in a name (e.g., `[abc]`, `[a-z]`).
//   - `{group1,group2,...}`: Matches any of the pattern groups.
//
// Case-insensitivity is the default behavior for matching.
package glob

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/filesystem"
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

type Glob struct {
	OriginalPattern string
	IgnoreCase      bool
	FS              filesystem.Filesystem
	platform        string
}

func (g *Glob) joinPath(elem1, elem2 string) string {
	if g.platform == "windows" {
		return filepath.Join(elem1, elem2)
	}
	return path.Join(elem1, elem2) // always “/”
}

func (g *Glob) parentDir(p string) string {
	if g.platform == "windows" {
		return filepath.Dir(p)
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
		return filepath.Join(cwd, g.normalizePathForFS(p)), nil
	}
	// unix use the slash variant only
	return path.Clean(path.Join(g.normalizePathForPattern(cwd), p)), nil
}

type GlobOption func(*Glob)

func WithIgnoreCase(v bool) GlobOption { return func(g *Glob) { g.IgnoreCase = v } }

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

func (g *Glob) String() string { return g.OriginalPattern }

func (g *Glob) ExpandNames() ([]string, error) {
	res, err := g.expandInternal(g.OriginalPattern, false)
	if err != nil && g.IgnoreCase { // tolerant mode
		slog.Warn("Ignoring malformed glob pattern",
			"pattern", g.OriginalPattern, "error", err)
		return []string{}, nil
	}
	return res, err
}

func (g *Glob) Expand() ([]string, error) {
	res, err := g.expandInternal(g.OriginalPattern, false)
	if err != nil && g.IgnoreCase {
		slog.Warn("Ignoring malformed glob pattern",
			"pattern", g.OriginalPattern, "error", err)
		return []string{}, nil
	}
	return res, err
}

func (g *Glob) tryWindowsCaseFold(absPath string, dirOnly bool) ([]string, error) {
	clean := filepath.Clean(absPath)
	vol := ""
	rest := clean
	if len(clean) >= 2 && clean[1] == ':' {
		vol = clean[:2]  // "C:"
		rest = clean[2:] // everything after the drive
	}
	rest = strings.TrimPrefix(rest, `\`)
	parts := strings.Split(rest, `\`)

	// Start at the root directory (e.g. "C:\").
	cur := vol + `\`

	// Walk every segment and pick the real-cased name.
	for _, p := range parts {
		if p == "" {
			continue // guard against "C:\"
		}
		entries, err := g.FS.ReadDir(cur) // list the directory we are in
		if err != nil {
			return nil, nil // unreadable -> give up
		}
		var matched string
		for _, de := range entries {
			if strings.EqualFold(de.Name(), p) { // case-insensitive compare
				matched = de.Name()
				break
			}
		}
		if matched == "" {
			return nil, nil // path element not found
		}
		cur = g.joinPath(cur, matched) // advance with correct case
	}

	// We have rebuilt the canonical cased path in cur.
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
		// Windows absolute paths: C:\... or \\... (UNC) or /... (converted from Unix-style)
		return (len(path) >= 3 && path[1] == ':' && (path[2] == '\\' || path[2] == '/')) ||
			strings.HasPrefix(path, "\\\\") ||
			strings.HasPrefix(path, "/")
	}
	// Unix absolute paths start with /
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
// It returns a slice of absolute paths for matched files or directories.
// `pattern` is the current glob pattern being processed.
// `dirOnly` specifies if only directories should be matched.
func (g *Glob) expandInternal(pattern string, dirOnly bool) ([]string, error) {
	if pattern == "" {
		return []string{}, nil
	}

	// Normalize pattern for internal processing (always use forward slashes)
	normalizedPattern := g.normalizePathForPattern(pattern)

	// Handle literal paths (no glob characters)
	if !strings.ContainsAny(normalizedPattern, string(globCharacters)) {
		absPath, err := g.absForPlatform(pattern)
		if err != nil {
			return nil, err
		}

		info, err := g.FS.Stat(absPath)
		if err == nil && (!dirOnly || info.IsDir()) {
			return []string{absPath}, nil
		}

		// Windows: retry with case-insensitive scan of the parent dir
		if g.platform == "windows" {
			if paths, _ := g.tryWindowsCaseFold(absPath, dirOnly); len(paths) > 0 {
				return paths, nil
			}
		}
		return []string{}, nil
	}

	// Split path into parent and child components
	parent := path.Dir(normalizedPattern)
	child := path.Base(normalizedPattern)

	// Handle root directory case
	if parent == "." && !g.isAbsolutePath(normalizedPattern) {
		cwd, err := g.FS.Getwd()
		if err != nil {
			return []string{}, fmt.Errorf("failed to get working directory: %w", err)
		}
		parent = g.normalizePathForPattern(cwd)
	}

	// Handle cross-separator brace expansion
	if strings.Count(child, "}") > strings.Count(child, "{") {
		return g.handleCrossSeparatorBrace(normalizedPattern, dirOnly)
	}

	// Handle recursive wildcard `**` as the final segment
	if child == "**" {
		parentDirs, err := g.expandInternal(g.normalizePathForFS(parent), true)
		if err != nil {
			return nil, err
		}
		var allResults []string
		seenPaths := make(map[string]bool)
		for _, pDir := range parentDirs {
			// Get ALL descendants, including the parent dir itself
			descendants, err := g.getRecursiveDirectoriesAndFiles(pDir, dirOnly)
			if err != nil {
				continue
			}
			// Include parent directory itself
			if !dirOnly {
				_, err := g.FS.Stat(pDir)
				if err == nil && !seenPaths[pDir] {
					allResults = append(allResults, pDir)
					seenPaths[pDir] = true
				}
			} else if !seenPaths[pDir] {
				allResults = append(allResults, pDir)
				seenPaths[pDir] = true
			}
			// Add descendants
			for _, d := range descendants {
				if !seenPaths[d] {
					allResults = append(allResults, d)
					seenPaths[d] = true
				}
			}
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
	return allResults, nil
}

// processPathSegment is the main workhorse for a standard `parent/child` pattern.
func (g *Glob) processPathSegment(parentPattern, childPattern string, dirOnly bool) ([]string, error) {
	// Convert parent pattern back to filesystem format for expansion
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
			if !dirOnly || entry.IsDir() {
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
			switch ros.OriginalRegexPattern {
			case `^\.$`: // Matches "."
				if !seenPaths[parentDir] {
					allMatches = append(allMatches, parentDir)
					seenPaths[parentDir] = true
				}
			case `^\.\.$`: // Matches ".."
				grandParentDir := g.parentDir(parentDir)
				if grandParentDir != parentDir { // Avoids root's parent being itself
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

	// Handle `**` as a complete segment
	if globSegment == "**" {
		regex.WriteString(".*") // `**` as a segment should match anything (including nothing)
	} else {
		// Replace `**` with a regex that matches any character (including path separators)
		// This is for patterns like `a**b.txt`
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
				regex.WriteString("[^/\\\\]*") // * matches anything except path separators (both / and \)
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
// It supports nested braces and multiple groups.
func ungroup(path string) ([]string, error) {
	if !strings.Contains(path, "{") {
		return []string{path}, nil
	}

	// This is a common algorithm for brace expansion:
	// Find the first top-level {...} group.
	// Expand this group.
	// For each expansion, prepend the prefix and recursively call ungroup on the (expansion + suffix).

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
			if level == 0 && firstOpenBrace != -1 { // Found a top-level group
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
				groupParts = append(groupParts, partBuilder.String()) // Add the last part

				// Recursively expand suffix first, as it applies to all parts of the current group.
				expandedSuffixes, err := ungroup(suffix)
				if err != nil {
					return nil, err
				}

				for _, gp := range groupParts {
					// Each part of the current group might itself contain groups or be literal.
					// So, we form `prefix + gp` and then expand that, then combine with expandedSuffixes.
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
				return results, nil // Processed the first top-level group
			}
		}
	}

	if level != 0 { // Unbalanced braces
		return nil, fmt.Errorf("unbalanced braces in pattern: %s", path)
	}

	// No top-level group found (e.g. "abc" or "a{b}c" where {b} was handled by inner recursion)
	return []string{path}, nil
}

// getRecursiveDirectoriesAndFiles is a helper for `**` when it's the last segment.
// It lists all files/directories under the root directory recursively.
// If dirOnly is true, only directories are returned.
// Returns absolute paths.
func (g *Glob) getRecursiveDirectoriesAndFiles(root string, dirOnly bool) ([]string, error) {
	var paths []string

	rootInfo, err := g.FS.Stat(root)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	if !rootInfo.IsDir() {
		return []string{}, nil
	}

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

			entryInfo, err := g.FS.Stat(nextPath)
			if err != nil {
				slog.Warn("Could not stat entry", "path", nextPath, "error", err)
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
// It takes a glob pattern and returns a slice of absolute paths to matching files and directories.
// Errors encountered during parts of the expansion (e.g., unreadable directory) are logged as warnings,
// and the function attempts to return successfully found matches.
// A fundamental error (e.g., invalid pattern syntax) will be returned as an error.
func GetFiles(pattern string) ([]string, error) {
	if pattern == "" {
		return []string{}, nil
	}

	g := NewGlob(pattern, nil)
	// Call ExpandNames which uses expandInternal.
	// expandInternal is designed to return errors for fundamental issues.
	return g.ExpandNames()
}
