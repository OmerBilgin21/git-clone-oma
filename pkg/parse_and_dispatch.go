package pkg

import (
	"fmt"
	"log"
	"os"
	"slices"
)

var OMA_IGNORE_DEFAULTS = []string{".git", ".oma", ".omaignore", ".gitignore"}

func ParseAndDispatch(args []string) {
	if slices.Contains(args, "init") {
		if len(args) > 2 {
			log.Fatal("illogical flags/commands type oma init --help for usage")
		} else if len(args) == 2 && args[1] == "--help" {
			log.Fatal("help docs, TBD")
		}

		currDir, err := os.Getwd()
		check(err, true)

		ignoreList := ParseOmaIgnore()

		var fileIngredients []FileIngredients
		WalkDirs(currDir, &fileIngredients, []string{}, ignoreList)

		fmt.Printf("fileIngredients: %+v\n", fileIngredients)
		return
	}

	if slices.Contains(args, "commit") {

	}
}
