package clients

type UpdateCoinsRequest struct {
	UserID string `json:"user_id"`
	Delta  int    `json:"delta"`
}
