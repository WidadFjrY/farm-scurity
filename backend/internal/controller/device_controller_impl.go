package controller

import (
	"farm-scurity/domain/web"
	"farm-scurity/internal/broker"
	"farm-scurity/internal/service"
	"farm-scurity/pkg/exception"
	"farm-scurity/pkg/helper"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DeviceControllerImpl struct {
	HistoryServ service.HistoryService
	PictureServ service.PictureService
	DeviceServ  service.DeviceService
}

func NewDeviceController(historyServ service.HistoryService, pictureServ service.PictureService, deviceServ service.DeviceService) DeviceController {
	return &DeviceControllerImpl{HistoryServ: historyServ, PictureServ: pictureServ, DeviceServ: deviceServ}
}

func (control *DeviceControllerImpl) Upload(ctx *gin.Context) {
	pictureId := ctx.Params.ByName("picture_id")

	contentType := ctx.GetHeader("Content-Type")
	if contentType != "image/jpeg" {
		panic(exception.NewBadRequestError("Invalid Content-Type! Expected image/jpeg"))
	}

	body, err := io.ReadAll(ctx.Request.Body)
	helper.Err(err)

	filePath := fmt.Sprintf("public/images/%s.jpg", pictureId)

	err = os.WriteFile(filePath, body, 0644)
	helper.Err(err)

	control.PictureServ.Save(ctx.Request.Context(), filePath, pictureId)

	helper.Response(ctx, http.StatusOK, "Ok", "success")
}

func (control *DeviceControllerImpl) MotionDetected(ctx *gin.Context) {
	var request web.MotionDetectedRequest
	pictureId := ctx.Params.ByName("picture_id")

	err := ctx.ShouldBind(&request)
	helper.Err(err)

	if request.MotionDetected {
		control.HistoryServ.Create(ctx.Request.Context(), "Gerakan Terdeteksi", pictureId, fmt.Sprintf("Gerakan Terdeteksi dari Sensor PIR dengan ID  %s", request.DeviceId))
		broker.MQTTRequest(web.MQTTRequest{
			ClientId: "SERVER",
			Topic:    "broker/farm-security/notification",
			Payload:  fmt.Sprintf("Gerakan Terdeteksi dari Sensor PIR dengan ID %s", request.DeviceId),
			MsgResp:  "ok",
		})
	}

	helper.Response(ctx, http.StatusOK, "Ok", "")
}

func (control *DeviceControllerImpl) GetDevices(ctx *gin.Context) {
	helper.Response(ctx, http.StatusOK, "Ok", control.DeviceServ.GetDevices(ctx.Request.Context()))
}

func (control *DeviceControllerImpl) SetIsActive(ctx *gin.Context) {
	var request web.SetIsActiveRequest
	ctx.ShouldBind(&request)

	isSuccess, _ := broker.MQTTRequest(web.MQTTRequest{
		ClientId: "SERVER",
		Topic:    "broker/farm-security",
		Payload:  fmt.Sprintf("DISABLE SENSOR ID: %s, IsActive: %s", request.ID, strconv.FormatBool(*request.IsActive)),
		MsgResp:  "ok",
	})

	if isSuccess {
		control.DeviceServ.SetIsActive(ctx.Request.Context(), request)
		helper.Response(ctx, http.StatusOK, "Ok", "")
		return
	}

	panic(exception.NewBadRequestError("Gagal melakukan operasi"))
}
