// Package enricher is responsible for augmenting the primary coverage data with
// additional metrics gathered from static analysis of the source code.
package enricher

import (
	"log/slog"
	"os"
	"runtime"
	"sync"

	"github.com/IgorBayerl/nanovision/analyzer"
	"github.com/IgorBayerl/nanovision/filereader"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/utils"
)

type Enricher struct {
	analyzers  []analyzer.Analyzer
	fileReader filereader.Reader
	logger     *slog.Logger
}

func New(analyzers []analyzer.Analyzer, fileReader filereader.Reader, logger *slog.Logger) *Enricher {
	return &Enricher{
		analyzers:  analyzers,
		fileReader: fileReader,
		logger:     logger,
	}
}

// findAnalyzerForFile iterates through the available analyzers to find one that
// supports the given file path.
//
// This allows the Enricher to be language-agnostic, dynamically selecting the
// correct tool (e.g., the Go analyzer for a '.go' file) for static analysis.
func (e *Enricher) findAnalyzerForFile(filePath string) analyzer.Analyzer {
	for _, analyzer := range e.analyzers {
		if analyzer.SupportsFile(filePath) {
			return analyzer
		}
	}
	return nil
}

// EnrichTree is the main entry point for the enrichment process. It traverses
// the entire model.SummaryTree, finds every file, and applies two key enhancements:
//
//   - It calculates the total number of lines in each source file, providing an
//     accurate denominator for 'total lines' metrics.
//   - It performs static code analysis on supported file types to extract
//     method-level details, such as cyclomatic complexity.
//
// This method modifies the tree in place, adding the new data directly to the
// FileNode objects.
func (e *Enricher) EnrichTree(tree *model.SummaryTree) {
	fileNodeMap := make(map[string]*model.FileNode)
	collectFiles(tree.Root, fileNodeMap)

	numWorkers := runtime.NumCPU()
	jobs := make(chan *model.FileNode, len(fileNodeMap))
	var wg sync.WaitGroup

	// Start worker goroutines
	for i := 0; i < numWorkers; i++ {
		go func() {
			for fileNode := range jobs {
				e.enrichFileNode(fileNode)
				wg.Done()
			}
		}()
	}

	// Send jobs to the workers
	for _, fileNode := range fileNodeMap {
		wg.Add(1)
		jobs <- fileNode
	}
	close(jobs)

	// Wait for all jobs to complete
	wg.Wait()
}

// enrichFileNode performs the enrichment process for a single file.
// This includes line counting and static code analysis.
// It is designed to be called concurrently.
func (e *Enricher) enrichFileNode(fileNode *model.FileNode) {
	path := fileNode.Path

	// Count the total number of lines in the source file.
	if abs, err := utils.FindFileInSourceDirs(path, []string{fileNode.SourceDir}, e.fileReader, e.logger); err == nil {
		if n, err := e.fileReader.CountLines(abs); err == nil {
			// Set the total lines. Aggregation is handled later, so no need to update parents here.
			fileNode.TotalLines = n
			fileNode.Metrics.TotalLines = n
		} else {
			e.logger.Warn("Could not count lines", "file", abs, "error", err)
		}
	} else {
		e.logger.Warn("Source file not found for line counting", "file", path, "error", err)
	}

	// Find a suitable analyzer for the file.
	analyzer := e.findAnalyzerForFile(path)
	if analyzer == nil {
		return // No analysis needed for this file type.
	}

	e.logger.Info("Analyzing file", "path", path, "analyzer", analyzer.Name())
	sourceBytes, err := e.readSourceFile(fileNode)
	if err != nil {
		e.logger.Warn("Could not read source file for analysis", "file", path, "error", err)
		return
	}

	analysis, err := analyzer.Analyze(sourceBytes)
	if err != nil {
		e.logger.Warn("Static analysis failed for file", "file", path, "error", err)
		return
	}

	e.applyAnalysisToFileNode(fileNode, analysis)
}

// readSourceFile locates and reads the content of a source file from disk.
// It uses the file's associated source directory to resolve its absolute path
// via utils.FindFileInSourceDirs. The file content is returned as a byte slice,
// which is the required input for the static code analyzers.
func (e *Enricher) readSourceFile(fileNode *model.FileNode) ([]byte, error) {
	sourceDirs := []string{fileNode.SourceDir}
	absPath, err := utils.FindFileInSourceDirs(fileNode.Path, sourceDirs, e.fileReader, e.logger)
	if err != nil {
		return nil, err
	}
	return os.ReadFile(absPath)
}

// applyAnalysisToFileNode translates the generic results from an analyzer into
// the specific data structures of the application's model.
//
// It iterates through the functions found by the analyzer, converts them into
// model.MethodMetrics, calculates their specific code coverage, and attaches
// them to the FileNode.
func (e *Enricher) applyAnalysisToFileNode(fileNode *model.FileNode, analysis analyzer.AnalysisResult) {
	var methodMetrics []model.MethodMetrics
	for _, funcMetric := range analysis.Functions {
		metric := model.MethodMetrics{
			Name:                 funcMetric.Name,
			StartLine:            funcMetric.Position.StartLine,
			EndLine:              funcMetric.Position.EndLine,
			CyclomaticComplexity: funcMetric.CyclomaticComplexity,
		}
		calculateMethodCoverage(fileNode, &metric)
		methodMetrics = append(methodMetrics, metric)
	}
	fileNode.Methods = methodMetrics
}

// calculateMethodCoverage computes the line and branch coverage for a single method
// by examining the coverage data of the lines within its start and end boundaries.
//
// This provides a more granular view than the overall file coverage, helping to
// identify specific functions that are poorly tested. For example, if a method
// spans lines 10 to 20, this function will sum the covered lines and branches
// only within that range from the parent file's line data.
func calculateMethodCoverage(file *model.FileNode, method *model.MethodMetrics) {
	for i := method.StartLine; i <= method.EndLine; i++ {
		if line, ok := file.Lines[i]; ok {
			if line.Hits >= 0 {
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

// collectFiles performs a recursive walk of the directory tree starting from a
// DirNode and populates a map with all the FileNode objects it finds. The map
// keys are the full file paths.
//
// This exists to simplify the enrichment process by providing a flat list of all
// files that need to be analyzed, avoiding the need to repeatedly traverse the
// tree structure.
func collectFiles(dir *model.DirNode, fileMap map[string]*model.FileNode) {
	for _, file := range dir.Files {
		fileMap[file.Path] = file
	}
	for _, subDir := range dir.Subdirs {
		collectFiles(subDir, fileMap)
	}
}
