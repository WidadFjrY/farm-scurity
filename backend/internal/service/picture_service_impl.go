package service

import (
	"context"
	"farm-scurity/domain/model"
	"farm-scurity/domain/web"
	"farm-scurity/internal/repository"
	"farm-scurity/pkg/helper"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

type PictureServiceImpl struct {
	DB   *gorm.DB
	Repo repository.PictureRepository
}

func NewPictureRepository(db *gorm.DB, repo repository.PictureRepository) PictureService {
	return &PictureServiceImpl{DB: db, Repo: repo}
}

func (serv *PictureServiceImpl) Save(ctx context.Context, filePath string, pictureId string) {
	rand.NewSource(time.Now().Unix())
	err := serv.DB.Transaction(func(tx *gorm.DB) error {
		serv.Repo.Save(ctx, tx, model.Picture{
			ID:  pictureId,
			URL: filePath,
		})
		return nil
	})
	helper.Err(err)
}

func (serv *PictureServiceImpl) GetAll(ctx context.Context) []web.GetPictureResponse {
	var pictures []web.GetPictureResponse

	err := serv.DB.Transaction(func(tx *gorm.DB) error {
		for _, picture := range serv.Repo.GetAll(ctx, tx) {
			pictures = append(pictures, web.GetPictureResponse{
				ID:        picture.ID,
				URL:       picture.URL,
				CreatedAt: picture.CreatedAt,
			})
		}
		return nil
	})
	helper.Err(err)

	return pictures
}

func (serv *PictureServiceImpl) GetById(ctx context.Context, pictureId string) web.GetPictureResponse {
	var picture model.Picture

	err := serv.DB.Transaction(func(tx *gorm.DB) error {
		picture = serv.Repo.GetById(ctx, tx, pictureId)
		return nil
	})
	helper.Err(err)

	return web.GetPictureResponse{
		ID:        picture.ID,
		URL:       picture.URL,
		CreatedAt: picture.CreatedAt,
	}
}

func (serv *PictureServiceImpl) GetLastPicture(ctx context.Context) web.GetPictureResponse {
	var picture model.Picture

	err := serv.DB.Transaction(func(tx *gorm.DB) error {
		picture = serv.Repo.GetLastPicture(ctx, tx)
		return nil
	})
	helper.Err(err)

	return web.GetPictureResponse{
		ID:        picture.ID,
		URL:       picture.URL,
		CreatedAt: picture.CreatedAt,
	}
}
