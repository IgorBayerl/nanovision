package htmlreport

// GetTranslations returns a map of localized strings.
// Values are taken from ReportGenerator.Core/Properties/ReportResources.resx
func GetTranslations() map[string]string {
	return map[string]string{
		// Existing from your provided file
		"Coverage":                "Coverage",
		"Summary":                 "Summary",
		"Assembly":                "Assembly",
		"Class":                   "Class", // Singular, used in Class Detail page
		"Filter":                  "Filter",
		"Name":                    "Name",
		"Covered":                 "Covered",
		"Uncovered":               "Uncovered",
		"Coverable":               "Coverable",
		"Total":                   "Total",
		"Average":                 "Average",
		"Lines":                   "Lines",
		"LineCoverage":            "Line coverage", // Used as Card Title AND Row Header
		"Branches":                "Branches",
		"BranchCoverage":          "Branch coverage", // Used as Card Title AND Row Header
		"Methods":                 "Methods",
		"MethodCoverage":          "Method coverage", // Used as Card Title
		"Metrics":                 "Metrics",
		"RiskHotspots":            "Risk Hotspots", // Used as H1 Title
		"Files":                   "Files",         // Generic "Files"
		"HistoricCoverage":        "Historic Coverage",
		"ExecutionTime":           "Execution time",
		"CyclomaticComplexity":    "Cyclomatic complexity",
		"CrapScore":               "CrapScore",
		"NPathComplexity":         "NPath complexity",
		"SequenceCoverage":        "Sequence coverage",
		"BranchCoverageNUnit":     "Branch coverage (NUnit)",
		"LineCoverageNUnit":       "Line coverage (NUnit)",
		"NotCovered":              "Not covered",
		"NotCoveredMessage":       "The element is not covered by any test.",
		"PartiallyCovered":        "Partially covered",
		"PartiallyCoveredMessage": "The element is only partially covered by tests.",
		"FullyCovered":            "Fully covered",
		"FullyCoveredMessage":     "The element is fully covered by tests.",
		"LoadingData":             "Loading data...",
		"NoCoverageData":          "No coverage data available.",
		"ShowHistoricChart":       "Show historic chart",
		"HideHistoricChart":       "Hide historic chart",
		"ChartLoading":            "Chart loading...",
		"ShowAll":                 "Show all",
		"ShowLess":                "Show less",
		"ShowMore":                "Show more",
		"OverallCoverage":         "Overall coverage",
		"ApplySettings":           "Apply settings",
		"Settings":                "Settings",
		"NoData":                  "No data available.",
		"ShowHelp":                "Show help",
		"HideHelp":                "Hide help",
		"CurrentBranch":           "Current branch",
		"CompareWithBranch":       "Compare with branch",
		"NoGitInfo":               "No Git information available for comparison.",
		"NoCommonCommits":         "No common commits found for comparison.",
		"History":                 "History", // Used as H1 Title
		"AllFiles":                "All files",
		"Percentage":              "Percentage",
		"FullMethodCoverage":      "Full method coverage",
		"NoGrouping":              "No grouping",
		"ByAssembly":              "By assembly",
		"ByNamespace":             "By namespace, Level:",

		// Angular-specific keys (must match Angular casing) TODO needs review
		"Grouping":                       "Grouping:",
		"CompareHistory":                 "Compare with:",
		"Date":                           "Date",
		"AllChanges":                     "All changes",
		"LineCoverageIncreaseOnly":       "Line coverage: Increase only",
		"LineCoverageDecreaseOnly":       "Line coverage: Decrease only",
		"BranchCoverageIncreaseOnly":     "Branch coverage: Increase only",
		"BranchCoverageDecreaseOnly":     "Branch coverage: Decrease only",
		"MethodCoverageIncreaseOnly":     "Method coverage: Increase only",
		"MethodCoverageDecreaseOnly":     "Method coverage: Decrease only",
		"FullMethodCoverageIncreaseOnly": "Full method coverage: Increase only",
		"FullMethodCoverageDecreaseOnly": "Full method coverage: Decrease only",
		"SelectCoverageTypes":            "Select coverage types",
		"SelectCoverageTypesAndMetrics":  "Select coverage types & metrics",
		"CoverageTypes":                  "Coverage types",

		// == Keys specifically needed for Summary Page Cards and Titles ==
		// Page Title
		"CoverageReport": "Coverage Report",

		// GitHub Buttons
		"StarTooltip":    "Star ReportGenerator on GitHub",
		"Star":           "Star",
		"SponsorTooltip": "Sponsor ReportGenerator on GitHub",
		"Sponsor":        "Sponsor",

		// Information Card (Title already covered by generic "Information" if added, or use "Summary")
		"Information":  "Information", // Card Title
		"Parser":       "Parser",
		"Assemblies2":  "Assemblies", // C# key for 'Assemblies' count in summary
		"Classes":      "Classes",    // Plural, for 'Classes' count in summary
		"Files2":       "Files",      // C# key for 'Files' count in summary
		"CoverageDate": "Coverage date",
		"Tag":          "Tag",

		// Line Coverage Card (Title "LineCoverage" is present)
		"CoveredLines":   "Covered lines",
		"UncoveredLines": "Uncovered lines",
		"CoverableLines": "Coverable lines",
		"TotalLines":     "Total lines",
		// "LineCoverage" already present for the last row's header

		// Branch Coverage Card (Title "BranchCoverage" is present)
		"CoveredBranches2": "Covered branches", // C# key for 'Covered branches' count
		"TotalBranches":    "Total branches",
		// "BranchCoverage" already present for the last row's header

		// Method Coverage Card (Title "MethodCoverage" is present)
		"CoveredCodeElements":           "Covered methods/properties",
		"FullCoveredCodeElements":       "Fully covered methods/properties",
		"TotalCodeElements":             "Total methods/properties",
		"CodeElementCoverageQuota2":     "Method/property coverage",      // C# key for overall method coverage %
		"FullCodeElementCoverageQuota2": "Full method/property coverage", // C# key for full method coverage %
		"MethodCoverageProVersion":      "This feature is only available for sponsors.",
		"MethodCoverageProButton":       "Upgrade to PRO version",

		// Section Titles / Paragraphs
		"NoRiskHotspots":      "No risk hotspots found.",
		"Coverage3":           "Coverage", // H1 Title for the main coverage table/list section
		"NoCoveredAssemblies": "No assemblies have been covered.",
		"GeneratedBy":         "Generated by",

		// For Class Detail Page
		"MethodsProperties": "Methods/Properties",
		"Files3":            "File(s)", // Used as H1 and in info card
		"File":              "File",    // Used like "File 0: path/to/file.cs"
		"NoFilesFound":      "No files found.",
		"Line":              "Line", // Header in source code table

		// == Angular-specific keys (must match Angular casing) ==
		"collapseAll":                    "Collapse all",
		"expandAll":                      "Expand all",
		"noGrouping":                     "No grouping",
		"byAssembly":                     "By assembly",
		"byNamespace":                    "By namespace, Level:",
		"grouping":                       "Grouping:",
		"compareHistory":                 "Compare with:",
		"date":                           "Date",
		"filter":                         "Filter",
		"coverage":                       "Line coverage",
		"branchCoverage":                 "Branch coverage",
		"methodCoverage":                 "Method coverage",
		"fullMethodCoverage":             "Full method coverage",
		"metrics":                        "Metrics",
		"selectCoverageTypes":            "Select coverage types",
		"selectCoverageTypesAndMetrics":  "Select coverage types & metrics",
		"name":                           "Name",
		"covered":                        "Covered",
		"uncovered":                      "Uncovered",
		"coverable":                      "Coverable",
		"total":                          "Total",
		"percentage":                     "Percentage",
		"allChanges":                     "All changes",
		"lineCoverageIncreaseOnly":       "Line coverage: Increase only",
		"lineCoverageDecreaseOnly":       "Line coverage: Decrease only",
		"branchCoverageIncreaseOnly":     "Branch coverage: Increase only",
		"branchCoverageDecreaseOnly":     "Branch coverage: Decrease only",
		"methodCoverageIncreaseOnly":     "Method coverage: Increase only",
		"methodCoverageDecreaseOnly":     "Method coverage: Decrease only",
		"fullMethodCoverageIncreaseOnly": "Full method coverage: Increase only",
		"fullMethodCoverageDecreaseOnly": "Full method coverage: Decrease only",
		// Ensure lowercase keys for Angular compatibility
		"methodCoverageProVersion": "This feature is only available for sponsors.",
		"coverageTypes":            "Coverage types",
		"history":                  "History",
	}

}
