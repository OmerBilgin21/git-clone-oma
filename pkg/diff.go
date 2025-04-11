package pkg

import (
	"fmt"
	"strings"
)

type Coordinate struct {
	startX int
	startY int
	endX   int
	endY   int
}

func Diff(oldStr string, newStr string) {
	oldArr, newArr := strings.Split(oldStr, ""), strings.Split(newStr, "")
	var slides []Coordinate
	var additions []Coordinate
	var deletions []Coordinate

	x, y := 0, 0
	for x < len(newArr) && y < len(oldArr) {

		// just slide
		if newArr[x] == oldArr[y] {
			startX, startY := x, y
			for x < len(newArr) && y < len(oldArr) && newArr[x] == oldArr[y] {
				x++
				y++
			}
			slides = append(slides, Coordinate{
				startX: startX,
				startY: startY,
				endX:   x,
				endY:   y,
			})

			// next step of newArr is matching the oldArr's char gotta add ma boi
		} else if x+1 < len(newArr) && newArr[x+1] == oldArr[y] {
			additions = append(additions, Coordinate{
				startX: x,
				startY: y,
				endX:   x + 1,
				endY:   y,
			})
			x++

			// next step of oldArr is matching the newArr's char gotta delete ma boi
		} else if y+1 < len(oldArr) && oldArr[y+1] == newArr[x] {
			deletions = append(deletions, Coordinate{
				startX: x,
				startY: y,
				endX:   x,
				endY:   y + 1,
			})
			y++

			// :fire: this is fine :fire:
		} else {
			additions = append(additions, Coordinate{
				startX: x,
				startY: y,
				endX:   x + 1,
				endY:   y,
			})

			deletions = append(deletions, Coordinate{
				startX: x,
				startY: y,
				endX:   x,
				endY:   y + 1,
			})
			x++
			y++
		}
	}

	fmt.Printf("slides: %+v\n", slides)
	fmt.Printf("additions: %+v\n", additions)
	fmt.Printf("deletions: %+v\n", deletions)
}
