package pkg

import (
	"context"
	"database/sql"
	"fmt"
	"oma/internal"
	"oma/internal/storage"
	"strconv"
)

func GitRevert(ctx context.Context, repoContainer *storage.RepositoryContainer, fileIngrediends *[]FileIngredients, backFlag internal.Flag) error {
	repoId, err := repoContainer.FileIORepository.GetRepositoryId()

	if err != nil {
		return err
	}

	fmt.Printf("backFlag.Value: %v\n", backFlag.Value)
	backAmount, err := strconv.Atoi(backFlag.Value)

	if err != nil {
		return fmt.Errorf("back flag's value must be an integer")
	}

	for _, file := range *fileIngrediends {
		repository, err := repoContainer.OmaRepository.GetLatestByFileName(ctx, sql.NullString{
			Valid:  true,
			String: file.fileName,
		}, repoId)

		if err != nil {
			// filename is an absolute path so this should work?
			if err := repoContainer.FileIORepository.DeleteFile(file.fileName); err != nil {
				return fmt.Errorf("file %v did not exist %v commits ago, however, the attempt of deleting it was not successful", file.fileName, backAmount)
			}
			// TODO: here I should remove the OmaRepository entry
		}

		maxVersion, _ := repoContainer.VersionsRepository.GetMaxVersionNumberForRepo(ctx, repository.ID)

		// means that this file did not have that many commits yet
		if maxVersion < backAmount {
			continue
		}

		// FIXME: why the latest? I don't know what I was thinking
		// it should be get latest X by repo id
		// X = backAmount
		versions, err := repoContainer.VersionsRepository.GetLatestByRepositoryId(ctx, repository.ID)

		if err != nil {
			continue
		}

		fmt.Printf("found version: %+v\nfor file: %v\n", versions, file.fileName)

		versionActions, err := repoContainer.VersionActionsRepository.GetByVersionId(ctx, versions.ID)

		if err != nil {
			return fmt.Errorf("there are versions defined for this file: %v, but no version actions?\nError:%w", file.fileName, err)
		}

		revertedFile, err := RebuildDiff(file.content, versionActions)

		if err != nil {
			return fmt.Errorf("error while rebuilding the old version of file: %v, error:\n%w", file.fileName, err)
		}

		err = repoContainer.FileIORepository.WriteToFile(file.fileName, revertedFile)

		if err != nil {
			return fmt.Errorf("error while writing the reverted file: %v, error:\n%w", file.fileName, err)
		}

		// TODO: here I should delete the versions (version actions are tied with a cascade so they'll be auto deleted)
	}

	return nil
}
