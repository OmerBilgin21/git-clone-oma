package pkg

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"oma/internal/storage"
)

func GitCommit(ctx context.Context, repoContainer *storage.RepositoryContainer, fileIngredients *[]FileIngredients) error {
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

		diffResult := GetDiff(foundRepo.CachedText.String, ingredient.content)

		if diffResult.error != nil {
			return err
		}

		if len(diffResult.additions) == 0 && len(diffResult.deletions) == 0 && len(diffResult.moves) == 0 {
			continue
		}

		fmt.Printf("changes detected, committing file: %v\n", ingredient.fileName)

		newVersion, err := repoContainer.VersionsRepository.Create(ctx, &storage.Versions{
			RepositoryId: foundRepo.ID,
		})

		if err != nil {
			return err
		}

		for _, addition := range diffResult.additions {
			_, err := repoContainer.VersionActionsRepository.Create(ctx, &storage.VersionActions{
				Dest:      addition.position,
				ActionKey: storage.AdditionKey,
				VersionId: newVersion.ID,
				Content:   addition.content,
			})

			if err != nil {
				return err
			}

		}

		for _, deletion := range diffResult.deletions {
			_, err := repoContainer.VersionActionsRepository.Create(ctx, &storage.VersionActions{
				Dest:      deletion.position,
				ActionKey: storage.DeletionKey,
				VersionId: newVersion.ID,
				Content:   deletion.content,
			})

			if err != nil {
				return err
			}

		}

		for _, move := range diffResult.moves {
			_, err := repoContainer.VersionActionsRepository.Create(ctx, &storage.VersionActions{
				Start:     sql.Null[int]{V: move.from, Valid: true},
				Dest:      move.to,
				Content:   move.content,
				ActionKey: storage.MoveKey,
				VersionId: newVersion.ID,
			})

			if err != nil {
				return err
			}
		}
	}

	log.Printf("diff committed successfully")
	return nil
}
