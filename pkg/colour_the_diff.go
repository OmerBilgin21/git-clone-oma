package pkg

import (
	"fmt"
	"strings"
)

var Red = "\033[31m"
var Green = "\033[32m"
var Reset = "\033[0m"

func ColourTheDiffs(
	additions []Coordinate,
	deletions []Coordinate,
	oldStr string,
	newStr string,
) (string, string) {
	o, n := strings.Split(oldStr, "\n"), strings.Split(newStr, "\n")
	fmt.Printf("start old: %v\n", len(o))
	fmt.Printf("start new: %v\n", len(n))
	oldArr, newArr := strings.Split(oldStr, ""), strings.Split(newStr, "")

	for _, addition := range additions {
		xPos := addition.StartX
		val := newArr[xPos]
		if xPos < len(newArr) {
			if val == "\n" && strLen(val) == 1 {
				newArr[xPos] = Green + "+\n" + Reset
				continue
			}
			newArr[xPos] = Green + val + Reset
		}
	}

	for _, deletion := range deletions {
		yPos := deletion.StartY
		val := oldArr[yPos]
		if yPos < len(oldArr) && strLen(val) == 1 {
			if val == "\n" {
				oldArr[yPos] = Red + "-\n" + Reset
			}
			oldArr[yPos] = Red + val + Reset
		}
	}

	return strings.Join(oldArr, ""), strings.Join(newArr, "")
}
