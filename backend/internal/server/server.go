package server

import (
	"farm-scurity/internal/app"
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
	UserDI := di.UserDI(db)
	DeviceDI := di.DeviceDI(db)

	router.HistoryRouter(gin, historiDI)
	router.UserRouter(gin, UserDI)
	router.DeviceRouter(gin, DeviceDI)

	gin.Static("api/public/", "./public/images")

	err := gin.Run("localhost:8080")
	helper.Err(err)
}
