package storage

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
)

func versionActionsToMap(data *VersionActions) map[string]any {
	return map[string]any{
		"start_x":    data.StartX,
		"start_y":    data.StartY,
		"end_x":      data.EndX,
		"end_y":      data.EndY,
		"action_key": data.ActionKey,
		"version_id": data.VersionId,
	}
}

type VersionActionsRepository interface {
	Create(ctx context.Context, data *VersionActions) (*VersionActions, error)
	// Get(ctx context.Context, id int) (*VersionActions, error)
	// GetLatestByRepositoryId(ctx context.Context, repoId int) (*[]VersionActions, *[]VersionActions, error)
}

type VersionActionsRepositoryImpl struct {
	db *sqlx.DB
}

func NewVersionActionsRepositoryImpl(db *sqlx.DB) *VersionActionsRepositoryImpl {
	return &VersionActionsRepositoryImpl{db: db}
}

func (self *VersionActionsRepositoryImpl) Create(ctx context.Context, data *VersionActions) (*VersionActions, error) {
	query, args, err := sq.Insert("version_actions").SetMap(versionActionsToMap(data)).Suffix("returning *").ToSql()

	if err != nil {
		log.Fatalf("error while generating the create version actions query, %v", err)
	}

	createdRepo := &VersionActions{}
	err = self.db.GetContext(ctx, createdRepo, query, args...)

	return createdRepo, err
}
