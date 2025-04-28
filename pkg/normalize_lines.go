package pkg

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/mattn/go-runewidth"
)

const Width = 50

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripANSI(input string) string {
	return ansiRegex.ReplaceAllString(input, "")
}

func strLen(str string) int {
	return runewidth.StringWidth(stripANSI(str))
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

func normalizeLines(oldText, newText string, width int) (string, string, error) {
	oldText = strings.ReplaceAll(oldText, "\r", "")
	newText = strings.ReplaceAll(newText, "\r", "")

	oldLines := strings.Split(oldText, "\n")
	newLines := strings.Split(newText, "\n")

	maxLines := max(len(oldLines), len(newLines))

	// make them 1:1 matches for lines by appending the shorter one with empty strings
	for len(oldLines) < maxLines {
		oldLines = append(oldLines, consolidateShortLine("", width))
	}
	for len(newLines) < maxLines {
		newLines = append(newLines, consolidateShortLine("", width))
	}

	for i := 0; i < maxLines; {
		var (
			side *[]string
			line string
		)

		if strLen(oldLines[i]) > width {
			side = &oldLines
			line = oldLines[i]
		} else if strLen(newLines[i]) > width {
			side = &newLines
			line = newLines[i]
		} else {
			newLines[i] = consolidateShortLine(newLines[i], width)
			oldLines[i] = consolidateShortLine(oldLines[i], width)

			i++
			continue
		}

		head := line[:width]
		tail := line[width:]
		head = consolidateShortLine(head, width)
		tail = consolidateShortLine(tail, width)

		*side = slices.Delete(*side, i, i+1)
		*side = slices.Insert(*side, i, head, tail)

		// preserve the overall string length after splitting
		if side == &oldLines {
			newLines = slices.Insert(newLines, i, consolidateShortLine("", width))
		} else if side == &newLines {
			oldLines = slices.Insert(oldLines, i, consolidateShortLine("", width))
		} else {
			return "", "", fmt.Errorf("string length consolidation failed")
		}

		if i+1 < len(newLines) && i+1 < len(oldLines) {
			oldLines[i] = consolidateShortLine(oldLines[i], width)
			oldLines[i+1] = consolidateShortLine(oldLines[i+1], width)
			newLines[i] = consolidateShortLine(newLines[i], width)
			newLines[i+1] = consolidateShortLine(newLines[i+1], width)
		}
	}

	return strings.Join(oldLines, "\n"), strings.Join(newLines, "\n"), nil
}
