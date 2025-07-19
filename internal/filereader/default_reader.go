package filereader

import (
	"io/fs"
	"os"
)

type DefaultReader struct{}

func NewDefaultReader() Reader {
	return &DefaultReader{}
}

func (dr *DefaultReader) ReadFile(path string) ([]string, error) {
	return ReadLinesInFile(path)
}

func (dr *DefaultReader) CountLines(path string) (int, error) {
	return CountLinesInFile(path)
}

func (dr *DefaultReader) Stat(name string) (fs.FileInfo, error) {
	return os.Stat(name)
}
