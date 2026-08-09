package entity

type DailyStat struct {
	UserID  string        `json:"user_id"`
	Streak  int32         `json:"streak"`
	Tasks   []TasksStat   `json:"tasks"`
	Pet     PetStat       `json:"pet"`
	Rewards []RewardsStat `json:"rewards"`
}

type PetStat struct {
	DailyGainedXP int64 `json:"daily_gained_xp"`
}

type TasksStat struct {
}

type RewardsStat struct {
	PromoCode   string `json:"promo_code"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
