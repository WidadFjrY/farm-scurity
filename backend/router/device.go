package router

import (
	"farm-scurity/internal/controller"

	"github.com/gin-gonic/gin"
)

func DeviceRouter(router *gin.Engine, control controller.DeviceController) {
	router.POST("/api/upload/:picture_id", control.Upload)
	router.POST("/api/motion-detected/:picture_id", control.MotionDetected)
	router.GET("/api/sensors/", control.GetDevices)
	router.PUT("/api/sensor/", control.SetIsActive)
}
