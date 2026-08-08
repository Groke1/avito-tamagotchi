package entity

type (
	UpdateCoinsRequest struct {
		UserID string `json:"user_id"`
		Delta  int    `json:"delta"`
	}
	UpdateCoinsResponse struct {
		UserID string `json:"user_id"`
		Coins  int64  `json:"coins"`
	}
)
