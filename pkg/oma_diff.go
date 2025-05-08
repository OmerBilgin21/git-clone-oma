package pkg

import (
	"context"
	"database/sql"
	"oma/internal/storage"
	"strings"
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
			versionActions, err := getAllVersionActionsForRepo(ctx, repoContainer, foundRepo.ID)
			if err != nil {
				return err
			}

			var rebuilt string
			RebuildDiff(strings.Split(foundRepo.CachedText.String, "\n"), versionActions, &rebuilt)

			if err := RenderDiffs(rebuilt, ingredient.content, ingredient.fileName); err != nil {
				return err
			}
		}
	}

	return nil
}
