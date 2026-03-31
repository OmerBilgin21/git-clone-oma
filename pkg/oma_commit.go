package pkg

import (
	"context"
	"fmt"
	"log"
	"oma/internal"
	"oma/internal/storage"
	"strings"
)

func createCache(repoContainer *storage.RepositoryContainer, ctx context.Context, ingredient internal.FileIngredients, repoId int) error {
	_, err := repoContainer.OmaRepository.Create(ctx, &storage.OmaRepository{
		FileName:   &ingredient.FileName,
		CachedText: &ingredient.Content,
		OmaRepoId:  repoId,
	})

	if err != nil {
		return fmt.Errorf("error while creating a new file cache:\n%w", err)
	}
	return nil
}

func GitCommit(ctx context.Context, repoContainer *storage.RepositoryContainer, fileIngredients *[]internal.FileIngredients, messageFlag internal.Flag) error {
	repoId, err := repoContainer.FileIORepository.GetRepositoryId()

	if err != nil {
		return err
	}

	existingCommitted := 0
	newCommitted := 0

	for _, ingredient := range *fileIngredients {
		foundRepo, err := repoContainer.OmaRepository.GetByFilename(ctx, ingredient.FileName, repoId)

		if err != nil {
			return fmt.Errorf("error while finding a repository for file: %v\nerror:\n%w\n", ingredient.FileName, err)
		}

		if foundRepo.ID == 0 {
			createCache(repoContainer, ctx, ingredient, repoId)
			newCommitted++
			continue
		}

		versionActions, err := getAllVersionActionsForRepo(ctx, repoContainer, foundRepo.ID)

		if err != nil {
			return err
		}

		var rebuilt string
		internal.RebuildDiff(strings.Split(*foundRepo.CachedText, "\n"), versionActions, &rebuilt)

		if rebuilt == ingredient.Content {
			continue
		}

		diffResult := internal.GetDiff(rebuilt, ingredient.Content, false)

		if len(diffResult.Actions) == 0 {
			continue
		}

		existingCommitted++

		newVersion, err := repoContainer.VersionsRepository.Create(ctx, &storage.Versions{
			RepositoryId: foundRepo.ID,
			Message:      messageFlag.Value,
		})

		if err != nil {
			return err
		}

		if err := createActions(ctx, repoContainer, diffResult.Actions, newVersion.ID); err != nil {
			return err
		}
	}

	if existingCommitted == 0 && newCommitted == 0 {
		log.Printf("no change!")
		return nil
	}

	log.Printf("diff committed successfully.\n%v known file(s) and %v new file(s)", existingCommitted, newCommitted)
	return nil
}
