package pkg

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"oma/internal/storage"
)

func GitCommit(ctx context.Context, repoContainer *storage.RepositoryContainer, fileIngredients *[]FileIngredients) error {
	for _, ingredient := range *fileIngredients {
		newres, err := repoContainer.OmaRepository.GetLatestByFileName(ctx, sql.NullString{
			String: ingredient.fileName,
			Valid:  true,
		})

		if err != nil {
			return err
		}

		if newres.ID == 0 {
			fmt.Printf("No previous version of the file, creating cache...")
			repoContainer.OmaRepository.Create(ctx, &storage.OmaRepository{
				FileName: sql.NullString{
					Valid:  true,
					String: ingredient.fileName,
				},
				CachedText: sql.NullString{
					Valid:  true,
					String: ingredient.content,
				},
			})

		} else {
			additions, deletions := GetDiff(newres.CachedText.String, ingredient.content)

			newVersion, err := repoContainer.VersionsRepository.Create(ctx, &storage.Versions{
				RepositoryId: newres.ID,
			})

			if err != nil {
				return err
			}

			for i := range additions {
				addition := additions[i]
				_, err := repoContainer.VersionActionsRepository.Create(ctx, &storage.VersionActions{
					StartX:    addition.StartX,
					StartY:    addition.StartY,
					EndX:      addition.EndX,
					EndY:      addition.EndY,
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
					StartX:    deletion.StartX,
					StartY:    deletion.StartY,
					EndX:      deletion.EndX,
					EndY:      deletion.EndY,
					ActionKey: storage.DeletionKey,
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
