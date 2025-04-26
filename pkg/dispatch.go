package pkg

import (
	"context"
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
		VersionActionsRepository: storage.NewVersionActionsRepositoryImpl(dbIns),
	}

	fileIngredients := walkDirsAndReadFiles()

	parseArgs := internal.NewCLIArgsParser(args)
	var cmd internal.Command
	var flags []internal.Flag
	err := parseArgs.GetCommand(&cmd)

	if err != nil {
		log.Fatalf("error while parsing the commands: %v", err)
	}

	err = parseArgs.GetFlags(&flags)

	if err != nil {
		log.Fatalf("error while parsing the flags: %v", err)
	}

	switch cmd {
	case internal.Init:
		GitInit(ctx, &repoContainer, &fileIngredients)
	case internal.Commit:
		GitCommit(ctx, &repoContainer, &fileIngredients)
	case internal.Diff:
		GitDiff(ctx, &repoContainer, &fileIngredients)
	}

}
