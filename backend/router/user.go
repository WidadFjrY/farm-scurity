package router

import (
	"farm-scurity/internal/controller"

	"github.com/gin-gonic/gin"
)

func UserRouter(router *gin.Engine, control controller.UserController) {
	router.POST("/api/capture", control.Capture)
	router.POST("/api/turn_on", control.TurnOn)
	router.POST("/api/turn_off", control.TurnOff)
}
