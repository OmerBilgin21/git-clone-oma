package pkg

import (
	"context"
	"database/sql"
	"oma/internal/storage"
)

func GitDiff(ctx context.Context, repoContainer *storage.RepositoryContainer, fileIngredients *[]FileIngredients) error {
	repoId, err := repoContainer.FileIORepository.GetRepositoryId()

	if err != nil {
		return err
	}

	for _, ingredient := range *fileIngredients {
		foundRepo, err := repoContainer.OmaRepository.GetByFilename(ctx, sql.NullString{
			String: ingredient.fileName,
			Valid:  true,
		}, repoId)

		if err != nil {
			return err
		}

		if foundRepo.ID == 0 {
			continue
		} else {
			if err := RenderDiffs(foundRepo.CachedText.String, ingredient.content, ingredient.fileName); err != nil {
				return err
			}
		}
	}

	return nil
}
