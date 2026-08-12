package domain

import (
	"math"
	"time"
)

type Pet struct {
	ID               int64
	UserID           string
	Name             string
	Level            int
	XP               int
	NextLevelXP      int
	Satiety          int
	Happiness        int
	CreatedAt        time.Time
	LastCalculatedAt time.Time
	LastFeedAt       *time.Time
	LastStrokeAt     *time.Time
	LastGainedXP     int
}

const (
	SatietyDropPerHour   = 5.0
	HappinessDropPerHour = 8.0

	MaxStatValue = 100
	MinStatValue = 0

	FeedCooldown   = 3 * time.Minute
	StrokeCooldown = 2 * time.Minute

	FeedSatietyIncrease     = 5
	StrokeHappinessIncrease = 3

	FeedXPAmount   = 4
	StrokeXPAmount = 10

	FeedCost = 5

	HungrySatietyThreshold  = 20
	SadnessSatietyThreshold = 50

	HappinessHighThreshold   = 60
	HappinessMediumThreshold = 45
	HappinessLowThreshold    = 20

	XPHighMultiplier     = 1.0
	XPMediumMultiplier   = 0.8
	XPLowMultiplier      = 0.6
	XPCriticalMultiplier = 0.5
)

func (p *Pet) Feed() ([]int, int, error) {
	if p.Satiety >= MaxStatValue {
		return []int{}, 0, ErrPetIsFull
	} else if p.LastFeedAt != nil && time.Since(*p.LastFeedAt) < FeedCooldown {
		return []int{}, 0, &ActionUnavailableError{
			RetryAfter: FeedCooldown - time.Since(*p.LastFeedAt),
		}
	}

	p.Satiety = min(p.Satiety+FeedSatietyIncrease, MaxStatValue)
	levelUps := p.AddXP(FeedXPAmount)
	p.LastFeedAt = new(time.Now())

	return levelUps, FeedCost, nil
}

func (p *Pet) Stroke() ([]int, error) {
	switch {
	case p.Happiness >= MaxStatValue:
		return []int{}, ErrPetIsTooHappy

	case p.Satiety < HungrySatietyThreshold:
		return []int{}, ErrPetIsTooHungry // TODO frontend integration

	case p.LastStrokeAt != nil && time.Since(*p.LastStrokeAt) < StrokeCooldown:
		return []int{}, &ActionUnavailableError{
			RetryAfter: StrokeCooldown - time.Since(*p.LastStrokeAt),
		}
	}

	p.Happiness = min(p.Happiness+StrokeHappinessIncrease, MaxStatValue)
	levelUp := p.AddXP(StrokeXPAmount)
	p.LastStrokeAt = new(time.Now())

	return levelUp, nil
}

func (p *Pet) xpMultiplier() float64 {
	switch {
	case p.Happiness >= HappinessHighThreshold:
		return XPHighMultiplier

	case p.Happiness >= HappinessMediumThreshold:
		return XPMediumMultiplier

	case p.Happiness >= HappinessLowThreshold:
		return XPLowMultiplier

	default:
		return XPCriticalMultiplier
	}
}

func (p *Pet) computeXP(amount int) int {
	return int(float64(amount) * p.xpMultiplier())
}

func (p *Pet) AddXP(baseAmount int) []int {
	amount := p.computeXP(baseAmount)

	p.XP += amount
	p.LastGainedXP = amount

	var levelUps []int

	for p.XP >= p.NextLevelXP {
		p.Level++
		p.NextLevelXP = 100 * p.Level * (p.Level + 1) / 2
		levelUps = append(levelUps, p.Level)
	}

	return levelUps
}

func (p *Pet) RecalculateState(now time.Time) bool {
	if p.LastCalculatedAt.IsZero() {
		p.LastCalculatedAt = now
		return false
	}

	duration := now.Sub(p.LastCalculatedAt)
	hoursPassed := duration.Hours()

	oldSatiety := p.Satiety
	oldHappiness := p.Happiness

	satietyDrop := int(math.Floor(hoursPassed * SatietyDropPerHour))
	if satietyDrop > 0 {
		p.Satiety = max(p.Satiety-satietyDrop, MinStatValue)
	}

	if p.Satiety < SadnessSatietyThreshold {
		happinessDrop := int(math.Floor(hoursPassed * HappinessDropPerHour))
		if happinessDrop > 0 {
			p.Happiness = max(p.Happiness-happinessDrop, MinStatValue)
		}
	}

	statsChanged := (p.Satiety != oldSatiety) || (p.Happiness != oldHappiness)
	if statsChanged {
		p.LastCalculatedAt = now
	}

	return statsChanged
}
