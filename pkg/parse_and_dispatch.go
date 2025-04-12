package pkg

import (
	"context"
	"fmt"
	"log"
	"os"
	"slices"
	"time"

	// "oma/internal/db"
	// "oma/internal/db/repositories"
	"oma/internal/db/models"
	"oma/internal/db/postgres"
	"oma/internal/db/repositories"

	"github.com/jmoiron/sqlx"
)

var OMA_IGNORE_DEFAULTS = []string{".git", ".oma", ".omaignore", ".gitignore"}

func ParseAndDispatch(args []string, dbIns *sqlx.DB) {
	defer dbIns.Close()
	if slices.Contains(args, "init") {
		if len(args) > 2 {
			log.Fatal("illogical flags/commands type oma init --help for usage")
		} else if len(args) == 2 && args[1] == "--help" {
			log.Fatal("help docs, TBD")
		}

		currDir, err := os.Getwd()
		check(err, true)

		ignoreList := ParseOmaIgnore()

		var fileIngredients []FileIngredients
		WalkDirs(currDir, &fileIngredients, []string{}, ignoreList)

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		omaRepo := postgres.NewOmaRepository(dbIns)
		repoContainer := repositories.RepositoryContainer{
			OmaRepository: omaRepo,
		}

		for _, entry := range fileIngredients {

			createdRepo, err := repoContainer.OmaRepository.Create(ctx, &models.OmaRepository{
				FileName:   &entry.fileName,
				CachedText: &entry.content,
			})
			check(err, true)

			fmt.Printf("saved: %v\n", *createdRepo.CachedText)
		}

		log.Print("repository initialized succesfully")
		return
	}

	if slices.Contains(args, "commit") {

	}
}
