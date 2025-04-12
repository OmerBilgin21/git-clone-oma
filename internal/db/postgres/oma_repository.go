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

func (r *OmaRepositoryImpl) Get(ctx context.Context, id int) (*models.OmaRepository, error) {
	query, _, err := sq.Select("repositories").Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		log.Fatalf("error while getting: %v", err)
	}

	foundRepo := &models.OmaRepository{}

	err = r.db.SelectContext(ctx, foundRepo, query, id)
	return foundRepo, err

}

func (r *OmaRepositoryImpl) Update(ctx context.Context, id int, data *models.OmaRepository) (*models.OmaRepository, error) {
	qb := sq.Update("repositories")

	fmt.Printf("&data: %v\n", *data)

	if data.FileName != nil {
		qb = qb.Set("filename", *data.FileName)
	}
	if data.CachedText != nil {
		fmt.Printf("*data.CachedText: %v\n", *data.CachedText)
		qb = qb.Set("cached_text", *data.CachedText)
	}

	qb = qb.Where(squirrel.Eq{"id": strconv.Itoa(id)}).Suffix("returning *")

	query, args, err := qb.ToSql()
	if err != nil {
		log.Fatalf("error while updating: %v\n", err)
	}

	fmt.Printf("query: %v\n", query)
	fmt.Printf("args: %v\n", args)

	updatedRepo := &models.OmaRepository{}
	err = r.db.GetContext(ctx, updatedRepo, query, args...)

	return updatedRepo, err
}

func (r *OmaRepositoryImpl) Delete(ctx context.Context, id int) error {
	query := `update repositories set deleted_at = now() where id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
