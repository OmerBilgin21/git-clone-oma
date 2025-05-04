package pkg

import (
	"context"
	"fmt"
	"oma/internal/storage"
)

func GitLog(ctx context.Context, repoContainer *storage.RepositoryContainer) error {
	repoId, err := repoContainer.FileIORepository.GetRepositoryId()

	if err != nil {
		return err
	}
	repositories, err := repoContainer.OmaRepository.GetAllByRepoId(ctx, repoId)
	if err != nil {
		return fmt.Errorf("no file has been found in this repository")
	}

	for _, repo := range repositories {
		versions, err := repoContainer.VersionsRepository.GetAllDistinctByRepoId(ctx, repo.ID)

		if err != nil {
			continue
		}

		for _, version := range versions {
			fmt.Printf("%s - %v\n", version.CreatedAt.Format("2006.01.02 - 15:04:05"), version.Message)
		}
	}
	return nil
}
