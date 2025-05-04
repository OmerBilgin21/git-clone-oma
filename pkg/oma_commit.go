package pkg

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"oma/internal"
	"oma/internal/storage"
)

func createActions(ctx context.Context, repoContainer *storage.RepositoryContainer, actions []Action, versionId int, key storage.Keys) error {
	for _, action := range actions {
		actionToCreate := storage.VersionActions{
			Dest:      action.to,
			ActionKey: key,
			VersionId: versionId,
			Content:   action.content,
		}

		if key == storage.MoveKey {
			actionToCreate.Start = sql.Null[int]{
				Valid: true,
				V:     action.from,
			}
		}

		_, err := repoContainer.VersionActionsRepository.Create(ctx, &actionToCreate)

		if err != nil {
			return err
		}
	}

	return nil
}

// FIXME: right now, it gets the diff between the current version and the cached version
// and then creates a commit based on that, it should be:
// first get the cached version, get the versions for that file, build the latest version
// and then find the diff between the latest built version and current version
// and then commit those
func GitCommit(ctx context.Context, repoContainer *storage.RepositoryContainer, fileIngredients *[]FileIngredients, messageFlag internal.Flag) error {
	repoId, err := repoContainer.FileIORepository.GetRepositoryId()

	if err != nil {
		return err
	}

	for _, ingredient := range *fileIngredients {
		foundRepo, err := repoContainer.OmaRepository.GetLatestByFileName(ctx, sql.NullString{
			String: ingredient.fileName,
			Valid:  true,
		}, repoId)

		if err != nil && foundRepo.ID == 0 {
			panic(fmt.Errorf("something went very wrong\nlease create an issue on GitHub: https://github.com/OmerBilgin21/git-clone-oma \nerror:\n%w", err))
		} else if err != nil {
			fmt.Printf("No previous version of the file, creating cache...")
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
				return fmt.Errorf("error while creating a new file cache:\n%v", err)
			}

			continue
		}

		diffResult := GetDiff(foundRepo.CachedText.String, ingredient.content, false)

		if diffResult.error != nil {
			return err
		}

		if len(diffResult.additions) == 0 && len(diffResult.deletions) == 0 && len(diffResult.moves) == 0 {
			continue
		}

		fmt.Printf("changes detected, committing file: %v\n", ingredient.fileName)

		newVersion, err := repoContainer.VersionsRepository.Create(ctx, &storage.Versions{
			RepositoryId: foundRepo.ID,
			Message:      messageFlag.Value,
		})

		if err != nil {
			return err
		}

		if err := createActions(ctx, repoContainer, diffResult.additions, newVersion.ID, storage.AdditionKey); err != nil {
			return err
		}

		if err := createActions(ctx, repoContainer, diffResult.deletions, newVersion.ID, storage.DeletionKey); err != nil {
			return err
		}

		if err := createActions(ctx, repoContainer, diffResult.moves, newVersion.ID, storage.MoveKey); err != nil {
			return err
		}
	}

	log.Printf("diff committed successfully")
	return nil
}
