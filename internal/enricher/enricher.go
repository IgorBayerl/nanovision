// internal/enricher/enricher.go
package enricher

import (
	"log/slog"
	"os"

	"github.com/IgorBayerl/AdlerCov/analyzer"
	"github.com/IgorBayerl/AdlerCov/filereader"
	"github.com/IgorBayerl/AdlerCov/internal/model"
	"github.com/IgorBayerl/AdlerCov/internal/utils"
)

// Enricher applies static analysis to the file nodes in a summary tree.
type Enricher struct {
	analyzers  []analyzer.Analyzer // The list of all available analyzers.
	fileReader filereader.Reader
	logger     *slog.Logger
}

// New creates a new Enricher, accepting a slice of all available analyzers.
func New(analyzers []analyzer.Analyzer, fileReader filereader.Reader, logger *slog.Logger) *Enricher {
	return &Enricher{
		analyzers:  analyzers,
		fileReader: fileReader,
		logger:     logger,
	}
}

// findAnalyzerForFile iterates through the available analyzers to find one that supports the file.
func (e *Enricher) findAnalyzerForFile(filePath string) analyzer.Analyzer {
	for _, analyzer := range e.analyzers {
		if analyzer.SupportsFile(filePath) {
			return analyzer
		}
	}
	return nil // No suitable analyzer found.
}

// EnrichTree iterates through all files in the tree and applies static analysis.
func (e *Enricher) EnrichTree(tree *model.SummaryTree) {
	fileNodeMap := make(map[string]*model.FileNode)
	collectFiles(tree.Root, fileNodeMap)

	for path, fileNode := range fileNodeMap {
		// Always compute total lines so the UI can show the "Total" column.
		if abs, err := utils.FindFileInSourceDirs(fileNode.Path, []string{fileNode.SourceDir}, e.fileReader, e.logger); err == nil {
			if n, err := e.fileReader.CountLines(abs); err == nil {
				// sync file + metrics and propagate delta upwards
				old := fileNode.Metrics.TotalLines
				fileNode.TotalLines = n
				fileNode.Metrics.TotalLines = n
				if delta := n - old; delta != 0 {
					// bump all ancestor folders
					for p := fileNode.Parent; p != nil; p = p.Parent {
						p.Metrics.TotalLines += delta
					}
					// bump the tree totals too
					tree.Metrics.TotalLines += delta
				}
			} else {
				e.logger.Warn("Could not count lines", "file", abs, "error", err)
			}
		} else {
			e.logger.Warn("Source file not found for line counting", "file", fileNode.Path, "error", err)
		}

		analyzer := e.findAnalyzerForFile(path)
		if analyzer == nil {
			continue // No analyzer available for this file type.
		}

		e.logger.Info("Analyzing file", "path", path, "analyzer", analyzer.Name())
		sourceBytes, err := e.readSourceFile(fileNode)
		if err != nil {
			e.logger.Warn("Could not read source file for analysis", "file", path, "error", err)
			continue
		}

		analysis, err := analyzer.Analyze(sourceBytes)
		if err != nil {
			e.logger.Warn("Static analysis failed for file", "file", path, "error", err)
			continue
		}

		e.applyAnalysisToFileNode(fileNode, analysis)
	}
}

func (e *Enricher) readSourceFile(fileNode *model.FileNode) ([]byte, error) {
	// The source dir associated with the file node is the most reliable one.
	sourceDirs := []string{fileNode.SourceDir}
	absPath, err := utils.FindFileInSourceDirs(fileNode.Path, sourceDirs, e.fileReader, e.logger)
	if err != nil {
		return nil, err
	}
	return os.ReadFile(absPath)
}

func (e *Enricher) applyAnalysisToFileNode(fileNode *model.FileNode, analysis analyzer.AnalysisResult) {
	var methodMetrics []model.MethodMetrics
	for _, funcMetric := range analysis.Functions {
		metric := model.MethodMetrics{
			Name:                 funcMetric.Name,
			StartLine:            funcMetric.Position.StartLine,
			EndLine:              funcMetric.Position.EndLine,
			CyclomaticComplexity: funcMetric.CyclomaticComplexity,
		}
		// Calculate coverage metrics for this specific method
		calculateMethodCoverage(fileNode, &metric)
		methodMetrics = append(methodMetrics, metric)
	}
	fileNode.Methods = methodMetrics
}

// calculateMethodCoverage computes the line and branch coverage for a single method.
func calculateMethodCoverage(file *model.FileNode, method *model.MethodMetrics) {
	for i := method.StartLine; i <= method.EndLine; i++ {
		if line, ok := file.Lines[i]; ok {
			if line.Hits >= 0 { // Is a coverable line
				method.LinesValid++
				if line.Hits > 0 {
					method.LinesCovered++
				}
			}
			method.BranchesValid += line.TotalBranches
			method.BranchesCovered += line.CoveredBranches
		}
	}
}

func collectFiles(dir *model.DirNode, fileMap map[string]*model.FileNode) {
	for _, file := range dir.Files {
		fileMap[file.Path] = file
	}
	for _, subDir := range dir.Subdirs {
		collectFiles(subDir, fileMap)
	}
}
