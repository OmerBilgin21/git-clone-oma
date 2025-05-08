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
// (or couldn't make tablewriter work in other words :D)
// so here we go:
func renderSideBySideDiff(oldColoured, newColoured, filename string) {
	separator := Reset + " | " + Reset
	oldLines, newLines := strings.Split(oldColoured, "\n"), strings.Split(newColoured, "\n")
	maxLines := max(len(oldLines), len(newLines))

	fa := []string{filename}
	recursiveEqualizer(&fa)
	headerName := fa[0]

	fmt.Printf("\n%s%s%s\n", headerName, separator, headerName)

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

func RenderDiffs(oldContent, newContent, filename string) error {
	diffResult := GetDiff(oldContent, newContent, true)

	if len(diffResult.Actions) > 0 {
		oldColoured, newColoured := ColourTheDiffs(diffResult.Actions, diffResult.NormalizedOld, diffResult.NormalizedNew)
		renderSideBySideDiff(oldColoured, newColoured, filename)
	}

	return nil
}
