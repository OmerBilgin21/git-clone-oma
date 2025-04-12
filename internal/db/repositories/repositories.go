package repositories

import (
	"oma/internal/db/postgres"
)

type RepositoryContainer struct {
	OmaRepository postgres.OmaRepositoryRepository
}
