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

		if err != nil {
			return err
		}

		if foundRepo.ID == 0 || foundRepo == nil {
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
		} else {
			additions, deletions, moves, _, _, err := GetDiff(foundRepo.CachedText.String, ingredient.content)

			if err != nil {
				return err
			}

			newVersion, err := repoContainer.VersionsRepository.Create(ctx, &storage.Versions{
				RepositoryId: foundRepo.ID,
			})

			if err != nil {
				return err
			}

			for i := range additions {
				addition := additions[i]
				_, err := repoContainer.VersionActionsRepository.Create(ctx, &storage.VersionActions{
					Dest:      addition,
					ActionKey: storage.AdditionKey,
					VersionId: newVersion.ID,
				})

				if err != nil {
					return err
				}

			}

			for i := range deletions {
				deletion := deletions[i]
				_, err := repoContainer.VersionActionsRepository.Create(ctx, &storage.VersionActions{
					Dest:      deletion,
					ActionKey: storage.DeletionKey,
					VersionId: newVersion.ID,
				})

				if err != nil {
					return err
				}

			}

			for i := range moves {
				move := moves[i]

				_, err := repoContainer.VersionActionsRepository.Create(ctx, &storage.VersionActions{
					Start:     sql.Null[int]{V: move.from, Valid: true},
					Dest:      move.to,
					ActionKey: storage.MoveKey,
					VersionId: newVersion.ID,
				})

				if err != nil {
					return err
				}
			}
		}
	}

	log.Printf("diff committed successfully")
	return nil
}
