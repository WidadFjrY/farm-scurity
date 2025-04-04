package controller

import (
	"farm-scurity/domain/web"
	"farm-scurity/internal/broker"
	"farm-scurity/internal/service"
	"farm-scurity/pkg/exception"
	"farm-scurity/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserControllerImpl struct {
	PictureServ service.PictureService
}

func NewUserController(pictureServ service.PictureService) UserController {
	return &UserControllerImpl{PictureServ: pictureServ}
}

func (controller *UserControllerImpl) Capture(ctx *gin.Context) {
	mqttRequest := web.MQTTRequest{
		ClientId: "SERVER",
		Topic:    "broker/farm-scurity",
		Payload:  "capture",
		MsgResp:  "ok",
	}

	respMQTT := broker.MQTTRequest(mqttRequest)

	if respMQTT {
		response := controller.PictureServ.GetLastPicture(ctx.Request.Context())
		helper.Response(ctx, http.StatusOK, "Ok", response)
	} else {
		panic(exception.NewBadRequestError("failed to capture"))
	}

}

func (controller *UserControllerImpl) TurnOn(ctx *gin.Context) {
	mqttRequest := web.MQTTRequest{
		ClientId: "SERVER",
		Topic:    "broker/farm-scurity",
		Payload:  "ALARM_ON",
		MsgResp:  "ok",
	}

	respMQTT := broker.MQTTRequest(mqttRequest)
	if respMQTT {
		helper.Response(ctx, http.StatusOK, "Ok", "")
	} else {
		panic(exception.NewBadRequestError("failed to activate speaker"))
	}

}

func (controller *UserControllerImpl) TurnOff(ctx *gin.Context) {
	mqttRequest := web.MQTTRequest{
		ClientId: "SERVER",
		Topic:    "broker/farm-scurity",
		Payload:  "ALARM_OFF",
		MsgResp:  "ok",
	}

	respMQTT := broker.MQTTRequest(mqttRequest)
	if respMQTT {
		helper.Response(ctx, http.StatusOK, "Ok", "")
	} else {
		panic(exception.NewBadRequestError("failed to deactivate speaker"))
	}

}
