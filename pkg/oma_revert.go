package pkg

import (
	"context"
	"database/sql"
	"fmt"
	"oma/internal/storage"
)

func GitRevert(ctx context.Context, repoContainer *storage.RepositoryContainer, fileIngrediends *[]FileIngredients) error {
	repoId, err := repoContainer.FileIORepository.GetRepositoryId()
	fmt.Printf("repoId: %v\n", repoId)

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

		if len(*versions) > 0 {
			fmt.Printf("versions: %v\n", versions)
		}

		// versionActions := repoContainer.VersionActionsRepository

	}

	return nil
}
