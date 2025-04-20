package storage

import "github.com/Masterminds/squirrel"

var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
