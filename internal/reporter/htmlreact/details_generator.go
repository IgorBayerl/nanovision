package htmlreact

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/IgorBayerl/AdlerCov/internal/model"
	"github.com/IgorBayerl/AdlerCov/internal/utils"
	"golang.org/x/net/html"
)

// generateDetailsPages iterates through all file nodes in the summary tree and creates a separate
// HTML details page for each one.
func generateDetailsPages(outputDir string, tree *model.SummaryTree) error {
	// A map to keep track of all file nodes in the tree.
	fileNodeMap := make(map[string]*model.FileNode)
	// Recursively collect all files from the tree.
	collectFiles(tree.Root, fileNodeMap)

	// Read the content of the embedded details.html file.
	detailsHTMLContent, err := readEmbeddedDetailsHTML()
	if err != nil {
		return err
	}

	// Process each file to generate its details page.
	for _, fileNode := range fileNodeMap {
		if err := createDetailPage(outputDir, fileNode, detailsHTMLContent); err != nil {
			// Log the error but continue processing other files.
			fmt.Fprintf(os.Stderr, "Warning: could not generate details page for '%s': %v\n", fileNode.Path, err)
		}
	}

	return nil
}

// createDetailPage generates a single HTML file with coverage details for a given file node.
func createDetailPage(outputDir string, fileNode *model.FileNode, detailsHTMLContent []byte) error {
	// Define the properties for the React component.
	props := map[string]interface{}{
		"fileName": fileNode.Path,
		"coverage": utils.CalculatePercentage(fileNode.Metrics.LinesCovered, fileNode.Metrics.LinesValid, 2),
		"risk":     getRiskStatus(utils.CalculatePercentage(fileNode.Metrics.LinesCovered, fileNode.Metrics.LinesValid, 2)),
	}

	// Convert the properties to a JSON string.
	propsJSON, err := json.Marshal(props)
	if err != nil {
		return fmt.Errorf("failed to marshal props to JSON: %w", err)
	}

	// Parse the HTML content and inject the data.
	modifiedHTML, err := injectDataIntoHTML(detailsHTMLContent, string(propsJSON))
	if err != nil {
		return err
	}

	// Create a safe and unique file name for the details page.
	detailsFileName := strings.ReplaceAll(fileNode.Path, "/", "_") + ".html"
	detailsFilePath := filepath.Join(outputDir, detailsFileName)

	// Write the modified HTML to the new file.
	return os.WriteFile(detailsFilePath, []byte(modifiedHTML), 0644)
}

// readEmbeddedDetailsHTML reads the content of the details.html file from the embedded file system.
func readEmbeddedDetailsHTML() ([]byte, error) {
	distFS, err := getReactDist()
	if err != nil {
		return nil, fmt.Errorf("failed to get embedded dist FS: %w", err)
	}

	file, err := distFS.Open("details.html")
	if err != nil {
		return nil, fmt.Errorf("failed to open embedded details.html: %w", err)
	}
	defer file.Close()

	return io.ReadAll(file)
}

// injectDataIntoHTML parses the HTML content, finds the React island, and updates its data-props attribute.
func injectDataIntoHTML(htmlContent []byte, propsJSON string) (string, error) {
	doc, err := html.Parse(bytes.NewReader(htmlContent))
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %w", err)
	}

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" {
			for _, a := range n.Attr {
				if a.Key == "data-island" && a.Val == "FileDetails" {
					// Found the target div, now find and update data-props.
					for j, attr := range n.Attr {
						if attr.Key == "data-props" {
							n.Attr[j].Val = propsJSON
							return
						}
					}
					// If data-props doesn't exist, add it.
					n.Attr = append(n.Attr, html.Attribute{Key: "data-props", Val: propsJSON})
					return
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(doc)

	var buf bytes.Buffer
	if err := html.Render(&buf, doc); err != nil {
		return "", fmt.Errorf("failed to render modified HTML: %w", err)
	}

	return buf.String(), nil
}

// collectFiles is a helper function to recursively gather all file nodes from the directory tree.
func collectFiles(dir *model.DirNode, fileMap map[string]*model.FileNode) {
	for _, file := range dir.Files {
		fileMap[file.Path] = file
	}
	for _, subDir := range dir.Subdirs {
		collectFiles(subDir, fileMap)
	}
}
