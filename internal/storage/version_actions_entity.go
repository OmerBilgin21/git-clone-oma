package storage

import (
	"database/sql"
	"time"
)

type Keys string

const (
	AdditionKey Keys = "addition"
	DeletionKey Keys = "deletion"
)

type VersionActions struct {
	ID        int          `db:"id"`
	CreatedAt time.Time    `db:"created_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
	Pos       int          `db:"pos"`
	ActionKey Keys         `db:"action_key"`
	VersionId int          `db:"version_id"`
	Content   string       `db:"content"`
}
