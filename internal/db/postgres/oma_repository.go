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

