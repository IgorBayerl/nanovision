package utils

import (
	"io"
	"os"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/htmlindex"
)

// DetectEncoding attempts to detect the encoding of a file.
// This is a simplified placeholder. Robust detection is complex.
// Returns UTF-8 as a default if detection fails or is ambiguous.
func DetectEncoding(filePath string) (encoding.Encoding, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Read a small chunk for BOM or content sniffing
	bom := make([]byte, 4)
	n, err := f.Read(bom)
	if err != nil && err != io.EOF {
		return nil, err
	}
	bom = bom[:n]

	// Basic BOM sniffing (UTF-8, UTF-16LE, UTF-16BE)
	if len(bom) >= 3 && bom[0] == 0xEF && bom[1] == 0xBB && bom[2] == 0xBF {
		return htmlindex.Get("utf-8")
	}
	if len(bom) >= 2 && bom[0] == 0xFF && bom[1] == 0xFE {
		return htmlindex.Get("utf-16le")
	}
	if len(bom) >= 2 && bom[0] == 0xFE && bom[1] == 0xFF {
		return htmlindex.Get("utf-16be")
	}

	// If no BOM, or for more complex detection, you'd use a library
	// or try-parse with common encodings.
	// For now, default to UTF-8.
	return htmlindex.Get("utf-8")
}
