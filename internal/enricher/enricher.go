// Package cache provides a content-addressed in-memory cache for static analysis results,
// backed by optional disk persistence.
//
// Cache keys are derived from file content hashes (SHA-256), meaning identical
// file contents always produce a cache hit regardless of filename or path. This
// makes it well-suited for incremental enrichment pipelines where source files
// rarely change between runs.
//
// # Basic Usage
//
//	manager, err := cache.NewManager("/path/to/cache/dir")
//	if err != nil { ... }
//
//	content, _ := os.ReadFile(path)
//
//	if cached, hit := manager.Get(content); hit {
//	    // use cached.TotalLines, cached.Result
//	} else {
//	    manager.Put(content, cache.CachedData{
//	        TotalLines: countLines(content),
//	        Result:     runAnalysis(content),
//	    })
//	}
//
//	// Persist to disk at the end of the run
//	if err := manager.Save(); err != nil { ... }
//
// # Persistence
//
// If a cache directory is provided to NewManager, the cache is loaded from disk
// on startup and can be flushed back via Save. If no directory is provided, or
// if initialization fails, the Manager operates as a pure in-memory cache with
// no persistence — callers should treat it as optional and non-fatal.
package enricher

import (
	"bytes"
	"log/slog"
	"os"
	"runtime"
	"sync"

	"github.com/IgorBayerl/nanovision/internal/analyzer"
	"github.com/IgorBayerl/nanovision/internal/cache"
	"github.com/IgorBayerl/nanovision/internal/filereader"
	"github.com/IgorBayerl/nanovision/internal/model"
	"github.com/IgorBayerl/nanovision/internal/utils"
)

type Enricher struct {
	analyzers    []analyzer.Analyzer
	fileReader   filereader.Reader
	logger       *slog.Logger
	cacheManager *cache.Manager
	ignoreCache  bool
}

func New(analyzers []analyzer.Analyzer, fileReader filereader.Reader, logger *slog.Logger, cacheDir string, ignoreCache bool) *Enricher {
	var cm *cache.Manager

	if cacheDir != "" && !ignoreCache {
		var err error
		cm, err = cache.NewManager(cacheDir)
		if err != nil {
			logger.Warn("Failed to initialize cache, proceeding without it", "error", err)
		}
	}

	return &Enricher{
		analyzers:    analyzers,
		fileReader:   fileReader,
		logger:       logger,
		cacheManager: cm,
		ignoreCache:  ignoreCache,
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

	if e.cacheManager != nil {
		if err := e.cacheManager.Save(); err != nil {
			e.logger.Warn("Failed to persist analysis cache", "error", err)
		} else {
			e.logger.Debug("Analysis cache saved successfully")
		}
	}
}

// enrichFileNode performs the enrichment process for a single file.
// This includes line counting and static code analysis.
// It is designed to be called concurrently.
func (e *Enricher) enrichFileNode(fileNode *model.FileNode) {
	path := fileNode.Path

	sourceDirs := []string{fileNode.SourceDir}
	absPath, err := utils.FindFileInSourceDirs(path, sourceDirs, e.fileReader, e.logger)
	if err != nil {
		e.logger.Warn("Source file not found", "file", path)
		return
	}

	// Read Content ONCE (needed for both hash and analysis)
	content, err := os.ReadFile(absPath)
	if err != nil {
		e.logger.Warn("Could not read source file", "file", path, "error", err)
		return
	}

	// CACHE CHECK
	if e.cacheManager != nil {
		if cached, hit := e.cacheManager.Get(content); hit {
			// CACHE HIT: Apply cached data and return
			fileNode.TotalLines = cached.TotalLines
			fileNode.Metrics.TotalLines = cached.TotalLines
			e.applyAnalysisToFileNode(fileNode, cached.Result)
			return
		}
	}

	// CACHE MISS: Perform Calculations

	// Calculate Total Lines (simple implementation matching standard text editors)
	totalLines := 0
	if len(content) > 0 {
		totalLines = bytes.Count(content, []byte{'\n'})
		if content[len(content)-1] != '\n' {
			totalLines++
		}
	}
	fileNode.TotalLines = totalLines
	fileNode.Metrics.TotalLines = totalLines

	// Run Static Analysis
	activeAnalyzer := e.findAnalyzerForFile(path)
	var analysisResult analyzer.AnalysisResult

	if activeAnalyzer != nil {
		e.logger.Info("Analyzing file", "language", activeAnalyzer.Name(), "file", path)

		// Use the renamed variable here
		var err error
		analysisResult, err = activeAnalyzer.Analyze(content)

		if err != nil {
			e.logger.Warn("Static analysis failed", "file", path, "error", err)
			return
		}
		e.applyAnalysisToFileNode(fileNode, analysisResult)
	} else {
		// If no analyzer found, cache just the line count to avoid re-reading file later
		if e.cacheManager != nil {
			e.cacheManager.Put(content, cache.CachedData{TotalLines: totalLines})
		}
		return
	}

	// UPDATE CACHE
	if e.cacheManager != nil {
		e.cacheManager.Put(content, cache.CachedData{
			TotalLines: totalLines,
			Result:     analysisResult,
		})
	}
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
