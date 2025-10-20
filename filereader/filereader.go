package filereader

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/htmlindex"
	"golang.org/x/text/transform"
)

type Reader interface {
	ReadFile(path string) ([]string, error)
	CountLines(path string) (int, error)
	Stat(name string) (fs.FileInfo, error)
	ReadDir(name string) ([]fs.DirEntry, error)
}

// DetectEncoding attempts to detect the encoding of a file by sniffing its BOM.
// It is moved here to break the import cycle with the utils package.
func DetectEncoding(filePath string) (encoding.Encoding, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	bom := make([]byte, 4)
	n, err := f.Read(bom)
	if err != nil && err != io.EOF {
		return nil, err
	}
	bom = bom[:n]

	if len(bom) >= 3 && bom[0] == 0xEF && bom[1] == 0xBB && bom[2] == 0xBF {
		return htmlindex.Get("utf-8")
	}
	if len(bom) >= 2 && bom[0] == 0xFF && bom[1] == 0xFE {
		return htmlindex.Get("utf-16le")
	}
	if len(bom) >= 2 && bom[0] == 0xFE && bom[1] == 0xFF {
		return htmlindex.Get("utf-16be")
	}

	return htmlindex.Get("utf-8")
}

func CountLinesInFile(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}
	return lineCount, scanner.Err()
}

func ReadLinesInFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Use the local DetectEncoding function
	detectedEncoding, err := DetectEncoding(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not detect encoding for %s: %v. Assuming UTF-8.\n", filePath, err)
	}

	var reader io.Reader = file
	if detectedEncoding != nil {
		_, seekErr := file.Seek(0, io.SeekStart)
		if seekErr != nil {
			return nil, fmt.Errorf("failed to seek file %s after encoding detection: %w", filePath, seekErr)
		}
		decoder := detectedEncoding.NewDecoder()
		reader = transform.NewReader(file, decoder)
	}

	var lines []string
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
