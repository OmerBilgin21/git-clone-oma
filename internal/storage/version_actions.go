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
		"content":    data.Content,
	}
}

type VersionActionsRepository interface {
	Create(ctx context.Context, data *VersionActions) (*VersionActions, error)
	GetByVersionId(ctx context.Context, versionId int) ([]VersionActions, error)
	DeleteByVersionId(ctx context.Context, versionId int) error
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
		return nil, fmt.Errorf("error while generating the create version actions query:\n%w", err)
	}

	createdRepo := &VersionActions{}
	err = versionActions.db.GetContext(ctx, createdRepo, query, args...)

	return createdRepo, err
}

func (versionActions *VersionActionsRepositoryImpl) GetByVersionId(ctx context.Context, versionId int) ([]VersionActions, error) {
	query, args, err := sq.Select("*").From("version_actions").Where(squirrel.Eq{
		"version_id": versionId,
	}).Where(squirrel.Expr("deleted_at IS NOT NULL")).ToSql()

	if err != nil {
		return nil, fmt.Errorf("error while generating the GetByVersionId query:\n%w", err)
	}

	foundVersionActions := []VersionActions{}

	err = versionActions.db.SelectContext(ctx, &foundVersionActions, query, args...)

	if err != nil {
		return nil, fmt.Errorf("error while finding version actions for version: %v\nerror:\n%w", versionId, err)
	}

	return foundVersionActions, err
}

func (versionActions *VersionActionsRepositoryImpl) DeleteByVersionId(ctx context.Context, versionId int) error {
	query := `update version_actions set deleted_at = now() where version_id = $1`
	_, err := versionActions.db.ExecContext(ctx, query, versionId)

	return err
}
