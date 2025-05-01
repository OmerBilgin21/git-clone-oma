package storage

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type VersionRepository interface {
	Create(ctx context.Context, data *Versions) (*Versions, error)
	Get(ctx context.Context, id int) (*Versions, error)
	GetLatestByRepositoryId(ctx context.Context, repoId int) (*[]Versions, error)
}

type VersionRepositoryImpl struct {
	db *sqlx.DB
}

func NewVersionRepository(db *sqlx.DB) *VersionRepositoryImpl {
	return &VersionRepositoryImpl{db: db}
}

func (versions *VersionRepositoryImpl) Create(ctx context.Context, data *Versions) (*Versions, error) {
	nextIdQuery, nextIdArgs, err := sq.Select("max(version_id)").From("versions").Where(squirrel.Eq{
		"repository_id": data.RepositoryId,
	}).ToSql()

	if err != nil {
		return nil, fmt.Errorf("error while generating the nextIdQuery for repository:\n%v\nerror:\n%v", data.RepositoryId, err)
	}

	var nextId int
	err = versions.db.GetContext(ctx, &nextId, nextIdQuery, nextIdArgs...)

	// means first version for a repo
	if err != nil {
		nextId = 1
	} else {
		nextId++
	}

	query, args, err := sq.Insert("versions").SetMap(map[string]any{
		"version_id":    nextId,
		"repository_id": data.RepositoryId,
	}).Suffix("returning *").ToSql()

	if err != nil {
		return nil, fmt.Errorf("error while generating the create versions query:\n%v", err)
	}

	createdRepo := &Versions{}
	err = versions.db.GetContext(ctx, createdRepo, query, args...)

	return createdRepo, err
}

func (versions *VersionRepositoryImpl) Get(ctx context.Context, id int) (*Versions, error) {
	query, args, err := sq.Select("*").From("versions").Where(squirrel.Eq{
		"id": id,
	}).ToSql()

	if err != nil {
		return nil, fmt.Errorf("error while generating the create versions query:\n%v", err)
	}

	foundRepo := &Versions{}
	err = versions.db.SelectContext(ctx, foundRepo, query, args...)

	return foundRepo, err
}

func (versions *VersionRepositoryImpl) GetLatestByRepositoryId(ctx context.Context, repoId int) (*[]Versions, error) {
	findLatestQuery, findLatestArgs, err := sq.Select("max(version_id)").From("versions").Where(squirrel.Eq{
		"repository_id": repoId,
	}).ToSql()

	if err != nil {
		return nil, fmt.Errorf("error while building the find latest version_id query for repository:\n%v\nerror:\n%v", repoId, err)
	}

	var latestVersionId int
	err = versions.db.GetContext(ctx, &latestVersionId, findLatestQuery, findLatestArgs...)

	if err != nil {
		return nil, fmt.Errorf("error while finding the latest version id for repository:\n%v\nerror:\n%v", repoId, err)
	}

	versionsQuery, versionsArgs, err := sq.Select("*").From("versions").Where(squirrel.Eq{
		"repository_id": repoId,
		"version_id":    latestVersionId,
	}).ToSql()

	if err != nil {
		return nil, fmt.Errorf("error while building GetLatestByRepositoryId query:\n%v", err)
	}

	foundVersions := []Versions{}
	err = versions.db.SelectContext(ctx, &foundVersions, versionsQuery, versionsArgs...)

	return &foundVersions, err
}
