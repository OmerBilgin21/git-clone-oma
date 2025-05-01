package pkg

import (
	"context"
	"fmt"
	"log"
	"oma/internal"
	"oma/internal/storage"
	_ "strings"
	"time"

	"github.com/jmoiron/sqlx"
)

func Dispatch(args []string, dbIns *sqlx.DB) {
	defer dbIns.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	repoContainer := storage.RepositoryContainer{
		OmaRepository:            storage.NewOmaRepository(dbIns),
		VersionsRepository:       storage.NewVersionRepository(dbIns),
		VersionActionsRepository: storage.NewVersionActionsRepository(dbIns),
		FileIORepository:         storage.NewFileIO(),
	}

	fileIngredients := walkDirsAndReadFiles()

	parseArgs := internal.NewCLIArgsParser(args)
	var cmd internal.Command
	var flags []internal.Flag
	err := parseArgs.GetCommand(&cmd)

	if err != nil {
		log.Fatalf("error while parsing the commands:\n%v", err)
	}

	err = parseArgs.GetFlags(&flags)

	if err != nil {
		log.Fatalf("error while parsing the flags:\n%v", err)
	}

	fmt.Printf("cmd: %v\n", cmd)
	fmt.Printf("flags: %+v\n", flags)

	switch cmd {
	case internal.Init:
		if err := GitInit(ctx, &repoContainer, &fileIngredients); err != nil {
			log.Fatalf("error while initialising repository:\n%v", err)
		}

	case internal.Commit:
		if err := GitCommit(ctx, &repoContainer, &fileIngredients); err != nil {
			log.Fatalf("error while committing your changes:\n%v", err)
		}

	case internal.Diff:
		if err := GitDiff(ctx, &repoContainer, &fileIngredients); err != nil {
			log.Fatalf("diff could not be displayed:\n%s", err)
		}
	case internal.Revert:
		if err := GitRevert(ctx, &repoContainer, &fileIngredients); err != nil {
			log.Fatalf("error while reverting:\n%v", err)
		}
	}

}
