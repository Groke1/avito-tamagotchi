package api

type CreatePetRequest struct {
	Name string `json:"name" validate:"required,min=2,max=25"`
}
