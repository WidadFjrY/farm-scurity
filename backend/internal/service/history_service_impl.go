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

type HistoryServiceImpl struct {
	DB   *gorm.DB
	Repo repository.HistoryRepository
}

func NewHistoryService(db *gorm.DB, repo repository.HistoryRepository) HistoryService {
	return &HistoryServiceImpl{DB: db, Repo: repo}
}

func (serv *HistoryServiceImpl) Create(ctx context.Context, operation string, pictureId string, description string) {
	rand.NewSource(time.Now().Unix())
	err := serv.DB.Transaction(func(tx *gorm.DB) error {
		serv.Repo.Create(ctx, tx, model.History{
			ID:          helper.GenerateRandomString(15),
			Description: description,
			PictureID:   &pictureId,
			Operation:   operation,
		})
		return nil
	})
	helper.Err(err)
}

func (serv *HistoryServiceImpl) GetAll(ctx context.Context) []web.HistoryResponse {
	var histories []web.HistoryResponse

	err := serv.DB.Transaction(func(tx *gorm.DB) error {
		for _, historiModel := range serv.Repo.GetAll(ctx, tx) {
			history := web.HistoryResponse{
				ID:          historiModel.ID,
				PictureID:   *historiModel.PictureID,
				Description: historiModel.Description,
				Operation:   historiModel.Operation,
				CreatedAt:   historiModel.CreatedAt.Format("15:04:01 02 Jan 2006"),
			}
			histories = append(histories, history)
		}
		return nil
	})

	helper.Err(err)

	return histories
}

func (serv *HistoryServiceImpl) GetById(ctx context.Context, historyId string) web.HistoryResponse {
	var history model.History

	err := serv.DB.Transaction(func(tx *gorm.DB) error {
		history = serv.Repo.GetById(ctx, tx, historyId)
		return nil
	})
	helper.Err(err)

	return web.HistoryResponse{
		ID:        history.ID,
		Operation: history.Operation,
		PictureID: *history.PictureID,
		CreatedAt: history.CreatedAt.Format("15:04:01 02 Jan 2006"),
	}
}

func (serv *HistoryServiceImpl) DeleteById(ctx context.Context, historyId string) {
	err := serv.DB.Transaction(func(tx *gorm.DB) error {
		serv.Repo.DeleteById(ctx, tx, historyId)
		return nil
	})
	helper.Err(err)
}
