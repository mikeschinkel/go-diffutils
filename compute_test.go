package diffutils

import (
	"testing"
)

func TestDiffLines_BothEmpty(t *testing.T) {
	content, err := DiffLines(nil, nil, "empty")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if content == nil {
		t.Fatal("expected non-nil content")
	}
	if len(content.Changes) != 0 {
		t.Errorf("expected 0 changes, got %d", len(content.Changes))
	}
	if content.Label != "empty" {
		t.Errorf("expected label 'empty', got %q", content.Label)
	}
}

func TestDiffLines_Identical(t *testing.T) {
	lines := []string{"line 1", "line 2", "line 3"}
	content, err := DiffLines(lines, lines, "identical")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(content.Changes) != 0 {
		t.Errorf("expected 0 changes for identical input, got %d", len(content.Changes))
	}
}

func TestDiffLines_AllNew(t *testing.T) {
	old := []string{}
	new := []string{"line 1", "line 2", "line 3"}
	content, err := DiffLines(old, new, "all-new")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(content.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(content.Changes))
	}
	c := content.Changes[0]
	if c.Type != LinesAdded {
		t.Errorf("expected LinesAdded, got %v", c.Type)
	}
	if c.NewRange.Count != 3 {
		t.Errorf("expected 3 new lines, got %d", c.NewRange.Count)
	}
}

func TestDiffLines_AllDeleted(t *testing.T) {
	old := []string{"line 1", "line 2"}
	new := []string{}
	content, err := DiffLines(old, new, "all-deleted")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(content.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(content.Changes))
	}
	c := content.Changes[0]
	if c.Type != LinesDeleted {
		t.Errorf("expected LinesDeleted, got %v", c.Type)
	}
	if c.OldRange.Count != 2 {
		t.Errorf("expected 2 old lines, got %d", c.OldRange.Count)
	}
}

func TestDiffLines_Mixed(t *testing.T) {
	old := []string{"alpha", "beta", "gamma", "delta"}
	new := []string{"alpha", "BETA", "gamma", "delta", "epsilon"}

	content, err := DiffLines(old, new, "mixed")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(content.OldLines) != 4 {
		t.Errorf("expected 4 old lines, got %d", len(content.OldLines))
	}
	if len(content.NewLines) != 5 {
		t.Errorf("expected 5 new lines, got %d", len(content.NewLines))
	}
	if len(content.Changes) == 0 {
		t.Fatal("expected at least 1 change")
	}

	// Verify we have both deletions and additions
	hasAdded := false
	hasDeleted := false
	for _, c := range content.Changes {
		switch c.Type {
		case LinesAdded:
			hasAdded = true
		case LinesDeleted:
			hasDeleted = true
		}
	}
	if !hasAdded {
		t.Error("expected at least one LinesAdded change")
	}
	if !hasDeleted {
		t.Error("expected at least one LinesDeleted change")
	}
}

func TestDiffLines_InsertInMiddle(t *testing.T) {
	old := []string{"line 1", "line 3"}
	new := []string{"line 1", "line 2", "line 3"}

	content, err := DiffLines(old, new, "insert-middle")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(content.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(content.Changes))
	}
	c := content.Changes[0]
	if c.Type != LinesAdded {
		t.Errorf("expected LinesAdded, got %v", c.Type)
	}
	if c.NewRange.Count != 1 {
		t.Errorf("expected 1 added line, got %d", c.NewRange.Count)
	}
	if c.NewRange.Start != 2 {
		t.Errorf("expected added line at position 2, got %d", c.NewRange.Start)
	}
}

func TestDiffLines_DeleteInMiddle(t *testing.T) {
	old := []string{"line 1", "line 2", "line 3"}
	new := []string{"line 1", "line 3"}

	content, err := DiffLines(old, new, "delete-middle")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(content.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(content.Changes))
	}
	c := content.Changes[0]
	if c.Type != LinesDeleted {
		t.Errorf("expected LinesDeleted, got %v", c.Type)
	}
	if c.OldRange.Count != 1 {
		t.Errorf("expected 1 deleted line, got %d", c.OldRange.Count)
	}
	if c.OldRange.Start != 2 {
		t.Errorf("expected deleted line at position 2, got %d", c.OldRange.Start)
	}
}

func TestChangeType_String(t *testing.T) {
	if LinesAdded.String() != "added" {
		t.Errorf("expected 'added', got %q", LinesAdded.String())
	}
	if LinesDeleted.String() != "deleted" {
		t.Errorf("expected 'deleted', got %q", LinesDeleted.String())
	}
	if ChangeType(99).String() != "unknown" {
		t.Errorf("expected 'unknown', got %q", ChangeType(99).String())
	}
}

func TestSequenceMatcher_Ratio(t *testing.T) {
	a := []string{"one", "two", "three"}
	b := []string{"one", "two", "three"}
	m := NewMatcher(a, b)
	ratio := m.Ratio()
	if ratio != 1.0 {
		t.Errorf("expected ratio 1.0 for identical sequences, got %f", ratio)
	}
}

func TestSequenceMatcher_GetOpCodes(t *testing.T) {
	a := []string{"one", "two", "three"}
	b := []string{"one", "TWO", "three"}
	m := NewMatcher(a, b)
	ops := m.GetOpCodes()

	if len(ops) == 0 {
		t.Fatal("expected non-empty opcodes")
	}

	// Should have: equal(one), replace(two→TWO), equal(three)
	hasReplace := false
	for _, op := range ops {
		if op.Tag == 'r' {
			hasReplace = true
		}
	}
	if !hasReplace {
		t.Error("expected a replace opcode")
	}
}

func TestGetUnifiedDiffString(t *testing.T) {
	diff := UnifiedDiff{
		A:        []string{"one\n", "two\n", "three\n"},
		B:        []string{"one\n", "TWO\n", "three\n"},
		FromFile: "a.txt",
		ToFile:   "b.txt",
		Context:  3,
	}
	result, err := GetUnifiedDiffString(diff)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == "" {
		t.Error("expected non-empty diff string")
	}
}
