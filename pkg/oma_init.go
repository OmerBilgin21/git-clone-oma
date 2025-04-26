package pkg

import (
	"context"
	"database/sql"
	"log"
	"oma/internal/storage"
)

func GitInit(ctx context.Context, repoContainer *storage.RepositoryContainer, fileIngredients *[]FileIngredients) error {
	for _, entry := range *fileIngredients {
		_, err := repoContainer.OmaRepository.Create(ctx, &storage.OmaRepository{
			FileName: sql.NullString{
				String: entry.fileName,
				Valid:  true,
			},

			CachedText: sql.NullString{
				String: entry.content,
				Valid:  true,
			},
		})

		if err != nil {
			return err
		}
	}

	log.Print("repository initialized succesfully")
	return nil
}
