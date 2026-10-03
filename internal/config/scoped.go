package config

import (
	"fmt"
	"io/fs"
	"maps"
	"os"
	"path"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/IgorBayerl/nanovision/internal/filtering"
)

// BandsFor returns the warning ranges of a file or folder: the root ranges,
// changed by every folder config on the way down to it. The nearest folder wins.
// The path is relative to the project root, "/" separated.
func (c *AppConfig) BandsFor(p string) StatusBands {
	var out StatusBands
	for _, o := range c.FolderBands {
		if p == o.Path || strings.HasPrefix(p, o.Path+"/") {
			if out == nil {
				out = maps.Clone(c.StatusBands)
			}
			maps.Copy(out, o.Bands)
		}
	}
	if out == nil {
		return c.StatusBands
	}
	return out
}

// resolveFolders turns the metric settings of the folder configs into
// FolderBands. It runs after resolveMetrics.
func (c *AppConfig) resolveFolders() error {
	c.FolderBands = nil
	for _, o := range c.Folders {
		bands := make(StatusBands)
		if err := scopedBands(o.Metrics.Files, FileMetricDefs, c.ActiveFileMetrics, o.File+": metrics.files", bands); err != nil {
			return err
		}
		if err := scopedBands(o.Metrics.Methods, MethodMetricDefs, c.ActiveMethodMetrics, o.File+": metrics.methods", bands); err != nil {
			return err
		}
		if len(bands) > 0 {
			source := o.File
			if rel, err := filepath.Rel(c.ProjectRoot, o.File); err == nil {
				source = filepath.ToSlash(rel)
			}
			c.FolderBands = append(c.FolderBands, FolderBand{Path: o.Path, Bands: bands, Source: source})
		}
	}
	// outer folders first, so BandsFor lets the nearest folder win
	sort.SliceStable(c.FolderBands, func(i, j int) bool {
		return len(c.FolderBands[i].Path) < len(c.FolderBands[j].Path)
	})
	return nil
}

// scopedBands reads the metric settings of a folder. A folder changes warning
// ranges only: which metrics exist is the same for the whole report.
func scopedBands(settings []MetricSetting, defs []MetricDef, active map[MetricKey]bool, where string, bands StatusBands) error {
	for _, s := range settings {
		def, ok := metricByName(defs, s.Name)
		if !ok {
			return fmt.Errorf("%s: unknown metric %q (expected one of: %s)", where, s.Name, strings.Join(metricNames(defs), ", "))
		}
		if !active[def.Key] {
			return fmt.Errorf("%s: metric %q is not enabled in the root config; a folder can change warning ranges only", where, s.Name)
		}
		if s.Warning == nil {
			return fmt.Errorf("%s: metric %q needs a warning range; a folder can change warning ranges only", where, s.Name)
		}
		if !def.Value && s.Warning.Max > 100 {
			return fmt.Errorf("%s: warning range %q of %s is above 100%%", where, s.Warning, s.Name)
		}
		bands[def.Key] = *s.Warning
	}
	return nil
}

// findNestedConfigs searches the project for config files in sub-folders and
// reads each as the config of its folder.
//
// ponytail: one walk of the whole project tree. On a very large workspace
// that takes seconds; nested_configs: false turns it off.
func findNestedConfigs(c *AppConfig) ([]FolderConfig, error) {
	var ignored []string
	ignored = append(ignored, c.FileFilters...)
	for _, pattern := range c.IgnoreFiles {
		ignored = append(ignored, "-"+pattern)
	}
	filter, err := filtering.NewDefaultFilter(ignored, true)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize file filter: %w", err)
	}
	outputDir, _ := filepath.Abs(c.OutputDir)

	var found []FolderConfig
	byFolder := make(map[string]string) // folder -> its config file
	err = filepath.WalkDir(c.ProjectRoot, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // an unreadable folder holds no config we can use
		}
		if d.IsDir() {
			name := d.Name()
			if p != c.ProjectRoot && (strings.HasPrefix(name, ".") || name == "node_modules" || p == outputDir) {
				return filepath.SkipDir
			}
			return nil
		}
		if !IsFileName(d.Name()) || filepath.Dir(p) == c.ProjectRoot {
			return nil
		}
		if other, twice := byFolder[filepath.Dir(p)]; twice {
			return sameFolderError(other, p)
		}
		byFolder[filepath.Dir(p)] = p
		rel, err := filepath.Rel(c.ProjectRoot, p)
		if err != nil {
			return nil
		}
		rel = filepath.ToSlash(rel)
		// a config inside ignored code is not ours, e.g. a vendored project
		if !filter.IsElementIncludedInReport(rel) {
			return nil
		}
		o, err := readNestedConfig(p)
		if err != nil {
			return err
		}
		o.Path = path.Dir(rel)
		found = append(found, o)
		return nil
	})
	return found, err
}

func readNestedConfig(file string) (FolderConfig, error) {
	o := FolderConfig{File: file}
	data, err := os.ReadFile(file)
	if err != nil {
		return o, fmt.Errorf("failed to read config file %s: %w", file, err)
	}
	keys, err := topLevelKeys(data)
	if err != nil {
		return o, fmt.Errorf("%s: %w", file, err)
	}
	scoped, root := yamlKeys(ScopedConfig{}), yamlKeys(AppConfig{})
	for _, k := range keys {
		switch {
		case slices.Contains(scoped, k):
		case slices.Contains(root, k) || removedKeys[k] != "":
			return o, fmt.Errorf("%s: %q can only be set in the root config file; a folder can set: %s. "+
				"If this folder is a separate project, add it to ignore_files of the root config",
				file, k, strings.Join(scoped, ", "))
		default:
			return o, fmt.Errorf("%s: unknown key %q (a folder can set: %s)", file, k, strings.Join(scoped, ", "))
		}
	}
	if err := decodeStrict(data, &o.ScopedConfig); err != nil {
		return o, fmt.Errorf("%s: %w", file, err)
	}
	return o, nil
}
