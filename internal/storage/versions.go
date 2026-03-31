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

type VersionRepository interface {
	Create(ctx context.Context, data *Versions) (*Versions, error)
	Get(ctx context.Context, id int) (*Versions, error)
	Delete(ctx context.Context, id int) error
	GetLatestXByRepoId(ctx context.Context, repoId, x int) ([]Versions, error)
	GetMaxVersionNumberForRepo(ctx context.Context, repoId int) (int, error)
	GetAllDistinctByRepoId(ctx context.Context, repoId int) ([]Versions, error)
	GetAllByRepoId(ctx context.Context, repoId int) ([]Versions, error)
}

type VersionRepositoryImpl struct {
	db *gorm.DB
}

func NewVersionRepository(db *gorm.DB) *VersionRepositoryImpl {
	return &VersionRepositoryImpl{db: db}
}

func (r *VersionRepositoryImpl) GetMaxVersionNumberForRepo(ctx context.Context, repoId int) (int, error) {
	var maxId int
	r.db.WithContext(ctx).Model(&Versions{}).Select("COALESCE(MAX(version_id), 0)").
		Where("repository_id = ?", repoId).Scan(&maxId)
	return maxId, nil
}

func (r *VersionRepositoryImpl) Create(ctx context.Context, data *Versions) (*Versions, error) {
	maxId, err := r.GetMaxVersionNumberForRepo(ctx, data.RepositoryId)
	if err != nil {
		return nil, fmt.Errorf("error getting max version number:\n%w", err)
	}

	data.VersionId = maxId + 1

	result := r.db.WithContext(ctx).Create(data)
	return data, result.Error
}

func (r *VersionRepositoryImpl) Get(ctx context.Context, id int) (*Versions, error) {
	var v Versions
	result := r.db.WithContext(ctx).First(&v, id)
	return &v, result.Error
}

func (r *VersionRepositoryImpl) GetLatestXByRepoId(ctx context.Context, repoId, x int) ([]Versions, error) {
	latestVersionId, err := r.GetMaxVersionNumberForRepo(ctx, repoId)
	if err != nil {
		return nil, fmt.Errorf("error finding latest version for repository %v:\n%w", repoId, err)
	}

	var versions []Versions
	result := r.db.WithContext(ctx).
		Where("repository_id = ? AND version_id <= ?", repoId, latestVersionId-x).
		Find(&versions)
	return versions, result.Error
}

func (r *VersionRepositoryImpl) GetAllDistinctByRepoId(ctx context.Context, repoId int) ([]Versions, error) {
	var versions []Versions
	result := r.db.WithContext(ctx).Where("repository_id = ?", repoId).Group("version_id").Find(&versions)
	return versions, result.Error
}

func (r *VersionRepositoryImpl) GetAllByRepoId(ctx context.Context, repoId int) ([]Versions, error) {
	var versions []Versions
	result := r.db.WithContext(ctx).Where("repository_id = ?", repoId).Find(&versions)
	return versions, result.Error
}

func (r *VersionRepositoryImpl) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&Versions{}, id).Error
}
