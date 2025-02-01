package controller

import (
	"encoding/json"
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
	PictureServ service.PictureService
}

func NewDeviceController(historyServ service.HistoryService, pictureServ service.PictureService) DeviceController {
	return &DeviceControllerImpl{HistoryServ: historyServ, PictureServ: pictureServ}
}

func (control *DeviceControllerImpl) Upload(ctx *gin.Context) {
	file, err := ctx.FormFile("image")
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
	control.PictureServ.Save(ctx.Request.Context(), filePath, false)
	helper.Response(ctx, http.StatusOK, "Ok", "success")
}

func (control *DeviceControllerImpl) MotionDetected(ctx *gin.Context) {
	metadataJSON := ctx.PostForm("metadata")
	var request web.MotionDetectedRequest

	if err := json.Unmarshal([]byte(metadataJSON), &request); err != nil {
		helper.Response(ctx, http.StatusBadRequest, "Invalid JSON metadata", nil)
		return
	}

	if request.MotionDetected {
		control.HistoryServ.Create(ctx.Request.Context(), "Gerakan Terdeteksi")
		broker.MQTTRequest(web.MQTTRequest{
			ClientId: "SERVER",
			Topic:    "broker/farm-scurity/notification",
			Payload:  "Gerakan Terdeteksi!",
			MsgResp:  "ok",
		})
	}

	file, err := ctx.FormFile("image")
	if err != nil {
		helper.Response(ctx, http.StatusBadRequest, "Image not found", nil)
		return
	}

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

	control.PictureServ.Save(ctx.Request.Context(), filePath, true)
}
