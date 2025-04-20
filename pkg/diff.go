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

		} else if x+1 < len(newArr) && newArr[x+1] == oldArr[y] {
			additions = append(additions, models.Coordinate{
				StartX: x,
				StartY: y,
				EndX:   x + 1,
				EndY:   y,
			})
			x++

		} else if y+1 < len(oldArr) && oldArr[y+1] == newArr[x] {
			deletions = append(deletions, models.Coordinate{
				StartX: x,
				StartY: y,
				EndX:   x,
				EndY:   y + 1,
			})
			y++

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
