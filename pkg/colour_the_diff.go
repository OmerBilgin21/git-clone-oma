package pkg

import (
	"oma/internal/storage"
	"strings"
)

var Red = "\033[31m"
var Green = "\033[32m"
var Reset = "\033[0m"

func ColourTheDiffs(
	actions []Action,
	oldStr string,
	newStr string,
) (string, string) {
	oldArr, newArr := strings.Split(oldStr, "\n"), strings.Split(newStr, "\n")

	for _, x := range actions {
		if x.ActionType == storage.AdditionKey {
			if x.Pos < len(newArr) {
				newArr[x.Pos] = Green + newArr[x.Pos] + Reset
			}
		} else {
			if x.Pos < len(oldArr) {
				oldArr[x.Pos] = Red + oldArr[x.Pos] + Reset
			}
		}
	}

	return strings.Join(oldArr, "\n"), strings.Join(newArr, "\n")
}
