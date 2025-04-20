package storage

type RepositoryContainer struct {
	OmaRepository            OmaRepoRepository
	VersionsRepository       VersionRepository
	VersionActionsRepository VersionActionsRepository
}
