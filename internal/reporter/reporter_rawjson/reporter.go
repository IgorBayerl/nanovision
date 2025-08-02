package reporter_rawjson

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/IgorBayerl/AdlerCov/internal/model"
	"github.com/IgorBayerl/AdlerCov/internal/reporter"
)

type RawJsonReportBuilder struct {
	outputDir string
}

func NewRawJsonReportBuilder(outputDir string) reporter.ReportBuilder {
	return &RawJsonReportBuilder{
		outputDir: outputDir,
	}
}

func (b *RawJsonReportBuilder) ReportType() string {
	return "RawJson"
}

func (b *RawJsonReportBuilder) CreateReport(tree *model.SummaryTree) error {
	fileName := "RawJson.json"
	targetPath := filepath.Join(b.outputDir, fileName)

	// Marshal the tree structure into an indented JSON format for readability.
	jsonData, err := json.MarshalIndent(tree, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal summary tree to JSON: %w", err)
	}

	// Write the JSON data to the file.
	err = os.WriteFile(targetPath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write JSON summary report to '%s': %w", targetPath, err)
	}

	return nil
}
