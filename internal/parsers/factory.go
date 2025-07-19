package parsers

import "fmt"

type ParserFactory struct {
	parsers []IParser
}

func NewParserFactory(parsers ...IParser) *ParserFactory {
	return &ParserFactory{
		parsers: parsers,
	}
}

func (f *ParserFactory) FindParserForFile(filePath string) (IParser, error) {
	for _, p := range f.parsers {
		if p.SupportsFile(filePath) {
			return p, nil
		}
	}
	return nil, fmt.Errorf("no suitable parser found for file: %s", filePath)
}
