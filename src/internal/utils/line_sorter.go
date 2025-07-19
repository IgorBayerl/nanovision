// Path: internal/utils/line_sorter.go (or your utils file)
package utils

import "sort"

// SortableByLineAndName defines an interface for items that can be sorted
// primarily by their first line number and secondarily by a sortable name.
type SortableByLineAndName interface {
	GetFirstLine() int
	GetSortableName() string // This should return a name suitable for lexicographical sorting.
}

// SortByLineAndName sorts a slice of items that implement SortableByLineAndName.
// Items are sorted by their first line number, then by their sortable name.
// Items with a FirstLine of 0 (or less) are typically pushed to the end or handled as per specific needs.
// This implementation pushes items with FirstLine 0 to the end if other items have non-zero FirstLine.
func SortByLineAndName[T SortableByLineAndName](slice []T) {
	sort.Slice(slice, func(i, j int) bool {
		itemI := slice[i]
		itemJ := slice[j]

		firstLineI := itemI.GetFirstLine()
		firstLineJ := itemJ.GetFirstLine()

		// Handle invalid/default line numbers by pushing them to the end.
		if firstLineI <= 0 && firstLineJ > 0 {
			return false // itemI (with 0 or less) comes after itemJ (with positive line)
		}
		if firstLineI > 0 && firstLineJ <= 0 {
			return true // itemI (with positive line) comes before itemJ (with 0 or less)
		}

		// If both are invalid/default or both are valid and equal, sort by name.
		if firstLineI == firstLineJ {
			return itemI.GetSortableName() < itemJ.GetSortableName()
		}

		// Both are valid and different, sort by line number.
		return firstLineI < firstLineJ
	})
}
