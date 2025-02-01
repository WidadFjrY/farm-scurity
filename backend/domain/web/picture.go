package web

import "time"

type GetPictureResponse struct {
	ID                    string    `json:"id"`
	URL                   string    `json:"url"`
	IsFromMotionDetection bool      `json:"is_from_motion_detection"`
	CreatedAt             time.Time `json:"created_at"`
}
