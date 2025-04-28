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

func consolidateShortLine(line string, Width int) string {
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

func onInitial(oldText, newText string) {

	oldText = strings.ReplaceAll(oldText, "\r", "")
	newText = strings.ReplaceAll(newText, "\r", "")

	oldLines := strings.Split(oldText, "\n")
	newLines := strings.Split(newText, "\n")

	maxLines := max(len(oldLines), len(newLines))

	// make them 1:1 matches for lines by appending the shorter one with empty strings
	for len(oldLines) < maxLines {
		oldLines = append(oldLines, consolidateShortLine("", Width))
	}
	for len(newLines) < maxLines {
		newLines = append(newLines, consolidateShortLine("", Width))
	}
}

// TODO: new recursive normalizer
func newTrial(oldArr, newArr []string, idx int, err error) ([]string, []string, int, error, bool) {
	if err != nil || (idx > len(newArr) && idx > len(oldArr)) {
		return oldArr, newArr, idx, err, true
	}

	if idx == 0 {
		onInitial(strings.Join(oldArr, "\n"), strings.Join(newArr, "\n"))
	}

	if idx < len(oldArr) {
		line := oldArr[idx]
		head := line[:Width]
		tail := line[Width:]
		head = consolidateShortLine(head, Width)
		tail = consolidateShortLine(tail, Width)

		// *side = slices.Delete(*side, i, i+1)
		// *side = slices.Insert(*side, i, head, tail)
		//
		// // preserve the overall string length after splitting
		// if side == &oldLines {
		// 	newLines = slices.Insert(newLines, i, consolidateShortLine("", Width))
		// } else if side == &newLines {
		// 	oldLines = slices.Insert(oldLines, i, consolidateShortLine("", Width))
		// } else {
		// 	return "", "", fmt.Errorf("string length consolidation failed")
		// }

		if strLen(head) <= Width && strLen(tail) <= Width {
			oldArr = slices.Delete(oldArr, idx, idx+1)
			oldArr = slices.Insert(oldArr, idx, head, tail)
			idx++
			newOld, newNew, index, err, shouldReturn := newTrial(oldArr, newArr, idx, err)
			if shouldReturn {
				return newOld, newNew, index, err, true
			}
		}

	}

	return oldArr, newArr, idx, err, false
}

func normalizeLines(oldText, newText string) (string, string, error) {
	oldText = strings.ReplaceAll(oldText, "\r", "")
	newText = strings.ReplaceAll(newText, "\r", "")

	oldLines := strings.Split(oldText, "\n")
	newLines := strings.Split(newText, "\n")

	maxLines := max(len(oldLines), len(newLines))

	// make them 1:1 matches for lines by appending the shorter one with empty strings
	for len(oldLines) < maxLines {
		oldLines = append(oldLines, consolidateShortLine("", Width))
	}
	for len(newLines) < maxLines {
		newLines = append(newLines, consolidateShortLine("", Width))
	}

	newMaxLines := max(len(oldLines), len(newLines))

	for i := 0; i < newMaxLines; {
		var (
			side *[]string
			line string
		)

		if strLen(oldLines[i]) > Width {
			side = &oldLines
			line = oldLines[i]
		} else if strLen(newLines[i]) > Width {
			side = &newLines
			line = newLines[i]
		} else {
			newLines[i] = consolidateShortLine(newLines[i], Width)
			oldLines[i] = consolidateShortLine(oldLines[i], Width)

			i++
			continue
		}

		head := line[:Width]
		tail := line[Width:]
		head = consolidateShortLine(head, Width)
		tail = consolidateShortLine(tail, Width)

		*side = slices.Delete(*side, i, i+1)
		*side = slices.Insert(*side, i, head, tail)

		// preserve the overall string length after splitting
		if side == &oldLines {
			newLines = slices.Insert(newLines, i, consolidateShortLine("", Width))
		} else if side == &newLines {
			oldLines = slices.Insert(oldLines, i, consolidateShortLine("", Width))
		} else {
			return "", "", fmt.Errorf("string length consolidation failed")
		}

		if i+1 < len(newLines) && i+1 < len(oldLines) {
			oldLines[i] = consolidateShortLine(oldLines[i], Width)
			oldLines[i+1] = consolidateShortLine(oldLines[i+1], Width)
			newLines[i] = consolidateShortLine(newLines[i], Width)
			newLines[i+1] = consolidateShortLine(newLines[i+1], Width)
		}
	}

	return strings.Join(oldLines, "\n"), strings.Join(newLines, "\n"), nil
}
