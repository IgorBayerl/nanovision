package htmlreact

import (
	"embed"
	"io/fs"
)

//go:embed all:assets/dist
var reactDist embed.FS

// getReactDist returns an fs rooted at the dist directory.
func getReactDist() (fs.FS, error) {
	return fs.Sub(reactDist, "assets/dist")
}
