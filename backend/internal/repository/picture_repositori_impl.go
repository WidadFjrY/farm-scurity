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

type PictureRepositoryImpl struct{}

func NewPictureRepository() PictureRepository {
	return &PictureRepositoryImpl{}
}

func (repo *PictureRepositoryImpl) Save(ctx context.Context, tx *gorm.DB, picture model.Picture) {
	err := tx.WithContext(ctx).Create(&picture).Error
	helper.Err(err)
}

func (repo *PictureRepositoryImpl) GetAll(ctx context.Context, tx *gorm.DB) []model.Picture {
	var pictures []model.Picture

	err := tx.WithContext(ctx).Find(&pictures).Error
	helper.Err(err)

	return pictures
}

func (repo *PictureRepositoryImpl) GetById(ctx context.Context, tx *gorm.DB, pictureId string) model.Picture {
	var picture model.Picture

	err := tx.WithContext(ctx).Where("id = ?", pictureId).First(&picture).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		panic(exception.NewNotFoundError(fmt.Sprintf("picture with id %s not found", pictureId)))
	}
	helper.Err(err)

	return picture
}

func (repo *PictureRepositoryImpl) GetLastPicture(ctx context.Context, tx *gorm.DB) model.Picture {
	var picture model.Picture

	err := tx.WithContext(ctx).Order("created_at DESC").First(&picture).Error
	helper.Err(err)

	return picture
}
