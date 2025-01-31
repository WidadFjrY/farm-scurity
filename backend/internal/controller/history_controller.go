package controller

import "github.com/gin-gonic/gin"

type HistoryController interface {
	GetAll(ctx *gin.Context)
	GetById(ctx *gin.Context)
	DeleteById(ctx *gin.Context)
}
