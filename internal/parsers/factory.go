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
		fmt.Println("[Find parser]: trying parser:", p.Name(), "...")
		if p.SupportsFile(filePath) {
			fmt.Println("[Find parser]: found compatible parser:", p.Name())
			return p, nil
		}
		fmt.Println("[Find parser]: ", p.Name(), " not compatible")
	}
	return nil, fmt.Errorf("no suitable parser found for file: %s", filePath)
}
