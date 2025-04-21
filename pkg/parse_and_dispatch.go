package pkg

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"oma/internal/storage"
	"os"
	"slices"
	_ "strings"
	"time"

	"github.com/jmoiron/sqlx"
)

var OMA_IGNORE_DEFAULTS = []string{".git", ".oma", ".omaignore", ".gitignore", "node_modules"}

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

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	repoContainer := storage.RepositoryContainer{
		OmaRepository:            storage.NewOmaRepository(dbIns),
		VersionsRepository:       storage.NewVersionRepository(dbIns),
		VersionActionsRepository: storage.NewVersionActionsRepositoryImpl(dbIns),
	}

	parseArgs := NewCLIArgsParser(args)
	parseArgs.Validate()
	command := parseArgs.GetCommand()
	namedArg, namedValue := parseArgs.GetPrefixAndValue("message")
	if command == "commit" {
		fmt.Printf("namedArg: %v\n", namedArg)
		fmt.Printf("namedValue: %v\n", namedValue)
	}

	if slices.Contains(args, "init") {
		if len(args) > 2 {
			log.Fatal("illogical flags/commands type oma init --help for usage")
		} else if len(args) == 2 && args[1] == "--help" {
			log.Fatal("help docs, TBD")
		}

		fileIngredients := walkDirsAndReturn()

		for _, entry := range fileIngredients {
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
			check(err, true)
		}

		log.Print("repository initialized succesfully")
		return
	} else if slices.Contains(args, "diff") {
		fileIngredients := walkDirsAndReturn()

		for _, ingredient := range fileIngredients {
			newres, err := repoContainer.OmaRepository.GetLatestByFileName(ctx, sql.NullString{
				String: ingredient.fileName,
				Valid:  true,
			})
			check(err, false)

			if newres.ID == 0 {
				repoContainer.OmaRepository.Create(ctx, &storage.OmaRepository{
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
				RenderDiffs(newres.CachedText.String, ingredient.content, newres.FileName.String, ingredient.fileName)
			}

		}
	} else if slices.Contains(args, "commit") {
		FileIngredients := walkDirsAndReturn()

		for _, ingredient := range FileIngredients {
			newres, err := repoContainer.OmaRepository.GetLatestByFileName(ctx, sql.NullString{
				String: ingredient.fileName,
				Valid:  true,
			})
			check(err, false)

			if newres.ID == 0 {
				fmt.Printf("No previous version of the file, creating cache...")
				repoContainer.OmaRepository.Create(ctx, &storage.OmaRepository{
					FileName: sql.NullString{
						Valid:  true,
						String: ingredient.fileName,
					},
					CachedText: sql.NullString{
						Valid:  true,
						String: ingredient.content,
					},
				})

				check(err, true)
			} else {
				additions, deletions := GetDiffs(newres.CachedText.String, ingredient.content)

				newVersion, err := repoContainer.VersionsRepository.Create(ctx, &storage.Versions{
					RepositoryId: newres.ID,
				})

				check(err, true)

				for i := range additions {
					addition := additions[i]
					_, err := repoContainer.VersionActionsRepository.Create(ctx, &storage.VersionActions{
						StartX:    addition.StartX,
						StartY:    addition.StartY,
						EndX:      addition.EndX,
						EndY:      addition.EndY,
						ActionKey: storage.AdditionKey,
						VersionId: newVersion.ID,
					})
					check(err, true)
				}

				for i := range deletions {
					deletion := deletions[i]
					_, err := repoContainer.VersionActionsRepository.Create(ctx, &storage.VersionActions{
						StartX:    deletion.StartX,
						StartY:    deletion.StartY,
						EndX:      deletion.EndX,
						EndY:      deletion.EndY,
						ActionKey: storage.DeletionKey,
						VersionId: newVersion.ID,
					})
					check(err, true)
				}

				log.Printf("Commit created succesfully for file: %v\n", ingredient.fileName)
			}
		}
	}
}
