package storage

import (
	"context"
	"github.com/jmoiron/sqlx"
	"log"
)

type VersionRepository interface {
	Create(ctx context.Context, data *Versions) (*Versions, error)
	// Get(ctx context.Context, id int) (*Versions, error)
	// Update(ctx context.Context, id int, data *Versions) (*Versions, error)
	// Delete(ctx context.Context, id int) error
}

type VersionRepositoryImpl struct {
	db *sqlx.DB
}

func versionsToMap(data *Versions) map[string]any {
	return map[string]any{
		"start_x":       data.StartX,
		"start_y":       data.StartY,
		"end_x":         data.EndX,
		"end_y":         data.EndY,
		"action_key":    data.ActionKey,
		"repository_id": data.RepositoryId,
	}
}

func NewVersionRepository(db *sqlx.DB) *VersionRepositoryImpl {
	return &VersionRepositoryImpl{db: db}
}

func (r *VersionRepositoryImpl) Create(ctx context.Context, data *Versions) (*Versions, error) {
	query, args, err := sq.Insert("versions").SetMap(versionsToMap(data)).Suffix("returning *").ToSql()

	if err != nil {
		log.Fatalf("error while generating the create versions query, %v", err)
	}

	createdRepo := &Versions{}
	err = r.db.GetContext(ctx, createdRepo, query, args...)

	return createdRepo, err
}
