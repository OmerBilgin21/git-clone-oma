package pkg

import (
	"context"
	"fmt"
	"log"
	"oma/internal"
	"oma/internal/storage"
	"oma/util"
	"os"
	"strings"
	_ "strings"
	"time"

	"gorm.io/gorm"
)

type DispatchCommand struct {
	db                 *gorm.DB
	omaRepo            *storage.OmaRepositoryImpl
	versionsRepo       *storage.VersionRepositoryImpl
	versionActionsRepo *storage.VersionActionsRepositoryImpl
	fileIO             *storage.FileIOImpl
	// slices are auto by reference in Go, no need for explicit ptr
	fileIngredients []internal.FileIngredient
}

func NewDispatchCommand(db *gorm.DB, omaRepo *storage.OmaRepositoryImpl, versionRepo *storage.VersionRepositoryImpl, versionActionsRepo *storage.VersionActionsRepositoryImpl, fileIO *storage.FileIOImpl) *DispatchCommand {
	return &DispatchCommand{
		db:                 db,
		omaRepo:            omaRepo,
		versionsRepo:       versionRepo,
		versionActionsRepo: versionActionsRepo,
		fileIO:             fileIO,
	}
}

func (d *DispatchCommand) Dispatch(args []string, dbIns *gorm.DB) {
	sqlDB, err := dbIns.DB()
	if err != nil {
		panic(err)
	}

	defer sqlDB.Close()

	dbIns.AutoMigrate(&storage.OmaRepository{}, &storage.Versions{}, &storage.VersionActions{})

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fileIngredients := d.WalkDirsAndReadFiles()
	d.fileIngredients = fileIngredients

	parseArgs := internal.NewCLIArgsParser(args)
	var cmd internal.Command
	err = parseArgs.GetCommand(&cmd)

	if err != nil {
		log.Fatalf("error while parsing the commands:\n%v", err)
	}

	switch cmd {
	case internal.Init:
		if err := d.GitInit(ctx); err != nil {
			log.Fatalf("error while initialising repository:\n%v", err)
		}

	case internal.Commit:
		messageFlag, err := parseArgs.GetFlag("message")
		if err != nil {
			log.Fatalf("%v\n", err)
		}
		if err := d.GitCommit(ctx, messageFlag); err != nil {
			log.Fatalf("error while committing your changes:\n%v", err)
		}

	case internal.Diff:
		if err := d.GitDiff(ctx); err != nil {
			log.Fatalf("diff could not be displayed:\n%s", err)
		}
	case internal.Revert:
		backFlag, err := parseArgs.GetFlag("back")
		if err != nil {
			log.Fatalf("%v\n", err)
		}
		if err := d.GitRevert(ctx, backFlag); err != nil {
			log.Fatalf("error while reverting:\n%v", err)
		}
	case internal.Log:
		if err := d.GitLog(ctx); err != nil {
			log.Fatalf("error while logging the commit history: %v", err)
		}
		// if err := Gi
	}

}

func (d *DispatchCommand) ParseOmaIgnore() []string {
	omaIgnoreBytes, err := os.ReadFile("./.omaignore")
	if err != nil {
		panic(err)
	}
	omaIgnore := string(omaIgnoreBytes)

	separatedArgs := util.PurifyReadResult(strings.Split(omaIgnore, "\n"))
	separatedArgs = append(separatedArgs, internal.OMA_IGNORE_DEFAULTS...)

	return separatedArgs
}

func (d *DispatchCommand) WalkDirsAndReadFiles() []internal.FileIngredient {
	currDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	ignoreList := d.ParseOmaIgnore()

	var fileIngredients []internal.FileIngredient
	internal.WalkDirs(currDir, fileIngredients, []string{}, ignoreList, d.fileIO)

	return fileIngredients
}

func (d *DispatchCommand) GetAllVersionActionsForRepo(ctx context.Context, repositoryId int) ([]storage.VersionActions, error) {
	allVersions, err := d.versionsRepo.GetAllByRepoId(ctx, repositoryId)

	if err != nil {
		return []storage.VersionActions{}, fmt.Errorf("error while trying to get versions for repository: %v\nerror:\n%w", repositoryId, err)
	}

	var versionActions []storage.VersionActions
	for _, version := range allVersions {
		versionActionsOfVersion, err := d.versionActionsRepo.GetByVersionId(ctx, version.ID)
		if err != nil {
			return []storage.VersionActions{}, fmt.Errorf("error while trying to get version actions for version: %v\nerror:\n%w", version.ID, err)
		}
		if len(versionActionsOfVersion) > 0 {
			versionActions = append(versionActions, versionActionsOfVersion...)
		}
	}

	return versionActions, nil
}
