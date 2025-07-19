package htmlreport

import (
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/assets"
	"golang.org/x/net/html"
)

// initializeAssets initializes and sets up all required assets for the HTML report.
// It copies static and Angular assets from the embedded filesystem to the output directory
// and parses the embedded Angular index.html to extract critical CSS and JavaScript file references.
// Returns an error if any critical operation fails.
func (b *HtmlReportBuilder) initializeAssets() error {
	if err := b.copyStaticAssets(); err != nil {
		return fmt.Errorf("failed to copy static assets: %w", err)
	}

	if err := b.copyAngularAssets(b.OutputDir); err != nil {
		return fmt.Errorf("failed to copy angular assets: %w", err)
	}

	angularFS, err := assets.AngularDist()
	if err != nil {
		return fmt.Errorf("could not get embedded angular assets: %w", err)
	}

	indexFileReader, err := angularFS.Open("index.html")
	if err != nil {
		return fmt.Errorf("failed to open embedded index.html: %w", err)
	}
	defer indexFileReader.Close()

	cssFile, runtimeJsFile, polyfillsJsFile, mainJsFile, err := b.parseAngularIndexHTML(indexFileReader)
	if err != nil {
		return fmt.Errorf("failed to parse embedded Angular index.html: %w", err)
	}

	if cssFile == "" || runtimeJsFile == "" || mainJsFile == "" {
		return fmt.Errorf("missing critical Angular assets from index.html (css: '%s', runtime: '%s', main: '%s')", cssFile, runtimeJsFile, mainJsFile)
	}

	b.angularCssFile = cssFile

	var jsBuilder strings.Builder

	runtimeContent, err := fs.ReadFile(angularFS, runtimeJsFile)
	if err != nil {
		return fmt.Errorf("failed to read embedded Angular runtime JS file %s: %w", runtimeJsFile, err)
	}
	jsBuilder.Write(runtimeContent)
	jsBuilder.WriteString(";\n\n") // Add semicolon for safe concatenation

	// Read Polyfills JS from the embedded FS (if it exists)
	if polyfillsJsFile != "" {
		polyfillsContent, err := fs.ReadFile(angularFS, polyfillsJsFile)
		if err != nil {
			return fmt.Errorf("failed to read embedded Angular polyfills JS file %s: %w", polyfillsJsFile, err)
		}
		jsBuilder.Write(polyfillsContent)
		jsBuilder.WriteString(";\n\n")
	}

	mainContent, err := fs.ReadFile(angularFS, mainJsFile)
	if err != nil {
		return fmt.Errorf("failed to read embedded Angular main JS file %s: %w", mainJsFile, err)
	}
	jsBuilder.Write(mainContent)
	jsBuilder.WriteString(";\n")

	// Write the combined JavaScript to a file in the output directory on the real disk.
	b.combinedAngularJsFile = "reportgenerator.combined.js"
	combinedJsPath := filepath.Join(b.OutputDir, b.combinedAngularJsFile)
	err = os.WriteFile(combinedJsPath, []byte(jsBuilder.String()), 0644)
	if err != nil {
		return fmt.Errorf("failed to write combined Angular JS file %s: %w", combinedJsPath, err)
	}

	// Clear out the individual JS file names as they are no longer needed for the template.
	b.angularRuntimeJsFile = ""
	b.angularPolyfillsJsFile = ""
	b.angularMainJsFile = ""

	return nil
}

// copyStaticAssets copies static asset files from the embedded filesystem
// to the report's output directory. It also combines custom CSS files into a single report.css.
func (b *HtmlReportBuilder) copyStaticAssets() error {
	angularComplementsFS, err := assets.AngularComplementaryAssets()
	if err != nil {
		return fmt.Errorf("could not get embedded angular assets: %w", err)
	}

	filesToCopy := []string{
		"custom.css",
		"custom.js",
		"chartist.min.css",
		"chartist.min.js",
		"custom-azurepipelines.css",
		"custom-azurepipelines_adaptive.css",
		"custom-azurepipelines_dark.css",
		"custom_adaptive.css",
		"custom_bluered.css",
		"custom_dark.css",
	}

	for _, fileName := range filesToCopy {
		destinationPath := filepath.Join(b.OutputDir, fileName)

		// Read the file content from the embedded filesystem.
		content, err := fs.ReadFile(angularComplementsFS, fileName)
		if err != nil {
			// It's okay to warn and skip if a non-critical asset is missing.
			slog.Warn("Failed to read embedded asset, skipping", "asset", fileName, "error", err)
			continue
		}

		// Write the content to the destination file on the real disk.
		if err := os.WriteFile(destinationPath, content, 0644); err != nil {
			return fmt.Errorf("failed to write asset %s to output directory: %w", destinationPath, err)
		}
	}

	// Combine custom.css and custom_dark.css into a single report.css file.
	customCSSBytes, err := fs.ReadFile(angularComplementsFS, "custom.css")
	if err != nil {
		slog.Warn("Failed to read custom.css for combining into report.css", "error", err)
	}

	customDarkCSSBytes, err := fs.ReadFile(angularComplementsFS, "custom_dark.css")
	if err != nil {
		slog.Warn("Failed to read custom_dark.css for combining into report.css", "error", err)
	}

	var combinedCSSBuilder strings.Builder
	combinedCSSBuilder.Write(customCSSBytes)
	combinedCSSBuilder.WriteString("\n")
	combinedCSSBuilder.Write(customDarkCSSBytes)

	if combinedCSSBuilder.Len() > 0 {
		err = os.WriteFile(filepath.Join(b.OutputDir, "report.css"), []byte(combinedCSSBuilder.String()), 0644)
		if err != nil {
			return fmt.Errorf("failed to write combined report.css: %w", err)
		}
	} else {
		slog.Warn("custom.css and custom_dark.css were not found; report.css may be missing or incomplete")
	}

	return nil
}

// copyAngularAssets recursively copies all files from the embedded Angular app's dist filesystem
// to the report's output directory on the real disk, preserving the directory structure.
func (b *HtmlReportBuilder) copyAngularAssets(outputDir string) error {
	// Get the embedded filesystem containing the compiled Angular application.
	angularDistFS, err := assets.AngularDist()
	if err != nil {
		return fmt.Errorf("could not get embedded angular assets: %w", err)
	}
	// Walk the embedded filesystem and copy each file and directory.
	// The root "." refers to the root of the embedded filesystem.
	return fs.WalkDir(angularDistFS, ".", func(path string, directoryEntry fs.DirEntry, walkError error) error {
		if walkError != nil {
			return fmt.Errorf("error accessing path %s during walk: %w", path, walkError)
		}

		destinationPath := filepath.Join(outputDir, path)

		if directoryEntry.IsDir() {
			// If it's a directory, create it on the real disk.
			if err := os.MkdirAll(destinationPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", destinationPath, err)
			}
		} else {
			// If it's a file, read it from the embedded FS and write it to the real disk.
			sourceFile, err := angularDistFS.Open(path)
			if err != nil {
				return fmt.Errorf("failed to open embedded file %s: %w", path, err)
			}
			defer sourceFile.Close()

			destinationFile, err := os.Create(destinationPath)
			if err != nil {
				return fmt.Errorf("failed to create destination file %s: %w", destinationPath, err)
			}
			defer destinationFile.Close()

			if _, err := io.Copy(destinationFile, sourceFile); err != nil {
				return fmt.Errorf("failed to copy file content to %s: %w", destinationPath, err)
			}
		}
		return nil
	})
}

// parseAngularIndexHTML parses the Angular index.html file and extracts references to
// critical assets including CSS and JavaScript files (runtime, polyfills, and main).
// Returns the file paths for CSS, runtime JS, polyfills JS, and main JS files.
// Returns an error if the file cannot be opened or parsed.
// parseAngularIndexHTML parses the content of the Angular index.html file from a reader
// and extracts references to critical assets: the main CSS file and JavaScript modules.
func (b *HtmlReportBuilder) parseAngularIndexHTML(reader io.Reader) (cssFile, runtimeJs, polyfillsJs, mainJs string, err error) {
	// The function no longer opens a file; it directly parses the provided reader.
	document, err := html.Parse(reader)
	if err != nil {
		return "", "", "", "", fmt.Errorf("failed to parse HTML from reader: %w", err)
	}

	var findAssets func(*html.Node)
	findAssets = func(node *html.Node) {
		if node.Type == html.ElementNode {
			if node.Data == "link" {
				isStylesheet := false
				var href string
				for _, attr := range node.Attr {
					if attr.Key == "rel" && attr.Val == "stylesheet" {
						isStylesheet = true
					}
					if attr.Key == "href" {
						href = attr.Val
					}
				}
				if isStylesheet && href != "" {
					cssFile = href
				}
			} else if node.Data == "script" {
				var src string
				isModule := false
				for _, attr := range node.Attr {
					if attr.Key == "src" {
						src = attr.Val
					}
					if attr.Key == "type" && attr.Val == "module" {
						isModule = true
					}
				}
				if src != "" && isModule {
					baseSrc := filepath.Base(src)
					if strings.HasPrefix(baseSrc, "runtime.") && strings.HasSuffix(baseSrc, ".js") {
						runtimeJs = src
					} else if strings.HasPrefix(baseSrc, "polyfills.") && strings.HasSuffix(baseSrc, ".js") {
						polyfillsJs = src
					} else if strings.HasPrefix(baseSrc, "main.") && strings.HasSuffix(baseSrc, ".js") {
						mainJs = src
					}
				}
			}
		}
		// Recursively traverse the HTML node tree.
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			findAssets(child)
		}
	}

	findAssets(document)
	return
}
