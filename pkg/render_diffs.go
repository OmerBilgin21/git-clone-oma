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
// (or couldn't make tablewriter work in other words :D)
// so here we go:
func renderSideBySideDiff(oldColoured, newColoured, filename string) {
	separator := Reset + " | " + Reset
	oldLines, newLines := strings.Split(oldColoured, "\n"), strings.Split(newColoured, "\n")
	maxLines := max(len(oldLines), len(newLines))

	// FIXME: don't do this, it'll be confusing
	// if there are files with same names within the same project
	filenameBasePath := filepath.Base(filename)
	headerName := consolidateShortLine(filenameBasePath)

	// TODO: break down headers into Width length string slices as well,
	// maybe they'll get too long due to their absolute path
	fmt.Printf("\n%s%s%s\n", headerName, separator, headerName)

	columnSeparator := strings.Repeat("-", Width*2)
	fmt.Printf(columnSeparator + "\n")

	for i := range maxLines {
		oldLine, newLine := oldLines[i], newLines[i]

		// TODO: check if rendering an empty red line for emty oldLine
		// and a green line for empty newLine makes sense or not?
		// I don't remember how well this works for just a line addition or removal
		if isVisuallyEmpty(oldLine) && isVisuallyEmpty(newLine) {
			continue
		}

		fmt.Printf("%s%s%s\n", oldLine, separator, newLine)
	}
}

func RenderDiffs(oldContent, newContent, filename string) error {
	diffResult := GetDiff(oldContent, newContent, true)

	if diffResult.error != nil {
		return fmt.Errorf("diff view would be broken therefore it won't be shown for this file: %v\n", diffResult.error)
	}

	if len(diffResult.additions) > 0 || len(diffResult.deletions) > 0 || len(diffResult.moves) > 0 {
		oldColoured, newColoured := ColourTheDiffs(diffResult.additions, diffResult.deletions, diffResult.moves, diffResult.normalizedOld, diffResult.normalizedNew)
		renderSideBySideDiff(oldColoured, newColoured, filename)
	}

	return nil
}
