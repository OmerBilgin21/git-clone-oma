package pkg

import (
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

func consolidateShortLine(line string) string {
	hasNewL := strings.Contains(line, "\n")
	text := strings.TrimSuffix(line, "\n")

	// 2 space indent for tabs
	text = strings.ReplaceAll(text, "\t", "  ")

	visible := strLen(text)

	if visible == Width || visible > Width {
		return line
	}

	padding := strings.Repeat(" ", Width-visible)
	text = text + padding

	if hasNewL {
		text += "\n"
	}

	return text
}

func postEqualizerCleanUp(oldText, newText string) (string, string) {
	oldText = strings.ReplaceAll(oldText, "\r", "")
	newText = strings.ReplaceAll(newText, "\r", "")

	oldLines := strings.Split(oldText, "\n")
	newLines := strings.Split(newText, "\n")

	maxLines := max(len(oldLines), len(newLines))

	// make them 1:1 matches for lines by appending the shorter one with empty strings
	for len(oldLines) < maxLines {
		oldLines = append(oldLines, consolidateShortLine(""))
	}
	for len(newLines) < maxLines {
		newLines = append(newLines, consolidateShortLine(""))
	}

	return strings.Join(oldLines, "\n"), strings.Join(newLines, "\n")
}

func recursiveEqualizer(arr *[]string) bool {
	for i, line := range *arr {
		if strLen(line) > Width {

			head := line[:Width]
			tail := line[Width:]

			head = consolidateShortLine(head)
			tail = consolidateShortLine(tail)

			*arr = slices.Delete(*arr, i, i+1)
			*arr = slices.Insert(*arr, i, head, tail)
			if res := recursiveEqualizer(arr); res {
				return true
			}
		} else {
			newLine := consolidateShortLine(line)
			*arr = slices.Delete(*arr, i, i+1)
			*arr = slices.Insert(*arr, i, newLine)
		}
	}

	return true
}

func normalizeLines(oldText, newText string) (string, string) {
	oldArr, newArr := strings.Split(oldText, "\n"), strings.Split(newText, "\n")
	recursiveEqualizer(&oldArr)
	recursiveEqualizer(&newArr)

	return postEqualizerCleanUp(strings.Join(oldArr, "\n"), strings.Join(newArr, "\n"))
}
