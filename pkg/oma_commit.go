package pkg

import (
	"context"
	"fmt"
	"log"
	"oma/internal"
	"oma/internal/storage"
	"strings"
)

func createActions(ctx context.Context, versionActionsRepo *storage.VersionActionsRepository, actions []internal.Action, versionId int) error {
	for _, action := range actions {
		actionToCreate := storage.VersionActions{
			Pos:       action.Pos,
			ActionKey: action.ActionType,
			VersionId: versionId,
			Content:   action.Content,
		}

		_, err := versionActionsRepo.Create(ctx, &actionToCreate)

		if err != nil {
			return err
		}
	}

	return nil
}

func createCache(ctx context.Context, omaRepo *storage.OmaRepositoryImpl, ingredient internal.FileIngredient, repoId int) error {
	_, err := omaRepo.Create(ctx, &storage.OmaRepository{
		FileName:   &ingredient.FileName,
		CachedText: &ingredient.Content,
		OmaRepoId:  repoId,
	})

	if err != nil {
		return fmt.Errorf("error while creating a new file cache:\n%w", err)
	}
	return nil
}

func (d *OmaVC) OmaCommit(ctx context.Context, messageFlag internal.Flag) error {
	repoId, err := d.fileIO.GetRepositoryId()

	if err != nil {
		return err
	}

	existingCommitted := 0
	newCommitted := 0

	for _, ingredient := range d.fileIngredients {
		foundRepo, err := d.omaRepo.GetByFilename(ctx, ingredient.FileName, repoId)

		if err != nil {
			return fmt.Errorf("error while finding a repository for file: %v\nerror:\n%w\n", ingredient.FileName, err)
		}

		if foundRepo.ID == 0 {
			createCache(ctx, d.omaRepo, ingredient, repoId)
			newCommitted++
			continue
		}

		versionActions, err := d.GetAllVersionActionsForRepo(ctx, foundRepo.ID)

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

		newVersion, err := d.versionsRepo.Create(ctx, &storage.Versions{
			RepositoryId: foundRepo.ID,
			Message:      messageFlag.Value,
		})

		if err != nil {
			return err
		}

		if err := createActions(ctx, d.versionActionsRepo, diffResult.Actions, newVersion.ID); err != nil {
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
