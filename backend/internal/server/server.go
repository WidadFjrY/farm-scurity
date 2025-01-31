package server

import (
	"farm-scurity/internal/app"
	"farm-scurity/internal/controller"
	"farm-scurity/internal/di"
	"farm-scurity/internal/middleware"
	"farm-scurity/pkg/helper"
	"farm-scurity/router"

	"github.com/gin-gonic/gin"
)

func Run() {
	db := app.NewDB()

	gin := gin.Default()
	gin.Use(middleware.ErrorHandling())

	historiDI := di.HistoryDI(db)

	userController := controller.NewUserController()

	router.HistoryRouter(gin, historiDI)
	router.UserRouter(gin, userController)

	err := gin.Run("localhost:8080")
	helper.Err(err)
}
