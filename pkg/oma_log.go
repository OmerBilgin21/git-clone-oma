package pkg

import (
	"context"
	"fmt"
)

func (d *OmaVC) OmaLog(ctx context.Context) error {
	repoId, err := d.fileIO.GetRepositoryId()

	if err != nil {
		return err
	}
	repositories, err := d.omaRepo.GetAllByRepoId(ctx, repoId)
	if err != nil {
		return fmt.Errorf("no file has been found in this repository")
	}

	for _, repo := range repositories {
		versions, err := d.versionsRepo.GetAllDistinctByRepoId(ctx, repo.ID)

		if err != nil {
			continue
		}

		for _, version := range versions {
			fmt.Printf("%s - %v\n", version.CreatedAt.Format("2006.01.02 - 15:04:05"), version.Message)
		}
	}
	return nil
}
