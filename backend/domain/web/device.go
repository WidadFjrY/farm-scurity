package web

import "time"

type MotionDetectedRequest struct {
	DeviceId       string    `json:"device_id"`
	Timestamp      time.Time `json:"timestamp"`
	MotionDetected bool      `json:"motion_detected"`
}
