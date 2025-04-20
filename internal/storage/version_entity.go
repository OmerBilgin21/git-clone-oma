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

type Versions struct {
	ID           int          `db:"id"`
	CreatedAt    time.Time    `db:"created_at"`
	DeletedAt    sql.NullTime `db:"deleted_at"`
	StartX       int          `db:"start_x"`
	StartY       int          `db:"start_y"`
	EndX         int          `db:"end_x"`
	EndY         int          `db:"end_y"`
	ActionKey    Keys         `db:"action_key"`
	RepositoryId int          `db:"repository_id"`
}
