package settings

// Settings corresponds to C#'s ReportGenerator.Core.Settings.
// It holds various global settings that control the behavior of the report generation.
type Settings struct {
	// NumberOfReportsParsedInParallel defines how many coverage report files are parsed simultaneously.
	// Default: 1
	NumberOfReportsParsedInParallel int

	// NumberOfReportsMergedInParallel defines how many parsed report results are merged in parallel.
	// Default: 1
	NumberOfReportsMergedInParallel int

	// MaximumNumberOfHistoricCoverageFiles defines the maximum number of older history files to process.
	// Default: 100
	MaximumNumberOfHistoricCoverageFiles int

	// CachingDurationOfRemoteFilesInMinutes defines how long (in minutes) downloaded source files are cached.
	// Default: 10080 (7 days)
	CachingDurationOfRemoteFilesInMinutes int

	// DisableRiskHotspots, if true, disables the calculation and display of risk hotspots.
	// Default: false
	DisableRiskHotspots bool

	// ExcludeTestProjects, if true, attempts to exclude assemblies/classes that look like test projects (primarily for Clover reports).
	// Default: false
	ExcludeTestProjects bool

	// CreateSubdirectoryForAllReportTypes, if true, creates a subdirectory for each report type in the target directory.
	// Default: false
	CreateSubdirectoryForAllReportTypes bool

	// CustomHeadersForRemoteFiles allows specifying custom HTTP headers for fetching remote source files.
	// Format: "Header1:Value1;Header2:Value2"
	// Default: ""
	CustomHeadersForRemoteFiles string

	// DefaultAssemblyName is used for reports (like GCov, LCov) that don't inherently contain assembly names.
	// Default: "Default"
	DefaultAssemblyName string

	// MaximumDecimalPlacesForCoverageQuotas controls the precision of displayed coverage percentages.
	// Default: 1
	MaximumDecimalPlacesForCoverageQuotas int

	// MaximumDecimalPlacesForPercentageDisplay controls the precision of displayed percentages.
	// Default: 0
	MaximumDecimalPlacesForPercentageDisplay int

	// HistoryFileNamePrefix is an optional prefix for history files.
	// Default: ""
	HistoryFileNamePrefix string

	// RawMode, if true, attempts to report on compiler-generated/nested classes separately rather than merging them into parent classes.
	// This is a PRO feature in C#.
	// Default: false
	RawMode bool

	// VerbosityLevelFromConfig is a placeholder if you decide to load verbosity from settings too,
	// though it's often handled by ReportConfiguration directly from command line.
	// VerbosityLevelFromConfig string
}

// NewSettings creates a new Settings instance with default values.
// These defaults mirror those in the C# project's appsettings.json or Settings.cs defaults.
func NewSettings() *Settings {
	return &Settings{
		NumberOfReportsParsedInParallel:          1,
		NumberOfReportsMergedInParallel:          1,
		MaximumNumberOfHistoricCoverageFiles:     100,
		CachingDurationOfRemoteFilesInMinutes:    7 * 24 * 60, // 10080 minutes = 7 days
		DisableRiskHotspots:                      false,
		ExcludeTestProjects:                      false,
		CreateSubdirectoryForAllReportTypes:      false,
		CustomHeadersForRemoteFiles:              "",
		DefaultAssemblyName:                      "Default",
		MaximumDecimalPlacesForCoverageQuotas:    1,
		MaximumDecimalPlacesForPercentageDisplay: 0,
		HistoryFileNamePrefix:                    "",
		RawMode:                                  false,
	}
}
