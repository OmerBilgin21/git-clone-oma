package pkg

import (
	"reflect"
	"testing"
)

func TestGetDiff(test *testing.T) {
	testCases := []struct {
		name          string
		oldStr        string
		newStr        string
		additions     []Action
		deletion      []Action
		moves         []Action
		normalizedOld string
		normalizedNew string
	}{
		{
			name:      "replace first line",
			oldStr:    "hi\nthere\nbrother\nhowdy",
			newStr:    "hello\nthere\nbrother\nhowdy",
			additions: []Action{{to: 0, content: "hello"}},
			deletion:  []Action{{to: 0, content: "hi"}},
			moves:     []Action{},
		},
		{
			name:      "delete first line",
			oldStr:    "hi\nthere\nbrother\nhowdy",
			newStr:    "there\nbrother\nhowdy",
			additions: []Action{},
			deletion:  []Action{{to: 0, content: "hi"}},
			moves:     []Action{{from: 1, to: 0, content: "there"}, {from: 2, to: 1, content: "brother"}, {from: 3, to: 2, content: "howdy"}},
		},
		{
			name:      "move first line to next",
			oldStr:    "hi\nthere\nbrother\nhowdy",
			newStr:    "there\nhi\nbrother\nhowdy",
			additions: []Action{},
			deletion:  []Action{},
			moves:     []Action{{from: 1, to: 0, content: "there"}, {from: 0, to: 1, content: "hi"}},
		},
		{
			name:      "no changes",
			oldStr:    "hi\nthere\nbrother\nhowdy",
			newStr:    "hi\nthere\nbrother\nhowdy",
			additions: []Action{},
			deletion:  []Action{},
			moves:     []Action{},
		},
	}

	for _, tCase := range testCases {
		test.Run(tCase.name, func(t *testing.T) {
			diffResult := GetDiff(tCase.oldStr, tCase.newStr, false)
			if diffResult.error != nil {
				t.Errorf("diff result failed with error:\n%v", diffResult.error)
				t.FailNow()
			}

			if len(diffResult.additions) != len(tCase.additions) || len(diffResult.deletions) != len(tCase.deletion) || len(diffResult.moves) != len(tCase.moves) {
				t.Errorf("addition, deletion or move count is wrong\nadditions:\n%+v\ndeletions:\n%+v\nmoves:\n%+v", diffResult.additions, diffResult.deletions, diffResult.moves)
				t.FailNow()
			}

			if len(tCase.additions) > 0 && !reflect.DeepEqual(diffResult.additions, tCase.additions) {
				t.Errorf("additions are not equal to expected case\ngot:\n%+v\nexpected:\n%+v", diffResult.additions, tCase.additions)
				t.FailNow()
			}

			if len(tCase.deletion) > 0 && !reflect.DeepEqual(diffResult.deletions, tCase.deletion) {
				t.Errorf("deletions are not equal to expected case\ngot:\n%+v\nexpected:\n%+v", diffResult.deletions, tCase.deletion)
				t.FailNow()
			}

			if len(tCase.moves) > 0 && !reflect.DeepEqual(diffResult.moves, tCase.moves) {
				t.Errorf("moves are not equal to expected case\ngot:\n%+v\nexpected:\n%+v", diffResult.moves, tCase.moves)
				t.FailNow()
			}
		})
	}
}
