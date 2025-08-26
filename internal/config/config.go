package config

import (
	"fmt"
	"strings"

	"github.com/IgorBayerl/AdlerCov/filtering"
	"github.com/IgorBayerl/AdlerCov/logging"
)

type ReportInputPair struct {
	ReportPattern string
	SourceDir     string
}

type RawConfigInput struct {
	ReportPatterns string
	SourceDirs     string
	ReportTypes    string
	FileFilters    string
	OutputDir      string
	Tag            string
	Title          string
	LogFile        string
	LogFormat      string
	Verbosity      string
	Verbose        bool
}

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

	MaximumDecimalPlacesForCoverageQuotas    int
	MaximumDecimalPlacesForPercentageDisplay int

	InputPairs []ReportInputPair
}

// BuildAppConfig creates and validates the application's configuration from the single RawConfigInput object.
// It is responsible for all parsing, validation, and transformation of CLI arguments.
func BuildAppConfig(rawInput RawConfigInput) (*AppConfig, error) {
	patterns := strings.Split(rawInput.ReportPatterns, ";")
	dirs := strings.Split(rawInput.SourceDirs, ";")

	if rawInput.ReportPatterns != "" && rawInput.SourceDirs != "" && len(patterns) != len(dirs) {
		return nil, fmt.Errorf(
			"mismatch between number of report patterns (%d) and source directories (%d)",
			len(patterns),
			len(dirs),
		)
	}

	fileFilter, err := filtering.NewDefaultFilter(strings.Split(rawInput.FileFilters, ";"), true)
	if err != nil {
		return nil, err
	}

	verbosity, _ := logging.ParseVerbosity(rawInput.Verbosity)
	if rawInput.Verbose {
		verbosity = logging.Verbose
	}

	resolvedPairs := resolveInputPairs(patterns, dirs)

	return &AppConfig{
		ReportPatterns:                           patterns,
		SourceDirs:                               dirs,
		ReportTypes:                              strings.Split(rawInput.ReportTypes, ","),
		FileFilters:                              strings.Split(rawInput.FileFilters, ";"),
		OutputDir:                                rawInput.OutputDir,
		Tag:                                      rawInput.Tag,
		Title:                                    rawInput.Title,
		LogFile:                                  rawInput.LogFile,
		LogFormat:                                rawInput.LogFormat,
		Verbosity:                                verbosity,
		FileFilterInstance:                       fileFilter,
		MaximumDecimalPlacesForCoverageQuotas:    1,
		MaximumDecimalPlacesForPercentageDisplay: 0,
		InputPairs:                               resolvedPairs,
	}, nil
}

// resolveInputPairs matches slices of report patterns and source directories into
// structured pairs. The function pairs them by their order (e.g., the first report
// is matched with the first directory).
//
// It also acts as a validation step by discarding any pairs that are incomplete
// after trimming whitespace (e.g., a pattern without a directory).
//
// Example:
//
//	Given the following slices:
//	- patterns: ["reportA.xml", "reportB.out", " "]
//	- dirs:     ["./src/projectA", "./src/projectB"]
//
//	It will produce:
//	[]ReportInputPair{
//	    {ReportPattern: "reportA.xml", SourceDir: "./src/projectA"},
//	    {ReportPattern: "reportB.out", SourceDir: "./src/projectB"},
//	}
//	// The third, empty report pattern is ignored.
func resolveInputPairs(patterns []string, dirs []string) []ReportInputPair {
	var pairs []ReportInputPair
	for i := range patterns {
		trimmedPattern := strings.TrimSpace(patterns[i])
		var trimmedSourceDir string
		if i < len(dirs) {
			trimmedSourceDir = strings.TrimSpace(dirs[i])
		}

		if trimmedPattern != "" && trimmedSourceDir != "" {
			pairs = append(pairs, ReportInputPair{
				ReportPattern: trimmedPattern,
				SourceDir:     trimmedSourceDir,
			})
		}
	}
	return pairs
}
