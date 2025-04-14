package controller

import (
	"farm-scurity/domain/web"
	"farm-scurity/internal/broker"
	"farm-scurity/internal/service"
	"farm-scurity/pkg/exception"
	"farm-scurity/pkg/helper"
	"net/http"
	"strings"

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
		Topic:    "broker/farm-security",
		Payload:  "TAKE_PHOTO",
		MsgResp:  "ok",
	}

	respMQTT, payload := broker.MQTTRequest(mqttRequest)

	var pictureId string
	parts := strings.Split(payload, "pictureId:")
	if len(parts) == 2 {
		pictureId = strings.TrimSpace(parts[1])
	}

	if respMQTT {
		helper.Response(ctx, http.StatusOK, "Ok", pictureId)
	} else {
		panic(exception.NewBadRequestError("failed to capture"))
	}

}

func (controller *UserControllerImpl) TurnOn(ctx *gin.Context) {
	mqttRequest := web.MQTTRequest{
		ClientId: "SERVER",
		Topic:    "broker/farm-security",
		Payload:  "ALARM_ON",
		MsgResp:  "ok",
	}

	respMQTT, _ := broker.MQTTRequest(mqttRequest)
	if respMQTT {
		helper.Response(ctx, http.StatusOK, "Ok", "")
	} else {
		panic(exception.NewBadRequestError("failed to activate speaker"))
	}

}

func (controller *UserControllerImpl) TurnOff(ctx *gin.Context) {
	mqttRequest := web.MQTTRequest{
		ClientId: "SERVER",
		Topic:    "broker/farm-security",
		Payload:  "ALARM_OFF",
		MsgResp:  "ok",
	}

	respMQTT, _ := broker.MQTTRequest(mqttRequest)
	if respMQTT {
		helper.Response(ctx, http.StatusOK, "Ok", "")
	} else {
		panic(exception.NewBadRequestError("failed to deactivate speaker"))
	}

}
