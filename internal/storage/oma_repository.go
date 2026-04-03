package storage

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type OmaRepository struct {
	ID         int            `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `gorm:"not null" json:"createdAt"`
	UpdatedAt  time.Time      `gorm:"not null" json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deletedAt"`
	OmaRepoId  int            `gorm:"not null" json:"omaRepoId"`
	FileName   *string        `gorm:"default:null" json:"fileName"`
	CachedText *string        `gorm:"default:null" json:"cachedText"`
}

type OmaRepositoryImpl struct {
	db *gorm.DB
}

func NewOmaRepository(db *gorm.DB) *OmaRepositoryImpl {
	return &OmaRepositoryImpl{db: db}
}

func (r *OmaRepositoryImpl) GetNextOmaRepoId(ctx context.Context) (int, error) {
	var maxId int
	r.db.WithContext(ctx).Model(&OmaRepository{}).Select("COALESCE(MAX(oma_repo_id), 0)").Scan(&maxId)
	return maxId + 1, nil
}

func (r *OmaRepositoryImpl) Create(ctx context.Context, data *OmaRepository) (*OmaRepository, error) {
	if data.FileName == nil {
		return nil, fmt.Errorf("illogical attempt of creating a repository")
	}

	err := r.db.WithContext(ctx).Create(data).Error
	return data, err
}

func (r *OmaRepositoryImpl) Get(ctx context.Context, id int) (*OmaRepository, error) {
	var repo OmaRepository
	err := r.db.WithContext(ctx).First(&repo, id).Error
	return &repo, err
}

func (r *OmaRepositoryImpl) GetByFilename(ctx context.Context, filename string, omaRepoId int) (*OmaRepository, error) {
	var repo OmaRepository
	err := r.db.WithContext(ctx).Where("file_name = ? AND oma_repo_id = ?", filename, omaRepoId).First(&repo).Error
	return &repo, err
}

func (r *OmaRepositoryImpl) GetMany(ctx context.Context, ids []int) (*[]OmaRepository, error) {
	var repos []OmaRepository
	err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&repos).Error
	return &repos, err
}

func (r *OmaRepositoryImpl) Update(ctx context.Context, id int, data *OmaRepository) (*OmaRepository, error) {
	updates := map[string]any{}
	if data.FileName != nil {
		updates["file_name"] = *data.FileName
	}
	if data.CachedText != nil {
		updates["cached_text"] = *data.CachedText
	}

	err := r.db.WithContext(ctx).Model(&OmaRepository{}).Where("id = ?", id).Updates(updates).Error
	if err != nil {
		return nil, err
	}

	var updated OmaRepository
	err = r.db.WithContext(ctx).First(&updated, id).Error
	return &updated, err
}

func (r *OmaRepositoryImpl) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&OmaRepository{}, id).Error
}

func (r *OmaRepositoryImpl) GetAllByRepoId(ctx context.Context, repoId int) (*[]OmaRepository, error) {
	var repos *[]OmaRepository
	err := r.db.WithContext(ctx).Where("oma_repo_id = ?", repoId).Find(&repos).Error
	return repos, err
}
