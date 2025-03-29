package di

import (
	"farm-scurity/internal/controller"
	"farm-scurity/internal/repository"
	"farm-scurity/internal/service"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func DeviceDI(db *gorm.DB, validator *validator.Validate) controller.DeviceController {
	histRepo := repository.NewHistoryRepository()
	histServ := service.NewHistoryService(db, histRepo)

	pictRepo := repository.NewPictureRepository()
	pictServ := service.NewPictureRepository(db, pictRepo)

	devcRepo := repository.NewDeviceRepository()
	devcServ := service.NewDaviceService(db, validator, devcRepo)

	cntrl := controller.NewDeviceController(histServ, pictServ, devcServ)

	return cntrl
}
