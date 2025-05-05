package pkg

import (
	"slices"
	"strings"
)

type Action struct {
	from    int
	to      int
	content string
}

type DiffResult struct {
	additions     []Action
	deletions     []Action
	moves         []Action
	normalizedOld string
	normalizedNew string
	error         error
}

func GetDiff(oldStr string, newStr string, visualMode bool) DiffResult {
	var oldArr []string
	var newArr []string
	normalizedOld, normalizedNew := normalizeLines(oldStr, newStr)

	if visualMode {
		oldArr, newArr = strings.Split(normalizedOld, "\n"), strings.Split(normalizedNew, "\n")
	} else {
		oldArr, newArr = strings.Split(oldStr, "\n"), strings.Split(newStr, "\n")
	}

	var additions []Action
	var deletions []Action
	var moves []Action

	oldMap := make(map[string]int)
	newMap := make(map[string]int)

	for i, l := range oldArr {
		oldMap[l] = i
	}

	for i, l := range newArr {
		newMap[l] = i
	}

	for x, n := range newArr {
		if _, exists := oldMap[n]; !exists {
			additions = append(additions, Action{
				to:      x,
				content: n,
			})
		}
	}

	for y, o := range oldArr {
		if _, exists := newMap[o]; !exists {
			deletions = append(deletions, Action{
				to:      y,
				content: o,
			})
		}
	}

	// in order to properly understand what is really moved, we need a temp version of the old
	// string where we apply the additions and deletions first
	temp := oldArr
	for _, add := range additions {
		temp = slices.Insert(temp, add.to, add.content)
	}

	for _, del := range deletions {
		temp = slices.Delete(temp, del.to, del.to+1)
	}

	tempMap := make(map[string]int)

	for i, l := range temp {
		tempMap[l] = i
	}

	for _, n := range newArr {
		if _, exists := tempMap[n]; exists && (tempMap[n] != newMap[n]) {
			toBeAdded := Action{from: tempMap[n], to: newMap[n], content: n}
			skip := false
			for _, move := range moves {
				if move.from == toBeAdded.to && move.to == toBeAdded.from {
					skip = true
				}
			}
			if !skip {
				moves = append(moves, toBeAdded)
			}
		}
	}

	return DiffResult{
		additions:     additions,
		deletions:     deletions,
		moves:         moves,
		normalizedOld: normalizedOld,
		normalizedNew: normalizedNew,
		error:         nil,
	}
}
