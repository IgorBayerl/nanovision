// Package sarif emits coverage diagnostics as a SARIF 2.1.0 file.
//
// SARIF (Static Analysis Results Interchange Format) is the industry-standard
// static-analysis result format: JetBrains IDEs read it natively, VSCode reads
// it via the Microsoft SARIF Viewer extension, and GitHub code scanning ingests
// it directly. Emitting SARIF lets nanovision surface coverage risks as inline
// editor problems without any custom per-editor code.
package sarif

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/diagnostics"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/reporter"
	"github.com/IgorBayerl/nanovision/internal/status"
)

const (
	sarifVersion = "2.1.0"
	sarifSchema  = "https://json.schemastore.org/sarif-2.1.0.json"
	toolName     = "nanovision"
	toolURI      = "https://github.com/IgorBayerl/nanovision"
	outputFile   = "nanovision.sarif.json"
)

type builder struct {
	outputDir string
	config    *config.AppConfig
	registry  map[config.MetricKey]status.Evaluator
}

// NewSarifReportBuilder creates a SARIF report builder. The registry supplies
// human-readable rule metadata for the SARIF "rules" section.
func NewSarifReportBuilder(outputDir string, cfg *config.AppConfig, registry map[config.MetricKey]status.Evaluator) reporter.ReportBuilder {
	return &builder{outputDir: outputDir, config: cfg, registry: registry}
}

func (b *builder) ReportType() string { return "Sarif" }

func (b *builder) CreateReport(tree *model.SummaryTree) error {
	diags := diagnostics.Extract(tree, b.config, b.registry)
	doc := b.toSarif(diags)

	outputPath := filepath.Join(b.outputDir, outputFile)
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create SARIF file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	if err := enc.Encode(doc); err != nil {
		return fmt.Errorf("failed to encode SARIF: %w", err)
	}
	return nil
}

// ---- SARIF 2.1.0 document model (minimal subset we need) ----

type sarifLog struct {
	Schema  string     `json:"$schema"`
	Version string     `json:"version"`
	Runs    []sarifRun `json:"runs"`
}

type sarifRun struct {
	Tool    sarifTool     `json:"tool"`
	Results []sarifResult `json:"results"`
}

type sarifTool struct {
	Driver sarifDriver `json:"driver"`
}

type sarifDriver struct {
	Name           string      `json:"name"`
	InformationURI string      `json:"informationUri"`
	Rules          []sarifRule `json:"rules"`
}

type sarifRule struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	ShortDescription sarifText `json:"shortDescription"`
}

type sarifResult struct {
	RuleID    string          `json:"ruleId"`
	Level     string          `json:"level"` // error | warning | note | none
	Message   sarifText       `json:"message"`
	Locations []sarifLocation `json:"locations"`
}

type sarifText struct {
	Text string `json:"text"`
}

type sarifLocation struct {
	PhysicalLocation sarifPhysicalLocation `json:"physicalLocation"`
}

type sarifPhysicalLocation struct {
	ArtifactLocation sarifArtifactLocation `json:"artifactLocation"`
	Region           sarifRegion           `json:"region"`
}

type sarifArtifactLocation struct {
	URI string `json:"uri"`
}

type sarifRegion struct {
	StartLine int `json:"startLine"`
	EndLine   int `json:"endLine"`
}

func (b *builder) toSarif(diags []diagnostics.Diagnostic) sarifLog {
	results := make([]sarifResult, 0, len(diags))
	rules := make([]sarifRule, 0)
	seenRule := make(map[string]bool)

	for _, d := range diags {
		if !seenRule[d.RuleID] {
			seenRule[d.RuleID] = true
			rules = append(rules, sarifRule{
				ID:               d.RuleID,
				Name:             d.RuleName,
				ShortDescription: sarifText{Text: b.ruleDescription(d)},
			})
		}
		results = append(results, sarifResult{
			RuleID:  d.RuleID,
			Level:   sarifLevel(d.Severity),
			Message: sarifText{Text: d.Message},
			Locations: []sarifLocation{{
				PhysicalLocation: sarifPhysicalLocation{
					ArtifactLocation: sarifArtifactLocation{URI: d.File},
					Region:           sarifRegion{StartLine: d.StartLine, EndLine: d.EndLine},
				},
			}},
		})
	}

	return sarifLog{
		Schema:  sarifSchema,
		Version: sarifVersion,
		Runs: []sarifRun{{
			Tool: sarifTool{Driver: sarifDriver{
				Name:           toolName,
				InformationURI: toolURI,
				Rules:          rules,
			}},
			Results: results,
		}},
	}
}

func (b *builder) ruleDescription(d diagnostics.Diagnostic) string {
	if ev, ok := b.registry[config.MetricKey(d.RuleID)]; ok {
		return ev.Description()
	}
	return d.RuleName
}

// sarifLevel maps nanovision severities onto SARIF's fixed level vocabulary.
func sarifLevel(s diagnostics.Severity) string {
	switch s {
	case diagnostics.SeverityError:
		return "error"
	case diagnostics.SeverityWarning:
		return "warning"
	default:
		return "note"
	}
}
