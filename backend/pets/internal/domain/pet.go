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
}

const (
	SatietyDropPerHour   = 5.0
	HappinessDropPerHour = 8.0
	MaxStatValue         = 100
	MinStatValue         = 0
)

func (p *Pet) Feed() (bool, error) {
	if p.Satiety >= MaxStatValue {
		return false, ErrPetIsFull
	}

	p.Satiety = min(p.Satiety+5, MaxStatValue)
	levelUp := p.AddXP(2)

	return levelUp, nil
}

func (p *Pet) Stroke() (bool, error) {
	if p.Happiness >= MaxStatValue {
		return false, ErrPetIsTooHappy
	} else if p.Satiety < 20 {
		return false, ErrPetIsTooHungry // TODO frontend integration
	}

	p.Happiness = min(p.Happiness+5, MaxStatValue)
	levelUp := p.AddXP(3)

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

	if p.Satiety < 40 {
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
