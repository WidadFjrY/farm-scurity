package controller

import "github.com/gin-gonic/gin"

type DeviceController interface {
	Upload(ctx *gin.Context)
	MotionDetected(ctx *gin.Context)
}
