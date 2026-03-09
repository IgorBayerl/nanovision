package diffapply

import (
	"log/slog"
	"testing"

	"github.com/IgorBayerl/nanovision/internal/model"
)

func TestBuildResolver(t *testing.T) {
	logger := slog.Default()

	tests := []struct {
		name         string
		diffPaths    []string
		fileIndex    map[string]*model.FileNode
		coveredSet   map[string]bool
		diffPath     string
		wantPath     string
		wantResolved bool
	}{
		{
			name: "H1: strip one component",
			diffPaths: []string{
				"backend/src/foo.go",
				"backend/src/bar.go",
			},
			fileIndex: map[string]*model.FileNode{
				"src/foo.go": {},
				"src/bar.go": {},
			},
			coveredSet:   map[string]bool{},
			diffPath:     "backend/src/foo.go",
			wantPath:     "src/foo.go",
			wantResolved: true,
		},
		{
			name: "H2: common monorepo prefix",
			diffPaths: []string{
				"//depot/main/backend/foo.go",
				"//depot/main/backend/bar.go",
			},
			fileIndex: map[string]*model.FileNode{
				"backend/foo.go": {},
				"backend/bar.go": {},
			},
			coveredSet:   map[string]bool{},
			diffPath:     "//depot/main/backend/foo.go",
			wantPath:     "backend/foo.go",
			wantResolved: true,
		},
		{
			name: "H3: suffix matching",
			diffPaths: []string{
				"unique/path/test.go",
			},
			fileIndex: map[string]*model.FileNode{
				"src/test.go":       {},
				"lib/other/test.go": {},
				"path/test.go":      {},
			},
			coveredSet:   map[string]bool{},
			diffPath:     "unique/path/test.go",
			wantPath:     "path/test.go",
			wantResolved: true,
		},
		{
			name: "H4: coverage tie-breaker",
			diffPaths: []string{
				"test.go",
			},
			fileIndex: map[string]*model.FileNode{
				"src/test.go": {},
				"lib/test.go": {},
			},
			coveredSet: map[string]bool{
				"lib/test.go": true,
			},
			diffPath:     "test.go",
			wantPath:     "lib/test.go",
			wantResolved: true,
		},
		{
			name: "ambiguous path",
			diffPaths: []string{
				"test.go",
			},
			fileIndex: map[string]*model.FileNode{
				"src/test.go": {},
				"lib/test.go": {},
			},
			coveredSet: map[string]bool{
				"src/test.go": true,
				"lib/test.go": true,
			},
			diffPath:     "test.go",
			wantResolved: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolver := BuildResolver(tt.diffPaths, tt.fileIndex, tt.coveredSet, logger)
			gotPath, ok := resolver.Resolve(tt.diffPath)
			if ok != tt.wantResolved {
				t.Errorf("Resolve() ok = %v, want %v", ok, tt.wantResolved)
			}
			if ok && gotPath != tt.wantPath {
				t.Errorf("Resolve() = %q, want %q", gotPath, tt.wantPath)
			}
		})
	}
}
