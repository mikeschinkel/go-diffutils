package diffutils

// DiffLines computes a diff between two line slices.
// Uses the absorbed go-difflib SequenceMatcher engine (Myers-like algorithm).
// Returns a DiffContent with the original lines and computed change blocks.
func DiffLines(old, new []string, label string) (content *DiffContent, err error) {
	content = &DiffContent{
		OldLines: old,
		NewLines: new,
		Label:    label,
	}

	if len(old) == 0 && len(new) == 0 {
		goto end
	}

	content.Changes, err = computeChanges(old, new)

end:
	return content, err
}

// computeChanges uses SequenceMatcher to find change blocks between two slices.
func computeChanges(old, new []string) (changes []DiffChange, err error) {
	matcher := NewMatcher(old, new)
	opCodes := matcher.GetOpCodes()

	for _, op := range opCodes {
		switch op.Tag {
		case 'r':
			// Replace: old lines deleted, new lines added
			changes = append(changes, DiffChange{
				Type: LinesDeleted,
				OldRange: LineRange{
					Start: op.I1 + 1,
					Count: op.I2 - op.I1,
				},
				NewRange: LineRange{
					Start: op.J1 + 1,
					Count: 0,
				},
			})
			changes = append(changes, DiffChange{
				Type: LinesAdded,
				OldRange: LineRange{
					Start: op.I2 + 1,
					Count: 0,
				},
				NewRange: LineRange{
					Start: op.J1 + 1,
					Count: op.J2 - op.J1,
				},
			})
		case 'd':
			// Delete: lines removed from old
			changes = append(changes, DiffChange{
				Type: LinesDeleted,
				OldRange: LineRange{
					Start: op.I1 + 1,
					Count: op.I2 - op.I1,
				},
				NewRange: LineRange{
					Start: op.J1 + 1,
					Count: 0,
				},
			})
		case 'i':
			// Insert: lines added to new
			changes = append(changes, DiffChange{
				Type: LinesAdded,
				OldRange: LineRange{
					Start: op.I1 + 1,
					Count: 0,
				},
				NewRange: LineRange{
					Start: op.J1 + 1,
					Count: op.J2 - op.J1,
				},
			})
		case 'e':
			// Equal: no change
		}
	}

	return changes, nil
}
