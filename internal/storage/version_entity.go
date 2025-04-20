package storage

import (
	"database/sql"
	"time"
)

type Versions struct {
	ID           int          `db:"id"`
	CreatedAt    time.Time    `db:"created_at"`
	DeletedAt    sql.NullTime `db:"deleted_at"`
	RepositoryId int          `db:"repository_id"`
}
