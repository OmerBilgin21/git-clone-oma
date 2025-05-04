package pkg

import (
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
		if y, exists := oldMap[n]; exists && (oldMap[n] != newMap[n]) {
			moves = append(moves, Action{from: y, to: x, content: n})
		} else if _, exists := oldMap[n]; exists && (oldMap[n] == newMap[n]) {
			continue
		} else {
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

	return DiffResult{
		additions:     additions,
		deletions:     deletions,
		moves:         moves,
		normalizedOld: normalizedOld,
		normalizedNew: normalizedNew,
		error:         nil,
	}
}
