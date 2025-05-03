package controller

import (
	"farm-scurity/internal/service"
	"farm-scurity/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HistoryControllerImpl struct {
	Serv service.HistoryService
}

func NewHistoryController(serv service.HistoryService) HistoryController {
	return &HistoryControllerImpl{Serv: serv}
}

func (controller *HistoryControllerImpl) GetAll(ctx *gin.Context) {

	response := controller.Serv.GetAll(ctx.Request.Context())
	helper.Response(ctx, http.StatusOK, "Ok", response)
}

func (controller *HistoryControllerImpl) GetById(ctx *gin.Context) {
	historyId := ctx.Params.ByName("historyId")

	response := controller.Serv.GetById(ctx.Request.Context(), historyId)
	helper.Response(ctx, http.StatusOK, "Ok", response)
}

func (controller *HistoryControllerImpl) DeleteById(ctx *gin.Context) {
	historyId := ctx.Params.ByName("historyId")

	controller.Serv.DeleteById(ctx.Request.Context(), historyId)
	helper.Response(ctx, http.StatusOK, "Ok", "Deleted")
}

func (controller *HistoryControllerImpl) UpdateIsRead(ctx *gin.Context) {
	historyId := ctx.Params.ByName("historyId")

	controller.Serv.UpdateIsRead(ctx.Request.Context(), historyId)
	helper.Response(ctx, http.StatusOK, "Ok", "updated")
}
