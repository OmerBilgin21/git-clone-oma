package storage

import (
	"context"
	"database/sql"
	_ "database/sql"
	"fmt"
	"strconv"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type OmaRepoRepository interface {
	Create(ctx context.Context, data *OmaRepository) (*OmaRepository, error)
	Get(ctx context.Context, id int) (*OmaRepository, error)
	GetMany(ctx context.Context, ids []int) (*[]OmaRepository, error)
	GetLatestByFileName(ctx context.Context, filename sql.NullString) (*OmaRepository, error)
	Update(ctx context.Context, id int, data *OmaRepository) (*OmaRepository, error)
	Delete(ctx context.Context, id int) error
}

type OmaRepositoryImpl struct {
	db *sqlx.DB
}

func NewOmaRepository(db *sqlx.DB) *OmaRepositoryImpl {
	return &OmaRepositoryImpl{db: db}
}

func (omaRepo *OmaRepositoryImpl) Create(ctx context.Context, data *OmaRepository) (*OmaRepository, error) {
	query := `insert into repositories (cached_text, filename) 
	values ($1, $2) 
	returning *`

	var cachedText, fileName any

	if data.CachedText.Valid {
		cachedText = data.CachedText.String
	} else {
		cachedText = nil
	}

	if data.FileName.Valid {
		fileName = data.FileName.String
	} else {
		fileName = nil
	}

	createdRepo := &OmaRepository{}

	err := omaRepo.db.GetContext(ctx, createdRepo, query, cachedText, fileName)

	if err != nil {
		return nil, fmt.Errorf("error while creating an oma repository %v", err)
	}

	return createdRepo, err
}

func (omaRepo *OmaRepositoryImpl) Get(ctx context.Context, id int) (*OmaRepository, error) {
	query, _, err := sq.Select("repositories").Where(squirrel.Eq{"id": id}).ToSql()

	if err != nil {
		return nil, fmt.Errorf("error while getting: %v", err)
	}

	foundRepo := &OmaRepository{}

	err = omaRepo.db.SelectContext(ctx, foundRepo, query, id)
	return foundRepo, err

}

func (omaRepo *OmaRepositoryImpl) GetLatestByFileName(ctx context.Context, filename sql.NullString) (*OmaRepository, error) {
	query := "select * from repositories where filename = $1 order by id limit 1"

	if !filename.Valid {
		return nil, fmt.Errorf("you can not search for a null file name")
	}

	foundRepo := []OmaRepository{}
	err := omaRepo.db.SelectContext(ctx, &foundRepo, query, filename)
	if err != nil || len(foundRepo) != 1 {
		return nil, err
	}

	return &foundRepo[0], err
}

func (omaRepo *OmaRepositoryImpl) GetMany(ctx context.Context, ids []int) (*[]OmaRepository, error) {
	query, args, err := sq.Select("*").From("repositories").Where(squirrel.Eq{"id": ids}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("error while getting: %v", err)
	}

	foundRepos := []OmaRepository{}
	err = omaRepo.db.SelectContext(ctx, &foundRepos, query, args...)

	return &foundRepos, err
}

func (omaRepo *OmaRepositoryImpl) Update(ctx context.Context, id int, data *OmaRepository) (*OmaRepository, error) {
	qb := sq.Update("repositories")

	if data.FileName.Valid {
		qb = qb.Set("filename", data.FileName.String)
	}
	if data.CachedText.Valid {
		qb = qb.Set("cached_text", data.CachedText.String)
	}

	qb = qb.Where(squirrel.Eq{"id": strconv.Itoa(id)}).Suffix("returning *")

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error while updating: %v\n", err)
	}

	updatedRepo := &OmaRepository{}
	err = omaRepo.db.GetContext(ctx, updatedRepo, query, args...)

	return updatedRepo, err
}

func (omaRepo *OmaRepositoryImpl) Delete(ctx context.Context, id int) error {
	query := `update repositories set deleted_at = now() where id = $1`
	_, err := omaRepo.db.ExecContext(ctx, query, id)
	return err
}
