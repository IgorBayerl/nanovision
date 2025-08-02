package config

import (
	"fmt"
	"strings"

	"github.com/IgorBayerl/AdlerCov/internal/filtering"
	"github.com/IgorBayerl/AdlerCov/internal/language"
	"github.com/IgorBayerl/AdlerCov/internal/logging"
)

// AppConfig holds the parsed and validated configuration for the application.
// This is the single source of truth that gets passed through the application.
type AppConfig struct {
	ReportPatterns []string
	SourceDirs     []string
	ReportTypes    []string
	FileFilters    []string
	OutputDir      string
	Tag            string
	Title          string
	LogFile        string
	LogFormat      string
	Verbosity      logging.VerbosityLevel

	FileFilterInstance filtering.IFilter

	// Settings (formerly from the settings package)
	MaximumDecimalPlacesForCoverageQuotas    int
	MaximumDecimalPlacesForPercentageDisplay int

	LangFactory *language.ProcessorFactory
}

// BuildAppConfig creates the definitive AppConfig from raw flag inputs.
// It handles all parsing and validation in one place.
func BuildAppConfig(
	reportPatterns, sourceDirs, reportTypes, fileFilters, outputDir, tag, title, logFile, logFormat string,
	verbosity logging.VerbosityLevel,
	langFactory *language.ProcessorFactory,
) (*AppConfig, error) {

	patterns := strings.Split(reportPatterns, ";")
	dirs := strings.Split(sourceDirs, ";")

	if reportPatterns != "" && len(patterns) != len(dirs) {
		return nil, fmt.Errorf(
			"mismatch between number of report patterns (%d) and source directories (%d)",
			len(patterns),
			len(dirs),
		)
	}

	fileFilter, err := filtering.NewDefaultFilter(strings.Split(fileFilters, ";"), true)
	if err != nil {
		return nil, err
	}

	return &AppConfig{
		ReportPatterns: patterns,
		SourceDirs:     dirs,
		ReportTypes:    strings.Split(reportTypes, ","),
		FileFilters:    strings.Split(fileFilters, ";"),
		OutputDir:      outputDir,
		Tag:            tag,
		Title:          title,
		LogFile:        logFile,
		LogFormat:      logFormat,
		Verbosity:      verbosity,

		FileFilterInstance: fileFilter,

		// Set default values for settings
		MaximumDecimalPlacesForCoverageQuotas:    1,
		MaximumDecimalPlacesForPercentageDisplay: 0,

		LangFactory: langFactory,
	}, nil
}
