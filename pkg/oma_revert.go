package pkg

import (
	"context"
	"database/sql"
	"fmt"
	"oma/internal/storage"
)

func GitRevert(ctx context.Context, repoContainer *storage.RepositoryContainer, fileIngrediends *[]FileIngredients) error {
	repoId, err := repoContainer.FileIORepository.GetRepositoryId()

	if err != nil {
		return err
	}

	for _, file := range *fileIngrediends {

		repository, err := repoContainer.OmaRepository.GetLatestByFileName(ctx, sql.NullString{
			Valid:  true,
			String: file.fileName,
		}, repoId)

		fmt.Printf("repository: %v\n", repository.OmaRepoId)

		if err != nil {
			return fmt.Errorf("no repository found for file: %v\nerror: %v\n", file.fileName, err)
		}

		versions, err := repoContainer.VersionsRepository.GetLatestByRepositoryId(ctx, repository.ID)

		if err != nil {
			return fmt.Errorf("no versions found for repository: %v\nerror: %v\n", repository.OmaRepoId, err)
		}

		fmt.Printf("versions: %v\n", versions)

		versionActions, err := repoContainer.VersionActionsRepository.GetByVersionId(ctx, versions.ID)

		// FIXME: don't return this, skip it.
		// this means there's no found version action, which is fine.
		// just skip it
		if err != nil {
			return err
		}

		// TODO: build old versions here

		fmt.Printf("len(versionActions): %v\n", len(versionActions))

		fmt.Printf("yo ended man\n")
	}

	return nil
}
