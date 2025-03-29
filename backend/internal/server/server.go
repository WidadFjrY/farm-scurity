package server

import (
	"farm-scurity/internal/app"
	"farm-scurity/internal/di"
	"farm-scurity/internal/middleware"
	"farm-scurity/pkg/helper"
	"farm-scurity/router"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Run() {
	db := app.NewDB()
	validator := validator.New()

	gin := gin.Default()
	gin.Use(middleware.ErrorHandling())

	historiDI := di.HistoryDI(db)
	UserDI := di.UserDI(db)
	DeviceDI := di.DeviceDI(db, validator)

	router.HistoryRouter(gin, historiDI)
	router.UserRouter(gin, UserDI)
	router.DeviceRouter(gin, DeviceDI)

	gin.Static("api/public/images", "./public/images")

	err := gin.Run(":8080")
	helper.Err(err)
}
