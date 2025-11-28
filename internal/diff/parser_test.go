package diff

import (
	"io"
	"log/slog"
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
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

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
				// All 3 lines are added
				if len(h.AddedLineOffsets) != 3 {
					t.Errorf("got %d added lines, want 3", len(h.AddedLineOffsets))
				}
				if len(h.ModifiedLineOffsets) != 0 {
					t.Errorf("got %d modified lines, want 0", len(h.ModifiedLineOffsets))
				}

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
				// Everything is considered "Added" for coverage reporting
				if len(h.ModifiedLineOffsets) != 0 {
					t.Errorf("got %d modified lines, want 0 (all should be added)", len(h.ModifiedLineOffsets))
				}
				if len(h.AddedLineOffsets) != 3 {
					t.Errorf("got %d added lines, want 3", len(h.AddedLineOffsets))
				}

				// Lines at offsets 1, 2, 3 are new/changed
				wantAdd := []int{1, 2, 3}
				if !reflect.DeepEqual(h.AddedLineOffsets, wantAdd) {
					t.Errorf("added line offsets = %v, want %v", h.AddedLineOffsets, wantAdd)
				}

				checkOffsets(t, h, "modified_file")
			},
		},
		{
			name: "perforce_style_diff",
			input: `--- //depot/main/src/App.cs#5	2023-10-25 14:00:00
+++ C:\Workspaces\Project\src\App.cs	2023-10-25 14:05:00
@@ -10,2 +10,2 @@
-        var x = 1;
+        var x = 2;`,
			validate: func(t *testing.T, d *DiffData) {
				if len(d.Files) != 1 {
					t.Fatalf("got %d files, want 1", len(d.Files))
				}
				f := d.Files[0]

				expectedOld := "//depot/main/src/App.cs"
				if f.OldPath != expectedOld {
					t.Errorf("old path mismatch: got %q, want %q", f.OldPath, expectedOld)
				}

				expectedNew := `C:\Workspaces\Project\src\App.cs`
				if f.NewPath != expectedNew {
					t.Errorf("new path mismatch: got %q, want %q", f.NewPath, expectedNew)
				}

				if len(f.Hunks) != 1 {
					t.Fatalf("got %d hunks, want 1", len(f.Hunks))
				}
				// Update: now classified as added
				if len(f.Hunks[0].AddedLineOffsets) != 1 {
					t.Error("expected 1 added line")
				}
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

				h2 := f.Hunks[1]
				if len(h2.ModifiedLineOffsets) != 0 {
					t.Errorf("hunk 2: got %d modified lines, want 0", len(h2.ModifiedLineOffsets))
				}
				if len(h2.AddedLineOffsets) != 1 {
					t.Errorf("hunk 2: got %d added lines, want 1", len(h2.AddedLineOffsets))
				}
				wantAdd2 := []int{1}
				if !reflect.DeepEqual(h2.AddedLineOffsets, wantAdd2) {
					t.Errorf("hunk 2: added line offsets = %v, want %v", h2.AddedLineOffsets, wantAdd2)
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
			// The parser is resilient and skips invalid hunks instead of returning an error.
			wantErr: false,
			validate: func(t *testing.T, d *DiffData) {
				if len(d.Files) != 1 {
					t.Fatalf("got %d files, want 1", len(d.Files))
				}
				// Verify that no hunks were parsed for the file
				if len(d.Files[0].Hunks) != 0 {
					t.Errorf("got %d hunks, want 0", len(d.Files[0].Hunks))
				}
			},
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
			diff, err := Parse(tmpFile, logger)
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

func TestParse_PerforceStyleMultiFile(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	// Minimal reproduction of the issue: multiple files without "diff --git" headers.
	input := `--- //depot/path/to/file1.go	2025-01-01 10:00:00.000000000 0100
+++ D:\local\path\to\file1.go	2025-01-01 10:00:00.000000000 0100
@@ -10,1 +10,1 @@
-var x = 1
+var x = 2
--- //depot/path/to/file2.go	2025-01-01 10:00:00.000000000 0100
+++ D:\local\path\to\file2.go	2025-01-01 10:00:00.000000000 0100
@@ -5,1 +5,1 @@
-func foo() {}
+func bar() {}
`

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "repro_minimal.diff")
	if err := os.WriteFile(tmpFile, []byte(input), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	diff, err := Parse(tmpFile, logger)
	if err != nil {
		t.Fatalf("Parse() unexpected error: %v", err)
	}

	// We expect 2 files to be parsed
	if len(diff.Files) != 2 {
		t.Errorf("got %d files, want 2", len(diff.Files))
	}

	if len(diff.Files) > 0 {
		want1 := "//depot/path/to/file1.go"
		if diff.Files[0].OldPath != want1 {
			t.Errorf("File 1 OldPath mismatch: got %q, want %q", diff.Files[0].OldPath, want1)
		}
	}
	if len(diff.Files) > 1 {
		want2 := "//depot/path/to/file2.go"
		if diff.Files[1].OldPath != want2 {
			t.Errorf("File 2 OldPath mismatch: got %q, want %q", diff.Files[1].OldPath, want2)
		}
	}
}
