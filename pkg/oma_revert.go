package pkg

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"oma/internal"
	"oma/internal/storage"
	"strconv"
	"strings"
)

func GitRevert(ctx context.Context, repoContainer *storage.RepositoryContainer, fileIngrediends *[]FileIngredients, backFlag internal.Flag) error {
	repoId, err := repoContainer.FileIORepository.GetRepositoryId()

	if err != nil {
		return err
	}

	backAmount, err := strconv.Atoi(backFlag.Value)

	if err != nil {
		return fmt.Errorf("back flag's value must be an integer")
	}

	for _, file := range *fileIngrediends {
		repository, err := repoContainer.OmaRepository.GetByFilename(ctx, sql.NullString{
			Valid:  true,
			String: file.fileName,
		}, repoId)

		if err != nil {
			return fmt.Errorf("error while getting repository: %v\nerror:\n%w", file.fileName, err)
		}

		if repository.ID == 0 {
			// filename is an absolute path so this should work?
			if err := repoContainer.FileIORepository.DeleteFile(file.fileName); err != nil {
				return fmt.Errorf("file %v did not exist %v commits ago, however, the attempt of deleting it was not successful", file.fileName, backAmount)
			}
			err = repoContainer.OmaRepository.Delete(ctx, repository.ID)
			if err != nil {
				return fmt.Errorf("error while deleting the file entry:\n%w", err)
			}
		}

		maxVersion, _ := repoContainer.VersionsRepository.GetMaxVersionNumberForRepo(ctx, repository.ID)

		// means that this file did not have that many commits yet
		if maxVersion < backAmount {
			continue
		}

		versions, err := repoContainer.VersionsRepository.GetLatestXByRepoId(ctx, repository.ID, backAmount)

		if err != nil {
			return fmt.Errorf("error while retrieving versions:\n%w", err)
		}

		for _, version := range versions {
			versionActions, err := repoContainer.VersionActionsRepository.GetByVersionId(ctx, version.ID)

			if err != nil {
				return fmt.Errorf("there are versions defined for this file: %v, but no version actions?\nError:%w", file.fileName, err)
			}

			oldVersion := strings.Split(file.content, "\n")
			var revertedFile string
			RecursiveRebuildDiff(oldVersion, versionActions, &revertedFile, true)

			err = repoContainer.FileIORepository.WriteToFile(file.fileName, revertedFile)

			if err != nil {
				return fmt.Errorf("error while writing the reverted file: %v, error:\n%w", file.fileName, err)
			}

			repoContainer.VersionsRepository.Delete(ctx, version.ID)
		}
	}

	log.Printf("%v commits reverted successfully", backAmount)
	return nil
}
