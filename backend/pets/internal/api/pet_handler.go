package api

import (
	"encoding/json"
	"net/http"

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
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Требуется повторная авторизация")
		return
	}

	pet, err := ph.service.GetPet(ctx, userID)
	if err != nil {
		writeError(w, http.StatusNotFound, "PET_NOT_FOUND", "Сначала создайте питомца")
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
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Требуется повторная авторизация")
		return
	}

	var req CreatePetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "Проверьте переданные данные")
		return
	}

	if err := ph.validator.Struct(req); err != nil {
		writeError(w, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "Проверьте переданные данные")
		return
	}

	pet, err := ph.service.CreatePet(ctx, req.Name, userID)
	if err != nil {
		writeError(w, http.StatusConflict, "PET_ALREADY_EXISTS", "У пользователя уже есть питомец")
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
