package controller

import (
	"farm-scurity/domain/web"
	"farm-scurity/internal/broker"
	"farm-scurity/internal/service"
	"farm-scurity/pkg/exception"
	"farm-scurity/pkg/helper"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type DeviceControllerImpl struct {
	HistoryServ service.HistoryService
}

func NewDeviceController(historyServ service.HistoryService) DeviceController {
	return &DeviceControllerImpl{HistoryServ: historyServ}
}

func (control *DeviceControllerImpl) Upload(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	helper.Err(err)

	allowedExtensionImage := map[string]bool{
		"jpg":  true,
		"png":  true,
		"jpeg": true,
	}

	extension := strings.Split(file.Filename, ".")[1]
	filePath := fmt.Sprintf("public/images/%s.%s", helper.GenerateRandomString(15), extension)

	if !allowedExtensionImage[extension] {
		panic(exception.NewBadRequestError("file not supported! file must be jpg/png/jpeg"))
	}

	err = ctx.SaveUploadedFile(file, filePath)
	helper.Err(err)
}

func (control *DeviceControllerImpl) MotionDetected(ctx *gin.Context) {
	var request web.MotionDetectedRequest
	if err := ctx.ShouldBind(&request); err != nil {
		helper.Response(ctx, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	if request.MotionDetected {
		control.HistoryServ.Create(ctx.Request.Context(), "Gerakan Terdeteksi")
		broker.MQTTRequest(web.MQTTRequest{
			ClientId: "SERVER",
			Topic:    "broker/farm-scurity/notification",
			Payload:  "Gerakan Tedeteksi!",
			MsgResp:  "ok",
		})
	}
}
