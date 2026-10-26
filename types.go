package diffutils

// DiffContent holds diff data for two versions of text.
// It is the canonical interchange type for all diff operations in this package.
type DiffContent struct {
	OldLines []string     // Full "before" content, one entry per line
	NewLines []string     // Full "after" content, one entry per line
	Changes  []DiffChange // Change blocks (additions and deletions)
	Label    string       // Display label (filename, path, etc.)
}

// DiffChange represents a contiguous block of added or deleted lines.
type DiffChange struct {
	Type     ChangeType
	OldRange LineRange
	NewRange LineRange
}

// ChangeType indicates whether lines were added or deleted in a diff.
type ChangeType int

const (
	// LinesAdded means lines exist in new but not old.
	LinesAdded ChangeType = iota
	// LinesDeleted means lines exist in old but not new.
	LinesDeleted
)

// String returns a human-readable representation of the change type.
func (ct ChangeType) String() string {
	switch ct {
	case LinesAdded:
		return "added"
	case LinesDeleted:
		return "deleted"
	default:
		return "unknown"
	}
}

// LineRange is a 1-based line range.
type LineRange struct {
	Start int // 1-based line number where range starts
	Count int // Number of lines in the range
}
