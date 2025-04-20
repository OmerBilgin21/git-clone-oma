package postgres

import "github.com/Masterminds/squirrel"

var Sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
