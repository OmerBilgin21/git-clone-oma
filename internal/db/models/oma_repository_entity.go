package models

import (
	"database/sql"
	"time"
)

type OmaRepository struct {
	ID         int            `db:"id"`
	CreatedAt  time.Time      `db:"created_at"`
	DeletedAt  sql.NullTime   `db:"deleted_at"`
	FileName   sql.NullString `db:"filename"`
	CachedText sql.NullString `db:"cached_text"`
}
