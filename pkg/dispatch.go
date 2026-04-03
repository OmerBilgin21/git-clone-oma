package pkg

import (
	"context"
	"oma/internal"
	"oma/internal/storage"
	"os"
	"slices"
	"strings"
	"time"

	"gorm.io/gorm"
)

type OmaVC struct {
	db                 *gorm.DB
	omaRepo            *storage.OmaRepositoryImpl
	versionsRepo       *storage.VersionRepository
	versionActionsRepo *storage.VersionActionsRepository
	fileIO             *storage.FileIOImpl
	// slices are auto by reference in Go, no need for explicit ptr
	fileIngredients []internal.FileIngredient
}

func NewOmaVC(db *gorm.DB, omaRepo *storage.OmaRepositoryImpl, versionRepo *storage.VersionRepository, versionActionsRepo *storage.VersionActionsRepository, fileIO *storage.FileIOImpl) *OmaVC {
	return &OmaVC{
		db:                 db,
		omaRepo:            omaRepo,
		versionsRepo:       versionRepo,
		versionActionsRepo: versionActionsRepo,
		fileIO:             fileIO,
	}
}

func (d *OmaVC) RunCMD(args []string, dbIns *gorm.DB) {
	sqlDB, err := dbIns.DB()

	if err != nil {
		internal.LogAndExit("error while getting the DB instance", err)
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
		internal.LogAndExit("error while parsing arguments", err)
	}

	switch cmd {
	case internal.Init:
		if err := d.OmaInit(ctx); err != nil {
			internal.LogAndExit("error while initialising repository", err)
		}

	case internal.Commit:
		messageFlag, err := parseArgs.GetFlag("message")
		if err != nil {
			internal.LogAndExit(err)
		}
		if err := d.OmaCommit(ctx, messageFlag); err != nil {
			internal.LogAndExit("error while committing your changes", err)
		}

	case internal.Diff:
		if err := d.OmaDiff(ctx); err != nil {
			internal.LogAndExit("diff could not be displayed", err)
		}

	case internal.Revert:
		backFlag, err := parseArgs.GetFlag("back")
		if err != nil {
			internal.LogAndExit("error while parsing arguments for revert", err)
		}
		if err := d.OmaRevert(ctx, backFlag); err != nil {
			internal.LogAndExit("error while reverting", err)
		}

	case internal.Log:
		if err := d.OmaLog(ctx); err != nil {
			internal.LogAndExit("error while logging the commit history", err)
		}
	}

}

func (d *OmaVC) ParseOmaIgnore() []string {
	omaIgnoreBytes, err := os.ReadFile("./.omaignore")
	if err != nil {
		panic(err)
	}

	omaIgnore := string(omaIgnoreBytes)

	lines := strings.Split(omaIgnore, "\n")
	if len(lines) > 0 {
		for i, elem := range lines {
			if elem == "" || elem == " " || elem == "\n" {
				lines = slices.Delete(lines, i, i+1)
			}
		}
	}

	lines = append(lines, internal.OMA_IGNORE_DEFAULTS...)

	return lines
}

func (d *OmaVC) WalkDirsAndReadFiles() []internal.FileIngredient {
	currDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	ignoreList := d.ParseOmaIgnore()

	var fileIngredients []internal.FileIngredient
	internal.WalkDirs(currDir, &fileIngredients, []string{}, ignoreList, d.fileIO)

	return fileIngredients
}

func (d *OmaVC) GetAllVersionActionsForRepo(ctx context.Context, repositoryId int) ([]storage.VersionActions, error) {
	allVersions, err := d.versionsRepo.GetAllByRepoId(ctx, repositoryId)

	if err != nil {
		return []storage.VersionActions{}, err
	}

	var versionActions []storage.VersionActions
	for _, version := range allVersions {
		versionActionsOfVersion, err := d.versionActionsRepo.GetByVersionId(ctx, version.ID)
		if err != nil {
			return []storage.VersionActions{}, err
		}
		if len(*versionActionsOfVersion) > 0 {
			versionActions = append(versionActions, *versionActionsOfVersion...)
		}
	}

	return versionActions, nil
}
