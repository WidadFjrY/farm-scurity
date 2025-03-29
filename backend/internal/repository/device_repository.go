package repository

import (
	"context"
	"farm-scurity/domain/model"

	"gorm.io/gorm"
)

type DeviceRepository interface {
	GetDevices(ctx context.Context, tx *gorm.DB) []model.Device
	SetIsActive(ctx context.Context, tx *gorm.DB, device model.Device)
}
