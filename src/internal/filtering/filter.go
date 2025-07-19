// Package filtering provides a robust mechanism for including or excluding elements
// (like assemblies, classes, or files) from a report based on user-defined patterns.
//
// The filtering logic is inspired by common .gitignore and build tool patterns,
// providing a familiar and powerful way for users to control the report's scope.
//
// # Filter Syntax
//
// Each filter is a string that must start with either a '+' (for inclusion)
// or a '-' (for exclusion). The rest of the string is a pattern that supports
// wildcards:
//
//   - `+` (Inclusion): If an element matches an inclusion filter, it is considered
//     for the report, unless it also matches an exclusion filter. If no inclusion
//     filters are provided, all elements are considered included by default.
//
//   - `-` (Exclusion): If an element matches an exclusion filter, it is always
//     removed from the report, regardless of any matching inclusion filters.
//     Exclusion takes precedence over inclusion.
//
//   - `*` (Wildcard): Matches zero or more characters. For example, `+MyProject.*`
//     will include all elements starting with "MyProject.".
//
//   - `?` (Wildcard): Matches exactly one character.
//
// All pattern matching is case-insensitive by default.
package filtering

import (
	"fmt"
	"regexp"
	"strings"
)

type IFilter interface {
	// IsElementIncludedInReport determines if an element with the given name
	// should be included in the report based on the configured filter rules.
	// The core logic is: an element is included if it matches any include filter
	// AND does not match any exclude filter.
	IsElementIncludedInReport(name string) bool

	// HasCustomFilters returns true if the filter was created with specific
	// user-defined rules (i.e., any '+' or '-' filters). It returns false if
	// the filter is using the default "include all" behavior.
	HasCustomFilters() bool
}

type DefaultFilter struct {
	includeFilters []*regexp.Regexp
	excludeFilters []*regexp.Regexp
	hasCustom      bool
}

// The optional `osIndependantPathSeparator` parameter, if true, treats both `/` and `\`
// as path separators in the patterns, making file filters work seamlessly across
// different operating systems.
func NewDefaultFilter(filters []string, osIndependantPathSeparator ...bool) (IFilter, error) {
	osPathSep := false
	if len(osIndependantPathSeparator) > 0 {
		osPathSep = osIndependantPathSeparator[0]
	}

	df := &DefaultFilter{}
	var errs []string

	for _, f := range filters {
		trimmedFilter := strings.TrimSpace(f)
		if trimmedFilter == "" {
			continue // Ignore empty strings
		}

		if strings.HasPrefix(trimmedFilter, "+") {
			re, err := createFilterRegex(trimmedFilter, osPathSep)
			if err != nil {
				errs = append(errs, fmt.Sprintf("invalid include filter '%s': %v", trimmedFilter, err))
				continue
			}
			df.includeFilters = append(df.includeFilters, re)
		} else if strings.HasPrefix(trimmedFilter, "-") {
			re, err := createFilterRegex(trimmedFilter, osPathSep)
			if err != nil {
				errs = append(errs, fmt.Sprintf("invalid exclude filter '%s': %v", trimmedFilter, err))
				continue
			}
			df.excludeFilters = append(df.excludeFilters, re)
		} else {
			errs = append(errs, fmt.Sprintf("filter '%s' must start with '+' or '-'", trimmedFilter))
		}
	}

	if len(errs) > 0 {
		return nil, fmt.Errorf("error creating filters: %s", strings.Join(errs, "; "))
	}

	df.hasCustom = len(df.includeFilters) > 0 || len(df.excludeFilters) > 0

	// If no include filters are specified, default to including everything.
	// This is the most intuitive behavior for users.
	if len(df.includeFilters) == 0 {
		re, _ := createFilterRegex("+*", false) // Default include all pattern
		df.includeFilters = append(df.includeFilters, re)
	}

	return df, nil
}

// An element is included if it matches at least one include filter and no exclude filters.
func (df *DefaultFilter) IsElementIncludedInReport(name string) bool {
	// Exclusion filters always take precedence.
	for _, excludeRe := range df.excludeFilters {
		if excludeRe.MatchString(name) {
			return false
		}
	}

	// If not excluded, check if it matches any inclusion filter.
	for _, includeRe := range df.includeFilters {
		if includeRe.MatchString(name) {
			return true
		}
	}

	// If it reached here, it means it was not excluded but did not match any include filter.
	// This can only happen if there are specific include filters and the element name didn't match any of them.
	return false
}

func (df *DefaultFilter) HasCustomFilters() bool {
	return df.hasCustom
}

// createFilterRegex converts a filter string (e.g., "+MyNamespace.*") to a regular expression.
// It handles escaping and wildcard conversion.
func createFilterRegex(filter string, osIndependantPathSeparator bool) (*regexp.Regexp, error) {
	if len(filter) == 0 {
		return nil, fmt.Errorf("empty filter string")
	}
	// Remove '+' or '-' prefix to get the pattern.
	pattern := filter[1:]

	// Validate for balanced brackets before quoting, as QuoteMeta would escape them.
	if strings.Count(pattern, "[") != strings.Count(pattern, "]") {
		return nil, fmt.Errorf("unbalanced brackets in filter pattern '%s'", pattern)
	}

	// First, escape any characters that have special meaning in regex.
	pattern = regexp.QuoteMeta(pattern)

	// Then, convert our simple wildcards to their regex equivalents.
	// `\*` becomes `.*` (match any characters).
	// `\?` becomes `.` (match any single character).
	pattern = strings.ReplaceAll(pattern, `\*`, ".*")
	pattern = strings.ReplaceAll(pattern, `\?`, ".")

	// If specified, make path separators in the pattern platform-agnostic.
	if osIndependantPathSeparator {
		// This replaces both `/` and `\` with a character class `[/\\]`
		// that matches either separator.
		pattern = strings.ReplaceAll(pattern, "/", `[/\\]`)
		pattern = strings.ReplaceAll(pattern, `\\\\`, `[/\\]`)
	}

	// Compile the final regex. It is case-insensitive `(?i)` and anchored
	// to match the entire string `^...$`.
	return regexp.Compile("(?i)^" + pattern + "$")
}
