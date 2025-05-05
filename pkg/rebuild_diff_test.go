package pkg

import (
	"database/sql"
	"math/rand"
	"oma/internal/storage"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestRebuildDiff(test *testing.T) {
	testCases := []struct {
		name   string
		oldStr string
		newStr string
	}{
		{
			name:   "add a word and move a word from line",
			oldStr: "hi\nthere\nbrother\nhowdy",
			newStr: "hi\nmy\nbrother\nthere\nhowdy",
		},
		{
			name:   "replace first line",
			oldStr: "hi\nthere\nbrother\nhowdy",
			newStr: "hello\nthere\nbrother\nhowdy",
		},
		{
			name:   "delete first line",
			oldStr: "hi\nthere\nbrother\nhowdy",
			newStr: "there\nbrother\nhowdy",
		},
		{
			name:   "move first line to next",
			oldStr: "hi\nthere\nbrother\nhowdy",
			newStr: "there\nhi\nbrother\nhowdy",
		},
		{
			name:   "no changes",
			oldStr: "hi\nthere\nbrother\nhowdy",
			newStr: "hi\nthere\nbrother\nhowdy",
		},
	}

	for _, tCase := range testCases {
		test.Run(tCase.name, func(t *testing.T) {
			diffResult := GetDiff(tCase.oldStr, tCase.newStr, false)

			if diffResult.error != nil {
				t.Errorf("diff result failed with error:\n%v", diffResult.error)
				t.FailNow()
			}

			var versionActions []storage.VersionActions

			for _, add := range diffResult.additions {
				versionActions = append(versionActions, storage.VersionActions{
					ID:        rand.Intn(100),
					CreatedAt: time.Now(),
					DeletedAt: sql.NullTime{
						Valid: false,
					},
					VersionId: 1,
					ActionKey: storage.AdditionKey,
					Dest:      add.to,
					Content:   add.content,
				})
			}

			for _, del := range diffResult.deletions {
				versionActions = append(versionActions, storage.VersionActions{
					ID:        rand.Intn(100),
					CreatedAt: time.Now(),
					DeletedAt: sql.NullTime{
						Valid: false,
					},
					VersionId: 1,
					ActionKey: storage.DeletionKey,
					Dest:      del.to,
					Content:   del.content,
				})
			}

			for _, move := range diffResult.moves {
				versionActions = append(versionActions, storage.VersionActions{
					ID:        rand.Intn(100),
					CreatedAt: time.Now(),
					DeletedAt: sql.NullTime{
						Valid: false,
					},
					VersionId: 1,
					ActionKey: storage.MoveKey,
					Start: sql.Null[int]{
						Valid: true,
						V:     move.from,
					},
					Dest:    move.to,
					Content: move.content,
				})
			}

			var rebuilt2 string
			oldArr := strings.Split(tCase.oldStr, "\n")
			RecursiveRebuildDiff(oldArr, versionActions, &rebuilt2, false)

			if !reflect.DeepEqual(rebuilt2, tCase.newStr) {
				t.Errorf("rebuilt string and new string are not equal\ngot:\n%v\nexpected:\n%v", rebuilt2, tCase.newStr)
			}

		})
	}
}
