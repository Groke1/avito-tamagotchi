package http

type dailyBonusRequest struct {
	UserID string `json:"user_id"`
	Streak int32  `json:"streak"`
	Coins  int32  `json:"coins"`
}

type petDailyStatRequest struct {
	UserID string `json:"user_id"`
}

type emptyResponse struct{}
