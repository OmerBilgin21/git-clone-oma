package storage

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func versionActionsToMap(data *VersionActions) map[string]any {
	return map[string]any{
		"start":      data.Start,
		"dest":       data.Dest,
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

func (versionActions *VersionActionsRepositoryImpl) Create(ctx context.Context, data *VersionActions) (*VersionActions, error) {
	query, args, err := sq.Insert("version_actions").SetMap(versionActionsToMap(data)).Suffix("returning *").ToSql()

	if err != nil {
		return nil, fmt.Errorf("error while generating the create version actions query: %v", err)
	}

	createdRepo := &VersionActions{}
	err = versionActions.db.GetContext(ctx, createdRepo, query, args...)

	return createdRepo, err
}
