package pkg

import (
	"strings"
)

var Red = "\033[31m"
var Green = "\033[32m"
var Reset = "\033[0m"
var Orange = "\033[33m"

func ColourTheDiffs(
	additions []AddOrDelete,
	deletions []AddOrDelete,
	moves []Move,
	oldStr string,
	newStr string,
) (string, string) {
	oldArr, newArr := strings.Split(oldStr, "\n"), strings.Split(newStr, "\n")

	for _, x := range additions {
		newArr[x.position] = Green + newArr[x.position] + Reset
	}

	for _, y := range deletions {
		oldArr[y.position] = Red + oldArr[y.position] + Reset
	}

	for _, m := range moves {
		newArr[m.to] = Orange + newArr[m.to] + Reset
	}

	return strings.Join(oldArr, "\n"), strings.Join(newArr, "\n")
}
