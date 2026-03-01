package htmlreact

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// FullReportData is the schema for the massive JSON payload
type FullReportData struct {
	Summary summaryV1             `json:"summary"`
	Details map[string]*detailsV1 `json:"details"`
}

// GenerateSingleFile creates one HTML file `index.html` with all data embedded,
// and copies the standard React assets over just like GenerateSummary.
func GenerateSingleFile(outDir string, summary summaryV1, details map[string]*detailsV1, logger Logger) error {
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return fmt.Errorf("create output dir %q: %w", outDir, err)
	}

	distFS, err := getReactDist()
	if err != nil {
		return fmt.Errorf("failed to get embedded dist FS: %w", err)
	}

	src, err := distFS.Open("index.html")
	if err != nil {
		return fmt.Errorf("open embedded index.html: %w", err)
	}
	defer src.Close()

	htmlBytes, err := io.ReadAll(src)
	if err != nil {
		return fmt.Errorf("read index.html: %w", err)
	}

	fullData := FullReportData{
		Summary: summary,
		Details: details,
	}

	var jsonBuf bytes.Buffer
	enc := json.NewEncoder(&jsonBuf)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(fullData); err != nil {
		return fmt.Errorf("encode full report data: %w", err)
	}

	scriptContent := fmt.Sprintf(
		"window.__NANOVISION_MODE__ = 'single'; window.__NANOVISION_FULL_DATA__ = %s;",
		jsonBuf.String(),
	)

	// Inject the script into <head>.
	modifiedHTML := strings.Replace(
		string(htmlBytes),
		"</head>",
		fmt.Sprintf("<script>%s</script></head>", scriptContent),
		1,
	)

	dest := filepath.Join(outDir, "index.html")
	if err := os.WriteFile(dest, []byte(modifiedHTML), 0o644); err != nil {
		return fmt.Errorf("write single file report: %w", err)
	}

	if err := copyDistAssetFiles(outDir, distFS, logger); err != nil {
		return fmt.Errorf("failed to copy dist assets: %w", err)
	}

	if logger != nil {
		logger.Infof("Generated single-file report at %s", dest)
	}
	return nil
}

func copyDistAssetFiles(outDir string, distFS fs.FS, logger Logger) error {
	return fs.WalkDir(distFS, ".", func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("walk dist at %q: %w", path, walkErr)
		}
		if path == "." || path == "index.html" || path == "details.html" {
			return nil // skip root meta entry and the html files since we just injected them
		}
		dest := filepath.Join(outDir, path)
		if d.IsDir() {
			if err := os.MkdirAll(dest, 0o755); err != nil {
				return fmt.Errorf("mkdir %q: %w", dest, err)
			}
			return nil
		}

		src, err := distFS.Open(path)
		if err != nil {
			return fmt.Errorf("open embedded %q: %w", path, err)
		}
		defer src.Close()

		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return fmt.Errorf("mkdir parent %q: %w", filepath.Dir(dest), err)
		}

		tmp := dest + ".tmp"
		dst, err := os.Create(tmp)
		if err != nil {
			return fmt.Errorf("create %q: %w", tmp, err)
		}

		_, copyErr := io.Copy(dst, src)
		closeErr := dst.Close()
		if copyErr != nil {
			_ = os.Remove(tmp)
			return fmt.Errorf("copy to %q: %w", tmp, copyErr)
		}
		if closeErr != nil {
			_ = os.Remove(tmp)
			return fmt.Errorf("close %q: %w", tmp, closeErr)
		}

		if err := os.Rename(tmp, dest); err != nil {
			_ = os.Remove(tmp)
			return fmt.Errorf("rename %q -> %q: %w", tmp, dest, err)
		}

		if logger != nil {
			logger.Debugf("wrote %s", dest)
		}
		return nil
	})
}
