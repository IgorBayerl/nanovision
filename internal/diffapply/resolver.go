package diffapply

import (
	"log/slog"
	"strings"

	"github.com/IgorBayerl/nanovision/internal/diff"
	"github.com/IgorBayerl/nanovision/internal/model"
)

// Resolver maps paths from diff files to paths in the coverage tree
type Resolver interface {
	// Resolve attempts to map a diff path to a tree path
	// Returns the resolved path and whether resolution was successful
	Resolve(diffPath string) (treeRel string, ok bool)
}

type resolverImpl struct {
	cache       map[string]string   // Cache of resolved paths
	warned      map[string]struct{} // Set of paths we've warned about
	fileIndex   map[string]*model.FileNode
	coveredSet  map[string]bool
	logger      *slog.Logger
	stripPrefix string   // For H2: common prefix to strip
	stripN      int      // For H1: number of components to strip (-pN)
	candidates  []string // For H3/H4: potential matches for suffix matching
}

// BuildResolver creates a new Resolver using multiple heuristics:
// H1: Strip N leading components that gives most exact matches
// H2: Remove longest common prefix from monorepo root
// H3: Match by basename and suffix components
// H4: Prefer paths with known coverage in ambiguous cases
func BuildResolver(dd *diff.DiffData, fileIndex map[string]*model.FileNode, coveredSet map[string]bool, logger *slog.Logger) Resolver {
	r := &resolverImpl{
		cache:      make(map[string]string),
		warned:     make(map[string]struct{}),
		fileIndex:  fileIndex,
		coveredSet: coveredSet,
		logger:     logger,
	}

	// Collect all diff paths and tree paths
	diffPaths := make([]string, 0, len(dd.Files))
	treePaths := make([]string, 0, len(fileIndex))
	for _, f := range dd.Files {
		normalized := diff.Normalize(f.NewPath)
		if logger != nil {
			logger.Debug("Normalized diff path", "raw", f.NewPath, "normalized", normalized)
		}
		diffPaths = append(diffPaths, normalized)
	}
	for path := range fileIndex {
		treePaths = append(treePaths, path)
	}

	// Try H1: Strip N path components
	bestN := r.findBestStripN(diffPaths, treePaths)
	if bestN > 0 {
		r.stripN = bestN
		return r
	}

	// Try H2: Common prefix
	if prefix := r.findCommonPrefix(diffPaths); prefix != "" {
		r.stripPrefix = prefix
		return r
	}

	// Store candidates for H3/H4
	r.candidates = treePaths
	return r
}

func (r *resolverImpl) Resolve(diffPath string) (string, bool) {
	// Check cache first
	if resolved, ok := r.cache[diffPath]; ok {
		return resolved, true
	}

	// Normalize the diff path
	diffPath = diff.Normalize(diffPath)

	// Check for a direct, exact match first before trying heuristics.
	if _, ok := r.fileIndex[diffPath]; ok {
		r.cache[diffPath] = diffPath
		return diffPath, true
	}

	// Try H1: Strip N components
	if r.stripN > 0 {
		parts := strings.Split(diffPath, "/")
		if len(parts) > r.stripN {
			result := strings.Join(parts[r.stripN:], "/")
			if _, ok := r.fileIndex[result]; ok {
				r.cache[diffPath] = result
				return result, true
			}
		}
	}

	// Try H2: Strip common prefix
	if r.stripPrefix != "" && strings.HasPrefix(diffPath, r.stripPrefix) {
		result := strings.TrimPrefix(diffPath, r.stripPrefix)
		if _, ok := r.fileIndex[result]; ok {
			r.cache[diffPath] = result
			return result, true
		}
	}

	// Try H3/H4: Suffix matching with coverage tie-breaking
	result, ok := r.resolveBySuffix(diffPath)
	if ok {
		r.cache[diffPath] = result
		return result, true
	}

	// Log warning for unmapped paths
	if _, warned := r.warned[diffPath]; !warned {
		r.logger.Warn("unable to map diff path", "path", diffPath)
		r.warned[diffPath] = struct{}{}
	}

	return "", false
}

func (r *resolverImpl) findBestStripN(diffPaths, treePaths []string) int {
	bestN := 0
	bestMatches := 0
	bestDepth := 0

	// Try stripping 1 to maxN components
	maxN := 5 // Reasonable limit for monorepo nesting
	for n := 1; n <= maxN; n++ {
		matches := 0
		totalDepth := 0

		for _, dp := range diffPaths {
			parts := strings.Split(dp, "/")
			if len(parts) <= n {
				continue
			}
			stripped := strings.Join(parts[n:], "/")
			for _, tp := range treePaths {
				if stripped == tp {
					matches++
					totalDepth += len(strings.Split(tp, "/"))
				}
			}
		}

		// Update best if we have more matches or same matches but deeper paths
		if matches > bestMatches || (matches == bestMatches && totalDepth > bestDepth) {
			bestN = n
			bestMatches = matches
			bestDepth = totalDepth
		}
	}

	return bestN
}

func (r *resolverImpl) findCommonPrefix(paths []string) string {
	if len(paths) == 0 {
		return ""
	}

	// Find the shortest path
	minLen := len(paths[0])
	for _, p := range paths[1:] {
		if len(p) < minLen {
			minLen = len(p)
		}
	}

	// Find common prefix
	prefix := ""
	for i := 0; i < minLen; i++ {
		char := paths[0][i]
		for _, p := range paths[1:] {
			if p[i] != char {
				return prefix
			}
		}
		prefix += string(char)
	}

	// Only use prefix up to last slash
	if lastSlash := strings.LastIndex(prefix, "/"); lastSlash != -1 {
		return prefix[:lastSlash+1]
	}
	return ""
}

func (r *resolverImpl) resolveBySuffix(diffPath string) (string, bool) {
	diffParts := strings.Split(diffPath, "/")
	baseName := diffParts[len(diffParts)-1]

	var bestMatch string
	var bestScore int
	var ambiguousCandidates []string

	for _, candidate := range r.candidates {
		if !strings.HasSuffix(candidate, baseName) {
			continue
		}

		// Score based on matching suffix components
		candParts := strings.Split(candidate, "/")
		score := 0
		for i := 1; i <= len(diffParts) && i <= len(candParts); i++ {
			if diffParts[len(diffParts)-i] == candParts[len(candParts)-i] {
				score++
			} else {
				break
			}
		}

		// H4: Add bonus for covered files
		if r.coveredSet[candidate] {
			score++
		}

		if score > bestScore {
			bestScore = score
			bestMatch = candidate
			ambiguousCandidates = []string{candidate}
		} else if score == bestScore && bestMatch != "" {
			// If both paths have coverage, or neither has coverage, it's truly ambiguous
			if r.coveredSet[candidate] == r.coveredSet[bestMatch] {
				ambiguousCandidates = append(ambiguousCandidates, candidate)
			} else if r.coveredSet[candidate] {
				// If the new candidate has coverage but best match doesn't, use the candidate
				bestMatch = candidate
				ambiguousCandidates = []string{candidate}
			}
			// If best match has coverage but candidate doesn't, keep best match
		}
	}

	// Check for ambiguity: multiple candidates with same score and coverage status
	if len(ambiguousCandidates) > 1 {
		if _, warned := r.warned[diffPath]; !warned {
			r.logger.Warn("ambiguous path resolution", "path", diffPath, "candidates", ambiguousCandidates)
			r.warned[diffPath] = struct{}{}
		}
		return "", false
	}

	if bestMatch != "" {
		return bestMatch, true
	}

	return "", false
}
