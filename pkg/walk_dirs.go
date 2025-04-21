package pkg

import (
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

	// we have come to an end node on our directory tree and it is not the root node, going one step back
	if rootDir != currDirName && slices.Contains(processedSteps, currDirName) {
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

	for _, fileEntry := range someList {
		if slices.Contains(ignoreList, fileEntry.Name()) {
			continue
		}

		if slices.Contains(processedSteps, fileEntry.Name()) {
			continue
		}
		// We faced a directory that we have not been inside: %v, going in...
		if fileEntry.IsDir() {
			curr = filepath.Join(curr, fileEntry.Name())
			if WalkDirs(curr, fileIngredientsPtr, processedSteps, ignoreList) {
				return true
			}
		}

		fileNameToProcess := filepath.Join(curr, fileEntry.Name())
		contentBytes, err := os.ReadFile(fileNameToProcess)
		check(err, false)

		content := strings.ReplaceAll(string(contentBytes), "\r", "")
		// content = strings.ReplaceAll(content, "\n", "\n")

		content = strings.ReplaceAll(content, "\t", "  ")

		*fileIngredientsPtr = append(*fileIngredientsPtr, FileIngredients{
			fileName: fileNameToProcess,
			content:  content,
		})

		processedSteps = append(processedSteps, fileEntry.Name())
	}

	// Everything has been processed in this dir
	// going back
	if currDirName != rootDir {
		processedSteps = append(processedSteps, currDirName)
		if WalkDirs(filepath.Join(curr, "../"), fileIngredientsPtr, processedSteps, ignoreList) {
			return true
		}
	}

	return true
}
