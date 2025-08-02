package lcov

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/IgorBayerl/AdlerCov/internal/model"
	"github.com/IgorBayerl/AdlerCov/internal/reporter"
)

type LcovReportBuilder struct {
	outputDir string
}

func NewLcovReportBuilder(outputDir string) reporter.ReportBuilder {
	return &LcovReportBuilder{
		outputDir: outputDir,
	}
}

func (b *LcovReportBuilder) ReportType() string {
	return "Lcov"
}

func (b *LcovReportBuilder) CreateReport(tree *model.SummaryTree) error {
	fileName := "lcov.info"
	targetPath := filepath.Join(b.outputDir, fileName)

	file, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("failed to create lcov report file '%s': %w", targetPath, err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// Collect all FileNodes by walking the tree.
	files := collectFileNodes(tree.Root)
	sort.Slice(files, func(i, j int) bool {
		return files[i].Path < files[j].Path
	})

	for _, fileNode := range files {
		if err := writeLcovFileSection(writer, fileNode); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Could not process file '%s' for LCOV report: %v\n", fileNode.Path, err)
		}
	}

	return nil
}

func writeLcovFileSection(writer *bufio.Writer, file *model.FileNode) error {
	if _, err := writer.WriteString(fmt.Sprintf("SF:%s\n", file.Path)); err != nil {
		return err
	}

	// FN/FNDA: Function data.
	// Sort methods by their starting line number for consistent output.
	sort.SliceStable(file.Methods, func(i, j int) bool {
		return file.Methods[i].StartLine < file.Methods[j].StartLine
	})
	fnf := len(file.Methods)
	fnh := 0
	for _, method := range file.Methods {
		if _, err := writer.WriteString(fmt.Sprintf("FN:%d,%s\n", method.StartLine, method.Name)); err != nil {
			return err
		}
		hitCount := 0
		if method.LineCoverage > 0 {
			hitCount = 1 // LCOV treats any coverage as 1 hit.
			fnh++
		}
		if _, err := writer.WriteString(fmt.Sprintf("FNDA:%d,%s\n", hitCount, method.Name)); err != nil {
			return err
		}
	}
	if _, err := writer.WriteString(fmt.Sprintf("FNF:%d\n", fnf)); err != nil {
		return err
	}
	if _, err := writer.WriteString(fmt.Sprintf("FNH:%d\n", fnh)); err != nil {
		return err
	}

	// DA/LF/LH: Line data.
	lf := file.Metrics.LinesValid
	lh := file.Metrics.LinesCovered
	for lineNum, lineMetrics := range file.Lines {
		if lineMetrics.Hits >= 0 {
			if _, err := writer.WriteString(fmt.Sprintf("DA:%d,%d\n", lineNum, lineMetrics.Hits)); err != nil {
				return err
			}
		}
	}
	if _, err := writer.WriteString(fmt.Sprintf("LF:%d\n", lf)); err != nil {
		return err
	}
	if _, err := writer.WriteString(fmt.Sprintf("LH:%d\n", lh)); err != nil {
		return err
	}

	// BRDA/BRF/BRH: Branch data.
	brf := file.Metrics.BranchesValid
	brh := file.Metrics.BranchesCovered
	if brf > 0 {
		for lineNum, lineMetrics := range file.Lines {
			if lineMetrics.TotalBranches > 0 {
				// This is a simplification. LCOV needs per-branch data, which our
				// new model currently aggregates. For now, we report the line-level aggregate.
				// Format: BRDA:<line>,<block>,<branch>,<hits>
				// We can represent covered branches as hit and uncovered as not hit.
				for i := 0; i < lineMetrics.TotalBranches; i++ {
					hits := "-"
					if i < lineMetrics.CoveredBranches {
						hits = "1"
					}
					if _, err := writer.WriteString(fmt.Sprintf("BRDA:%d,%d,%d,%s\n", lineNum, 0, i, hits)); err != nil {
						return err
					}
				}
			}
		}
		if _, err := writer.WriteString(fmt.Sprintf("BRF:%d\n", brf)); err != nil {
			return err
		}
		if _, err := writer.WriteString(fmt.Sprintf("BRH:%d\n", brh)); err != nil {
			return err
		}
	}

	if _, err := writer.WriteString("end_of_record\n"); err != nil {
		return err
	}

	return nil
}

// collectFileNodes is a new helper to walk the tree and gather all file nodes.
func collectFileNodes(dir *model.DirNode) []*model.FileNode {
	var files []*model.FileNode
	for _, file := range dir.Files {
		files = append(files, file)
	}
	for _, subDir := range dir.Subdirs {
		files = append(files, collectFileNodes(subDir)...)
	}
	return files
}
