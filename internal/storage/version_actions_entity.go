package storage

import (
	"database/sql"
	"time"
)

type Keys string

const (
	AdditionKey Keys = "addition"
	DeletionKey Keys = "deletion"
	MoveKey     Keys = "move"
)

type VersionActions struct {
	ID        int           `db:"id"`
	CreatedAt time.Time     `db:"created_at"`
	DeletedAt sql.NullTime  `db:"deleted_at"`
	Start     sql.Null[int] `db:"start"`
	Dest      int           `db:"dest"`
	ActionKey Keys          `db:"action_key"`
	VersionId int           `db:"version_id"`
}
