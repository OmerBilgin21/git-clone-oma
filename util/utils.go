package util

import (
	"slices"
)

func PurifyReadResult(lines []string) []string {
	if len(lines) > 0 {
		for i, elem := range lines {
			if elem == "" || elem == " " || elem == "\n" {
				lines = slices.Delete(lines, i, i+1)
			}
		}
	}

	return lines
}
