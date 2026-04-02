package main

import (
	"oma/internal"
	"os"
)

func main() {
	for _, target := range []string{".oma"} {
		if _, err := os.Stat(target); err == nil {
			if err := os.RemoveAll(target); err != nil {
				internal.Logger("error while removing ", target, err)
			}
			internal.Logger("removed", target)
		}
	}

	internal.Logger("reset complete \n")
}
