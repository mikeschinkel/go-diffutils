package diffutils

import "errors"

var (
	// ErrDiffUtils is the layer sentinel for all diffutils errors.
	ErrDiffUtils = errors.New("diffutils")

	// ErrEmptyInput indicates that both old and new inputs are empty.
	ErrEmptyInput = errors.New("empty input")

	// ErrNilContent indicates a nil DiffContent was passed.
	ErrNilContent = errors.New("nil content")
)
