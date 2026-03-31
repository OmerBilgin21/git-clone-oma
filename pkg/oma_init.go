package pkg

import (
	"context"
	"fmt"
	"log"
	"oma/internal"
	"oma/internal/storage"
	"strings"
)

func GitInit(ctx context.Context, repoContainer *storage.RepositoryContainer, fileIngredients *[]internal.FileIngredients) error {
	var randomCreatedRepo *storage.OmaRepository
	nextId, err := repoContainer.OmaRepository.GetNextOmaRepoId(ctx)

	if err != nil {
		return err
	}

	for i, entry := range *fileIngredients {
		createdRepo, err := repoContainer.OmaRepository.Create(ctx, &storage.OmaRepository{
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

	err = repoContainer.FileIORepository.CreateRepoInitInfo(randomCreatedRepo.OmaRepoId)

	if err != nil {
		if strings.Contains(err.Error(), "exists") {
			return fmt.Errorf("can not initialize repository, it is already done")
		}
		return fmt.Errorf("error while creating repo init file:\n%v", err)
	}

	log.Print("repository initialized succesfully")
	return nil
}
