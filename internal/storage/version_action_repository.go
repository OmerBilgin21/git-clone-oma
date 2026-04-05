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

type VersionActionsRepository struct {
	db *gorm.DB
}

func NewVersionActionsRepository(db *gorm.DB) *VersionActionsRepository {
	return &VersionActionsRepository{db: db}
}

func (r *VersionActionsRepository) Create(ctx context.Context, data *VersionActions) error {
	return r.db.WithContext(ctx).Create(data).Error
}

func (r *VersionActionsRepository) GetByVersionId(ctx context.Context, versionId int) (*[]VersionActions, error) {
	var actions []VersionActions
	err := r.db.WithContext(ctx).Where("version_id = ?", versionId).Order("id asc").Find(&actions).Error
	return &actions, err
}

func (r *VersionActionsRepository) DeleteByVersionId(ctx context.Context, versionId int) error {
	return r.db.WithContext(ctx).Where("version_id = ?", versionId).Delete(&VersionActions{}).Error
}
