package model

// ChangeKind represents the type of change detected for a file
type ChangeKind int

const (
	// ChangeKindNone indicates no changes detected
	ChangeKindNone ChangeKind = iota
	// ChangeKindAdded indicates the file was added
	ChangeKindAdded
	// ChangeKindModified indicates the file was modified
	ChangeKindModified
)

// String returns the string representation of ChangeKind
func (k ChangeKind) String() string {
	switch k {
	case ChangeKindAdded:
		return "added"
	case ChangeKindModified:
		return "modified"
	default:
		return ""
	}
}

// MarshalJSON implements json.Marshaler interface
func (k ChangeKind) MarshalJSON() ([]byte, error) {
	return []byte(`"` + k.String() + `"`), nil
}

// UnmarshalJSON implements json.Unmarshaler interface
func (k *ChangeKind) UnmarshalJSON(data []byte) error {
	// Remove quotes
	s := string(data[1 : len(data)-1])
	switch s {
	case "added":
		*k = ChangeKindAdded
	case "modified":
		*k = ChangeKindModified
	default:
		*k = ChangeKindNone
	}
	return nil
}

// DiffInfo stores information about changes detected in a file
type DiffInfo struct {
	// Kind indicates the type of change detected
	Kind ChangeKind `json:"kind,omitempty"`
	// AddedLines maps line numbers to whether they were added (true if added)
	AddedLines map[int]bool `json:"addedLines,omitempty"`
	// ModifiedLines maps line numbers to whether they were modified (true if modified)
	ModifiedLines map[int]bool `json:"modifiedLines,omitempty"`
}
