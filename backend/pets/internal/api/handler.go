package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/domain"
	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/service"

	"github.com/go-playground/validator/v10"
)

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

	userID, ok := ctx.Value("user_id").(string)
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

	writeJsonResponse(w, http.StatusOK, petResponse)
}

func (ph *PetHandler) CreatePet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := ctx.Value("user_id").(string)
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

	writeJsonResponse(w, http.StatusCreated, petResponse)
}

func (ph *PetHandler) FeedPet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		writeError(w, ErrUnauthorized)
		return
	}

	pet, err := ph.service.FeedPet(ctx, userID)
	if errors.Is(err, domain.ErrPetNotFound) {
		writeError(w, ErrPetNotFound)
		return
	} else if err != nil {
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

	writeJsonResponse(w, http.StatusOK, petResponse)
}

func (ph *PetHandler) StrokePet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		writeError(w, ErrUnauthorized)
		return
	}

	pet, err := ph.service.StrokePet(ctx, userID)
	if errors.Is(err, domain.ErrPetNotFound) {
		writeError(w, ErrPetNotFound)
		return
	} else if errors.Is(err, domain.ErrUnavailableAction) {
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

	writeJsonResponse(w, http.StatusOK, petResponse)
}

func (ph *PetHandler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		writeError(w, ErrUnauthorized)
		return
	}

	limit := 20
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err != nil {
			if l >= 1 && l <= 100 {
				limit = l
			}
		}
	}

	leaderboardItems, currentUser, err := ph.service.GetLeaderboard(ctx, limit, userID)
	if err != nil {
		writeError(w, ErrInternalError)
		return
	}

	leaderboard := make([]LeaderboardItemResponse, len(leaderboardItems))
	for i, item := range leaderboardItems {
		leaderboard[i] = LeaderboardItemResponse{
			Rank:     item.Rank,
			Level:    item.Level,
			UserName: item.UserName,
			PetName:  item.PetName,
		}
	}

	leaderboardResponse := LeaderboardResponse{
		Items: leaderboard,
		CurrentUser: LeaderboardItemResponse{
			Rank:     currentUser.Rank,
			Level:    currentUser.Level,
			UserName: currentUser.UserName,
			PetName:  currentUser.PetName,
		},
	}

	writeJsonResponse(w, http.StatusOK, leaderboardResponse)
}

func (ph *PetHandler) DailyBonus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req BonusXpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, ErrValidationError)
		return
	}

	err := ph.service.ClaimDailyBonus(ctx, req.Streak, req.UserID)
	if err != nil {
		writeError(w, ErrInternalError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (ph *PetHandler) UpdateXP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req UpdateXPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, ErrValidationError)
		return
	}

	_, err := ph.service.GrantXP(ctx, req.XP, req.UserID)
	if err != nil {
		writeError(w, ErrInternalError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
