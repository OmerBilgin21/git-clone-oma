package pkg

import (
	"oma/internal/storage"
	"slices"
	"strings"
)

type Action struct {
	Pos        int
	Content    string
	ActionType storage.Keys
}

type DiffResult struct {
	Actions       []Action
	NormalizedOld string
	NormalizedNew string
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

	temp := slices.Clone(oldArr)
	var actions []Action

	newCtr, tempCtr := 0, 0
	for newCtr < len(newArr) && tempCtr < len(temp) {
		if newArr[newCtr] == temp[tempCtr] {
			newCtr++
			tempCtr++
			continue
		}
		found := -1
		for j := tempCtr + 1; j < len(temp); j++ {
			if temp[j] == newArr[newCtr] {
				found = j
				break
			}
		}
		if found >= 0 {
			// temp[tempCtr:found] are deletions
			for k := tempCtr; k < found; k++ {
				actions = append(actions, Action{Pos: tempCtr, Content: temp[tempCtr], ActionType: storage.DeletionKey})
				temp = slices.Delete(temp, tempCtr, tempCtr+1)
			}
		} else {
			// couldn't find the newStr's current line in the rest of temp, add all of em
			actions = append(actions, Action{Pos: tempCtr, Content: newArr[newCtr], ActionType: storage.AdditionKey})
			temp = slices.Insert(temp, tempCtr, newArr[newCtr])
			newCtr++
			tempCtr++
		}
	}

	// ran out of newArr, delete leftovers
	for tempCtr < len(temp) {
		actions = append(actions, Action{Pos: tempCtr, Content: temp[tempCtr], ActionType: storage.DeletionKey})
		temp = slices.Delete(temp, tempCtr, tempCtr+1)
	}

	// ran out of temp, insert remaining
	for ; newCtr < len(newArr); newCtr++ {
		actions = append(actions, Action{Pos: tempCtr, Content: newArr[newCtr], ActionType: storage.AdditionKey})
		temp = slices.Insert(temp, tempCtr, newArr[newCtr])
		tempCtr++
	}

	return DiffResult{
		Actions:       actions,
		NormalizedOld: normalizedOld,
		NormalizedNew: normalizedNew,
	}
}
