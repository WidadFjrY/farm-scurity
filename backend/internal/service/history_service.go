package service

import (
	"context"
	"farm-scurity/domain/web"
)

type HistoryService interface {
	Create(ctx context.Context, operation string)
	GetAll(ctx context.Context) []web.HistoryResponse
	GetById(ctx context.Context, historyId string) web.HistoryResponse
	DeleteById(ctx context.Context, historyId string)
}
