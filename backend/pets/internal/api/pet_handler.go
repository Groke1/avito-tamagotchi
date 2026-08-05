package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/cayman444/avito-gamification-hackathon/blob/main/backend/pets/internal/domain"
	"github.com/cayman444/avito-gamification-hackathon/blob/main/backend/pets/internal/service"

	"github.com/go-playground/validator/v10"
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
	Name string `json:"name" validate:"required,min=2,max=25"`
}

type PetHandler struct {
	service   *service.PetService
	validator *validator.Validate
}

func NewPetHandler(service *service.PetService) *PetHandler {
	return &PetHandler{
		service:   service,
		validator: validator.New(),
	}
}

func (ph *PetHandler) GetPet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := ctx.Value("user_id").(int64)
	if !ok {
		writeError(w, ErrUnauthorized)
		return
	}

	pet, err := ph.service.GetPet(ctx, userID)
	if err != nil {
		writeError(w, ErrPetNotFound)
		return
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
		writeError(w, ErrUnauthorized)
		return
	}

	var req CreatePetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, ErrValidationError)
		return
	} else if err := ph.validator.Struct(req); err != nil {
		writeError(w, ErrValidationError)
		return
	}

	pet, err := ph.service.CreatePet(ctx, req.Name, userID)
	if err != nil {
		writeError(w, ErrPetAlreadyExists)
		return
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

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(petResponse)
}

func (ph *PetHandler) FeedPet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := ctx.Value("user_id").(int64)
	if !ok {
		writeError(w, ErrUnauthorized)
		return
	}

	pet, err := ph.service.FeedPet(ctx, userID)
	if errors.Is(err, domain.ErrPetNotFound) {
		writeError(w, ErrPetNotFound)
		return
	} else if errors.Is(err, domain.ErrForbiddenAction) {
		writeError(w, ErrUnavailableAction)
		return
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

func (ph *PetHandler) StrokePet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := ctx.Value("user_id").(int64)
	if !ok {
		writeError(w, ErrUnauthorized)
		return
	}

	pet, err := ph.service.StrokePet(ctx, userID)
	if errors.Is(err, domain.ErrPetNotFound) {
		writeError(w, ErrPetNotFound)
		return
	} else if errors.Is(err, domain.ErrForbiddenAction) {
		writeError(w, ErrUnavailableAction)
		return
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
