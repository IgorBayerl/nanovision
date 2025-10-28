/*
This script automates fetching and setting up Tree-sitter grammars for a Go project.
It is a necessary workaround for `go mod vendor`, which fails to include the required C source files.

The script performs three main steps:

1.  **Fetch Runtimes:** Downloads the `go-tree-sitter` Go wrapper and the core `tree-sitter` C library, arranging them into a single, CGO-compatible package.

2.  **Apply CGO Fixes:** To prevent common build errors, it performs two critical fixes:
  - Deletes `lib.c` (a unity build file) to avoid "multiple definition" linker errors.
  - Patches Go files to remove the corresponding `#include "lib.c"` line.

3.  **Fetch Grammars:** It reads `grammars.yaml` and downloads each specified language grammar, extracting its C source (`src/`) and Go bindings (`bindings/go/`).
*/
package main

import (
	"archive/zip"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

const (
	defaultConfigPath = "grammars.yaml"

	// Fixed sources for runtime.
	goTreeSitterRepo = "https://github.com/tree-sitter/go-tree-sitter"
	treeSitterCore   = "https://github.com/tree-sitter/tree-sitter"

	// Pin if you prefer reproducibility. "v0.24.0" works well; "master" also OK.
	goTreeSitterRef = "master"
	coreRef         = "master"
)

type Config struct {
	SchemaVersion int    `yaml:"nanovision_grammars version"`
	BaseDir       string `yaml:"base_dir"`
	Packages      []Pkg  `yaml:"packages"`
}
type Pkg struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
	Ref  string `yaml:"ref"`
}

func main() {
	force := false
	flag.BoolVar(&force, "f", false, "force re-download (overwrite existing)")
	flag.BoolVar(&force, "force", false, "force re-download (overwrite existing)")
	flag.Parse()

	cfg, err := loadConfig(defaultConfigPath)
	die(err)
	if strings.TrimSpace(cfg.BaseDir) == "" {
		die(errors.New("base_dir is required in YAML"))
	}

	// 1) Install go-tree-sitter runtime with proper layout (no patches).
	die(installGoTreeSitter(cfg.BaseDir, force))

	// 2) Install grammars as-is (no patching): src/** + bindings/go/**
	for _, p := range cfg.Packages {
		die(installGrammar(cfg.BaseDir, p, force))
	}

	fmt.Println("All done.")
}

// ---------------- go-tree-sitter runtime ----------------

func installGoTreeSitter(baseDir string, force bool) error {
	dest := filepath.Join(baseDir, mustRepoTail(goTreeSitterRepo))
	if exists(dest) && !force {
		fmt.Printf("skip %s (already exists)\n", dest)
		return nil
	}
	_ = os.RemoveAll(dest)
	die(os.MkdirAll(dest, 0o755))

	// 1) Fetch go-tree-sitter's .go files and its own C helper files.
	gtZipURL := mustCodeload(goTreeSitterRepo, "master")
	gtData := mustHTTP(gtZipURL)
	gtZR := mustZip(gtData)
	gtTop := topDir(gtZR)
	if gtTop == "" {
		return errors.New("go-tree-sitter zip missing top-level dir")
	}

	var goFiles []string
	for _, f := range gtZR.File {
		name := toSlash(f.Name)
		if !strings.HasPrefix(name, gtTop+"/") {
			continue
		}
		rel := strings.TrimPrefix(name, gtTop+"/")

		if !strings.Contains(rel, "/") && strings.HasSuffix(rel, ".go") && !strings.HasSuffix(rel, "_test.go") {
			outPath := filepath.Join(dest, rel)
			die(extractFile(gtZR, name, outPath))
			goFiles = append(goFiles, outPath)
		} else if !strings.Contains(rel, "/") && (rel == "allocator.c" || rel == "allocator.h") {
			outPath := filepath.Join(dest, rel)
			fmt.Printf("extracting helper %s\n", rel)
			die(extractFile(gtZR, name, outPath))
		}
	}

	// 2) Fetch tree-sitter core and replicate its src and include structure.
	coreZipURL := mustCodeload(treeSitterCore, "master")
	coreData := mustHTTP(coreZipURL)
	coreZR := mustZip(coreData)
	coreTop := topDir(coreZR)
	if coreTop == "" {
		return errors.New("tree-sitter core zip missing top-level dir")
	}

	tsIncludeDir := filepath.Join(dest, "tree_sitter")
	die(os.MkdirAll(tsIncludeDir, 0o755))

	for _, f := range coreZR.File {
		name := toSlash(f.Name)
		if !strings.HasPrefix(name, coreTop+"/") {
			continue
		}
		rel := strings.TrimPrefix(name, coreTop+"/")

		// --- ROBUST LOGIC HERE ---
		// If a file is in lib/src/, copy it to the destination, preserving its sub-path.
		// This handles files at the root (e.g., lib.c) and in subdirectories (e.g., unicode/utf8.h, portable/endian.h).
		if strings.HasPrefix(rel, "lib/src/") {
			// e.g., "lib/src/portable/endian.h" -> "portable/endian.h"
			subPath := strings.TrimPrefix(rel, "lib/src/")
			if subPath != "" && !strings.HasSuffix(subPath, "/") { // Ignore the src directory itself
				out := filepath.Join(dest, filepath.FromSlash(subPath))
				die(extractDirAware(coreZR, name, out))
			}
		}

		// Handle public headers from lib/include
		if strings.HasPrefix(rel, "lib/include/tree_sitter/") {
			out := filepath.Join(tsIncludeDir, filepath.Base(rel))
			die(extractFile(coreZR, name, out))
		}
	}

	// 3) *** MODIFIED STEP ***
	// Delete the unity build file `lib.c` to prevent "multiple definition" errors,
	// then patch the Go files to remove the corresponding #include directive.
	libCPath := filepath.Join(dest, "lib.c")
	if exists(libCPath) {
		fmt.Printf("deleting unity build file: %s\n", libCPath)
		die(os.Remove(libCPath))
	}

	cflagsRegex := regexp.MustCompile(`(?m)^// #cgo CFLAGS: -I\./tree-sitter/lib/include\s*$`)
	cIncludeRegex := regexp.MustCompile(`(?m)^/\* #include ".*\.c" \*/\s*$`)
	// ADDED: Regex to find and remove the #include "lib.c" line.
	libCIncludeRegex := regexp.MustCompile(`(?m)^\s*#include "lib.c".*\n?`)

	for _, goFile := range goFiles {
		content, err := os.ReadFile(goFile)
		die(err)
		newContent := cflagsRegex.ReplaceAll(content, []byte{})
		newContent = cIncludeRegex.ReplaceAll(newContent, []byte{})
		// ADDED: Apply the new patch to remove the include line.
		newContent = libCIncludeRegex.ReplaceAll(newContent, []byte{})

		if len(newContent) < len(content) {
			fmt.Printf("patched %s\n", goFile)
			die(os.WriteFile(goFile, newContent, 0o644))
		}
	}

	fmt.Printf("installed %s\n", dest)
	return nil
}

// ---------------- grammars ----------------

func installGrammar(baseDir string, p Pkg, force bool) error {
	if strings.TrimSpace(p.URL) == "" || strings.TrimSpace(p.Ref) == "" {
		return fmt.Errorf("invalid package (url/ref required): %+v", p)
	}
	tail, err := repoTail(p.URL)
	if err != nil {
		return fmt.Errorf("parse repo %q: %w", p.URL, err)
	}
	dest := filepath.Join(baseDir, tail)
	if exists(dest) && !force {
		fmt.Printf("skip %s (already exists)\n", dest)
		return nil
	}
	_ = os.RemoveAll(dest)
	die(os.MkdirAll(dest, 0o755))

	zURL := mustCodeload(p.URL, p.Ref)
	data := mustHTTP(zURL)
	zr := mustZip(data)
	top := topDir(zr)
	if top == "" {
		return errors.New("grammar zip missing top-level dir")
	}

	foundSrc, foundGo := false, false
	for _, f := range zr.File {
		name := toSlash(f.Name)
		if !strings.HasPrefix(name, top+"/") {
			continue
		}
		rel := strings.TrimPrefix(name, top+"/")

		switch {
		case strings.HasPrefix(rel, "src/"):
			out := filepath.Join(dest, filepath.FromSlash(rel))
			die(extractDirAware(zr, name, out))
			foundSrc = true
		case strings.HasPrefix(rel, "bindings/go/"):
			if strings.HasSuffix(toSlash(name), "_test.go") {
				continue
			}
			out := filepath.Join(dest, filepath.FromSlash(rel))
			die(extractDirAware(zr, name, out))
			foundGo = true

		case rel == "LICENSE" || rel == "License" || rel == "license":
			die(extractFile(zr, name, filepath.Join(dest, "LICENSE")))
		}
	}
	if !foundSrc {
		return fmt.Errorf("%s: src/ not found", p.URL)
	}
	if !foundGo {
		return fmt.Errorf("%s: bindings/go/ not found", p.URL)
	}

	fmt.Printf("installed %s\n", dest)
	return nil

}

// ---------------- helpers ----------------

func loadConfig(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c Config
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func repoTail(raw string) (string, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}
	if !strings.EqualFold(u.Host, "github.com") {
		return "", fmt.Errorf("only github.com supported: %s", raw)
	}
	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid GitHub repo path: %s", u.Path)
	}
	return parts[len(parts)-1], nil
}
func mustRepoTail(raw string) string {
	t, err := repoTail(raw)
	if err != nil {
		die(err)
	}
	return t
}

func codeloadURL(repo, ref string) (string, error) {
	u, err := url.Parse(repo)
	if err != nil {
		return "", err
	}
	if !strings.EqualFold(u.Host, "github.com") {
		return "", fmt.Errorf("only github.com supported: %s", repo)
	}
	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid GitHub repo path: %s", u.Path)
	}
	return fmt.Sprintf("https://codeload.github.com/%s/%s/zip/%s", parts[0], parts[1], ref), nil
}
func mustCodeload(repo, ref string) string {
	u, err := codeloadURL(repo, ref)
	if err != nil {
		die(err)
	}
	return u
}

func httpGet(u string) ([]byte, error) {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "tspack/1.0 (+https://github.com/)")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, fmt.Errorf("GET %s: %s\n%s", u, resp.Status, string(body))
	}
	return io.ReadAll(resp.Body)
}
func mustHTTP(u string) []byte {
	b, err := httpGet(u)
	if err != nil {
		die(err)
	}
	return b
}

func mustZip(data []byte) *zip.Reader {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		die(err)
	}
	return zr
}

func topDir(zr *zip.Reader) string {
	for _, f := range zr.File {
		name := toSlash(f.Name)
		if i := strings.IndexRune(name, '/'); i > 0 {
			return name[:i]
		}
	}
	return ""
}

func extractDirAware(zr *zip.Reader, zipName, outPath string) error {
	if strings.HasSuffix(zipName, "/") {
		return os.MkdirAll(outPath, 0o755)
	}
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}
	return extractFile(zr, zipName, outPath)
}

func extractFile(zr *zip.Reader, zipName, outPath string) error {
	rc, f, err := openZipFile(zr, zipName)
	if err != nil {
		return err
	}
	defer rc.Close()
	if f.FileInfo().IsDir() {
		return os.MkdirAll(outPath, 0o755)
	}
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}
	tmp := outPath + ".tmp"
	out, err := os.Create(tmp)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, rc); err != nil {
		out.Close()
		_ = os.Remove(tmp)
		return err
	}
	out.Close()
	return os.Rename(tmp, outPath)
}

func openZipFile(zr *zip.Reader, zipName string) (io.ReadCloser, *zip.File, error) {
	target := toSlash(zipName)
	for _, f := range zr.File {
		if toSlash(f.Name) == target {
			rc, err := f.Open()
			return rc, f, err
		}
	}
	return nil, nil, fmt.Errorf("file %s not found in zip", zipName)
}

func toSlash(s string) string { return filepath.ToSlash(s) }
func exists(path string) bool { _, err := os.Stat(path); return err == nil }

func die(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
