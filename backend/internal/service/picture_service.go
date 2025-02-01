package service

import (
	"context"
	"farm-scurity/domain/web"
)

type PictureService interface {
	Save(ctx context.Context, filePath string, isFromMotionDetection bool)
	GetById(ctx context.Context, pictureId string) web.GetPictureResponse
	GetAll(ctx context.Context) []web.GetPictureResponse
	GetLastPicture(ctx context.Context) web.GetPictureResponse
}
