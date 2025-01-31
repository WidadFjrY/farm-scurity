package web

import "time"

type HistoryResponse struct {
	ID         string    `json:"id"`
	Operation  string    `json:"oprration"`
	URLPicture string    `json:"url_picture"`
	CreatedAt  time.Time `json:"created_at"`
}
