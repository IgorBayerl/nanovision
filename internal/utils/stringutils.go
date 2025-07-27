package utils

import (
	"fmt"
	"regexp"
	"strings"
)

// SplitThatEnsuresGlobsAreSafe splits a string by any of the given separators,
// but does not split within brace-delimited glob patterns like {group1,group2}.
// Based on: Palmmedia.ReportGenerator.Core.Common.StringExtensions.cs (SplitThatEnsuresGlobsAreSafe method)
func SplitThatEnsuresGlobsAreSafe(s string, separators []rune) []string {
	if len(separators) == 0 {
		return []string{s}
	}

	var parts []string
	var currentPart strings.Builder
	braceLevel := 0

	isSeparator := func(r rune) bool {
		for _, sep := range separators {
			if r == sep {
				return true
			}
		}
		return false
	}

	for _, char := range s {
		if char == '{' {
			braceLevel++
			currentPart.WriteRune(char)
		} else if char == '}' {
			if braceLevel > 0 { // Only decrement if we are inside braces
				braceLevel--
			}
			currentPart.WriteRune(char)
		} else if isSeparator(char) && braceLevel == 0 {
			if currentPart.Len() > 0 { // Add part if not empty
				parts = append(parts, strings.TrimSpace(currentPart.String()))
				currentPart.Reset()
			}
		} else {
			currentPart.WriteRune(char)
		}
	}

	if currentPart.Len() > 0 { // Add the last part
		parts = append(parts, strings.TrimSpace(currentPart.String()))
	}

	// Handle case where input is empty or only separators
	if len(s) > 0 && len(parts) == 0 && currentPart.Len() == 0 {
		return []string{""} // if s was e.g. ";", parts would be empty. C# returns {""}.
	} else if len(parts) == 0 && currentPart.Len() == 0 && len(s) == 0 { // s is empty
		return []string{""} // C# returns {""} for empty input
	}

	return parts
}

// FilterToRegex converts a simple filter pattern (+/- prefix, * wildcard) to a regex.
// E.g., "+MyAssembly.*" becomes "^MyAssembly\..*$" (case-insensitive).
// Returns the regex and a boolean indicating if it's an inclusion (true) or exclusion (false) filter.
// Based on: Palmmedia.ReportGenerator.Core.Parser.Filtering.DefaultFilter.cs (CreateFilterRegex method)
// Original C# logic involves Regex.Escape and specific replacements for '*'
// This Go version uses regexp.QuoteMeta and similar replacements.
func FilterToRegex(filterPattern string) (*regexp.Regexp, bool, error) {
	if len(filterPattern) < 2 || (filterPattern[0] != '+' && filterPattern[0] != '-') {
		return nil, false, fmt.Errorf("invalid filter pattern: '%s'. Must start with '+' or '-'", filterPattern)
	}

	isInclude := filterPattern[0] == '+'
	pattern := filterPattern[1:]

	// Escape regex special characters first
	pattern = regexp.QuoteMeta(pattern)

	// Then convert glob-like wildcards '*' and '?'
	// C# original: filter = filter.Replace("*", "$$$*"); ... filter = Regex.Escape(filter); filter = filter.Replace(@"\$\$\$\*", ".*");
	// Go: QuoteMeta escapes '*', so we replace `\*` with `.*`
	pattern = strings.ReplaceAll(pattern, `\*`, ".*")
	pattern = strings.ReplaceAll(pattern, `\?`, ".") // QuoteMeta escapes '?', so replace `\?` with `.`

	// Anchor the pattern and make it case-insensitive
	regexString := "(?i)^" + pattern + "$"
	re, err := regexp.Compile(regexString)
	if err != nil {
		return nil, false, fmt.Errorf("failed to compile filter regex for '%s': %w", filterPattern, err)
	}
	return re, isInclude, nil
}

var (
	// Based on: Palmmedia.ReportGenerator.Core.Reporting.Builders.Rendering.StringHelper.cs (ReplaceInvalidPathChars method)
	// Original C# Regex: Regex.Replace(path, "[^\\w^\\.]", "_")
	// Go version allows hyphens explicitly and uses `+` to match one or more.
	invalidPathCharsRegex = regexp.MustCompile(`[^\w.-]+`)
)

// ReplaceInvalidPathChars replaces characters in a path that are not word characters, dots, or hyphens with an underscore.
func ReplaceInvalidPathChars(path string) string {
	return invalidPathCharsRegex.ReplaceAllString(path, "_")
}

// GetShortMethodName creates a shorter, display-friendly version of a full method name.
// It replaces complex signatures with "()" or "(...)".
// E.g., "MyMethod(System.String, System.Int32)" becomes "MyMethod(...)".
// E.g., "MyMethod()" remains "MyMethod()".
// E.g., "MyMethod" becomes "MyMethod" (if no parentheses were present).
// Based on logic in: Palmmedia.ReportGenerator.Core.Parser.CoberturaParser (GetShortMethodName method, though it was private there)
// and similar logic in other parts of the C# codebase for display names.
func GetShortMethodName(fullName string) string {
	indexOpen := strings.Index(fullName, "(")

	if indexOpen <= 0 { // No opening parenthesis or it's the first character (unlikely for valid method names)
		return fullName
	}

	// Find the matching closing parenthesis. This is a simplification and assumes no nested parentheses in the signature part itself.
	// For more complex scenarios (e.g. generic types with angle brackets in signature), this might need refinement.
	indexClose := strings.Index(fullName[indexOpen:], ")")
	if indexClose == -1 { // No closing parenthesis found after open
		return fullName // Or perhaps append "()" if it's clearly a method name missing them
	}
	indexClose += indexOpen // Adjust indexClose to be relative to the start of fullName

	var signature string
	if indexClose > indexOpen+1 { // Signature is not just "()"
		signature = "(...)"
	} else { // Signature is "()"
		signature = "()"
	}

	return fullName[:indexOpen] + signature
}
