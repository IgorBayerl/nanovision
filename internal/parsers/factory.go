package parsers

import (
	"fmt"
	"log/slog"
)

type ParserFactory struct {
	parsers []IParser
	logger  *slog.Logger
}

func NewParserFactory(logger *slog.Logger, parsers ...IParser) *ParserFactory {
	if logger == nil {
		logger = slog.Default()
	}
	return &ParserFactory{
		parsers: parsers,
		logger:  logger,
	}
}

func (f *ParserFactory) FindParserForFile(filePath string) (IParser, error) {
	for _, p := range f.parsers {
		f.logger.Debug("Checking parser compatibility", "language", p.Name(), "file", filePath)

		if p.SupportsFile(filePath) {
			f.logger.Info("Parser found", "language", p.Name(), "file", filePath)
			return p, nil
		}
	}
	return nil, fmt.Errorf("no suitable parser found for file: %s", filePath)
}
