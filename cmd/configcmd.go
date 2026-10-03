package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/IgorBayerl/nanovision/internal/config"
)

// The config reference of the website. A test fails when the file is stale.
//go:generate go run . config schema -o ../docs/src/generated/config.json

const configUsage = `Usage: nanovision config <command> [flags]

Commands:
  docs           print every config key and every metric, with its documentation
  check          read the root config and the configs of its sub-folders, and report errors
  show [folder]  print the settings that apply to a folder, and which file sets each
  schema         print the config reference as JSON (the website reads it)

Flags:
  -config <file> root config file (default: the one in this folder)
  -o <file>      schema: write to a file instead of standard output
`

// runConfigCommand is `nanovision config ...`. Every command only reads.
func runConfigCommand(args []string) int {
	if len(args) == 0 || args[0] == "-h" || args[0] == "-help" || args[0] == "--help" {
		fmt.Fprint(os.Stderr, configUsage)
		return 2
	}
	cmd := args[0]
	fs := flag.NewFlagSet("config "+cmd, flag.ExitOnError)
	fs.Usage = func() { fmt.Fprint(os.Stderr, configUsage) }
	configPath := fs.String("config", "", "")
	out := fs.String("o", "", "")
	fs.Parse(args[1:])

	switch cmd {
	case "docs":
		printConfigDocs(os.Stdout, config.Schema())
	case "schema":
		data, _ := json.MarshalIndent(config.Schema(), "", "  ")
		data = append(data, '\n')
		if *out == "" {
			os.Stdout.Write(data)
			return 0
		}
		if err := os.MkdirAll(filepath.Dir(*out), 0o755); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		if err := os.WriteFile(*out, data, 0o644); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
	case "check", "show":
		if *configPath == "" {
			if found, err := config.FindFile("."); err == nil && found == "" {
				fmt.Fprintf(os.Stderr, "No config file (%s) in this folder; pass -config <file>.\n", config.FileNames[0])
				return 1
			}
		}
		cfg, err := config.Load(*configPath, config.RawConfigInput{})
		if err != nil {
			fmt.Fprintln(os.Stderr, "Configuration error:", err)
			return 1
		}
		if cmd == "check" {
			printConfigCheck(os.Stdout, cfg)
		} else {
			printConfigShow(os.Stdout, cfg, fs.Arg(0))
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown config command %q\n\n%s", cmd, configUsage)
		return 2
	}
	return 0
}

func printMetrics(w io.Writer) {
	schema := config.Schema()
	for _, scope := range []struct{ key, title string }{
		{"files", "File and folder metrics (yaml: metrics.files)"},
		{"methods", "Method metrics (yaml: metrics.methods)"},
	} {
		fmt.Fprintf(w, "\n%s\n", scope.title)
		for _, m := range schema.Metrics[scope.key] {
			fmt.Fprintf(w, "  %-32s %s\n", m.Name, m.Doc)
		}
	}
}

func printConfigDocs(w io.Writer, schema config.SchemaDoc) {
	fmt.Fprintln(w, "nanovision.yaml reference")
	fmt.Fprintln(w, "\nA key marked [root] is valid only in the root config file. The other keys")
	fmt.Fprintln(w, "can also be set for one folder: put a nanovision.yaml inside the folder.")
	fmt.Fprintln(w, "It applies to that folder and everything below it.")
	fmt.Fprintln(w)
	var print func(fields []config.Field, rootOnly bool)
	print = func(fields []config.Field, rootOnly bool) {
		for _, f := range fields {
			line := f.Key + "  (" + f.Type + ")"
			if len(f.Values) > 0 {
				line += "  " + strings.Join(f.Values, " | ")
			}
			if f.Default != "" {
				line += "  default: " + f.Default
			}
			if f.RootOnly || rootOnly {
				line += "  [root]"
			}
			fmt.Fprintf(w, "%s\n    %s\n", line, f.Doc)
			print(f.Fields, f.RootOnly || rootOnly)
		}
	}
	print(schema.Fields, false)

	fmt.Fprintln(w, "\nOutput formats (yaml: report_types)")
	for _, f := range schema.OutputFormats {
		fmt.Fprintf(w, "  %-32s %s Writes %s.\n", f.Name, f.Doc, f.Writes)
	}
	printMetrics(w)
}

func printConfigCheck(w io.Writer, cfg *config.AppConfig) {
	fmt.Fprintf(w, "OK  %s\n", cfg.ConfigFile)
	for _, o := range cfg.Folders {
		fmt.Fprintf(w, "OK  %s  (folder %s)\n", relTo(cfg.ProjectRoot, o.File), o.Path)
	}
	fmt.Fprintf(w, "\n%d reports, %d file metrics, %d method metrics, %d folders with their own settings.\n",
		len(cfg.InputPairs), len(cfg.FileMetrics), len(cfg.MethodMetrics), len(cfg.Folders))
}

// printConfigShow prints what applies to one folder, and where it comes from.
func printConfigShow(w io.Writer, cfg *config.AppConfig, folder string) {
	folder = strings.Trim(filepath.ToSlash(filepath.Clean(folder)), "/")
	if folder == "." {
		folder = ""
	}
	name := folder
	if name == "" {
		name = "the project root"
	}
	fmt.Fprintf(w, "Settings that apply to %s\n", name)

	bands := cfg.BandsFor(folder)
	origin := func(key config.MetricKey) string {
		from := cfg.ConfigFile
		for _, o := range cfg.FolderBands {
			if _, sets := o.Bands[key]; sets && (folder == o.Path || strings.HasPrefix(folder, o.Path+"/")) {
				from = o.Source
			}
		}
		return from
	}
	for _, scope := range []struct {
		title string
		keys  []config.MetricKey
	}{{"File metrics", cfg.FileMetrics}, {"Method metrics", cfg.MethodMetrics}} {
		fmt.Fprintf(w, "\n%s\n", scope.title)
		for _, key := range scope.keys {
			def, _ := config.Metric(key)
			if band, ok := bands[key]; ok {
				fmt.Fprintf(w, "  %-32s warning %-10s set in %s\n", def.Name, band, origin(key))
			} else {
				fmt.Fprintf(w, "  %-32s no warning range\n", def.Name)
			}
		}
	}

	fmt.Fprintln(w, "\nIgnored files")
	for _, p := range cfg.IgnoreFiles {
		fmt.Fprintf(w, "  %-32s set in %s\n", p, cfg.ConfigFile)
	}
	for _, o := range cfg.Folders {
		if folder == o.Path || strings.HasPrefix(folder, o.Path+"/") {
			for _, p := range o.IgnoreFiles {
				fmt.Fprintf(w, "  %-32s set in %s\n", o.Path+"/"+p, relTo(cfg.ProjectRoot, o.File))
			}
		}
	}
}

func relTo(root, file string) string {
	if rel, err := filepath.Rel(root, file); err == nil {
		return filepath.ToSlash(rel)
	}
	return file
}
