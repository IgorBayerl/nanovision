package diff

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// checkOffsets verifies that offsets are within valid range and don't overlap
func checkOffsets(t *testing.T, h Hunk, prefix string) {
	// Create a map to check for duplicate offsets
	seen := make(map[int]bool)

	// Check added lines
	for _, offset := range h.AddedLineOffsets {
		if offset < 0 || offset >= h.NewLines {
			t.Errorf("%s: added line offset %d outside valid range [0,%d)", prefix, offset, h.NewLines)
		}
		if seen[offset] {
			t.Errorf("%s: duplicate offset %d", prefix, offset)
		}
		seen[offset] = true
	}

	// Check modified lines
	for _, offset := range h.ModifiedLineOffsets {
		if offset < 0 || offset >= h.NewLines {
			t.Errorf("%s: modified line offset %d outside valid range [0,%d)", prefix, offset, h.NewLines)
		}
		if seen[offset] {
			t.Errorf("%s: duplicate offset %d", prefix, offset)
		}
		seen[offset] = true
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  bool
		validate func(*testing.T, *DiffData)
	}{
		{
			name: "new_file",
			input: `diff --git a/test.txt b/test.txt
new file mode 100644
--- /dev/null
+++ b/test.txt
@@ -0,0 +1,3 @@
+line1
+line2
+line3`,
			validate: func(t *testing.T, d *DiffData) {
				if len(d.Files) != 1 {
					t.Fatalf("got %d files, want 1", len(d.Files))
				}
				f := d.Files[0]
				if f.Kind != "added" {
					t.Errorf("got kind %q, want 'added'", f.Kind)
				}
				if f.OldPath != "/dev/null" {
					t.Errorf("got old path %q, want '/dev/null'", f.OldPath)
				}
				if f.NewPath != "test.txt" {
					t.Errorf("got new path %q, want 'test.txt'", f.NewPath)
				}

				h := f.Hunks[0]
				if len(h.AddedLineOffsets) != 3 {
					t.Errorf("got %d added lines, want 3", len(h.AddedLineOffsets))
				}
				if len(h.ModifiedLineOffsets) != 0 {
					t.Errorf("got %d modified lines, want 0", len(h.ModifiedLineOffsets))
				}

				// Check offsets are correct (should be 0, 1, 2 for consecutive additions)
				want := []int{0, 1, 2}
				if !reflect.DeepEqual(h.AddedLineOffsets, want) {
					t.Errorf("added line offsets = %v, want %v", h.AddedLineOffsets, want)
				}

				checkOffsets(t, h, "new_file")
			},
		},
		{
			name: "modified_file",
			input: `diff --git a/test.txt b/test.txt
--- a/test.txt
+++ b/test.txt
@@ -1,3 +1,4 @@
 unchanged1
-old2
+new2
-old3
+new3
+added4`,
			validate: func(t *testing.T, d *DiffData) {
				if len(d.Files) != 1 {
					t.Fatalf("got %d files, want 1", len(d.Files))
				}
				f := d.Files[0]
				if f.Kind != "modified" {
					t.Errorf("got kind %q, want 'modified'", f.Kind)
				}
				h := f.Hunks[0]
				if len(h.ModifiedLineOffsets) != 2 {
					t.Errorf("got %d modified lines, want 2", len(h.ModifiedLineOffsets))
				}
				if len(h.AddedLineOffsets) != 1 {
					t.Errorf("got %d added lines, want 1", len(h.AddedLineOffsets))
				}

				// Modified lines should be at positions 1 and 2, added line at 3
				wantMod := []int{1, 2}
				wantAdd := []int{3}
				if !reflect.DeepEqual(h.ModifiedLineOffsets, wantMod) {
					t.Errorf("modified line offsets = %v, want %v", h.ModifiedLineOffsets, wantMod)
				}
				if !reflect.DeepEqual(h.AddedLineOffsets, wantAdd) {
					t.Errorf("added line offsets = %v, want %v", h.AddedLineOffsets, wantAdd)
				}

				checkOffsets(t, h, "modified_file")
			},
		},
		{
			name: "multiple_hunks",
			input: `diff --git a/test.txt b/test.txt
--- a/test.txt
+++ b/test.txt
@@ -1,2 +1,3 @@
 unchanged1
+added2
 unchanged2
@@ -10,2 +11,3 @@
 unchanged3
-old4
+new4
 unchanged5`,
			validate: func(t *testing.T, d *DiffData) {
				if len(d.Files) != 1 {
					t.Fatalf("got %d files, want 1", len(d.Files))
				}
				f := d.Files[0]
				if len(f.Hunks) != 2 {
					t.Fatalf("got %d hunks, want 2", len(f.Hunks))
				}

				// First hunk: one added line at position 1
				h1 := f.Hunks[0]
				if len(h1.AddedLineOffsets) != 1 {
					t.Errorf("hunk 1: got %d added lines, want 1", len(h1.AddedLineOffsets))
				}
				wantAdd := []int{1}
				if !reflect.DeepEqual(h1.AddedLineOffsets, wantAdd) {
					t.Errorf("hunk 1: added line offsets = %v, want %v", h1.AddedLineOffsets, wantAdd)
				}
				checkOffsets(t, h1, "multiple_hunks[0]")

				// Second hunk: one modified line at position 1
				h2 := f.Hunks[1]
				if len(h2.ModifiedLineOffsets) != 1 {
					t.Errorf("hunk 2: got %d modified lines, want 1", len(h2.ModifiedLineOffsets))
				}
				wantMod := []int{1}
				if !reflect.DeepEqual(h2.ModifiedLineOffsets, wantMod) {
					t.Errorf("hunk 2: modified line offsets = %v, want %v", h2.ModifiedLineOffsets, wantMod)
				}
				checkOffsets(t, h2, "multiple_hunks[1]")
			},
		},
		{
			name: "renamed_file",
			input: `diff --git a/old.txt b/new.txt
--- a/old.txt
+++ b/new.txt
@@ -1,2 +1,3 @@
 unchanged1
-old2
+new2
+added3`,
			validate: func(t *testing.T, d *DiffData) {
				if len(d.Files) != 1 {
					t.Fatalf("got %d files, want 1", len(d.Files))
				}
				f := d.Files[0]
				if f.Kind != "modified" {
					t.Errorf("got kind %q, want 'modified'", f.Kind)
				}
				if f.NewPath != "new.txt" {
					t.Errorf("got new path %q, want 'new.txt'", f.NewPath)
				}
			},
		},
		{
			name: "invalid_hunk_header",
			input: `diff --git a/test.txt b/test.txt
--- a/test.txt
+++ b/test.txt
@@ broken header format @@
+new line
@@ -1 invalid +2,3 @@
+another line`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary file with the test input
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "test.diff")
			if err := os.WriteFile(tmpFile, []byte(tt.input), 0644); err != nil {
				t.Fatalf("failed to write test file: %v", err)
			}

			// Parse the diff
			diff, err := Parse(tmpFile)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			// Run validation
			if tt.validate != nil {
				tt.validate(t, diff)
			}
		})
	}
}
