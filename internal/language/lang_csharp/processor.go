package lang_csharp

import (
	"regexp"
	"strings"

	"github.com/IgorBayerl/AdlerCov/internal/language"
	"github.com/IgorBayerl/AdlerCov/internal/model"
	"github.com/IgorBayerl/AdlerCov/internal/utils"
)

var (
	// Matches methods with a return type. Catches name in group 1.
	// Now handles generic method names like `MyMethod<T>`.
	csharpMethodRegex = regexp.MustCompile(
		`^\s*(?:[\w\s\.<>\[\]\?]+)\s+([\w_]+(?:<[\w\s,]+>)?)\s*\([^)]*\)`,
	)
	// Matches constructors (no return type). Catches name in group 1.
	csharpConstructorRegex = regexp.MustCompile(
		`^\s*(?:public|private|protected|internal|static)?\s*([\w_]+)\s*\([^)]*\)`,
	)
	// Matches properties. Catches name in group 1. Makes the opening brace optional on the same line.
	csharpPropertyRegex = regexp.MustCompile(
		`^\s*(?:[\w\s\.<>\[\]\?]+)\s+([\w_]+)\s*`,
	)
	// GUARD REGEX: Matches class/struct/interface/enum definitions to EXCLUDE them from method processing.
	typeDefinitionRegex = regexp.MustCompile(`^\s*(?:public|private|internal|protected|internal)?\s*(?:static|sealed|abstract)?\s*\b(class|struct|interface|enum)\s+`)
)

type CSharpProcessor struct{}

func NewCSharpProcessor() language.Processor {
	return &CSharpProcessor{}
}

func (p *CSharpProcessor) Name() string {
	return "C#"
}

func (p *CSharpProcessor) Detect(filePath string) bool {
	return strings.HasSuffix(strings.ToLower(filePath), ".cs")
}

func (p *CSharpProcessor) AnalyzeFile(filePath string, sourceLines []string) ([]model.MethodMetrics, error) {
	var methods []model.MethodMetrics
	processedLines := make(map[int]bool)
	inBlockComment := false

	for i, line := range sourceLines {
		if processedLines[i+1] {
			continue
		}

		// Handle block comment state transitions at the start of the line
		originalLine := line
		if inBlockComment {
			if endIdx := strings.Index(originalLine, "*/"); endIdx != -1 {
				inBlockComment = false
				line = originalLine[endIdx+2:] // Process the rest of the line
			} else {
				continue // Still in a block comment, skip the whole line
			}
		}

		// Remove comments for cleaner regex matching
		if startIdx := strings.Index(line, "/*"); startIdx != -1 {
			if endIdx := strings.Index(line[startIdx:], "*/"); endIdx == -1 {
				inBlockComment = true
			}
			// This part is complex; for now, we simplify by just ignoring lines with unclosed block comments.
			// A full parser would handle this better, but for regex this is a safe compromise.
		}
		if startIdx := strings.Index(line, "//"); startIdx != -1 {
			line = line[:startIdx]
		}

		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			continue
		}

		// --- The crucial guard clause ---
		// If the line is a type definition, skip it entirely.
		if typeDefinitionRegex.MatchString(trimmedLine) {
			continue
		}

		var methodName string
		// Try matching in order of specificity to avoid ambiguity
		if csharpMethodRegex.MatchString(trimmedLine) && strings.Contains(trimmedLine, "(") {
			match := csharpMethodRegex.FindStringSubmatch(trimmedLine)
			if len(match) > 1 {
				methodName = match[1]
			}
		} else if csharpConstructorRegex.MatchString(trimmedLine) && strings.Contains(trimmedLine, "(") {
			match := csharpConstructorRegex.FindStringSubmatch(trimmedLine)
			if len(match) > 1 {
				methodName = match[1]
			}
		} else if csharpPropertyRegex.MatchString(trimmedLine) && strings.Contains(originalLine, "{") {
			match := csharpPropertyRegex.FindStringSubmatch(trimmedLine)
			if len(match) > 1 {
				methodName = match[1]
			}
		}

		if methodName != "" {
			startLine := i + 1

			// Find the line containing the opening brace, which might be the current line or a subsequent one.
			braceLineIndex := -1
			for j := i; j < len(sourceLines); j++ {
				if strings.Contains(sourceLines[j], "{") {
					braceLineIndex = j
					break
				}
			}

			if braceLineIndex != -1 {
				endLine, found := utils.FindMatchingBrace(sourceLines, braceLineIndex)
				if found {
					methods = append(methods, model.MethodMetrics{
						Name:      methodName,
						StartLine: startLine, // The line with the signature
						EndLine:   endLine,   // The line with the closing brace
					})
					// Mark the entire block as processed to prevent inner members from being re-matched.
					for k := startLine; k <= endLine; k++ {
						processedLines[k] = true
					}
				}
			}
		}
	}
	return methods, nil
}
