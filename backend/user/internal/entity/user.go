package entity

type User struct {
	ID            string
	Username      string
	Email         string
	Password      string
	PasswordHash  string
	Coins         uint64
	CurrentStreak uint64
}
