package controller

import (
	"farm-scurity/domain/web"
	"farm-scurity/internal/broker"
	"farm-scurity/pkg/exception"
	"farm-scurity/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserControllerImpl struct{}

func NewUserController() UserController {
	return &UserControllerImpl{}
}

func (controller *UserControllerImpl) Capture(ctx *gin.Context) {
	mqttRequest := web.MQTTRequest{
		ClientId: "SERVER",
		Topic:    "broker/farm-scurity",
		Payload:  "capture",
		MsgResp:  "ok",
	}

	response := broker.MQTTRequest(mqttRequest)
	if response {
		helper.Response(ctx, http.StatusOK, "Ok", "")
	}

	panic(exception.NewBadRequestError("failed to capture"))
}

func (controller *UserControllerImpl) TurnOn(ctx *gin.Context) {
	mqttRequest := web.MQTTRequest{
		ClientId: "SERVER",
		Topic:    "broker/farm-scurity",
		Payload:  "speaker_on",
		MsgResp:  "ok",
	}

	response := broker.MQTTRequest(mqttRequest)
	if response {
		helper.Response(ctx, http.StatusOK, "Ok", "")
	}

	panic(exception.NewBadRequestError("failed to activate speaker"))
}

func (controller *UserControllerImpl) TurnOff(ctx *gin.Context) {
	mqttRequest := web.MQTTRequest{
		ClientId: "SERVER",
		Topic:    "broker/farm-scurity",
		Payload:  "speaker_off",
		MsgResp:  "ok",
	}

	response := broker.MQTTRequest(mqttRequest)
	if response {
		helper.Response(ctx, http.StatusOK, "Ok", "")
	}

	panic(exception.NewBadRequestError("failed to deactivate speaker"))
}
