package di

import (
	"farm-scurity/internal/controller"
	"farm-scurity/internal/repository"
	"farm-scurity/internal/service"

	"gorm.io/gorm"
)

func HistoryDI(db *gorm.DB) controller.HistoryController {
	repo := repository.NewHistoryRepository()
	serv := service.NewHistoryService(db, repo)
	cntrl := controller.NewHistoryController(serv)

	return cntrl
}
