package assert

import (
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// Package assert contains a handful of **test-only** helpers that wrap
// `github.com/google/go-cmp/cmp` so you can write one-liner assertions
// with rich, colourised diffs and zero boiler-plate.

// ---------------------------------------------------------------------------
// Common cmp.Options, exposed so callers can reuse them.
// ---------------------------------------------------------------------------

// ToSlash converts every string (or every string inside a slice / map / struct
// field) to the Unix path separator before comparison.  Handy when the same
// test runs on both Windows and POSIX machines.
var ToSlash = cmp.Transformer("toSlash", filepath.ToSlash)

// SortPaths orders string slices lexicographically before comparing them.
// Combine it with ToSlash to make path-list comparisons order-insensitive.
var SortPaths = cmpopts.SortSlices(func(a, b string) bool { return a < b })

// CmpPaths is a convenience bundle of ToSlash + SortPaths pass it to
// Equal whenever you compare path lists.
//
//	extra := assert.CmpPaths // save typing
//	assert.Equal(t, want, got, extra...)
var CmpPaths = cmp.Options{ToSlash, SortPaths}

// ---------------------------------------------------------------------------
// Lightweight assertion helpers
// ---------------------------------------------------------------------------

// Equal fails the test if *want* and *got* differ.  The additional `cmp.Option`
// parameters let you tailor the comparison at the call-site:
//
//	assert.Equal(t, want, got, assert.CmpPaths...) // ignore order + separators
//
// Generics (`T any`) mean the same function works for slices, structs, maps,
// scalars—anything `cmp` supports.
func Equal[T any](t *testing.T, want, got T, opts ...cmp.Option) {
	t.Helper()

	// Always allow comparison of unexported struct fields; users can override
	// by passing cmp.Exporter or similar in opts.
	opts = append(opts, cmp.AllowUnexported())

	if diff := cmp.Diff(want, got, opts...); diff != "" {
		t.Fatalf("(-want +got):\n%s", diff)
	}
}

// NoError aborts the test if *err* is non-nil—syntactic sugar for the common
// case:
//
//	foo, err := DoSomething()
//	assert.NoError(t, err)
func NoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
