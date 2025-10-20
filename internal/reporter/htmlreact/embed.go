package htmlreact

import (
	"embed"
	"io/fs"
)

// This next line will throw an error if there is no content inside of assets/dist
//
//go:embed all:assets/dist
var reactDist embed.FS

// getReactDist returns an fs rooted at the dist directory.
func getReactDist() (fs.FS, error) {
	return fs.Sub(reactDist, "assets/dist")
}
