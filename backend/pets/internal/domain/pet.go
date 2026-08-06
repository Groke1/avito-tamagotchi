package domain

import (
	"time"
)

type Pet struct {
	ID          int64
	UserID      string
	Name        string
	Level       int
	XP          int
	NextLevelXP int
	Satiety     int
	Happiness   int
	CreatedAt   time.Time
}

func (p *Pet) Feed() (bool, error) {
	if p.Satiety >= 100 {
		return false, ErrPetIsFull
	}

	p.Satiety = min(p.Satiety+5, 100)
	p.XP += 2

	if p.XP >= p.NextLevelXP {
		p.Level += 1
		p.NextLevelXP = p.NextLevelXP * p.Level * p.Level
		return true, nil
	}

	return false, nil
}

func (p *Pet) Stroke() (bool, error) {
	if p.Happiness >= 100 {
		return false, ErrPetIsTooHappy
	} else if p.Satiety < 20 {
		return false, ErrPetIsTooHungry // TODO frontend integration
	}

	p.Happiness = min(p.Happiness+5, 100)
	p.XP += 3

	if p.XP >= p.NextLevelXP {
		p.Level += 1
		p.NextLevelXP *= p.Level * p.Level
		return true, nil
	}

	return false, nil
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
