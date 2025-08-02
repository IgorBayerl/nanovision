package lang_default_test

import (
	"testing"

	"github.com/IgorBayerl/AdlerCov/internal/language/lang_default"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultProcessor_Detect(t *testing.T) {
	// Arrange
	p := lang_default.NewDefaultProcessor()

	// Act & Assert
	// The default processor should never detect any file type.
	assert.False(t, p.Detect("file.cs"))
	assert.False(t, p.Detect("file.go"))
	assert.False(t, p.Detect("file.txt"))
	assert.False(t, p.Detect(""))
}

func TestDefaultProcessor_AnalyzeFile(t *testing.T) {
	// Arrange
	p := lang_default.NewDefaultProcessor()
	sourceCode := `
	Some random text file content.
	This could be anything.
	`
	sourceLines := []string{sourceCode}

	// Act
	methods, err := p.AnalyzeFile("test.txt", sourceLines)

	// Assert
	// The default processor should always return an empty slice and no error.
	require.NoError(t, err)
	require.NotNil(t, methods)
	assert.Empty(t, methods, "Default processor should not find any methods.")
}
