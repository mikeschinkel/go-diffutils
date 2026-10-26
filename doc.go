// Package diffutils provides canonical types and algorithms for computing
// and representing text diffs in Go.
//
// The core diff engine is based on pmezard/go-difflib (MIT licensed),
// which implements a port of Python's difflib SequenceMatcher. This
// package absorbs that code and adds:
//
//   - Canonical interchange types ([DiffContent], [DiffChange], [LineRange])
//   - A high-level [DiffLines] function for comparing string slices
//   - Interfaces ([DiffSource], [InlineHighlighter]) for extensibility
//
// Adapter sub-packages bridge major Go diff libraries:
//
//   - [github.com/mikeschinkel/go-diffutils/sgdiff] — sourcegraph/go-diff adapter
//   - [github.com/mikeschinkel/go-diffutils/dmp] — sergi/go-diff inline highlighter
//
// # Stability
//
// This package is provisional as of v0.1.0. The public API may change in
// minor releases until promoted to stable.
package diffutils
