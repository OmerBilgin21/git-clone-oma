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
	GetLatestByRepositoryId(ctx context.Context, repoId int) (*Versions, error)
	GetMaxVersionNumberForRepo(ctx context.Context, repoId int) (int, error)
}

type VersionRepositoryImpl struct {
	db *sqlx.DB
}

func NewVersionRepository(db *sqlx.DB) *VersionRepositoryImpl {
	return &VersionRepositoryImpl{db: db}
}

func (versions *VersionRepositoryImpl) GetMaxVersionNumberForRepo(ctx context.Context, repoId int) (int, error) {
	maxIdQuery, maxIdArgs, err := sq.Select("max(version_id)").From("versions").Where(squirrel.Eq{
		"repository_id": repoId,
	}).ToSql()

	if err != nil {
		return -1, fmt.Errorf("error while generating the GetMaxVersionNumberForRepo for repository:\n%v\nerror:\n%w", repoId, err)
	}

	var maxId int

	err = versions.db.GetContext(ctx, &maxId, maxIdQuery, maxIdArgs...)

	if err != nil {
		return 0, err
	}

	return maxId, nil
}

func (versions *VersionRepositoryImpl) Create(ctx context.Context, data *Versions) (*Versions, error) {
	nextId, err := versions.GetMaxVersionNumberForRepo(ctx, data.RepositoryId)

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
		return nil, fmt.Errorf("error while generating the create versions query:\n%w", err)
	}

	createdRepo := &Versions{}
	err = versions.db.GetContext(ctx, createdRepo, query, args...)

	if err != nil {
		return nil, err
	}

	return createdRepo, nil
}

func (versions *VersionRepositoryImpl) Get(ctx context.Context, id int) (*Versions, error) {
	query, args, err := sq.Select("*").From("versions").Where(squirrel.Eq{
		"id": id,
	}).ToSql()

	if err != nil {
		return nil, fmt.Errorf("error while generating the create versions query:\n%w", err)
	}

	foundRepo := &Versions{}
	err = versions.db.SelectContext(ctx, foundRepo, query, args...)

	if err != nil {
		return nil, err
	}

	return foundRepo, nil
}

func (versions *VersionRepositoryImpl) GetLatestByRepositoryId(ctx context.Context, repoId int) (*Versions, error) {
	latestVersionId, err := versions.GetMaxVersionNumberForRepo(ctx, repoId)

	if err != nil {
		return nil, fmt.Errorf("error while finding the latest version id for repository:\n%v\nerror:\n%w", repoId, err)
	}

	versionsQuery, versionsArgs, err := sq.Select("*").From("versions").Where(squirrel.Eq{
		"repository_id": repoId,
		"version_id":    latestVersionId,
	}).ToSql()

	if err != nil {
		return nil, fmt.Errorf("error while building GetLatestByRepositoryId query:\n%w", err)
	}

	foundVersions := []Versions{}
	err = versions.db.SelectContext(ctx, &foundVersions, versionsQuery, versionsArgs...)

	if len(foundVersions) != 1 || err != nil {
		return nil, fmt.Errorf("something went very wrong, please create an issue. error:\n%w", err)
	}

	return &foundVersions[0], nil
}
