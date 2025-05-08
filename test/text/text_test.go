package text

import (
	"database/sql"
	"oma/internal/storage"
	"oma/pkg"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestText(test *testing.T) {

	testCases := []struct {
		name          string
		oldFile       string
		newFile       string
		expectActions bool
	}{
		{
			name:          "short to short",
			oldFile:       "./shortOld.txt",
			newFile:       "./shortNew.txt",
			expectActions: true,
		},
		{
			name:          "short to long",
			oldFile:       "./shortOld.txt",
			newFile:       "./longNew.txt",
			expectActions: true,
		},
		{
			name:          "long to short",
			oldFile:       "./longOld.txt",
			newFile:       "./shortNew.txt",
			expectActions: true,
		},
		{
			name:          "long to long",
			oldFile:       "./longOld.txt",
			newFile:       "./longNew.txt",
			expectActions: true,
		},
		{
			name:          "empty to short",
			oldFile:       "./emptyOld.txt",
			newFile:       "./shortNew.txt",
			expectActions: true,
		},
		{
			name:          "empty to long",
			oldFile:       "./emptyOld.txt",
			newFile:       "./longNew.txt",
			expectActions: true,
		},
		{
			name:          "short to empty",
			oldFile:       "./shortOld.txt",
			newFile:       "./emptyNew.txt",
			expectActions: true,
		},
		{
			name:          "long to empty",
			oldFile:       "./longOld.txt",
			newFile:       "./emptyNew.txt",
			expectActions: true,
		},
		{
			name:          "weird to weird",
			oldFile:       "./weirdOld.txt",
			newFile:       "./weirdNew.txt",
			expectActions: true,
		},
		{
			name:          "empty to empty",
			oldFile:       "./emptyOld.txt",
			newFile:       "./emptyNew.txt",
			expectActions: false,
		},
		{
			name:          "identical",
			oldFile:       "./identicalOld.txt",
			newFile:       "./identicalNew.txt",
			expectActions: false,
		},
	}

	fileIoRepo := storage.NewFileIO()

	for _, tCase := range testCases {
		test.Run(tCase.name, func(t *testing.T) {

			oldStr, err := fileIoRepo.ReadFile(tCase.oldFile)
			if err != nil {
				t.Logf("test old file couldn't be read")
				t.FailNow()
			}

			newStr, err := fileIoRepo.ReadFile(tCase.newFile)
			if err != nil {
				t.Logf("test new file couldn't be read")
				t.FailNow()
			}

			diffResult := pkg.GetDiff(oldStr, newStr, false)

			if (!tCase.expectActions && len(diffResult.Actions) != 0) || (tCase.expectActions && len(diffResult.Actions) == 0) {
				t.Logf("unexpected action amount: %+v", diffResult.Actions)
				t.FailNow()
			}

			versionActions := []storage.VersionActions{}
			for i, action := range diffResult.Actions {
				versionActions = append(versionActions, storage.VersionActions{
					ID:        i + 1,
					CreatedAt: time.Now(),
					DeletedAt: sql.NullTime{
						Valid: false,
					},
					VersionId: 1,
					ActionKey: action.ActionType,
					Pos:       action.Pos,
					Content:   action.Content,
				})
			}

			var rebuilt string
			pkg.RebuildDiff(strings.Split(oldStr, "\n"), versionActions, &rebuilt)

			if !reflect.DeepEqual(rebuilt, newStr) {
				t.Logf("rebuilt string and newStr are not equal for case: %v", tCase.name)
				t.FailNow()
			}
		})
	}
}
