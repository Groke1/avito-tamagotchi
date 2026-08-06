package api

type CreatePetRequest struct {
	Name string `json:"name" validate:"required,min=2,max=25"`
}
<<<<<<< HEAD

type BonusXpRequest struct {
	UserID string `json:"user_id"`
	Streak int    `json:"streak"`
}
=======
>>>>>>> 9f0afb9c68d0604e731ec3d40cd30366c4e2a04f
