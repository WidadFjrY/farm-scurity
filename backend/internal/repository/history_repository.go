package repository

import (
	"context"
	"farm-scurity/domain/model"

	"gorm.io/gorm"
)

type HistoryRepository interface {
	Create(ctx context.Context, tx *gorm.DB, history model.History)
	GetAll(ctx context.Context, tx *gorm.DB) []model.History
	GetById(ctx context.Context, tx *gorm.DB, historyId string) model.History
	UpdateIsRead(ctx context.Context, tx *gorm.DB, historyId string)
	DeleteById(ctx context.Context, tx *gorm.DB, historyId string)
	DeleteAll(ctx context.Context, tx *gorm.DB)
}
