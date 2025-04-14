package service

import (
	"context"
	"farm-scurity/domain/model"
	"farm-scurity/domain/web"
	"farm-scurity/internal/repository"
	"farm-scurity/pkg/helper"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type DeviceServiceImpl struct {
	DB        *gorm.DB
	Validator *validator.Validate
	Repo      repository.DeviceRepository
}

func NewDaviceService(db *gorm.DB, validator *validator.Validate, repo repository.DeviceRepository) DeviceService {
	return &DeviceServiceImpl{DB: db, Validator: validator, Repo: repo}
}

func (serv *DeviceServiceImpl) GetDevices(ctx context.Context) []web.Device {
	var devices []web.Device

	helper.Err(serv.DB.Transaction(func(tx *gorm.DB) error {
		for _, device := range serv.Repo.GetDevices(ctx, tx) {
			devices = append(devices, web.Device{
				ID:       device.ID,
				Location: device.Location,
				IsActive: device.IsActive,
			})
		}
		return nil
	}))

	return devices
}

func (serv *DeviceServiceImpl) SetIsActive(ctx context.Context, request web.SetIsActiveRequest) {
	err := serv.Validator.Struct(&request)
	helper.Err(err)

	helper.Err(serv.DB.Transaction(func(tx *gorm.DB) error {
		serv.Repo.SetIsActive(ctx, tx, model.Device{
			ID:       request.ID,
			IsActive: *request.IsActive,
		})
		return nil
	}))
}
