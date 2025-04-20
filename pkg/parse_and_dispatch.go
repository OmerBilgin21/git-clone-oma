package pkg

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"oma/internal/storage"
	"os"
	"slices"
	"strings"
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
				additions, deletions := Diff(newres.CachedText.String, ingredient.content)
				oldColoured, newColoured := ColourTheDiffs(additions, deletions, newres.CachedText.String, ingredient.content)
				if len(additions) > 0 || len(deletions) > 0 {
					RenderDiffs(oldColoured, newColoured)
				}
			}

		}
	} else if slices.Contains(args, "commit") {
		slices.Sort(args)
		commitIndex, wasFound := slices.BinarySearchFunc(args, "message", func(s1, target string) int {
			if strings.HasPrefix(s1, "message") {
				return 0
			}

			if s1 > target {
				return 1
			}
			return -1
		})
		// commitIndex, wasFound := slices.BinarySearch(args, "message")
		fmt.Printf("commitIndex: %v\n", commitIndex)
		fmt.Printf("wasFound: %v\n", wasFound)
		// if
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
				additions, deletions := Diff(newres.CachedText.String, ingredient.content)

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
