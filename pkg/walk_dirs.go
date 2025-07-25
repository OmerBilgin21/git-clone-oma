package pkg

import (
	"oma/internal/storage"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

var OMA_IGNORE_DEFAULTS = []string{".git", ".oma", ".omaignore", ".gitignore", "node_modules"}

type FileIngredients struct {
	fileName string
	content  string
}

// at the time I implemented this I didn't know
// path/filepath.WalkDir existed :D
// but hey, it works, so no need to change it
func WalkDirs(curr string, fileIngredientsPtr *[]FileIngredients, processedSteps []string, ignoreList []string, repoContainer *storage.RepositoryContainer) bool {
	dirIngredientList, err := os.ReadDir(curr)
	check(err, true)

	rootDirPath, err := os.Getwd()
	check(err, true)
	rootDir := filepath.Base(rootDirPath)
	currDirName := filepath.Base(curr)

	// we have come to an end node on our directory tree and it is not the root node, going one step back
	if rootDir != currDirName && slices.Contains(processedSteps, currDirName) {
		processedSteps = append(processedSteps, currDirName)
		if WalkDirs(filepath.Join(curr, "../"), fileIngredientsPtr, processedSteps, ignoreList, repoContainer) {
			return true
		}
	} else if currDirName == rootDir && len(processedSteps) > 0 {
		var foundElems = 0
		for _, entry := range dirIngredientList {
			compareElem := strings.Trim(entry.Name(), " ")

			if slices.Contains(processedSteps, compareElem) || slices.Contains(ignoreList, compareElem) {
				foundElems += 1
			}
		}

		if foundElems == len(dirIngredientList) {
			return true
		}
	}

	for _, fileEntry := range dirIngredientList {
		if slices.Contains(ignoreList, fileEntry.Name()) {
			continue
		}

		if slices.Contains(processedSteps, fileEntry.Name()) {
			continue
		}
		// We faced a directory that we have not been inside: %v, going in...
		if fileEntry.IsDir() {
			curr = filepath.Join(curr, fileEntry.Name())
			if WalkDirs(curr, fileIngredientsPtr, processedSteps, ignoreList, repoContainer) {
				return true
			}
		}

		fileNameToProcess := filepath.Join(curr, fileEntry.Name())
		content, err := repoContainer.FileIORepository.ReadFile(fileNameToProcess)
		check(err, false)

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
		if WalkDirs(filepath.Join(curr, "../"), fileIngredientsPtr, processedSteps, ignoreList, repoContainer) {
			return true
		}
	}

	return true
}
