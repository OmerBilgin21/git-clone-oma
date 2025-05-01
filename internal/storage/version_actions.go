package storage

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
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
	GetByVersionId(ctx context.Context, versionId int) ([]VersionActions, error)
}

type VersionActionsRepositoryImpl struct {
	db *sqlx.DB
}

func NewVersionActionsRepository(db *sqlx.DB) *VersionActionsRepositoryImpl {
	return &VersionActionsRepositoryImpl{db: db}
}

func (versionActions *VersionActionsRepositoryImpl) Create(ctx context.Context, data *VersionActions) (*VersionActions, error) {
	query, args, err := sq.Insert("version_actions").SetMap(versionActionsToMap(data)).Suffix("returning *").ToSql()

	if err != nil {
		return nil, fmt.Errorf("error while generating the create version actions query:\n%v", err)
	}

	createdRepo := &VersionActions{}
	err = versionActions.db.GetContext(ctx, createdRepo, query, args...)

	return createdRepo, err
}

func (versionActions *VersionActionsRepositoryImpl) GetByVersionId(ctx context.Context, versionId int) ([]VersionActions, error) {
	query, _, err := sq.Select("*").From("version_actions").Where(squirrel.Eq{
		"version_id": versionId,
	}).ToSql()

	if err != nil {
		return nil, fmt.Errorf("error while generating the GetByVersionId query:\n%v", err)
	}

	foundVersionActions := []VersionActions{}

	err = versionActions.db.SelectContext(ctx, foundVersionActions, query)

	if err != nil {
		return nil, fmt.Errorf("error while finding version actions for version: %v\nerror:\n%v", versionId, err)
	}

	return foundVersionActions, err
}
