package pkg

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type FileIngredients struct {
	fileName string
	content  string
}

func WalkDirs(curr string, fileIngredientsPtr *[]FileIngredients, processedSteps []string, ignoreList []string) bool {
	someList, err := os.ReadDir(curr)
	check(err, true)

	rootDirPath, err := os.Getwd()
	check(err, true)
	rootDir := filepath.Base(rootDirPath)
	currDirName := filepath.Base(curr)

	if rootDir != currDirName && slices.Contains(processedSteps, currDirName) {
		fmt.Printf("we have come to an end node on our directory tree and it is not the root node, going one step back!\n")
		processedSteps = append(processedSteps, currDirName)
		if WalkDirs(filepath.Join(curr, "../"), fileIngredientsPtr, processedSteps, ignoreList) {
			return true
		}
	} else if currDirName == rootDir && len(processedSteps) > 0 {
		var foundElems = 0
		for _, entry := range someList {
			compareElem := strings.Trim(entry.Name(), " ")

			if slices.Contains(processedSteps, compareElem) || slices.Contains(ignoreList, compareElem) {
				foundElems += 1
			}
		}

		if foundElems == len(someList) {
			return true
		}
	}

	for _, idk := range someList {
		if slices.Contains(ignoreList, idk.Name()) {
			continue
		}

		if slices.Contains(processedSteps, idk.Name()) {
			continue
		}

		if idk.IsDir() {
			fmt.Printf("We faced a directory that we have not been inside: %v, going in...\n", idk.Name())

			curr = filepath.Join(curr, idk.Name())
			if WalkDirs(curr, fileIngredientsPtr, processedSteps, ignoreList) {
				return true
			}
		}

		fileNameToProcess := filepath.Join(curr, idk.Name())
		contentBytes, err := os.ReadFile(fileNameToProcess)
		check(err, false)

		content := string(contentBytes)

		*fileIngredientsPtr = append(*fileIngredientsPtr, FileIngredients{
			fileName: fileNameToProcess,
			content:  content,
		})

		processedSteps = append(processedSteps, idk.Name())
	}

	if currDirName != rootDir {
		fmt.Printf("Everything has been processed here: %v, and it's not the root directory, going back...\n", currDirName)
		processedSteps = append(processedSteps, currDirName)
		if WalkDirs(filepath.Join(curr, "../"), fileIngredientsPtr, processedSteps, ignoreList) {
			return true
		}
	}

	return true
}
