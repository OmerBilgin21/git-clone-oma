package pkg

import (
	"context"
	"database/sql"
	"log"
	"oma/internal/db/models"
	"oma/internal/db/postgres"
	"oma/internal/db/repositories"
	"os"
	"slices"
	"time"

	"github.com/jmoiron/sqlx"
)

var OMA_IGNORE_DEFAULTS = []string{".git", ".oma", ".omaignore", ".gitignore"}

func walkDirsAndReturn() []FileIngredients {
	currDir, err := os.Getwd()
	check(err, true)

	ignoreList := ParseOmaIgnore()

	var fileIngredients []FileIngredients
	WalkDirs(currDir, &fileIngredients, []string{}, ignoreList)
	return fileIngredients
}

func ParseAndDispatch(args []string, dbIns *sqlx.DB) {
	defer dbIns.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	omaRepo := postgres.NewOmaRepository(dbIns)
	repoContainer := repositories.RepositoryContainer{
		OmaRepository: omaRepo,
	}

	if slices.Contains(args, "init") {
		if len(args) > 2 {
			log.Fatal("illogical flags/commands type oma init --help for usage")
		} else if len(args) == 2 && args[1] == "--help" {
			log.Fatal("help docs, TBD")
		}

		fileIngredients := walkDirsAndReturn()

		for _, entry := range fileIngredients {
			_, err := repoContainer.OmaRepository.Create(ctx, &models.OmaRepository{
				FileName: sql.NullString{
					String: entry.fileName,
					Valid:  true,
				},

				CachedText: sql.NullString{
					String: entry.content,
					Valid:  true,
				},
			})
			check(err, true)
		}

		log.Print("repository initialized succesfully")
		return
	}

	if slices.Contains(args, "diff") {
		fileIngredients := walkDirsAndReturn()

		for _, ingredient := range fileIngredients {
			newres, err := omaRepo.GetLatestByFileName(ctx, sql.NullString{
				String: ingredient.fileName,
				Valid:  true,
			})
			check(err, false)

			if newres.ID == 0 {
				omaRepo.Create(ctx, &models.OmaRepository{
					FileName: sql.NullString{
						String: ingredient.fileName,
						Valid:  true,
					},
					CachedText: sql.NullString{
						Valid:  true,
						String: ingredient.content,
					},
				})
			} else {
				additions, deletions := Diff(newres.CachedText.String, ingredient.content)
				oldColoured, newColoured := ColourTheDiffs(additions, deletions, newres.CachedText.String, ingredient.content)
				if len(additions) > 0 || len(deletions) > 0 {
					RenderDiffs(oldColoured, newColoured)
				}
			}

		}
	}
}
