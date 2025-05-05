package pkg

import (
	"fmt"
	"oma/internal/storage"
	"slices"
	"strings"
)

func reconcileVersionActionsAfterApply(remaining []storage.VersionActions, wasDone storage.VersionActions, reverseMode bool) []storage.VersionActions {
	reconciled := remaining
	// TODO: pretty obvious, invert
	if reverseMode {
		fmt.Printf("not implemented yet!\n")
	} else {
		for i, action := range remaining {
			if wasDone.ActionKey == storage.DeletionKey {
				if (action.ActionKey == storage.AdditionKey || action.ActionKey == storage.DeletionKey) && action.Dest >= wasDone.Dest {
					temp := remaining[i]
					remaining = slices.Delete(remaining, i, i+1)
					temp.Dest--
					remaining = slices.Insert(remaining, i, temp)
				}
			} else if wasDone.ActionKey == storage.AdditionKey {
				if (action.ActionKey == storage.AdditionKey || action.ActionKey == storage.DeletionKey) && action.Dest >= wasDone.Dest {
					temp := remaining[i]
					remaining = slices.Delete(remaining, i, i+1)
					temp.Dest++
					remaining = slices.Insert(remaining, i, temp)
				}
			}
		}
	}

	return reconciled
}

func RecursiveRebuildDiff(oldArr []string, versionActions []storage.VersionActions, newVersion *string, revertMode bool) bool {
	// TODO: again, pretty obvious
	if len(versionActions) > 0 && revertMode {
		fmt.Printf("not implemented yet!\n")
	} else if len(versionActions) > 0 && !revertMode {
		action := versionActions[0]

		if action.ActionKey == storage.DeletionKey {
			oldArr = slices.Delete(oldArr, action.Dest, action.Dest+1)
		} else if action.ActionKey == storage.AdditionKey {
			oldArr = slices.Insert(oldArr, action.Dest, action.Content)
		} else {
			oldArr = slices.Delete(oldArr, action.Start.V, action.Start.V+1)
			oldArr = slices.Insert(oldArr, action.Dest, action.Content)
		}

		remainders := slices.Delete(versionActions, 0, 1)
		reconciled := reconcileVersionActionsAfterApply(remainders, action, revertMode)
		*newVersion = strings.Join(oldArr, "\n")
		if res := RecursiveRebuildDiff(oldArr, reconciled, newVersion, revertMode); res {
			return true
		}
	}

	*newVersion = strings.Join(oldArr, "\n")
	return true
}
