package postgres

import (
	"context"
	"github.com/jmoiron/sqlx"
	"log"
	"oma/internal/db/models"
)

type VersionRepository interface {
	Create(ctx context.Context, data *models.Versions) (*models.Versions, error)
	// Get(ctx context.Context, id int) (*models.Versions, error)
	// Update(ctx context.Context, id int, data *models.Versions) (*models.Versions, error)
	// Delete(ctx context.Context, id int) error
}

type VersionRepositoryImpl struct {
	db *sqlx.DB
}

func versionsToMap(data *models.Versions) map[string]any {
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

func (r *VersionRepositoryImpl) Create(ctx context.Context, data *models.Versions) (*models.Versions, error) {
	query, args, err := Sq.Insert("versions").SetMap(versionsToMap(data)).Suffix("returning *").ToSql()

	if err != nil {
		log.Fatalf("error while generating the create versions query, %v", err)
	}

	createdRepo := &models.Versions{}
	err = r.db.GetContext(ctx, createdRepo, query, args...)

	return createdRepo, err
}
