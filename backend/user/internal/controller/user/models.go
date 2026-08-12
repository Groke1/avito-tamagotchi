package user

type profileResponse struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Coins    uint64 `json:"coins"`
	Streak   uint64 `json:"streak"`
}

type usernamesRequest struct {
	UserIDs []string `json:"user_ids"`
}

type usernameResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type usernamesResponse struct {
	Users []usernameResponse `json:"users"`
}

type updateCoinsRequest struct {
	UserID string `json:"user_id"`
	Delta  int64  `json:"delta"`
}

type updateCoinsResponse struct {
	UserID string `json:"user_id"`
	Coins  uint64 `json:"coins"`
}

type actionRequest struct {
	UserID     string `json:"user_id"`
	OccurredAt string `json:"occurred_at"`
}
