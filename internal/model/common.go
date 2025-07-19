package model

// LineVisitStatus indicates the coverage status of a line in a source file.
type LineVisitStatus int

const (
	// NotCoverable means the line cannot be covered (e.g., comments, empty lines).
	NotCoverable LineVisitStatus = iota
	// NotCovered means the line is coverable but was not executed.
	NotCovered
	// PartiallyCovered means the line is a branch point and only some branches were executed.
	PartiallyCovered
	// Covered means the line was fully executed (or all branches covered if it's a branch point).
	Covered
)

// MetricStatus, CodeElementType etc. can also go here if not already present.
