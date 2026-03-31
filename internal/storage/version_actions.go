package storage

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Keys string

const (
	AdditionKey Keys = "addition"
	DeletionKey Keys = "deletion"
)

type VersionActions struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `gorm:"not null" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"not null" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
	Pos       int            `gorm:"not null" json:"pos"`
	ActionKey Keys           `gorm:"not null" json:"actionKey"`
	VersionId int            `gorm:"not null" json:"versionId"`
	Content   string         `gorm:"not null" json:"content"`
}

type VersionActionsRepository interface {
	Create(ctx context.Context, data *VersionActions) (*VersionActions, error)
	GetByVersionId(ctx context.Context, versionId int) ([]VersionActions, error)
	DeleteByVersionId(ctx context.Context, versionId int) error
}

type VersionActionsRepositoryImpl struct {
	db *gorm.DB
}

func NewVersionActionsRepository(db *gorm.DB) *VersionActionsRepositoryImpl {
	return &VersionActionsRepositoryImpl{db: db}
}

func (r *VersionActionsRepositoryImpl) Create(ctx context.Context, data *VersionActions) (*VersionActions, error) {
	result := r.db.WithContext(ctx).Create(data)
	return data, result.Error
}

func (r *VersionActionsRepositoryImpl) GetByVersionId(ctx context.Context, versionId int) ([]VersionActions, error) {
	var actions []VersionActions
	result := r.db.WithContext(ctx).Where("version_id = ?", versionId).Find(&actions)
	return actions, result.Error
}

func (r *VersionActionsRepositoryImpl) DeleteByVersionId(ctx context.Context, versionId int) error {
	return r.db.WithContext(ctx).Where("version_id = ?", versionId).Delete(&VersionActions{}).Error
}
