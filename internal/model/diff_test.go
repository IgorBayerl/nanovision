package model

import (
	"encoding/json"
	"testing"
)

func TestChangeKind_String(t *testing.T) {
	tests := []struct {
		kind ChangeKind
		want string
	}{
		{ChangeKindNone, ""},
		{ChangeKindAdded, "added"},
		{ChangeKindModified, "modified"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.kind.String(); got != tt.want {
				t.Errorf("ChangeKind.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChangeKind_JSONSerialization(t *testing.T) {
	tests := []struct {
		name string
		kind ChangeKind
		want string
	}{
		{"none", ChangeKindNone, `""`},
		{"added", ChangeKindAdded, `"added"`},
		{"modified", ChangeKindModified, `"modified"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test marshaling
			data, err := json.Marshal(tt.kind)
			if err != nil {
				t.Fatalf("Marshal error: %v", err)
			}
			if string(data) != tt.want {
				t.Errorf("Marshal = %v, want %v", string(data), tt.want)
			}

			// Test unmarshaling
			var kind ChangeKind
			if err := json.Unmarshal([]byte(tt.want), &kind); err != nil {
				t.Fatalf("Unmarshal error: %v", err)
			}
			if kind != tt.kind {
				t.Errorf("Unmarshal = %v, want %v", kind, tt.kind)
			}
		})
	}
}

func TestDiffInfo_JSONSerialization(t *testing.T) {
	tests := []struct {
		name     string
		node     *FileNode
		wantDiff bool
	}{
		{
			name: "file_without_diff",
			node: &FileNode{
				Name: "test.go",
				Path: "/test.go",
			},
			wantDiff: false,
		},
		{
			name: "file_with_empty_diff",
			node: &FileNode{
				Name: "test.go",
				Path: "/test.go",
				Diff: &DiffInfo{},
			},
			wantDiff: true,
		},
		{
			name: "file_with_added_lines",
			node: &FileNode{
				Name: "test.go",
				Path: "/test.go",
				Diff: &DiffInfo{
					Kind: ChangeKindAdded,
					AddedLines: map[int]bool{
						1: true,
						5: true,
					},
				},
			},
			wantDiff: true,
		},
		{
			name: "file_with_modified_lines",
			node: &FileNode{
				Name: "test.go",
				Path: "/test.go",
				Diff: &DiffInfo{
					Kind: ChangeKindModified,
					ModifiedLines: map[int]bool{
						10: true,
						20: true,
					},
				},
			},
			wantDiff: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test marshaling
			data, err := json.Marshal(tt.node)
			if err != nil {
				t.Fatalf("Marshal error: %v", err)
			}

			// Test unmarshaling
			var node FileNode
			if err := json.Unmarshal(data, &node); err != nil {
				t.Fatalf("Unmarshal error: %v", err)
			}

			// Verify diff presence
			hasDiff := node.Diff != nil
			if hasDiff != tt.wantDiff {
				t.Errorf("Diff presence = %v, want %v", hasDiff, tt.wantDiff)
			}

			// If we expect a diff, verify its contents
			if tt.wantDiff && tt.node.Diff != nil {
				if node.Diff.Kind != tt.node.Diff.Kind {
					t.Errorf("Kind = %v, want %v", node.Diff.Kind, tt.node.Diff.Kind)
				}

				// Check added lines
				if len(tt.node.Diff.AddedLines) > 0 {
					for line := range tt.node.Diff.AddedLines {
						if !node.Diff.AddedLines[line] {
							t.Errorf("Added line %d not found after unmarshal", line)
						}
					}
				}

				// Check modified lines
				if len(tt.node.Diff.ModifiedLines) > 0 {
					for line := range tt.node.Diff.ModifiedLines {
						if !node.Diff.ModifiedLines[line] {
							t.Errorf("Modified line %d not found after unmarshal", line)
						}
					}
				}
			}
		})
	}
}

func TestDiffInfo_NilSafety(t *testing.T) {
	// Test nil maps don't panic
	diff := &DiffInfo{}

	// These should not panic
	_ = diff.AddedLines[1]    // Reading from nil map
	_ = diff.ModifiedLines[1] // Reading from nil map

	// Initialize maps
	diff.AddedLines = make(map[int]bool)
	diff.ModifiedLines = make(map[int]bool)

	// Test sparse line numbers
	diff.AddedLines[1] = true
	diff.AddedLines[100] = true
	diff.ModifiedLines[50] = true
	diff.ModifiedLines[200] = true

	// Verify values
	if !diff.AddedLines[1] {
		t.Error("Expected line 1 to be added")
	}
	if !diff.AddedLines[100] {
		t.Error("Expected line 100 to be added")
	}
	if !diff.ModifiedLines[50] {
		t.Error("Expected line 50 to be modified")
	}
	if !diff.ModifiedLines[200] {
		t.Error("Expected line 200 to be modified")
	}
}
