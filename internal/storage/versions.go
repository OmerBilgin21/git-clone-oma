package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"log"
)

type VersionRepository interface {
	Create(ctx context.Context, data *Versions) (*Versions, error)
	Get(ctx context.Context, id int) (*Versions, error)
	GetLatestByRepositoryId(ctx context.Context, repoId int) (*[]Versions, *[]Versions, error)
}

type VersionRepositoryImpl struct {
	db *sqlx.DB
}

func NewVersionRepository(db *sqlx.DB) *VersionRepositoryImpl {
	return &VersionRepositoryImpl{db: db}
}

func (r *VersionRepositoryImpl) Create(ctx context.Context, data *Versions) (*Versions, error) {
	query, args, err := sq.Insert("versions").Columns("repository_id").Values(data.RepositoryId).Suffix("returning *").ToSql()

	if err != nil {
		log.Fatalf("error while generating the create versions query, %v", err)
	}

	createdRepo := &Versions{}
	err = r.db.GetContext(ctx, createdRepo, query, args...)

	return createdRepo, err
}

func (r *VersionRepositoryImpl) Get(ctx context.Context, id int) (*Versions, error) {
	query, args, err := sq.Select("*").From("versions").Where(squirrel.Eq{
		"id": id,
	}).ToSql()

	if err != nil {
		log.Fatalf("error while generating the create versions query, %v", err)
	}

	foundRepo := &Versions{}
	err = r.db.SelectContext(ctx, foundRepo, query, args...)

	return foundRepo, err
}

func (r *VersionRepositoryImpl) GetLatestByRepositoryId(ctx context.Context, repoId int) (*[]Versions, *[]Versions, error) {
	additionQuery, additionArgs, err := sq.Select("*").From("versions").Where(squirrel.Eq{
		"repository_id": repoId,
		"action_key":    AdditionKey,
	}).Limit(1).ToSql()

	deletionQuery, deletionArgs, err := sq.Select("*").From("versions").Where(squirrel.Eq{
		"repository_id": repoId,
		"action_key":    DeletionKey,
	}).Limit(1).ToSql()

	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("error while building GetLatestByRepositoryId query: %s\n", err))
	}

	additions := []Versions{}
	err = r.db.GetContext(ctx, &additions, additionQuery, additionArgs...)

	deletions := []Versions{}
	err = r.db.GetContext(ctx, &deletions, deletionQuery, deletionArgs...)

	return nil, nil, nil
}
