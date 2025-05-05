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
	GetNextOmaRepoId(ctx context.Context) (int, error)
	Create(ctx context.Context, data *OmaRepository) (*OmaRepository, error)
	Get(ctx context.Context, id int) (*OmaRepository, error)
	GetMany(ctx context.Context, ids []int) (*[]OmaRepository, error)
	GetByFilename(ctx context.Context, filename sql.NullString, omaRepoId int) (*OmaRepository, error)
	Update(ctx context.Context, id int, data *OmaRepository) (*OmaRepository, error)
	Delete(ctx context.Context, id int) error
	GetAllByRepoId(ctx context.Context, repoId int) ([]OmaRepository, error)
}

type OmaRepositoryImpl struct {
	db *sqlx.DB
}

func NewOmaRepository(db *sqlx.DB) *OmaRepositoryImpl {
	return &OmaRepositoryImpl{db: db}
}

func (omaRepo *OmaRepositoryImpl) GetNextOmaRepoId(ctx context.Context) (int, error) {
	nextIdQuery, _, err := sq.Select("max(oma_repo_id)").From("repositories").ToSql()
	var id int

	if err != nil {
		return -1, fmt.Errorf("error while getting the next id for oma repository:\n%v", err)
	}

	err = omaRepo.db.GetContext(ctx, &id, nextIdQuery)

	if err != nil {
		id = 1
	} else {
		id++
	}

	return id, nil
}

func (omaRepo *OmaRepositoryImpl) Create(ctx context.Context, data *OmaRepository) (*OmaRepository, error) {
	if !data.FileName.Valid {
		return nil, fmt.Errorf("illogical attempt of creating a repository")
	}

	query, args, err := sq.Insert("repositories").SetMap(map[string]any{
		"filename":    data.FileName.String,
		"cached_text": data.CachedText.String,
		"oma_repo_id": data.OmaRepoId,
	}).Suffix("returning *").ToSql()

	createdRepo := &OmaRepository{}

	err = omaRepo.db.GetContext(ctx, createdRepo, query, args...)

	if err != nil {
		return nil, fmt.Errorf("error while creating an oma repository:\n%v", err)
	}

	return createdRepo, err
}

func (omaRepo *OmaRepositoryImpl) Get(ctx context.Context, id int) (*OmaRepository, error) {
	query, _, err := sq.Select("repositories").Where(squirrel.Eq{"id": id}).ToSql()

	if err != nil {
		return nil, fmt.Errorf("error while getting:\n%v", err)
	}

	foundRepo := &OmaRepository{}

	err = omaRepo.db.SelectContext(ctx, foundRepo, query, id)
	return foundRepo, err

}

func (omaRepo *OmaRepositoryImpl) GetByFilename(ctx context.Context, filename sql.NullString, omaRepoId int) (*OmaRepository, error) {
	if !filename.Valid {
		return nil, fmt.Errorf("you can not search for a nil file name\n")
	}

	query, args, err := sq.Select("*").From("repositories").Where(squirrel.Eq{
		"filename":    filename.String,
		"oma_repo_id": omaRepoId,
	}).ToSql()

	if err != nil {
		return nil, fmt.Errorf("error while generating the GetByFilename query:\n%v", err)
	}

	foundRepo := []OmaRepository{}

	err = omaRepo.db.SelectContext(ctx, &foundRepo, query, args...)

	if err != nil {
		return nil, err
	}

	if len(foundRepo) != 1 {
		return nil, fmt.Errorf("could not find a repo for given file-repository ID combination, repository ID: %v, filename: %v", omaRepoId, filename.String)
	}

	return &foundRepo[0], err
}

func (omaRepo *OmaRepositoryImpl) GetMany(ctx context.Context, ids []int) (*[]OmaRepository, error) {
	query, args, err := sq.Select("*").From("repositories").Where(squirrel.Eq{"id": ids}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("error at GetMany:\n%v", err)
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
		return nil, fmt.Errorf("error while updating:\n%v", err)
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

func (omaRepo *OmaRepositoryImpl) GetAllByRepoId(ctx context.Context, repoId int) ([]OmaRepository, error) {
	query, args, err := sq.Select("*").From("repositories").Where(squirrel.Eq{
		"oma_repo_id": repoId,
	}).ToSql()

	if err != nil {
		return nil, fmt.Errorf("error while building GetAllByRepoId query, error:\n%w", err)
	}

	var foundRepos []OmaRepository

	err = omaRepo.db.SelectContext(ctx, &foundRepos, query, args...)

	if err != nil {
		return nil, fmt.Errorf("error while finding the repositories for: %v, error:\n%w", repoId, err)
	}

	return foundRepos, nil
}
