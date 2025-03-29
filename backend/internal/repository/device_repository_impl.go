package repository

import (
	"context"
	"farm-scurity/domain/model"
	"farm-scurity/pkg/helper"

	"gorm.io/gorm"
)

type DeviceRepositoryImpl struct{}

func NewDeviceRepository() DeviceRepository {
	return &DeviceRepositoryImpl{}
}

func (repo *DeviceRepositoryImpl) GetDevices(ctx context.Context, tx *gorm.DB) []model.Device {
	var devices []model.Device
	helper.Err(tx.WithContext(ctx).Find(&devices).Error)
	return devices
}

func (repo *DeviceRepositoryImpl) SetIsActive(ctx context.Context, tx *gorm.DB, device model.Device) {
	helper.Err(tx.WithContext(ctx).Table("devices").Where("id = ?", device.ID).Updates(map[string]interface{}{
		"is_active": device.IsActive,
	}).Error)
}
