package config

import (
	"fmt"
	"strings"
)

// OutputFormat describes one value of report_types. This table is the only
// place an output format is named and documented: the config validation, the
// generated docs and the website configurator read it. The code that writes
// the format is in cmd/main.go (reportWriters), under the same name.
type OutputFormat struct {
	Name   string
	Writes string // the file it writes into output_dir
	Doc    string
}

var OutputFormats = []OutputFormat{
	{
		Name:   "Html",
		Writes: "index.html",
		Doc:    "Interactive web report."},
	{
		Name:   "HtmlUnified",
		Writes: "index.html",
		Doc:    "The web report as one file."},
	{
		Name:   "TextSummary",
		Writes: "Summary.txt",
		Doc:    "Summary table. Also printed to the console."},
	{
		Name:   "Lcov",
		Writes: "lcov.info",
		Doc:    "Merged coverage in lcov format."},
	{
		Name:   "RawJson",
		Writes: "RawJson.json",
		Doc:    "All data as JSON."},
	{
		Name:   "Sarif",
		Writes: "nanovision.sarif.json",
		Doc:    "Problems as SARIF, for IDEs and code scanning."},
	{
		Name:   "Annotations",
		Writes: "Annotations.md",
		Doc:    "Problems as a Markdown list."},
}

func outputFormatNames() []string {
	names := make([]string, len(OutputFormats))
	for i, f := range OutputFormats {
		names[i] = f.Name
	}
	return names
}

// checkReportTypes rejects a report type the tool does not write.
func checkReportTypes(types []string) error {
	names := outputFormatNames()
	for _, t := range types {
		t = strings.TrimSpace(t)
		if t == "HtmlReview" {
			return fmt.Errorf("report type HtmlReview is now the Changes tab of Html and HtmlUnified; " +
				"for a report of only the changed files set diff.only_changed or pass -only-changed")
		}
		known := false
		for _, n := range names {
			known = known || n == t
		}
		if !known {
			return fmt.Errorf("unknown report type %q (expected one of: %s)", t, strings.Join(names, ", "))
		}
	}
	return nil
}
