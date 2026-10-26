package dmp

import (
	"strings"

	"github.com/mikeschinkel/go-diffutils"
	dmp "github.com/sergi/go-diff/diffmatchpatch"
)

// Ensure Highlighter implements InlineHighlighter.
var _ diffutils.InlineHighlighter = (*Highlighter)(nil)

// HighlighterArgs holds configuration for NewHighlighter.
type HighlighterArgs struct {
	// DeletedColor is the ANSI color code for deleted characters (default: red "31").
	DeletedColor string
	// InsertedColor is the ANSI color code for inserted characters (default: green "32").
	InsertedColor string
}

// Highlighter implements [diffutils.InlineHighlighter] using diff-match-patch
// for character-level diffs between two lines.
type Highlighter struct {
	deletedColor  string
	insertedColor string
}

// NewHighlighter creates a Highlighter with optional configuration.
func NewHighlighter(args *HighlighterArgs) *Highlighter {
	h := &Highlighter{
		deletedColor:  "31", // Red
		insertedColor: "32", // Green
	}
	if args != nil {
		if args.DeletedColor != "" {
			h.deletedColor = args.DeletedColor
		}
		if args.InsertedColor != "" {
			h.insertedColor = args.InsertedColor
		}
	}
	return h
}

// HighlightInline returns old and new lines with ANSI escape sequences
// marking deleted and inserted characters respectively.
func (h *Highlighter) HighlightInline(oldLine, newLine string) (oldResult, newResult string) {
	differ := dmp.New()
	diffs := differ.DiffMain(oldLine, newLine, false)
	diffs = differ.DiffCleanupSemantic(diffs)

	var oldBuf, newBuf strings.Builder

	for _, d := range diffs {
		switch d.Type {
		case dmp.DiffEqual:
			oldBuf.WriteString(d.Text)
			newBuf.WriteString(d.Text)
		case dmp.DiffDelete:
			oldBuf.WriteString("\x1b[" + h.deletedColor + "m")
			oldBuf.WriteString(d.Text)
			oldBuf.WriteString("\x1b[0m")
		case dmp.DiffInsert:
			newBuf.WriteString("\x1b[" + h.insertedColor + "m")
			newBuf.WriteString(d.Text)
			newBuf.WriteString("\x1b[0m")
		}
	}

	return oldBuf.String(), newBuf.String()
}
