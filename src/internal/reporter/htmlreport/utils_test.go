package htmlreport

import (
	"strings"
	"testing"

	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/model"
)

// TestGenerateUniqueFilename tests the generateUniqueFilename function.
func TestGenerateUniqueFilename(t *testing.T) {
	tests := []struct {
		name              string
		assemblyShortName string
		className         string
		existingFilenames map[string]struct{}
		want              string
		wantExistingCount int // How many entries we expect in existingFilenames after the call
	}{
		{
			name:              "simple case, no existing",
			assemblyShortName: "MyAssembly",
			className:         "MyClass",
			existingFilenames: make(map[string]struct{}),
			want:              "MyAssemblyMyClass.html",
			wantExistingCount: 1,
		},
		{
			name:              "with namespace, no existing",
			assemblyShortName: "MyAssembly",
			className:         "MyNamespace.Core.MyClass",
			existingFilenames: make(map[string]struct{}),
			want:              "MyAssemblyMyClass.html",
			wantExistingCount: 1,
		},
		{
			name:              "filename collision, multiple existing",
			assemblyShortName: "MyAssembly",
			className:         "MyClass",
			existingFilenames: map[string]struct{}{
				"myassemblymyclass.html":  {},
				"myassemblymyclass2.html": {},
			},
			want:              "MyAssemblyMyClass3.html",
			wantExistingCount: 3,
		},
		{
			name:              "empty assembly name",
			assemblyShortName: "",
			className:         "MyNamespace.MyClass",
			existingFilenames: make(map[string]struct{}),
			want:              "MyClass.html",
			wantExistingCount: 1,
		},
		{
			name:              "empty class name (after processing)",
			assemblyShortName: "MyAssembly",
			className:         "MyNamespace.", // results in empty processedClassName
			existingFilenames: make(map[string]struct{}),
			want:              "MyAssembly.html", // baseName becomes just assemblyShortName
			wantExistingCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateUniqueFilename(tt.assemblyShortName, tt.className, tt.existingFilenames)
			if got != tt.want {
				t.Errorf("generateUniqueFilename() got = %v, want %v", got, tt.want)
			}
			if len(tt.existingFilenames) != tt.wantExistingCount {
				t.Errorf("generateUniqueFilename() modified existingFilenames to count %d, want %d. Map: %v", len(tt.existingFilenames), tt.wantExistingCount, tt.existingFilenames)
			}
			// Check if the generated filename (lowercase) is indeed in the map
			if _, ok := tt.existingFilenames[strings.ToLower(tt.want)]; !ok {
				t.Errorf("generateUniqueFilename() expected filename %s (lowercase) to be in existingFilenames map, but it was not. Map: %v", strings.ToLower(tt.want), tt.existingFilenames)
			}
		})
	}
}
func TestCountTotalClasses(t *testing.T) {
	tests := []struct {
		name       string
		assemblies []model.Assembly
		want       int
	}{
		{"no assemblies", []model.Assembly{}, 0},
		{"one assembly no classes", []model.Assembly{{Name: "A1"}}, 0},
		{
			"multiple assemblies with classes",
			[]model.Assembly{
				{Name: "A1", Classes: []model.Class{{Name: "C1"}, {Name: "C2"}}},
				{Name: "A2", Classes: []model.Class{{Name: "C3"}}},
			},
			3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := countTotalClasses(tt.assemblies); got != tt.want {
				t.Errorf("countTotalClasses() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCountUniqueFiles(t *testing.T) {
	tests := []struct {
		name       string
		assemblies []model.Assembly
		want       int
	}{
		{"no assemblies", []model.Assembly{}, 0},
		{
			"single file",
			[]model.Assembly{{
				Classes: []model.Class{{
					Files: []model.CodeFile{{Path: "file1.cs"}},
				}},
			}},
			1,
		},
		{
			"multiple unique files",
			[]model.Assembly{{
				Classes: []model.Class{
					{Files: []model.CodeFile{{Path: "file1.cs"}}},
					{Files: []model.CodeFile{{Path: "file2.cs"}}},
				},
			}},
			2,
		},
		{
			"duplicate files across classes/assemblies",
			[]model.Assembly{
				{Classes: []model.Class{
					{Files: []model.CodeFile{{Path: "file1.cs"}}},
				}},
				{Classes: []model.Class{
					{Files: []model.CodeFile{{Path: "file1.cs"}}}, // Duplicate
					{Files: []model.CodeFile{{Path: "file2.cs"}}},
				}},
			},
			2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := countUniqueFiles(tt.assemblies); got != tt.want {
				t.Errorf("countUniqueFiles() = %v, want %v", got, tt.want)
			}
		})
	}
}
