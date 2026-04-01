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
	fileIO := storage.NewFileIO()

	dispatch := pkg.NewDispatchCommand(
		db,
		omaRepo,
		versionRepo,
		versionActionsRepo,
		fileIO,
	)

	dispatch.Dispatch(cliArgs, db)
}
