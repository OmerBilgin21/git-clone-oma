package pkg

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"oma/internal/storage"
	"strings"
)

func GitInit(ctx context.Context, repoContainer *storage.RepositoryContainer, fileIngredients *[]FileIngredients) error {
	var randomCreatedRepo *storage.OmaRepository
	nextId, err := repoContainer.OmaRepository.GetNextOmaRepoId(ctx)

	if err != nil {
		return err
	}

	for i, entry := range *fileIngredients {
		createdRepo, err := repoContainer.OmaRepository.Create(ctx, &storage.OmaRepository{
			FileName: sql.NullString{
				String: entry.fileName,
				Valid:  true,
			},
			CachedText: sql.NullString{
				String: entry.content,
				Valid:  true,
			},
			OmaRepoId: nextId,
		})

		if i == 0 {
			randomCreatedRepo = createdRepo
		}

		if err != nil {
			return err
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
