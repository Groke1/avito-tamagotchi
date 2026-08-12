package domain

import "encoding/json"

type JourneyReward struct {
	Coins int    `json:"coins"`
	Item  string `json:"item,omitempty"`
}

type JourneyResult struct {
	Location string        `json:"location"`
	Events   []string      `json:"events"`
	Reward   JourneyReward `json:"reward"`
}

type JourneyStory struct {
	Title  string `json:"title"`
	Story  string `json:"story"`
	Teaser string `json:"teaser"`
}

type PetMemory struct {
	Personality string          `json:"personality"`
	Summary     string          `json:"summary"`
	Characters  json.RawMessage `json:"characters,omitempty"`
	Storylines  json.RawMessage `json:"storylines,omitempty"`
}

type JourneyGenerationInput struct {
	Journey       JourneyResult
	Memory        PetMemory
	RecentStories []JourneyStory // ожидается 2–3 последних, не больше
}
