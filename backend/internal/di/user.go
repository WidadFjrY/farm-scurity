package di

import (
	"farm-scurity/internal/controller"
	"farm-scurity/internal/repository"
	"farm-scurity/internal/service"

	"gorm.io/gorm"
)

func UserDI(db *gorm.DB) controller.UserController {
	repo := repository.NewPictureRepository()
	serv := service.NewPictureRepository(db, repo)
	cntrl := controller.NewUserController(serv)

	return cntrl
}
