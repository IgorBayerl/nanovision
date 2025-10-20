package htmlreact

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type Logger interface {
	Debugf(string, ...any)
	Infof(string, ...any)
	Errorf(string, ...any)
}

// GenerateSummary copies the embedded React dist and writes data.js.
// If logger is nil, logging is skipped.
func GenerateSummary(outDir string, data summaryV1, logger Logger) error {
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return fmt.Errorf("create output dir %q: %w", outDir, err)
	}
	if err := copyDist(outDir, logger); err != nil {
		return fmt.Errorf("copy dist to %q: %w", outDir, err)
	}
	if err := writeSummaryDataJS(outDir, data); err != nil {
		return fmt.Errorf("write data.js: %w", err)
	}
	return nil
}

func copyDist(outDir string, logger Logger) error {
	distFS, err := getReactDist()
	if err != nil {
		return fmt.Errorf("failed to get embedded dist FS: %w", err)
	}
	return fs.WalkDir(distFS, ".", func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("walk dist at %q: %w", path, walkErr)
		}
		if path == "." {
			return nil // skip root meta entry
		}
		dest := filepath.Join(outDir, path)
		if d.IsDir() {
			if err := os.MkdirAll(dest, 0o755); err != nil {
				return fmt.Errorf("mkdir %q: %w", dest, err)
			}
			return nil
		}

		// file
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
