package pkg

import (
	"fmt"
	"os"
	"slices"
)

var OmaIgnoreDefaults = []string{".git", ".oma", ".omaignore", ".gitignore"}

func ParseAndDispatch(args []string) {
	if slices.Contains(args, "init") {

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
