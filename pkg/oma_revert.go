package pkg

import (
	"context"
	"fmt"
	"log"
	"oma/internal"
	"strconv"
	"strings"
)

func (d *OmaVC) OmaRevert(ctx context.Context, backFlag internal.Flag) error {
	repoId, err := d.fileIO.GetRepositoryId()

	if err != nil {
		return err
	}

	backAmount, err := strconv.Atoi(backFlag.Value)

	if err != nil {
		return fmt.Errorf("back flag's value must be an integer")
	}

	for _, file := range d.fileIngredients {
		repository, err := d.omaRepo.GetByFilename(ctx, file.FileName, repoId)

		if err != nil {
			return fmt.Errorf("error while getting repository: %v\nerror:\n%w", file.FileName, err)
		}

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

		oldVersion := strings.Split(file.Content, "\n")
		var revertedFile string
		internal.RebuildDiff(oldVersion, versionActions, &revertedFile)

		err = d.fileIO.WriteToFile(file.FileName, revertedFile)

		if err != nil {
			return fmt.Errorf("error while writing the reverted file: %v, error:\n%w", file.FileName, err)
		}
	}

	log.Printf("%v commits reverted successfully", backAmount)
	return nil
}
