package diffutils

// DiffSource produces diff content from some origin.
// Implementations might read from Git, files, databases, etc.
type DiffSource interface {
	DiffContent() (*DiffContent, error)
}

// InlineHighlighter provides character-level highlighting within changed lines.
// Implementations typically use diff-match-patch or similar algorithms to
// identify which characters changed between two versions of a line.
type InlineHighlighter interface {
	HighlightInline(oldLine, newLine string) (oldResult, newResult string)
}
