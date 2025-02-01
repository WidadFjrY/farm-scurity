package router

import (
	"farm-scurity/internal/controller"

	"github.com/gin-gonic/gin"
)

func DeviceRouter(router *gin.Engine, control controller.DeviceController) {
	router.POST("/api/upload", control.Upload)
	router.POST("/api/motion-detected", control.MotionDetected)
}
