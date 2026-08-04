package api

import (
	"encoding/json"
	"net/http"

	"github.com/cayman444/avito-gamification-hackathon/blob/main/backend/pets/internal/service"
)

type PetResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Level       int    `json:"level"`
	XP          int    `json:"xp"`
	NextLevelXP int    `json:"next_level_xp"`
	Satiety     int    `json:"satiety"`
	Happiness   int    `json:"happiness"`
}

type CreatePetRequest struct {
	Name string `json:"name"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type PetHandler struct {
	service *service.PetService
}

func NewPetHandler(service *service.PetService) *PetHandler {
	return &PetHandler{
		service: service,
	}
}

func (ph *PetHandler) GetPet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := ctx.Value("user_id").(int64)

	if !ok {
		// TODO
		return
	}

	pet, err := ph.service.GetPet(ctx, userID)
	if err != nil {
		// TODO
	}

	petResponse := PetResponse{
		ID:          pet.ID,
		Name:        pet.Name,
		Level:       pet.Level,
		XP:          pet.XP,
		NextLevelXP: pet.NextLevelXP,
		Satiety:     pet.Satiety,
		Happiness:   pet.Happiness,
	}

	json.NewEncoder(w).Encode(petResponse)
}

func (ph *PetHandler) CreatePet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := ctx.Value("user_id").(int64)
	if !ok {
		// TODO
		return
	}

	var req CreatePetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// TODO
		return
	}

	pet, err := ph.service.CreatePet(ctx, req.Name, userID)
	if err != nil {
		// TODO
	}

	petResponse := PetResponse{
		ID:          pet.ID,
		Name:        pet.Name,
		Level:       pet.Level,
		XP:          pet.XP,
		NextLevelXP: pet.NextLevelXP,
		Satiety:     pet.Satiety,
		Happiness:   pet.Happiness,
	}

	json.NewEncoder(w).Encode(petResponse)
}
