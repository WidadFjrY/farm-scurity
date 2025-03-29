package web

import "time"

type HistoryResponse struct {
	ID        string    `json:"id"`
	Operation string    `json:"operation"`
	PictureID string    `json:"picture_id"`
	CreatedAt time.Time `json:"created_at"`
}
