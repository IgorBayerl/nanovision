package config

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"go.yaml.in/yaml/v3"
)

func defaultCLIInput() RawConfigInput {
	return RawConfigInput{
		ReportTypes: "TextSummary,Html",
		OutputDir:   "coverage-report",
		LogFormat:   "text",
		Verbosity:   "Info",
	}
}

// writeProject writes files (path -> content) into a new folder and returns
// the path of its root nanovision.yaml.
func writeProject(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for name, content := range files {
		p := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return filepath.Join(root, "nanovision.yaml")
}

const rootYAML = `
reports:
  - coverage.out
  - path: build/cobertura.xml
    source: services/backend
    name: backend
ignore_files:
  - "vendor/**"
metrics:
  files:
    - name: methods_hit
    - name: statement_coverage
      warning: "60..75"
    - name: max_complexity
      warning: "10..15"
  methods:
    - name: complexity
      warning: "10..150"
`

func TestLoad_ReportsMetricsAndFolders(t *testing.T) {
	configPath := writeProject(t, map[string]string{
		"nanovision.yaml": rootYAML,
		"services/backend/nanovision.yaml": `
metrics:
  files:
    - name: statement_coverage
      warning: "80..90"
`,
		// the nearest folder config wins
		"services/backend/api/nanovision.yaml": `
reports:
  - path: api.out
    name: api
ignore_files:
  - "gen/**"
metrics:
  files:
    - name: statement_coverage
      warning: "90..95"
  methods:
    - name: complexity
      warning: "5..8"
`,
		// inside ignored code: a vendored project, not a folder config
		"vendor/lib/nanovision.yaml": "title: other project\n",
	})
	cfg, err := Load(configPath, defaultCLIInput())
	if err != nil {
		t.Fatal(err)
	}

	if want := []MetricKey{MethodsHit, StatementCoverage, MaxCyclomaticComplexity}; !reflect.DeepEqual(cfg.FileMetrics, want) {
		t.Errorf("file metrics keep the listed order: got %v", cfg.FileMetrics)
	}
	if want := []MetricKey{CyclomaticComplexity}; !reflect.DeepEqual(cfg.MethodMetrics, want) {
		t.Errorf("method metrics: got %v", cfg.MethodMetrics)
	}
	if !cfg.ActiveFileMetrics[MethodsHit] || cfg.ActiveFileMetrics[LineCoverage] {
		t.Errorf("active file metrics: got %v", cfg.ActiveFileMetrics)
	}

	for path, want := range map[string]Band{
		"cmd/main.go":                      {60, 75},
		"services/backend":                 {80, 90},
		"services/backend/db/store.go":     {80, 90},
		"services/backend/api/handler.go":  {90, 95},
		"services/backendish/unrelated.go": {60, 75},
	} {
		if got := cfg.BandsFor(path)[StatementCoverage]; got != want {
			t.Errorf("statement band of %s: got %v, want %v", path, got, want)
		}
	}
	api := cfg.BandsFor("services/backend/api/handler.go")
	if api[CyclomaticComplexity] != (Band{5, 8}) || api[MaxCyclomaticComplexity] != (Band{10, 15}) {
		t.Errorf("a folder changes only the ranges it names: got %v", api)
	}

	// the paths of a folder config become relative to the project root
	apiDir := "services/backend/api"
	want := []ReportInputPair{
		{ReportPattern: "coverage.out", SourceDir: "."},
		{ReportPattern: "build/cobertura.xml", SourceDir: "services/backend", Name: "backend"},
		{ReportPattern: apiDir + "/api.out", SourceDir: apiDir, Name: "api"},
	}
	if !reflect.DeepEqual(cfg.InputPairs, want) {
		t.Errorf("input pairs:\n got %+v\nwant %+v", cfg.InputPairs, want)
	}

	for path, included := range map[string]bool{
		"services/backend/api/handler.go": true,
		"services/backend/api/gen/pb.go":  false,
		"services/backend/gen/pb.go":      true,
		"vendor/lib/lib.go":               false,
	} {
		if got := cfg.FileFilterInstance.IsElementIncludedInReport(path); got != included {
			t.Errorf("%s included: got %v, want %v", path, got, included)
		}
	}
}

func TestLoad_Errors(t *testing.T) {
	reports := "reports: [coverage.out]\n"
	for name, tc := range map[string]struct {
		files map[string]string
		want  string
	}{
		"removed key": {
			map[string]string{"nanovision.yaml": reports + "status_bands:\n  line_coverage: \"60..75\"\n"},
			`"status_bands" is no longer a config key`,
		},
		"unknown key": {
			map[string]string{"nanovision.yaml": reports + "titel: x\n"},
			`unknown key "titel"`,
		},
		"unknown metric": {
			map[string]string{"nanovision.yaml": reports + "metrics:\n  methods:\n    - name: methods_hit\n"},
			`metrics.methods: unknown metric "methods_hit" (expected one of: statement_coverage,`,
		},
		"metric listed twice": {
			map[string]string{"nanovision.yaml": reports + "metrics:\n  files:\n    - name: methods_hit\n    - name: methods_hit\n"},
			"listed twice",
		},
		"percentage above 100": {
			map[string]string{"nanovision.yaml": reports + "metrics:\n  files:\n    - name: methods_hit\n      warning: \"60..175\"\n"},
			"above 100%",
		},
		"bad range": {
			map[string]string{"nanovision.yaml": reports + "metrics:\n  files:\n    - name: methods_hit\n      warning: sixty\n"},
			`invalid range "sixty"`,
		},
		"unknown report key": {
			map[string]string{"nanovision.yaml": "reports:\n  - path: a.out\n    sorce: .\n"},
			`unknown key "sorce" in a report`,
		},
		"no reports": {
			map[string]string{"nanovision.yaml": "title: x\n"},
			"at least one report",
		},
		"root-only key in a folder": {
			map[string]string{"nanovision.yaml": reports, "sub/nanovision.yaml": "output_dir: out\n"},
			`"output_dir" can only be set in the root config file`,
		},
		"folder enables a metric": {
			map[string]string{
				"nanovision.yaml":     reports + "metrics:\n  files:\n    - name: methods_hit\n",
				"sub/nanovision.yaml": "metrics:\n  files:\n    - name: line_coverage\n      warning: \"1..2\"\n",
			},
			`metric "line_coverage" is not enabled in the root config`,
		},
		"overrides key of older versions": {
			map[string]string{"nanovision.yaml": reports + "overrides:\n  - path: cmd\n"},
			"put a nanovision.yaml inside the folder",
		},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := Load(writeProject(t, tc.files), defaultCLIInput())
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Errorf("got error %v, want it to contain %q", err, tc.want)
			}
		})
	}
}

// The root and the sub-folders accept the same file names.
func TestLoad_FileNames(t *testing.T) {
	for _, name := range FileNames {
		t.Run(name, func(t *testing.T) {
			configPath := writeProject(t, map[string]string{
				"nanovision.yaml": "reports: [coverage.out]\nmetrics:\n  files:\n    - name: methods_hit\n",
				"sub/" + name:     "metrics:\n  files:\n    - name: methods_hit\n      warning: \"1..2\"\n",
			})
			cfg, err := Load(configPath, defaultCLIInput())
			if err != nil {
				t.Fatal(err)
			}
			if got := cfg.BandsFor("sub/a.go")[MethodsHit]; got != (Band{1, 2}) {
				t.Errorf("sub/%s was not read: got %v", name, got)
			}

			// the same name is found as the root config of a folder
			root := filepath.Dir(configPath)
			found, err := FindFile(filepath.Join(root, "sub"))
			if err != nil || filepath.Base(found) != name {
				t.Errorf("FindFile: got %q, %v", found, err)
			}
		})
	}

	configPath := writeProject(t, map[string]string{
		"nanovision.yaml":      "reports: [coverage.out]\n",
		"sub/nanovision.yaml":  "ignore_files: [a]\n",
		"sub/.nanovision.yaml": "ignore_files: [b]\n",
	})
	if _, err := Load(configPath, defaultCLIInput()); err == nil || !strings.Contains(err.Error(), "same folder") {
		t.Errorf("two config files in one folder must be an error, got %v", err)
	}
	if _, err := FindFile(filepath.Join(filepath.Dir(configPath), "sub")); err == nil {
		t.Error("FindFile must reject two config files in one folder")
	}
}

func TestLoad_NestedConfigsOff(t *testing.T) {
	configPath := writeProject(t, map[string]string{
		"nanovision.yaml":     "reports: [coverage.out]\nnested_configs: false\n",
		"sub/nanovision.yaml": "output_dir: out\n",
	})
	if _, err := Load(configPath, defaultCLIInput()); err != nil {
		t.Errorf("with nested_configs off the file in sub/ is not read: %v", err)
	}
}

// loadFlags is Load without a config file.
func loadFlags(cfg *AppConfig, cli RawConfigInput) error {
	if err := cfg.mergeCliOverrides(cli); err != nil {
		return err
	}
	if err := cfg.validate(); err != nil {
		return err
	}
	return cfg.computeDerivedFields()
}

func TestConfig_DefaultsAndFlags(t *testing.T) {
	cli := defaultCLIInput()
	cli.Reports = []string{"a.xml;b.out"}
	cfg := GetDefaultConfig()
	if err := loadFlags(cfg, cli); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(cfg.FileMetrics, DefaultFileMetrics) || !reflect.DeepEqual(cfg.MethodMetrics, DefaultMethodMetrics) {
		t.Errorf("without a metrics section every metric is on: got %v and %v", cfg.FileMetrics, cfg.MethodMetrics)
	}
	if want := []ReportInputPair{{ReportPattern: "a.xml", SourceDir: "."}, {ReportPattern: "b.out", SourceDir: "."}}; !reflect.DeepEqual(cfg.InputPairs, want) {
		t.Errorf("the source folder defaults to '.': got %+v", cfg.InputPairs)
	}

	cli.SourceDirs = "src"
	cli.FileMetrics = []string{"line_coverage", "statement_coverage"}
	cli.MethodMetrics = []string{"complexity"}
	cli.StatusBands = []string{"line_coverage=40..80", "methods.complexity=10..15"}
	cfg = GetDefaultConfig()
	if err := loadFlags(cfg, cli); err != nil {
		t.Fatal(err)
	}
	if cfg.InputPairs[0].SourceDir != "src" || cfg.InputPairs[1].SourceDir != "src" {
		t.Errorf("one -sourcedirs value serves every report: got %+v", cfg.InputPairs)
	}
	cli.StatusBands = append(cli.StatusBands, "methods_hit=50..60")
	cfg = GetDefaultConfig()
	if err := loadFlags(cfg, cli); err != nil {
		t.Fatal(err)
	}
	if want := []MetricKey{LineCoverage, StatementCoverage, MethodsHit}; !reflect.DeepEqual(cfg.FileMetrics, want) {
		t.Errorf("-file-metrics, plus the metric a -threshold names: got %v", cfg.FileMetrics)
	}
	if want := (StatusBands{LineCoverage: {40, 80}, CyclomaticComplexity: {10, 15}, MethodsHit: {50, 60}}); !reflect.DeepEqual(cfg.StatusBands, want) {
		t.Errorf("-threshold: got %v", cfg.StatusBands)
	}

	for name, mutate := range map[string]func(*RawConfigInput){
		"source count":      func(c *RawConfigInput) { c.SourceDirs = "a;b;c" },
		"unknown metric":    func(c *RawConfigInput) { c.FileMetrics = []string{"lines_coverage"} },
		"unknown threshold": func(c *RawConfigInput) { c.StatusBands = []string{"complexity=1..2"} },
	} {
		bad := cli
		mutate(&bad)
		if err := loadFlags(GetDefaultConfig(), bad); err == nil {
			t.Errorf("%s: expected an error", name)
		}
	}
}

func TestSchema_EveryKeyIsDocumented(t *testing.T) {
	schema := Schema()
	for _, f := range Flatten(schema.Fields) {
		if f.Doc == "" {
			t.Errorf("config key %s has no doc tag", f.Key)
		}
	}
	for scope, docs := range schema.Metrics {
		for _, m := range docs {
			if m.Doc == "" || m.Label == "" {
				t.Errorf("metric %s.%s needs a Label and a Doc", scope, m.Name)
			}
		}
	}

	byKey := map[string]Field{}
	for _, f := range Flatten(schema.Fields) {
		byKey[f.Key] = f
	}
	if f := byKey["output_dir"]; f.Default != "coverage-report" || !f.RootOnly {
		t.Errorf("output_dir: got %+v", f)
	}
	if f := byKey["ignore_files"]; f.RootOnly || f.Type != "list of strings" {
		t.Errorf("ignore_files is a folder setting: got %+v", f)
	}
	if f := byKey["metrics.files[].warning"]; f.Type != "range" {
		t.Errorf("metrics.files[].warning: got %+v", f)
	}
	if f := byKey["vcs.type"]; !reflect.DeepEqual(f.Values, []string{"git", "perforce"}) || f.Default != "" {
		t.Errorf("vcs.type: got %+v", f)
	}
}

func TestConfig_History(t *testing.T) {
	cfg := GetDefaultConfig()
	if err := yaml.Unmarshal([]byte("history:\n  store: .nanovision\nvcs:\n  type: Git\n"), cfg); err != nil {
		t.Fatal(err)
	}
	cli := defaultCLIInput()
	cli.Reports, cli.SourceDirs = []string{"report.xml"}, "."
	cli.RunKind, cli.Revision, cli.Profile = "Submit", " //game/main@118432 ", "unit-win64"

	if err := loadFlags(cfg, cli); err != nil {
		t.Fatal(err)
	}

	if cfg.History.Store != ".nanovision" || cfg.History.IsServer() {
		t.Errorf("store: got %+v", cfg.History)
	}
	if cfg.History.MaxDistance != 50 || cfg.History.KeepLocal != 20 {
		t.Errorf("YAML without these keys must keep the defaults, got %+v", cfg.History)
	}
	if cfg.History.Profile != "unit-win64" {
		t.Errorf("the flag wins over the default profile, got %q", cfg.History.Profile)
	}
	if cfg.VCS.Type != "git" {
		t.Errorf("vcs type is case-insensitive, got %q", cfg.VCS.Type)
	}
	if cfg.Run.Kind != "submit" || cfg.Run.Revision != "//game/main@118432" {
		t.Errorf("run flags: got %+v", cfg.Run)
	}
	if !(HistoryConfig{Store: "https://nv.studio.local:7070"}).IsServer() {
		t.Error("an http URL is a team server")
	}
}

func TestConfig_HistoryValidation(t *testing.T) {
	for name, mutate := range map[string]func(*AppConfig, *RawConfigInput){
		"bad profile":  func(c *AppConfig, _ *RawConfigInput) { c.History.Profile = "unit win64" },
		"bad vcs":      func(c *AppConfig, _ *RawConfigInput) { c.VCS.Type = "svn" },
		"bad run kind": func(_ *AppConfig, cli *RawConfigInput) { cli.RunKind = "nightly" },
		"bad distance": func(c *AppConfig, _ *RawConfigInput) { c.History.MaxDistance = -1 },
	} {
		t.Run(name, func(t *testing.T) {
			cfg := GetDefaultConfig()
			cli := defaultCLIInput()
			cli.Reports, cli.SourceDirs = []string{"report.xml"}, "."
			mutate(cfg, &cli)
			if err := loadFlags(cfg, cli); err == nil {
				t.Error("expected a validation error")
			}
		})
	}
}

func TestConfig_ReportAndMetricFlags(t *testing.T) {
	cli := defaultCLIInput()
	cli.Reports = []string{"unit.out,source=src,name=unit", "build/**/*.gcov"}
	cli.FileMetrics = []string{"methods_hit", "statement_coverage=60..75"}
	cli.MethodMetrics = []string{"complexity=10..15"}
	cfg := GetDefaultConfig()
	if err := loadFlags(cfg, cli); err != nil {
		t.Fatal(err)
	}
	want := []ReportInputPair{
		{ReportPattern: "unit.out", SourceDir: "src", Name: "unit"},
		{ReportPattern: "build/**/*.gcov", SourceDir: "."},
	}
	if !reflect.DeepEqual(cfg.InputPairs, want) {
		t.Errorf("-report: got %+v", cfg.InputPairs)
	}
	if want := []MetricKey{MethodsHit, StatementCoverage}; !reflect.DeepEqual(cfg.FileMetrics, want) {
		t.Errorf("-file-metric keeps the order of the flags: got %v", cfg.FileMetrics)
	}
	if want := (StatusBands{StatementCoverage: {60, 75}, CyclomaticComplexity: {10, 15}}); !reflect.DeepEqual(cfg.StatusBands, want) {
		t.Errorf("a metric flag carries its warning range: got %v", cfg.StatusBands)
	}

	cli.Reports = []string{"unit.out,sorce=src"}
	if err := loadFlags(GetDefaultConfig(), cli); err == nil {
		t.Error("an unknown part of -report is an error")
	}
}

// -set reaches every plain config key, so no key needs the config file.
func TestConfig_SetFlag(t *testing.T) {
	cli := defaultCLIInput()
	cli.Reports = []string{"a.out"}
	cli.Set = []string{
		"title=123",
		"history.keep_local=5",
		"review.gate.patch_statement_coverage=80",
		"diff.only_changed=true",
		"nested_configs=false",
		"ignore_files=vendor/**",
		"ignore_files=**/*_test.go",
		"vcs.base_branch=origin/dev",
	}
	cfg := GetDefaultConfig()
	if err := loadFlags(cfg, cli); err != nil {
		t.Fatal(err)
	}
	if cfg.Title != "123" || cfg.History.KeepLocal != 5 || !cfg.Diff.OnlyChanged || cfg.NestedConfigs || cfg.VCS.BaseBranch != "origin/dev" {
		t.Errorf("got title %q, keep_local %d, only_changed %v, nested %v, base branch %q",
			cfg.Title, cfg.History.KeepLocal, cfg.Diff.OnlyChanged, cfg.NestedConfigs, cfg.VCS.BaseBranch)
	}
	if g := cfg.Review.Gate.PatchStatementCoverage; g == nil || *g != 80 {
		t.Errorf("review.gate.patch_statement_coverage: got %v", g)
	}
	if want := []string{"vendor/**", "**/*_test.go"}; !reflect.DeepEqual(cfg.IgnoreFiles, want) {
		t.Errorf("a list key takes one -set for each item: got %v", cfg.IgnoreFiles)
	}

	for _, bad := range []string{"titel=x", "overrides=x", "review=x", "history.keep_local=many", "novalue"} {
		cli.Set = []string{bad}
		if err := loadFlags(GetDefaultConfig(), cli); err == nil {
			t.Errorf("-set %s: expected an error", bad)
		}
	}

	// every key the form of the website can edit is settable
	for _, f := range Flatten(Schema().Fields) {
		if strings.Contains(f.Key, "[]") || f.Type == "object" || f.Type == "list" {
			continue
		}
		if _, known := map[string]bool{"string": true, "boolean": true, "number": true, "list of strings": true}[f.Type]; !known {
			t.Errorf("config key %s has type %q, which -set cannot set", f.Key, f.Type)
		}
	}
}
