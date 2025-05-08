package pkg

import (
	"oma/internal/storage"
	"slices"
	"strings"
)

func RebuildDiff(oldArr []string, versionActions []storage.VersionActions, newVersion *string) bool {
	if len(versionActions) > 0 {
		action := versionActions[0]

		if action.ActionKey == storage.DeletionKey {
			oldArr = slices.Delete(oldArr, action.Pos, action.Pos+1)
		} else if action.ActionKey == storage.AdditionKey {
			oldArr = slices.Insert(oldArr, action.Pos, action.Content)
		}

		remainders := slices.Delete(versionActions, 0, 1)
		*newVersion = strings.Join(oldArr, "\n")
		if res := RebuildDiff(oldArr, remainders, newVersion); res {
			return true
		}
	}

	*newVersion = strings.Join(oldArr, "\n")
	return true
}
