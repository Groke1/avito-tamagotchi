package http

type dailyBonusRequest struct {
	UserID string `json:"user_id"`
	Streak int32  `json:"streak"`
}

type emptyResponse struct{}
