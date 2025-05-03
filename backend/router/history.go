package router

import (
	"farm-scurity/internal/controller"

	"github.com/gin-gonic/gin"
)

func HistoryRouter(router *gin.Engine, cntrl controller.HistoryController) {
	router.GET("/api/histories/", cntrl.GetAll)
	router.DELETE("/api/history/:historyId", cntrl.DeleteById)
	router.PUT("/api/history/:historyId", cntrl.UpdateIsRead)
}
