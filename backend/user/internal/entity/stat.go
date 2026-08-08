package entity

type DailyStat struct {
	UserID  string
	Streak  int32
	Tasks   TasksStat
	Pet     PetStat
	Rewards RewardsStat
}

type PetStat struct {
}

type TasksStat struct {
}

type RewardsStat struct {
	PromoCode   string
	Name        string
	Description string
}
