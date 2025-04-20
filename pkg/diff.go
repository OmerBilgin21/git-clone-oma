package pkg

import (
	"oma/internal/db/models"
	"strings"
)

func Diff(oldStr string, newStr string) ([]models.Coordinate, []models.Coordinate) {
	oldArr, newArr := strings.Split(oldStr, ""), strings.Split(newStr, "")
	var additions []models.Coordinate
	var deletions []models.Coordinate

	x, y := 0, 0
	for x < len(newArr) && y < len(oldArr) {

		// just slide
		if newArr[x] == oldArr[y] {
			for x < len(newArr) && y < len(oldArr) && newArr[x] == oldArr[y] {
				x++
				y++
			}

			// next step of newArr is matching the oldArr's char gotta add ma boi
		} else if x+1 < len(newArr) && newArr[x+1] == oldArr[y] {
			additions = append(additions, models.Coordinate{
				StartX: x,
				StartY: y,
				EndX:   x + 1,
				EndY:   y,
			})
			x++

			// next step of oldArr is matching the newArr's char gotta delete ma boi
		} else if y+1 < len(oldArr) && oldArr[y+1] == newArr[x] {
			deletions = append(deletions, models.Coordinate{
				StartX: x,
				StartY: y,
				EndX:   x,
				EndY:   y + 1,
			})
			y++

			// :fire: this is fine :fire:
		} else {
			additions = append(additions, models.Coordinate{
				StartX: x,
				StartY: y,
				EndX:   x + 1,
				EndY:   y,
			})

			deletions = append(deletions, models.Coordinate{
				StartX: x,
				StartY: y,
				EndX:   x,
				EndY:   y + 1,
			})
			x++
			y++
		}
	}

	return additions, deletions
}
