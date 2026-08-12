package domain

import "context"

type JourneyResult struct {
	Location  string
	Events    []string
	Reward    JourneyReward
	FinalMood string
}

type JourneyReward struct {
	Coins int
	Item  string
}

type JourneyStory struct {
	Title  string
	Text   string
	Teaser string
}

type JourneyStoryGenerator interface {
	Generate(ctx context.Context, journey JourneyResult) (JourneyStory, error)
}
