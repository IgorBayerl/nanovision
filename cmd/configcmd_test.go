package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/IgorBayerl/nanovision/internal/config"
)

// The website reads docs/src/generated/config.json. It is generated from the
// config struct and the metric tables, so it must never be edited by hand or
// left behind after a config change.
func TestGeneratedConfigReferenceIsCurrent(t *testing.T) {
	want, err := json.MarshalIndent(config.Schema(), "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile("../docs/src/generated/config.json")
	if err != nil {
		t.Fatalf("%v\nrun: go generate ./...", err)
	}
	// git may check the file out with CRLF line ends
	got = bytes.ReplaceAll(got, []byte("\r\n"), []byte("\n"))
	if !bytes.Equal(bytes.TrimSpace(got), want) {
		t.Fatal("docs/src/generated/config.json is stale; run: go generate ./...")
	}
}

func TestConfigDocsListEveryKeyAndMetric(t *testing.T) {
	var out bytes.Buffer
	printConfigDocs(&out, config.Schema())
	text := out.String()
	for _, want := range []string{"output_dir", "[root]", "metrics.files[].warning", "review.gate.patch_statement_coverage", "statement_coverage", "crap_score"} {
		if !strings.Contains(text, want) {
			t.Errorf("config docs lack %q", want)
		}
	}
}

// An output format is one entry in config.OutputFormats plus one writer.
func TestEveryOutputFormatHasAWriter(t *testing.T) {
	for _, f := range config.OutputFormats {
		if _, ok := reportWriters[f.Name]; !ok {
			t.Errorf("output format %s has no writer in reportWriters", f.Name)
		}
		if f.Doc == "" || f.Writes == "" {
			t.Errorf("output format %s needs a Doc and the file it Writes", f.Name)
		}
	}
	if len(reportWriters) != len(config.OutputFormats) {
		t.Error("a writer in reportWriters has no entry in config.OutputFormats")
	}
}

// The website builds command lines from the flag tags of the config struct.
func TestFlagTagsNameRealFlags(t *testing.T) {
	parseAndBindFlags()
	tagged := 0
	for _, f := range config.Flatten(config.Schema().Fields) {
		if f.Flag == "" {
			continue
		}
		tagged++
		if flag.Lookup(f.Flag) == nil {
			t.Errorf("config key %s names the flag -%s, which does not exist", f.Key, f.Flag)
		}
	}
	if tagged == 0 {
		t.Error("no config key has a flag tag")
	}
}
