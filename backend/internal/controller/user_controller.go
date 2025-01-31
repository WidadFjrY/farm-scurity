package controller

import "github.com/gin-gonic/gin"

type UserController interface {
	Capture(ctx *gin.Context)
	TurnOn(ctx *gin.Context)
	TurnOff(ctx *gin.Context)
}
