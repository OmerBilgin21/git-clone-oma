package postgres

import (
	"context"
	_ "database/sql"
	"fmt"
	"log"
	"oma/internal/db/models"
	"strconv"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type OmaRepositoryRepository interface {
	Create(ctx context.Context, data *models.OmaRepository) (*models.OmaRepository, error)
	Get(ctx context.Context, id int) (*models.OmaRepository, error)
	Update(ctx context.Context, id int, data *models.OmaRepository) (*models.OmaRepository, error)
	Delete(ctx context.Context, id int) error
}

type OmaRepositoryImpl struct {
	db *sqlx.DB
}

var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

func NewOmaRepository(db *sqlx.DB) *OmaRepositoryImpl {
	return &OmaRepositoryImpl{db: db}
}

func (r *OmaRepositoryImpl) Create(ctx context.Context, data *models.OmaRepository) (*models.OmaRepository, error) {
	query := `insert into repositories (cached_text, filename) 
	values ($1, $2) 
	returning *`

	createdRepo := &models.OmaRepository{}

	fmt.Printf("data.CachedText: %v\n", *data.CachedText)
	err := r.db.GetContext(ctx, createdRepo, query, *data.CachedText, *data.FileName)
	if err != nil {
		log.Print(err)
	}

	return createdRepo, err
}
