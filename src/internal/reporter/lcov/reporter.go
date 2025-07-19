package lcov

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/model"
	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/reporter"
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

func (b *LcovReportBuilder) CreateReport(summary *model.SummaryResult) error {
	fileName := "lcov.info"
	targetPath := filepath.Join(b.outputDir, fileName)

	file, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("failed to create lcov report file '%s': %w", targetPath, err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	files := getAllFiles(summary.Assemblies)
	sort.Slice(files, func(i, j int) bool {
		return files[i].Path < files[j].Path
	})

	for _, fileAnalysis := range files {
		if err := writeLcovFileSection(writer, fileAnalysis); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Could not process file '%s' for LCOV report: %v\n", fileAnalysis.Path, err)
		}
	}

	return nil
}

func writeLcovFileSection(writer *bufio.Writer, file *model.CodeFile) error {
	// TN: Test Name (optional, often omitted in summary files)
	if _, err := writer.WriteString(fmt.Sprintf("SF:%s\n", file.Path)); err != nil {
		return err
	}

	// FN/FNDA: Function and Function execution data
	// Sort code elements by line number for deterministic output.
	sort.SliceStable(file.CodeElements, func(i, j int) bool {
		return file.CodeElements[i].FirstLine < file.CodeElements[j].FirstLine
	})
	fnf := len(file.CodeElements) // Total functions found
	fnh := 0                      // Total functions hit (covered)

	for _, codeElement := range file.CodeElements {
		if _, err := writer.WriteString(fmt.Sprintf("FN:%d,%s\n", codeElement.FirstLine, codeElement.FullName)); err != nil {
			return err
		}

		// Determine if the function was hit.
		// A function is considered "hit" if at least one of its coverable lines was executed.
		isHit := false
		for _, line := range file.Lines {
			if line.Number >= codeElement.FirstLine && line.Number <= codeElement.LastLine {
				if line.Hits > 0 {
					isHit = true
					break
				}
			}
		}

		hitCount := 0
		if isHit {
			hitCount = 1 // LCOV FNDA typically uses 1 for hit, 0 for not hit.
			fnh++
		}
		if _, err := writer.WriteString(fmt.Sprintf("FNDA:%d,%s\n", hitCount, codeElement.FullName)); err != nil {
			return err
		}
	}

	if _, err := writer.WriteString(fmt.Sprintf("FNF:%d\n", fnf)); err != nil {
		return err
	}
	if _, err := writer.WriteString(fmt.Sprintf("FNH:%d\n", fnh)); err != nil {
		return err
	}

	// DA: Line data
	lf := 0 // Lines found (coverable)
	lh := 0 // Lines hit (covered)

	// Sort lines to ensure deterministic output for DA records
	sortedLines := make([]model.Line, len(file.Lines))
	copy(sortedLines, file.Lines)
	sort.SliceStable(sortedLines, func(i, j int) bool {
		return sortedLines[i].Number < sortedLines[j].Number
	})

	for _, line := range sortedLines {
		if line.LineVisitStatus != model.NotCoverable {
			lf++
			if line.Hits > 0 {
				lh++
			}
			if _, err := writer.WriteString(fmt.Sprintf("DA:%d,%d\n", line.Number, line.Hits)); err != nil {
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

	// BRDA/BRF/BRH: Branch data
	brf := 0 // Branches found
	brh := 0 // Branches hit
	for _, line := range sortedLines {
		if line.IsBranchPoint && len(line.Branch) > 0 {
			for blockIdx, branchDetail := range line.Branch {
				brf++
				// In LCOV, the 4th parameter is hits, or '-' if never taken.
				hits := "-"
				if branchDetail.Visits > 0 {
					brh++
					hits = strconv.Itoa(branchDetail.Visits)
				}
				// BRDA:<line_number>,<block_number>,<branch_number>,<hits>
				// block_number is usually 0 unless you have complex logic.
				// branch_number is the index of the branch on that line.
				if _, err := writer.WriteString(fmt.Sprintf("BRDA:%d,%d,%d,%s\n", line.Number, 0, blockIdx, hits)); err != nil {
					return err
				}
			}
		}
	}

	if brf > 0 {
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

// getAllFiles collects and deduplicates all CodeFile objects from the summary.
func getAllFiles(assemblies []model.Assembly) []*model.CodeFile {
	fileMap := make(map[string]*model.CodeFile)

	for _, assembly := range assemblies {
		for _, class := range assembly.Classes {
			for i := range class.Files {
				file := &class.Files[i]
				if _, exists := fileMap[file.Path]; !exists {
					fileMap[file.Path] = file
				}
			}
		}
	}

	files := make([]*model.CodeFile, 0, len(fileMap))
	for _, file := range fileMap {
		files = append(files, file)
	}
	return files
}
