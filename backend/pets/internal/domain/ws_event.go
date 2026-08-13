package domain

type WsEvent string

const (
	EventPetUpdated        WsEvent = "pet.updated"
	EventLeaderboardChange WsEvent = "leaderboard.position_changed"
	EventLevelUp           WsEvent = "pet.leveled_up"
	EventStreakReward      WsEvent = "streak.rewards"
	EventTripCompleted     WsEvent = "trip.completed"
)
