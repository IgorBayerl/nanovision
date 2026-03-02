package cache

import "time"

// CacheValidator determines if a cached entry is still valid.
//
// In this system, static analysis results are cached using the hash of
// the file contents as the primary key. However, if the rules for the
// static analyzers change (e.g. by updating tree-sitter queries),
// previous cached results become stale and inaccurate even if the target
// file's source code hasn't changed.
//
// This interface allows validating existing cache entries against the
// current execution environment metadata.
//
// If you need to change the validation strategy (e.g. adding a new
// property like 'EnvironmentType' to invalidate cache differently
// between dev and prod paths), you can:
// 1. Add the new field to both BuildMetadata and CacheMetadata.
// 2. Adjust or create a new implementation of CacheValidator below
//    that accounts for your new rules in the `IsValid` method.
type CacheValidator interface {
	// IsValid checks if the cached data is still applicable given current context
	IsValid(metadata CacheMetadata, currentMetadata BuildMetadata) bool
}

// BuildMetadata represents the build-time context
type BuildMetadata struct {
	CommitHash      string // Injected at build time
	AnalyzerVersion string // Version of your analysis engine
}

// CacheMetadata is stored WITH the cached data
type CacheMetadata struct {
	CommitHash      string
	AnalyzerVersion string
	CachedAt        time.Time
}

type StrictValidator struct {
	CurrentBuildMetadata BuildMetadata
}

// IsValid returns true only if BOTH commit and analyzer versions match
func (sv *StrictValidator) IsValid(metadata CacheMetadata, _ BuildMetadata) bool {
	return metadata.CommitHash == sv.CurrentBuildMetadata.CommitHash &&
		metadata.AnalyzerVersion == sv.CurrentBuildMetadata.AnalyzerVersion
}

type RelaxedValidator struct {
	CurrentBuildMetadata BuildMetadata
}

// IsValid returns true if only analyzer version matches (useful for dev)
func (rv *RelaxedValidator) IsValid(metadata CacheMetadata, _ BuildMetadata) bool {
	return metadata.AnalyzerVersion == rv.CurrentBuildMetadata.AnalyzerVersion
}

// DevValidator always invalidates on dev builds
type DevValidator struct{}

func (dv *DevValidator) IsValid(metadata CacheMetadata, _ BuildMetadata) bool {
	return false // Always invalidate in dev
}
