package di

import (
	"farm-scurity/internal/controller"
	"farm-scurity/internal/repository"
	"farm-scurity/internal/service"

	"gorm.io/gorm"
)

func DeviceDI(db *gorm.DB) controller.DeviceController {
	histRepo := repository.NewHistoryRepository()
	histServ := service.NewHistoryService(db, histRepo)

	pictRepo := repository.NewPictureRepository()
	pictServ := service.NewPictureRepository(db, pictRepo)

	cntrl := controller.NewDeviceController(histServ, pictServ)

	return cntrl
}
