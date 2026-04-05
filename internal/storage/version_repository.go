package storage

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Versions struct {
	ID           int            `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `gorm:"not null" json:"createdAt"`
	UpdatedAt    time.Time      `gorm:"not null" json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deletedAt"`
	VersionId    int            `gorm:"not null" json:"versionId"`
	RepositoryId int            `gorm:"not null" json:"repositoryId"`
	Message      string         `gorm:"not null" json:"message"`
}

type VersionRepository struct {
	db *gorm.DB
}

func NewVersionRepository(db *gorm.DB) *VersionRepository {
	return &VersionRepository{db: db}
}

func (r *VersionRepository) GetMaxVersionNumberForRepo(ctx context.Context, repoId int) (int, error) {
	var maxId int
	r.db.WithContext(ctx).Model(&Versions{}).Select("COALESCE(MAX(version_id), 0)").
		Where("repository_id = ?", repoId).Scan(&maxId)
	return maxId, nil
}

func (r *VersionRepository) Create(ctx context.Context, data *Versions) (*Versions, error) {
	maxId, err := r.GetMaxVersionNumberForRepo(ctx, data.RepositoryId)
	if err != nil {
		return nil, fmt.Errorf("error getting max version number:\n%w", err)
	}

	data.VersionId = maxId + 1

	result := r.db.WithContext(ctx).Create(data)
	return data, result.Error
}

func (r *VersionRepository) Get(ctx context.Context, id int) (*Versions, error) {
	var v Versions
	result := r.db.WithContext(ctx).First(&v, id)
	return &v, result.Error
}

func (r *VersionRepository) GetLatestXByRepoId(ctx context.Context, repoId, x int) ([]Versions, error) {
	latestVersionId, err := r.GetMaxVersionNumberForRepo(ctx, repoId)
	if err != nil {
		return nil, fmt.Errorf("error finding latest version for repository %v:\n%w", repoId, err)
	}

	var versions []Versions
	result := r.db.WithContext(ctx).
		Where("repository_id = ? AND version_id > ?", repoId, latestVersionId-x).
		Find(&versions)
	return versions, result.Error
}

func (r *VersionRepository) GetAllDistinctByRepoId(ctx context.Context, repoId int) ([]Versions, error) {
	var versions []Versions
	result := r.db.WithContext(ctx).Where("repository_id = ?", repoId).Group("version_id").Find(&versions)
	return versions, result.Error
}

func (r *VersionRepository) GetAllByRepoId(ctx context.Context, repoId int) ([]Versions, error) {
	var versions []Versions
	result := r.db.WithContext(ctx).Where("repository_id = ?", repoId).Order("id asc").Find(&versions)
	return versions, result.Error
}

func (r *VersionRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&Versions{}, id).Error
}
