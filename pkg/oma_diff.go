package pkg

import (
	"context"
	"database/sql"
	"oma/internal/storage"
)

func GitDiff(ctx context.Context, repoContainer *storage.RepositoryContainer, fileIngredients *[]FileIngredients) error {
	for _, ingredient := range *fileIngredients {
		newres, err := repoContainer.OmaRepository.GetLatestByFileName(ctx, sql.NullString{
			String: ingredient.fileName,
			Valid:  true,
		})

		if err != nil {
			return err
		}

		if newres.ID == 0 {
			repoContainer.OmaRepository.Create(ctx, &storage.OmaRepository{
				FileName: sql.NullString{
					String: ingredient.fileName,
					Valid:  true,
				},
				CachedText: sql.NullString{
					Valid:  true,
					String: ingredient.content,
				},
			})
		} else {
			RenderDiffs(newres.CachedText.String, ingredient.content, newres.FileName.String, ingredient.fileName)
		}
	}

	return nil
}
