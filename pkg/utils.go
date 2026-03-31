package pkg

import (
	"context"
	"fmt"
	"oma/internal"
	"oma/internal/storage"
	"os"
	"slices"
	"strings"
)

func createActions(ctx context.Context, repoContainer *storage.RepositoryContainer, actions []internal.Action, versionId int) error {
	for _, action := range actions {
		actionToCreate := storage.VersionActions{
			Pos:       action.Pos,
			ActionKey: action.ActionType,
			VersionId: versionId,
			Content:   action.Content,
		}

		_, err := repoContainer.VersionActionsRepository.Create(ctx, &actionToCreate)

		if err != nil {
			return err
		}
	}

	return nil
}

func getAllVersionActionsForRepo(ctx context.Context, repoContainer *storage.RepositoryContainer, repositoryId int) ([]storage.VersionActions, error) {
	allVersions, err := repoContainer.VersionsRepository.GetAllByRepoId(ctx, repositoryId)

	if err != nil {
		return []storage.VersionActions{}, fmt.Errorf("error while trying to get versions for repository: %v\nerror:\n%w", repositoryId, err)
	}

	var versionActions []storage.VersionActions
	for _, version := range allVersions {
		versionActionsOfVersion, err := repoContainer.VersionActionsRepository.GetByVersionId(ctx, version.ID)
		if err != nil {
			return []storage.VersionActions{}, fmt.Errorf("error while trying to get version actions for version: %v\nerror:\n%w", version.ID, err)
		}
		if len(versionActionsOfVersion) > 0 {
			versionActions = append(versionActions, versionActionsOfVersion...)
		}
	}

	return versionActions, nil
}

func walkDirsAndReadFiles(repoContainer *storage.RepositoryContainer) []internal.FileIngredients {
	currDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	ignoreList := parseOmaIgnore()

	var fileIngredients []internal.FileIngredients
	internal.WalkDirs(currDir, &fileIngredients, []string{}, ignoreList, repoContainer)

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
	if err != nil {
		panic(err)
	}
	omaIgnore := string(omaIgnoreBytes)

	separatedArgs := purifyReadResult(strings.Split(omaIgnore, "\n"))
	separatedArgs = append(separatedArgs, internal.OMA_IGNORE_DEFAULTS...)

	return separatedArgs
}
