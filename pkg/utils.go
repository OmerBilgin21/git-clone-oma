package pkg

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type FileIngredients struct {
	fileName string
	content  string
}

func PurifyReadResult(lines []string) []string {

	if len(lines) > 0 {
		for i, elem := range lines {
			if elem == "" || elem == " " || elem == "\n" {
				lines = slices.Delete(lines, i, i+1)
			}
		}
	}

	return lines
}

func ParseOmaIgnore() []string {
	omaIgnoreBytes, err := os.ReadFile("./.omaignore")
	check(err, false)
	omaIgnore := string(omaIgnoreBytes)

	separatedArgs := PurifyReadResult(strings.Split(omaIgnore, "\n"))
	separatedArgs = append(separatedArgs, OmaIgnoreDefaults...)

	return separatedArgs
}

func FindIndex() {}

func check(err error, fail bool) {
	if err != nil {
		if fail {
			log.Fatal(err)
		}
		fmt.Printf("err: %v\n", err)
	}
}
