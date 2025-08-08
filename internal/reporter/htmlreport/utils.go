package htmlreport

import (
	"fmt"
	"math"
	"path/filepath"
	"strings"

	"github.com/IgorBayerl/AdlerCov/internal/utils"
)

const maxFilenameLengthBase = 95

func countTotalClasses(assemblies []Assembly) int {
	count := 0
	for _, asm := range assemblies {
		count += len(asm.Classes)
	}
	return count
}

func countUniqueFiles(assemblies []Assembly) int {
	if len(assemblies) == 0 {
		return 0
	}

	var allFiles []CodeFile
	for _, asm := range assemblies {
		for _, cls := range asm.Classes {
			allFiles = append(allFiles, cls.Files...)
		}
	}

	distinctFiles := utils.DistinctBy(allFiles, func(file CodeFile) string {
		return file.Path // Assuming Path is the unique key
	})

	return len(distinctFiles)
}

func lineVisitStatusToString(status LineVisitStatus) string { // Changed parameter type
	switch status {
	case Covered: // Use Covered
		return "green"
	case NotCovered: // Use NotCovered
		return "red"
	case PartiallyCovered: // Use PartiallyCovered
		return "orange"
	default: // NotCoverable
		return "gray"
	}
}

// generateUniqueFilename creates a sanitized and unique HTML filename for a class.
// It takes assembly and class names, and a map of existing filenames to ensure uniqueness.
// The existingFilenames map is modified by this function.
func generateUniqueFilename(
	targetFilePath string, // This is the full path to the source file, e.g., "MyProject/Services/MyService.cs".
	existingFilenames map[string]struct{},
) string {
	// 1. Remove the file extension from the path first.
	// E.g., "MyProject/Services/MyService.cs" becomes "MyProject/Services/MyService".
	fileExtension := filepath.Ext(targetFilePath)
	baseName := strings.TrimSuffix(targetFilePath, fileExtension)

	// 2. Replace all path separators and any remaining dots with underscores.
	// This flattens the entire path into a single, valid string.
	// E.g., "MyProject/Services/MyService" becomes "MyProject_Services_MyService".
	sanitizedName := strings.ReplaceAll(baseName, "/", "_")
	sanitizedName = strings.ReplaceAll(sanitizedName, "\\", "_")
	sanitizedName = strings.ReplaceAll(sanitizedName, ".", "_")

	// 3. Perform a final sanitization for any other invalid characters (like spaces or symbols).
	// This ensures the filename is as clean as possible.
	sanitizedName = utils.ReplaceInvalidPathChars(sanitizedName)

	// 4. Truncate the name if it's excessively long to prevent issues with filesystem limits.
	if len(sanitizedName) > maxFilenameLengthBase {
		if maxFilenameLengthBase > 50 {
			// Preserve the start and end of the name, as they are often the most unique parts.
			sanitizedName = sanitizedName[:50] + sanitizedName[len(sanitizedName)-(maxFilenameLengthBase-50):]
		} else {
			sanitizedName = sanitizedName[:maxFilenameLengthBase]
		}
	}

	// 5. Handle any potential collisions by appending a counter.
	fileName := sanitizedName + ".html"
	normalizedFileNameToCheck := strings.ToLower(fileName)
	counter := 1

	_, exists := existingFilenames[normalizedFileNameToCheck]
	for exists {
		counter++
		fileName = fmt.Sprintf("%s_%d.html", sanitizedName, counter) // Use underscore before number
		normalizedFileNameToCheck = strings.ToLower(fileName)
		_, exists = existingFilenames[normalizedFileNameToCheck]
	}

	// 6. Add the new unique filename to the map to prevent future collisions.
	existingFilenames[normalizedFileNameToCheck] = struct{}{}

	return fileName
}

// getCoverageBarValue snaps a coverage percentage (0-100) to the nearest available CSS class value.
func getCoverageBarValue(coverage float64) int {
	if math.IsNaN(coverage) || coverage < 0 {
		return -1 // Special value for the template to hide the bar
	}

	rounded := int(math.Round(coverage))

	if rounded <= 0 {
		return 0 // Fully covered
	}
	if rounded <= 10 {
		return 10
	}
	if rounded <= 20 {
		return 20
	}
	if rounded <= 30 {
		return 30
	}
	if rounded <= 40 {
		return 40
	}
	if rounded <= 50 {
		return 50
	}
	if rounded <= 60 {
		return 60
	}
	if rounded <= 70 {
		return 70
	}
	if rounded <= 80 {
		return 80
	}
	if rounded <= 90 {
		return 90
	}

	return 100
}
