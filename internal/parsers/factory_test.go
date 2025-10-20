package parsers_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/IgorBayerl/AdlerCov/filereader"
	"github.com/IgorBayerl/AdlerCov/internal/parsers"
	"github.com/IgorBayerl/AdlerCov/internal/parsers/parser_cobertura"
	"github.com/IgorBayerl/AdlerCov/internal/parsers/parser_gocover"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ParserFactory_FindParserForFile(t *testing.T) {
	tmpDir := t.TempDir()

	coberturaFile := filepath.Join(tmpDir, "cobertura.xml")
	_ = os.WriteFile(coberturaFile, []byte(`<?xml version="1.0" ?><coverage/>`), 0644)

	gocoverFile := filepath.Join(tmpDir, "gocover.out")
	_ = os.WriteFile(gocoverFile, []byte(`mode: set`), 0644)

	unknownFile := filepath.Join(tmpDir, "unknown.txt")
	_ = os.WriteFile(unknownFile, []byte(`some data`), 0644)

	fileReader := filereader.NewDefaultReader()
	factory := parsers.NewParserFactory(
		parser_cobertura.NewCoberturaParser(fileReader),
		parser_gocover.NewGoCoverParser(fileReader),
	)

	testCases := []struct {
		name         string
		filePath     string
		expectedType string
		expectError  bool
	}{
		{
			name:         "Should select CoberturaParser for Cobertura XML",
			filePath:     coberturaFile,
			expectedType: "Cobertura",
			expectError:  false,
		},
		{
			name:         "Should select GoCoverParser for Go cover profile",
			filePath:     gocoverFile,
			expectedType: "GoCover",
			expectError:  false,
		},
		{
			name:        "Should return error for unknown file type",
			filePath:    unknownFile,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser, err := factory.FindParserForFile(tc.filePath)

			if tc.expectError {
				assert.Error(t, err, "Expected an error for unsupported file type")
				assert.Nil(t, parser, "Parser should be nil on error")
			} else {
				require.NoError(t, err, "Expected no error for supported file type")
				require.NotNil(t, parser, "Parser should not be nil for supported file type")
				assert.Equal(t, tc.expectedType, parser.Name(), "The wrong parser type was selected")
			}
		})
	}
}
