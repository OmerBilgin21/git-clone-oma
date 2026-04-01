package pkg

import (
	"context"
	"oma/internal"
	"strings"
)

func (d *DispatchCommand) GitDiff(ctx context.Context) error {
	repoId, err := d.fileIO.GetRepositoryId()

	if err != nil {
		return err
	}

	for _, ingredient := range d.fileIngredients {
		foundRepo, err := d.omaRepo.GetByFilename(ctx, ingredient.FileName, repoId)

		if err != nil {
			return err
		}

		if foundRepo.ID == 0 {
			continue
		} else {
			versionActions, err := d.GetAllVersionActionsForRepo(ctx, foundRepo.ID)
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
