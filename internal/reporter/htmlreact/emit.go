package htmlreact

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func writeSummaryDataJS(outDir string, data summaryV1) error {
	var jsonBuf bytes.Buffer
	enc := json.NewEncoder(&jsonBuf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(data); err != nil {
		return fmt.Errorf("marshal summary JSON: %w", err)
	}

	var buf bytes.Buffer
	buf.WriteString("window.__ADLERCOV_SUMMARY__=")
	buf.Write(bytes.TrimSpace(jsonBuf.Bytes()))
	buf.WriteString(";")

	dest := filepath.Join(outDir, "data.js")
	tmp := dest + ".tmp"
	if err := os.WriteFile(tmp, buf.Bytes(), 0o644); err != nil {
		return fmt.Errorf("write temp %q: %w", tmp, err)
	}
	if err := os.Rename(tmp, dest); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("rename %q -> %q: %w", tmp, dest, err)
	}
	return nil
}
