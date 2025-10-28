package utils

// FindMatchingBrace scans source code lines to find the closing brace '}'
// that matches the first opening brace '{' found on or after the startLineIndex.
// It correctly handles nested braces, line comments, block comments, and string literals.
func FindMatchingBrace(sourceLines []string, startLineIndex int) (int, bool) {
	braceLevel := 0
	inBlockComment := false
	inString := false
	foundFirstBrace := false

	for i := startLineIndex; i < len(sourceLines); i++ {
		line := sourceLines[i]
		for j := 0; j < len(line); j++ {
			char := line[j]

			// Handle state transitions for comments and strings
			if inBlockComment {
				if char == '*' && j+1 < len(line) && line[j+1] == '/' {
					inBlockComment = false
					j++ // Skip the '/'
				}
				continue
			}
			if inString {
				if char == '\\' { // Handle escaped quotes
					j++
				} else if char == '"' {
					inString = false
				}
				continue
			}

			// Check for start of comments or strings
			if char == '/' && j+1 < len(line) {
				if line[j+1] == '/' { // Line comment
					goto nextLine
				}
				if line[j+1] == '*' { // Block comment
					inBlockComment = true
					j++
					continue
				}
			}
			if char == '"' {
				inString = true
				continue
			}

			// Process braces only if not inside a comment or string
			if char == '{' {
				braceLevel++
				foundFirstBrace = true
			} else if char == '}' {
				braceLevel--
			}

			if foundFirstBrace && braceLevel == 0 {
				return i + 1, true // Return 1-based line number
			}
		}
	nextLine:
	}

	return -1, false // No matching brace found
}
