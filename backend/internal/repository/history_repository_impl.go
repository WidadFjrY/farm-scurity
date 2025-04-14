package repository

import (
	"context"
	"errors"
	"farm-scurity/domain/model"
	"farm-scurity/pkg/exception"
	"farm-scurity/pkg/helper"
	"fmt"

	"gorm.io/gorm"
)

type HistoryRepositoryImpl struct{}

func NewHistoryRepository() HistoryRepository {
	return &HistoryRepositoryImpl{}
}

func (repo *HistoryRepositoryImpl) Create(ctx context.Context, tx *gorm.DB, history model.History) {
	err := tx.WithContext(ctx).Create(&history).Error
	helper.Err(err)
}

func (repo *HistoryRepositoryImpl) GetAll(ctx context.Context, tx *gorm.DB) []model.History {
	var histories []model.History

	err := tx.WithContext(ctx).Order("created_at DESC").Find(&histories).Error
	helper.Err(err)

	return histories
}

func (repo *HistoryRepositoryImpl) GetById(ctx context.Context, tx *gorm.DB, historyId string) model.History {
	var history model.History

	err := tx.WithContext(ctx).Where("id = ?", historyId).Preload("Picture").First(&history).Error
	helper.Err(err)

	return history
}

func (repo *HistoryRepositoryImpl) DeleteById(ctx context.Context, tx *gorm.DB, historyId string) {
	err := tx.WithContext(ctx).Table("histories").Where("id = ?", historyId).Delete(nil).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		panic(exception.NewNotFoundError(fmt.Sprintf("history with id %s not found", historyId)))
	}

	helper.Err(err)
}
