package pkg

import (
	"context"
	"oma/internal"
	"oma/internal/storage"
	"strings"
)

func GitDiff(ctx context.Context, repoContainer *storage.RepositoryContainer, fileIngredients *[]internal.FileIngredients) error {
	repoId, err := repoContainer.FileIORepository.GetRepositoryId()

	if err != nil {
		return err
	}

	for _, ingredient := range *fileIngredients {
		foundRepo, err := repoContainer.OmaRepository.GetByFilename(ctx, ingredient.FileName, repoId)

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
			internal.RebuildDiff(strings.Split(*foundRepo.CachedText, "\n"), versionActions, &rebuilt)

			if err := internal.RenderDiffs(rebuilt, ingredient.Content, ingredient.FileName); err != nil {
				return err
			}
		}
	}

	return nil
}
