package service

import (
	"context"
	"farm-scurity/domain/web"
)

type HistoryService interface {
	Create(ctx context.Context, operation string, pictureId string, description string)
	GetAll(ctx context.Context) []web.HistoryResponse
	GetById(ctx context.Context, historyId string) web.HistoryResponse
	UpdateIsRead(ctx context.Context, historyId string)
	DeleteById(ctx context.Context, historyId string)
	DeleteAll(ctx context.Context)
}
