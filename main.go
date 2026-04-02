package main

import (
	"oma/internal"
	"oma/internal/storage"
	"oma/pkg"
	"os"
)

func main() {
	cliArgs := os.Args[1:]
	db := internal.GetDb()

	omaRepo := storage.NewOmaRepository(db)
	versionRepo := storage.NewVersionRepository(db)
	versionActionsRepo := storage.NewVersionActionsRepository(db)
	fileIO, err := storage.NewFileIO()

	if err != nil {
		internal.LogAndExit("error while getting fileIO", err)
	}

	omaVC := pkg.NewOmaVC(
		db,
		omaRepo,
		versionRepo,
		versionActionsRepo,
		fileIO,
	)

	omaVC.RunCMD(cliArgs, db)
}
