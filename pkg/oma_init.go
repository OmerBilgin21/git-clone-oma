package pkg

import (
	"context"
	"fmt"
	"oma/internal"
	"oma/internal/storage"
	"strings"
)

func (d *OmaVC) OmaInit(ctx context.Context) error {
	var randomCreatedRepo *storage.OmaRepository
	nextId, err := d.omaRepo.GetNextOmaRepoId(ctx)

	if err != nil {
		return err
	}

	for i, entry := range d.fileIngredients {
		createdRepo, err := d.omaRepo.Create(ctx, &storage.OmaRepository{
			FileName:   &entry.FileName,
			CachedText: &entry.Content,
			OmaRepoId:  nextId,
		})

		if err != nil {
			return err
		}

		if i == 0 {
			randomCreatedRepo = createdRepo
		}

	}

	err = d.fileIO.CreateRepoInitInfo(randomCreatedRepo.OmaRepoId)

	if err != nil {
		if strings.Contains(err.Error(), "exists") {
			return fmt.Errorf("can not initialize repository, it is already done")
		}
		return fmt.Errorf("error while creating repo init file:\n%v", err)
	}

	internal.Logger("repository initialized succesfully")
	return nil
}
