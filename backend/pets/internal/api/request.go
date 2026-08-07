package api

type CreatePetRequest struct {
	Name string `json:"name" validate:"required,min=2,max=25"`
}

type BonusXpRequest struct {
	UserID string `json:"user_id"`
	Streak int    `json:"streak"`
}
