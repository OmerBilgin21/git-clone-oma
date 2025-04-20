package pkg

import (
	"fmt"
	"github.com/mattn/go-runewidth"
	"regexp"
	"slices"
	"strings"
)

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripANSI(input string) string {
	return ansiRegex.ReplaceAllString(input, "")
}

func strLen(str string) int {
	return runewidth.StringWidth(stripANSI(str))
}

func isVisuallyEmpty(s string) bool {
	s = stripANSI(s)
	s = strings.ReplaceAll(s, "\t", "")
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\r", "")

	return strings.TrimSpace(s) == ""
}

func consolidateShortLine(line string, width int) string {
	hasNewL := strings.Contains(line, "\n")
	text := strings.TrimSuffix(line, "\n")

	// 2 space indent for tabs
	text = strings.ReplaceAll(text, "\t", "  ")

	visible := strLen(text)
	if visible == width || visible > width {
		return line
	}

	padding := strings.Repeat(" ", width-visible)
	text = text + padding

	if hasNewL {
		text += "\n"
	}

	return text
}

// I couldn't find a nice table rendering display tool for large text files
// (or couldn't make tablewriter work in other words)
// so here we go:
func renderSideBySideDiff(oldColoured, newColoured string, headerWidth int) {
	oldColoured = strings.ReplaceAll(oldColoured, "\r", "")
	newColoured = strings.ReplaceAll(newColoured, "\r", "")

	oldLines := strings.Split(oldColoured, "\n")
	newLines := strings.Split(newColoured, "\n")

	maxLines := max(len(oldLines), len(newLines))

	// make them 1:1 matches for lines by appending the shorter one with empty strings
	for len(oldLines) < maxLines {
		oldLines = append(oldLines, "")
	}
	for len(newLines) < maxLines {
		newLines = append(newLines, "")
	}

	for i := 0; i < maxLines; {
		var (
			side *[]string
			line string
		)

		if strLen(oldLines[i]) > headerWidth {
			side = &oldLines
			line = oldLines[i]
		} else if strLen(newLines[i]) > headerWidth {
			side = &newLines
			line = newLines[i]
		} else {
			newLines[i] = consolidateShortLine(newLines[i], headerWidth)
			oldLines[i] = consolidateShortLine(oldLines[i], headerWidth)

			i++
			continue
		}

		head := line[:headerWidth]
		tail := line[headerWidth:]

		*side = slices.Delete(*side, i, i+1)
		*side = slices.Insert(*side, i, head, tail)

		if side == &oldLines {
			newLines = slices.Insert(newLines, i, "\n")
		} else if side == &newLines {
			oldLines = slices.Insert(oldLines, i, "\n")
		}

		if i+1 < len(newLines[i]) && i+1 < len(oldLines[i]) {
			oldLines[i] = consolidateShortLine(oldLines[i], headerWidth)
			oldLines[i+1] = consolidateShortLine(oldLines[i+1], headerWidth)
			newLines[i] = consolidateShortLine(newLines[i], headerWidth)
			newLines[i+1] = consolidateShortLine(newLines[i+1], headerWidth)
		}
	}

	columnSeparator := strings.Repeat("-", headerWidth*2)
	fmt.Printf("%-*s | %-*s", headerWidth, "OLD TEXT", headerWidth, "NEW TEXT")
	fmt.Printf(columnSeparator + "\n")

	for i := range maxLines {
		oldLine, newLine := strings.TrimSuffix(oldLines[i], "\n"), strings.TrimSuffix(newLines[i], "\n")
		oldLine, newLine = strings.TrimSuffix(oldLine, "\r"), strings.TrimSuffix(newLine, "\r")

		if isVisuallyEmpty(oldLine) && isVisuallyEmpty(newLine) {
			continue
		}

		if strLen(oldLine) != headerWidth || strLen(newLine) != headerWidth {
			fmt.Printf("MISMATCH @ line %d: old=%d, new=%d\n", i, strLen(oldLine), strLen(newLine))
		}

		// FIXME: there's a bug with every line having an extra EOF when visualizing right now
		// not due to actual files having those or the newline below -dunno why-
		fmt.Printf("%s | %s\n", oldLine, newLine)
		// fmt.Printf("%q | %q\n", oldLine, newLine)
	}

}

func RenderDiffs(oldColoured string, newColoured string) {
	renderSideBySideDiff(oldColoured, newColoured, 75)
}
