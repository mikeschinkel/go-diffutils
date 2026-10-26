package sgdiff

import (
	"strings"

	"github.com/mikeschinkel/go-diffutils"
	"github.com/sourcegraph/go-diff/diff"
)

// ParseUnifiedDiff parses unified diff text (standard or Git extended)
// into DiffContent slices. Since unified diffs contain only hunks (not full
// file content), OldLines and NewLines are reconstructed from hunk data.
func ParseUnifiedDiff(text string) (contents []*diffutils.DiffContent, err error) {
	var fileDiffs []*diff.FileDiff

	fileDiffs, err = diff.ParseMultiFileDiff([]byte(text))
	if err != nil {
		err = NewErr(ErrParseDiff, "text_len", len(text), err)
		goto end
	}

	for _, fd := range fileDiffs {
		var dc *diffutils.DiffContent
		dc, err = fileDiffToContent(fd)
		if err != nil {
			goto end
		}
		contents = append(contents, dc)
	}

end:
	return contents, err
}

// ParseUnifiedDiffWithContent overlays parsed hunks onto full file content,
// producing a DiffContent with complete OldLines and NewLines.
func ParseUnifiedDiffWithContent(diffText string, oldFull, newFull []string) (content *diffutils.DiffContent, err error) {
	var fileDiffs []*diff.FileDiff

	fileDiffs, err = diff.ParseMultiFileDiff([]byte(diffText))
	if err != nil {
		err = NewErr(ErrParseDiff, "text_len", len(diffText), err)
		goto end
	}

	if len(fileDiffs) == 0 {
		content = &diffutils.DiffContent{
			OldLines: oldFull,
			NewLines: newFull,
		}
		goto end
	}

	content = &diffutils.DiffContent{
		OldLines: oldFull,
		NewLines: newFull,
		Label:    extractLabel(fileDiffs[0]),
	}

	for _, hunk := range fileDiffs[0].Hunks {
		content.Changes = append(content.Changes, hunkToChanges(hunk)...)
	}

end:
	return content, err
}

// fileDiffToContent converts a single parsed FileDiff into DiffContent.
// Since unified diffs only contain hunks, we reconstruct partial line content.
func fileDiffToContent(fd *diff.FileDiff) (content *diffutils.DiffContent, err error) {
	content = &diffutils.DiffContent{
		Label: extractLabel(fd),
	}

	// Reconstruct lines from hunks
	var oldLines, newLines []string
	for _, hunk := range fd.Hunks {
		bodyLines := strings.Split(string(hunk.Body), "\n")
		for _, line := range bodyLines {
			if len(line) == 0 {
				continue
			}
			switch line[0] {
			case ' ':
				oldLines = append(oldLines, line[1:])
				newLines = append(newLines, line[1:])
			case '-':
				oldLines = append(oldLines, line[1:])
			case '+':
				newLines = append(newLines, line[1:])
			}
		}
		content.Changes = append(content.Changes, hunkToChanges(hunk)...)
	}

	content.OldLines = oldLines
	content.NewLines = newLines

	return content, nil
}

// hunkToChanges converts a single hunk into DiffChange entries.
func hunkToChanges(hunk *diff.Hunk) []diffutils.DiffChange {
	var changes []diffutils.DiffChange
	bodyLines := strings.Split(string(hunk.Body), "\n")

	oldLine := int(hunk.OrigStartLine)
	newLine := int(hunk.NewStartLine)

	i := 0
	for i < len(bodyLines) {
		line := bodyLines[i]
		if len(line) == 0 {
			i++
			continue
		}

		switch line[0] {
		case ' ':
			// Context line — advance both
			oldLine++
			newLine++
			i++
		case '-':
			// Count contiguous deletions
			start := oldLine
			count := 0
			for i < len(bodyLines) && len(bodyLines[i]) > 0 && bodyLines[i][0] == '-' {
				count++
				oldLine++
				i++
			}
			changes = append(changes, diffutils.DiffChange{
				Type: diffutils.LinesDeleted,
				OldRange: diffutils.LineRange{
					Start: start,
					Count: count,
				},
				NewRange: diffutils.LineRange{
					Start: newLine,
					Count: 0,
				},
			})
		case '+':
			// Count contiguous additions
			start := newLine
			count := 0
			for i < len(bodyLines) && len(bodyLines[i]) > 0 && bodyLines[i][0] == '+' {
				count++
				newLine++
				i++
			}
			changes = append(changes, diffutils.DiffChange{
				Type: diffutils.LinesAdded,
				OldRange: diffutils.LineRange{
					Start: oldLine,
					Count: 0,
				},
				NewRange: diffutils.LineRange{
					Start: start,
					Count: count,
				},
			})
		default:
			i++
		}
	}

	return changes
}

// extractLabel returns a display label from a FileDiff.
func extractLabel(fd *diff.FileDiff) string {
	if fd.NewName != "" && fd.NewName != "/dev/null" {
		return strings.TrimPrefix(fd.NewName, "b/")
	}
	if fd.OrigName != "" && fd.OrigName != "/dev/null" {
		return strings.TrimPrefix(fd.OrigName, "a/")
	}
	return ""
}
