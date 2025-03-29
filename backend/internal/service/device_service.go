package service

import (
	"context"
	"farm-scurity/domain/web"
)

type DeviceService interface {
	GetDevices(ctx context.Context) []web.Device
	SetIsActive(ctx context.Context, request web.SetIsActiveRequest)
}
