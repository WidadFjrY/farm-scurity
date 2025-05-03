package web

type HistoryResponse struct {
	ID          string `json:"id"`
	Operation   string `json:"operation"`
	Description string `json:"description"`
	IsRead      bool   `json:"is_read"`
	PictureID   string `json:"picture_id"`
	CreatedAt   string `json:"created_at"`
}
