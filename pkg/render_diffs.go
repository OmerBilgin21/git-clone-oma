package pkg

import (
	"errors"
	"fmt"
	"strings"
)

func isVisuallyEmpty(s string) bool {
	s = stripANSI(s)
	s = strings.ReplaceAll(s, "\t", "")
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\r", "")

	return strings.TrimSpace(s) == ""
}

// I couldn't find a nice table rendering display tool for large text files
// (or couldn't make tablewriter work in other words)
// so here we go:
func renderSideBySideDiff(oldColoured, newColoured, oldName, newName string, headerWidth int) {
	oldLines, newLines := strings.Split(oldColoured, "\n"), strings.Split(newColoured, "\n")
	maxLines := max(len(oldLines), len(newLines))

	fmt.Printf("%-*s | %-*s", headerWidth, oldName, headerWidth, newName)
	columnSeparator := strings.Repeat("-", headerWidth*2)
	fmt.Printf(columnSeparator + "\n")

	for i := range maxLines {
		oldLine, newLine := oldLines[i], newLines[i]

		if isVisuallyEmpty(oldLine) && isVisuallyEmpty(newLine) {
			continue
		}

		// if strLen(oldLine) != headerWidth || strLen(newLine) != headerWidth {
		// 	fmt.Printf("MISMATCH @ line %d: old=%d, new=%d\n", i, strLen(oldLine), strLen(newLine))
		// }

		// FIXME: there's a bug with every line having an extra EOF when visualizing right now
		// not due to actual files having those or the newline below -dunno why-
		fmt.Printf("%s | %s\n", oldLine, newLine)
		// fmt.Printf("%q | %q\n", oldLine, newLine)
	}
}

func RenderDiffs(oldContent, newContent, oldName, newName string) {
	additions, deletions := GetDiffs(oldContent, newContent)
	if len(additions) > 0 || len(deletions) > 0 {
		var headerWidth = 50
		normalizedOld, normalizedNew, err := normalizeLines(oldContent, newContent, headerWidth)

		if err != nil {
			check(errors.Join(err, errors.New("diff view would be broken therefore it won't be shown for this file")), false)
			return
		}

		oldColoured, newColoured := ColourTheDiffs(additions, deletions, normalizedOld, normalizedNew)
		renderSideBySideDiff(oldColoured, newColoured, oldName, newName, headerWidth)
	}
}
