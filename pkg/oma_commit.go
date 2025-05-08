package pkg

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"oma/internal"
	"oma/internal/storage"
	"strings"
)

func createCache(repoContainer *storage.RepositoryContainer, ctx context.Context, ingredient FileIngredients, repoId int) error {
	_, err := repoContainer.OmaRepository.Create(ctx, &storage.OmaRepository{
		FileName: sql.NullString{
			Valid:  true,
			String: ingredient.fileName,
		},
		CachedText: sql.NullString{
			Valid:  true,
			String: ingredient.content,
		},
		OmaRepoId: repoId,
	})

	if err != nil {
		return fmt.Errorf("error while creating a new file cache:\n%w", err)
	}
	return nil
}

func GitCommit(ctx context.Context, repoContainer *storage.RepositoryContainer, fileIngredients *[]FileIngredients, messageFlag internal.Flag) error {
	repoId, err := repoContainer.FileIORepository.GetRepositoryId()

	if err != nil {
		return err
	}

	existingCommitted := 0
	newCommitted := 0

	for _, ingredient := range *fileIngredients {
		foundRepo, err := repoContainer.OmaRepository.GetByFilename(ctx, sql.NullString{
			String: ingredient.fileName,
			Valid:  true,
		}, repoId)

		if err != nil {
			return fmt.Errorf("error while finding a repository for file: %v\nerror:\n%w\n", ingredient.fileName, err)
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
		RebuildDiff(strings.Split(foundRepo.CachedText.String, "\n"), versionActions, &rebuilt)

		if rebuilt == ingredient.content {
			continue
		}

		diffResult := GetDiff(rebuilt, ingredient.content, false)

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
