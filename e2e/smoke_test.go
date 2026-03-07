package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// projectRoot returns the absolute path to the nanovision repository root.
func projectRoot(t *testing.T) string {
	t.Helper()
	// This file lives at <root>/e2e/smoke_test.go — go up one level.
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("unable to determine project root via runtime.Caller")
	}
	return filepath.Dir(filepath.Dir(thisFile))
}

// buildBinary compiles the nanovision CLI into a temp directory and returns
// the path to the resulting executable.
func buildBinary(t *testing.T, root string) string {
	t.Helper()
	tmpDir := t.TempDir()
	binaryName := "nanovision"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}
	binaryPath := filepath.Join(tmpDir, binaryName)

	cmd := exec.Command("go", "build", "-mod=vendor", "-o", binaryPath, filepath.Join(root, "cmd", "main.go"))
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to build binary: %v\n%s", err, out)
	}
	return binaryPath
}

func TestSmokeDefaultMetrics(t *testing.T) {
	root := projectRoot(t)
	binary := buildBinary(t, root)

	coberturaXML := filepath.Join(root, "demo_projects", "cpp", "report", "cobertura", "cobertura.xml")
	sourceDir := filepath.Join(root, "demo_projects", "cpp", "project")
	outDir := t.TempDir()

	cmd := exec.Command(binary,
		"-report="+coberturaXML,
		"-sourcedirs="+sourceDir,
		"-reporttypes=TextSummary,Html",
		"-output="+outDir,
	)
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("nanovision exited with error: %v\nOutput:\n%s", err, out)
	}

	// ---- Assert text output ----
	txtPath := filepath.Join(outDir, "Summary.txt")
	txtBytes, err := os.ReadFile(txtPath)
	if err != nil {
		t.Fatalf("Summary.txt not found: %v", err)
	}
	txtContent := string(txtBytes)

	// Default file metrics include statement_coverage and branch_coverage
	// (line_coverage is NOT in defaults). Check that the defaults render.
	for _, want := range []string{"Statement Coverage", "Branch Coverage"} {
		if !strings.Contains(txtContent, want) {
			t.Errorf("text output missing %q\nContent:\n%s", want, txtContent)
		}
	}

	// ---- Assert HTML output ----
	htmlPath := filepath.Join(outDir, "index.html")
	htmlBytes, err := os.ReadFile(htmlPath)
	if err != nil {
		t.Fatalf("index.html not found: %v", err)
	}
	htmlContent := string(htmlBytes)

	// The HTML JSON embeds metric keys as object keys (e.g. "statement_coverage":).
	for _, wantKey := range []string{`"statement_coverage"`, `"branch_coverage"`} {
		if !strings.Contains(htmlContent, wantKey) {
			t.Errorf("HTML output missing metric key %s", wantKey)
		}
	}
}

func TestSmokeLineCoverageOnly(t *testing.T) {
	root := projectRoot(t)
	binary := buildBinary(t, root)

	coberturaXML := filepath.Join(root, "demo_projects", "cpp", "report", "cobertura", "cobertura.xml")
	sourceDir := filepath.Join(root, "demo_projects", "cpp", "project")
	outDir := t.TempDir()

	cmd := exec.Command(binary,
		"-report="+coberturaXML,
		"-sourcedirs="+sourceDir,
		"-reporttypes=TextSummary,Html",
		"-output="+outDir,
		"-file-metrics=line_coverage",
	)
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("nanovision exited with error: %v\nOutput:\n%s", err, out)
	}

	// ---- Assert text output has only line coverage ----
	txtBytes, err := os.ReadFile(filepath.Join(outDir, "Summary.txt"))
	if err != nil {
		t.Fatalf("Summary.txt not found: %v", err)
	}
	txtContent := string(txtBytes)

	if !strings.Contains(txtContent, "Line Coverage:") {
		t.Errorf("text output missing 'Line Coverage:'")
	}
	for _, notWant := range []string{"Statement Coverage:", "Branch Coverage:", "(Stmt)", "(Branch)"} {
		if strings.Contains(txtContent, notWant) {
			t.Errorf("text output should NOT contain %q when only line_coverage is configured\nContent:\n%s", notWant, txtContent)
		}
	}

	// ---- Assert HTML output has only line coverage ----
	htmlBytes, err := os.ReadFile(filepath.Join(outDir, "index.html"))
	if err != nil {
		t.Fatalf("index.html not found: %v", err)
	}
	htmlContent := string(htmlBytes)

	if !strings.Contains(htmlContent, `"line_coverage"`) {
		t.Errorf("HTML output missing metric key \"line_coverage\"")
	}
	for _, notWantKey := range []string{`"statement_coverage"`, `"branch_coverage"`} {
		if strings.Contains(htmlContent, notWantKey) {
			t.Errorf("HTML output should NOT contain %s when only line_coverage is configured", notWantKey)
		}
	}
}
