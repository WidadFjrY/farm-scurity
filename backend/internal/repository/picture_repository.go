package repository

import (
	"context"
	"farm-scurity/domain/model"

	"gorm.io/gorm"
)

type PictureRepository interface {
	Save(ctx context.Context, tx *gorm.DB, picture model.Picture)
	GetAll(ctx context.Context, tx *gorm.DB) []model.Picture
	GetById(ctx context.Context, tx *gorm.DB, pictureId string) model.Picture
	GetLastPicture(ctx context.Context, tx *gorm.DB) model.Picture
}
