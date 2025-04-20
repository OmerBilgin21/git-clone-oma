package pkg

import (
	"oma/internal/db/models"
	"strings"
)

var Red = "\033[31m"
var Green = "\033[32m"
var Reset = "\033[0m"

func ColourTheDiffs(
	additions []models.Coordinate,
	deletions []models.Coordinate,
	oldStr string,
	newStr string,
) (string, string) {
	oldArr, newArr := strings.Split(oldStr, ""), strings.Split(newStr, "")

	for _, addition := range additions {
		xPos := addition.StartX
		if xPos < len(newArr) {
			newArr[xPos] = Green + newArr[xPos] + Reset
		}
	}

	for _, deletion := range deletions {
		yPos := deletion.StartY
		if yPos < len(oldArr) {
			oldArr[yPos] = Red + oldArr[yPos] + Reset
		}
	}

	return strings.Join(oldArr, ""), strings.Join(newArr, "")
}
