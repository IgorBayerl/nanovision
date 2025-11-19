package diff

import "testing"

func TestNormalize(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "git_style_a_prefix",
			path: "a/src/x.go",
			want: "src/x.go",
		},
		{
			name: "git_style_b_prefix",
			path: "b/src/x.go",
			want: "src/x.go",
		},
		{
			name: "windows_backslashes",
			path: `b\src\x.go`,
			want: "src/x.go",
		},
		{
			name: "leading_dot_slash",
			path: "./foo/bar.go",
			want: "foo/bar.go",
		},
		{
			name: "multiple_prefixes",
			path: "a/b/./src/x.go",
			want: "b/src/x.go",
		},
		{
			name: "clean_path",
			path: "src/x.go",
			want: "src/x.go",
		},
		{
			name: "absolute_windows_path",
			path: `C:\repo\src\x.go`,
			want: "C:/repo/src/x.go",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.path)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}
