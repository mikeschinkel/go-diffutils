package dmp

import (
	"strings"
	"testing"
)

func TestHighlighter_HighlightInline_Basic(t *testing.T) {
	h := NewHighlighter(nil)
	oldResult, newResult := h.HighlightInline("hello world", "hello there")

	// Old result should contain "world" highlighted as deleted
	if !strings.Contains(oldResult, "world") {
		t.Error("expected old result to contain 'world'")
	}
	// New result should contain "there" highlighted as inserted
	if !strings.Contains(newResult, "there") {
		t.Error("expected new result to contain 'there'")
	}

	// Both should contain ANSI escape codes
	if !strings.Contains(oldResult, "\x1b[") {
		t.Error("expected ANSI codes in old result")
	}
	if !strings.Contains(newResult, "\x1b[") {
		t.Error("expected ANSI codes in new result")
	}
}

func TestHighlighter_HighlightInline_Identical(t *testing.T) {
	h := NewHighlighter(nil)
	oldResult, newResult := h.HighlightInline("same line", "same line")

	// No highlighting for identical lines
	if strings.Contains(oldResult, "\x1b[") {
		t.Error("expected no ANSI codes for identical lines (old)")
	}
	if strings.Contains(newResult, "\x1b[") {
		t.Error("expected no ANSI codes for identical lines (new)")
	}
	if oldResult != "same line" {
		t.Errorf("expected 'same line', got %q", oldResult)
	}
}

func TestHighlighter_HighlightInline_Empty(t *testing.T) {
	h := NewHighlighter(nil)
	oldResult, newResult := h.HighlightInline("", "added")

	if oldResult != "" {
		t.Errorf("expected empty old result, got %q", oldResult)
	}
	if !strings.Contains(newResult, "added") {
		t.Error("expected new result to contain 'added'")
	}
}

func TestHighlighter_CustomColors(t *testing.T) {
	h := NewHighlighter(&HighlighterArgs{
		DeletedColor:  "91", // Bright red
		InsertedColor: "92", // Bright green
	})
	oldResult, newResult := h.HighlightInline("foo", "bar")

	if !strings.Contains(oldResult, "\x1b[91m") {
		t.Error("expected bright red ANSI code in old result")
	}
	if !strings.Contains(newResult, "\x1b[92m") {
		t.Error("expected bright green ANSI code in new result")
	}
}
