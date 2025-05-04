package pkg

import (
	"fmt"
	"oma/internal/storage"
	"slices"
	"strings"
)

func RebuildDiff(oldVersion string, versionActions []storage.VersionActions) (string, error) {
	var rebuilt []string
	oldArr := strings.Split(oldVersion, "\n")

	for _, action := range versionActions {
		if action.ActionKey == storage.AdditionKey {
			// don't forget, we'll delete here
			// fmt.Printf("addition: %+v\n", addition)
			if action.Dest > len(oldArr) || action.Dest+1 > len(oldArr) {
				fmt.Printf("len oldArr: %v\n", len(oldArr))
				continue
			}
			fmt.Println("continues")
			rebuilt = slices.Delete(oldArr, action.Dest, action.Dest+1)
		} else if action.ActionKey == storage.DeletionKey {
			// don't forget we'll add here
			// fmt.Printf("deletion: %+v\n", deletion)
			if action.Dest > len(oldArr) || action.Dest+1 > len(oldArr) {
				fmt.Printf("len oldArr: %v\n", len(oldArr))
				continue
			}
			fmt.Println("continues")
			rebuilt = slices.Insert(oldArr, action.Dest, action.Content)
		} else {
			// don't forget to use the from and to inverted here
			// fmt.Printf("move: %+v\n", move)
			if action.Dest > len(oldArr) || action.Dest+1 > len(oldArr) {
				fmt.Printf("len oldArr: %v\n", len(oldArr))
				continue
			}
			if action.Start.V > len(oldArr) || action.Start.V+1 > len(oldArr) {
				fmt.Printf("len oldArr: %v\n", len(oldArr))
				continue
			}
			fmt.Println("continues")
			rebuilt = slices.Delete(oldArr, action.Dest, action.Dest+1)
			rebuilt = slices.Insert(oldArr, action.Start.V, action.Content)
		}
	}

	return strings.Join(rebuilt, "\n"), nil
}
