package assets

import (
	"embed"
	"io/fs"
)

//go:embed all:angular_report_complement/*
var htmlReportAssets embed.FS

//go:embed all:angular_frontend_spa/dist/*
var angularDistAssets embed.FS

// Contains embed of the css minified and some .js files to handle graphs generation
func AngularComplementaryAssets() (fs.FS, error) {
	return fs.Sub(htmlReportAssets, "angular_report_complement")
}

// Contains embed of the build of the angular project
func AngularDist() (fs.FS, error) {
	return fs.Sub(angularDistAssets, "angular_frontend_spa/dist")
}