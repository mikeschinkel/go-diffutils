package sgdiff

import (
	"testing"
)

func TestParseUnifiedDiff_Simple(t *testing.T) {
	diffText := `--- a/hello.go
+++ b/hello.go
@@ -1,3 +1,4 @@
 package main

+import "fmt"
 func main() {
`

	contents, err := ParseUnifiedDiff(diffText)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(contents) != 1 {
		t.Fatalf("expected 1 file diff, got %d", len(contents))
	}

	dc := contents[0]
	if dc.Label != "hello.go" {
		t.Errorf("expected label 'hello.go', got %q", dc.Label)
	}
	if len(dc.Changes) == 0 {
		t.Error("expected at least 1 change")
	}
}

func TestParseUnifiedDiff_MultiFile(t *testing.T) {
	diffText := `--- a/file1.go
+++ b/file1.go
@@ -1,2 +1,3 @@
 package main
+// added comment

--- a/file2.go
+++ b/file2.go
@@ -1,3 +1,2 @@
 package util
-// removed comment

`

	contents, err := ParseUnifiedDiff(diffText)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(contents) != 2 {
		t.Fatalf("expected 2 file diffs, got %d", len(contents))
	}
	if contents[0].Label != "file1.go" {
		t.Errorf("expected label 'file1.go', got %q", contents[0].Label)
	}
	if contents[1].Label != "file2.go" {
		t.Errorf("expected label 'file2.go', got %q", contents[1].Label)
	}
}

func TestParseUnifiedDiff_Empty(t *testing.T) {
	contents, err := ParseUnifiedDiff("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(contents) != 0 {
		t.Errorf("expected 0 contents for empty input, got %d", len(contents))
	}
}

func TestParseUnifiedDiffWithContent(t *testing.T) {
	diffText := `--- a/test.go
+++ b/test.go
@@ -2,2 +2,3 @@
 line2
+inserted
 line3
`
	oldFull := []string{"line1", "line2", "line3"}
	newFull := []string{"line1", "line2", "inserted", "line3"}

	content, err := ParseUnifiedDiffWithContent(diffText, oldFull, newFull)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if content == nil {
		t.Fatal("expected non-nil content")
	}
	if len(content.OldLines) != 3 {
		t.Errorf("expected 3 old lines, got %d", len(content.OldLines))
	}
	if len(content.NewLines) != 4 {
		t.Errorf("expected 4 new lines, got %d", len(content.NewLines))
	}
	if len(content.Changes) == 0 {
		t.Error("expected at least 1 change")
	}
}
