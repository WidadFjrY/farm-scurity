package controller

import (
	"context"
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
	pictureId := ctx.Param("picture_id")

	file, header, err := ctx.Request.FormFile("imageFile")
	if err != nil {
		helper.Err(err)
		return
	}
	defer file.Close()

	if header.Header.Get("Content-Type") != "image/jpeg" {
		panic(exception.NewBadRequestError("Invalid Content-Type! Expected image/jpeg"))
	}

	filePath := fmt.Sprintf("public/images/%s.jpg", pictureId)
	out, err := os.Create(filePath)
	if err != nil {
		helper.Err(err)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		helper.Err(err)
		return
	}

	control.PictureServ.Save(ctx.Request.Context(), filePath, pictureId)
	helper.Response(ctx, http.StatusOK, "Ok", "success")
}

func (control *DeviceControllerImpl) MotionDetected(ctx *gin.Context) {
	var request web.MotionDetectedRequest
	pictureId := ctx.Params.ByName("picture_id")

	err := ctx.ShouldBind(&request)
	helper.Err(err)

	if request.MotionDetected {
		go func(pictureId string) {
			control.HistoryServ.Create(context.Background(), "Gerakan Terdeteksi", pictureId, fmt.Sprintf("Gerakan Terdeteksi dari Sensor PIR dengan ID  %s", request.DeviceId))
		}(pictureId)
		broker.MQTTRequest(web.MQTTRequest{
			ClientId: "SERVER",
			Topic:    "bido_dihara/broker/farm-security/notification",
			Payload:  fmt.Sprintf("Gerakan Terdeteksi dari Sensor PIR dengan ID %s", request.DeviceId),
			MsgResp:  "ok",
		}, true)
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
		Topic:    "bido_dihara/broker/farm-security",
		Payload:  fmt.Sprintf("DISABLE SENSOR ID: %s, IsActive: %s", request.ID, strconv.FormatBool(*request.IsActive)),
		MsgResp:  "ok",
	}, false)

	if isSuccess {
		control.DeviceServ.SetIsActive(ctx.Request.Context(), request)
		helper.Response(ctx, http.StatusOK, "Ok", "")
		return
	}

	panic(exception.NewBadRequestError("Gagal melakukan operasi"))
}
