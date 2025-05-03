package web

type MotionDetectedRequest struct {
	DeviceId       string `json:"device_id"`
	MotionDetected bool   `json:"motion_detected"`
}

type Device struct {
	ID       string `json:"id"`
	Location string `json:"location"`
	IsActive bool   `json:"is_active"`
}

type SetIsActiveRequest struct {
	ID       string `json:"id" validate:"required"`
	IsActive *bool  `json:"is_active" validate:"required"`
}
