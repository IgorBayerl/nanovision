package testutil

import (
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
)

// CmpPaths is a cmp.Option bundle that normalizes path separators to '/'
// and sorts string slices before comparison. This is ideal for comparing
// lists of file paths in a cross-platform, order-insensitive way.
var CmpPaths = cmp.Options{
	cmp.Transformer("toSlash", filepath.ToSlash),
	cmpopts.SortSlices(func(a, b string) bool { return a < b }),
}

// PathsMatch asserts that two string slices containing file paths are equivalent,
// ignoring slice order and path separator differences (e.g., '\' vs '/').
//
// It uses go-cmp for its powerful diffing capabilities and reports the
// failure using the testify/assert framework for consistent test output.
//
// Returns true if the paths match, false otherwise.
func PathsMatch(t *testing.T, want, got []string, msgAndArgs ...interface{}) bool {
	t.Helper()

	// Use go-cmp to compare the slices with our custom options
	diff := cmp.Diff(want, got, CmpPaths...)

	if diff == "" {
		return true // The slices are equivalent, assertion passes.
	}

	// The slices are different. Report the failure using testify.
	// This ensures the test output format is consistent with other testify assertions.
	assert.Fail(t, "Path lists do not match.", "(-want +got):\n%s", diff)
	return false
}
