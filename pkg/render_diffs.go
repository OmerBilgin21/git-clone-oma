package pkg

import (
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

		// FIXME: there's a bug with every line having an extra EOF when visualizing right now
		// not due to actual files having those or the newline below -dunno why-
		fmt.Printf("%s | %s\n", oldLine, newLine)
		// fmt.Printf("%q | %q\n", oldLine, newLine)
	}
}

func RenderDiffs(oldContent, newContent, oldName, newName string) error {
	additions, deletions := GetDiff(oldContent, newContent)

	if len(additions) > 0 || len(deletions) > 0 {
		var headerWidth = 50
		normalizedOld, normalizedNew, err := normalizeLines(oldContent, newContent, headerWidth)

		if err != nil {
			return fmt.Errorf("diff view would be broken therefore it won't be shown for this file: %v", err)
		}

		oldColoured, newColoured := ColourTheDiffs(additions, deletions, normalizedOld, normalizedNew)
		renderSideBySideDiff(oldColoured, newColoured, oldName, newName, headerWidth)
	}

	return nil
}
