package diffutils

import (
	"errors"
	"testing"
)

var errWriteFailed = errors.New("write failed")

// failingWriter fails every write. A diff smaller than bufio's buffer reaches
// it only when the buffer is flushed, so it catches a dropped flush error.
type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) {
	return 0, errWriteFailed
}

func TestWriteUnifiedDiff_ReportsFlushError(t *testing.T) {
	err := WriteUnifiedDiff(failingWriter{}, UnifiedDiff{
		A:        []string{"a\n"},
		B:        []string{"b\n"},
		FromFile: "old",
		ToFile:   "new",
	})
	if !errors.Is(err, errWriteFailed) {
		t.Errorf("expected %v, got %v", errWriteFailed, err)
	}
}

func TestWriteContextDiff_ReportsFlushError(t *testing.T) {
	err := WriteContextDiff(failingWriter{}, ContextDiff{
		A:        []string{"a\n"},
		B:        []string{"b\n"},
		FromFile: "old",
		ToFile:   "new",
	})
	if !errors.Is(err, errWriteFailed) {
		t.Errorf("expected %v, got %v", errWriteFailed, err)
	}
}
