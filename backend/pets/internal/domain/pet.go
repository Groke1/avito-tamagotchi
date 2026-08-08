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
	LastFeedAt       time.Time
	LastStrokeAt     time.Time
}

const (
	SatietyDropPerHour   = 5.0
	HappinessDropPerHour = 8.0

	MaxStatValue = 100
	MinStatValue = 0

	FeedCooldown   = 3 * time.Minute
	StrokeCooldown = 2 * time.Minute

	FeedSatietyIncrease     = 5
	StrokeHappinessIncrease = 5

	FeedXPAmount   = 2
	StrokeXPAmount = 3

	FeedCost = 5

	HungrySatietyThreshold  = 20
	SadnessSatietyThreshold = 40
)

func (p *Pet) Feed() (bool, int, error) {
	if p.Satiety >= MaxStatValue {
		return false, 0, ErrPetIsFull
	} else if time.Since(p.LastFeedAt) < FeedCooldown {
		return false, 0, &ActionUnavailableError{
			RetryAfter: FeedCooldown - time.Since(p.LastFeedAt),
		}
	}

	p.Satiety = min(p.Satiety+5, MaxStatValue)
	levelUp := p.AddXP(FeedXPAmount)
	p.LastFeedAt = time.Now()

	return levelUp, FeedCost, nil
}

func (p *Pet) Stroke() (bool, error) {
	if p.Happiness >= MaxStatValue {
		return false, ErrPetIsTooHappy
	} else if p.Satiety < HungrySatietyThreshold {
		return false, ErrPetIsTooHungry // TODO frontend integration
	} else if time.Since(p.LastStrokeAt) < StrokeCooldown {
		return false, &ActionUnavailableError{
			RetryAfter: StrokeCooldown - time.Since(p.LastStrokeAt),
		}
	}

	p.Happiness = min(p.Happiness+5, MaxStatValue)
	levelUp := p.AddXP(StrokeXPAmount)
	p.LastStrokeAt = time.Now()

	return levelUp, nil
}

func (p *Pet) AddXP(amount int) bool {
	p.XP += amount
	if p.XP >= p.NextLevelXP {
		p.Level += 1
		p.NextLevelXP *= p.Level * p.Level
		return true
	}

	return false
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
