package pkg

import (
	"fmt"
	"path/filepath"
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
func renderSideBySideDiff(oldColoured, newColoured, oldName, newName string) {
	separator := Reset + " | " + Reset
	oldLines, newLines := strings.Split(oldColoured, "\n"), strings.Split(newColoured, "\n")
	maxLines := max(len(oldLines), len(newLines))

	fmt.Printf("\n%s%s%s\n", consolidateShortLine(filepath.Base(oldName)), separator, filepath.Base(newName))

	columnSeparator := strings.Repeat("-", Width*2)
	fmt.Printf(columnSeparator + "\n")

	for i := range maxLines {
		oldLine, newLine := oldLines[i], newLines[i]

		if isVisuallyEmpty(oldLine) && isVisuallyEmpty(newLine) {
			continue
		}

		fmt.Printf("%s%s%s\n", oldLine, separator, newLine)
	}
}

func RenderDiffs(oldContent, newContent, oldName, newName string) error {
	additions, deletions, moves, normalizedOld, normalizedNew, err := GetDiff(oldContent, newContent)

	if err != nil {
		return fmt.Errorf("diff view would be broken therefore it won't be shown for this file: %v", err)
	}

	if len(additions) > 0 || len(deletions) > 0 || len(moves) > 0 {
		oldColoured, newColoured := ColourTheDiffs(additions, deletions, moves, normalizedOld, normalizedNew)
		renderSideBySideDiff(oldColoured, newColoured, oldName, newName)
	}

	return nil
}
