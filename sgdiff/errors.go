package sgdiff

import "errors"

var (
	// ErrSgDiff is the layer sentinel for sgdiff errors.
	ErrSgDiff = errors.New("sgdiff")

	// ErrParseDiff indicates a failure to parse unified diff text.
	ErrParseDiff = errors.New("parse diff")
)
