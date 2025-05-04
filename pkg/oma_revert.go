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
			return fmt.Errorf("no entry found for file: %v\nplease commit your changes before trying to revert", file.fileName)
		}

		maxVersion, _ := repoContainer.VersionsRepository.GetMaxVersionNumberForRepo(ctx, repository.ID)

		// means that this file did not have that many commits yet
		if maxVersion < backAmount {
			continue
		}

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
	}

	return nil
}
