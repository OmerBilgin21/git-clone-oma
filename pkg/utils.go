package pkg

import (
	"log"
	"os"
	"slices"
	"strings"
)

func walkDirsAndReadFiles() []FileIngredients {
	currDir, err := os.Getwd()
	check(err, true)

	ignoreList := parseOmaIgnore()

	var fileIngredients []FileIngredients
	WalkDirs(currDir, &fileIngredients, []string{}, ignoreList)

	return fileIngredients
}

func purifyReadResult(lines []string) []string {
	if len(lines) > 0 {
		for i, elem := range lines {
			if elem == "" || elem == " " || elem == "\n" {
				lines = slices.Delete(lines, i, i+1)
			}
		}
	}

	return lines
}

func parseOmaIgnore() []string {
	omaIgnoreBytes, err := os.ReadFile("./.omaignore")
	check(err, false)
	omaIgnore := string(omaIgnoreBytes)

	separatedArgs := purifyReadResult(strings.Split(omaIgnore, "\n"))
	separatedArgs = append(separatedArgs, OMA_IGNORE_DEFAULTS...)

	return separatedArgs
}

func check(err error, fail bool) {
	if err != nil {
		if fail {
			log.Fatal(err)
		}
		log.Printf("err: %v\n", err)
	}
}
