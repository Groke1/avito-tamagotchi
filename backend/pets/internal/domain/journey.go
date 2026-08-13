package domain

type JourneyReward struct {
	RewardXP    int32   `json:"reward_xp"`
	RewardCoins int32   `json:"reward_coins"`
	RewardCode  *string `json:"reward_code"`
}

type JourneyResult struct {
	Location string        `json:"location"`
	Events   []string      `json:"events"`
	Reward   JourneyReward `json:"reward"`
}

type JourneyStory struct {
	Title string `json:"title"`
	Story string `json:"story"`
	// Teaser string `json:"teaser"`
}

type JourneyGenerationInput struct {
	Journey JourneyResult
	Memory  []PetTrip // ожидается 2–3 последних, не больше
}
