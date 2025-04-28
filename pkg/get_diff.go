package pkg

import (
	"strings"
)

type Move struct {
	from int
	to   int
}

func GetDiff(oldStr string, newStr string) ([]int, []int, []Move, string, string, error) {
	normalizedOld, normalizedNew, err := normalizeLines(oldStr, newStr)

	if err != nil {
		return []int{}, []int{}, []Move{}, "", "", err
	}

	oldArr, newArr := strings.Split(normalizedOld, "\n"), strings.Split(normalizedNew, "\n")
	var additions []int
	var deletions []int
	var moves []Move

	oldMap := make(map[string]int)
	newMap := make(map[string]int)

	for i, l := range oldArr {
		oldMap[l] = i
	}

	for i, l := range newArr {
		newMap[l] = i
	}

	for x, n := range newArr {
		if y, exists := oldMap[n]; exists && (oldMap[n] != newMap[n]) {
			moves = append(moves, Move{from: y, to: x})
		} else if _, exists := oldMap[n]; exists && (oldMap[n] == newMap[n]) {
			continue
		} else {
			additions = append(additions, x)
		}
	}

	for y, o := range oldArr {
		if _, exists := newMap[o]; !exists {
			deletions = append(deletions, y)
		}
	}

	return additions, deletions, moves, normalizedOld, normalizedNew, nil
}
