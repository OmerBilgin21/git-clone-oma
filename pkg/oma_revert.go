package pkg

import (
	"context"
	"fmt"
	"oma/internal"
	"strconv"
	"strings"
)

func (d *OmaVC) getDiffOfEverythingAgainstCurrentState(ctx context.Context) (int, int, error) {

	repoId, err := d.fileIO.GetRepositoryId()

	if err != nil {
		return 0, 0, err
	}

	existingCommitted := 0
	newCommitted := 0

	for _, ingredient := range d.fileIngredients {
		foundRepo, err := d.omaRepo.GetByFilename(ctx, ingredient.FileName, *repoId)

		if err != nil {
			return 0, 0, fmt.Errorf("error while finding a repository for file: %v\nerror:\n%w\n", ingredient.FileName, err)
		}

		if foundRepo.ID == 0 {
			newCommitted++
			continue
		}

		versionActions, err := d.GetAllVersionActionsForRepo(ctx, foundRepo.ID)

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
	}

	return existingCommitted, newCommitted, nil
}

func (d *OmaVC) OmaRevert(ctx context.Context, backFlag internal.Flag) error {
	repoId, err := d.fileIO.GetRepositoryId()

	if err != nil {
		return err
	}

	backAmount, err := strconv.Atoi(backFlag.Value)

	if err != nil {
		return fmt.Errorf("back flag's value must be an integer")
	}

	existing, newFiles, err := d.getDiffOfEverythingAgainstCurrentState(ctx)

	if err != nil {
		return fmt.Errorf("there was an error while diffing the current state of your repository to Oma's last recorded state")
	}

	if existing > 0 || newFiles > 0 {
		return fmt.Errorf("there are %v new and/or %v existing files with modifications, please commit them first", newFiles, existing)
	}

	maxAcrossAll, err := d.omaRepo.GetMaxVersionCountByOmaRepoId(ctx, *repoId)
	if err != nil {
		return fmt.Errorf("error while retrieving max version count:\n%w", err)
	}

	if backAmount > maxAcrossAll {
		return fmt.Errorf("cannot revert %v commit(s), maximum amount of versions in this repository for a file is %v", backAmount, maxAcrossAll)
	}

	for _, file := range d.fileIngredients {
		repository, err := d.omaRepo.GetByFilename(ctx, file.FileName, *repoId)

		if err != nil {
			return fmt.Errorf("error while getting repository: %v\nerror:\n%w", file.FileName, err)
		}

		// did not exist before
		if repository.ID == 0 {
			// filename is an absolute path so this should work?
			if err := d.fileIO.DeleteFile(file.FileName); err != nil {
				return fmt.Errorf("file %v did not exist %v commits ago, however, the attempt of deleting it was not successful", file.FileName, backAmount)
			}
			err = d.omaRepo.Delete(ctx, repository.ID)
			if err != nil {
				return fmt.Errorf("error while deleting the file entry:\n%w", err)
			}
			continue
		}

		maxVersion, _ := d.versionsRepo.GetMaxVersionNumberForRepo(ctx, repository.ID)

		// means that this file did not have that many commits yet
		if maxVersion < backAmount {
			continue
		}

		versions, err := d.versionsRepo.GetLatestXByRepoId(ctx, repository.ID, backAmount)

		if err != nil {
			return fmt.Errorf("error while retrieving versions:\n%w", err)
		}

		for _, version := range versions {
			err = d.versionsRepo.Delete(ctx, version.ID)
			if err != nil {
				return err
			}

			err = d.versionActionsRepo.DeleteByVersionId(ctx, version.ID)
			if err != nil {
				return err
			}
		}

		versionActions, err := d.GetAllVersionActionsForRepo(ctx, repository.ID)
		if err != nil {
			return err
		}

		oldVersion := strings.Split(*repository.CachedText, "\n")
		var revertedFile string
		internal.RebuildDiff(oldVersion, versionActions, &revertedFile)

		err = d.fileIO.WriteToFile(file.FileName, revertedFile)

		if err != nil {
			return fmt.Errorf("error while writing the reverted file: %v, error:\n%w", file.FileName, err)
		}
	}

	internal.Logger("commits reverted successfully", backAmount)
	return nil
}
